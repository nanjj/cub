package main

import (
	"sync"

	"github.com/nanjj/cub/tasks"
)

type Routes struct {
	sync.Map
}

func (r *Routes) Add(target, via string) {
	r.Store(target, via)
}

func (r *Routes) Get(target string) (via string) {
	if v, ok := r.Load(target); ok {
		if via, ok = v.(string); ok {
			return
		}
	}
	return
}

func (r *Routes) Dispatch(targets tasks.Targets) (vias map[string]tasks.Targets) {
	vias = map[string]tasks.Targets{}
	for _, target := range targets {
		via := r.Get(target)
		if ss, ok := vias[via]; ok {
			ss = append(ss, target)
			vias[via] = ss
		} else {
			vias[via] = []string{target}
		}
	}
	return
}
