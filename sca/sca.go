package sca

import (
	"context"
	"log"

	"github.com/ugorji/go/codec"
	"nanomsg.org/go/mangos/v2"
)

var (
	cbor = &codec.CborHandle{}
)

type Targets []string

func (targets Targets) All() bool {
	return targets != nil && len(targets) == 0
}

func (targets Targets) Local() bool {
	return targets == nil
}

func (targets *Targets) ToAll() {
	*targets = []string{}
}

func (targets *Targets) ToLocal() {
	*targets = nil
}

func (t Targets) Dup() (dup Targets) {
	if t == nil {
		return
	}
	dup = append([]string{}, t...)
	return
}

type Head struct {
	Id       int64   `codec:"id"`
	Receiver Targets `codec:"receiver"`
	Sender   Targets `codec:"sender"`
}

func (h Head) Dup() (dup Head) {
	dup = Head{
		Id:       h.Id,
		Receiver: h.Receiver.Dup(),
		Sender:   h.Sender.Dup(),
	}
	return
}

//go:generate codecgen -o cg_$GOFILE $GOFILE
type Event struct {
	Head
	Action   string            `codec:"action"`
	Carrier  map[string]string `codec:"carrier"`
	Payload  DataObject        `codec:"payload"`
	Callback string            `codec:"callback"`
}

type Handler func(context.Context, DataObject) (DataObject, error)

func (e *Event) Dup() (dup *Event) {
	if e == nil {
		return
	}
	dup = &Event{
		Head:     e.Head.Dup(),
		Action:   e.Action,
		Payload:  e.Payload.Dup(),
		Callback: e.Callback,
	}
	return
}

func (e *Event) Send(sock mangos.Socket) (err error) {
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
