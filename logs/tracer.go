package logs

import (
	"io"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go/config"
)

// NewTracer - returns a new tracer
func NewTracer(name string, opts ...config.Option) (tracer opentracing.Tracer, closer io.Closer, err error) {
	c, err := config.FromEnv()
	if err != nil {
		return
	}
	if name != "" {
		c.ServiceName = name
	}
	tracer, closer, err = c.NewTracer(opts...)
	return
}
