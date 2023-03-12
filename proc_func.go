package proc

type ProcFunc func(*Next) error

var _ Proc = ProcFunc(nil)

func (p ProcFunc) Run(next *Next) error {
	return p(next)
}
