package artifex

import (
	"errors"
	"fmt"
	"github.com/robfig/cron/v3"
	"time"
)

// Dispatcher maintains a pool for available workers
// and a job queue that workers will process
type Dispatcher struct {
	maxWorkers int
	maxQueue   int
	workers    []*Worker
	tickers    []*DispatchTicker
	crons      []*DispatchCron
	workerPool chan chan Job
	jobQueue   chan Job
	quit       chan bool
	active     bool
}

// NewDispatcher creates a new dispatcher with the given
// number of workers and buffers the job queue based on maxQueue.
// It also initializes the channels for the worker pool and job queue
func NewDispatcher(maxWorkers int, maxQueue int) *Dispatcher {
	return &Dispatcher{
		maxWorkers: maxWorkers,
		maxQueue:   maxQueue,
	}
}

// Start creates and starts workers, adding them to the worker pool.
// Then, it starts a select loop to wait for job to be dispatched
// to available workers
func (d *Dispatcher) Start() {
	d.workers = []*Worker{}
	d.tickers = []*DispatchTicker{}
	d.crons = []*DispatchCron{}
	d.workerPool = make(chan chan Job, d.maxWorkers)
	d.jobQueue = make(chan Job, d.maxQueue)
	d.quit = make(chan bool)

	for i := 0; i < d.maxWorkers; i++ {
		worker := NewWorker(d.workerPool)
		worker.Start()
		d.workers = append(d.workers, worker)
	}

	d.active = true

	go func() {
		for {
			select {
			case job := <-d.jobQueue:
				go func(job Job) {
					jobChannel := <-d.workerPool
					jobChannel <- job
				}(job)
			case <-d.quit:
				return
			}
		}
	}()
}

// Stop ends execution for all workers/tickers and
// closes all channels, then removes all workers/tickers
func (d *Dispatcher) Stop() {
	if !d.active {
		return
	}

	d.active = false

	for i := range d.workers {
		d.workers[i].Stop()
	}

	for i := range d.tickers {
		d.tickers[i].Stop()
	}

	for i := range d.crons {
		d.crons[i].Stop()
	}

	d.workers = []*Worker{}
	d.tickers = []*DispatchTicker{}
	d.crons = []*DispatchCron{}
	d.quit <- true
}

// Dispatch pushes the given job into the job queue.
// The first available worker will perform the job
func (d *Dispatcher) Dispatch(run func()) error {
	if !d.active {
		return errors.New("dispatcher is not active")
	}

	d.jobQueue <- Job{Run: run}
	return nil
}

// DispatchIn pushes the given job into the job queue
// after the given duration has elapsed
func (d *Dispatcher) DispatchIn(run func(), duration time.Duration) error {
	if !d.active {
		return errors.New("dispatcher is not active")
	}

	go func() {
		time.Sleep(duration)
		d.jobQueue <- Job{Run: run}
	}()

	return nil
}

// DispatchEvery pushes the given job into the job queue
// continuously at the given interval
func (d *Dispatcher) DispatchEvery(run func(), interval time.Duration) (*DispatchTicker, error) {
	if !d.active {
		return nil, errors.New("dispatcher is not active")
	}

	t := time.NewTicker(interval)
	dt := &DispatchTicker{ticker: t, quit: make(chan bool)}
	d.tickers = append(d.tickers, dt)

	go func() {
		for {
			select {
			case <-t.C:
				d.jobQueue <- Job{Run: run}
			case <-dt.quit:
				return
			}
		}
	}()

	return dt, nil
}

// DispatchEvery pushes the given job into the job queue
// each time the cron definition is met
func (d *Dispatcher) DispatchCron(run func(), cronStr string) (*DispatchCron, error) {
	if !d.active {
		return nil, errors.New("dispatcher is not active")
	}

	dc := &DispatchCron{cron: cron.New(cron.WithSeconds())}
	d.crons = append(d.crons, dc)

	_, err := dc.cron.AddFunc(cronStr, func() {
		d.jobQueue <- Job{Run: run}
	})

	if err != nil {
		return nil, errors.New("invalid cron definition")
	}

	dc.cron.Start()
	return dc, nil
}

// DispatchBatch pushes multiple given jobs into the job queue
// optionally separated from each other by a given time gap
func (d *Dispatcher) DispatchBatch(runs []func(), gap time.Duration) error {
	for !d.active {
		return errors.New("dispatcher is not active")
	}

	for i, run := range runs {
		if err := d.DispatchIn(run, time.Duration(int64(i)*int64(gap))); err != nil {
			return err
		}
	}

	return nil
}

// DispatchCronBatch pushes a job into the queue to dispatch a batch of jobs once the given definition is met.
// Jobs can optionally be separated from each other by a given time gap.
func (d *Dispatcher) DispatchCronBatch(runs []func(), cronStr string, gap time.Duration) (*DispatchCron, error) {
	for !d.active {
		return nil, errors.New("dispatcher is not active")
	}

	return d.DispatchCron(func() {
		if err := d.DispatchBatch(runs, gap); err != nil {
			fmt.Println(fmt.Errorf("failed to dispatch batch of jobs, %w", err))
		}
	}, cronStr)
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

// DispatchCron represents a dispatched cron job
// that executes using cron expression formats.
type DispatchCron struct {
	cron *cron.Cron
}

// Stops ends the execution cycle for the given cron.
func (c *DispatchCron) Stop() {
	c.cron.Stop()
}
