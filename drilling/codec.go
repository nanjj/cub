package drilling

import (
	"time"
)

//go:generate codecgen -o codec_gen.go codec.go
type TestingEvent struct {
	Id        int       `json:"id"`
	Kind      int       `json:"kind"`
	CreatedAt time.Time `json:"created_at"`
	Target    string    `json:"target"`
	Source    string    `json:"source"`
}
