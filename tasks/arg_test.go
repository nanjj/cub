package tasks

import (
	"testing"

	"github.com/ugorji/go/codec"
)

func TestArgsString(t *testing.T) {
	s := "hello"
	out := make([]byte, 0, 1024)
	enc := codec.NewEncoderBytes(&out, cbor)
	if err := enc.Encode(&s); err != nil {
		t.Fatal(err)
	}
	t.Log(len(out))
	t.Log(string(out))
}

func BenchmarkArgs(b *testing.B) {
	out := make([]byte, 0, 8)
	s := "hello"
	b.Run("Arg", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			if string(Arg(s)) != s {
				b.Fatal()
			}
		}
	})
	b.Run("String", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			if s != s {
				b.Fatal()
			}
		}
	})
	b.Run("Codec", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			out = out[:]
			enc := codec.NewEncoderBytes(&out, cbor)
			if err := enc.Encode(s); err != nil {
				b.Fatal(err)
			}
			dec := codec.NewDecoderBytes(out, cbor)
			hello := ""
			if err := dec.Decode(&hello); err != nil {
				b.Fatal(err)
			}
			if hello != "hello" {
				b.Fatal()
			}
		}
	})
}
