package proc

type ProcFunc func([]Proc) ([]Proc, error)

var _ Proc = ProcFunc(nil)

func (p ProcFunc) Run(procs []Proc) ([]Proc, error) {
	return p(procs)
}
