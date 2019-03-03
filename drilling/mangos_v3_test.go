package drilling

import (
	"fmt"
	"testing"
	"time"

	"encoding/json"

	"golang.org/x/sync/errgroup"
	"nanomsg.org/go/mangos/v3"
	"nanomsg.org/go/mangos/v3/protocol/pair"
	"nanomsg.org/go/mangos/v3/protocol/rep"
	"nanomsg.org/go/mangos/v3/protocol/req"
	_ "nanomsg.org/go/mangos/v3/transport/inproc"
	_ "nanomsg.org/go/mangos/v3/transport/tcp"
)

// Request/reply pattern
func TestMangosV3ReqRep(t *testing.T) {
	const (
		listen = "tcp://127.0.0.1:40899"
	)
	type Message struct {
		Index    int           `json:"index"`
		Duration time.Duration `json:"duration"`
		Input    string        `json:"input"`
		Output   string        `json:"output"`
	}
	start := func() mangos.Socket {
		var (
			b   []byte
			err error
		)
		sock, err := rep.NewSocket()
		if err != nil {
			t.Fatal(err)
		}
		if err = sock.Listen(listen); err != nil {
			t.Fatal(err)
		}
		go func() {
			for {
				b, err = sock.Recv()
				if err != nil {
					return
				}
				msg := &Message{}
				if err = json.Unmarshal(b, msg); err != nil {
					t.Fatal(err)
				}
				msg.Output = msg.Input
				msg.Index++
				if b, err = json.Marshal(msg); err != nil {
					t.Fatal(err)
				}
				if msg.Duration != 0 {
					time.Sleep(msg.Duration)
				}
				if err = sock.Send(b); err != nil {
					return
				}
			}
		}()
		return sock
	}
	server := start()
	defer func() {
		if err := server.Close(); err != nil {
			t.Log(err)
		}
	}()
	sock, err := req.NewSocket()
	if err != nil {
		t.Fatal(err)
	}
	if err = sock.Dial(listen); err != nil {
		t.Fatal(err)
	}
	var (
		b []byte
	)
	sendAndRecv := func(sock mangos.Socket, i int) (reply *Message, err error) {
		msg := &Message{Index: i, Input: "input", Duration: time.Microsecond * time.Duration(i)}
		if b, err = json.Marshal(msg); err != nil {
			return
		}
		if err = sock.Send(b); err != nil {
			return
		}
		reply = &Message{}
		if b, err = sock.Recv(); err != nil {
			return
		}
		if err = json.Unmarshal(b, reply); err != nil {
			return
		}
		return
	}
	// send and receive 10 times
	for i := 0; i < 10; i++ {
		msg, err := sendAndRecv(sock, i)
		if err != nil {
			t.Fatal(msg, err)
		}
		if msg.Output == "" || msg.Output != msg.Input || msg.Index != i+1 {
			t.Fatal(i, msg)
		}
	}

	// send 10 times first
	for i := 0; i < 10; i++ {
		msg := &Message{Index: i, Input: "input"}
		if b, err = json.Marshal(msg); err != nil {
			t.Fatal(i, err)
		}
		if err = sock.Send(b); err != nil {
			t.Fatal(i, err)
		}
	}
	{ // only the 10th message being received
		b, err = sock.Recv()
		msg := &Message{}
		if err = json.Unmarshal(b, msg); err != nil {
			t.Fatal(err)
		}
		if msg.Output == "" || msg.Output != msg.Input || msg.Index != 10 {
			t.Log(msg)
		}
	}
	if err = sock.Close(); err != nil {
		t.Fatal(err)
	}

	run := func(n int) (err error) {
		sock, err := req.NewSocket()
		if err != nil {
			return
		}
		if err = sock.Dial(listen); err != nil {
			return
		}
		for i := 0; i < 10; i++ {
			var msg *Message
			msg, err = sendAndRecv(sock, i+n*10)
			if err != nil {
				return
			}
			if msg.Input == "" || msg.Input != msg.Output || msg.Index != i+1+n*10 {
				err = fmt.Errorf("%d:%v", i, msg)
				return
			}
			t.Logf("client %d received message %v", n, msg)
		}
		err = sock.Close()
		return
	}

	rg := errgroup.Group{}
	rg.Go(func() error { return run(1) })
	rg.Go(func() error { return run(2) })
	if err = rg.Wait(); err != nil {
		t.Fatal(err)
	}
}

// pair pattern test
func TestMangosV3Pair(t *testing.T) {
	const (
		step2addr = "inproc://step2"
		step3addr = "inproc://step3"
	)
	closes := make(chan func() error, 10)
	dial := func(addr string) (sock mangos.Socket, err error) {
		if sock, err = pair.NewSocket(); err != nil {
			return
		}
		err = sock.Dial(addr)
		return
	}
	notify := func(sock mangos.Socket) (err error) {
		err = sock.Send([]byte("ready"))
		return
	}

	listen := func(addr string) (sock mangos.Socket, err error) {
		if sock, err = pair.NewSocket(); err != nil {
			return
		}
		err = sock.Listen(addr)
		return
	}

	wait := func(sock mangos.Socket) (err error) {
		var (
			b []byte
		)
		if b, err = sock.Recv(); err != nil {
			return
		}
		if "ready" != string(b) {
			err = fmt.Errorf("Not ready")
		}
		return
	}

	shutdown := func(sock mangos.Socket) (err error) {
		closes <- sock.Close
		return nil
	}

	g := errgroup.Group{}
	step1 := func() (err error) {
		t.Log("enter step 1")
		t.Log("step 1 is ready")
		sock, err := dial(step2addr)
		defer shutdown(sock)
		err = notify(sock)
		return
	}

	step2 := func() (err error) {
		t.Log("enter step2")
		sock, err := listen(step2addr)
		if err != nil {
			return
		}
		defer shutdown(sock)
		g.Go(step1)
		if err = wait(sock); err != nil {
			return err
		}
		t.Log("step2 is ready")
		sock, err = dial(step3addr)
		if err != nil {
			return
		}
		err = notify(sock)
		defer shutdown(sock)
		return
	}

	step3 := func() (err error) {
		t.Log("enter step3")
		sock, err := listen(step3addr)
		if err != nil {
			return
		}
		g.Go(step2)
		if err = wait(sock); err != nil {
			return
		}
		t.Log("step3 is ready")
		return
	}
	g.Go(step3)
	if err := g.Wait(); err != nil {
		t.Fatal(err)
	}
	for len(closes) > 0 {
		(<-closes)()
	}
}
