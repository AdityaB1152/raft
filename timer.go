package main

import (
	"time"
)

func NewTimer() *Timer {
	return &Timer{
		C: make(chan struct{}, 1),
	}
}

type Timer struct {
	C      chan struct{}
	Timer  *time.Timer
	Period time.Duration
}

func (t *Timer) Start(period time.Duration) {
	t.Period = period
	t.Timer = time.AfterFunc(t.Period, func() {
		t.C <- struct{}{}
	})
}

func (t *Timer) Reset() {
	if t.Timer != nil {
		t.Timer.Stop()
	}
	t.Start(t.Period)
}

func (t *Timer) Stop() {
	if t.Timer != nil {
		t.Timer.Stop()
	}
}

func (t *Timer) OnTimeout(f func()) {
	go func() {
		for range t.C {
			f()
		}
	}()
}
