package oldworld

import (
	"context"
	"sync"
)

type Worker struct {
	quitFn context.CancelFunc
	wg     sync.WaitGroup
	ijobs  chan func()
	ojobs  chan func()
}

// NewWorker creates a new worker that can hold maxJobs number of jobs
// at a given time. The worker's persistent goroutine is started.
func NewWorker(bufSize int) *Worker {
	ctx, quitFn := context.WithCancel(context.Background())
	w := Worker{
		// TODO ask tyler about buffering here
		ijobs:  make(chan func(), bufSize),
		ojobs:  make(chan func(), bufSize),
		quitFn: quitFn,
		wg:     sync.WaitGroup{},
	}
	w.wg.Add(2)
	go w.start(ctx)
	go w.bufferJobs(ctx)
	return &w
}

func (w *Worker) start(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			w.wg.Done()
			return
		case job, ok := <-w.ojobs:
			if !ok {
				panic("read on closed channel")
			}
			job()
		}
	}
}

// bufferJobs buffers jobs and preserves FIFO order.
func (w *Worker) bufferJobs(ctx context.Context) {
	// TODO use linked list instead, see stdlib container/list
	var buf []func()
	// use dummy job
	var j func()
	// use dummy ojobs
	ojobs := make(chan func())
	for {
		select {
		case <-ctx.Done():
			w.wg.Done()
			return
		case inJob, ok := <-w.ijobs:
			if !ok {
				panic("read from closed channel")
			}
			if j == nil {
				// set up proper job
				j = inJob
				// set up proper ojobs
				ojobs = w.ojobs
			} else {
				buf = append(buf, inJob)
			}
		case ojobs <- j:
			if len(buf) == 0 {
				// use dummy job
				j = nil
				// use dummy ojobs
				ojobs = make(chan func())
			} else {
				j, buf = buf[0], buf[1:]
			}
		}
	}
}

// Quit will cancel the worker's persistent goroutines.
func (w *Worker) Quit() {
	w.quitFn()
	w.wg.Wait()
	close(w.ijobs)
	close(w.ojobs)
}

// AddJob gives a job to the worker.
func (w *Worker) AddJob(fn func()) {
	// TODO gaurantee that read from ijobs will prioritize over ojobs write in bufferJobs()
	w.ijobs <- fn
}
