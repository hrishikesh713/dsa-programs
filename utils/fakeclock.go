package utils

import (
	"sync"
	"time"
)

type Clocker interface {
	Now() time.Time
	Advance(d time.Duration) time.Time
	Sleep(d time.Duration)
	Sub(t time.Time) time.Duration
}

type RealClock struct{}

func NewRealClock() *RealClock {
	return new(RealClock)
}

func (rc *RealClock) Now() time.Time {
	return time.Now()
}

func (rc *RealClock) Sleep(d time.Duration) {
	time.Sleep(d)
}

func (rc *RealClock) Advance(d time.Duration) time.Time {
	time.Sleep(d)
	return time.Now()
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

func (rc *RealClock) Sub(t time.Time) time.Duration {
	return time.Since(t)
}

func (fc *FakeClock) Sleep(d time.Duration) {
	<-fc.after(d)
}

func (fc *FakeClock) Sub(t time.Time) time.Duration {
	return fc.Curr.Sub(t)
}

func (fc *FakeClock) after(d time.Duration) <-chan time.Time {
	fc.mu.Lock()
	defer fc.mu.Unlock()
	ft := FakeTimer{Notify: make(chan time.Time, 1), SchedTime: fc.Curr.Add(d)}
	fc.Timers = append(fc.Timers, ft)
	return ft.Notify
}

func (fc *FakeClock) Advance(d time.Duration) time.Time {
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
	return fc.Curr
}
