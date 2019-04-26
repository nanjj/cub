package sca

import (
	"context"
	"io"
	"log"

	"github.com/opentracing/opentracing-go"
	slog "github.com/opentracing/opentracing-go/log"
	"github.com/uber/jaeger-client-go/config"
)

type SpanLogger struct {
	opentracing.Span
	*log.Logger
}

func (w *SpanLogger) Write(b []byte) (n int, err error) {
	w.LogFields(slog.String("message", string(b)))
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
func StartSpanFromContextWithTracer(c context.Context, tracer opentracing.Tracer, name string, opts ...opentracing.StartSpanOption) (sl *SpanLogger, ctx context.Context) {
	if c == nil {
		c = context.Background()
	}
	ctx = c
	if parentSpan := SpanFromContext(ctx); parentSpan != nil {
		opts = append(opts, opentracing.ChildOf(parentSpan.Context()))
	}
	sp := tracer.StartSpan(name, opts...)
	sl = &SpanLogger{Span: sp}
	logger := log.New(sl, "", log.Lshortfile)
	sl.Logger = logger
	ctx = opentracing.ContextWithSpan(ctx, sl)
	return
}

func StartSpanFromCarrier(carrier map[string]string, tracer opentracing.Tracer, name string) (sl *SpanLogger, ctx context.Context) {
	ctx = WithValues(context.Background(), carrier)
	opts := []opentracing.StartSpanOption{}
	if sc, err := Extract(tracer, carrier); err == nil {
		opts = append(opts, opentracing.ChildOf(sc))
	}
	sp := tracer.StartSpan(name, opts...)
	sl = &SpanLogger{Span: sp}
	logger := log.New(sl, "", log.Lshortfile)
	sl.Logger = logger
	ctx = opentracing.ContextWithSpan(ctx, sl)
	return
}

func SpanFromContext(ctx context.Context) (sp *SpanLogger) {
	if span := opentracing.SpanFromContext(ctx); span != nil {
		var ok bool
		if sp, ok = span.(*SpanLogger); ok {
			return
		}
	}
	return
}
