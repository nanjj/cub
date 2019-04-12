package main

import (
	retry "github.com/avast/retry-go"
	"nanomsg.org/go/mangos/v2"
)

func retryDial(sock mangos.Socket, addr string) (err error) {
	return retry.Do(func() error { return sock.Dial(addr) })
}

func retryListen(sock mangos.Socket, addr string) (err error) {
	return retry.Do(func() error { return sock.Listen(addr) })
}

func retrySend(sock mangos.Socket, b []byte) (err error) {
	return retry.Do(func() error { return sock.Send(b) })
}

func retryRecv(sock mangos.Socket) (b []byte, err error) {
	err = retry.Do(func() (err error) {
		b, err = sock.Recv()
		return
	})
	return
}
