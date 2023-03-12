package proc

type Proc interface {
	Run(*Next) error
}

type Next struct {
	procs []Proc
}

func (n *Next) Add(procs ...Proc) {
	n.procs = append(n.procs, procs...)
}

func (n *Next) reset() {
	n.procs = n.procs[:0]
}
