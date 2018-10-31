package queue

import (
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
	"time"
)

func TestQueue(t *testing.T) {
	a := 0
	b := 0
	c := 0

	d := NewDispatcher(10, 3)
	d.Run()

	d.Dispatch(Job{Run: func() {
		a = 1
	}})

	d.Dispatch(Job{Run: func() {
		b = 2
	}})

	d.Dispatch(Job{Run: func() {
		c = 3
	}})

	time.Sleep(time.Second * 3)
	assert.Equal(t, 1, a)
	assert.Equal(t, 2, b)
	assert.Equal(t, 3, c)
}

func TestQueue_Mutex(t *testing.T) {
	n := 100
	mutex := &sync.Mutex{}

	d := NewDispatcher(10, 3)
	d.Run()

	var v []int

	for i := 0; i < n; i++ {
		d.Dispatch(Job{Run: func() {
			mutex.Lock()
			v = append(v, 0)
			mutex.Unlock()
		}})
	}

	time.Sleep(time.Second * 3)
	assert.Equal(t, n, len(v))
}