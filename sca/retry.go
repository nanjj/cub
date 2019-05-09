package sca

import (
	retry "github.com/avast/retry-go"
	"nanomsg.org/go/mangos/v2"
)

var (
	Retry = RetryOpts{retry.LastErrorOnly(true), retry.Attempts(5)}
)

type RetryOpts []retry.Option

func (opts RetryOpts) Dial(sock mangos.Socket, addr string) (err error) {
	return retry.Do(func() error { return sock.Dial(addr) }, opts...)
}

func (opts RetryOpts) Listen(sock mangos.Socket, addr string) (err error) {
	return retry.Do(func() error { return sock.Listen(addr) }, opts...)
}

func (opts RetryOpts) Send(sock mangos.Socket, b []byte) (err error) {
	return retry.Do(func() error { return sock.Send(b) }, opts...)
}

func (opts RetryOpts) Recv(sock mangos.Socket) (b []byte, err error) {
	err = retry.Do(func() (err error) {
		b, err = sock.Recv()
		return
	}, opts...)
	return
}
