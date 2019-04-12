package tasks

import (
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

//go:generate codecgen -o cg_$GOFILE $GOFILE
type Task struct {
	Targets Targets  `codec:"targets"`
	Id      int64    `codec:"id"`
	Name    string   `codec:"name"`
	Args    []string `codec:"args"`
}

func (task *Task) Dup() (dup *Task) {
	dup = &Task{
		Id:   task.Id,
		Name: task.Name,
		Args: task.Args,
	}
	dup.Targets = append(dup.Targets, task.Targets...)
	return
}

func (task *Task) Send(sock mangos.Socket) (err error) {
	out := make([]byte, 0, 1024)
	enc := codec.NewEncoderBytes(&out, cbor)
	if err = enc.Encode(task); err != nil {
		log.Println(err)
		return
	}
	if err = RetrySend(sock, out); err != nil {
		log.Println(err)
		return
	}
	return
}

func (task *Task) Recv(sock mangos.Socket) (err error) {
	in, err := RetryRecv(sock)
	if err != nil {
		log.Println(err)
		return
	}
	dec := codec.NewDecoderBytes(in, cbor)
	if err = dec.Decode(task); err != nil {
		log.Println(err)
		return
	}
	return
}
