package artifex

// Job represents a runnable process, where Start
// will be executed by a worker via the dispatch queue
type Job struct {
	Run func()
}
