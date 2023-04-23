package proc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type intProc int

func (intProc) Run(next *Next) error {
	return nil
}

func TestProcQueue(t *testing.T) {
	q := NewProcQueue()
	defer q.end()
	for i := 0; i < maxProcQueuePartCapacity*128; i++ {
		q.enqueue(intProc(i))
		p, ok := q.dequeue()
		assert.True(t, ok)
		assert.Equal(t, intProc(i), p)
		_, ok = q.dequeue()
		assert.False(t, ok)
	}
}

func BenchmarkProcQueueEnqueue(b *testing.B) {
	q := NewProcQueue()
	defer q.end()
	for i := 0; i < b.N; i++ {
		q.enqueue(nil)
	}
}

func BenchmarkProcQueueEnqueueAndDequeue(b *testing.B) {
	q := NewProcQueue()
	defer q.end()
	for i := 0; i < b.N; i++ {
		q.enqueue(nil)
		_, ok := q.dequeue()
		if !ok {
			b.Fatal()
		}
	}
}

func TestProcQueueRun(t *testing.T) {

	var adder func(i int, p *int) ProcFunc
	adder = func(i int, p *int) ProcFunc {
		if i == 0 {
			return func(next *Next) error {
				return nil
			}
		}
		return func(next *Next) error {
			*p++
			next.Add(adder(i-1, p))
			return nil
		}
	}

	n := 0
	queue := NewProcQueue(adder(5, &n))
	if err := queue.RunAll(); err != nil {
		t.Fatal(err)
	}
	if n != 5 {
		t.Fatal()
	}

}

func BenchmarkProcQueueRun(b *testing.B) {
	var proc ProcFunc
	i := 0
	proc = func(next *Next) error {
		i++
		if i == b.N {
			return nil
		}
		next.Add(proc)
		return nil
	}
	queue := NewProcQueue(proc)
	b.ResetTimer()
	if err := queue.RunAll(); err != nil {
		b.Fatal(err)
	}
}
