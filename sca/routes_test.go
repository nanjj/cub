package sca_test

import (
	"reflect"
	"testing"

	"github.com/nanjj/cub/sca"
	"github.com/nanjj/cub/sdo"
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
		name    string
		data    map[string]string
		targets sdo.Targets
		local   bool
		ups     sdo.Targets
		vias    map[string]sdo.Targets
	}{
		{"me", nil, nil, true, nil, nil},
		{"me", nil, []string{"me"}, true, nil, nil},
		{"me", map[string]string{}, []string{}, true, nil, nil},
		{"me", map[string]string{"a": "b"}, []string{}, true, nil, map[string]sdo.Targets{"b": []string{}}},
		{"me", map[string]string{"a": "b"}, []string{""}, false, []string{""}, nil},
		{"me", map[string]string{"a": "b"}, []string{"", "me", "a"}, true, []string{""}, map[string]sdo.Targets{"b": []string{"a"}}},
		{"me", map[string]string{"a": "b", "c": "d"}, []string{}, true, nil, map[string]sdo.Targets{"b": []string{}, "d": []string{}}},
		{"me", map[string]string{"a": "b"}, []string{"a"}, false, nil, map[string]sdo.Targets{"b": []string{"a"}}},
		{"me", map[string]string{"a": "b"}, []string{"a", "me"}, true, nil, map[string]sdo.Targets{"b": []string{"a"}}},
		{"me", map[string]string{"a": "b", "c": "d"}, []string{"a", "c"}, false, nil, map[string]sdo.Targets{"b": []string{"a"}, "d": []string{"c"}}},
		{"me", map[string]string{"a": "b", "c": "b"}, []string{"a", "c"}, false, nil, map[string]sdo.Targets{"b": []string{"a", "c"}}},
		{"me", map[string]string{"a": "b", "c": "b"}, []string{"a", "c", "e"}, false, []string{"e"}, map[string]sdo.Targets{"b": {"a", "c"}}},
	}

	for _, tc := range tcs {
		t.Run("", func(t *testing.T) {
			r := &sca.Rms{Name: tc.name}
			for k, v := range tc.data {
				r.AddRoute(k, v)
				r.AddMember(v, nil)
			}
			if tc.vias == nil {
				tc.vias = map[string]sdo.Targets{}
			}
			targets := tc.targets
			local, ups, vias := r.Dispatch(targets)
			if local != tc.local {
				t.Fatal(local, tc.local)
			}
			if !reflect.DeepEqual(ups, tc.ups) {
				t.Fatal(ups, tc.ups)
			}
			if !reflect.DeepEqual(vias, tc.vias) {
				t.Fatal(targets, vias, tc.vias)
			}
		})
	}
}

func TestRmsHasMember(t *testing.T) {
	type T struct {
		data map[string]string
		want bool
	}
	tcs := map[string]T{
		"1": T{},
		"2": T{map[string]string{}, false},
		"3": T{map[string]string{"a": "b"}, true},
		"4": T{map[string]string{"a": "b", "c": "d"}, true},
	}
	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			rms := sca.Rms{}
			for k, v := range tc.data {
				rms.AddRoute(k, v)
				rms.AddMember(v, nil)
			}
			if want := rms.HasMember(); want != tc.want {
				t.Fatal(want, tc)
			}
		})
	}
}
