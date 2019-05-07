package sca

import (
	"context"
	"fmt"
	"time"

	"github.com/nanjj/cub/logs"
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
	tracer  *logs.Tracer
	g       errgroup.Group
}

func NewRunner(cfg *Config) (r *Runner, err error) {
	name, listen, leader := cfg.RunnerName, cfg.RunnerListen, cfg.LeaderListen
	tracer, err := logs.NewTracer(name,
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
		tracer:  tracer,
	}
	r.AddAction("login", r.Login)
	r.AddAction("logout", r.Logout)
	r.AddAction("route", r.Route)
	r.AddAction("derour", r.Derour)
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
		// Login
		if err = SendEvent(ctx, sock,
			&Event{
				Action: "login",
				Payload: Payload{
					DataObject(name), DataObject(listen)},
				From: name,
			}); err != nil {
			sp.Error("failed to login", zap.Stack("stack"), zap.Error(err))
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
		sp, ctx := logs.StartSpanFromCarrier(e.Carrier, r.tracer, "Recv")
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

func (r *Runner) Derour(ctx context.Context, req Payload) (rep Payload, err error) {
	sp, ctx := logs.StartSpanFromContext(ctx, "Derour")
	defer sp.Finish()
	for i := range req {
		r.rms.DelRoute(string(req[i]))
	}
	// tell leader
	if leader := r.leader; leader != nil {
		if err = SendEvent(ctx, leader.Client, &Event{
			Action:  "derour",
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
	//Use global tracer
	sp, ctx := logs.StartSpanFromContextWithTracer(context.Background(), r.tracer, "Close")
	defer r.tracer.Close()
	defer sp.Finish()
	if leader := r.leader; leader != nil {
		name := r.Name()
		sock := leader.Client
		defer sock.Close()
		if err = SendEvent(ctx, sock,
			&Event{
				Action:  "logout",
				Payload: Payload{DataObject(name)},
				From:    name,
			}); err != nil {
			sp.Error("failed to logout", zap.Stack("stack"), zap.Error(err))
			return
		}
	}
	for m := range r.Members() {
		if sock := r.Member(m); sock != nil {
			defer sock.Close()
		}
	}
	defer time.Sleep(time.Millisecond * 10)
	return
}

func (r *Runner) Name() string { return r.self.Name }

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

// Login name listen
func (r *Runner) Login(ctx context.Context, req Payload) (rep Payload, err error) {
	sp, ctx := logs.StartSpanFromContext(ctx, "Login")
	defer sp.Finish()
	if len(req) != 2 {
		sp.Error("Bad request")
		err = fmt.Errorf("bad request")
		return
	}
	name := string(req[0])
	listen := string(req[1])
	// add new member
	sock, ok := r.rms.GetMember(name)
	if ok && listen != "" {
		sock.Close()
		sock = nil
	}
	if sock == nil {
		if sock, err = NewClient(listen); err != nil {
			sp.Error(err.Error())
			return
		}
		r.rms.AddMember(name, sock)
	}
	req[1] = DataObject(name)
	r.Route(ctx, req)
	return
}

// Logout name
func (r *Runner) Logout(ctx context.Context, req Payload) (rep Payload, err error) {
	sp, ctx := logs.StartSpanFromContext(ctx, "Logout")
	defer sp.Finish()
	if len(req) != 1 {
		sp.Error("Bad request")
		err = fmt.Errorf("bad request")
		return
	}
	name := string(req[0])
	targets := r.rms.Targets(name)
	req = make([]DataObject, len(targets))
	for i := range targets {
		req[i] = DataObject(targets[i])
	}
	r.Derour(ctx, req)

	if sock, ok := r.rms.GetMember(name); ok {
		r.rms.DelMember(name)
		if sock != nil {
			sock.Close()
		}
	}
	return
}
