package sca

// All - all nodes in the subtree
// a -> b -> c
//      | -> d
//      ` -> e -> f
//  if targets being exectured on node b, the action will be token on
//  node b, c, d,  e, and f
func (targets Targets) All() bool {
	return targets != nil && len(targets) == 0
}

// ToAll - change to `All'.
func (targets *Targets) ToAll() {
	*targets = []string{}
}

// Local - local node of the subtree
// if `targets` being executed on node b, b is the local node.
func (targets Targets) Local() bool {
	return targets == nil
}

// ToLocal - change to local
func (targets *Targets) ToLocal() {
	*targets = nil
}

// Clone
func (t Targets) Clone() (dup Targets) {
	if t == nil {
		return
	}
	dup = append([]string{}, t...)
	return
}
