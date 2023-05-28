package proc

import (
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
	put   func() bool
}

var procQueuePartPool = pr2.NewPool(
	256,
	func() *procQueuePart {
		return &procQueuePart{
			procs: make([]Proc, 0, maxProcQueuePartCapacity),
		}
	},
)

func NewProcQueue(procs ...Proc) *ProcQueue {
	var part *procQueuePart
	put := procQueuePartPool.Get(&part)
	if part.put == nil {
		part.put = put
	}
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
		var newPart *procQueuePart
		put := procQueuePartPool.Get(&newPart)
		if newPart.put == nil {
			newPart.put = put
		}
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

func (q *ProcQueue) Run(next *Next) {
	if q.empty() {
		return
	}
	next.reset()
	proc, _ := q.dequeue()
	proc.Run(next)
	for _, newProc := range next.procs {
		q.enqueue(newProc)
	}
}

func (q *ProcQueue) RunAll() {
	var next *Next
	put := nextsPool.Get(&next)
	defer put()
	for {
		if q.empty() {
			return
		}
		q.Run(next)
	}
}
