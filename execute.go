package proc

import "github.com/reusee/pr2"

func Execute(proc Proc) error {
	buf, put := procsPool.Get()
	defer put()
	queue := newProcQueue()
	queue.enqueue(proc)
	for !queue.empty() {
		proc, _ := queue.dequeue()
		buf = buf[:0]
		err := proc.Run(&buf)
		if err != nil {
			return err
		}
		for _, newProc := range buf {
			queue.enqueue(newProc)
		}
	}
	return nil
}

var procsPool = pr2.NewPool(
	128,
	func(_ pr2.PoolPutFunc) []Proc {
		return make([]Proc, 0, 8)
	},
)
