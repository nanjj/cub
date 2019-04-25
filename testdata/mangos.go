package drilling

import (
	"fmt"
	"os"
)

var (
	defers  = make(chan func() error, 1024)
	listens = make(chan bool, 1024)
)

func Do() {
	for len(defers) > 0 {
		(<-defers)()
	}
}

func Defer(f func() error) {
	defers <- f
}

func Add() {
	listens <- true
}

func Wait(n int) {
	for i := 0; i < n; i++ {
		<-listens
	}
}

func tLog(a ...interface{}) {
	fmt.Fprintln(os.Stderr, a...)
}

func tLogf(format string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, format, a...)
}
