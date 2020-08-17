[![GoDoc](http://godoc.org/github.com/mborders/artifex?status.png)](http://godoc.org/github.com/mborders/artifex)
[![Build Status](https://travis-ci.org/mborders/artifex.svg?branch=master)](https://travis-ci.org/mborders/artifex)
[![Go Report Card](https://goreportcard.com/badge/github.com/mborders/artifex)](https://goreportcard.com/report/github.com/mborders/artifex)
[![codecov](https://codecov.io/gh/mborders/artifex/branch/master/graph/badge.svg)](https://codecov.io/gh/mborders/artifex)

# artifex

Simple in-memory job queue for Golang using worker-based dispatching

Documentation here: https://godoc.org/github.com/mborders/artifex

Cron jobs use the robfig/cron library: https://godoc.org/github.com/robfig/cron

## Example Usage

```go
// 10 workers, 100 max in job queue
d := artifex.NewDispatcher(10, 100)
d.Start()

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

// Returns a DispatchCron
dc, err := d.DispatchCron(func() {
  // do something every 1s
}, "*/1 * * * * *")

// Stop a given DispatchCron
dc.Stop()

// Stop a dispatcher and all its workers/tickers
d.Stop()
```
