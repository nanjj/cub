package logs_test

import (
	"testing"

	"github.com/nanjj/cub/logs"
	"github.com/uber/jaeger-client-go/config"
)

func TestNewTracer(t *testing.T) {
	tr, err := logs.NewTracer("tracer",
		config.Tag("runner", "127.0.0.1:54321"),
		config.Tag("leader", "127.0.0.1:54312"))
	if err != nil {
		t.Fatal()
	}
	defer tr.Close()
	sp := tr.StartSpan("TestNewTracer")
	sp.Finish()
}
