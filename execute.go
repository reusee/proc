package proc

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
