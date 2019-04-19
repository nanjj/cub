package sca_test

import (
	"testing"

	"github.com/nanjj/cub/sca"
	"github.com/opentracing/opentracing-go"
)

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
			ctx := sca.WithValues(nil, values)
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
	err := sca.Inject(tracer, sp.Context(), carrier)
	if err != nil {
		t.Fatal()
	}
	t.Log(carrier)
	if len(carrier) != 1 {
		t.Fatal(carrier)
	}
	sc, err := sca.Extract(tracer, carrier)
	if err != nil {
		t.Fatal(err)
	}
	ctx := sca.WithValues(nil, carrier)
	sp, ctx = opentracing.StartSpanFromContextWithTracer(ctx, tracer, "TestInject", opentracing.FollowsFrom(sc))
	defer sp.Finish()
}
