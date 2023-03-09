package proc

import "github.com/reusee/pr2"

type procQueue struct {
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

func newProcQueue() *procQueue {
	part, _ := procQueuePartPool.Get()
	return &procQueue{
		head: part,
		tail: part,
	}
}

func (p *procQueue) empty() bool {
	return p.head == p.tail &&
		p.head.begin == len(p.head.procs)
}

const maxProcQueuePartCapacity = 64

func (p *procQueuePart) reset() {
	p.procs = p.procs[:0]
	p.begin = 0
	p.next = nil
}

func (p *procQueue) enqueue(proc Proc) {
	if len(p.head.procs) >= maxProcQueuePartCapacity {
		// extend
		newPart, _ := procQueuePartPool.Get()
		newPart.reset()
		p.head.next = newPart
		p.head = newPart
	}

	p.head.procs = append(p.head.procs, proc)
}

func (p *procQueue) dequeue() (Proc, bool) {
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

func (p *procQueue) end() {
	for p.tail != nil {
		part := p.tail
		p.tail = p.tail.next
		part.put()
	}
}
