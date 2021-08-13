package concurrency

import (
	"sync"
	"sync/atomic"
)

// AtomicBool supports safe concurrency when setting, checking, and waiting for a boolean value.
type AtomicBool struct {
	value int32
	c     *sync.Cond
	m     *sync.Mutex
}

// Set sets the boolean to given value.
func (ab *AtomicBool) Set(ok bool) {
	var val int32
	if ok {
		val = 1
	} else {
		val = 0
	}

	atomic.StoreInt32(&ab.value, val)
	ab.c.Broadcast()
}

// Value returns the current status.
func (ab *AtomicBool) Value() bool {
	return atomic.LoadInt32(&ab.value)&1 == 1
}

// WaitForTrue causes the caller to block until the AtomicBool is set to true.
func (ab *AtomicBool) WaitForTrue() bool {
	ab.c.L.Lock()
	for !ab.Value() {
		ab.c.Wait()
	}
	ab.c.L.Unlock()

	return true
}

// NewAtomicBool returns a ready for use AtomicBool.
func NewAtomicBool() *AtomicBool {
	m := &sync.Mutex{}
	c := sync.NewCond(m)
	return &AtomicBool{
		m: m,
		c: c,
	}
}
