package sca_test

import (
	"context"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/nanjj/cub/logs"
	"github.com/nanjj/cub/sca"
	"golang.org/x/sync/errgroup"
)

func TestRunnerJoin(t *testing.T) {
	name := func(id int) string {
		return fmt.Sprintf("r%d", id)
	}
	addr := func(id int) string {
		return fmt.Sprintf("tcp://127.0.0.1:100%d", id)
	}
	newConfig := func(name, listen, leader string) *sca.Config {
		return &sca.Config{RunnerName: name, RunnerListen: listen, LeaderListen: leader}
	}
	// r11 <- r21
	var g errgroup.Group
	r11, err := sca.NewRunner(newConfig(name(11), addr(11), ""))
	if err != nil {
		t.Fatal(err)
	}
	defer r11.Close()
	g.Go(r11.Run)
	r21, err := sca.NewRunner(newConfig(name(21), addr(21), addr(11)))
	if err != nil {
		t.Fatal(err)
	}
	defer r21.Close()
	g.Go(r21.Run)

	time.Sleep(time.Millisecond * 1000)

	if members := r21.Members(); len(members) != 0 {
		t.Fatal(members)
	}

	if r21.Leader() == nil {
		t.Fatal()
	}

	if r11.Leader() != nil {
		t.Fatal(r11.Leader())
	}
	if routes := r11.Routes(); len(routes) != 1 && routes["r21"] != name(21) {
		t.Fatal(routes)
	}

	if members := r11.Members(); len(members) != 1 {
		t.Fatal(r11)
	}
	// ping r11
	ch := make(chan time.Time, 1024)
	ping := func(ctx context.Context, req sca.Payload) (rep sca.Payload, err error) {
		sp, ctx := logs.StartSpanFromContext(ctx, "ping")
		defer sp.Finish()
		ch <- time.Now()
		return
	}
	r11.AddAction("ping", ping)
	r21.AddAction("ping", ping)
	event := &sca.Event{
		Action: "ping",
	}
	startTime := time.Now()
	if err := sca.SendEvent(context.Background(), r21.Leader(), event); err != nil {
		t.Fatal(err)
	}
	endTime := <-ch
	t.Log(endTime.Sub(startTime))
	event.Receiver = []string{name(11)}
	startTime = time.Now()
	if err := sca.SendEvent(context.Background(), r21.Leader(), event); err != nil {
		t.Fatal(err)
	}
	endTime = <-ch
	t.Log(endTime.Sub(startTime))
	// ping all
	event.Receiver.ToAll()
	startTime = time.Now()
	if err := sca.SendEvent(context.Background(), r21.Leader(), event); err != nil {
		t.Fatal(err)
	}
	endTime = <-ch
	t.Log(endTime.Sub(startTime))
	endTime = <-ch
	t.Log(endTime.Sub(startTime))
	// join r31
	// r11 <- r21 <- r31
	r31, err := sca.NewRunner(newConfig(name(31), addr(31), addr(21)))
	if err != nil {
		t.Fatal(err)
	}
	defer r31.Close()
	g.Go(r31.Run)
	// wait a while
	time.Sleep(time.Millisecond * 1000)
	// Check points:
	// | runner | routes           | members |
	// |--------+------------------+---------|
	// | r11    | r21:r21, r31:r21 | r21     |
	// | r21    | r31:r31          | r31     |
	// | r31    | -                | -       |
	//
	checkRoutes := func(r *sca.Runner, want map[string]string) {
		if routes := r.Routes(); !reflect.DeepEqual(want, routes) {
			t.Fatal(r.Name(), routes)
		}
	}
	checkRoutes(r11, map[string]string{
		"r21": "r21",
		"r31": "r21",
	})
	checkRoutes(r21, map[string]string{"r31": "r31"})
	checkRoutes(r31, map[string]string{})
	checkMembers := func(r *sca.Runner, want []string) {
		if members := r.Members(); !reflect.DeepEqual(want, members) {
			t.Fatal(r.Name(), members)
		}
	}
	checkMembers(r11, []string{"r21"})
	checkMembers(r21, []string{"r31"})
	checkMembers(r31, []string{})
	// ping r31 from r11
	r31.AddAction("ping", ping)
	event.Receiver = []string{r31.Name()}
	startTime = time.Now()
	if err := sca.SendEvent(context.Background(), r21.Leader(), event); err != nil {
		t.Fatal(err)
	}
	endTime = <-ch
	t.Log(endTime.Sub(startTime))
	// ping all
	event.Receiver.ToAll()
	startTime = time.Now()
	if err := sca.SendEvent(context.Background(), r21.Leader(), event); err != nil {
		t.Fatal(err)
	}
	// r11 pong
	endTime = <-ch
	t.Log(endTime.Sub(startTime))
	// r21 pong
	endTime = <-ch
	t.Log(endTime.Sub(startTime))
	// r31 pong
	endTime = <-ch
	t.Log(endTime.Sub(startTime))
	time.Sleep(time.Millisecond * 100)
	if len(ch) != 0 {
		t.Fatal()
	}
	// join r32 to r21
	r32, err := sca.NewRunner(newConfig(name(32), addr(32), addr(21)))
	if err != nil {
		t.Fatal(err)
	}
	defer r32.Close()
	g.Go(r32.Run)
	//r11 -> r21 -> r31
	//        |
	//        +---> r32
	//
	//| runner | routes                    | members   |
	//|--------+---------------------------+-----------|
	//| r11    | {r21:r21,r31:r21,r32:r21} | [r21]     |
	//| r21    | {r31:r31,r32:r32}         | [r31,r32] |
	//| r31    | -                         | -         |
	//| r32    | -                         | -         |
	time.Sleep(time.Millisecond * 100)
	checkRoutes(r11, map[string]string{
		"r21": "r21",
		"r31": "r21",
		"r32": "r21",
	})
	checkRoutes(r21, map[string]string{
		"r31": "r31",
		"r32": "r32",
	})
	checkRoutes(r31, map[string]string{})
	checkRoutes(r32, map[string]string{})
	checkMembers(r11, []string{"r21"})
	checkMembers(r21, []string{"r31", "r32"})
	checkMembers(r31, []string{})
	checkMembers(r32, []string{})
	// ping all
	r32.AddAction("ping", ping)
	event.Receiver.ToAll()
	startTime = time.Now()
	if err := sca.SendEvent(context.Background(), r21.Leader(), event); err != nil {
		t.Fatal(err)
	}
	// r11 pong
	endTime = <-ch
	t.Log(endTime.Sub(startTime))
	// r21 pong
	endTime = <-ch
	t.Log(endTime.Sub(startTime))
	// r31 pong
	endTime = <-ch
	t.Log(endTime.Sub(startTime))
	// r32 pong
	endTime = <-ch
	t.Log(endTime.Sub(startTime))
	time.Sleep(time.Millisecond * 100)
	if len(ch) != 0 {
		t.Fatal()
	}
}
