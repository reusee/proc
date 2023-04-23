package proc

import (
	"io"

	"github.com/reusee/pr2"
)

type ProcQueue struct {
	head *procQueuePart
	tail *procQueuePart
}

type procQueuePart struct {
	procs []Proc
	begin int
	next  *procQueuePart
	put   pr2.PoolPutFunc
}

var procQueuePartPool = pr2.NewPool(
	256,
	func(put pr2.PoolPutFunc) *procQueuePart {
		return &procQueuePart{
			procs: make([]Proc, 0, maxProcQueuePartCapacity),
			put:   put,
		}
	},
)

func NewProcQueue(procs ...Proc) *ProcQueue {
	part, _ := procQueuePartPool.Get()
	queue := &ProcQueue{
		head: part,
		tail: part,
	}
	for _, proc := range procs {
		queue.enqueue(proc)
	}
	return queue
}

func (p *ProcQueue) empty() bool {
	return p.head == p.tail &&
		p.head.begin == len(p.head.procs)
}

const maxProcQueuePartCapacity = 64

func (p *procQueuePart) reset() {
	p.procs = p.procs[:0]
	p.begin = 0
	p.next = nil
}

func (p *ProcQueue) enqueue(proc Proc) {
	if len(p.head.procs) >= maxProcQueuePartCapacity {
		// extend
		newPart, _ := procQueuePartPool.Get()
		newPart.reset()
		p.head.next = newPart
		p.head = newPart
	}

	p.head.procs = append(p.head.procs, proc)
}

func (p *ProcQueue) dequeue() (Proc, bool) {
	if p.empty() {
		return nil, false
	}

	if p.tail.begin >= len(p.tail.procs) {
		// shrink
		if p.tail.next == nil {
			panic("impossible")
		}
		part := p.tail
		p.tail = p.tail.next
		part.put()
	}

	proc := p.tail.procs[p.tail.begin]
	p.tail.begin++
	return proc, true
}

func (p *ProcQueue) end() {
	for p.tail != nil {
		part := p.tail
		p.tail = p.tail.next
		part.put()
	}
}

var _ Proc = new(ProcQueue)

func (q *ProcQueue) Run(next *Next) error {
	if q.empty() {
		return io.EOF
	}
	next.reset()
	proc, _ := q.dequeue()
	err := proc.Run(next)
	if err != nil {
		return err
	}
	for _, newProc := range next.procs {
		q.enqueue(newProc)
	}
	return nil
}

func (q *ProcQueue) RunAll() error {
	next, put := nextsPool.Get()
	defer put()
	for {
		err := q.Run(next)
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
	}
	return nil
}
