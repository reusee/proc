package proc

type ProcFunc func(*Next)

var _ Proc = ProcFunc(nil)

func (p ProcFunc) Step(next *Next) {
	p(next)
}
