package proc

type ProcFunc func(*[]Proc) error

var _ Proc = ProcFunc(nil)

func (p ProcFunc) Run(procs *[]Proc) error {
	return p(procs)
}
