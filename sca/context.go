package sca

import (
	"context"

	"github.com/opentracing/opentracing-go"
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
