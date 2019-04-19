package sca

import (
	"context"
	"log"

	"github.com/opentracing/opentracing-go"
	"github.com/ugorji/go/codec"
	"nanomsg.org/go/mangos/v2"
)

func (h Head) Dup() (dup Head) {
	dup = Head{
		Id:       h.Id,
		Receiver: h.Receiver.Dup(),
		Sender:   h.Sender.Dup(),
	}
	return
}

func (e *Event) Dup() (dup *Event) {
	if e == nil {
		return
	}
	dup = &Event{
		Head:     e.Head.Dup(),
		Action:   e.Action,
		Callback: e.Callback,
	}
	l := len(e.Payload)
	payload := make([]DataObject, l)
	for i := 0; i < l; i++ {
		payload[i] = e.Payload[i].Dup()
	}
	dup.Payload = payload
	carrier := map[string]string{}
	for k, v := range e.Carrier {
		carrier[k] = v
	}
	dup.Carrier = carrier
	return
}

func (e *Event) Emit(ctx context.Context, sock mangos.Socket) (err error) {
	sp, ctx := StartSpanFromContext(ctx, "Emit")
	defer sp.Finish()
	tracer := sp.Tracer()
	if e.Carrier == nil {
		e.Carrier = map[string]string{}
	}
	Inject(tracer, sp.Context(), e.Carrier)
	out := make([]byte, 0, 1024)
	enc := codec.NewEncoderBytes(&out, cbor)
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

func (e *Event) Recv(sock mangos.Socket) (err error) {
	in, err := RetryRecv(sock)
	if err != nil {
		log.Println(err)
		return
	}
	dec := codec.NewDecoderBytes(in, cbor)
	if err = dec.Decode(e); err != nil {
		log.Println(err)
		return
	}
	return
}

func (e *Event) SpanContext(tracer opentracing.Tracer, name string) (sp opentracing.Span, ctx context.Context) {
	ctx = WithValues(context.Background(), e.Carrier)
	opts := []opentracing.StartSpanOption{}
	if sc, err := Extract(tracer, e.Carrier); err == nil {
		opts = append(opts, opentracing.ChildOf(sc))
	}
	sp = tracer.StartSpan(name, opts...)
	ctx = opentracing.ContextWithSpan(ctx, sp)
	return
}
