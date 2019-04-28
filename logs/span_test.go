package logs_test

import (
	"context"
	"io/ioutil"
	"log"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/nanjj/cub/logs"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
)

func BenchmarkLogging(b *testing.B) {
	b.Run("ZapLogSpan", func(b *testing.B) {
		b.ReportAllocs()
		sp1, ctx := logs.StartSpanFromContext(context.Background(), b.Name())
		defer sp1.Finish()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			func(ctx context.Context) {
				sp1, ctx := logs.StartSpanFromContext(ctx, "f1")
				sp1.Info("BenchmarkStartSpanFromContext")
				defer sp1.Finish()
			}(ctx)
		}
	})

	b.Run("ZapSpanNoLog", func(b *testing.B) {
		b.ReportAllocs()
		sp1, ctx := logs.StartSpanFromContext(context.Background(), b.Name())
		defer sp1.Finish()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			func(ctx context.Context) {
				sp1, ctx := logs.StartSpanFromContext(ctx, "f1")
				defer sp1.Finish()
			}(ctx)
		}
	})

	b.Run("StdLog", func(b *testing.B) {
		b.ReportAllocs()
		logger := log.New(ioutil.Discard, "", log.LstdFlags)

		for i := 0; i < b.N; i++ {
			func(ctx context.Context) {
				logger.Println("BenchmarkStartSpanFromContext")
			}(context.Background())
		}
	})
	b.Run("ZapNop", func(b *testing.B) {
		b.ReportAllocs()
		logger := zap.NewNop()
		for i := 0; i < b.N; i++ {
			func(ctx context.Context) {
				logger.Info("BenchmarkStartSpanFromContext")
			}(context.Background())
		}
	})
	b.Run("ZapLogNoSpan", func(b *testing.B) {
		b.ReportAllocs()
		var (
			ctx context.Context
			sp1 *logs.SpanLogger
		)
		name := b.Name()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			if i%90 == 0 {
				sp1, ctx = logs.StartSpanFromContext(context.Background(), name)
				defer sp1.Finish()
			}
			func(ctx context.Context) {
				sp1 := logs.SpanFromContext(ctx)
				sp1.Info("BenchmarkStartSpanFromContext")
			}(ctx)
		}
	})
}

func TestContext(t *testing.T) {
	tcs := []struct {
		values map[string]string
	}{
		{map[string]string{"a": "b"}},
		{map[string]string{"a": "b", "c": "d"}},
		{nil},
	}

	for _, tc := range tcs {
		t.Run("", func(t *testing.T) {
			values := tc.values
			ctx := logs.WithValues(nil, values)
			for k, v := range values {
				if value, ok := ctx.Value(k).(string); !ok || value != v {
					t.Log(values)
					t.Log(ctx)
					t.Log(ok, k, v, value)
					t.Fatal()
				}
			}
		})
	}
}

func TestSpanContext(t *testing.T) {
	tracer := opentracing.GlobalTracer()
	sp1 := tracer.StartSpan("TestSpanContext")
	defer sp1.Finish()
	carrier := map[string]string{}
	err := logs.Inject(tracer, sp1.Context(), carrier)
	if err != nil {
		t.Fatal()
	}
	t.Log(carrier)
	if len(carrier) != 1 {
		t.Fatal(carrier)
	}
	sc, err := logs.Extract(tracer, carrier)
	if err != nil {
		t.Fatal(err)
	}
	ctx := logs.WithValues(nil, carrier)
	sp1, ctx = opentracing.StartSpanFromContextWithTracer(ctx, tracer, "TestInject", opentracing.FollowsFrom(sc))
	defer sp1.Finish()
}

func TestTracerFromContext(t *testing.T) {
	global := opentracing.GlobalTracer()
	defer func() { opentracing.SetGlobalTracer(global) }()
	ctx := context.Background()
	tr := logs.TracerFromContext(ctx)
	if tr != global {
		t.Fatal(global, tr)
	}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	sp1 := NewMockSpan(ctrl)
	sp1.EXPECT().Tracer().Return(NewMockTracer(ctrl)).Times(1)
	sp1.EXPECT().Finish().Times(1)
	defer sp1.Finish()
	ctx = opentracing.ContextWithSpan(ctx, sp1)
	tr = logs.TracerFromContext(ctx)
	if tr == global {
		t.Fatal(global, tr)
	}
}

func TestSpanFromContext(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	sp1 := NewMockSpan(ctrl)
	ctx := context.Background()
	ctx = opentracing.ContextWithSpan(ctx, sp1)
	span := logs.SpanFromContext(ctx)
	if span.Span != sp1 {
		t.Fatal(sp1, span)
	}
	ctx = opentracing.ContextWithSpan(ctx, span)
	if logs.SpanFromContext(ctx) != span {
		t.Fatal(span)
	}
}

func TestStartSpanFromContext(t *testing.T) {
	global := opentracing.GlobalTracer()
	defer func() { opentracing.SetGlobalTracer(global) }()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	tr := NewMockTracer(ctrl)
	opentracing.SetGlobalTracer(tr)
	sp1 := NewMockSpan(ctrl)
	sp2 := NewMockSpan(ctrl)
	sp3 := NewMockSpan(ctrl)
	sp1c := NewMockSpanContext(ctrl)

	gomock.InOrder(
		tr.EXPECT().StartSpan(t.Name()).Return(sp1).Times(1),
		sp1.EXPECT().Tracer().Return(tr).Times(1),
		sp1.EXPECT().Context().Return(sp1c).Times(1),
		tr.EXPECT().StartSpan("Child", gomock.Any()).Return(sp2).Times(1),
		tr.EXPECT().StartSpan("Root").Return(sp3).Times(1),
		sp3.EXPECT().Finish().Times(1),
		sp2.EXPECT().Finish().Times(1),
		sp1.EXPECT().Finish().Times(1),
	)
	span, ctx := logs.StartSpanFromContext(nil, t.Name())
	if span.Span != sp1 || ctx == nil {
		t.Fatal(sp1, span, ctx)
	}
	defer span.Finish()
	span, ctx = logs.StartSpanFromContext(ctx, "Child")
	defer span.Finish()
	span, ctx = logs.StartSpanFromContextWithTracer(nil, tr, "Root")
	defer span.Finish()
}

func TestStartSpanFromCarrier(t *testing.T) {
	global := opentracing.GlobalTracer()
	defer func() {
		opentracing.SetGlobalTracer(global)
	}()
	carrier := map[string]string{
		"uber-trace-id": "4e5fda89553c4ebe:4e5fda89553c4ebe:0:1",
	}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	tr := NewMockTracer(ctrl)
	sp1 := NewMockSpan(ctrl)
	sp1c := NewMockSpanContext(ctrl)
	gomock.InOrder(
		tr.EXPECT().Extract(opentracing.TextMap,
			opentracing.TextMapCarrier(carrier)).Return(
			sp1c, nil).Times(1),
		tr.EXPECT().StartSpan(t.Name(), gomock.Any()).Return(sp1).Times(1),
		sp1.EXPECT().Finish().Times(1),
	)
	span, ctx := logs.StartSpanFromCarrier(carrier, tr, t.Name())
	defer span.Finish()
	if s, ok := ctx.Value("uber-trace-id").(string); !ok || s != "4e5fda89553c4ebe:4e5fda89553c4ebe:0:1" {
		t.Fatal(s, ok)
	}
	if span.Span != sp1 {
		t.Fatal(span, sp1)
	}
}
