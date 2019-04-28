package logs

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

// TracerFromContext get tracer from context or global tracer
func TracerFromContext(ctx context.Context) (tracer opentracing.Tracer) {
	if sp := opentracing.SpanFromContext(ctx); sp != nil {
		tracer = sp.Tracer()
	} else {
		tracer = opentracing.GlobalTracer()
	}
	return
}

//StartSpanFromContext -
func StartSpanFromContext(c context.Context, name string, opts ...opentracing.StartSpanOption) (sl *SpanLogger, ctx context.Context) {
	if c == nil {
		c = context.Background()
	}
	ctx = c
	tracer := TracerFromContext(ctx)
	sl, ctx = StartSpanFromContextWithTracer(ctx, tracer, name, opts...)
	return
}

// StartSpanFromContextWithTracer -
func StartSpanFromContextWithTracer(c context.Context, tracer opentracing.Tracer, name string, opts ...opentracing.StartSpanOption) (logger *SpanLogger, ctx context.Context) {
	if c == nil {
		c = context.Background()
	}
	ctx = c
	if parentSpan := SpanFromContext(ctx); parentSpan != nil {
		opts = append(opts, opentracing.ChildOf(parentSpan.Context()))
	}
	span := tracer.StartSpan(name, opts...)
	logger = NewSpanLogger(span)
	ctx = opentracing.ContextWithSpan(ctx, logger)
	return
}

func StartSpanFromCarrier(carrier map[string]string, tracer opentracing.Tracer, name string) (logger *SpanLogger, ctx context.Context) {
	ctx = WithValues(context.Background(), carrier)
	opts := []opentracing.StartSpanOption{}
	if sc, err := Extract(tracer, carrier); err == nil {
		opts = append(opts, opentracing.ChildOf(sc))
	}
	span := tracer.StartSpan(name, opts...)
	logger = NewSpanLogger(span)
	ctx = opentracing.ContextWithSpan(ctx, logger)
	return
}

func SpanFromContext(ctx context.Context) (sp *SpanLogger) {
	if span := opentracing.SpanFromContext(ctx); span != nil {
		ok := false
		if sp, ok = span.(*SpanLogger); ok {
			return
		} else {
			sp = NewSpanLogger(span)
		}
	}
	return
}
