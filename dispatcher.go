package queue

import "fmt"

type Dispatcher struct {
	WorkerPool chan chan Job
	JobQueue   chan Job
	maxWorkers int
}

func NewDispatcher(maxWorkers int, maxQueue int) *Dispatcher {
	return &Dispatcher{
		WorkerPool: make(chan chan Job, maxWorkers),
		JobQueue:   make(chan Job, maxQueue),
		maxWorkers: maxWorkers,
	}
}

func (d *Dispatcher) Run() {
	for i := 0; i < d.maxWorkers; i++ {
		worker := NewWorker(fmt.Sprintf("%d", i), d.WorkerPool)
		worker.Start()
	}

	go d.start()
}

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
