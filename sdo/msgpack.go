package sdo

import (
	"github.com/ugorji/go/codec"
)

var (
	msgpack = &codec.MsgpackHandle{}
)

func Encode(v interface{}) (b []byte, err error) {
	switch val := v.(type) {
	case string:
		b = []byte(val)
	case []byte:
		if v != nil {
			b = make([]byte, len(val))
			copy(b, val)
		}
	default: // msgpack
		out := make([]byte, 0, 256)
		enc := codec.NewEncoderBytes(&out, msgpack)
		if err = enc.Encode(val); err == nil {
			b = out
		}
	}
	return
}

func Decode(v interface{}, b []byte) (err error) {
	switch val := v.(type) {
	case *string:
		*val = string(b)
	case *[]byte:
		*val = make([]byte, len(b))
		copy(*val, b)
	default:
		dec := codec.NewDecoderBytes(b, msgpack)
		err = dec.Decode(val)
	}
	return
}
