package ratestor

import (
	"testing"
	"time"
)

func TestAllow(t *testing.T) {
	rs := NewRateStor()
	defer rs.Close()

	err := rs.Allow("key1", time.Millisecond, 2)
	if err != nil {
		t.Errorf("Expected no error, but got %v", err)
	}

	err = rs.Allow("key1", time.Millisecond, 2)
	if err != nil {
		t.Errorf("Expected no error, but got %v", err)
	}

	err = rs.Allow("key1", time.Millisecond, 2)
	if err != ErrRateLimitExceeded {
		t.Errorf("Expected error %v, but got %v", ErrRateLimitExceeded, err)
	}

	// Test case 2: Allow after expiration
	time.Sleep(2 * time.Millisecond)
	err = rs.Allow("key1", time.Millisecond, 1)
	if err != nil {
		t.Errorf("Expected no error, but got %v", err)
	}
}

func TestGCRun(t *testing.T) {
	rs := NewRateStor(WithGCInterval(2 * time.Millisecond))
	defer rs.Close()

	rs.Allow("key1", 3*time.Millisecond, 1)

	rs.lock.Lock()
	_, ok := rs.rates["key1"]
	if !ok {
		t.Errorf("Expected key1 to be present in the map")
	}
	rs.lock.Unlock()

	time.Sleep(2 * time.Millisecond)

	rs.lock.Lock()
	_, ok = rs.rates["key1"]
	if !ok {
		t.Errorf("Expected key1 to be present in the map")
	}
	rs.lock.Unlock()

	time.Sleep(2 * time.Millisecond)

	rs.lock.Lock()
	_, ok = rs.rates["key1"]
	if ok {
		t.Errorf("Expected key1 to be removed from the map")
	}
	rs.lock.Unlock()
}
