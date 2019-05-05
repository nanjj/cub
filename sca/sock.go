package sca

import (
	"context"

	"github.com/nanjj/cub/logs"
	"github.com/ugorji/go/codec"
	"go.uber.org/zap"
	"nanomsg.org/go/mangos/v2"
	"nanomsg.org/go/mangos/v2/protocol/pull"
	"nanomsg.org/go/mangos/v2/protocol/push"
)

func SendEvent(ctx context.Context, sock mangos.Socket, e *Event) (err error) {
	sp, ctx := logs.StartSpanFromContext(ctx, "SendEvent")
	defer sp.Finish()
	tracer := sp.Tracer()
	if e.Carrier == nil {
		e.Carrier = map[string]string{}
	}
	logs.Inject(tracer, sp.Context(), e.Carrier)
	out := make([]byte, 0, 1024)
	enc := codec.NewEncoderBytes(&out, msgpack)
	if err = enc.Encode(e); err != nil {
		sp.Error("Failed to encode", zap.Stack("stack"), zap.Error(err))
		return
	}
	if err = RetrySend(sock, out); err != nil {
		sp.Error("failed to send", zap.Stack("stack"), zap.Error(err))
		return
	}
	return
}

func RecvEvent(ctx context.Context, sock mangos.Socket, e *Event) (err error) {
	sp, ctx := logs.StartSpanFromContext(ctx, "RecvEvent")
	defer sp.Finish()
	in, err := RetryRecv(sock)
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

func NewClient(listen string) (sock mangos.Socket, err error) {
	sock, err = push.NewSocket()
	if err != nil {
		return
	}
	err = RetryDial(sock, listen)
	return
}

func NewServer(listen string) (sock mangos.Socket, err error) {
	sock, err = pull.NewSocket()
	if err != nil {
		return
	}
	err = RetryListen(sock, listen)
	return

}
