package proc

import "github.com/reusee/pr2"

type Proc interface {
	Step(*Control)
}

type Control struct {
	procs []Proc
}

func (n *Control) Next(proc Proc) {
	n.procs = append(n.procs, proc)
}

func (n *Control) reset() {
	n.procs = n.procs[:0]
}

var controlsPool = pr2.NewPool(
	128,
	func() *Control {
		return &Control{
			procs: make([]Proc, 0, 8),
		}
	},
)
