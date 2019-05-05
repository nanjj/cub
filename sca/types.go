package sca

//go:generate codecgen -o cg_$GOFILE $GOFILE

type Targets []string
type DataObject []byte
type Payload []DataObject

type Event struct {
	Id       int64             `codec:"id"`
	To       Targets           `codec:"to"`
	From     string            `codec:"from"`
	Action   string            `codec:"action"`
	Carrier  map[string]string `codec:"carrier"`
	Payload  Payload           `codec:"payload"`
	Callback string            `codec:"callback"`
}
