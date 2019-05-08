package sca_test

import (
	"os"
	"testing"
	"time"

	"github.com/nanjj/cub/logs"
	"github.com/opentracing/opentracing-go"
)

func TestMain(m *testing.M) {
	os.Exit(func() int {
		if tr, err := logs.NewTracer("ScaTest"); err == nil {
			defer time.Sleep(time.Millisecond)
			defer tr.Close()
			opentracing.SetGlobalTracer(tr)
		}
		return m.Run()
	}())
}
