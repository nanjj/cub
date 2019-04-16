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
	codecSelferCcUTF87565 = 1
	codecSelferCcRAW7565  = 255
	// ----- value types used ----
	codecSelferValueTypeArray7565  = 10
	codecSelferValueTypeMap7565    = 9
	codecSelferValueTypeString7565 = 6
	codecSelferValueTypeInt7565    = 2
	codecSelferValueTypeUint7565   = 3
	codecSelferValueTypeFloat7565  = 4
	codecSelferBitsize7565         = uint8(32 << (^uint(0) >> 63))
)

var (
	errCodecSelferOnlyMapOrArrayEncodeToStruct7565 = errors.New(`only encoded map or array can be decoded into a struct`)
)

type codecSelfer7565 struct{}

func init() {
	if codec1978.GenVersion != 10 {
		_, file, _, _ := runtime.Caller(0)
		panic("codecgen version mismatch: current: 10, need " + strconv.FormatInt(int64(codec1978.GenVersion), 10) + ". Re-generate file: " + file)
	}
	if false {
		var _ byte = 0 // reference the types, but skip this branch at build/run time
	}
}

func (x Targets) CodecEncodeSelf(e *codec1978.Encoder) {
	var h codecSelfer7565
	z, r := codec1978.GenHelperEncoder(e)
	_, _, _ = h, z, r
	if x == nil {
		r.EncodeNil()
	} else {
		if false {
		} else if yyxt1 := z.Extension(z.I2Rtid(x)); yyxt1 != nil {
			z.EncExtension(x, yyxt1)
		} else {
			h.encTargets((Targets)(x), e)
		}
	}
}

func (x *Targets) CodecDecodeSelf(d *codec1978.Decoder) {
	var h codecSelfer7565
	z, r := codec1978.GenHelperDecoder(d)
	_, _, _ = h, z, r
	if false {
	} else if yyxt1 := z.Extension(z.I2Rtid(x)); yyxt1 != nil {
		z.DecExtension(x, yyxt1)
	} else {
		h.decTargets((*Targets)(x), d)
	}
}

func (x *Head) CodecEncodeSelf(e *codec1978.Encoder) {
	var h codecSelfer7565
	z, r := codec1978.GenHelperEncoder(e)
	_, _, _ = h, z, r
	if x == nil {
		r.EncodeNil()
	} else {
		if false {
		} else if yyxt1 := z.Extension(z.I2Rtid(x)); yyxt1 != nil {
			z.EncExtension(x, yyxt1)
		} else {
			yysep2 := !z.EncBinary()
			yy2arr2 := z.EncBasicHandle().StructToArray
			_, _ = yysep2, yy2arr2
			const yyr2 bool = false // struct tag has 'toArray'
			if yyr2 || yy2arr2 {
				r.WriteArrayStart(3)
			} else {
				r.WriteMapStart(3)
			}
			if yyr2 || yy2arr2 {
				r.WriteArrayElem()
				if false {
				} else {
					r.EncodeInt(int64(x.Id))
				}
			} else {
				r.WriteMapElemKey()
				if z.IsJSONHandle() {
					z.WriteStr("\"id\"")
				} else {
					r.EncodeStringEnc(codecSelferCcUTF87565, `id`)
				}
				r.WriteMapElemValue()
				if false {
				} else {
					r.EncodeInt(int64(x.Id))
				}
			}
			if yyr2 || yy2arr2 {
				r.WriteArrayElem()
				if x.Receiver == nil {
					r.EncodeNil()
				} else {
					x.Receiver.CodecEncodeSelf(e)
				}
			} else {
				r.WriteMapElemKey()
				if z.IsJSONHandle() {
					z.WriteStr("\"receiver\"")
				} else {
					r.EncodeStringEnc(codecSelferCcUTF87565, `receiver`)
				}
				r.WriteMapElemValue()
				if x.Receiver == nil {
					r.EncodeNil()
				} else {
					x.Receiver.CodecEncodeSelf(e)
				}
			}
			if yyr2 || yy2arr2 {
				r.WriteArrayElem()
				if x.Sender == nil {
					r.EncodeNil()
				} else {
					x.Sender.CodecEncodeSelf(e)
				}
			} else {
				r.WriteMapElemKey()
				if z.IsJSONHandle() {
					z.WriteStr("\"sender\"")
				} else {
					r.EncodeStringEnc(codecSelferCcUTF87565, `sender`)
				}
				r.WriteMapElemValue()
				if x.Sender == nil {
					r.EncodeNil()
				} else {
					x.Sender.CodecEncodeSelf(e)
				}
			}
			if yyr2 || yy2arr2 {
				r.WriteArrayEnd()
			} else {
				r.WriteMapEnd()
			}
		}
	}
}

func (x *Head) CodecDecodeSelf(d *codec1978.Decoder) {
	var h codecSelfer7565
	z, r := codec1978.GenHelperDecoder(d)
	_, _, _ = h, z, r
	if false {
	} else if yyxt1 := z.Extension(z.I2Rtid(x)); yyxt1 != nil {
		z.DecExtension(x, yyxt1)
	} else {
		yyct2 := r.ContainerType()
		if yyct2 == codecSelferValueTypeMap7565 {
			yyl2 := r.ReadMapStart()
			if yyl2 == 0 {
				r.ReadMapEnd()
			} else {
				x.codecDecodeSelfFromMap(yyl2, d)
			}
		} else if yyct2 == codecSelferValueTypeArray7565 {
			yyl2 := r.ReadArrayStart()
			if yyl2 == 0 {
				r.ReadArrayEnd()
			} else {
				x.codecDecodeSelfFromArray(yyl2, d)
			}
		} else {
			panic(errCodecSelferOnlyMapOrArrayEncodeToStruct7565)
		}
	}
}

func (x *Head) codecDecodeSelfFromMap(l int, d *codec1978.Decoder) {
	var h codecSelfer7565
	z, r := codec1978.GenHelperDecoder(d)
	_, _, _ = h, z, r
	var yyhl3 bool = l >= 0
	for yyj3 := 0; ; yyj3++ {
		if yyhl3 {
			if yyj3 >= l {
				break
			}
		} else {
			if r.CheckBreak() {
				break
			}
		}
		r.ReadMapElemKey()
		yys3 := z.StringView(r.DecodeStringAsBytes())
		r.ReadMapElemValue()
		switch yys3 {
		case "id":
			if r.TryDecodeAsNil() {
				x.Id = 0
			} else {
				x.Id = (int64)(r.DecodeInt64())
			}
		case "receiver":
			if r.TryDecodeAsNil() {
				x.Receiver = nil
			} else {
				x.Receiver.CodecDecodeSelf(d)
			}
		case "sender":
			if r.TryDecodeAsNil() {
				x.Sender = nil
			} else {
				x.Sender.CodecDecodeSelf(d)
			}
		default:
			z.DecStructFieldNotFound(-1, yys3)
		} // end switch yys3
	} // end for yyj3
	r.ReadMapEnd()
}

func (x *Head) codecDecodeSelfFromArray(l int, d *codec1978.Decoder) {
	var h codecSelfer7565
	z, r := codec1978.GenHelperDecoder(d)
	_, _, _ = h, z, r
	var yyj7 int
	var yyb7 bool
	var yyhl7 bool = l >= 0
	yyj7++
	if yyhl7 {
		yyb7 = yyj7 > l
	} else {
		yyb7 = r.CheckBreak()
	}
	if yyb7 {
		r.ReadArrayEnd()
		return
	}
	r.ReadArrayElem()
	if r.TryDecodeAsNil() {
		x.Id = 0
	} else {
		x.Id = (int64)(r.DecodeInt64())
	}
	yyj7++
	if yyhl7 {
		yyb7 = yyj7 > l
	} else {
		yyb7 = r.CheckBreak()
	}
	if yyb7 {
		r.ReadArrayEnd()
		return
	}
	r.ReadArrayElem()
	if r.TryDecodeAsNil() {
		x.Receiver = nil
	} else {
		x.Receiver.CodecDecodeSelf(d)
	}
	yyj7++
	if yyhl7 {
		yyb7 = yyj7 > l
	} else {
		yyb7 = r.CheckBreak()
	}
	if yyb7 {
		r.ReadArrayEnd()
		return
	}
	r.ReadArrayElem()
	if r.TryDecodeAsNil() {
		x.Sender = nil
	} else {
		x.Sender.CodecDecodeSelf(d)
	}
	for {
		yyj7++
		if yyhl7 {
			yyb7 = yyj7 > l
		} else {
			yyb7 = r.CheckBreak()
		}
		if yyb7 {
			break
		}
		r.ReadArrayElem()
		z.DecStructFieldNotFound(yyj7-1, "")
	}
	r.ReadArrayEnd()
}

func (x *Event) CodecEncodeSelf(e *codec1978.Encoder) {
	var h codecSelfer7565
	z, r := codec1978.GenHelperEncoder(e)
	_, _, _ = h, z, r
	if x == nil {
		r.EncodeNil()
	} else {
		if false {
		} else if yyxt1 := z.Extension(z.I2Rtid(x)); yyxt1 != nil {
			z.EncExtension(x, yyxt1)
		} else {
			yysep2 := !z.EncBinary()
			yy2arr2 := z.EncBasicHandle().StructToArray
			_, _ = yysep2, yy2arr2
			const yyr2 bool = false // struct tag has 'toArray'
			if yyr2 || yy2arr2 {
				r.WriteArrayStart(7)
			} else {
				r.WriteMapStart(7)
			}
			if yyr2 || yy2arr2 {
				r.WriteArrayElem()
				if false {
				} else {
					r.EncodeInt(int64(x.Id))
				}
			} else {
				r.WriteMapElemKey()
				if z.IsJSONHandle() {
					z.WriteStr("\"id\"")
				} else {
					r.EncodeStringEnc(codecSelferCcUTF87565, `id`)
				}
				r.WriteMapElemValue()
				if false {
				} else {
					r.EncodeInt(int64(x.Id))
				}
			}
			if yyr2 || yy2arr2 {
				r.WriteArrayElem()
				if x.Receiver == nil {
					r.EncodeNil()
				} else {
					x.Receiver.CodecEncodeSelf(e)
				}
			} else {
				r.WriteMapElemKey()
				if z.IsJSONHandle() {
					z.WriteStr("\"receiver\"")
				} else {
					r.EncodeStringEnc(codecSelferCcUTF87565, `receiver`)
				}
				r.WriteMapElemValue()
				if x.Receiver == nil {
					r.EncodeNil()
				} else {
					x.Receiver.CodecEncodeSelf(e)
				}
			}
			if yyr2 || yy2arr2 {
				r.WriteArrayElem()
				if x.Sender == nil {
					r.EncodeNil()
				} else {
					x.Sender.CodecEncodeSelf(e)
				}
			} else {
				r.WriteMapElemKey()
				if z.IsJSONHandle() {
					z.WriteStr("\"sender\"")
				} else {
					r.EncodeStringEnc(codecSelferCcUTF87565, `sender`)
				}
				r.WriteMapElemValue()
				if x.Sender == nil {
					r.EncodeNil()
				} else {
					x.Sender.CodecEncodeSelf(e)
				}
			}
			if yyr2 || yy2arr2 {
				r.WriteArrayElem()
				if false {
				} else {
					if z.EncBasicHandle().StringToRaw {
						r.EncodeStringBytesRaw(z.BytesView(string(x.Action)))
					} else {
						r.EncodeStringEnc(codecSelferCcUTF87565, string(x.Action))
					}
				}
			} else {
				r.WriteMapElemKey()
				if z.IsJSONHandle() {
					z.WriteStr("\"action\"")
				} else {
					r.EncodeStringEnc(codecSelferCcUTF87565, `action`)
				}
				r.WriteMapElemValue()
				if false {
				} else {
					if z.EncBasicHandle().StringToRaw {
						r.EncodeStringBytesRaw(z.BytesView(string(x.Action)))
					} else {
						r.EncodeStringEnc(codecSelferCcUTF87565, string(x.Action))
					}
				}
			}
			if yyr2 || yy2arr2 {
				r.WriteArrayElem()
				if x.Carrier == nil {
					r.EncodeNil()
				} else {
					if false {
					} else {
						z.F.EncMapStringStringV(x.Carrier, e)
					}
				}
			} else {
				r.WriteMapElemKey()
				if z.IsJSONHandle() {
					z.WriteStr("\"carrier\"")
				} else {
					r.EncodeStringEnc(codecSelferCcUTF87565, `carrier`)
				}
				r.WriteMapElemValue()
				if x.Carrier == nil {
					r.EncodeNil()
				} else {
					if false {
					} else {
						z.F.EncMapStringStringV(x.Carrier, e)
					}
				}
			}
			if yyr2 || yy2arr2 {
				r.WriteArrayElem()
				if x.Payload == nil {
					r.EncodeNil()
				} else {
					yysf19 := &x.Payload
					yysf19.CodecEncodeSelf(e)
				}
			} else {
				r.WriteMapElemKey()
				if z.IsJSONHandle() {
					z.WriteStr("\"payload\"")
				} else {
					r.EncodeStringEnc(codecSelferCcUTF87565, `payload`)
				}
				r.WriteMapElemValue()
				if x.Payload == nil {
					r.EncodeNil()
				} else {
					yysf20 := &x.Payload
					yysf20.CodecEncodeSelf(e)
				}
			}
			if yyr2 || yy2arr2 {
				r.WriteArrayElem()
				if false {
				} else {
					if z.EncBasicHandle().StringToRaw {
						r.EncodeStringBytesRaw(z.BytesView(string(x.Callback)))
					} else {
						r.EncodeStringEnc(codecSelferCcUTF87565, string(x.Callback))
					}
				}
			} else {
				r.WriteMapElemKey()
				if z.IsJSONHandle() {
					z.WriteStr("\"callback\"")
				} else {
					r.EncodeStringEnc(codecSelferCcUTF87565, `callback`)
				}
				r.WriteMapElemValue()
				if false {
				} else {
					if z.EncBasicHandle().StringToRaw {
						r.EncodeStringBytesRaw(z.BytesView(string(x.Callback)))
					} else {
						r.EncodeStringEnc(codecSelferCcUTF87565, string(x.Callback))
					}
				}
			}
			if yyr2 || yy2arr2 {
				r.WriteArrayEnd()
			} else {
				r.WriteMapEnd()
			}
		}
	}
}

func (x *Event) CodecDecodeSelf(d *codec1978.Decoder) {
	var h codecSelfer7565
	z, r := codec1978.GenHelperDecoder(d)
	_, _, _ = h, z, r
	if false {
	} else if yyxt1 := z.Extension(z.I2Rtid(x)); yyxt1 != nil {
		z.DecExtension(x, yyxt1)
	} else {
		yyct2 := r.ContainerType()
		if yyct2 == codecSelferValueTypeMap7565 {
			yyl2 := r.ReadMapStart()
			if yyl2 == 0 {
				r.ReadMapEnd()
			} else {
				x.codecDecodeSelfFromMap(yyl2, d)
			}
		} else if yyct2 == codecSelferValueTypeArray7565 {
			yyl2 := r.ReadArrayStart()
			if yyl2 == 0 {
				r.ReadArrayEnd()
			} else {
				x.codecDecodeSelfFromArray(yyl2, d)
			}
		} else {
			panic(errCodecSelferOnlyMapOrArrayEncodeToStruct7565)
		}
	}
}

func (x *Event) codecDecodeSelfFromMap(l int, d *codec1978.Decoder) {
	var h codecSelfer7565
	z, r := codec1978.GenHelperDecoder(d)
	_, _, _ = h, z, r
	var yyhl3 bool = l >= 0
	for yyj3 := 0; ; yyj3++ {
		if yyhl3 {
			if yyj3 >= l {
				break
			}
		} else {
			if r.CheckBreak() {
				break
			}
		}
		r.ReadMapElemKey()
		yys3 := z.StringView(r.DecodeStringAsBytes())
		r.ReadMapElemValue()
		switch yys3 {
		case "id":
			if r.TryDecodeAsNil() {
				x.Head.Id = 0
			} else {
				x.Id = (int64)(r.DecodeInt64())
			}
		case "receiver":
			if r.TryDecodeAsNil() {
				x.Head.Receiver = nil
			} else {
				x.Receiver.CodecDecodeSelf(d)
			}
		case "sender":
			if r.TryDecodeAsNil() {
				x.Head.Sender = nil
			} else {
				x.Sender.CodecDecodeSelf(d)
			}
		case "action":
			if r.TryDecodeAsNil() {
				x.Action = ""
			} else {
				x.Action = (string)(r.DecodeString())
			}
		case "carrier":
			if r.TryDecodeAsNil() {
				x.Carrier = nil
			} else {
				if false {
				} else {
					z.F.DecMapStringStringX(&x.Carrier, d)
				}
			}
		case "payload":
			if r.TryDecodeAsNil() {
				x.Payload = nil
			} else {
				x.Payload.CodecDecodeSelf(d)
			}
		case "callback":
			if r.TryDecodeAsNil() {
				x.Callback = ""
			} else {
				x.Callback = (string)(r.DecodeString())
			}
		default:
			z.DecStructFieldNotFound(-1, yys3)
		} // end switch yys3
	} // end for yyj3
	r.ReadMapEnd()
}

func (x *Event) codecDecodeSelfFromArray(l int, d *codec1978.Decoder) {
	var h codecSelfer7565
	z, r := codec1978.GenHelperDecoder(d)
	_, _, _ = h, z, r
	var yyj12 int
	var yyb12 bool
	var yyhl12 bool = l >= 0
	yyj12++
	if yyhl12 {
		yyb12 = yyj12 > l
	} else {
		yyb12 = r.CheckBreak()
	}
	if yyb12 {
		r.ReadArrayEnd()
		return
	}
	r.ReadArrayElem()
	if r.TryDecodeAsNil() {
		x.Head.Id = 0
	} else {
		x.Id = (int64)(r.DecodeInt64())
	}
	yyj12++
	if yyhl12 {
		yyb12 = yyj12 > l
	} else {
		yyb12 = r.CheckBreak()
	}
	if yyb12 {
		r.ReadArrayEnd()
		return
	}
	r.ReadArrayElem()
	if r.TryDecodeAsNil() {
		x.Head.Receiver = nil
	} else {
		x.Receiver.CodecDecodeSelf(d)
	}
	yyj12++
	if yyhl12 {
		yyb12 = yyj12 > l
	} else {
		yyb12 = r.CheckBreak()
	}
	if yyb12 {
		r.ReadArrayEnd()
		return
	}
	r.ReadArrayElem()
	if r.TryDecodeAsNil() {
		x.Head.Sender = nil
	} else {
		x.Sender.CodecDecodeSelf(d)
	}
	yyj12++
	if yyhl12 {
		yyb12 = yyj12 > l
	} else {
		yyb12 = r.CheckBreak()
	}
	if yyb12 {
		r.ReadArrayEnd()
		return
	}
	r.ReadArrayElem()
	if r.TryDecodeAsNil() {
		x.Action = ""
	} else {
		x.Action = (string)(r.DecodeString())
	}
	yyj12++
	if yyhl12 {
		yyb12 = yyj12 > l
	} else {
		yyb12 = r.CheckBreak()
	}
	if yyb12 {
		r.ReadArrayEnd()
		return
	}
	r.ReadArrayElem()
	if r.TryDecodeAsNil() {
		x.Carrier = nil
	} else {
		if false {
		} else {
			z.F.DecMapStringStringX(&x.Carrier, d)
		}
	}
	yyj12++
	if yyhl12 {
		yyb12 = yyj12 > l
	} else {
		yyb12 = r.CheckBreak()
	}
	if yyb12 {
		r.ReadArrayEnd()
		return
	}
	r.ReadArrayElem()
	if r.TryDecodeAsNil() {
		x.Payload = nil
	} else {
		x.Payload.CodecDecodeSelf(d)
	}
	yyj12++
	if yyhl12 {
		yyb12 = yyj12 > l
	} else {
		yyb12 = r.CheckBreak()
	}
	if yyb12 {
		r.ReadArrayEnd()
		return
	}
	r.ReadArrayElem()
	if r.TryDecodeAsNil() {
		x.Callback = ""
	} else {
		x.Callback = (string)(r.DecodeString())
	}
	for {
		yyj12++
		if yyhl12 {
			yyb12 = yyj12 > l
		} else {
			yyb12 = r.CheckBreak()
		}
		if yyb12 {
			break
		}
		r.ReadArrayElem()
		z.DecStructFieldNotFound(yyj12-1, "")
	}
	r.ReadArrayEnd()
}

func (x codecSelfer7565) encTargets(v Targets, e *codec1978.Encoder) {
	var h codecSelfer7565
	z, r := codec1978.GenHelperEncoder(e)
	_, _, _ = h, z, r
	r.WriteArrayStart(len(v))
	for _, yyv1 := range v {
		r.WriteArrayElem()
		if false {
		} else {
			if z.EncBasicHandle().StringToRaw {
				r.EncodeStringBytesRaw(z.BytesView(string(yyv1)))
			} else {
				r.EncodeStringEnc(codecSelferCcUTF87565, string(yyv1))
			}
		}
	}
	r.WriteArrayEnd()
}

func (x codecSelfer7565) decTargets(v *Targets, d *codec1978.Decoder) {
	var h codecSelfer7565
	z, r := codec1978.GenHelperDecoder(d)
	_, _, _ = h, z, r

	yyv1 := *v
	yyh1, yyl1 := z.DecSliceHelperStart()
	var yyc1 bool
	_ = yyc1
	if yyl1 == 0 {
		if yyv1 == nil {
			yyv1 = []string{}
			yyc1 = true
		} else if len(yyv1) != 0 {
			yyv1 = yyv1[:0]
			yyc1 = true
		}
	} else {
		yyhl1 := yyl1 > 0
		var yyrl1 int
		_ = yyrl1
		if yyhl1 {
			if yyl1 > cap(yyv1) {
				yyrl1 = z.DecInferLen(yyl1, z.DecBasicHandle().MaxInitLen, 16)
				if yyrl1 <= cap(yyv1) {
					yyv1 = yyv1[:yyrl1]
				} else {
					yyv1 = make([]string, yyrl1)
				}
				yyc1 = true
			} else if yyl1 != len(yyv1) {
				yyv1 = yyv1[:yyl1]
				yyc1 = true
			}
		}
		var yyj1 int
		// var yydn1 bool
		for yyj1 = 0; (yyhl1 && yyj1 < yyl1) || !(yyhl1 || r.CheckBreak()); yyj1++ { // bounds-check-elimination
			if yyj1 == 0 && yyv1 == nil {
				if yyhl1 {
					yyrl1 = z.DecInferLen(yyl1, z.DecBasicHandle().MaxInitLen, 16)
				} else {
					yyrl1 = 8
				}
				yyv1 = make([]string, yyrl1)
				yyc1 = true
			}
			yyh1.ElemContainerState(yyj1)

			var yydb1 bool
			if yyj1 >= len(yyv1) {
				yyv1 = append(yyv1, "")
				yyc1 = true

			}
			if yydb1 {
				z.DecSwallow()
			} else {
				if r.TryDecodeAsNil() {
					yyv1[yyj1] = ""
				} else {
					yyv1[yyj1] = (string)(r.DecodeString())
				}

			}

		}
		if yyj1 < len(yyv1) {
			yyv1 = yyv1[:yyj1]
			yyc1 = true
		} else if yyj1 == 0 && yyv1 == nil {
			yyv1 = make([]string, 0)
			yyc1 = true
		}
	}
	yyh1.End()
	if yyc1 {
		*v = yyv1
	}
}
