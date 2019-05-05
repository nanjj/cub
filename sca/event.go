package sca

func (e *Event) Clone() (dup *Event) {
	if e == nil {
		return
	}
	dup = &Event{
		Id:       e.Id,
		To:       e.To.Clone(),
		From:     e.From,
		Action:   e.Action,
		Callback: e.Callback,
	}
	if e.Payload != nil {
		l := len(e.Payload)
		payload := make([]DataObject, l)
		for i := 0; i < l; i++ {
			payload[i] = e.Payload[i].Clone()
		}
		dup.Payload = payload
	}
	if e.Carrier != nil {
		carrier := map[string]string{}
		for k, v := range e.Carrier {
			carrier[k] = v
		}
		dup.Carrier = carrier
	}
	return
}
