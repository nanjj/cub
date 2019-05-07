package logs

import (
	"io"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go/config"
)

type Tracer struct {
	opentracing.Tracer
	io.Closer
}

// NewTracer - returns a new tracer
func NewTracer(name string, opts ...config.Option) (tracer *Tracer, err error) {
	c, err := config.FromEnv()
	if err != nil {
		return
	}
	if name != "" {
		c.ServiceName = name
	}
	tr, closer, err := c.NewTracer(opts...)
	if err == nil {
		tracer = &Tracer{
			Tracer: tr,
			Closer: closer,
		}
	}
	return
}
