package main

import (
	"reflect"
	"testing"
)

func TestDispatch(t *testing.T) {
	tcs := []struct {
		data    map[string]string
		targets []string
		vias    map[string]Targets
	}{
		{},
		{map[string]string{}, []string{}, map[string]Targets{}},
		{map[string]string{"a": "b"}, []string{"a"}, map[string]Targets{"b": []string{"a"}}},
		{map[string]string{"a": "b", "c": "d"}, []string{"a", "c"}, map[string]Targets{"b": []string{"a"}, "d": []string{"c"}}},
		{map[string]string{"a": "b", "c": "b"}, []string{"a", "c"}, map[string]Targets{"b": []string{"a", "c"}}},
		{map[string]string{"a": "b", "c": "b"}, []string{"a", "c", "e"}, map[string]Targets{"b": {"a", "c"}, "": {"e"}}},
	}

	for _, tc := range tcs {
		t.Run("", func(t *testing.T) {
			r := &Routes{}
			for k, v := range tc.data {
				r.Add(k, v)
			}
			if tc.vias == nil {
				tc.vias = map[string]Targets{}
			}
			targets := tc.targets
			vias := r.Dispatch(targets)
			if !reflect.DeepEqual(vias, tc.vias) {
				t.Fatal(vias, tc.vias)
			}
		})
	}
}
