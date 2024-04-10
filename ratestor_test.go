package ratestor

import (
	"testing"
	"time"
)

func TestAllow(t *testing.T) {
	rs := NewRateStor()
	defer rs.Close()

	if err := rs.Allow("key1", time.Millisecond, 2); err != nil {
		t.Errorf("Expected no error, but got %v", err)
	}

	if err := rs.Allow("key1", time.Millisecond, 2); err != nil {
		t.Errorf("Expected no error, but got %v", err)
	}

	if err := rs.Allow("key1", time.Millisecond, 2); err != ErrRateLimitExceeded {
		t.Errorf("Expected error %v, but got %v", ErrRateLimitExceeded, err)
	}

	// Test case 2: Allow after expiration
	time.Sleep(2 * time.Millisecond)

	if err := rs.Allow("key1", time.Millisecond, 1); err != nil {
		t.Errorf("Expected no error, but got %v", err)
	}
}

func TestGCRun(t *testing.T) {
	rs := NewRateStor(WithGCInterval(2*time.Millisecond), WithGCBatchSize(10))
	defer rs.Close()

	_ = rs.Allow("key1", 3*time.Millisecond, 1)

	rs.lock.Lock()

	if _, ok := rs.rates["key1"]; !ok {
		t.Errorf("Expected key1 to be present in the map")
	}
	rs.lock.Unlock()

	time.Sleep(2 * time.Millisecond)

	rs.lock.Lock()

	if _, ok := rs.rates["key1"]; !ok {
		t.Errorf("Expected key1 to be present in the map")
	}
	rs.lock.Unlock()

	time.Sleep(4 * time.Millisecond)

	rs.lock.Lock()

	if _, ok := rs.rates["key1"]; ok {
		t.Errorf("Expected key1 to be removed from the map")
	}

	rs.lock.Unlock()
}

func TestClose(t *testing.T) {
	rs := NewRateStor()
	rs.Close()

	if err := rs.Allow("key", time.Millisecond, 1); err != ErrRateStorClosed {
		t.Errorf("Expected error %v, but got %v", ErrRateStorClosed, err)
	}
}
