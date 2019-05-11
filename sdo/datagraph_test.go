package sdo_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/nanjj/cub/sdo"
)

func TestDataGraphLoadBytes(t *testing.T) {
	tcs := []struct {
		summary sdo.Summary
		payload sdo.Payload
	}{
		{},
		{sdo.Summary{}, sdo.Payload{}},
		{sdo.Summary{
			Action: "ping",
		}, sdo.Payload{}},
		{sdo.Summary{
			Action: "ping",
		}, sdo.Payload{sdo.DataObject("time")}},
		{sdo.Summary{
			Action: "ping",
		}, sdo.Payload{sdo.DataObject("time"), sdo.DataObject("members")}},
	}
	for i, tc := range tcs {
		name := fmt.Sprintf("%02d", i)
		t.Run(name, func(t *testing.T) {
			g1 := sdo.DataGraph{
				Summary: tc.summary,
				Payload: tc.payload,
			}
			b, err := g1.Bytes()
			if err != nil {
				t.Fatal(err)
			}
			t.Logf("len=%d, \ncontent=%x\n%s", len(b), b, string(b))
			g2 := sdo.DataGraph{}
			if err = g2.Load(b); err != nil {
				t.Fatal(err)
			}

			if !reflect.DeepEqual(g1.Summary, g2.Summary) {
				t.Log(g1.Summary, g1.Summary.Lens == nil)
				t.Log(g2.Summary, g2.Summary.Lens == nil)
				t.Fatal()
			}
			if !reflect.DeepEqual(g1.Payload, g2.Payload) {
				t.Logf("%v,%v", g1.Payload, g1.Payload == nil)
				t.Logf("%v,%v", g2.Payload, g2.Payload == nil)
				t.Fatal()
			}
			if !reflect.DeepEqual(g1, g2) {
				t.Log(g1)
				t.Log(g2)
				t.Fatal()
			}
		})
	}
}

func TestDatagraphClone(t *testing.T) {
	tcs := []struct {
		name string
		g    *sdo.DataGraph
	}{
		{"nil", nil},
		{"empty", &sdo.DataGraph{}},
		{"PayloadEmpty", &sdo.DataGraph{Payload: sdo.Payload{}}},
		{"PayloadOne", &sdo.DataGraph{Payload: sdo.Payload{sdo.DataObject("one")}}},
		{"PayloadTwo", &sdo.DataGraph{Payload: sdo.Payload{
			sdo.DataObject("one"),
			sdo.DataObject("two"),
		}}},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			g := tc.g.Clone()
			if !reflect.DeepEqual(g, tc.g) {
				t.Fatalf("%v\n%v", g, tc.g)
			}
		})
	}
}
func TestDataGraphLoad(t *testing.T) {
	g := &sdo.DataGraph{}
	if err := g.Load([]byte{}); err != sdo.ErrInvalidData {
		t.Fatal(err)
	}

	if err := g.Load([]byte("invalid")); err == nil {
		t.Fatal(err)
	}
	b := []byte{}
	fmt.Sscanf("87a16900a174c0a166a0a161a470696e67a162a0a163c0a16c940474696d65", "%x", &b)
	if len(b) == 0 {
		t.Fatal(b)
	}
	if err := g.Load(b); err != sdo.ErrIncompletData {
		t.Fatal(err)
	}
}
