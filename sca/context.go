package sca

import (
	"context"
	"io"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go/config"
)

func WithValues(c context.Context,
	values map[string]string) (ctx context.Context) {
	if c == nil {
		ctx = context.Background()
	}
	for k, v := range values {
		ctx = context.WithValue(ctx, k, v)
	}
	return
}

func Extract(tracer opentracing.Tracer,
	carrier map[string]string) (sc opentracing.SpanContext, err error) {
	sc, err = tracer.Extract(opentracing.TextMap,
		opentracing.TextMapCarrier(carrier))
	return
}

func Inject(tracer opentracing.Tracer,
	sc opentracing.SpanContext, carrier map[string]string) (err error) {
	err = tracer.Inject(sc, opentracing.TextMap,
		opentracing.TextMapCarrier(carrier))
	return
}

func TracerFromContext(ctx context.Context) (tracer opentracing.Tracer) {
	if sp := opentracing.SpanFromContext(ctx); sp != nil {
		tracer = sp.Tracer()
	} else {
		tracer = opentracing.GlobalTracer()
	}
	return
}

func StartSpanFromContext(c context.Context, name string) (sp opentracing.Span, ctx context.Context) {
	if c == nil {
		c = context.Background()
	}
	ctx = c
	tracer := TracerFromContext(ctx)
	sp, ctx = opentracing.StartSpanFromContextWithTracer(ctx, tracer, name)
	return
}

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
