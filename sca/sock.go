package sca

import (
	"context"
	"log"

	"github.com/ugorji/go/codec"
	"nanomsg.org/go/mangos/v2"
)

func SendEvent(ctx context.Context, sock mangos.Socket, e *Event) (err error) {
	sp, ctx := StartSpanFromContext(ctx, "SendEvent")
	defer sp.Finish()
	tracer := sp.Tracer()
	if e.Carrier == nil {
		e.Carrier = map[string]string{}
	}
	Inject(tracer, sp.Context(), e.Carrier)
	out := make([]byte, 0, 1024)
	enc := codec.NewEncoderBytes(&out, msgpack)
	if err = enc.Encode(e); err != nil {
		log.Println(err)
		return
	}
	if err = RetrySend(sock, out); err != nil {
		log.Println(err)
		return
	}
	return
}

func RecvEvent(ctx context.Context, sock mangos.Socket, e *Event) (err error) {
	sp, ctx := StartSpanFromContext(ctx, "RecvEvent")
	defer sp.Finish()
	in, err := RetryRecv(sock)
	if err != nil {
		sp.Println(err)
		return
	}
	dec := codec.NewDecoderBytes(in, msgpack)
	if err = dec.Decode(e); err != nil {
		sp.Println(err)
		return
	}
	return
}