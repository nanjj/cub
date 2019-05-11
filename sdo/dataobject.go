package sdo

type DataObject []byte

func (d *DataObject) Encode(v interface{}) (err error) {
	*d, err = Encode(v)
	return
}

func (d DataObject) Decode(i interface{}) (err error) {
	err = Decode(i, d)
	return
}

func (d DataObject) Clone() (dup DataObject) {
	if d == nil {
		return
	}
	dup = make([]byte, len(d))
	copy(dup, d)
	return
}
