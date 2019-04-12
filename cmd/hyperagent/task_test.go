package main

import "testing"

func TestTargets(t *testing.T) {
	var targets Targets
	if targets != nil || !targets.Local() {
		t.Fatal(targets)
	}
	targets = []string{}
	if targets == nil || !targets.All() {
		t.Fatal(targets)
	}
	targets.ToLocal()
	if !targets.Local() {
		t.Fatal(targets)
	}
	targets.ToAll()
	if !targets.All() {
		t.Fatal(targets)
	}
}
