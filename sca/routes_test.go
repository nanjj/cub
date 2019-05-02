package sca_test

import (
	"reflect"
	"testing"

	"github.com/nanjj/cub/sca"
	"golang.org/x/sync/errgroup"
)

func TestRoutesRace(t *testing.T) {
	data := map[string]string{
		"r1": "r2",
		"r3": "r2",
		"r4": "r2",
		"r5": "r2",
	}
	var g errgroup.Group
	routes := &sca.Rms{}
	// Add
	g.Go(func() error {
		for i := 0; i < 100; i++ {
			for k, v := range data {
				routes.AddRoute(k, v)
			}
		}
		return nil
	})
	// Get
	g.Go(func() error {
		for i := 0; i < 100; i++ {
			for k, v := range data {
				if via := routes.GetRoute(k); via != "" && via != v {
					t.Fatal(via, v)
				}
			}
		}
		return nil
	})
	// Delete
	g.Go(func() error {
		for i := 0; i < 100; i++ {
			for k := range data {
				routes.DelRoute(k)
			}
		}
		return nil
	})
}

func TestRoutesDispatch(t *testing.T) {
	tcs := []struct {
		data    map[string]string
		targets []string
		vias    map[string]sca.Targets
	}{
		{},
		{map[string]string{}, []string{}, map[string]sca.Targets{}},
		{map[string]string{"a": "b"}, []string{"a"}, map[string]sca.Targets{"b": []string{"a"}}},
		{map[string]string{"a": "b", "c": "d"}, []string{"a", "c"}, map[string]sca.Targets{"b": []string{"a"}, "d": []string{"c"}}},
		{map[string]string{"a": "b", "c": "b"}, []string{"a", "c"}, map[string]sca.Targets{"b": []string{"a", "c"}}},
		{map[string]string{"a": "b", "c": "b"}, []string{"a", "c", "e"}, map[string]sca.Targets{"b": {"a", "c"}, "": {"e"}}},
	}

	for _, tc := range tcs {
		t.Run("", func(t *testing.T) {
			r := &sca.Rms{}
			for k, v := range tc.data {
				r.AddRoute(k, v)
			}
			if tc.vias == nil {
				tc.vias = map[string]sca.Targets{}
			}
			targets := tc.targets
			vias := r.Dispatch(targets)
			if !reflect.DeepEqual(vias, tc.vias) {
				t.Fatal(vias, tc.vias)
			}
		})
	}
}
