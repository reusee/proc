package proc

type Proc interface {
	Run([]Proc) ([]Proc, error)
}
