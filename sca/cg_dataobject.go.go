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
	codecSelferCcUTF83653 = 1
	codecSelferCcRAW3653  = 255
	// ----- value types used ----
	codecSelferValueTypeArray3653  = 10
	codecSelferValueTypeMap3653    = 9
	codecSelferValueTypeString3653 = 6
	codecSelferValueTypeInt3653    = 2
	codecSelferValueTypeUint3653   = 3
	codecSelferValueTypeFloat3653  = 4
	codecSelferBitsize3653         = uint8(32 << (^uint(0) >> 63))
)

var (
	errCodecSelferOnlyMapOrArrayEncodeToStruct3653 = errors.New(`only encoded map or array can be decoded into a struct`)
)

type codecSelfer3653 struct{}

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
	var h codecSelfer3653
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
	var h codecSelfer3653
	z, r := codec1978.GenHelperDecoder(d)
	_, _, _ = h, z, r
	if false {
	} else if yyxt1 := z.Extension(z.I2Rtid(x)); yyxt1 != nil {
		z.DecExtension(x, yyxt1)
	} else {
		h.decDataObject((*DataObject)(x), d)
	}
}

func (x codecSelfer3653) encDataObject(v DataObject, e *codec1978.Encoder) {
	var h codecSelfer3653
	z, r := codec1978.GenHelperEncoder(e)
	_, _, _ = h, z, r
	r.EncodeStringBytesRaw([]byte(v))
}

func (x codecSelfer3653) decDataObject(v *DataObject, d *codec1978.Decoder) {
	var h codecSelfer3653
	z, r := codec1978.GenHelperDecoder(d)
	_, _, _ = h, z, r
	*v = r.DecodeBytes(*((*[]byte)(v)), false)
}
