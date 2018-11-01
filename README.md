[![GoDoc](http://godoc.org/github.com/alphaterra/artifex?status.png)](http://godoc.org/github.com/alphaterra/artifex)
[![Build Status](https://travis-ci.org/alphaterra/artifex.svg?branch=master)](https://travis-ci.org/alphaterra/artifex)
[![Go Report Card](https://goreportcard.com/badge/github.com/alphaterra/artifex)](https://goreportcard.com/report/github.com/alphaterra/artifex)
[![codecov](https://codecov.io/gh/alphaterra/artifex/branch/master/graph/badge.svg)](https://codecov.io/gh/alphaterra/artifex)

# artifex

Simple in-memory job queue for Golang using worker-based dispatching

Documentation here: https://godoc.org/github.com/alphaterra/artifex

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
