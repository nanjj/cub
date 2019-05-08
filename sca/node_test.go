package sca_test

import (
	"context"
	"testing"
	"time"

	"github.com/nanjj/cub/sca"
	"golang.org/x/sync/errgroup"
	"nanomsg.org/go/mangos"
)

func TestNodeRecv(t *testing.T) {
	listen := "tcp://127.0.0.1:12345"
	var g errgroup.Group
	n1 := sca.Node{Name: "node1", Listen: listen}
	defer g.Wait()
	defer n1.Close()
	n2 := sca.Node{Name: "node2", Listen: listen}
	defer n2.Close()
	e := &sca.Event{}
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
	e := &sca.Event{} // die event
	if err := n1.Send(ctx, e); err != sca.ErrorDial {
		t.Fatal(err)
	}
}
