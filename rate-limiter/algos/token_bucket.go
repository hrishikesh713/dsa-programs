package algos

import (
	"fmt"
	"time"

	"github.com/hrishikesh713/dsa-programs/utils"
)

type TokenBucketRateLimiter struct {
	Rate        int // tokens/second
	Bursts      int // max number of tokens allowed
	WindowStart time.Time
	State       map[string]*tenantState
	Clock       utils.Clocker
}

type tenantState struct {
	tokens         int
	LastRefillTime time.Time
}

type Request struct {
	TenantID string
	Value    any
}

func NewTokenBucketRateLimiter(rate int, bursts int, fakeClock bool) (*TokenBucketRateLimiter, error) {
	if fakeClock {
		return &TokenBucketRateLimiter{Rate: rate, Bursts: bursts, State: make(map[string]*tenantState), Clock: utils.NewFakeClock()}, nil
	}
	return &TokenBucketRateLimiter{Rate: rate, Bursts: bursts, State: make(map[string]*tenantState), Clock: utils.NewRealClock()}, nil
}

func (tb *TokenBucketRateLimiter) RefillTokens(tenantID string) error {
	ts, ok := tb.State[tenantID]
	if !ok {
		return fmt.Errorf("tenant not found :  %s", tenantID)
	}
	now := tb.Clock.Now()
	elapsed := int(now.Sub(ts.LastRefillTime).Seconds())

	tokensToRefill := min((ts.tokens + (tb.Rate * elapsed)), tb.Bursts)
	ts.tokens = tokensToRefill
	ts.LastRefillTime = now
	return nil
}

func (tb *TokenBucketRateLimiter) Allow(r *Request) (bool, error) {
	_, ok := tb.State[r.TenantID]
	if !ok {
		return false, fmt.Errorf("cannot find tenantID %s", r.TenantID)
	}
	_ = tb.RefillTokens(r.TenantID)
	t := tb.State[r.TenantID]
	if t.tokens <= 0 {
		return false, nil
	}
	tb.State[r.TenantID].tokens--
	return true, nil
}

func (tb *TokenBucketRateLimiter) SetupTokenBucket(tenantID string) error {
	if tb.State == nil {
		return fmt.Errorf("token bucket not initialized")
	}
	tb.State[tenantID] = &tenantState{tokens: tb.Rate, LastRefillTime: tb.Clock.Now()}
	return nil
}
