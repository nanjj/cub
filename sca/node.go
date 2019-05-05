package sca

import (
	"nanomsg.org/go/mangos/v2"
)

type Node struct {
	Name   string
	Listen string
	Server mangos.Socket
	Client mangos.Socket
}
