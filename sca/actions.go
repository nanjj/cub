package sca

import (
	"context"
	"fmt"
	"math/rand"
	"sort"
	"sync"
	"time"
)

type Action func(context.Context, Payload) (Payload, error)
type Actions struct{ sync.Map }

func init() {
	rand.Seed(time.Now().UnixNano())
}

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

func (m *Actions) New(action Action) (name string) {
	for {
		name = fmt.Sprintf("cb-%05d", rand.Intn(99999))
		if _, ok := m.Get(name); !ok {
			m.Add(name, action)
			return
		}
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
