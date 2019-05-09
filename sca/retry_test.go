package sca_test

import (
	"io"
	"testing"
	"time"

	"github.com/avast/retry-go"
	"github.com/nanjj/cub/sca"
	"golang.org/x/sync/errgroup"
	"nanomsg.org/go/mangos/v2"
	"nanomsg.org/go/mangos/v2/protocol/pull"
	"nanomsg.org/go/mangos/v2/protocol/push"
)

func TestRetry(t *testing.T) {
	const (
		addr = "tcp://127.0.0.1:55555"
	)
	oldRetry := sca.Retry
	defer func() { sca.Retry = oldRetry }()
	sca.Retry = sca.RetryOpts{retry.Attempts(10), retry.Delay(time.Millisecond * 10)}
	tcs := []struct {
		dial   func(mangos.Socket, string) error
		listen func(mangos.Socket, string) error
		send   func(mangos.Socket, []byte) error
		recv   func(mangos.Socket) ([]byte, error)
		noerr  bool
	}{
		{
			dial:   sca.Retry.Dial,
			listen: sca.Retry.Listen,
			send:   sca.Retry.Send,
			recv:   sca.Retry.Recv,
			noerr:  true,
		},
		{
			dial:   func(sock mangos.Socket, listen string) error { return sock.Dial(listen) },
			listen: func(sock mangos.Socket, listen string) error { return sock.Listen(listen) },
			send:   func(sock mangos.Socket, b []byte) error { return sock.Send(b) },
			recv:   func(sock mangos.Socket) ([]byte, error) { return sock.Recv() },
			noerr:  false,
		},
	}
	for _, tc := range tcs {
		t.Run("", func(t *testing.T) {

			var (
				closers = make(chan io.Closer, 1024)
				results = make(chan string, 1024)
				g       = errgroup.Group{}
			)
			defer func() {
				l := len(closers)
				for i := 0; i < l; i++ {
					closer := <-closers
					closer.Close()
				}
			}()

			// Dial and Send
			g.Go(func() (err error) {
				var sock mangos.Socket
				sock, err = push.NewSocket()
				if err != nil {
					t.Log(err)
					return
				}
				closers <- sock
				err = tc.dial(sock, addr)
				if err != nil {
					t.Log(err)
					return
				}
				g.Go(func() (err error) {
					err = tc.send(sock, []byte("hello"))
					if err != nil {
						t.Log(err)
					}
					return
				})
				return
			})
			time.Sleep(time.Millisecond * 10)
			//Listen and Recv
			g.Go(func() (err error) {
				var sock mangos.Socket
				sock, err = pull.NewSocket()
				if err != nil {
					t.Log(err)
					return
				}
				sock.SetOption(mangos.OptionRecvDeadline, time.Millisecond*10)
				closers <- sock
				err = tc.listen(sock, addr)
				if err != nil {
					t.Log(err)
					return
				}
				g.Go(func() (err error) {
					var b []byte
					b, err = tc.recv(sock)
					if err != nil {
						t.Log(err)
						return
					}
					results <- string(b)
					return
				})
				return
			})
			err := g.Wait()
			if tc.noerr != (err == nil) {
				t.Fatal(err)
			}
			if len(results) != 0 {
				if result := <-results; result != "hello" {
					t.Fatal(result)
				}
			}
		})
	}
}
