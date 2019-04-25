package sca

func (h Head) Dup() (dup Head) {
	dup = Head{
		Id:       h.Id,
		Receiver: h.Receiver.Dup(),
		Sender:   h.Sender.Dup(),
	}
	return
}

func (e *Event) Dup() (dup *Event) {
	if e == nil {
		return
	}
	dup = &Event{
		Head:     e.Head.Dup(),
		Action:   e.Action,
		Callback: e.Callback,
	}
	l := len(e.Payload)
	payload := make([]DataObject, l)
	for i := 0; i < l; i++ {
		payload[i] = e.Payload[i].Dup()
	}
	dup.Payload = payload
	carrier := map[string]string{}
	for k, v := range e.Carrier {
		carrier[k] = v
	}
	dup.Carrier = carrier
	return
}
