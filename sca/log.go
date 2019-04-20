package sca

import (
	"context"
	"io"
	"log"
	"os"
	"strings"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go/config"
)

type spanWritter struct {
	span opentracing.Span
}

type SpanLogger struct {
	opentracing.Span
	*log.Logger
}

func (w *spanWritter) Write(b []byte) (n int, err error) {
	if sp := w.span; sp != nil {
		s := strings.TrimSpace(string(b))
		w.span.LogKV("message", s)
		n = len(b)
	} else {
		n, err = os.Stderr.Write(b)
	}
	return
}

// TracerFromContext get tracer from context or global tracer
func TracerFromContext(ctx context.Context) (tracer opentracing.Tracer) {
	if sp := opentracing.SpanFromContext(ctx); sp != nil {
		tracer = sp.Tracer()
	} else {
		tracer = opentracing.GlobalTracer()
	}
	return
}

// NewTracer - returns a new tracer
func NewTracer(name, runner, leader string) (tracer opentracing.Tracer, closer io.Closer, err error) {
	c, err := config.FromEnv()
	if err != nil {
		return
	}
	if name != "" {
		c.ServiceName = name
	}
	if runner != "" && leader != "" {
		tags := map[string]string{
			"runner": runner,
			"leader": leader,
		}
		for k, v := range tags {
			c.Tags = append(c.Tags, opentracing.Tag{k, v})
		}
	}
	tracer, closer, err = c.NewTracer()
	return

}

//StartSpanFromContext -
func StartSpanFromContext(c context.Context, name string) (sl *SpanLogger, ctx context.Context) {
	if c == nil {
		c = context.Background()
	}
	ctx = c
	tracer := TracerFromContext(ctx)
	sl, ctx = StartSpanFromContextWithTracer(ctx, tracer, name)
	return
}

// StartSpanFromContextWithTracer -
func StartSpanFromContextWithTracer(c context.Context, tracer opentracing.Tracer, name string) (sl *SpanLogger, ctx context.Context) {
	if c == nil {
		c = context.Background()
	}
	ctx = c
	sp, ctx := opentracing.StartSpanFromContextWithTracer(ctx, tracer, name)
	logger := log.New(&spanWritter{sp}, "", log.Lshortfile)
	sl = &SpanLogger{sp, logger}
	return
}
