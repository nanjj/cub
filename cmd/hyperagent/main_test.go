package main

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/avast/retry-go"
	"github.com/nanjj/cub/sca"
	"golang.org/x/sync/errgroup"
)

func TestRunHyperAgentE(t *testing.T) {
	os.Setenv("HYPERAGENT_RUNNER_IP", "127.0.0.1")
	oldRetry := sca.Retry
	defer func() { sca.Retry = oldRetry }()
	sca.Retry = append(sca.Retry, retry.Attempts(3), retry.Delay(time.Millisecond*10))
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
	go main()
	node := sca.Node{Name: name, Listen: listen}
	node.Exit(context.Background())
	if err := g.Wait(); err != nil {
		t.Fatal(err)
	}
}
