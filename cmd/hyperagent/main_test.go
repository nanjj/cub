package main

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/nanjj/cub/sca"
	"golang.org/x/sync/errgroup"
)

func TestMain(m *testing.M) {
	rt := func() int {
		os.Setenv("HYPERAGENT_RUNNER_IP", "127.0.0.1")
		return m.Run()
	}()
	os.Exit(rt)
}

func TestRunHyperAgentE(t *testing.T) {
	var g errgroup.Group
	g.Go(func() error {
		main()
		return nil
	})
	for i := 0; i < 10; i++ {
		listen, name := os.Getenv("HYPERAGENT_RUNNER_LISTEN"), os.Getenv("HYPERAGENT_RUNNER_NAME")
		if listen != "" && name != "" {
			break
		}
		time.Sleep(time.Millisecond)
	}
	listen, name := os.Getenv("HYPERAGENT_RUNNER_LISTEN"), os.Getenv("HYPERAGENT_RUNNER_NAME")
	if listen == "" || name == "" {
		t.Fatal()
	}
	main()
	node := sca.Node{Name: name, Listen: listen}
	node.Exit(context.Background())
	if err := g.Wait(); err != nil {
		t.Fatal(err)
	}
}
