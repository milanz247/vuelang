// Package ratelimit provides a simple token-bucket rate limiter keyed by IP.
// It is purposefully dependency-free and safe for concurrent use.
package ratelimit

import (
	"sync"
	"time"
)

// Limiter tracks request counts per key within a sliding window.
type Limiter struct {
	mu      sync.Mutex
	buckets map[string]*bucket
	max     int
	window  time.Duration
}

type bucket struct {
	count     int
	resetAt   time.Time
}

// New creates a Limiter that allows max requests per window duration.
func New(max int, window time.Duration) *Limiter {
	l := &Limiter{
		buckets: make(map[string]*bucket),
		max:     max,
		window:  window,
	}
	go l.cleanup()
	return l
}

// Allow returns true if the given key is within its request budget.
func (l *Limiter) Allow(key string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()
	b, ok := l.buckets[key]
	if !ok || now.After(b.resetAt) {
		l.buckets[key] = &bucket{count: 1, resetAt: now.Add(l.window)}
		return true
	}
	if b.count >= l.max {
		return false
	}
	b.count++
	return true
}

// cleanup sweeps expired buckets every 5 minutes to prevent unbounded growth.
func (l *Limiter) cleanup() {
	ticker := time.NewTicker(5 * time.Minute)
	for range ticker.C {
		now := time.Now()
		l.mu.Lock()
		for k, b := range l.buckets {
			if now.After(b.resetAt) {
				delete(l.buckets, k)
			}
		}
		l.mu.Unlock()
	}
}
