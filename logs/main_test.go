package logs_test

import (
	"os"
	"testing"
	"time"

	"github.com/nanjj/cub/logs"
	"github.com/opentracing/opentracing-go"
)

func TestMain(m *testing.M) {
	os.Exit(func() int {
		if tr, err := logs.NewTracer("LogsTest"); err == nil {
			opentracing.SetGlobalTracer(tr)
			defer time.Sleep(time.Millisecond)
			defer tr.Close()
		}
		return m.Run()
	}())
}
