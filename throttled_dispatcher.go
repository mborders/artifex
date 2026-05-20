package artifex

import "time"

// ThrottleDispatcher wraps a Dispatcher with rate-limiting behavior,
// allowing at most maxJobs to be dispatched per interval.
type ThrottleDispatcher struct {
	*Dispatcher
	maxJobs  int
	interval time.Duration
	tokens   chan struct{}
	quit     chan bool
}

// NewThrottledDispatcher creates a new dispatcher that limits job execution
// to at most maxJobs per interval (e.g., 5 jobs per time.Minute).
func NewThrottledDispatcher(maxWorkers int, maxQueue int, maxJobs int, interval time.Duration) *ThrottleDispatcher {
	tokens := make(chan struct{}, maxJobs)
	// Pre-fill the token bucket
	for i := 0; i < maxJobs; i++ {
		tokens <- struct{}{}
	}
	return &ThrottleDispatcher{
		Dispatcher: NewDispatcher(maxWorkers, maxQueue),
		maxJobs:    maxJobs,
		interval:   interval,
		tokens:     tokens,
		quit:       make(chan bool),
	}
}

// Start begins the throttled dispatcher and starts the token replenishment ticker.
func (td *ThrottleDispatcher) Start() {
	td.Dispatcher.Start()
	go td.replenish()
}

// Stop halts the throttled dispatcher and stops token replenishment.
func (td *ThrottleDispatcher) Stop() {
	td.quit <- true
	td.Dispatcher.Stop()
}

// replenish refills the token bucket at the configured rate.
// Tokens are added in batches of maxJobs every interval.
func (td *ThrottleDispatcher) replenish() {
	ticker := time.NewTicker(td.interval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			for i := 0; i < td.maxJobs; i++ {
				select {
				case td.tokens <- struct{}{}:
				default:
					// Bucket is full; discard extra tokens
				}
			}
		case <-td.quit:
			return
		}
	}
}

// Dispatch waits for a throttle token before queuing the job,
// ensuring no more than maxJobs are dispatched per interval.
func (td *ThrottleDispatcher) Dispatch(run func()) error {
	<-td.tokens // Block until a token is available
	return td.Dispatcher.Dispatch(run)
}
