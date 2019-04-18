package sca

import "github.com/ugorji/go/codec"

var (
	cbor = &codec.CborHandle{}
)

func (d *DataObject) Encode(i interface{}) (err error) {
	switch v := i.(type) {
	case string:
		*d = []byte(v)
	case []byte:
		if v != nil {
			*d = make([]byte, len(v))
			copy(*d, v)
		}
	default: // cbor
		out := make([]byte, 0, 256)
		enc := codec.NewEncoderBytes(&out, cbor)
		if err = enc.Encode(v); err == nil {
			*d = out
		}
	}
	return
}

func (d DataObject) Decode(i interface{}) (err error) {
	switch v := i.(type) {
	case *string:
		*v = string(d)
	case *[]byte:
		*v = make([]byte, len(d))
		copy(*v, d)
	default:
		dec := codec.NewDecoderBytes(d, cbor)
		err = dec.Decode(v)
	}
	return
}

func (d DataObject) Dup() (dup DataObject) {
	if d == nil {
		return
	}
	dup = make([]byte, len(d))
	copy(dup, d)
	return
}
