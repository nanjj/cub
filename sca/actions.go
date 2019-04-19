package sca

import (
	"context"
	"sort"
	"sync"
)

type Action func(context.Context, Payload) (Payload, error)
type Actions struct{ sync.Map }

func (m *Actions) Get(name string) (a Action, ok bool) {
	v, ok := m.Load(name)
	if ok && v != nil {
		if a, ok = v.(Action); ok {
			return
		}
	}
	return
}

func (m *Actions) Add(name string, action Action) {
	if action != nil {
		m.Store(name, action)
	}
}

func (m *Actions) Names() (names []string) {
	names = []string{}
	m.Range(func(k, v interface{}) bool {
		if v == nil {
			return true
		}
		if a, ok := v.(Action); !ok || a == nil {
			return true
		}
		if name, ok := k.(string); ok {
			names = append(names, name)
		}
		return true
	})
	sort.Strings(names) // sort to stable order
	return
}
