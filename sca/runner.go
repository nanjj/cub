package sca

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/nanjj/cub/logs"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go/config"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"nanomsg.org/go/mangos/v2"
	_ "nanomsg.org/go/mangos/v2/transport/all"
)

const (
	ROOT   = "root"
	LEAF   = "leaf"
	BRANCH = "branch"
)

type Runner struct {
	leader  *Node
	self    *Node
	actions *Actions
	rms     *Rms
	tracer  opentracing.Tracer
	closers []io.Closer
	g       errgroup.Group
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
		self:    &Node{Name: name, Listen: listen},
		actions: &Actions{},
		rms:     &Rms{Name: name},
		closers: []io.Closer{closer},
		tracer:  tracer,
	}
	r.closers = append(r.closers, closer)
	r.AddAction("join", r.rms.Join)
	r.AddAction("route", r.Route)
	r.AddAction("ping", r.Ping)
	sock, err := NewServer(listen)
	if err != nil {
		sp.Error("Failed to listen", zap.Stack("stack"), zap.Error(err))
		return
	}
	r.self.Server = sock
	r.g.Go(r.run)
	if leader != "" {
		if sock, err = NewClient(leader); err != nil {
			sp.Error("Failed to dial leader", zap.Stack("stack"), zap.Error(err))
			return
		}
		// Join
		if err = SendEvent(ctx, sock,
			&Event{
				Action: "join",
				Payload: Payload{
					DataObject(name), DataObject(listen)},
				From: name,
			}); err != nil {
			sp.Error("failed to join", zap.Stack("stack"), zap.Error(err))
			return
		}
		//Add route
		if err = SendEvent(ctx, sock,
			&Event{
				Action: "route",
				Payload: Payload{
					DataObject(name), DataObject(name)},
				From: name,
			}); err != nil {
			sp.Error(err.Error())
			return
		}
		r.leader = &Node{Listen: listen, Client: sock}
	}
	return
}

func (r *Runner) run() (err error) {
	for {
		e := &Event{}
		if err = RecvEvent(context.Background(), r.self.Server, e); err != nil {
			continue
		}
		sp, ctx := logs.StartSpanFromCarrier(e.Carrier, r.Tracer(), "Recv")
		if err = r.Handle(ctx, e); err != nil {
			sp.Finish()
			continue
		}
		sp.Finish()
	}
	return
}

func (r *Runner) Handle(ctx context.Context, e *Event) (err error) {
	sp, ctx := logs.StartSpanFromContext(ctx, "Handle")
	defer sp.Finish()
	targets := e.To
	local, ups, vias := r.rms.Dispatch(targets)
	if l := len(ups); l != 0 {
		if leader := r.leader; leader != nil {
			dup := e.Clone()
			dup.To = ups
			if err = SendEvent(ctx, leader.Client, dup); err != nil {
				sp.Error("Failed to send", zap.Stack("stack"), zap.Error(err))
				return
			}
		} else {
			if l == 1 && ups[0] == "" {
				local = true
			} else {
				err = fmt.Errorf("targets %v not found", ups)
				sp.Error("targets not found", zap.Stack("stack"), zap.Error(err))
				return
			}
		}
	}

	// forward to members
	for via, downs := range vias {
		if member, ok := r.rms.GetMember(via); ok {
			dup := e.Clone()
			dup.To = downs
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
			ack := &Event{Carrier: e.Carrier, Action: callback}
			ack.To = Targets{e.From}
			ack.From = r.Name()
			ack.Payload = rep
			if err = r.Handle(ctx, ack); err != nil {
				sp.Error("Failed to send leader", zap.Stack("stack"), zap.Error(err))
			}
		}
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
	name = r.self.Name
	rep = Payload{
		DataObject(name),
		DataObject(r.NodeType()),
	}
	// tell leader
	if leader := r.leader; leader != nil {
		req[0] = DataObject(name)
		if err = SendEvent(ctx, leader.Client, &Event{
			Action:  "route",
			Payload: req,
		}); err != nil {
			sp.Error(err.Error())
			return
		}
	}
	return
}

func (r *Runner) Members() (members Set) {
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

func (r *Runner) AddCallback(action Action) (name string) {
	return r.actions.New(action)
}

func (r *Runner) RemoveAction(name string) {
	r.actions.Delete(name)
}

func (r *Runner) Leader() (leader mangos.Socket) {
	if r.leader != nil {
		leader = r.leader.Client
	}
	return
}

func (r *Runner) Self() mangos.Socket {
	return r.self.Server
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
	for m := range r.Members() {
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

func (r *Runner) Name() string { return r.self.Name }

func (r *Runner) Tracer() opentracing.Tracer { return r.tracer }

func (r *Runner) Wait() error {
	return r.g.Wait()
}

func (r *Runner) NodeType() string {
	if r.leader == nil {
		return ROOT
	} else if r.rms.HasMember() {
		return BRANCH
	} else {
		return LEAF
	}
}
