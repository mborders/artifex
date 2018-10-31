package queue

// Job represents a runnable process, where Run
// will be executed by a worker via the dispatch queue
type Job struct {
	Run func()
}
