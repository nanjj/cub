package sca

import (
	"sync"

	"github.com/nanjj/cub/sdo"
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

func (r *Rms) AddMember(name string, node *Node) {
	r.m.Store(name, node)
}

func (r *Rms) DelMember(name string) {
	r.m.Delete(name)
}

func (r *Rms) GetMember(name string) (node *Node, ok bool) {
	v, ok := r.m.Load(name)
	if ok {
		node, ok = v.(*Node)
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

func (r *Rms) Dispatch(targets sdo.Targets) (local bool, ups sdo.Targets, vias map[string]sdo.Targets) {
	vias = map[string]sdo.Targets{}
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
