package proc

import (
	"context"
	"sync"
)

type Worker struct {
	mutex       sync.Mutex
	cond        *sync.Cond
	queue       *ProcQueue
	ctx         context.Context
	cancel      context.CancelFunc
	newProcChan chan Proc
	next        Next
}

func NewWorker(ctx context.Context) *Worker {
	ctx, cancel := context.WithCancel(ctx)
	ret := &Worker{
		queue:       NewProcQueue(),
		ctx:         ctx,
		cancel:      cancel,
		newProcChan: make(chan Proc),
	}
	ret.cond = sync.NewCond(&ret.mutex)
	go ret.start(ctx)
	return ret
}

func (w *Worker) start(ctx context.Context) {
	defer func() {
		w.mutex.Lock()
		defer w.mutex.Unlock()
		w.queue = nil
		w.cond.Broadcast()
	}()

loop:
	for {

		if w.queue.empty() {
			select {
			case <-ctx.Done():
				break loop
			case proc := <-w.newProcChan:
				if proc != nil {
					w.queue.enqueue(proc)
				}
			}

		} else {
			select {
			case <-ctx.Done():
				break loop
			case proc := <-w.newProcChan:
				if proc != nil {
					w.queue.enqueue(proc)
				}
			default:
				w.next.reset()
				proc, _ := w.queue.dequeue()
				proc.Step(&w.next)
				for _, newProc := range w.next.procs {
					w.queue.enqueue(newProc)
				}
			}
		}

	}

	// do rest works
	w.queue.Run()

}

func (w *Worker) Do(proc Proc) error {
	select {
	case <-w.ctx.Done():
		return w.ctx.Err()
	case w.newProcChan <- proc:
	}
	return nil
}

func (w *Worker) Close() error {
	w.cancel()
	w.mutex.Lock()
	defer w.mutex.Unlock()
	for w.queue != nil {
		w.cond.Wait()
	}
	return nil
}
