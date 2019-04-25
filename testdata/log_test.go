package drilling

import (
	"context"
	"log"
	"strings"
	"testing"

	"github.com/opentracing/opentracing-go"
)

type SpanWritter struct {
	span opentracing.Span
}

func (w *SpanWritter) Write(b []byte) (n int, err error) {
	s := strings.TrimSpace(string(b))
	w.span.LogKV("message", s)
	n = len(b)
	return
}

type SpanLogger struct {
	opentracing.Span
	*log.Logger
}

func (sp SpanLogger) MarkError() {
	sp.SetTag("error", true)
}

func TracerFromContext(ctx context.Context) (tracer opentracing.Tracer) {
	if sp := opentracing.SpanFromContext(ctx); sp != nil {
		tracer = sp.Tracer()
	} else {
		tracer = opentracing.GlobalTracer()
	}
	return
}

func StartSpanFromContext(c context.Context, name string) (sl *SpanLogger, ctx context.Context) {
	if c == nil {
		c = context.Background()
	}
	ctx = c
	tracer := TracerFromContext(ctx)
	sp, ctx := opentracing.StartSpanFromContextWithTracer(ctx, tracer, name)
	logger := log.New(&SpanWritter{sp}, "", log.Lshortfile)
	sl = &SpanLogger{sp, logger}
	return
}

func TestStartSpanFromContextLogger(t *testing.T) {
	logger, ctx := StartSpanFromContext(context.Background(), "TestStartSpanFromContextLogger")
	logger.Println("hello")
	logger.Println("something wrong")
	defer logger.Finish()
	logger2, ctx := StartSpanFromContext(ctx, "TestStartSpanFromContextLogger2")
	defer logger2.Finish()
	logger2.Println("why?")
}

func TestSpanStdLog(t *testing.T) {
	sp, _ := StartSpanFromContext(context.Background(), "TestSpanStdLog")
	defer sp.Finish()
	sp.Println("invalid parameter")
	sp.MarkError()
	// Trace: http://127.0.0.1:16686/trace/3ce2a393b55083cb
}
