package drilling

import (
	"time"
)

//go:generate codecgen -o cg_$GOFILE $GOFILE
type TestingEvent struct {
	Id        int       `codec:"id"`
	Kind      int       `codec:"kind"`
	CreatedAt time.Time `codec:"created_at"`
	Target    string    `codec:"target"`
	Source    string    `codec:"source"`
}

type TestingHead struct {
	Targets  []int64 `codec:"targets"`
	Callback int64   `codec:"callback"`
}

type TestingBody struct {
	Action  string
	Command string
	Args    []string
}

type TestingHeadBody struct {
	TestingHead
	TestingBody
}
