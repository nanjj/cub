package logs_test

import (
	"context"
	"io/ioutil"
	"log"
	"testing"

	"github.com/nanjj/cub/logs"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
)

func BenchmarkLogging(b *testing.B) {
	b.Run("ZapLogSpan", func(b *testing.B) {
		b.ReportAllocs()
		sp, ctx := logs.StartSpanFromContext(context.Background(), b.Name())
		defer sp.Finish()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			func(ctx context.Context) {
				sp, ctx := logs.StartSpanFromContext(ctx, "f1")
				sp.Info("BenchmarkStartSpanFromContext")
				defer sp.Finish()
			}(ctx)
		}
	})

	b.Run("ZapSpanNoLog", func(b *testing.B) {
		b.ReportAllocs()
		sp, ctx := logs.StartSpanFromContext(context.Background(), b.Name())
		defer sp.Finish()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			func(ctx context.Context) {
				sp, ctx := logs.StartSpanFromContext(ctx, "f1")
				defer sp.Finish()
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
			sp  *logs.SpanLogger
		)
		name := b.Name()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			if i%90 == 0 {
				sp, ctx = logs.StartSpanFromContext(context.Background(), name)
				defer sp.Finish()
			}
			func(ctx context.Context) {
				sp := logs.SpanFromContext(ctx)
				sp.Info("BenchmarkStartSpanFromContext")
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
	sp := tracer.StartSpan("TestSpanContext")
	defer sp.Finish()
	carrier := map[string]string{}
	err := logs.Inject(tracer, sp.Context(), carrier)
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
	sp, ctx = opentracing.StartSpanFromContextWithTracer(ctx, tracer, "TestInject", opentracing.FollowsFrom(sc))
	defer sp.Finish()
}
