package queue

import "fmt"

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
