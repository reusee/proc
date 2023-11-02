package proc

type ProcFunc func(*Control)

var _ Proc = ProcFunc(nil)

func (p ProcFunc) Step(ctrl *Control) {
	p(ctrl)
}
