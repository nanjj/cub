package sca

import (
	"testing"

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
			ctx := WithValues(nil, values)
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

func TestExtract(t *testing.T) {
	tracer := opentracing.GlobalTracer()
	sp := tracer.StartSpan("TestExtract")
	defer sp.Finish()
	carrier, err := Inject(tracer, sp.Context())
	if err != nil {
		t.Fatal()
	}
	t.Log(carrier)
	if len(carrier) != 1 {
		t.Fatal(carrier)
	}
	sc, err := Extract(tracer, carrier)
	if err != nil {
		t.Fatal(err)
	}
	ctx := WithValues(nil, carrier)
	sp, ctx = opentracing.StartSpanFromContextWithTracer(ctx, tracer, "TestInject", opentracing.FollowsFrom(sc))
	defer sp.Finish()
}
