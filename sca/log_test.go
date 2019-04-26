package sca_test

import (
	"context"
	"io/ioutil"
	"log"
	"testing"

	"github.com/nanjj/cub/sca"
)

func BenchmarkStartSpanFromContext(b *testing.B) {

	f1 := func(ctx context.Context) {
		sp, ctx := sca.StartSpanFromContext(ctx, "f1")
		sp.Println("BenchmarkStartSpanFromContext")
		defer sp.Finish()
	}

	logger := log.New(ioutil.Discard, "", log.Lshortfile)

	f2 := func(ctx context.Context) {
		logger.Println("BenchmarkStartSpanFromContext")
	}

	f3 := func(ctx context.Context) {
		sp := sca.SpanFromContext(ctx)
		sp.Println("BenchmarkStartSpanFromContext")
	}
	b.Run("Span", func(b *testing.B) {
		b.ReportAllocs()
		sp, ctx := sca.StartSpanFromContext(context.Background(), "Span")
		defer sp.Finish()
		for i := 0; i < b.N; i++ {
			f1(ctx)
		}
	})
	b.Run("StdLog", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			f2(context.Background())
		}
	})
	b.Run("SpanFromContext", func(b *testing.B) {
		b.ReportAllocs()
		var (
			ctx context.Context
			sp  *sca.SpanLogger
		)
		for i := 0; i < b.N; i++ {
			if i%90 == 0 {
				sp, ctx = sca.StartSpanFromContext(context.Background(), "SpanFromContext")
				defer sp.Finish()
			}
			f3(ctx)
		}

	})
}
