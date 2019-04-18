package sca

import (
	"log"

	"github.com/ugorji/go/codec"
	"nanomsg.org/go/mangos/v2"
)

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
