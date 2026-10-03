package artifex

import (
	"sync"
	"testing"
	"time"
)

func TestDispatchTickerStop(t *testing.T) {
	tests := []struct {
		name string
		stop func(*Dispatcher, *DispatchTicker)
	}{
		{
			name: "repeated ticker stop",
			stop: func(d *Dispatcher, dt *DispatchTicker) {
				dt.Stop()
				dt.Stop()
				d.Stop()
			},
		},
		{
			name: "ticker before dispatcher",
			stop: func(d *Dispatcher, dt *DispatchTicker) {
				dt.Stop()
				d.Stop()
			},
		},
		{
			name: "dispatcher before ticker",
			stop: func(d *Dispatcher, dt *DispatchTicker) {
				d.Stop()
				dt.Stop()
			},
		},
		{
			name: "concurrent ticker and dispatcher stops",
			stop: func(d *Dispatcher, dt *DispatchTicker) {
				var wg sync.WaitGroup
				start := make(chan struct{})
				for i := 0; i < 8; i++ {
					wg.Add(1)
					go func() {
						defer wg.Done()
						<-start
						dt.Stop()
					}()
				}
				wg.Add(1)
				go func() {
					defer wg.Done()
					<-start
					d.Stop()
				}()
				close(start)
				wg.Wait()
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := NewDispatcher(1, 1)
			d.Start()
			dt, err := d.DispatchEvery(func() {}, time.Hour)
			if err != nil {
				t.Fatal(err)
			}

			done := make(chan struct{})
			go func() {
				tt.stop(d, dt)
				close(done)
			}()

			select {
			case <-done:
			case <-time.After(time.Second):
				t.Fatal("Stop did not return")
			}
		})
	}
}
