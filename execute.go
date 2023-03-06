package proc

func Execute(proc Proc) error {
	queue := newProcQueue()
	queue.enqueue(proc)
	for !queue.empty() {
		proc, _ := queue.dequeue()
		newProcs, err := proc.Run()
		if err != nil {
			return err
		}
		for _, newProc := range newProcs {
			queue.enqueue(newProc)
		}
	}
	return nil
}
