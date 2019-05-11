package sca_test

import (
	"context"
	"testing"
	"time"

	retry "github.com/avast/retry-go"
	"github.com/nanjj/cub/sca"
	"github.com/nanjj/cub/sdo"
	"golang.org/x/sync/errgroup"
	"nanomsg.org/go/mangos"
)

func TestNodeRecv(t *testing.T) {
	oldRetry := sca.Retry
	defer func() { sca.Retry = oldRetry }()
	sca.Retry = append(sca.Retry, retry.Attempts(1), retry.Delay(time.Millisecond))
	listen := "tcp://127.0.0.1:12345"
	var g errgroup.Group
	n1 := sca.Node{Name: "node1", Listen: listen}
	defer g.Wait()
	defer n1.Close()
	n2 := sca.Node{Name: "node2", Listen: listen}
	defer n2.Close()
	e := &sdo.DataGraph{}
	ctx := context.Background()
	g.Go(func() (err error) {
		if err = n1.Recv(ctx, e); err != mangos.ErrClosed {
			t.Fatal(err)
		}
		return
	})
	time.Sleep(time.Millisecond * 10)
	if err := n2.Recv(ctx, e); err != sca.ErrorListen {
		t.Fatal(err)
	}
}

func TestNodeSend(t *testing.T) {
	listen := "tcp://127.0.0.1:12345"
	n1 := sca.Node{Name: "node1", Listen: listen}
	defer n1.Close()
	ctx := context.Background()
	e := &sdo.DataGraph{} // die event
	if err := n1.Send(ctx, e); err != sca.ErrorDial {
		t.Fatal(err)
	}
}
