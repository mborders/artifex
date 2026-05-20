package artifex

import (
	"sync/atomic"
	"testing"
	"time"
)

func TestThrottledDispatcher_LimitsRate(t *testing.T) {
	// Allow 3 jobs per 500ms window
	td := NewThrottledDispatcher(5, 100, 3, 500*time.Millisecond)
	td.Start()
	defer td.Stop()

	var count int32

	// Dispatch 6 jobs; first 3 should run immediately, next 3 after ~500ms
	for i := 0; i < 6; i++ {
		go func() {
			td.Dispatch(func() {
				atomic.AddInt32(&count, 1)
			})
		}()
	}

	// After a short wait, only first batch should have been dispatched
	time.Sleep(100 * time.Millisecond)
	firstBatch := atomic.LoadInt32(&count)
	if firstBatch > 3 {
		t.Errorf("expected at most 3 jobs in first window, got %d", firstBatch)
	}

	// Wait for second window to open
	time.Sleep(600 * time.Millisecond)
	total := atomic.LoadInt32(&count)
	if total != 6 {
		t.Errorf("expected 6 total jobs, got %d", total)
	}
}
