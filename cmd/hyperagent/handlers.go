package main

import (
	"log"
)

func (r *Runner) Ping(args []string) error {
	log.Println("ping", r.Name(), args)
	return nil
}
