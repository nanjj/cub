package sca

func (h Head) Clone() (dup Head) {
	dup = Head{
		Id:       h.Id,
		Receiver: h.Receiver.Clone(),
		Sender:   h.Sender.Clone(),
	}
	return
}

func (e *Event) Clone() (dup *Event) {
	if e == nil {
		return
	}
	dup = &Event{
		Head:     e.Head.Clone(),
		Action:   e.Action,
		Callback: e.Callback,
	}
	l := len(e.Payload)
	payload := make([]DataObject, l)
	for i := 0; i < l; i++ {
		payload[i] = e.Payload[i].Clone()
	}
	dup.Payload = payload
	carrier := map[string]string{}
	for k, v := range e.Carrier {
		carrier[k] = v
	}
	dup.Carrier = carrier
	return
}
