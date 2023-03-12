package proc

import "github.com/reusee/pr2"

func Execute(proc Proc) error {
	next, put := nextsPool.Get()
	defer put()
	queue := newProcQueue()
	queue.enqueue(proc)
	for !queue.empty() {
		proc, _ := queue.dequeue()
		next.reset()
		err := proc.Run(next)
		if err != nil {
			return err
		}
		for _, newProc := range next.procs {
			queue.enqueue(newProc)
		}
	}
	return nil
}

var nextsPool = pr2.NewPool(
	128,
	func(_ pr2.PoolPutFunc) *Next {
		return &Next{
			procs: make([]Proc, 0, 8),
		}
	},
)
