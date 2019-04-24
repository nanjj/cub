package drilling

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go/config"
	"github.com/uber/jaeger-client-go/thrift-gen/jaeger"
)

func init() {
	os.Setenv("JAEGER_SAMPLER_TYPE", "const")
	os.Setenv("JAEGER_SAMPLER_PARAM", "1")
	os.Setenv("JAEGER_REPORTER_MAX_QUEUE_SIZE", "64")
	os.Setenv("JAEGER_REPORTER_FLUSH_INTERVAL", "10s")
	os.Setenv("JAEGER_TAGS", "runner=r1")
}

func NewTracer(name string, options ...string) (tr opentracing.Tracer, closer io.Closer, err error) {
	var (
		cfg *config.Configuration
	)
	if cfg, err = config.FromEnv(); err != nil {
		return
	}
	cfg.ServiceName = name
	if len(options) != 0 {
		cfg.Reporter.LocalAgentHostPort = options[0]
	}
	if tr, closer, err = cfg.NewTracer(); err != nil {
		return
	}
	return
}

func TestMain(m *testing.M) {
	rt := func() int {
		tr, closer, err := NewTracer("TestGlobal")
		if err != nil {
			panic(err)
		}
		defer closer.Close()
		opentracing.SetGlobalTracer(tr)
		return m.Run()
	}()
	os.Exit(rt)
}

func TestRootSpan(t *testing.T) {
	rootSpan := opentracing.StartSpan("RootSpan")
	defer rootSpan.Finish()
}

func TestChildOfSpan(t *testing.T) {
	ctx := context.Background()
	spanA, ctx := opentracing.StartSpanFromContext(ctx, "A")
	defer spanA.Finish()
	spanB, ctxB := opentracing.StartSpanFromContext(ctx, "B")
	defer spanB.Finish()
	spanD, _ := opentracing.StartSpanFromContext(ctxB, "D")
	defer spanD.Finish()
	spanC, ctxC := opentracing.StartSpanFromContext(ctx, "C")
	defer spanC.Finish()
	spanE, _ := opentracing.StartSpanFromContext(ctxC, "E")
	defer spanE.Finish()
	spanF, _ := opentracing.StartSpanFromContext(ctxC, "F")
	defer spanF.Finish()
}

func TestFollowsFromSpan(t *testing.T) {
	ChildOfSpan := opentracing.StartSpanFromContext
	FollowsFrom := func(c context.Context, name string) (sp opentracing.Span, ctx context.Context) {
		sp = opentracing.SpanFromContext(c)
		tr := sp.Tracer()
		sp = tr.StartSpan(name, opentracing.FollowsFrom(sp.Context()))
		ctx = opentracing.ContextWithSpan(c, sp)
		return
	}

	h := func(ctx context.Context) {
		sp, ctx := FollowsFrom(ctx, "H")
		defer sp.Finish()
	}

	g := func(ctx context.Context) {
		sp, ctx := FollowsFrom(ctx, "G")
		defer sp.Finish()
		go h(ctx)
	}

	f := func(ctx context.Context) {
		sp, ctx := ChildOfSpan(ctx, "F")
		defer sp.Finish()
		go g(ctx)
	}
	e := func(ctx context.Context) {
		sp, ctx := ChildOfSpan(ctx, "E")
		defer sp.Finish()
	}
	d := func(ctx context.Context) {
		sp, ctx := ChildOfSpan(ctx, "D")
		defer sp.Finish()
	}
	c := func(ctx context.Context) {
		sp, ctx := ChildOfSpan(ctx, "C")
		defer sp.Finish()
		e(ctx)
		f(ctx)
	}
	b := func(ctx context.Context) {
		sp, ctx := ChildOfSpan(ctx, "B")
		defer sp.Finish()
		d(ctx)
	}
	a := func(ctx context.Context) {
		sp, ctx := ChildOfSpan(ctx, "A")
		defer sp.Finish()
		b(ctx)
		c(ctx)
	}
	a(context.Background())
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
	t.Skip()
	sp, _ := opentracing.StartSpanFromContext(context.Background(), "RootSpan")
	defer sp.Finish()
	t.Log(sp.Context())
	traceID := strings.Split(fmt.Sprintf("%v", sp.Context()), ":")[0]
	err := Open(fmt.Sprintf("http://127.0.0.1:16686/trace/%s", traceID))
	if err != nil {
		t.Fatal(err)
	}
}

func TestSpanContextInject(t *testing.T) {
	sp, _ := opentracing.StartSpanFromContext(context.Background(), "A")
	defer sp.Finish()
	tr := sp.Tracer()
	carrier := opentracing.TextMapCarrier{}
	if err := tr.Inject(sp.Context(), opentracing.TextMap, carrier); err != nil {
		t.Fatal(err)
	}
	t.Log(carrier)
}

func TestSpanContextExtract(t *testing.T) {
	carrier := opentracing.TextMapCarrier{
		"uber-trace-id": "650e10af6420a88:650e10af6420a88:0:1",
	}
	tr := opentracing.GlobalTracer()
	spanContext, err := tr.Extract(opentracing.TextMap, carrier)
	if err != nil {
		t.Fatal(err)
	}
	sp := tr.StartSpan("B", opentracing.FollowsFrom(spanContext))
	defer sp.Finish()
	t.Log(sp.Context())
}

func TestSpanLogKV(t *testing.T) {
	sp, _ := opentracing.StartSpanFromContext(context.Background(), "RootSpan")
	defer sp.Finish()
	sp.LogKV("message", "debugging")
	// http://127.0.0.1:16686/trace/3edbcc362b5a07bf
}

func TestSpanTags(t *testing.T) {
	sp, _ := opentracing.StartSpanFromContext(context.Background(), "RootSpan")
	sp.SetTag("error", true)
	defer sp.Finish()
}

func TestThriftOverUDP(t *testing.T) {
	closers := make(chan io.Closer, 1024)
	defer func() {
		for i := 0; i < len(closers); i++ {
			c := <-closers
			c.Close()
		}
	}()
	ch := make(chan []byte, 1024)
	addr := "127.0.0.1:6381"
	AgentProxy(addr, ch, closers)
	tr, closer, err := NewTracer("AgentProxy", addr)
	if err != nil {
		t.Fatal(err)
	}
	closers <- closer
	// Set Global Tracer
	opentracing.SetGlobalTracer(tr)
	ctx := context.Background()
	sp, ctx := StartSpanFromContext(ctx, "TestThriftOverUDP")
	sp.Println("something wrong")
	sp.MarkError()
	sp.Finish()    // finish span
	closer.Close() // close tracer
	// Now get thrifts
	b := <-ch
	emitBatch := &EmitBatch{}
	if err := emitBatch.Decode(b); err != nil {
		t.Fatal(err)
	}
	t.Log(emitBatch)
}

func TestThiftModifyAndResent(t *testing.T) {
	data := []byte(`
{
  "name": "emitBatch",
  "seqid": 1,
  "typeid": 4,
  "args": {
    "batch": {
      "process": {
        "serviceName": "AgentProxy",
        "tags": [
          {
            "key": "runner",
            "vType": "STRING",
            "vStr": "r1"
          },
          {
            "key": "jaeger.version",
            "vType": "STRING",
            "vStr": "Go-2.16.0"
          },
          {
            "key": "hostname",
            "vType": "STRING",
            "vStr": "nanjj.cn.ibm.com"
          },
          {
            "key": "ip",
            "vType": "STRING",
            "vStr": "9.119.157.182"
          },
          {
            "key": "client-uuid",
            "vType": "STRING",
            "vStr": "438da1dfd4090f87"
          }
        ]
      },
      "spans": [
        {
          "traceIdLow": 8235869911999135402,
          "traceIdHigh": 0,
          "spanId": 8235869911999135402,
          "parentSpanId": 0,
          "operationName": "TestThriftOverUDP",
          "flags": 1,
          "startTime": 1556005904696354,
          "duration": 31,
          "tags": [
            {
              "key": "sampler.type",
              "vType": "STRING",
              "vStr": "const"
            },
            {
              "key": "sampler.param",
              "vType": "BOOL",
              "vBool": true
            },
            {
              "key": "error",
              "vType": "BOOL",
              "vBool": true
            }
          ],
          "logs": [
            {
              "timestamp": 1556005904696383,
              "fields": [
                {
                  "key": "message",
                  "vType": "STRING",
                  "vStr": "jaeger_test.go:214: something wrong"
                }
              ]
            }
          ]
        }
      ]
    }
  }
}`)
	emitBatch := &EmitBatch{}
	if err := json.Unmarshal(data, emitBatch); err != nil {
		t.Fatal(err)
	}
	tags := emitBatch.Args.Batch.Spans[0].Tags
	tenantId := "eadfea84-659d-11e9-b918-5babd48eca8c"
	tags = append(tags, &jaeger.Tag{
		Key:   "tenantId",
		VStr:  &tenantId,
		VType: jaeger.TagType_STRING})
	emitBatch.Args.Batch.Spans[0].Tags = tags
	b, err := emitBatch.Encode()
	if err != nil {
		t.Fatal(err)
	}
	addr := "127.0.0.1:6831"
	writer, err := NewUDPWriter(addr)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := writer.Write(b); err != nil {
		t.Fatal(err)
	}
	writer.Close()
	// trace: http://127.0.0.1:16686/trace/724bb00ca9093eaa
}

func TestSimulateNoRootSpan(t *testing.T) {
	sp, ctx := StartSpanFromContext(context.Background(), "RootSpan")
	sp, ctx = StartSpanFromContext(ctx, "ChildSpan")
	t.Log(sp.Context())
	defer sp.Finish()
}

func TestDefer(t *testing.T) {
	defer t.Log("1")
	defer t.Log("2")
	defer t.Log("3")
}
