package main

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/nanjj/cub/tasks"
	"nanomsg.org/go/mangos/v2"
	"nanomsg.org/go/mangos/v2/protocol/pull"
	"nanomsg.org/go/mangos/v2/protocol/push"
	_ "nanomsg.org/go/mangos/v2/transport/all"
)

const (
	RunnerListen = "runner.listen"
	RunnerName   = "runner.name"
	RunnerIp     = "runner.ip"
	LeaderListen = "leader.listen"
)

type Runner struct {
	name     string
	listen   string
	leader   mangos.Socket
	self     mangos.Socket
	members  map[string]mangos.Socket
	handlers map[string]tasks.Handler
	routes   *Routes
}

func NewRunner(name, listen, leader string) (runner *Runner, err error) {
	runner = &Runner{
		name:     name,
		listen:   listen,
		members:  map[string]mangos.Socket{},
		handlers: map[string]tasks.Handler{},
		routes:   &Routes{},
	}
	runner.AddHandler("join", runner.Join)
	runner.AddHandler("ping", runner.Ping)
	sock, err := pull.NewSocket()
	if err != nil {
		return
	}

	if err = tasks.RetryListen(sock, listen); err != nil {
		return
	}
	runner.self = sock
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
		runner.leader = sock
	}
	return
}

func (runner *Runner) Run() (err error) {
	for {
		if err = runner.Handle(); err != nil {
			return
		}
	}
	return
}

func (runner *Runner) Handle() (err error) {
	task := &tasks.Task{}
	if err = task.Recv(runner.self); err != nil {
		return
	}
	local := false
	targets := task.Targets
	vias := map[string]tasks.Targets{}
	if targets.Local() {
		local = true
	} else if targets.All() {
		local = true
		for k, _ := range runner.members {
			vias[k] = targets
		}
	} else {
		vias = runner.routes.Dispatch(targets)
	}
	for k, v := range vias {
		if k == "" {
			if leader := runner.Leader(); leader == nil {
				local = true
				continue
			} else {
				var tgts tasks.Targets
				for _, tgt := range v {
					if string(tgt) != runner.Name() {
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
		if member, ok := runner.members[k]; ok {
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
		if f, ok := runner.handlers[name]; ok {
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
			if err = ack.Send(runner.Leader()); err != nil {
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
func (runner *Runner) Join(args []tasks.Arg) (ack []tasks.Arg, err error) {
	l := len(args)
	if l < 3 {
		err = fmt.Errorf("bad request")
		return
	}
	name := string(args[0])
	listen := string(args[1])
	// add new member
	sock, ok := runner.members[name]
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
		runner.members[name] = sock
	}

	// update routes
	for i := 2; i < l; i++ {
		target := string(args[i])
		runner.routes.Add(target, name)
	}
	// tell leader
	if leader := runner.leader; leader != nil {
		args[0] = tasks.Arg(runner.name)
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
	kill := func(sock mangos.Socket) {
		if sock != nil {
			if err := sock.Close(); err != nil {
				errs = append(errs, err.Error())
			}
		}
	}
	kill(r.self)
	kill(r.leader)
	for _, sock := range r.members {
		kill(sock)
	}
	if errs != nil {
		err = errors.New(strings.Join(errs, "\n"))
	}
	return
}

func (r *Runner) Name() string { return r.name }
