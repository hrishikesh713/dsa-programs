package utils

import (
	"sync"
	"time"
)

type Clock interface {
	Now() time.Time
	Advance(d time.Duration)
	Sleep(d time.Duration)
	After(d time.Duration) <-chan time.Time
}

type FakeTimer struct {
	SchedTime time.Time
	Notify    chan time.Time
}

type FakeClock struct {
	mu     sync.Mutex
	Curr   time.Time
	Timers []FakeTimer
}

func NewFakeClock() *FakeClock {
	return new(FakeClock)
}

func (fc *FakeClock) Now() time.Time {
	fc.mu.Lock()
	defer fc.mu.Unlock()
	return fc.Curr
}

func (fc *FakeClock) Sleep(d time.Duration) {
	<-fc.After(d)
}

func (fc *FakeClock) After(d time.Duration) <-chan time.Time {
	fc.mu.Lock()
	defer fc.mu.Unlock()
	ft := FakeTimer{Notify: make(chan time.Time, 1), SchedTime: fc.Curr.Add(d)}
	fc.Timers = append(fc.Timers, ft)
	return ft.Notify
}

func (fc *FakeClock) Advance(d time.Duration) {
	fc.mu.Lock()
	defer fc.mu.Unlock()
	fc.Curr = fc.Curr.Add(d)
	remaining := make([]FakeTimer, 0, len(fc.Timers))
	for _, v := range fc.Timers {
		if !v.SchedTime.After(fc.Curr) {
			select {
			case v.Notify <- fc.Curr:
			default:
			}
		} else {
			remaining = append(remaining, v)
		}
	}
	fc.Timers = remaining
}
