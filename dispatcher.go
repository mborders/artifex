package queue

import (
	"fmt"
	"time"
)

// Dispatcher maintains a pool for available workers
// and a job queue that workers will process
type Dispatcher struct {
	WorkerPool chan chan Job
	JobQueue   chan Job
	maxWorkers int
}

// NewDispatcher creates a new dispatcher with the given
// number of works and buffers the job queue based on maxQueue.
// It also initializes the channels for the worker pool and job queue
func NewDispatcher(maxWorkers int, maxQueue int) *Dispatcher {
	return &Dispatcher{
		WorkerPool: make(chan chan Job, maxWorkers),
		JobQueue:   make(chan Job, maxQueue),
		maxWorkers: maxWorkers,
	}
}

// Run creates and starts workers, adding them to the worker pool.
// Then, it starts a select loop to wait for job to be dispatched
// to available workers
func (d *Dispatcher) Run() {
	for i := 0; i < d.maxWorkers; i++ {
		worker := NewWorker(fmt.Sprintf("%d", i), d.WorkerPool)
		worker.Start()
	}

	go d.start()
}

// Dispatch pushes the given job into the job queue.
// The first available worker will perform the job
func (d *Dispatcher) Dispatch(job Job) {
	d.JobQueue <- job
}

// DispatchIn pushes the given job into the job queue
// after the given duration has elapsed
func (d *Dispatcher) DispatchIn(job Job, duration time.Duration) {
	go func() {
		time.Sleep(duration)
		d.JobQueue <- job
	}()
}

// DispatchEvery pushes the given job into the job queue
// continuously at the given interval
func (d *Dispatcher) DispatchEvery(job Job, interval time.Duration) *DispatchTicker {
	t := time.NewTicker(interval)
	dt := &DispatchTicker{ticker: t, quit: make(chan bool)}

	go func() {
		for {
			select {
			case <-t.C:
				d.JobQueue <- job
			case <-dt.quit:
				return
			}
		}
	}()

	return dt
}

func (d *Dispatcher) start() {
	for {
		select {
		case job := <-d.JobQueue:
			go func(job Job) {
				jobChannel := <-d.WorkerPool
				jobChannel <- job
			}(job)
		}
	}
}

// DispatchTicker represents a dispatched job ticker
// that executes on a given interval. This provides
// a means for stopping the execution cycle from continuing.
type DispatchTicker struct {
	ticker *time.Ticker
	quit   chan bool
}

// Stop ends the execution cycle for the given ticker.
func (dt *DispatchTicker) Stop() {
	dt.ticker.Stop()
	dt.quit <- true
}
