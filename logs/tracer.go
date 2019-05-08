package logs

import (
	"io"
	"os"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go/config"
)

func init() {
	const (
		envJaegerSamplerType           = "JAEGER_SAMPLER_TYPE"
		envJaegerSamplerParam          = "JAEGER_SAMPLER_PARAM"
		envJaegerReporterMaxQueueSize  = "JAEGER_REPORTER_MAX_QUEUE_SIZE"
		envJaegerReporterFlushInterval = "JAEGER_REPORTER_FLUSH_INTERVAL"
	)
	os.Setenv(envJaegerSamplerType, "const")
	os.Setenv(envJaegerSamplerParam, "1")
	if os.Getenv(envJaegerReporterMaxQueueSize) == "" {
		os.Setenv(envJaegerReporterMaxQueueSize, "64")
	}
	if os.Getenv(envJaegerReporterFlushInterval) == "" {
		os.Setenv(envJaegerReporterFlushInterval, "10s")
	}
}

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
