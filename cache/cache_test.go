package cache

import (
	"testing"
	"time"
)

func TestCache(t *testing.T) {
	now := time.Now()
	c := New(time.Second)
	var tests = []struct {
		key       string
		value     interface{}
		ok        bool
		ttl       time.Duration
		nowOffset time.Duration
	}{
		{"k1", 1, true, time.Minute, 0},
		{"k2", 2, true, time.Minute, time.Minute},
		{"k3", nil, false, time.Minute, -time.Second * 61},
		{"k4", nil, false, time.Second * 5, -time.Second * 6},
	}
	for i, tt := range tests {
		c.now = func() time.Time { return now.Add(tt.nowOffset) }
		c.Set(tt.key, tt.value, tt.ttl)
		c.now = func() time.Time { return now }
		v, ok := c.Get(tt.key)
		if ok != tt.ok || v != tt.value {
			t.Errorf("#%d: Get(%q) = (%v, %t), want (%v, %t)", i, tt.key, v, ok, tt.value, tt.ok)
		}
	}
	c.evictExpired()
	if got, want := c.Len(), 2; got != want {
		t.Errorf("Len() = %d, want %d", got, want)
	}
}
