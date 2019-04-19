package sca_test

import (
	"io"
	"os"
	"testing"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go/config"
)

func TestMain(m *testing.M) {
	rt := func() int {
		os.Setenv("JAEGER_SERVICE_NAME", "ScaTest")
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
