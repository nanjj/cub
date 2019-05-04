package sca_test

import (
	"reflect"
	"testing"

	"github.com/nanjj/cub/sca"
)

func TestTargetsToAll(t *testing.T) {
	toall := sca.Targets([]string{})
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
	t.Log(reflect.DeepEqual(map[string]sca.Targets{"": []string{}}, map[string]sca.Targets{"": []string{}}))
}
