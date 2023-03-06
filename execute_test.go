package proc

import "testing"

func TestExecute(t *testing.T) {

	var adder func(i int, p *int) ProcFunc
	adder = func(i int, p *int) ProcFunc {
		if i == 0 {
			return func(procs []Proc) ([]Proc, error) {
				return nil, nil
			}
		}
		return func(procs []Proc) ([]Proc, error) {
			*p++
			procs = append(procs, adder(i-1, p))
			return procs, nil
		}
	}

	n := 0
	if err := Execute(adder(5, &n)); err != nil {
		t.Fatal(err)
	}
	if n != 5 {
		t.Fatal()
	}

}

func BenchmarkExecute(b *testing.B) {
	var proc ProcFunc
	i := 0
	proc = func(procs []Proc) ([]Proc, error) {
		i++
		if i == b.N {
			return nil, nil
		}
		procs = append(procs, proc)
		return procs, nil
	}
	Execute(proc)
}
