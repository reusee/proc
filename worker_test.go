package proc

import (
	"context"
	"errors"
	"io"
	"sync"
	"sync/atomic"
	"testing"
)

func TestWorker(t *testing.T) {
	n := 0
	max := 1024
	var proc ProcFunc
	proc = func(next *Next) error {
		n++
		if n == max {
			return nil
		}
		next.Add(proc)
		return nil
	}

	w := NewWorker(context.Background())
	if err := w.Do(proc); err != nil {
		t.Fatal(err)
	}
	if err := w.Close(); err != nil {
		t.Fatal(err)
	}
	if err := w.Close(); err != nil {
		t.Fatal(err)
	}
	if n != max {
		t.Fatalf("got %d", n)
	}

	err := w.Do(nil)
	if !errors.Is(err, context.Canceled) {
		t.Fatal()
	}
}

func TestWorkerError(t *testing.T) {
	w := NewWorker(context.Background())
	if err := w.Do(ProcFunc(func(_ *Next) error {
		return io.EOF
	})); err != nil {
		t.Fatal(err)
	}
	err := w.Close()
	if !errors.Is(err, io.EOF) {
		t.Fatal()
	}
	err = w.Close()
	if !errors.Is(err, io.EOF) {
		t.Fatal()
	}
}

func TestWorkerConcurrent(t *testing.T) {
	w := NewWorker(context.Background())
	wg := new(sync.WaitGroup)
	n := 1024
	var c atomic.Int64
	wg.Add(n)
	for i := 0; i < n; i++ {
		go func() {
			defer wg.Done()
			if err := w.Do(ProcFunc(func(_ *Next) error {
				c.Add(1)
				return nil
			})); err != nil {
				panic(err)
			}
		}()
	}
	wg.Wait()
	if err := w.Close(); err != nil {
		t.Fatal(err)
	}
	if c.Load() != 1024 {
		t.Fatal()
	}
}

func TestWorkerCanceledCtx(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	w := NewWorker(ctx)
	cancel()
	err := w.Do(nil)
	if !errors.Is(err, context.Canceled) {
		t.Fatal()
	}
	if err := w.Close(); err != nil {
		t.Fatal(err)
	}

	ctx, cancel = context.WithCancel(context.Background())
	cancel()
	w = NewWorker(ctx)
	err = w.Do(nil)
	if !errors.Is(err, context.Canceled) {
		t.Fatal()
	}
	if err := w.Close(); err != nil {
		t.Fatal(err)
	}
}

func BenchmarkWorker(b *testing.B) {
	for i := 0; i < b.N; i++ {
		w := NewWorker(context.Background())
		if err := w.Do(ProcFunc(func(_ *Next) error {
			return nil
		})); err != nil {
			b.Fatal(err)
		}
		if err := w.Close(); err != nil {
			b.Fatal(err)
		}
	}
}
