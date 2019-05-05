package sca

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/nanjj/cub/logs"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go/config"
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
	actions *Actions
	rms     *Rms
	tracer  opentracing.Tracer
	closers []io.Closer
}

func NewRunner(cfg *Config) (r *Runner, err error) {
	name, listen, leader := cfg.RunnerName, cfg.RunnerListen, cfg.LeaderListen
	tracer, closer, err := logs.NewTracer(name,
		config.Tag("runner", listen),
		config.Tag("leader", leader))
	if err != nil {
		panic(err)
	}
	ctx := context.Background()
	sp, ctx := logs.StartSpanFromContextWithTracer(ctx, tracer, "NewRunner")
	defer sp.Finish()
	r = &Runner{
		name:    name,
		listen:  listen,
		actions: &Actions{},
		rms:     &Rms{Name: name},
		closers: []io.Closer{closer},
		tracer:  tracer,
	}
	r.closers = append(r.closers, closer)
	r.AddAction("join", r.rms.Join)
	r.AddAction("route", r.Route)
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
		if err = SendEvent(ctx, sock,
			&Event{
				Action: "join",
				Payload: Payload{
					DataObject(name), DataObject(listen), DataObject(name)},
			}); err != nil {
			sp.Error(err.Error())
			return
		}
		if err = SendEvent(ctx, sock,
			&Event{
				Action: "route",
				Payload: Payload{
					DataObject(name), DataObject(name)},
			}); err != nil {
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
	// local := false
	targets := e.Receiver
	local, ups, vias := r.rms.Dispatch(targets)
	if len(ups) != 0 {
		if leader := r.leader; leader != nil {
			dup := e.Clone()
			dup.Receiver = ups
			if err = SendEvent(ctx, leader, dup); err != nil {
				sp.Error("Failed to send", zap.Stack("stack"), zap.Error(err))
				return
			}
		} else {
			local = true
		}
	}

	// forward to members
	for via, downs := range vias {
		if member, ok := r.rms.GetMember(via); ok {
			dup := e.Clone()
			dup.Receiver = downs
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
	t := DataObject{}
	err = t.Encode(time.Now().UTC())
	if err == nil {
		rep = Payload{DataObject(r.Name()), t}
	}
	return
}

func (r *Runner) Route(ctx context.Context, req Payload) (rep Payload, err error) {
	sp, ctx := logs.StartSpanFromContext(ctx, "Route")
	defer sp.Finish()
	name := string(req[0])
	l := len(req)
	// update routes
	for i := 1; i < l; i++ {
		target := string(req[i])
		r.rms.AddRoute(target, name)
	}
	// tell leader
	if leader := r.leader; leader != nil {
		req[0] = DataObject(r.name)
		if err = SendEvent(ctx, leader, &Event{
			Action:  "route",
			Payload: req,
		}); err != nil {
			sp.Error(err.Error())
			return
		}
	}
	return
}

func (r *Runner) Members() (members []string) {
	members = r.rms.Members()
	return
}

func (r *Runner) Routes() (routes map[string]string) {
	routes = r.rms.Routes()
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
	if sock, ok := r.rms.GetMember(name); ok {
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
