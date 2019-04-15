package drilling

import (
	"context"
	"io"
	"os"
	"testing"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go/config"
)

func TestMain(m *testing.M) {
	rt := func() int {
		os.Setenv("JAEGER_SERVICE_NAME", "HyperagentTest")
		os.Setenv("JAEGER_SAMPLER_TYPE", "const")
		os.Setenv("JAEGER_SAMPLER_PARAM", "1")
		os.Setenv("JAEGER_REPORTER_MAX_QUEUE_SIZE", "64")
		os.Setenv("JAEGER_REPORTER_FLUSH_INTERVAL", "10s")
		os.Setenv("JAEGER_TAGS", "runner=r1")
		var (
			cfg    *config.Configuration
			err    error
			tr     opentracing.Tracer
			closer io.Closer
		)

		if cfg, err = config.FromEnv(); err != nil {
			panic(err)
		}

		if tr, closer, err = cfg.NewTracer(); err != nil {
			panic(err)
		}
		defer closer.Close()
		opentracing.SetGlobalTracer(tr)
		return m.Run()
	}()
	os.Exit(rt)
}

func TestStartSpan(t *testing.T) {
	sp0 := opentracing.StartSpan("TestStartSpan")
	defer sp0.Finish()
	sp0.LogKV("name", "sp0")
}

func TestStartSpanFromContext(t *testing.T) {
	ctx := context.Background()
	sp1, ctx := opentracing.StartSpanFromContext(ctx, "sp1")
	defer sp1.Finish()
	sp1.LogKV("name", "sp1")
	sp2, ctx := opentracing.StartSpanFromContext(ctx, "sp2")
	sp2.LogKV("name", "sp2")
	defer sp2.Finish()
}

func TestSpanContext(t *testing.T) {
	ctx1 := context.Background()
	sp1, ctx1 := opentracing.StartSpanFromContext(ctx1, "client")
	spctx := sp1.Context()
	tr1 := sp1.Tracer()
	//carrier := opentracing.HTTPHeadersCarrier{}
	carrier := opentracing.TextMapCarrier{}
	if err := tr1.Inject(spctx, opentracing.TextMap, &carrier); err != nil {
		t.Fatal(err)
	}
	t.Log(carrier, spctx)
	// spctx.ForeachBaggageItem(func(k, v string) bool {
	// 	t.Log(k, v)
	// 	return true
	// })
	sp1.Finish()
	ctx2 := context.Background()
	tr2 := opentracing.GlobalTracer()
	sm, err := tr2.Extract(opentracing.TextMap, carrier)
	if err != nil {
		t.Fatal(err)
	}
	sp2, _ := opentracing.StartSpanFromContext(ctx2, "server", opentracing.ChildOf(sm))
	sp2.Finish()
	sp3, _ := opentracing.StartSpanFromContext(ctx2, "server2", opentracing.FollowsFrom(sm))
	sp3.Finish()
}

func TestDefer(t *testing.T) {
	defer t.Log("1")
	defer t.Log("2")
	defer t.Log("3")
}
