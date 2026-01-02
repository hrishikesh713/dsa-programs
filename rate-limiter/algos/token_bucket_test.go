package algos

import (
	"fmt"
	"testing"
	"time"

	"github.com/hrishikesh713/dsa-programs/utils"
)

func TestTokenBucketInitialization(t *testing.T) {
	// TODO:
	// 1. add slog
	tokens := 1000
	bursts := 10000
	_, err := NewTokenBucketRateLimiter(tokens, bursts, true)
	if err != nil {
		t.Errorf("cannot initialize token bucket with %d tokens", tokens)
	}
}

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

func NewTokenBucketRateLimiter(rate int, bursts int, fakeClock bool) (*TokenBucketRateLimiter, error) {
	if fakeClock {
		return &TokenBucketRateLimiter{Rate: rate, Bursts: bursts, State: make(map[string]*tenantState), Clock: utils.NewFakeClock()}, nil
	}
	return &TokenBucketRateLimiter{Rate: rate, Bursts: bursts, State: make(map[string]*tenantState), Clock: utils.NewRealClock()}, nil
}

type Request struct {
	TenantID string
	Value    any
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

func TestTokenBucketAllow(t *testing.T) {
	tenantID := "foo"
	tb, _ := NewTokenBucketRateLimiter(100, 1000, true)
	err := tb.SetupTokenBucket(tenantID)
	if err != nil {
		t.Fatalf("cannot set up Token bucket")
	}
	req := Request{TenantID: tenantID, Value: "bar"}
	for i := range 99 {
		allowed, allowedErr := tb.Allow(&req)
		if allowedErr != nil {
			t.Errorf("cannot use the rate limiter for tenantID %s", req.TenantID)
		}
		if !allowed {
			t.Errorf("request %d should have been allowed", i)
		}
	}
}

func TestTokenBucketFull(t *testing.T) {
	tenantID := "foo"
	tb, _ := NewTokenBucketRateLimiter(100, 1000, true)
	err := tb.SetupTokenBucket(tenantID)
	if err != nil {
		t.Fatalf("cannot set up Token bucket")
	}
	req := Request{TenantID: "foo", Value: "bar"}
	for i := range 100 {
		allowed, allowedErr := tb.Allow(&req)
		if allowedErr != nil {
			t.Errorf("cannot use the rate limiter for tenantID %s", req.TenantID)
		}
		if !allowed {
			t.Errorf("request %d should have been allowed", i)
		}
	}
	tb.Clock.Advance(time.Duration(99) * time.Millisecond)
	if allowed, _ := tb.Allow(&req); allowed {
		t.Errorf("this request should not have been allowed")
	}
}

func TestTokenBucketMultiTenantAllow(t *testing.T) {
	tenantIDBar := "bar"
	tenantIDFoo := "foo"
	tb, _ := NewTokenBucketRateLimiter(100, 1000, true)
	tenants := []string{"foo", "bar"}
	for _, tenant := range tenants {
		if err := tb.SetupTokenBucket(tenant); err != nil {
			t.Errorf("error while setting up tenant %s", tenant)
		}
	}
	reqFoo := Request{TenantID: tenantIDFoo, Value: "bar"}
	reqBar := Request{TenantID: tenantIDBar, Value: "baz"}
	for i := range 99 {
		allowed, allowedErr := tb.Allow(&reqFoo)
		if allowedErr != nil {
			t.Errorf("cannot use the rate limiter for tenantID %s", reqFoo.TenantID)
		}
		if !allowed {
			t.Errorf("request %d should have been allowed for tenant %s", i, reqFoo.TenantID)
		}
	}
	for i := range 99 {
		allowed, allowedErr := tb.Allow(&reqBar)
		if allowedErr != nil {
			t.Errorf("cannot use the rate limiter for tenantID %s", reqBar.TenantID)
		}
		if !allowed {
			t.Errorf("request %d should have been allowed for tenant %s", i, reqBar.TenantID)
		}
	}
}
