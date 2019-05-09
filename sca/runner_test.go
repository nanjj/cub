package sca_test

import (
	"context"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/nanjj/cub/logs"
	"github.com/nanjj/cub/sca"
)

func _testRunnerName(id int) string {
	return fmt.Sprintf("r%d", id)
}
func _testRunnerAddr(id int) string {
	if id <= 0 {
		return ""
	}
	return fmt.Sprintf("tcp://127.0.0.1:100%d", id)
}

func _testNewConfig(name, listen, leader string) *sca.Config {
	return &sca.Config{RunnerName: name, RunnerListen: listen, LeaderListen: leader}
}

func _testNewRunner(id int, leader int) (r *sca.Runner, err error) {
	return sca.NewRunner(_testNewConfig(_testRunnerName(id), _testRunnerAddr(id), _testRunnerAddr(leader)))
}

func _testWaitReady(l *sca.Runner, id int, exist bool) {
	name := _testRunnerName(id)
	for i := 0; i < 10; i++ {
		if _, ok := l.Routes()[name]; ok == exist {
			break
		}
		time.Sleep(time.Millisecond * 10)
	}
}

func TestRunnerLogin(t *testing.T) {
	n := _testRunnerName
	r := _testNewRunner
	w := _testWaitReady
	// r11 <- r21
	r11, err := r(11, 0)
	if err != nil {
		t.Fatal(err)
	}
	defer r11.Close()
	r21, err := r(21, 11)
	if err != nil {
		t.Fatal(err)
	}
	defer r21.Close()

	w(r11, 21, true)

	if members := r21.Members(); len(members) != 0 {
		t.Fatal(members)
	}

	if r21.Leader() == nil {
		t.Fatal()
	}

	if r11.Leader() != nil {
		t.Fatal(r11.Leader())
	}
	if routes := r11.Routes(); len(routes) != 1 && routes["r21"] != n(21) {
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

	if err := r21.Leader().Send(context.Background(), event); err != nil {
		t.Fatal(err)
	}
	endTime := <-ch
	t.Log(endTime.Sub(startTime))
	event.To = []string{n(11)}
	startTime = time.Now()
	if err := r21.Leader().Send(context.Background(), event); err != nil {
		t.Fatal(err)
	}
	endTime = <-ch
	t.Log(endTime.Sub(startTime))
	// ping all
	event.To.ToDown()
	startTime = time.Now()
	if err := r21.Leader().Send(context.Background(), event); err != nil {
		t.Fatal(err)
	}
	endTime = <-ch
	t.Log(endTime.Sub(startTime))
	endTime = <-ch
	t.Log(endTime.Sub(startTime))
	// login r31
	// r11 <- r21 <- r31
	r31, err := r(31, 21)
	if err != nil {
		t.Fatal(err)
	}
	defer r31.Close()
	// wait a while
	w(r21, 31, true)
	w(r11, 31, true)
	// Check points:
	// | runner | routes           | members |
	// |--------+------------------+---------|
	// | r11    | r21:r21, r31:r21 | r21     |
	// | r21    | r31:r31          | r31     |
	// | r31    | -                | -       |
	//
	checkRoutes := func(r *sca.Runner, want map[string]string) {
		if routes := r.Routes(); !reflect.DeepEqual(want, routes) {
			t.Fatal(r.Name(), routes, want)
		}
	}
	checkRoutes(r11, map[string]string{
		"r21": "r21",
		"r31": "r21",
	})
	checkRoutes(r21, map[string]string{"r31": "r31"})
	checkRoutes(r31, map[string]string{})
	checkMembers := func(r *sca.Runner, want sca.Set) {
		if members := r.Members(); !reflect.DeepEqual(want, members) {
			t.Fatal("node:", r.Name(), "actual:", members, "want:", want)
		}
	}
	checkMembers(r11, sca.NewSet("r21"))
	checkMembers(r21, sca.NewSet("r31"))
	checkMembers(r31, sca.NewSet())
	// ping r31 from r11
	r31.AddAction("ping", ping)
	event.To = []string{r31.Name()}
	startTime = time.Now()
	if err := r21.Leader().Send(context.Background(), event); err != nil {
		t.Fatal(err)
	}
	endTime = <-ch
	t.Log(endTime.Sub(startTime))
	// ping all
	event.To.ToDown()
	startTime = time.Now()
	if err := r21.Leader().Send(context.Background(), event); err != nil {
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
	time.Sleep(time.Millisecond * 10)
	if len(ch) != 0 {
		t.Fatal()
	}
	// login r32 to r21
	r32, err := r(32, 21)
	if err != nil {
		t.Fatal(err)
	}
	defer r32.Close()
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
	w(r11, 32, true)
	w(r21, 32, true)
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
	checkMembers(r11, sca.NewSet("r21"))
	checkMembers(r21, sca.NewSet("r31", "r32"))
	checkMembers(r31, sca.NewSet())
	checkMembers(r32, sca.NewSet())
	// ping all
	r32.AddAction("ping", ping)
	event.To.ToDown()
	startTime = time.Now()
	if err := r21.Leader().Send(context.Background(), event); err != nil {
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
	time.Sleep(time.Millisecond * 10)
	if len(ch) != 0 {
		t.Fatal()
	}
}

func TestRunnerLogout(t *testing.T) {
	r := _testNewRunner
	w := _testWaitReady
	rr := func(id, pid int) *sca.Runner {
		runner, err := r(id, pid)
		if err != nil {
			t.Fatal(err)
		}
		return runner
	}
	r10 := rr(10, 0)
	r20 := rr(20, 10)
	w(r10, 20, true)
	routes := r10.Routes()
	if len(routes) != 1 {
		t.Fatal()
	}

	r20.Close()
	for i := 0; i < 10; i++ {
		routes := r10.Routes()
		if len(routes) == 0 {
			break
		}
		time.Sleep(time.Millisecond * 10)
	}
	routes = r10.Routes()
	if len(routes) != 0 {
		t.Fatal(routes, r10.Members())
	}
	r10.Close()
	r10.Wait()
}

func TestRunnerNodeType(t *testing.T) {
	r := _testNewRunner
	n := _testRunnerName
	rr := func(id, pid int) *sca.Runner {
		runner, err := r(id, pid)
		if err != nil {
			t.Fatal(err)
		}
		if pid != 0 {
			for i := 0; i < 10; i++ {
				members := runner.Members()
				if members != nil && members[n(pid)] {
					break
				}
				time.Sleep(time.Millisecond)
			}
		}
		return runner
	}
	r10 := rr(10, 0) // root
	defer r10.Close()
	r20 := rr(20, 10) // branch
	defer r20.Close()
	r30 := rr(30, 20) // leaf
	defer r30.Close()
	if nodeType := r10.NodeType(); nodeType != sca.ROOT {
		t.Fatal(nodeType)
	}
	if nodeType := r20.NodeType(); nodeType != sca.BRANCH {
		t.Fatal(nodeType)
	}
	if nodeType := r30.NodeType(); nodeType != sca.LEAF {
		t.Fatal(nodeType)
	}

}
