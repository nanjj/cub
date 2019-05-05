package sca

import (
	"context"
	"fmt"
	"sync"

	"github.com/nanjj/cub/logs"
	"nanomsg.org/go/mangos/v2"
	"nanomsg.org/go/mangos/v2/protocol/push"
)

type Rms struct {
	Name string
	r    sync.Map
	m    sync.Map
}

func (r *Rms) AddRoute(target, via string) {
	r.r.Store(target, via)
}

func (r *Rms) GetRoute(target string) (via string) {
	if v, ok := r.r.Load(target); ok {
		if via, ok = v.(string); ok {
			return
		}
	}
	return
}

func (r *Rms) DelRoute(target string) {
	r.r.Delete(target)
}

func (r *Rms) Routes() (routes map[string]string) {
	routes = map[string]string{}
	f := func(k, v interface{}) bool {
		if key, ok := k.(string); ok {
			if value, ok := v.(string); ok {
				routes[key] = value
			}
		}
		return true
	}
	r.r.Range(f)
	return
}

func (r *Rms) AddMember(name string, sock mangos.Socket) {
	r.m.Store(name, sock)
}

func (r *Rms) GetMember(name string) (sock mangos.Socket, ok bool) {
	v, ok := r.m.Load(name)
	if ok {
		sock, ok = v.(mangos.Socket)
	}
	return
}

func (r *Rms) Members() (names []string) {
	names = []string{}
	r.m.Range(func(k, v interface{}) bool {
		if name, ok := k.(string); ok {
			names = append(names, name)
		}
		return true
	})
	return
}

func (r *Rms) Dispatch(targets Targets) (local bool, ups Targets, vias map[string]Targets) {
	vias = map[string]Targets{}
	if targets.Local() {
		local = true
		return
	}

	if targets.Down() {
		local = true
		for _, member := range r.Members() {
			vias[member] = targets
		}
		return
	}

	for _, target := range targets {
		via := r.GetRoute(target)
		if via == "" {
			if target == r.Name {
				local = true
			} else {
				ups = append(ups, target)
			}
		} else {
			if ss, ok := vias[via]; ok {
				ss = append(ss, target)
				vias[via] = ss
			} else {
				vias[via] = []string{target}
			}
		}
	}
	return
}

// name,listen, members...
func (r *Rms) Join(ctx context.Context, req Payload) (rep Payload, err error) {
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
	sock, ok := r.GetMember(name)
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
		r.AddMember(name, sock)
	}
	return
}
