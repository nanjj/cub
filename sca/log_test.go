package sca_test

import (
	"context"
	"log"
	"testing"

	"github.com/nanjj/cub/sca"
)

func BenchmarkStartSpanFromContext(b *testing.B) {
	f1 := func(ctx context.Context) {
		sp, ctx := sca.StartSpanFromContext(ctx, "BenchmarkStartSpanFromContext")
		sp.Println("BenchmarkStartSpanFromContext")
		defer sp.Finish()
	}
	f2 := func(ctx context.Context) {
		log.Println("BenchmarkStartSpanFromContext")
	}

	b.Run("Span", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			f1(context.Background())
		}
	})
	b.Run("Log", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			f2(context.Background())
		}
	})
}
