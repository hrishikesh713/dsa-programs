// Package algos implements various rate limiting algorithms.
package algos

import (
	"fmt"
	"time"
)

type FixedWindowRateLimiter struct {
	WindowDuration time.Duration
	State          map[string]*RateLimit
}

type RateLimit struct {
	Counter     int
	WindowStart time.Time
	VolumeLimit int
}

func NewFixedWindowRateLimiter(windowDuration time.Duration) *FixedWindowRateLimiter {
	return &FixedWindowRateLimiter{WindowDuration: windowDuration, State: make(map[string]*RateLimit)}
}

func (fw *FixedWindowRateLimiter) SetRateLimit(tenantID string, volumeLimit int) {
	if v, ok := fw.State[tenantID]; ok {
		v.VolumeLimit = volumeLimit
		return
	}
	rl := RateLimit{WindowStart: time.Now().Truncate(fw.WindowDuration), Counter: 0, VolumeLimit: volumeLimit}
	fw.State[tenantID] = &rl
}

func (fw *FixedWindowRateLimiter) Allow(tenantID string) (bool, error) {
	now := time.Now()
	nowFloored := now.Truncate(fw.WindowDuration)
	if _, ok := fw.State[tenantID]; !ok {
		return false, fmt.Errorf("tenantID %s not found", tenantID)
	}
	rl := fw.State[tenantID]
	if now.Sub(rl.WindowStart) >= fw.WindowDuration {
		rl.WindowStart = nowFloored
		rl.Counter = 1
		return true, nil
	}
	if rl.Counter >= rl.VolumeLimit {
		return false, nil
	}
	rl.Counter++
	return true, nil
}
