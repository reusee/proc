package proc

type Proc interface {
	Run() ([]Proc, error)
}
