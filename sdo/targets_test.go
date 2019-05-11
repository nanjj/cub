package sdo_test

import (
	"reflect"
	"testing"

	"github.com/nanjj/cub/sdo"
)

func TestTargetsToAll(t *testing.T) {
	toall := sdo.Targets([]string{})
	if !toall.Down() {
		t.Fatal(toall)
	}
	toall.ToLocal()
	if !toall.Local() {
		t.Fatal(toall)
	}
	toall.ToDown()
	if !toall.Down() {
		t.Fatal(toall)
	}
}

func TestTargetsDeepEqual(t *testing.T) {
	t.Log(reflect.DeepEqual([]string{}, []string{}))
	t.Log(reflect.DeepEqual(map[string]sdo.Targets{"": []string{}}, map[string]sdo.Targets{"": []string{}}))
}

func TestTargetsClone(t *testing.T) {
	type T struct {
		name    string
		targets sdo.Targets
	}
	tcs := []T{{"nil", nil},
		{"empty", sdo.Targets{}},
		{"one", sdo.Targets{"one"}},
		{"two", sdo.Targets{"one", "two"}}}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			want := tc.targets.Clone()
			if !reflect.DeepEqual(want, tc.targets) {
				t.Fatalf("%v\n%v", want, tc.targets)
			}
			if len(want) > 0 {
				want[0] = "modified"
				if reflect.DeepEqual(want, tc.targets) {
					t.Fatalf("%v\n%v", want, tc.targets)
				}
			}
		})
	}
}
