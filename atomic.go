package concurrency

import (
	"context"
	"log"
	"sync"
	"sync/atomic"
)

// AtomicBool supports safe concurrency when setting, checking, and waiting for a boolean value.
type AtomicBool struct {
	value int32

	wait chan struct{}
	waitLock *sync.RWMutex
}

// Set sets the boolean to given value.
func (ab *AtomicBool) Set(ok bool) {
	if ok {
		ab.setTrue()
	} else {
		ab.setFalse()
	}
}

func (ab *AtomicBool) setTrue() {
	if ab.Value() {
		return
	}

	ab.waitLock.Lock()
	defer ab.waitLock.Unlock()
	var val int32 = 1
	atomic.StoreInt32(&ab.value, val)
	close(ab.wait)
	ab.wait = make(chan struct{})
}

func (ab *AtomicBool) setFalse() {
	if !ab.Value() {
		return
	}

	var val int32 = 0
	atomic.StoreInt32(&ab.value, val)
}

// Value returns the current status.
func (ab *AtomicBool) Value() bool {
	return atomic.LoadInt32(&ab.value)&1 == 1
}

func (ab *AtomicBool) WaitWithContext(ctx context.Context) bool {
	ab.waitLock.RLock()
	ch := ab.wait
	ab.waitLock.RUnlock()
	for {
		select {
		case <-ctx.Done():
			log.Println("ctx done")
			return ab.Value()
		case <-ch:
			log.Println("received from wait chan")
			return ab.Value()
		}
	}
}

// Wait causes the caller to block until the AtomicBool is set to true.
func (ab *AtomicBool) Wait() bool {
	return ab.WaitWithContext(context.Background())
}

// NewAtomicBool returns a ready for use AtomicBool.
func NewAtomicBool() *AtomicBool {
	return &AtomicBool{
		waitLock: &sync.RWMutex{},
		wait: make(chan struct{}),
	}
}
