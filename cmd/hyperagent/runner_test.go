package main

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"golang.org/x/sync/errgroup"
)

func TestRunnerJoin(t *testing.T) {
	name := func(id int) string {
		return fmt.Sprintf("r%d", id)
	}
	addr := func(id int) string {
		return fmt.Sprintf("tcp://127.0.0.1:100%d", id)
	}
	// r11 <- r21
	var g errgroup.Group
	r11, err := NewRunner(name(11), addr(11), "")
	if err != nil {
		t.Fatal(err)
	}
	defer r11.Close()
	g.Go(r11.Run)
	r21, err := NewRunner(name(21), addr(21), addr(11))
	if err != nil {
		t.Fatal(err)
	}
	defer r21.Close()
	g.Go(r21.Run)

	time.Sleep(time.Millisecond * 100)

	if members := r21.Members(); len(members) != 0 {
		t.Fatal(r21.members)
	}

	if r21.leader == nil {
		t.Fatal()
	}

	if r11.leader != nil {
		t.Fatal(r11.leader)
	}
	if routes := r11.Routes(); len(routes) != 1 && routes["r21"] != name(21) {
		t.Fatal(routes)
	}

	if members := r11.Members(); len(members) != 1 {
		t.Fatal(r11)
	}
	// ping r11
	ch := make(chan time.Time, 1024)
	ping := func(args []string) (err error) {
		ch <- time.Now()
		return
	}
	r11.AddHandler("ping", ping)
	r21.AddHandler("ping", ping)
	task := &Task{
		Name: "ping",
	}
	startTime := time.Now()
	if err := task.Send(r21.Leader()); err != nil {
		t.Fatal(err)
	}
	endTime := <-ch
	t.Log(endTime.Sub(startTime))
	task.Targets = []string{name(11)}
	startTime = time.Now()
	if err := task.Send(r21.Leader()); err != nil {
		t.Fatal(err)
	}
	endTime = <-ch
	t.Log(endTime.Sub(startTime))
	// ping all
	task.Targets.ToAll()
	startTime = time.Now()
	if err := task.Send(r21.Leader()); err != nil {
		t.Fatal(err)
	}
	endTime = <-ch
	t.Log(endTime.Sub(startTime))
	endTime = <-ch
	t.Log(endTime.Sub(startTime))
	// join r31
	// r11 <- r21 <- r31
	r31, err := NewRunner(name(31), addr(31), addr(21))
	if err != nil {
		t.Fatal(err)
	}
	defer r31.Close()
	g.Go(r31.Run)
	// wait a while
	time.Sleep(time.Millisecond * 100)
	// Check points:
	// | runner | routes           | members |
	// |--------+------------------+---------|
	// | r11    | r21:r21, r31:r21 | r21     |
	// | r21    | r31:r31          | r31     |
	// | r31    | -                | -       |
	//
	checkRoutes := func(r *Runner, want map[string]string) {
		if routes := r.Routes(); !reflect.DeepEqual(want, routes) {
			t.Fatal(r.name, routes)
		}
	}
	checkRoutes(r11, map[string]string{
		"r21": "r21",
		"r31": "r21",
	})
	checkRoutes(r21, map[string]string{"r31": "r31"})
	checkRoutes(r31, map[string]string{})
	checkMembers := func(r *Runner, want []string) {
		if members := r.Members(); !reflect.DeepEqual(want, members) {
			t.Fatal(r.name, members)
		}
	}
	checkMembers(r11, []string{"r21"})
	checkMembers(r21, []string{"r31"})
	checkMembers(r31, []string{})
	// ping r31 from r11
	r31.AddHandler("ping", ping)
	task.Targets = []string{r31.Name()}
	startTime = time.Now()
	if err := task.Send(r21.Leader()); err != nil {
		t.Fatal(err)
	}
	endTime = <-ch
	t.Log(endTime.Sub(startTime))
	// ping all
	task.Targets.ToAll()
	startTime = time.Now()
	if err := task.Send(r21.Leader()); err != nil {
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
	r32, err := NewRunner(name(32), addr(32), addr(21))
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
	r32.AddHandler("ping", ping)
	task.Targets.ToAll()
	startTime = time.Now()
	if err := task.Send(r21.Leader()); err != nil {
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
