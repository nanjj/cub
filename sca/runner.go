package sca

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/nanjj/cub/logs"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
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
		panic(err)
	}
	ctx := context.Background()
	sp, ctx := logs.StartSpanFromContextWithTracer(ctx, tracer, "NewRunner")
	r = &Runner{
		name:    name,
		listen:  listen,
		members: &Members{},
		actions: &Actions{},
		routes:  &Routes{},
		closers: []io.Closer{closer},
		tracer:  tracer,
	}
	r.closers = append(r.closers, closer)
	r.AddAction("join", r.Join)
	r.AddAction("ping", r.Ping)
	sock, err := pull.NewSocket()
	if err != nil {
		return
	}
	if err = RetryListen(sock, listen); err != nil {
		sp.Error("Failed to listen", zap.Stack("stack"), zap.Error(err))
		return
	}
	r.self = sock
	if leader != "" {
		if sock, err = push.NewSocket(); err != nil {
			sp.Fatal(err.Error())
			return
		}
		if err = RetryDial(sock, leader); err != nil {
			sp.Fatal(err.Error())
			return
		}
		e := &Event{
			Action:  "join",
			Payload: Payload{DataObject(name), DataObject(listen), DataObject(name)},
		}
		if err = SendEvent(ctx, sock, e); err != nil {
			sp.Error(err.Error())
			return
		}
		r.leader = sock
	}
	return
}

func (r *Runner) Run() (err error) {
	for {
		e := &Event{}
		if err = RecvEvent(context.Background(), r.self, e); err != nil {
			continue
		}
		if err = r.Handle(e); err != nil {
			continue
		}
	}
	return
}

func (r *Runner) Handle(e *Event) (err error) {
	tracer := r.Tracer()
	sp, ctx := logs.StartSpanFromCarrier(e.Carrier, tracer, "Recv")
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
					if err = SendEvent(ctx, leader, dup); err != nil {
						sp.Error("Failed to send", zap.Stack("stack"), zap.Error(err))
						return
					}
				}
			}
			continue
		}
		if member, ok := r.members.Get(k); ok {
			dup := e.Dup()
			dup.Receiver = v
			if err = SendEvent(ctx, member, dup); err != nil {
				sp.Error("Failed to send", zap.Stack("stack"), zap.Error(err))
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
				sp.Error("Failed to run action", zap.Stack("stack"), zap.Error(err))
			}
		} else {
			err = fmt.Errorf("No action found")
			sp.Error("Failed to find action", zap.Stack("stack"), zap.Error(err))
		}
		callback := e.Callback
		if callback != "" && err == nil {
			ack := Event{Action: callback}
			ack.Receiver = e.Sender
			ack.Payload = rep
			if err = SendEvent(ctx, r.Leader(), &ack); err != nil {
				sp.Error("Failed to send leader", zap.Stack("stack"), zap.Error(err))
			}
		}
	}
	return
}

// Ping
func (r *Runner) Ping(ctx context.Context, req Payload) (rep Payload, err error) {
	sp, ctx := logs.StartSpanFromContext(ctx, "Ping")
	defer sp.Finish()
	sp.Info("Ping", zap.String("name", r.Name()))
	return
}

// name,listen, members...
func (r *Runner) Join(ctx context.Context, req Payload) (rep Payload, err error) {
	sp, ctx := logs.StartSpanFromContext(ctx, "Join")
	defer sp.Finish()
	l := len(req)
	if l < 3 {
		err = fmt.Errorf("bad request")
		sp.Error(err.Error())
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
			sp.Error(err.Error())
			return
		}
		if err = RetryDial(sock, listen); err != nil {
			sp.Error(err.Error())
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
		if err = SendEvent(ctx, leader, e); err != nil {
			sp.Error(err.Error())
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
