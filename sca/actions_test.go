package sca_test

import (
	"context"
	"reflect"
	"sort"
	"testing"

	"github.com/nanjj/cub/sca"
)

func TestActions(t *testing.T) {
	ping := func(context.Context, sca.Payload) (sca.Payload, error) {
		return nil, nil
	}
	tcs := []struct {
		actions map[string]sca.Action
		names   []string
	}{
		{},
		{map[string]sca.Action{
			"nil": nil}, nil,
		},
		{map[string]sca.Action{
			"ping": ping}, []string{"ping"},
		},
	}
	for _, tc := range tcs {
		t.Run("", func(t *testing.T) {
			if tc.names == nil {
				tc.names = []string{}
			}
			m := &sca.Actions{}
			for k, v := range tc.actions {
				m.Add(k, v)
			}
			names := m.Names()
			sort.Strings(tc.names)
			if !reflect.DeepEqual(names, tc.names) {
				t.Log(names)
				t.Log(tc.names)
				t.Fatal()
			}
			for _, name := range names {
				if a, ok := m.Get(name); !ok || a == nil {
					t.Fatal(name, ok, a)
				}
			}
		})

	}
}

func TestActionsNew(t *testing.T) {
	actions := sca.Actions{}
	name := actions.New(nil)
	if len(name) != 8 {
		t.Fatal(name)
	}
}
