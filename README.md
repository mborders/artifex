[![GoDoc](http://godoc.org/github.com/borderstech/artifex?status.png)](http://godoc.org/github.com/borderstech/artifex)
[![Build Status](https://travis-ci.org/borderstech/artifex.svg?branch=master)](https://travis-ci.org/borderstech/artifex)
[![Go Report Card](https://goreportcard.com/badge/github.com/borderstech/artifex)](https://goreportcard.com/report/github.com/borderstech/artifex)
[![codecov](https://codecov.io/gh/borderstech/artifex/branch/master/graph/badge.svg)](https://codecov.io/gh/borderstech/artifex)

# artifex

Simple in-memory job queue for Golang using worker-based dispatching

Documentation here: https://godoc.org/github.com/borderstech/artifex

## Example Usage

```go
// 10 workers, 100 max in job queue
d := artifex.NewDispatcher(10, 100)
d.Run()

d.Dispatch(func() {
  // do something
})

err := d.DispatchIn(func() {
  // do something in 500ms
}, time.Millisecond*500)

// Returns a DispatchTicker
dt, err := d.DispatchEvery(func() {
  // do something every 250ms
}, time.Millisecond*250)

// Stop a given DispatchTicker
dt.Stop()

// Stop a dispatcher and all its workers/tickers
d.Stop()
```
