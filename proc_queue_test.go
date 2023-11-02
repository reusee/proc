package proc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type intProc int

func (intProc) Step(ctrl *Control) {
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
			return func(ctrl *Control) {
			}
		}
		return func(ctrl *Control) {
			*p++
			ctrl.Next(adder(i-1, p))
		}
	}

	n := 0
	queue := NewProcQueue(adder(5, &n))
	queue.Run()
	if n != 5 {
		t.Fatal()
	}

}

func BenchmarkProcQueueRun(b *testing.B) {
	var proc ProcFunc
	i := 0
	proc = func(ctrl *Control) {
		i++
		if i == b.N {
			return
		}
		ctrl.Next(proc)
	}
	queue := NewProcQueue(proc)
	b.ResetTimer()
	queue.Run()
}
