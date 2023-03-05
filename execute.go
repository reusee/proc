package proc

import "container/list"

func Execute(proc Proc) error {
	procs := list.New()
	procs.PushBack(proc)
	for front := procs.Front(); front != nil; front = procs.Front() {
		procs.Remove(front)
		newProcs, err := front.Value.(Proc).Run()
		if err != nil {
			return err
		}
		for _, newProc := range newProcs {
			procs.PushBack(newProc)
		}
	}
	return nil
}
