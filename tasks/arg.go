package tasks

type Arg []byte

type Input struct {
	Err  string `codec:"err"`
	Args []Arg  `codec:"args"`
}
