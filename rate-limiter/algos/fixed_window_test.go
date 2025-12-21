package algos

import (
	"testing"
	"testing/synctest"
	"time"
)

func TestFixedWindowRateLimiterSingleTenant(t *testing.T) {
	requestLimit := 600
	windowDuration := time.Minute
	limiter := NewFixedWindowRateLimiter(requestLimit, windowDuration)
	tenantID := "foo"

	t.Run("allow requests within limit", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			for i := range requestLimit {
				allowed := limiter.Allow(tenantID)
				if !allowed {
					t.Errorf("Request %d not allowed", i)
				}
			}
		})
	})

	t.Run("denies requests excedding limit", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			allowed := limiter.Allow(tenantID)
			if allowed {
				t.Error("Request should not be allowed")
			}
		})
	})

	t.Run("resets after window duration", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			time.Sleep(windowDuration)
			allowed := limiter.Allow(tenantID)
			if !allowed {
				t.Error("request should have been allowed when new window is rolled over")
			}
		})
	})
}

// ============================================================================
// BONUS CHALLENGES (Try these after completing the basic test!)
// ============================================================================
//
// 1. Test multiple users/clients with separate rate limits
// 2. Test concurrent requests using goroutines
// 3. Test partial window consumption (use 2 requests, wait, use 3 more)
// 4. Test edge cases (0 limit, negative duration, etc.)
//
// ============================================================================
