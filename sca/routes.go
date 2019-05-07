package sca

import (
	"sync"

	"nanomsg.org/go/mangos/v2"
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

func (r *Rms) Targets(via string) (names []string) {
	f := func(k, v interface{}) bool {
		if name, ok := k.(string); ok {
			if value, ok := v.(string); ok {
				if value == via {
					names = append(names, name)
				}
			}
		}
		return true
	}
	r.r.Range(f)
	return
}

func (r *Rms) DelRoute(keys ...string) {
	for i := range keys {
		r.r.Delete(keys[i])
	}
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

func (r *Rms) DelMember(name string) {
	r.m.Delete(name)
}

func (r *Rms) GetMember(name string) (sock mangos.Socket, ok bool) {
	v, ok := r.m.Load(name)
	if ok {
		sock, ok = v.(mangos.Socket)
	}
	return
}

func (r *Rms) Members() (names Set) {
	names = NewSet()
	r.m.Range(func(k, v interface{}) bool {
		if name, ok := k.(string); ok {
			names.Add(name)
		}
		return true
	})
	return
}

func (r *Rms) HasMember() (ok bool) {
	r.m.Range(func(k, v interface{}) bool {
		ok = true
		return false
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
		for member := range r.Members() {
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
