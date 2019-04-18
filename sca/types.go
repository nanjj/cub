package sca

//go:generate codecgen -o cg_$GOFILE $GOFILE

type Targets []string
type DataObject []byte
type Payload []DataObject

type Head struct {
	Id       int64   `codec:"id"`
	Receiver Targets `codec:"receiver"`
	Sender   Targets `codec:"sender"`
}

type Event struct {
	Head
	Action   string            `codec:"action"`
	Carrier  map[string]string `codec:"carrier"`
	Payload  Payload           `codec:"payload"`
	Callback string            `codec:"callback"`
}
