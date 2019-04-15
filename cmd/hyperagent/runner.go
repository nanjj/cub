package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/nanjj/cub/tasks"
	"github.com/opentracing/opentracing-go"
	"nanomsg.org/go/mangos/v2"
	"nanomsg.org/go/mangos/v2/protocol/pull"
	"nanomsg.org/go/mangos/v2/protocol/push"
	_ "nanomsg.org/go/mangos/v2/transport/all"
)

type Runner struct {
	name     string
	listen   string
	leader   mangos.Socket
	self     mangos.Socket
	members  map[string]mangos.Socket
	handlers map[string]tasks.Handler
	routes   *Routes
	tracer   opentracing.Tracer
	closers  []io.Closer
}

func NewRunner(cfg *Config) (r *Runner, err error) {
	name, listen, leader := cfg.RunnerName, cfg.RunnerListen, cfg.LeaderListen
	r = &Runner{
		name:     name,
		listen:   listen,
		members:  map[string]mangos.Socket{},
		handlers: map[string]tasks.Handler{},
		routes:   &Routes{},
		closers:  []io.Closer{},
	}
	if jaeger := cfg.Jaeger; jaeger != nil {
		var closer io.Closer
		if r.tracer, closer, err = cfg.Jaeger.NewTracer(); err != nil {
			return
		}
		r.closers = append(r.closers, closer)
	}
	r.AddHandler("join", r.Join)
	r.AddHandler("ping", r.Ping)
	sock, err := pull.NewSocket()
	if err != nil {
		return
	}
	if err = tasks.RetryListen(sock, listen); err != nil {
		return
	}
	r.self = sock
	if leader != "" {
		if sock, err = push.NewSocket(); err != nil {
			return
		}
		if err = tasks.RetryDial(sock, leader); err != nil {
			return
		}
		task := &tasks.Task{
			Name: "join",
			Args: []tasks.Arg{tasks.Arg(name), tasks.Arg(listen), tasks.Arg(name)},
		}
		if err = task.Send(sock); err != nil {
			return
		}
		r.leader = sock
	}
	return
}

func (r *Runner) Run() (err error) {
	for {
		if err = r.Handle(); err != nil {
			return
		}
	}
	return
}

func (r *Runner) Handle() (err error) {
	task := &tasks.Task{}
	if err = task.Recv(r.self); err != nil {
		return
	}
	local := false
	targets := task.Targets
	vias := map[string]tasks.Targets{}
	if targets.Local() {
		local = true
	} else if targets.All() {
		local = true
		for k, _ := range r.members {
			vias[k] = targets
		}
	} else {
		vias = r.routes.Dispatch(targets)
	}
	for k, v := range vias {
		if k == "" {
			if leader := r.Leader(); leader == nil {
				local = true
				continue
			} else {
				var tgts tasks.Targets
				for _, tgt := range v {
					if string(tgt) != r.Name() {
						tgts = append(tgts, tgt)
					} else {
						local = true
					}
				}
				if len(tgts) != 0 {
					dup := task.Dup()
					dup.Targets = tgts
					if err = dup.Send(leader); err != nil {
						log.Println(err)
						return
					}
				}
			}
			continue
		}
		if member, ok := r.members[k]; ok {
			dup := task.Dup()
			dup.Targets = v
			if err = dup.Send(member); err != nil {
				log.Println(err)
				return
			}
		} else {
			return
		}
	}
	if local { // handle local
		name := task.Name
		args := task.Args
		var rep []tasks.Arg
		if f, ok := r.handlers[name]; ok {
			if rep, err = f(args); err != nil {
				log.Println(err)
			}
		} else {
			err = fmt.Errorf("No handler found")
			log.Println(err)
		}
		reply := task.Reply
		if reply != "" && err == nil {
			ack := tasks.Task{
				Name:    reply,
				Targets: task.Ack,
				Args:    rep,
			}
			if err = ack.Send(r.Leader()); err != nil {
				log.Println(err)
			}
		}
	}
	return
}

// Ping
func (r *Runner) Ping(args []tasks.Arg) (ack []tasks.Arg, err error) {
	log.Println("ping", r.Name(), args)
	return
}

// name,listen, members...
func (r *Runner) Join(args []tasks.Arg) (ack []tasks.Arg, err error) {
	l := len(args)
	if l < 3 {
		err = fmt.Errorf("bad request")
		return
	}
	name := string(args[0])
	listen := string(args[1])
	// add new member
	sock, ok := r.members[name]
	if ok && listen != "" {
		sock.Close()
		sock = nil
	}
	if sock == nil {
		if sock, err = push.NewSocket(); err != nil {
			return
		}
		if err = tasks.RetryDial(sock, listen); err != nil {
			return
		}
		r.members[name] = sock
	}

	// update routes
	for i := 2; i < l; i++ {
		target := string(args[i])
		r.routes.Add(target, name)
	}
	// tell leader
	if leader := r.leader; leader != nil {
		args[0] = tasks.Arg(r.name)
		args[1] = tasks.Arg("")
		task := &tasks.Task{
			Name: "join",
			Args: args,
		}
		if err = task.Send(leader); err != nil {
			return
		}
	}
	return
}

func (r *Runner) Members() (members []string) {
	members = []string{}
	for k, _ := range r.members {
		members = append(members, k)
	}
	return
}

func (r *Runner) Handlers() (handlers []string) {
	for k, _ := range r.handlers {
		handlers = append(handlers, k)
	}
	return
}

func (r *Runner) Routes() (routes map[string]string) {
	routes = map[string]string{}
	f := func(k, v interface{}) bool {
		if key, ok := k.(string); ok {
			if value, ok := v.(string); ok {
				routes[key] = value
			}
		}
		return true
	}
	r.routes.Range(f)
	return
}

func (r *Runner) AddHandler(name string, handler tasks.Handler) {
	r.handlers[name] = handler
}

func (r *Runner) Leader() mangos.Socket {
	return r.leader
}

func (r *Runner) Self() mangos.Socket {
	return r.self
}

func (r *Runner) Member(name string) mangos.Socket {
	if sock, ok := r.members[name]; ok {
		return sock
	}
	return nil
}

func (r *Runner) Close() (err error) {
	var errs []string
	closers := []io.Closer{
		r.Self(),
		r.Leader(),
	}
	for _, m := range r.Members() {
		closers = append(closers, r.Member(m))
	}
	closers = append(closers, r.closers...)
	for i := 0; i < len(closers); i++ {
		c := closers[i]
		if c != nil {
			if err := c.Close(); err != nil {
				errs = append(errs, err.Error())
			}
		}
	}
	if errs != nil {
		err = errors.New(strings.Join(errs, "\n"))
	}
	return
}

func (r *Runner) Name() string { return r.name }
