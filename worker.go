package queue

// Worker attaches to a provided worker pool, and
// looks for jobs on its job channel
type Worker struct {
	ID         string
	WorkerPool chan chan Job
	JobChannel chan Job
	quit       chan bool
}

// NewWorker creates a new worker using the given ID and
// attaches to the provided worker pool. It also initializes
// the job/quit channels
func NewWorker(id string, workerPool chan chan Job) Worker {
	return Worker{
		ID:         id,
		WorkerPool: workerPool,
		JobChannel: make(chan Job),
		quit:       make(chan bool)}
}

// Start initializes a select loop to listen for jobs to execute
func (w Worker) Start() {
	go func() {
		for {
			w.WorkerPool <- w.JobChannel

			select {
			case job := <-w.JobChannel:
				job.Run()
			case <-w.quit:
				return
			}
		}
	}()
}

// Stop will end the job select loop for the worker
func (w Worker) Stop() {
	go func() {
		w.quit <- true
	}()
}
