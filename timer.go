package main

import (
	"sync"
	"time"
)

type Timer struct {
	timer    *time.Timer
	randTime time.Duration
	mu       sync.Mutex
	callback func()
}

func NewTimer() *Timer {
	return &Timer{}
}

func (t *Timer) Start(randTime time.Duration) {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.randTime = randTime
	t.timer = time.AfterFunc(randTime, func() {
		if t.callback != nil {
			t.callback()
		}
	})
}

func (t *Timer) Stop() {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.timer != nil {
		t.timer.Stop()
	}
}

func (t *Timer) Reset() {
	t.Stop()
	t.Start(t.randTime)
}

func (t *Timer) OnTimeout(callback func()) {
	t.callback = callback
}
