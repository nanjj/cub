package sca_test

import (
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
