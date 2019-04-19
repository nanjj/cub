package sca

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/opentracing/opentracing-go"
	"nanomsg.org/go/mangos/v2"
	"nanomsg.org/go/mangos/v2/protocol/pull"
	"nanomsg.org/go/mangos/v2/protocol/push"
	_ "nanomsg.org/go/mangos/v2/transport/all"
)

type Runner struct {
	name    string
	listen  string
	leader  mangos.Socket
	self    mangos.Socket
	members *Members
	actions *Actions
	routes  *Routes
	tracer  opentracing.Tracer
	closers []io.Closer
}

func NewRunner(cfg *Config) (r *Runner, err error) {
	name, listen, leader := cfg.RunnerName, cfg.RunnerListen, cfg.LeaderListen
	tracer, closer, err := NewTracer(name, listen, leader)
	if err != nil {
		log.Fatal(err)
		return
	}
	ctx := context.Background()
	sp := tracer.StartSpan("NewRunner")
	defer sp.Finish()
	r = &Runner{
		name:    name,
		listen:  listen,
		members: &Members{},
		actions: &Actions{},
		routes:  &Routes{},
		closers: []io.Closer{closer},
		tracer:  tracer,
	}
	ctx = opentracing.ContextWithSpan(ctx, sp)
	r.closers = append(r.closers, closer)
	r.AddAction("join", r.Join)
	r.AddAction("ping", r.Ping)
	sock, err := pull.NewSocket()
	if err != nil {
		return
	}
	if err = RetryListen(sock, listen); err != nil {
		sp.LogKV("error", err)
		log.Fatal(err)
		return
	}
	r.self = sock
	if leader != "" {
		if sock, err = push.NewSocket(); err != nil {
			sp.LogKV("error", err)
			log.Fatal(err)
			return
		}
		if err = RetryDial(sock, leader); err != nil {
			sp.LogKV("error", err)
			log.Fatal(err)
			return
		}
		e := &Event{
			Action:  "join",
			Payload: Payload{DataObject(name), DataObject(listen), DataObject(name)},
		}
		if err = Send(ctx, sock, e); err != nil {
			sp.SetTag("error", true)
			sp.LogKV("error", err)
			return
		}
		r.leader = sock
	}
	return
}

func (r *Runner) Run() (err error) {
	for {
		e := &Event{}
		if err = e.Recv(r.self); err != nil {
			log.Println(err)
			continue
		}
		if err = r.Handle(e); err != nil {
			log.Println(err)
			continue
		}
	}
	return
}

func (r *Runner) Handle(e *Event) (err error) {
	tracer := r.Tracer()
	sp, ctx := e.SpanContext(tracer, "Recv")
	defer sp.Finish()
	local := false
	targets := e.Receiver
	vias := map[string]Targets{}
	if targets.Local() {
		local = true
	} else if targets.All() {
		local = true
		for _, k := range r.members.Names() {
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
				var tgts Targets
				for _, tgt := range v {
					if string(tgt) != r.Name() {
						tgts = append(tgts, tgt)
					} else {
						local = true
					}
				}
				if len(tgts) != 0 {
					dup := e.Dup()
					dup.Receiver = tgts
					if err = Send(ctx, leader, dup); err != nil {
						log.Println(err)
						return
					}
				}
			}
			continue
		}
		if member, ok := r.members.Get(k); ok {
			dup := e.Dup()
			dup.Receiver = v
			if err = Send(ctx, member, dup); err != nil {
				log.Println(err)
				return
			}
		} else {
			return
		}
	}
	if local { // handle local
		action := e.Action
		req := e.Payload
		var rep Payload
		if f, ok := r.actions.Get(action); ok {
			if rep, err = f(ctx, req); err != nil {
				log.Println(err)
			}
		} else {
			err = fmt.Errorf("No action found")
			log.Println(err)
		}
		callback := e.Callback
		if callback != "" && err == nil {
			ack := Event{Action: callback}
			ack.Receiver = e.Sender
			ack.Payload = rep
			if err = Send(ctx, r.Leader(), &ack); err != nil {
				log.Println(err)
			}
		}
	}
	return
}

// Ping
func (r *Runner) Ping(ctx context.Context, req Payload) (rep Payload, err error) {
	sp, ctx := opentracing.StartSpanFromContextWithTracer(ctx, r.tracer, fmt.Sprintf("Ping@%s", r.Name()))
	defer sp.Finish()
	log.Println("ping", r.Name(), req)
	return
}

// name,listen, members...
func (r *Runner) Join(ctx context.Context, req Payload) (rep Payload, err error) {
	sp, ctx := StartSpanFromContext(ctx, "Join")
	defer sp.Finish()
	l := len(req)
	if l < 3 {
		err = fmt.Errorf("bad request")
		sp.LogKV("error", err)
		return
	}
	name := string(req[0])
	listen := string(req[1])
	// add new member
	sock, ok := r.members.Get(name)
	if ok && listen != "" {
		sock.Close()
		sock = nil
	}
	if sock == nil {
		if sock, err = push.NewSocket(); err != nil {
			sp.LogKV("error", err)
			return
		}
		if err = RetryDial(sock, listen); err != nil {
			sp.LogKV("error", err)
			return
		}
		r.members.Add(name, sock)
	}

	// update routes
	for i := 2; i < l; i++ {
		target := string(req[i])
		r.routes.Add(target, name)
	}
	// tell leader
	if leader := r.leader; leader != nil {
		req[0] = DataObject(r.name)
		req[1] = DataObject("")
		e := &Event{
			Action:  "join",
			Payload: req,
		}
		if err = Send(ctx, leader, e); err != nil {
			sp.LogKV("error", err)
			return
		}
	}
	return
}

func (r *Runner) Members() (members []string) {
	members = r.members.Names()
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

func (r *Runner) AddAction(name string, action Action) {
	r.actions.Add(name, action)
}

func (r *Runner) Leader() mangos.Socket {
	return r.leader
}

func (r *Runner) Self() mangos.Socket {
	return r.self
}

func (r *Runner) Member(name string) mangos.Socket {
	if sock, ok := r.members.Get(name); ok {
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

func (r *Runner) Tracer() opentracing.Tracer { return r.tracer }
