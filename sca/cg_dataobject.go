// Code generated by codecgen - DO NOT EDIT.

package sca

import (
	"errors"
	codec1978 "github.com/ugorji/go/codec"
	"runtime"
	"strconv"
)

const (
	// ----- content types ----
	codecSelferCcUTF86867 = 1
	codecSelferCcRAW6867  = 255
	// ----- value types used ----
	codecSelferValueTypeArray6867  = 10
	codecSelferValueTypeMap6867    = 9
	codecSelferValueTypeString6867 = 6
	codecSelferValueTypeInt6867    = 2
	codecSelferValueTypeUint6867   = 3
	codecSelferValueTypeFloat6867  = 4
	codecSelferBitsize6867         = uint8(32 << (^uint(0) >> 63))
)

var (
	errCodecSelferOnlyMapOrArrayEncodeToStruct6867 = errors.New(`only encoded map or array can be decoded into a struct`)
)

type codecSelfer6867 struct{}

func init() {
	if codec1978.GenVersion != 10 {
		_, file, _, _ := runtime.Caller(0)
		panic("codecgen version mismatch: current: 10, need " + strconv.FormatInt(int64(codec1978.GenVersion), 10) + ". Re-generate file: " + file)
	}
	if false {
		var _ byte = 0 // reference the types, but skip this branch at build/run time
	}
}

func (x DataObject) CodecEncodeSelf(e *codec1978.Encoder) {
	var h codecSelfer6867
	z, r := codec1978.GenHelperEncoder(e)
	_, _, _ = h, z, r
	if x == nil {
		r.EncodeNil()
	} else {
		if false {
		} else if yyxt1 := z.Extension(z.I2Rtid(x)); yyxt1 != nil {
			z.EncExtension(x, yyxt1)
		} else {
			h.encDataObject((DataObject)(x), e)
		}
	}
}

func (x *DataObject) CodecDecodeSelf(d *codec1978.Decoder) {
	var h codecSelfer6867
	z, r := codec1978.GenHelperDecoder(d)
	_, _, _ = h, z, r
	if false {
	} else if yyxt1 := z.Extension(z.I2Rtid(x)); yyxt1 != nil {
		z.DecExtension(x, yyxt1)
	} else {
		h.decDataObject((*DataObject)(x), d)
	}
}

func (x codecSelfer6867) encDataObject(v DataObject, e *codec1978.Encoder) {
	var h codecSelfer6867
	z, r := codec1978.GenHelperEncoder(e)
	_, _, _ = h, z, r
	r.EncodeStringBytesRaw([]byte(v))
}

func (x codecSelfer6867) decDataObject(v *DataObject, d *codec1978.Decoder) {
	var h codecSelfer6867
	z, r := codec1978.GenHelperDecoder(d)
	_, _, _ = h, z, r
	*v = r.DecodeBytes(*((*[]byte)(v)), false)
}