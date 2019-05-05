package sca

type Set map[string]bool

func NewSet(elems ...string) (set Set) {
	set = map[string]bool{}
	for _, elem := range elems {
		set[elem] = true
	}
	return
}

func (set Set) Add(elem string) {
	set[elem] = true
}
