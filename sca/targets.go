package sca

// Down - all nodes in the subtree
// a -> b -> c
//      | -> d
//      ` -> e -> f
//  if targets being exectured on node b, the action will be token on
//  node b, c, d,  e, and f
func (t Targets) Down() bool {
	return t != nil && len(t) == 0
}

// ToDown - change to `Down'.
func (t *Targets) ToDown() {
	*t = []string{}
}

// Local - local node of the subtree
// if `targets` being executed on node b, b is the local node.
func (t Targets) Local() bool {
	return t == nil
}

// ToLocal - change to local
func (t *Targets) ToLocal() {
	*t = nil
}

// Clone
func (t Targets) Clone() (dup Targets) {
	if t == nil {
		return
	}
	dup = append([]string{}, t...)
	return
}
