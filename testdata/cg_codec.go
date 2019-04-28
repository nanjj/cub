// Code generated by codecgen - DO NOT EDIT.

package drilling

import (
	"errors"
	codec1978 "github.com/ugorji/go/codec"
	"runtime"
	"strconv"
	"time"
)

const (
	// ----- content types ----
	codecSelferCcUTF87934 = 1
	codecSelferCcRAW7934  = 255
	// ----- value types used ----
	codecSelferValueTypeArray7934  = 10
	codecSelferValueTypeMap7934    = 9
	codecSelferValueTypeString7934 = 6
	codecSelferValueTypeInt7934    = 2
	codecSelferValueTypeUint7934   = 3
	codecSelferValueTypeFloat7934  = 4
	codecSelferBitsize7934         = uint8(32 << (^uint(0) >> 63))
)

var (
	errCodecSelferOnlyMapOrArrayEncodeToStruct7934 = errors.New(`only encoded map or array can be decoded into a struct`)
)

type codecSelfer7934 struct{}

func init() {
	if codec1978.GenVersion != 10 {
		_, file, _, _ := runtime.Caller(0)
		panic("codecgen version mismatch: current: 10, need " + strconv.FormatInt(int64(codec1978.GenVersion), 10) + ". Re-generate file: " + file)
	}
	if false {
		var _ byte = 0 // reference the types, but skip this branch at build/run time
		var v0 time.Time
		_ = v0
	}
}

func (x *TestingEvent) CodecEncodeSelf(e *codec1978.Encoder) {
	var h codecSelfer7934
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
				r.WriteArrayStart(5)
			} else {
				r.WriteMapStart(5)
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
					r.EncodeStringEnc(codecSelferCcUTF87934, `id`)
				}
				r.WriteMapElemValue()
				if false {
				} else {
					r.EncodeInt(int64(x.Id))
				}
			}
			if yyr2 || yy2arr2 {
				r.WriteArrayElem()
				if false {
				} else {
					r.EncodeInt(int64(x.Kind))
				}
			} else {
				r.WriteMapElemKey()
				if z.IsJSONHandle() {
					z.WriteStr("\"kind\"")
				} else {
					r.EncodeStringEnc(codecSelferCcUTF87934, `kind`)
				}
				r.WriteMapElemValue()
				if false {
				} else {
					r.EncodeInt(int64(x.Kind))
				}
			}
			if yyr2 || yy2arr2 {
				r.WriteArrayElem()
				if false {
				} else if !z.EncBasicHandle().TimeNotBuiltin {
					r.EncodeTime(x.CreatedAt)
				} else if yyxt10 := z.Extension(z.I2Rtid(x.CreatedAt)); yyxt10 != nil {
					z.EncExtension(x.CreatedAt, yyxt10)
				} else if z.EncBinary() {
					z.EncBinaryMarshal(x.CreatedAt)
				} else if !z.EncBinary() && z.IsJSONHandle() {
					z.EncJSONMarshal(x.CreatedAt)
				} else {
					z.EncFallback(x.CreatedAt)
				}
			} else {
				r.WriteMapElemKey()
				r.EncodeStringEnc(codecSelferCcUTF87934, `created_at`)
				r.WriteMapElemValue()
				if false {
				} else if !z.EncBasicHandle().TimeNotBuiltin {
					r.EncodeTime(x.CreatedAt)
				} else if yyxt11 := z.Extension(z.I2Rtid(x.CreatedAt)); yyxt11 != nil {
					z.EncExtension(x.CreatedAt, yyxt11)
				} else if z.EncBinary() {
					z.EncBinaryMarshal(x.CreatedAt)
				} else if !z.EncBinary() && z.IsJSONHandle() {
					z.EncJSONMarshal(x.CreatedAt)
				} else {
					z.EncFallback(x.CreatedAt)
				}
			}
			if yyr2 || yy2arr2 {
				r.WriteArrayElem()
				if false {
				} else {
					if z.EncBasicHandle().StringToRaw {
						r.EncodeStringBytesRaw(z.BytesView(string(x.Target)))
					} else {
						r.EncodeStringEnc(codecSelferCcUTF87934, string(x.Target))
					}
				}
			} else {
				r.WriteMapElemKey()
				if z.IsJSONHandle() {
					z.WriteStr("\"target\"")
				} else {
					r.EncodeStringEnc(codecSelferCcUTF87934, `target`)
				}
				r.WriteMapElemValue()
				if false {
				} else {
					if z.EncBasicHandle().StringToRaw {
						r.EncodeStringBytesRaw(z.BytesView(string(x.Target)))
					} else {
						r.EncodeStringEnc(codecSelferCcUTF87934, string(x.Target))
					}
				}
			}
			if yyr2 || yy2arr2 {
				r.WriteArrayElem()
				if false {
				} else {
					if z.EncBasicHandle().StringToRaw {
						r.EncodeStringBytesRaw(z.BytesView(string(x.Source)))
					} else {
						r.EncodeStringEnc(codecSelferCcUTF87934, string(x.Source))
					}
				}
			} else {
				r.WriteMapElemKey()
				if z.IsJSONHandle() {
					z.WriteStr("\"source\"")
				} else {
					r.EncodeStringEnc(codecSelferCcUTF87934, `source`)
				}
				r.WriteMapElemValue()
				if false {
				} else {
					if z.EncBasicHandle().StringToRaw {
						r.EncodeStringBytesRaw(z.BytesView(string(x.Source)))
					} else {
						r.EncodeStringEnc(codecSelferCcUTF87934, string(x.Source))
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

func (x *TestingEvent) CodecDecodeSelf(d *codec1978.Decoder) {
	var h codecSelfer7934
	z, r := codec1978.GenHelperDecoder(d)
	_, _, _ = h, z, r
	if false {
	} else if yyxt1 := z.Extension(z.I2Rtid(x)); yyxt1 != nil {
		z.DecExtension(x, yyxt1)
	} else {
		yyct2 := r.ContainerType()
		if yyct2 == codecSelferValueTypeMap7934 {
			yyl2 := r.ReadMapStart()
			if yyl2 == 0 {
				r.ReadMapEnd()
			} else {
				x.codecDecodeSelfFromMap(yyl2, d)
			}
		} else if yyct2 == codecSelferValueTypeArray7934 {
			yyl2 := r.ReadArrayStart()
			if yyl2 == 0 {
				r.ReadArrayEnd()
			} else {
				x.codecDecodeSelfFromArray(yyl2, d)
			}
		} else {
			panic(errCodecSelferOnlyMapOrArrayEncodeToStruct7934)
		}
	}
}

func (x *TestingEvent) codecDecodeSelfFromMap(l int, d *codec1978.Decoder) {
	var h codecSelfer7934
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
				x.Id = (int)(z.C.IntV(r.DecodeInt64(), codecSelferBitsize7934))
			}
		case "kind":
			if r.TryDecodeAsNil() {
				x.Kind = 0
			} else {
				x.Kind = (int)(z.C.IntV(r.DecodeInt64(), codecSelferBitsize7934))
			}
		case "created_at":
			if r.TryDecodeAsNil() {
				x.CreatedAt = time.Time{}
			} else {
				if false {
				} else if !z.DecBasicHandle().TimeNotBuiltin {
					x.CreatedAt = r.DecodeTime()
				} else if yyxt7 := z.Extension(z.I2Rtid(x.CreatedAt)); yyxt7 != nil {
					z.DecExtension(x.CreatedAt, yyxt7)
				} else if z.DecBinary() {
					z.DecBinaryUnmarshal(&x.CreatedAt)
				} else if !z.DecBinary() && z.IsJSONHandle() {
					z.DecJSONUnmarshal(&x.CreatedAt)
				} else {
					z.DecFallback(&x.CreatedAt, false)
				}
			}
		case "target":
			if r.TryDecodeAsNil() {
				x.Target = ""
			} else {
				x.Target = (string)(r.DecodeString())
			}
		case "source":
			if r.TryDecodeAsNil() {
				x.Source = ""
			} else {
				x.Source = (string)(r.DecodeString())
			}
		default:
			z.DecStructFieldNotFound(-1, yys3)
		} // end switch yys3
	} // end for yyj3
	r.ReadMapEnd()
}

func (x *TestingEvent) codecDecodeSelfFromArray(l int, d *codec1978.Decoder) {
	var h codecSelfer7934
	z, r := codec1978.GenHelperDecoder(d)
	_, _, _ = h, z, r
	var yyj10 int
	var yyb10 bool
	var yyhl10 bool = l >= 0
	yyj10++
	if yyhl10 {
		yyb10 = yyj10 > l
	} else {
		yyb10 = r.CheckBreak()
	}
	if yyb10 {
		r.ReadArrayEnd()
		return
	}
	r.ReadArrayElem()
	if r.TryDecodeAsNil() {
		x.Id = 0
	} else {
		x.Id = (int)(z.C.IntV(r.DecodeInt64(), codecSelferBitsize7934))
	}
	yyj10++
	if yyhl10 {
		yyb10 = yyj10 > l
	} else {
		yyb10 = r.CheckBreak()
	}
	if yyb10 {
		r.ReadArrayEnd()
		return
	}
	r.ReadArrayElem()
	if r.TryDecodeAsNil() {
		x.Kind = 0
	} else {
		x.Kind = (int)(z.C.IntV(r.DecodeInt64(), codecSelferBitsize7934))
	}
	yyj10++
	if yyhl10 {
		yyb10 = yyj10 > l
	} else {
		yyb10 = r.CheckBreak()
	}
	if yyb10 {
		r.ReadArrayEnd()
		return
	}
	r.ReadArrayElem()
	if r.TryDecodeAsNil() {
		x.CreatedAt = time.Time{}
	} else {
		if false {
		} else if !z.DecBasicHandle().TimeNotBuiltin {
			x.CreatedAt = r.DecodeTime()
		} else if yyxt14 := z.Extension(z.I2Rtid(x.CreatedAt)); yyxt14 != nil {
			z.DecExtension(x.CreatedAt, yyxt14)
		} else if z.DecBinary() {
			z.DecBinaryUnmarshal(&x.CreatedAt)
		} else if !z.DecBinary() && z.IsJSONHandle() {
			z.DecJSONUnmarshal(&x.CreatedAt)
		} else {
			z.DecFallback(&x.CreatedAt, false)
		}
	}
	yyj10++
	if yyhl10 {
		yyb10 = yyj10 > l
	} else {
		yyb10 = r.CheckBreak()
	}
	if yyb10 {
		r.ReadArrayEnd()
		return
	}
	r.ReadArrayElem()
	if r.TryDecodeAsNil() {
		x.Target = ""
	} else {
		x.Target = (string)(r.DecodeString())
	}
	yyj10++
	if yyhl10 {
		yyb10 = yyj10 > l
	} else {
		yyb10 = r.CheckBreak()
	}
	if yyb10 {
		r.ReadArrayEnd()
		return
	}
	r.ReadArrayElem()
	if r.TryDecodeAsNil() {
		x.Source = ""
	} else {
		x.Source = (string)(r.DecodeString())
	}
	for {
		yyj10++
		if yyhl10 {
			yyb10 = yyj10 > l
		} else {
			yyb10 = r.CheckBreak()
		}
		if yyb10 {
			break
		}
		r.ReadArrayElem()
		z.DecStructFieldNotFound(yyj10-1, "")
	}
	r.ReadArrayEnd()
}

func (x *TestingHead) CodecEncodeSelf(e *codec1978.Encoder) {
	var h codecSelfer7934
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
				r.WriteArrayStart(2)
			} else {
				r.WriteMapStart(2)
			}
			if yyr2 || yy2arr2 {
				r.WriteArrayElem()
				if x.Targets == nil {
					r.EncodeNil()
				} else {
					if false {
					} else {
						z.F.EncSliceInt64V(x.Targets, e)
					}
				}
			} else {
				r.WriteMapElemKey()
				if z.IsJSONHandle() {
					z.WriteStr("\"targets\"")
				} else {
					r.EncodeStringEnc(codecSelferCcUTF87934, `targets`)
				}
				r.WriteMapElemValue()
				if x.Targets == nil {
					r.EncodeNil()
				} else {
					if false {
					} else {
						z.F.EncSliceInt64V(x.Targets, e)
					}
				}
			}
			if yyr2 || yy2arr2 {
				r.WriteArrayElem()
				if false {
				} else {
					r.EncodeInt(int64(x.Callback))
				}
			} else {
				r.WriteMapElemKey()
				if z.IsJSONHandle() {
					z.WriteStr("\"callback\"")
				} else {
					r.EncodeStringEnc(codecSelferCcUTF87934, `callback`)
				}
				r.WriteMapElemValue()
				if false {
				} else {
					r.EncodeInt(int64(x.Callback))
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

func (x *TestingHead) CodecDecodeSelf(d *codec1978.Decoder) {
	var h codecSelfer7934
	z, r := codec1978.GenHelperDecoder(d)
	_, _, _ = h, z, r
	if false {
	} else if yyxt1 := z.Extension(z.I2Rtid(x)); yyxt1 != nil {
		z.DecExtension(x, yyxt1)
	} else {
		yyct2 := r.ContainerType()
		if yyct2 == codecSelferValueTypeMap7934 {
			yyl2 := r.ReadMapStart()
			if yyl2 == 0 {
				r.ReadMapEnd()
			} else {
				x.codecDecodeSelfFromMap(yyl2, d)
			}
		} else if yyct2 == codecSelferValueTypeArray7934 {
			yyl2 := r.ReadArrayStart()
			if yyl2 == 0 {
				r.ReadArrayEnd()
			} else {
				x.codecDecodeSelfFromArray(yyl2, d)
			}
		} else {
			panic(errCodecSelferOnlyMapOrArrayEncodeToStruct7934)
		}
	}
}

func (x *TestingHead) codecDecodeSelfFromMap(l int, d *codec1978.Decoder) {
	var h codecSelfer7934
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
		case "targets":
			if r.TryDecodeAsNil() {
				x.Targets = nil
			} else {
				if false {
				} else {
					z.F.DecSliceInt64X(&x.Targets, d)
				}
			}
		case "callback":
			if r.TryDecodeAsNil() {
				x.Callback = 0
			} else {
				x.Callback = (int64)(r.DecodeInt64())
			}
		default:
			z.DecStructFieldNotFound(-1, yys3)
		} // end switch yys3
	} // end for yyj3
	r.ReadMapEnd()
}

func (x *TestingHead) codecDecodeSelfFromArray(l int, d *codec1978.Decoder) {
	var h codecSelfer7934
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
		x.Targets = nil
	} else {
		if false {
		} else {
			z.F.DecSliceInt64X(&x.Targets, d)
		}
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
		x.Callback = 0
	} else {
		x.Callback = (int64)(r.DecodeInt64())
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

func (x *TestingBody) CodecEncodeSelf(e *codec1978.Encoder) {
	var h codecSelfer7934
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
					if z.EncBasicHandle().StringToRaw {
						r.EncodeStringBytesRaw(z.BytesView(string(x.Action)))
					} else {
						r.EncodeStringEnc(codecSelferCcUTF87934, string(x.Action))
					}
				}
			} else {
				r.WriteMapElemKey()
				if z.IsJSONHandle() {
					z.WriteStr("\"Action\"")
				} else {
					r.EncodeStringEnc(codecSelferCcUTF87934, `Action`)
				}
				r.WriteMapElemValue()
				if false {
				} else {
					if z.EncBasicHandle().StringToRaw {
						r.EncodeStringBytesRaw(z.BytesView(string(x.Action)))
					} else {
						r.EncodeStringEnc(codecSelferCcUTF87934, string(x.Action))
					}
				}
			}
			if yyr2 || yy2arr2 {
				r.WriteArrayElem()
				if false {
				} else {
					if z.EncBasicHandle().StringToRaw {
						r.EncodeStringBytesRaw(z.BytesView(string(x.Command)))
					} else {
						r.EncodeStringEnc(codecSelferCcUTF87934, string(x.Command))
					}
				}
			} else {
				r.WriteMapElemKey()
				if z.IsJSONHandle() {
					z.WriteStr("\"Command\"")
				} else {
					r.EncodeStringEnc(codecSelferCcUTF87934, `Command`)
				}
				r.WriteMapElemValue()
				if false {
				} else {
					if z.EncBasicHandle().StringToRaw {
						r.EncodeStringBytesRaw(z.BytesView(string(x.Command)))
					} else {
						r.EncodeStringEnc(codecSelferCcUTF87934, string(x.Command))
					}
				}
			}
			if yyr2 || yy2arr2 {
				r.WriteArrayElem()
				if x.Args == nil {
					r.EncodeNil()
				} else {
					if false {
					} else {
						z.F.EncSliceStringV(x.Args, e)
					}
				}
			} else {
				r.WriteMapElemKey()
				if z.IsJSONHandle() {
					z.WriteStr("\"Args\"")
				} else {
					r.EncodeStringEnc(codecSelferCcUTF87934, `Args`)
				}
				r.WriteMapElemValue()
				if x.Args == nil {
					r.EncodeNil()
				} else {
					if false {
					} else {
						z.F.EncSliceStringV(x.Args, e)
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

func (x *TestingBody) CodecDecodeSelf(d *codec1978.Decoder) {
	var h codecSelfer7934
	z, r := codec1978.GenHelperDecoder(d)
	_, _, _ = h, z, r
	if false {
	} else if yyxt1 := z.Extension(z.I2Rtid(x)); yyxt1 != nil {
		z.DecExtension(x, yyxt1)
	} else {
		yyct2 := r.ContainerType()
		if yyct2 == codecSelferValueTypeMap7934 {
			yyl2 := r.ReadMapStart()
			if yyl2 == 0 {
				r.ReadMapEnd()
			} else {
				x.codecDecodeSelfFromMap(yyl2, d)
			}
		} else if yyct2 == codecSelferValueTypeArray7934 {
			yyl2 := r.ReadArrayStart()
			if yyl2 == 0 {
				r.ReadArrayEnd()
			} else {
				x.codecDecodeSelfFromArray(yyl2, d)
			}
		} else {
			panic(errCodecSelferOnlyMapOrArrayEncodeToStruct7934)
		}
	}
}

func (x *TestingBody) codecDecodeSelfFromMap(l int, d *codec1978.Decoder) {
	var h codecSelfer7934
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
		case "Action":
			if r.TryDecodeAsNil() {
				x.Action = ""
			} else {
				x.Action = (string)(r.DecodeString())
			}
		case "Command":
			if r.TryDecodeAsNil() {
				x.Command = ""
			} else {
				x.Command = (string)(r.DecodeString())
			}
		case "Args":
			if r.TryDecodeAsNil() {
				x.Args = nil
			} else {
				if false {
				} else {
					z.F.DecSliceStringX(&x.Args, d)
				}
			}
		default:
			z.DecStructFieldNotFound(-1, yys3)
		} // end switch yys3
	} // end for yyj3
	r.ReadMapEnd()
}

func (x *TestingBody) codecDecodeSelfFromArray(l int, d *codec1978.Decoder) {
	var h codecSelfer7934
	z, r := codec1978.GenHelperDecoder(d)
	_, _, _ = h, z, r
	var yyj8 int
	var yyb8 bool
	var yyhl8 bool = l >= 0
	yyj8++
	if yyhl8 {
		yyb8 = yyj8 > l
	} else {
		yyb8 = r.CheckBreak()
	}
	if yyb8 {
		r.ReadArrayEnd()
		return
	}
	r.ReadArrayElem()
	if r.TryDecodeAsNil() {
		x.Action = ""
	} else {
		x.Action = (string)(r.DecodeString())
	}
	yyj8++
	if yyhl8 {
		yyb8 = yyj8 > l
	} else {
		yyb8 = r.CheckBreak()
	}
	if yyb8 {
		r.ReadArrayEnd()
		return
	}
	r.ReadArrayElem()
	if r.TryDecodeAsNil() {
		x.Command = ""
	} else {
		x.Command = (string)(r.DecodeString())
	}
	yyj8++
	if yyhl8 {
		yyb8 = yyj8 > l
	} else {
		yyb8 = r.CheckBreak()
	}
	if yyb8 {
		r.ReadArrayEnd()
		return
	}
	r.ReadArrayElem()
	if r.TryDecodeAsNil() {
		x.Args = nil
	} else {
		if false {
		} else {
			z.F.DecSliceStringX(&x.Args, d)
		}
	}
	for {
		yyj8++
		if yyhl8 {
			yyb8 = yyj8 > l
		} else {
			yyb8 = r.CheckBreak()
		}
		if yyb8 {
			break
		}
		r.ReadArrayElem()
		z.DecStructFieldNotFound(yyj8-1, "")
	}
	r.ReadArrayEnd()
}

func (x *TestingHeadBody) CodecEncodeSelf(e *codec1978.Encoder) {
	var h codecSelfer7934
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
				r.WriteArrayStart(5)
			} else {
				r.WriteMapStart(5)
			}
			if yyr2 || yy2arr2 {
				r.WriteArrayElem()
				if x.Targets == nil {
					r.EncodeNil()
				} else {
					if false {
					} else {
						z.F.EncSliceInt64V(x.Targets, e)
					}
				}
			} else {
				r.WriteMapElemKey()
				if z.IsJSONHandle() {
					z.WriteStr("\"targets\"")
				} else {
					r.EncodeStringEnc(codecSelferCcUTF87934, `targets`)
				}
				r.WriteMapElemValue()
				if x.Targets == nil {
					r.EncodeNil()
				} else {
					if false {
					} else {
						z.F.EncSliceInt64V(x.Targets, e)
					}
				}
			}
			if yyr2 || yy2arr2 {
				r.WriteArrayElem()
				if false {
				} else {
					r.EncodeInt(int64(x.Callback))
				}
			} else {
				r.WriteMapElemKey()
				if z.IsJSONHandle() {
					z.WriteStr("\"callback\"")
				} else {
					r.EncodeStringEnc(codecSelferCcUTF87934, `callback`)
				}
				r.WriteMapElemValue()
				if false {
				} else {
					r.EncodeInt(int64(x.Callback))
				}
			}
			if yyr2 || yy2arr2 {
				r.WriteArrayElem()
				if false {
				} else {
					if z.EncBasicHandle().StringToRaw {
						r.EncodeStringBytesRaw(z.BytesView(string(x.Action)))
					} else {
						r.EncodeStringEnc(codecSelferCcUTF87934, string(x.Action))
					}
				}
			} else {
				r.WriteMapElemKey()
				if z.IsJSONHandle() {
					z.WriteStr("\"Action\"")
				} else {
					r.EncodeStringEnc(codecSelferCcUTF87934, `Action`)
				}
				r.WriteMapElemValue()
				if false {
				} else {
					if z.EncBasicHandle().StringToRaw {
						r.EncodeStringBytesRaw(z.BytesView(string(x.Action)))
					} else {
						r.EncodeStringEnc(codecSelferCcUTF87934, string(x.Action))
					}
				}
			}
			if yyr2 || yy2arr2 {
				r.WriteArrayElem()
				if false {
				} else {
					if z.EncBasicHandle().StringToRaw {
						r.EncodeStringBytesRaw(z.BytesView(string(x.Command)))
					} else {
						r.EncodeStringEnc(codecSelferCcUTF87934, string(x.Command))
					}
				}
			} else {
				r.WriteMapElemKey()
				if z.IsJSONHandle() {
					z.WriteStr("\"Command\"")
				} else {
					r.EncodeStringEnc(codecSelferCcUTF87934, `Command`)
				}
				r.WriteMapElemValue()
				if false {
				} else {
					if z.EncBasicHandle().StringToRaw {
						r.EncodeStringBytesRaw(z.BytesView(string(x.Command)))
					} else {
						r.EncodeStringEnc(codecSelferCcUTF87934, string(x.Command))
					}
				}
			}
			if yyr2 || yy2arr2 {
				r.WriteArrayElem()
				if x.Args == nil {
					r.EncodeNil()
				} else {
					if false {
					} else {
						z.F.EncSliceStringV(x.Args, e)
					}
				}
			} else {
				r.WriteMapElemKey()
				if z.IsJSONHandle() {
					z.WriteStr("\"Args\"")
				} else {
					r.EncodeStringEnc(codecSelferCcUTF87934, `Args`)
				}
				r.WriteMapElemValue()
				if x.Args == nil {
					r.EncodeNil()
				} else {
					if false {
					} else {
						z.F.EncSliceStringV(x.Args, e)
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

func (x *TestingHeadBody) CodecDecodeSelf(d *codec1978.Decoder) {
	var h codecSelfer7934
	z, r := codec1978.GenHelperDecoder(d)
	_, _, _ = h, z, r
	if false {
	} else if yyxt1 := z.Extension(z.I2Rtid(x)); yyxt1 != nil {
		z.DecExtension(x, yyxt1)
	} else {
		yyct2 := r.ContainerType()
		if yyct2 == codecSelferValueTypeMap7934 {
			yyl2 := r.ReadMapStart()
			if yyl2 == 0 {
				r.ReadMapEnd()
			} else {
				x.codecDecodeSelfFromMap(yyl2, d)
			}
		} else if yyct2 == codecSelferValueTypeArray7934 {
			yyl2 := r.ReadArrayStart()
			if yyl2 == 0 {
				r.ReadArrayEnd()
			} else {
				x.codecDecodeSelfFromArray(yyl2, d)
			}
		} else {
			panic(errCodecSelferOnlyMapOrArrayEncodeToStruct7934)
		}
	}
}

func (x *TestingHeadBody) codecDecodeSelfFromMap(l int, d *codec1978.Decoder) {
	var h codecSelfer7934
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
		case "targets":
			if r.TryDecodeAsNil() {
				x.TestingHead.Targets = nil
			} else {
				if false {
				} else {
					z.F.DecSliceInt64X(&x.Targets, d)
				}
			}
		case "callback":
			if r.TryDecodeAsNil() {
				x.TestingHead.Callback = 0
			} else {
				x.Callback = (int64)(r.DecodeInt64())
			}
		case "Action":
			if r.TryDecodeAsNil() {
				x.TestingBody.Action = ""
			} else {
				x.Action = (string)(r.DecodeString())
			}
		case "Command":
			if r.TryDecodeAsNil() {
				x.TestingBody.Command = ""
			} else {
				x.Command = (string)(r.DecodeString())
			}
		case "Args":
			if r.TryDecodeAsNil() {
				x.TestingBody.Args = nil
			} else {
				if false {
				} else {
					z.F.DecSliceStringX(&x.Args, d)
				}
			}
		default:
			z.DecStructFieldNotFound(-1, yys3)
		} // end switch yys3
	} // end for yyj3
	r.ReadMapEnd()
}

func (x *TestingHeadBody) codecDecodeSelfFromArray(l int, d *codec1978.Decoder) {
	var h codecSelfer7934
	z, r := codec1978.GenHelperDecoder(d)
	_, _, _ = h, z, r
	var yyj11 int
	var yyb11 bool
	var yyhl11 bool = l >= 0
	yyj11++
	if yyhl11 {
		yyb11 = yyj11 > l
	} else {
		yyb11 = r.CheckBreak()
	}
	if yyb11 {
		r.ReadArrayEnd()
		return
	}
	r.ReadArrayElem()
	if r.TryDecodeAsNil() {
		x.TestingHead.Targets = nil
	} else {
		if false {
		} else {
			z.F.DecSliceInt64X(&x.Targets, d)
		}
	}
	yyj11++
	if yyhl11 {
		yyb11 = yyj11 > l
	} else {
		yyb11 = r.CheckBreak()
	}
	if yyb11 {
		r.ReadArrayEnd()
		return
	}
	r.ReadArrayElem()
	if r.TryDecodeAsNil() {
		x.TestingHead.Callback = 0
	} else {
		x.Callback = (int64)(r.DecodeInt64())
	}
	yyj11++
	if yyhl11 {
		yyb11 = yyj11 > l
	} else {
		yyb11 = r.CheckBreak()
	}
	if yyb11 {
		r.ReadArrayEnd()
		return
	}
	r.ReadArrayElem()
	if r.TryDecodeAsNil() {
		x.TestingBody.Action = ""
	} else {
		x.Action = (string)(r.DecodeString())
	}
	yyj11++
	if yyhl11 {
		yyb11 = yyj11 > l
	} else {
		yyb11 = r.CheckBreak()
	}
	if yyb11 {
		r.ReadArrayEnd()
		return
	}
	r.ReadArrayElem()
	if r.TryDecodeAsNil() {
		x.TestingBody.Command = ""
	} else {
		x.Command = (string)(r.DecodeString())
	}
	yyj11++
	if yyhl11 {
		yyb11 = yyj11 > l
	} else {
		yyb11 = r.CheckBreak()
	}
	if yyb11 {
		r.ReadArrayEnd()
		return
	}
	r.ReadArrayElem()
	if r.TryDecodeAsNil() {
		x.TestingBody.Args = nil
	} else {
		if false {
		} else {
			z.F.DecSliceStringX(&x.Args, d)
		}
	}
	for {
		yyj11++
		if yyhl11 {
			yyb11 = yyj11 > l
		} else {
			yyb11 = r.CheckBreak()
		}
		if yyb11 {
			break
		}
		r.ReadArrayElem()
		z.DecStructFieldNotFound(yyj11-1, "")
	}
	r.ReadArrayEnd()
}