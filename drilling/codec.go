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
