package sca

import (
	retry "github.com/avast/retry-go"
	"nanomsg.org/go/mangos/v2"
)

func RetryDial(sock mangos.Socket, addr string) (err error) {
	return retry.Do(func() error { return sock.Dial(addr) },
		retry.LastErrorOnly(true), retry.Attempts(5))
}

func RetryListen(sock mangos.Socket, addr string) (err error) {
	return retry.Do(func() error { return sock.Listen(addr) },
		retry.LastErrorOnly(true), retry.Attempts(5))
}

func RetrySend(sock mangos.Socket, b []byte) (err error) {
	return retry.Do(func() error { return sock.Send(b) },
		retry.LastErrorOnly(true), retry.Attempts(5))
}

func RetryRecv(sock mangos.Socket) (b []byte, err error) {
	err = retry.Do(func() (err error) {
		b, err = sock.Recv()
		return
	}, retry.LastErrorOnly(true), retry.Attempts(5))
	return
}
