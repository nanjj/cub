package sca_test

import (
	"context"
	"testing"
	"time"

	"github.com/nanjj/cub/logs"
	"github.com/nanjj/cub/sca"
	"github.com/nanjj/cub/sdo"
)

func TestRunnerPingSelf(t *testing.T) {
	r := _testNewRunner
	n := _testRunnerName
	rr := func(id, pid int) *sca.Runner {
		runner, err := r(id, pid)
		if err != nil {
			t.Fatal(err)
		}
		return runner
	}
	r11 := rr(11, 0)
	defer r11.Close()
	ctx := context.Background()
	req := sdo.Payload{}
	// ping inside
	startTime := time.Now().UTC()
	rep, err := r11.Ping(ctx, req)
	if err != nil {
		t.Fatal(err)
	}

	if l := len(rep); l != 2 {
		t.Fatal(l)
	}

	if name := string(rep[0]); name != n(11) {
		t.Fatal(name)
	}
	endTime := time.Time{}
	if err := rep[1].Decode(&endTime); err != nil {
		t.Fatal(err)
	}
	t.Log(endTime.Sub(startTime))
	// ping outside
	results := make(chan sdo.Payload, 1024)
	r11.AddAction("pong", func(ctx context.Context, req sdo.Payload) (rep sdo.Payload, err error) {
		sp, ctx := logs.StartSpanFromContext(ctx, "Pong")
		defer sp.Finish()
		results <- req
		return
	})
	startTime = time.Now().UTC()
	if err = r11.Self().Send(ctx, &sdo.DataGraph{
		Summary: sdo.Summary{
			Action:   "ping",
			Callback: "pong"},
	}); err != nil {
		t.Fatal(err)
	}
	rep = <-results
	endTime = time.Time{}
	if string(rep[0]) != r11.Name() {
		t.Fatal(rep)
	}
	if err = rep[1].Decode(&endTime); err != nil {
		t.Fatal(err)
	}
	t.Log(endTime.Sub(startTime))
}

func TestRunnerPingTree(t *testing.T) {
	r := _testNewRunner
	n := _testRunnerName
	rr := func(id, pid int) *sca.Runner {
		runner, err := r(id, pid)
		if err != nil {
			t.Fatal(err)
		}
		return runner
	}
	//r11 - r21 - r32
	//      r22
	//      r23 - r31
	r11 := rr(11, 0) // top
	defer r11.Close()
	r21 := rr(21, 11)
	defer r21.Close()
	r32 := rr(32, 21)
	defer r32.Close()
	r22 := rr(22, 11)
	defer r22.Close()
	r23 := rr(23, 11)
	defer r23.Close()
	r31 := rr(31, 23)
	defer r31.Close()
	if r11 == nil || r21 == nil || r22 == nil || r23 == nil || r31 == nil || r32 == nil {
		t.Fatal(r11, r21, r22, r23, r31, r32)
	}
	// wait r31, r32 login r11 ready
	for {
		routes := r11.Routes()
		if _, ok := routes[n(31)]; ok {
			if _, ok = routes[n(32)]; ok {
				break
			}
		}
		time.Sleep(time.Millisecond)
	}
	results := make(chan sdo.Payload, 1024)
	pongCb := func(ctx context.Context,
		req sdo.Payload) (rep sdo.Payload, err error) {
		sp, ctx := logs.StartSpanFromContext(ctx, "Pong")
		defer sp.Finish()
		results <- req
		return
	}
	pong := r11.AddCallback(pongCb)
	// Now ping r31 from r11
	if err := r11.Self().Send(context.Background(), &sdo.DataGraph{
		Summary: sdo.Summary{
			Action:   "ping",
			To:       sdo.Targets{n(31)},
			From:     n(11),
			Callback: pong,
		},
	}); err != nil {
		t.Fatal(err)
	}
	rep := <-results
	if len(rep) != 2 {
		t.Fatal(rep)
	}
	if err := r32.Self().Send(context.Background(), &sdo.DataGraph{
		Summary: sdo.Summary{
			Action:   "ping",
			To:       sdo.Targets{n(31)},
			From:     n(32),
			Callback: r32.AddCallback(pongCb),
		},
	}); err != nil {
		t.Fatal(err)
	}
	rep = <-results
	if len(rep) != 2 {
		t.Fatal(rep)
	}
	// Now ping r31 and r32 from r22
	if err := r22.Self().Send(context.Background(), &sdo.DataGraph{
		Summary: sdo.Summary{Action: "ping",
			To:       sdo.Targets{n(31), n(32)},
			From:     n(22),
			Callback: r22.AddCallback(pongCb),
		},
	}); err != nil {
		t.Fatal(err)
	}
	rep = <-results
	if len(rep) != 2 {
		t.Fatal(rep)
	}
	rep = <-results
	if len(rep) != 2 {
		t.Fatal(rep)
	}
}
