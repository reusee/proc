package proc

import "github.com/reusee/pr2"

type Proc interface {
	Run(*Next)
}

type Next struct {
	procs []Proc
}

func (n *Next) Add(proc Proc) {
	n.procs = append(n.procs, proc)
}

func (n *Next) reset() {
	n.procs = n.procs[:0]
}

var nextsPool = pr2.NewPool(
	128,
	func() *Next {
		return &Next{
			procs: make([]Proc, 0, 8),
		}
	},
)
