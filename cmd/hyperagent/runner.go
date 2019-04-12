package main

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"nanomsg.org/go/mangos/v2"
	"nanomsg.org/go/mangos/v2/protocol/pull"
	"nanomsg.org/go/mangos/v2/protocol/push"
	_ "nanomsg.org/go/mangos/v2/transport/all"
)

const (
	RunnerListen = "runner.listen"
	RunnerName   = "runner.name"
	RunnerIp     = "runner.ip"
	LeaderListen = "leader.listen"
)

type Handler func([]string) error

type Runner struct {
	name     string
	listen   string
	leader   mangos.Socket
	self     mangos.Socket
	members  map[string]mangos.Socket
	handlers map[string]Handler
	routes   *Routes
}

func NewRunner(name, listen, leader string) (runner *Runner, err error) {
	runner = &Runner{
		name:     name,
		listen:   listen,
		members:  map[string]mangos.Socket{},
		handlers: map[string]Handler{},
		routes:   &Routes{},
	}
	runner.AddHandler("join", runner.Join)
	runner.AddHandler("ping", runner.Ping)
	sock, err := pull.NewSocket()
	if err != nil {
		return
	}

	if err = retryListen(sock, listen); err != nil {
		return
	}
	runner.self = sock
	if leader != "" {
		if sock, err = push.NewSocket(); err != nil {
			return
		}
		if err = retryDial(sock, leader); err != nil {
			return
		}
		task := &Task{
			Name: "join",
			Args: []string{name, listen, name},
		}
		if err = task.Send(sock); err != nil {
			return
		}
		runner.leader = sock
	}
	return
}

func (runner *Runner) Run() (err error) {
	for {
		if err = runner.Handle(); err != nil {
			return
		}
	}
	return
}

func (runner *Runner) Handle() (err error) {
	task := &Task{}
	if err = task.Recv(runner.self); err != nil {
		return
	}
	local := false
	targets := task.Targets
	vias := map[string]Targets{}
	if targets.Local() {
		local = true
	} else if targets.All() {
		local = true
		for k, _ := range runner.members {
			vias[k] = targets
		}
	} else {
		vias = runner.routes.Dispatch(targets)
	}
	for k, v := range vias {
		if k == "" {
			local = true
			continue
		}
		if member, ok := runner.members[k]; ok {
			dup := task.Dup()
			dup.Targets = v
			if err = dup.Send(member); err != nil {
				return
			}
		} else {
			return
		}
	}
	if local { // handle local
		name := task.Name
		args := task.Args
		if f, ok := runner.handlers[name]; ok {
			if err = f(args); err != nil {
				log.Println(err)
				return
			}
		} else {
			log.Println("No handler found")
		}
	}
	return
}

// name,listen, members...
func (runner *Runner) Join(args []string) (err error) {
	l := len(args)
	if l < 3 {
		return fmt.Errorf("bad request")
	}
	name := args[0]
	listen := args[1]
	// add new member
	sock, ok := runner.members[name]
	if ok && listen != "" {
		sock.Close()
		sock = nil
	}
	if sock == nil {
		if sock, err = push.NewSocket(); err != nil {
			return
		}
		if err = retryDial(sock, listen); err != nil {
			return
		}
		runner.members[name] = sock
	}

	// update routes
	for i := 2; i < l; i++ {
		target := args[i]
		runner.routes.Add(target, name)
	}
	// tell leader
	if leader := runner.leader; leader != nil {
		args[0] = runner.name
		args[1] = ""
		task := &Task{
			Name: "join",
			Args: args,
		}
		if err = task.Send(leader); err != nil {
			return
		}
	}
	return
}

func (r *Runner) Members() (members []string) {
	members = []string{}
	for k, _ := range r.members {
		members = append(members, k)
	}
	return
}

func (r *Runner) Handlers() (handlers []string) {
	for k, _ := range r.handlers {
		handlers = append(handlers, k)
	}
	return
}

func (r *Runner) Routes() (routes map[string]string) {
	routes = map[string]string{}
	f := func(k, v interface{}) bool {
		if key, ok := k.(string); ok {
			if value, ok := v.(string); ok {
				routes[key] = value
			}
		}
		return true
	}
	r.routes.Range(f)
	return
}

func (r *Runner) AddHandler(name string, handler Handler) {
	r.handlers[name] = handler
}

func (r *Runner) Leader() mangos.Socket {
	return r.leader
}

func (r *Runner) Self() mangos.Socket {
	return r.self
}

func (r *Runner) Member(name string) mangos.Socket {
	if sock, ok := r.members[name]; ok {
		return sock
	}
	return nil
}

func (r *Runner) Close() (err error) {
	var errs []string
	kill := func(sock mangos.Socket) {
		if sock != nil {
			if err := sock.Close(); err != nil {
				errs = append(errs, err.Error())
			}
		}
	}
	kill(r.self)
	kill(r.leader)
	for _, sock := range r.members {
		kill(sock)
	}
	if errs != nil {
		err = errors.New(strings.Join(errs, "\n"))
	}
	return
}

func (r *Runner) Name() string { return r.name }
