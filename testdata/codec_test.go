package drilling

import (
	"encoding/json"
	"reflect"
	"testing"
	"time"

	"github.com/ugorji/go/codec"
)

func TestCodecCborTime(t *testing.T) {
	cborHand := &codec.MsgpackHandle{}
	input := time.Now().Round(0)
	out := make([]byte, 32)
	enc := codec.NewEncoderBytes(&out, cborHand)
	if err := enc.Encode(input); err != nil {
		t.Fatal(err)
	}
	in := out
	dec := codec.NewDecoderBytes(in, cborHand)
	result := time.Time{}
	if err := dec.Decode(&result); err != nil {
		t.Fatal(err)
	}
	if !result.Equal(input) {
		t.Fatal(result, input)
	}
}

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
	cborHand := &codec.CborHandle{}
	out := make([]byte, 0, 1024)
	enc := codec.NewEncoderBytes(&out, cborHand)
	now := time.Now().UTC()
	v := &TestingEvent{1, 1, now, "hyper", "scheduler"}
	if err := enc.Encode(v); err != nil {
		t.Fatal(err)
	}
	t.Log(len(out), string(out))
	dec := codec.NewDecoderBytes(out, cborHand)
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
	msgpackHand := &codec.MsgpackHandle{}
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
	b.Run("StdJson", func(b *testing.B) {
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

	b.Run("CborMsgpack", func(b *testing.B){
		benchCodec(b, msgpackHand)
	})
}

func TestCborTime(t *testing.T) {
	cborHand := &codec.CborHandle{}
	out := make([]byte, 0, 16)
	enc := codec.NewEncoderBytes(&out, cborHand)
	now := time.Now().Round(0)
	if err := enc.Encode(&now); err != nil {
		t.Fatal(err)
	}
	t.Log(len(out))
	dec := codec.NewDecoderBytes(out, cborHand)
	want := time.Time{}
	if err := dec.Decode(&want); err != nil {
		t.Fatal(err)
	}
	if !want.Equal(now) {
		t.Log(want)
		t.Log(now)
		t.Fatal()
	}
}

var testingData = func() []byte {
	targets := []int64{1, 2, 3, 4, 5, 6, 7, 8, 9}
	m := &TestingHeadBody{
		TestingHead{
			targets,
			10,
		},
		TestingBody{"Script", "launch_vm.sh", []string{"hypercube01", "4", "4096", "40", "192.168.0.24/24"}},
	}
	b := make([]byte, 0, 1024)
	enc := codec.NewEncoderBytes(&b, &codec.CborHandle{})
	enc.Encode(m)
	return b
}()

func TestCborPartDecode(t *testing.T) {
	cborHand := &codec.CborHandle{}
	head := &TestingHead{}
	dec := codec.NewDecoderBytes(testingData, cborHand)
	if err := dec.Decode(head); err != nil {
		t.Fatal(err)
	}
	targets := []int64{1, 2, 3, 4, 5, 6, 7, 8, 9}
	if !reflect.DeepEqual(head.Targets, targets) {
		t.Fatal(head)
	}
}

func BenchmarkPartlyDecode(b *testing.B) {
	cborHand := &codec.CborHandle{}
	b.Run("HeadOnly", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			v := &TestingHead{}
			dec := codec.NewDecoderBytes(testingData, cborHand)
			if err := dec.Decode(v); err != nil {
				b.Fatal(err)
			}
		}
	})
	b.Run("WithBody", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			v := &TestingHeadBody{}
			dec := codec.NewDecoderBytes(testingData, cborHand)
			if err := dec.Decode(v); err != nil {
				b.Fatal(err)
			}
		}
	})
}
