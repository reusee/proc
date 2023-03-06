package proc

import "testing"

func TestExecute(t *testing.T) {

	var adder func(i int, p *int) ProcFunc
	adder = func(i int, p *int) ProcFunc {
		if i == 0 {
			return func() ([]Proc, error) {
				return nil, nil
			}
		}
		return func() ([]Proc, error) {
			*p++
			return []Proc{adder(i-1, p)}, nil
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
	proc = func() ([]Proc, error) {
		i++
		if i == b.N {
			return nil, nil
		}
		return []Proc{proc}, nil
	}
	Execute(proc)
}
