// ratelimiter/fixed_window_test.go
package ratelimiter

import (
	"testing"
	"time"

	"github.com/neelp03/throttlex/store"
)

func TestFixedWindowLimiter(t *testing.T) {
	memStore := store.NewMemoryStore()
	limiter, err := NewFixedWindowLimiter(memStore, 5, time.Second*1)
	if err != nil {
		t.Errorf("Failed to create rate limiter: %v", err)
	}
	key := "user1"

	// Simulate 5 allowed requests
	for i := 0; i < 5; i++ {
		allowed, err := limiter.Allow(key)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if !allowed {
			t.Errorf("Request %d should be allowed", i+1)
		}
	}

	// 6th request should be blocked
	allowed, err := limiter.Allow(key)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if allowed {
		t.Errorf("6th request should not be allowed")
	}

	// Wait for the window to expire
	time.Sleep(time.Second * 1)

	// Next request should be allowed after window resets
	allowed, err = limiter.Allow(key)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if !allowed {
		t.Errorf("Request after window reset should be allowed")
	}

	// Wait for the window to expire
	time.Sleep(time.Second * 1)

	// Edge case: negative limit
	_, err = NewFixedWindowLimiter(memStore, -1, time.Second)
	if err == nil {
		t.Error("Expected error when limit is set to a negative value")
	}

	// Edge case: zero window duration
	_, err = NewFixedWindowLimiter(memStore, 5, 0)
	if err == nil {
		t.Error("Expected error when window duration is set to zero")
	}
}

func TestFixedWindowLimiter_MultipleIPs(t *testing.T) {
	memStore := store.NewMemoryStore()
	limiter, err := NewFixedWindowLimiter(memStore, 10, time.Second*5)
	if err != nil {
		t.Fatalf("Failed to create FixedWindowLimiter: %v", err)
	}

	// Define a set of IPs to test
	ips := []string{
		"192.168.1.1", "192.168.1.2", "192.168.1.3",
		"192.168.1.4", "192.168.1.5", "192.168.1.6",
		"192.168.1.7", "192.168.1.8", "192.168.1.9",
		"192.168.1.10",
	}

	// Allow 10 requests for each IP and ensure each IP is allowed up to the limit
	for _, ip := range ips {
		for i := 0; i < 10; i++ {
			allowed, err := limiter.Allow(ip)
			if err != nil {
				t.Errorf("Unexpected error for IP %s: %v", ip, err)
			}
			if !allowed {
				t.Errorf("Request %d should be allowed for IP %s", i+1, ip)
			}
		}
		// 11th request should be blocked for each IP
		allowed, err := limiter.Allow(ip)
		if err != nil {
			t.Errorf("Unexpected error for IP %s: %v", ip, err)
		}
		if allowed {
			t.Errorf("11th request should not be allowed for IP %s", ip)
		}
	}

	// Wait for the window to expire, then check again
	time.Sleep(time.Second * 5)
	for _, ip := range ips {
		allowed, err := limiter.Allow(ip)
		if err != nil {
			t.Errorf("Unexpected error after window reset for IP %s: %v", ip, err)
		}
		if !allowed {
			t.Errorf("Request after window reset should be allowed for IP %s", ip)
		}
	}

	// Edge cases:
	// 1. Check behavior with an empty IP string
	allowed, err := limiter.Allow("")
	if err == nil || allowed {
		t.Errorf("Expected error or disallowed access for empty IP")
	}

	// 2. Check behavior with non-standard IP format
	invalidIP := "invalidIP"
	allowed, err = limiter.Allow(invalidIP)
	if err != nil {
		t.Errorf("Unexpected error for invalid IP %s: %v", invalidIP, err)
	}
	if !allowed {
		t.Errorf("Expected first request to be allowed for invalid IP format %s", invalidIP)
	}

	// 3. Check large number of requests for a single IP beyond the threshold
	singleIP := "192.168.2.1"
	for i := 0; i < 20; i++ {
		allowed, err := limiter.Allow(singleIP)
		if err != nil {
			t.Errorf("Unexpected error for IP %s: %v", singleIP, err)
		}
		if i < 10 && !allowed {
			t.Errorf("Request %d should be allowed for IP %s", i+1, singleIP)
		}
		if i >= 10 && allowed {
			t.Errorf("Request %d should be denied for IP %s after limit is reached", i+1, singleIP)
		}
	}
}

func TestFixedWindowLimiter_InvalidKey(t *testing.T) {
	memStore := store.NewMemoryStore()
	limiter, err := NewFixedWindowLimiter(memStore, 5, time.Second*1)
	if err != nil {
		t.Fatalf("Failed to create FixedWindowLimiter: %v", err)
	}

	// Test cases for invalid keys
	testCases := []struct {
		key     string
		name    string
		wantErr bool
	}{
		{"", "Empty Key", true},
		{"invalidKey!", "Invalid Format Key", true},
		{
			"thisisaverylongkeythatshouldnotbeallowedbecauseitexceedsthe256characterlimitandweareusingitjustfortestpurposesbutintheenditshoulddefinitelytriggeranerrorsoheregoesthelongstringwiththemaximumlengthofcharactersallowedbytheapplicationaddedextracharacterstomakeitsureitisover256characters",
			"Overly Long Key", true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			allowed, err := limiter.Allow(tc.key)
			if tc.wantErr {
				if err == nil {
					t.Errorf("Expected error for %s but got none", tc.name)
				}
				if allowed {
					t.Errorf("Expected disallowed access for %s, but got allowed", tc.name)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error for %s: %v", tc.name, err)
				}
				if !allowed {
					t.Errorf("Expected allowed access for %s, but got disallowed", tc.name)
				}
			}
		})
	}
}
