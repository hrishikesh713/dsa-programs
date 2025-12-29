package algos

import (
	"testing"
	"testing/synctest"
	"time"
)

func TestFixedWindowRateLimiterSingleTenant(t *testing.T) {
	windowDuration := time.Minute
	limiter := NewFixedWindowRateLimiter(windowDuration)
	tests := []struct {
		TestName        string
		TenantID        string
		RequestLimit    int
		NewRequestLimit int
	}{
		{
			TestName:        "Tenant foo",
			TenantID:        "foo",
			RequestLimit:    600,
			NewRequestLimit: 899,
		},
		{
			TestName:        "Tenant bar",
			TenantID:        "bar",
			RequestLimit:    300,
			NewRequestLimit: 499,
		},
	}

	for _, tt := range tests {
		synctest.Test(t, func(t *testing.T) {
			limiter.SetRateLimit(tt.TenantID, tt.RequestLimit)
			for i := range tt.RequestLimit {
				allowed, _ := limiter.Allow(tt.TenantID)
				if !allowed {
					t.Errorf("%s: Request %d not allowed", tt.TestName, i)
				}
			}
			allowed, _ := limiter.Allow(tt.TenantID)
			if allowed {
				t.Errorf("%s: Request should not be allowed", tt.TestName)
			}
			time.Sleep(windowDuration)
			for i := range tt.RequestLimit {
				allowed, _ := limiter.Allow(tt.TenantID)
				if !allowed {
					t.Errorf("%s: Request %d not allowed after window reset", tt.TestName, i)
				}
			}
			limiter.SetRateLimit(tt.TenantID, tt.NewRequestLimit)
			time.Sleep(windowDuration)
			for i := range tt.NewRequestLimit {
				allowed, _ := limiter.Allow(tt.TenantID)
				if !allowed {
					t.Errorf("%s: Request %d not allowed after limit change", tt.TestName, i)
				}
			}
		})
	}
}
