package sca

import (
	"context"
	"fmt"
	"time"

	"github.com/nanjj/cub/logs"
	"github.com/uber/jaeger-client-go/config"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	_ "nanomsg.org/go/mangos/v2/transport/all"
)

const (
	ROOT   = "root"
	LEAF   = "leaf"
	BRANCH = "branch"
	EXIT   = ""
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
	r.g.Go(r.run)
	if leader != "" {
		r.leader = &Node{Listen: leader}
		// Login
		if err = r.leader.Login(ctx, name, listen); err != nil {
			sp.Error("failed to login", zap.Stack("stack"), zap.Error(err))
			return
		}
	}
	return
}

func (r *Runner) run() (err error) {
	for {
		e := &Event{}
		if err = r.self.Recv(context.Background(), e); err != nil {
			continue
		}
		if e.Action == EXIT {
			break
		}
		r.g.Go(func() (err error) {
			sp, ctx := logs.StartSpanFromCarrier(e.Carrier, r.tracer, "Run")
			defer sp.Finish()
			if err = r.Handle(ctx, e); err != nil {
				sp.Error("Failed to handle event", zap.Error(err))
				return
			}
			return
		})
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
			if err = r.leader.Send(ctx, dup); err != nil {
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
			if err = member.Send(ctx, dup); err != nil {
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
		if err = r.leader.Send(ctx, &Event{
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
		if err = r.leader.Send(ctx, &Event{
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

func (r *Runner) Leader() (leader *Node) {
	return r.leader
}

func (r *Runner) Self() *Node {
	return r.self
}

func (r *Runner) Member(name string) *Node {
	if member, ok := r.rms.GetMember(name); ok {
		return member
	}
	return nil
}

func (r *Runner) Close() (err error) {
	//Use global tracer
	sp, ctx := logs.StartSpanFromContextWithTracer(context.Background(), r.tracer, "Close")
	defer r.tracer.Close()
	defer sp.Finish()
	defer r.self.Close()
	if leader := r.leader; leader != nil {
		defer leader.Close()
		name := r.Name()
		if err = leader.Logout(ctx, name); err != nil {
			sp.Error("failed to logout", zap.Stack("stack"), zap.Error(err))
			return
		}
	}
	for name := range r.Members() {
		if node := r.Member(name); node != nil {
			defer node.Close()
		}
	}
	r.self.Exit(ctx)
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
	member, ok := r.rms.GetMember(name)
	if ok && listen != "" {
		member.Close()
		member = nil
	}
	if member == nil {
		member = &Node{Name: name, Listen: listen}
		r.rms.AddMember(name, member)
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

	if node, ok := r.rms.GetMember(name); ok {
		r.rms.DelMember(name)
		if node != nil {
			node.Close()
		}
	}
	return
}
