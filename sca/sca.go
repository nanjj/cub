package sca

type DataObject []byte

func (d *DataObject) String() string {
	return string(*d)
}

func (d *DataObject) Bytes() []byte {
	return (*d)
}

func (d *DataObject) Decode(v interface{}) (err error) {
	return
}

func (d *DataObject) Encode(v interface{}) (err error) {
	return
}

type Targets []string

func (targets Targets) All() bool {
	return targets != nil && len(targets) == 0
}

func (targets Targets) Local() bool {
	return targets == nil
}

func (targets *Targets) ToAll() {
	*targets = []string{}
}

func (targets *Targets) ToLocal() {
	*targets = nil
}

type Head struct {
	Id       int64   `codec:"id"`
	Receiver Targets `codec:"receiver"`
	Sender   Targets `codec:"sender"`
}

//go:generate codecgen -o cg_$GOFILE $GOFILE
type Message struct {
	Head
	Action   string            `codec:"action"`
	Carrier  map[string]string `codec:"carrier"`
	Payload  DataObject        `codec:"payload"`
	Callback string            `codec:"callback"`
}
