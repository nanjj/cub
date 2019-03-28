package main

import (
	"fmt"
	"testing"
)

func TestColor(t *testing.T) {
	tcs := []struct {
		color Color
		want  string
	}{
		{Red, "Red"},
		{Green, "Green"},
		{Yellow, "Yellow"},
	}
	for _, tc := range tcs {
		want := fmt.Sprint(tc.color)
		if want != tc.want {
			t.Fatal(tc.color, want, tc.want)
		}
	}
}
