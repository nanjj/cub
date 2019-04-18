package sca

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

func (t Targets) Dup() (dup Targets) {
	if t == nil {
		return
	}
	dup = append([]string{}, t...)
	return
}

func (h Head) Dup() (dup Head) {
	dup = Head{
		Id:       h.Id,
		Receiver: h.Receiver.Dup(),
		Sender:   h.Sender.Dup(),
	}
	return
}
