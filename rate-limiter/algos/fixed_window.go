package algos

import "time"

type FixedWindowRateLimiter struct {
	VolumeLimit    int
	WindowDuration time.Duration
	State          map[string]*RateLimit
}

type RateLimit struct {
	Counter     int
	WindowStart time.Time
}

func NewFixedWindowRateLimiter(volumeLimit int, windowDuration time.Duration) *FixedWindowRateLimiter {
	return &FixedWindowRateLimiter{VolumeLimit: volumeLimit, WindowDuration: windowDuration, State: make(map[string]*RateLimit)}
}

func (fw *FixedWindowRateLimiter) Allow(tenantID string) bool {
	now := time.Now()
	nowFloored := now.Truncate(fw.WindowDuration)
	if _, ok := fw.State[tenantID]; !ok {
		s := RateLimit{WindowStart: nowFloored, Counter: 1}
		fw.State[tenantID] = &s
		return true
	}
	rl := fw.State[tenantID]
	if now.Sub(rl.WindowStart) >= fw.WindowDuration {
		rl.WindowStart = nowFloored
		rl.Counter = 1
		return true
	}
	if rl.Counter >= fw.VolumeLimit {
		return false
	}
	rl.Counter++
	return true
}
