package sca_test

import (
	"os"
	"testing"

	"github.com/nanjj/cub/logs"
	"github.com/opentracing/opentracing-go"
)

func TestMain(m *testing.M) {
	rt := func() int {
		if tr, err := logs.NewTracer("ScaTest"); err == nil {
			defer tr.Close()
			opentracing.SetGlobalTracer(tr)
		}
		return m.Run()
	}()
	os.Exit(rt)
}
