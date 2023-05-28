package proc

type ProcFunc func(*Next)

var _ Proc = ProcFunc(nil)

func (p ProcFunc) Run(next *Next) {
	p(next)
}
