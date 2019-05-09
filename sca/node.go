package sca

import (
	"context"
	"errors"
	"sync"

	"github.com/nanjj/cub/logs"
	"github.com/ugorji/go/codec"
	"go.uber.org/zap"
	"nanomsg.org/go/mangos/protocol/pull"
	"nanomsg.org/go/mangos/protocol/push"
	"nanomsg.org/go/mangos/v2"
)

var (
	ErrorListen = errors.New("Failed to listen")
	ErrorDial   = errors.New("Failed to dial")
	ErrorSock   = errors.New("Failed to new socket")
	ErrorDecode = errors.New("Failed to decode")
	ErrorEncode = errors.New("Failed to encode")
	ErrorSend   = errors.New("Failed to send")
	ErrorRecv   = errors.New("Failed to receive")
)

type Node struct {
	Name    string
	Listen  string
	puller  mangos.Socket
	pusher  mangos.Socket
	pullerf sync.Once
	pusherf sync.Once
}

func (n *Node) Close() (err error) {
	if n == nil {
		return
	}
	return Kill(n.pusher, n.puller)
}

func (n *Node) Send(ctx context.Context, e *Event) (err error) {
	sp, ctx := logs.StartSpanFromContext(ctx, "Send")
	defer sp.Finish()
	n.pusherf.Do(func() {
		if n.pusher, err = push.NewSocket(); err != nil {
			sp.Error("Failed to new socket",
				zap.Error(err),
				zap.Stack("stack"))
			return
		}
		listen := n.Listen
		if err = Retry.Dial(n.pusher, listen); err != nil {
			sp.Error("Failed to dial",
				zap.Error(err),
				zap.Stack("stack"),
				zap.String("listen", listen))
			n.pusher = nil
			return
		}
	})
	sock := n.pusher
	if err != nil || sock == nil {
		err = ErrorDial
		return
	}
	tracer := sp.Tracer()
	if e.Carrier == nil {
		e.Carrier = map[string]string{}
	}
	logs.Inject(tracer, sp.Context(), e.Carrier)
	out := make([]byte, 0, 1024)
	enc := codec.NewEncoderBytes(&out, msgpack)
	if err = enc.Encode(e); err != nil {
		sp.Error("Failed to encode", zap.Stack("stack"), zap.Error(err))
		err = ErrorEncode
		return
	}
	if err = Retry.Send(sock, out); err != nil {
		sp.Error("failed to send", zap.Stack("stack"), zap.Error(err))
		err = ErrorSend
		return
	}
	return
}

func (n *Node) Recv(ctx context.Context, e *Event) (err error) {
	sp, ctx := logs.StartSpanFromContext(ctx, "Recv")
	defer sp.Finish()
	n.pullerf.Do(func() {
		if n.puller, err = pull.NewSocket(); err != nil {
			sp.Error("Failed to new socket", zap.Stack("stack"), zap.Error(err))
			return
		}
		listen := n.Listen
		if err = Retry.Listen(n.puller, listen); err != nil {
			sp.Error("Failed to listen",
				zap.Stack("stack"), zap.Error(err),
				zap.String("listen", listen))
			n.puller = nil
			return
		}
	})
	sock := n.puller
	if err != nil || sock == nil {
		err = ErrorListen
		return
	}
	in, err := Retry.Recv(sock)
	if err != nil {
		sp.Error("Failed to receive", zap.Stack("stack"), zap.Error(err))
		return
	}
	dec := codec.NewDecoderBytes(in, msgpack)
	if err = dec.Decode(e); err != nil {
		sp.Error("Failed to decode", zap.Stack("stack"), zap.Error(err))
		return
	}
	return
}

func (n *Node) Exit(ctx context.Context) (err error) {
	err = n.Send(ctx, &Event{Action: EXIT})
	return
}

func (n *Node) Login(ctx context.Context, name, listen string) (err error) {
	err = n.Send(ctx,
		&Event{
			Action: "login",
			Payload: Payload{
				DataObject(name),
				DataObject(listen)},
			From: name})
	return
}

func (n *Node) Logout(ctx context.Context, name string) (err error) {
	err = n.Send(ctx, &Event{
		Action:  "logout",
		Payload: Payload{DataObject(name)},
		From:    name,
	})
	return
}
