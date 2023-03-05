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
