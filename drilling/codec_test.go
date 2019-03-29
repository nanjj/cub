package drilling

import (
	"encoding/json"
	"reflect"
	"testing"
	"time"

	"github.com/ugorji/go/codec"
)

func TestCodecJsonUsage(t *testing.T) {
	jsonHand := &codec.JsonHandle{
		Indent: 2,
	}
	out := []byte{}
	enc := codec.NewEncoderBytes(&out, jsonHand)
	now := time.Now().Round(0) // strip monotonic time
	v := &TestingEvent{1, 1, now, "hyper", "scheduler"}
	t.Log(now)
	if err := enc.Encode(v); err != nil {
		t.Fatal(err)
	}
	t.Log(len(out))
	dec := codec.NewDecoderBytes(out, jsonHand)
	want := &TestingEvent{}
	if err := dec.Decode(want); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(v, want) {
		t.Log(v)
		t.Log(want)
		t.Fatal()
	}
}

func TestCodecCbor(t *testing.T) {
	jsonHand := &codec.CborHandle{}
	out := make([]byte, 0, 1024)
	enc := codec.NewEncoderBytes(&out, jsonHand)
	now := time.Now().UTC()
	v := &TestingEvent{1, 1, now, "hyper", "scheduler"}
	if err := enc.Encode(v); err != nil {
		t.Fatal(err)
	}
	t.Log(len(out), string(out))
	dec := codec.NewDecoderBytes(out, jsonHand)
	want := &TestingEvent{}
	if err := dec.Decode(want); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(v, want) {
		t.Log(v)
		t.Log(want)
		t.Fatal()
	}
}

func BenchmarkCodecJsonCbor(b *testing.B) {
	now := time.Now().Round(0)
	v := &TestingEvent{1, 1, now, "hyper", "scheduler"}
	jsonHand := &codec.JsonHandle{}
	cborHand := &codec.CborHandle{}
	out := make([]byte, 0, 1024)
	benchCodec := func(b *testing.B, h codec.Handle) {
		for i := 0; i < b.N; i++ {
			out = out[0:0]
			enc := codec.NewEncoderBytes(&out, h)
			if err := enc.Encode(v); err != nil {
				b.Fatal(err)
			}
			in := out
			dec := codec.NewDecoderBytes(in, h)
			want := &TestingEvent{}
			if err := dec.Decode(want); err != nil {
				b.Fatal(err)
			}
			if !want.CreatedAt.Equal(v.CreatedAt) {
				b.Log(v)
				b.Log(want)
				b.Fatal()
			}
		}
	}
	b.Run("CodecJson", func(b *testing.B) {
		benchCodec(b, jsonHand)
	})
	b.Run("EncodingJson", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			if out, err := json.Marshal(v); err != nil {
				b.Fatal(out, err)
			} else {
				in := out
				want := &TestingEvent{}
				if err := json.Unmarshal(in, want); err != nil {
					b.Fatal(err)
				} else if !want.CreatedAt.Equal(v.CreatedAt) {
					b.Log(v)
					b.Log(want)
					b.Fatal()
				}
			}
		}
	})
	b.Run("CodecCbor", func(b *testing.B) {
		benchCodec(b, cborHand)
	})
}
