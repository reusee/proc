package proc

type Proc interface {
	Step(*Control)
}
