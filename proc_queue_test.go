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
	q := newProcQueue()
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
	q := newProcQueue()
	defer q.end()
	for i := 0; i < b.N; i++ {
		q.enqueue(nil)
	}
}

func BenchmarkProcQueueEnqueueAndDequeue(b *testing.B) {
	q := newProcQueue()
	defer q.end()
	for i := 0; i < b.N; i++ {
		q.enqueue(nil)
		_, ok := q.dequeue()
		if !ok {
			b.Fatal()
		}
	}
}
