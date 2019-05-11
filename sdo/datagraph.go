package sdo

type Payload []DataObject

type DataGraph struct {
	Summary
	Payload Payload
}

func (g *DataGraph) Load(b []byte) (err error) {
	if len(b) == 0 {
		err = ErrInvalidData
		return
	}
	sum := g.Summary
	err = Decode(&sum, b)
	if err != nil {
		return
	}
	g.Summary = sum
	ls := sum.Lens
	if ls == nil {
		return
	}
	n := len(ls)
	if n == 0 {
		g.Payload = Payload{}
		return
	}
	l := len(b)
	data := make([]DataObject, n)
	to := l
	for i := n - 1; i >= 0; i-- {
		from := to - ls[i]
		if from <= 0 {
			err = ErrIncompletData
			return
		}
		data[i] = DataObject(b[from:to])
		to = from
	}
	g.Payload = data
	return
}

func (g *DataGraph) Bytes() (b []byte, err error) {
	d := make([]byte, 0, 1024)
	payload := g.Payload
	n := len(payload)
	if n == 0 {
		if payload == nil {
			g.Summary.Lens = nil
		} else {
			g.Summary.Lens = []int{}
		}
	} else {
		ls := make([]int, n)
		for i := 0; i < n; i++ {
			ls[i] = len(payload[i])
			d = append(d, []byte(payload[i])...)
		}
		g.Summary.Lens = ls
	}
	sum := g.Summary
	h, err := Encode(&sum)
	if err != nil {
		return
	}
	hl := len(h)
	b = make([]byte, len(d)+hl)
	copy(b, h)
	copy(b[hl:], d)
	return
}

func (g *DataGraph) Clone() (dup *DataGraph) {
	if g == nil {
		return
	}
	dup = &DataGraph{
		Summary: g.Summary.Clone(),
	}
	if g.Payload != nil {
		dup.Payload = Payload{}
		for i := range g.Payload {
			dup.Payload = append(dup.Payload, g.Payload[i].Clone())
		}
	}
	return
}
