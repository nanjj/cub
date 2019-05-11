package sdo

type Targets []string

//go:generate codecgen -o cg_$GOFILE $GOFILE
type Summary struct {
	Id       int64             `codec:"i"`
	To       Targets           `codec:"t"`
	From     string            `codec:"f"`
	Action   string            `codec:"a"`
	Callback string            `codec:"b"`
	Carrier  map[string]string `codec:"c"`
	Lens     []int             `codec:"l"` // Lens internal use, do not modify
}

func (sum Summary) Clone() (dup Summary) {
	dup = Summary{
		Id:       sum.Id,
		To:       sum.To.Clone(),
		From:     sum.From,
		Action:   sum.Action,
		Callback: sum.Callback,
	}
	if sum.Carrier != nil {
		dup.Carrier = map[string]string{}
		for k, v := range sum.Carrier {
			dup.Carrier[k] = v
		}
	}
	if sum.Lens != nil {
		dup.Lens = make([]int, len(sum.Lens))
		copy(dup.Lens, sum.Lens)
	}
	return
}
