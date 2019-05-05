package sca_test

import (
	"reflect"
	"testing"

	"github.com/nanjj/cub/sca"
)

func TestEventClone(t *testing.T) {
	h1 := &sca.Event{}
	h2 := h1.Clone()
	if h1 == h2 {
		t.Fatal(&h1, &h2)
	}
	if !reflect.DeepEqual(h1, h2) {
		t.Fatal(h1, h2)
	}
}
