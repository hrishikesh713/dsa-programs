package algos

import (
	"testing"
	"time"
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
