package sca

import (
	"sync"

	"nanomsg.org/go/mangos/v2"
)

type Members struct {
	sync.Map
}

func (m *Members) Names() (names []string) {
	names = []string{}
	m.Range(func(k, v interface{}) bool {
		if name, ok := k.(string); ok {
			names = append(names, name)
		}
		return true
	})
	return
}

func (m *Members) Get(name string) (sock mangos.Socket, ok bool) {
	v, ok := m.Load(name)
	if ok {
		sock, ok = v.(mangos.Socket)
	}
	return
}

func (m *Members) Add(name string, sock mangos.Socket) {
	m.Store(name, sock)
}
