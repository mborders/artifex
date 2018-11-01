[![GoDoc](http://godoc.org/github.com/mborders/job-queue?status.png)](http://godoc.org/github.com/mborders/job-queue)
[![Build Status](https://travis-ci.org/mborders/job-queue.svg?branch=master)](https://travis-ci.org/mborders/job-queue)
[![Go Report Card](https://goreportcard.com/badge/github.com/mborders/job-queue)](https://goreportcard.com/report/github.com/mborders/job-queue)
[![codecov](https://codecov.io/gh/mborders/artifex/branch/master/graph/badge.svg)](https://codecov.io/gh/mborders/artifex)

# artifex

Simple in-memory job queue for Golang using worker-based dispatching

Documentation here: https://godoc.org/github.com/mborders/job-queue

## Example Usage

```go
// 10 workers, 100 max in job queue
d := artifex.NewDispatcher(10, 100)
d.Run()

d.Dispatch(artifex.Job{Run: func() {
  // do something
}})

err := d.DispatchIn(artifex.Job{Run: func() {
  // do something in 500ms
}}, time.Millisecond*500)

// Returns a DispatchTicker
dt, err := d.DispatchEvery(artifex.Job{Run: func() {
  // do something every 250ms
}}, time.Millisecond*250)

// Stop a given DispatchTicker
dt.Stop()

// Stop a dispatcher and all its workers
d.Stop()
```
