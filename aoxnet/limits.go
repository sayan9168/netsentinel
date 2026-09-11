package aoxnet

import (
	"fmt"
	"sync"
	"time"
)

// Limits protects a connection from unbounded protocol resource usage.
type Limits struct {
	MaxPayload       uint32
	MaxFrameRate     uint32
	RateWindow       time.Duration
	IdleTimeout      time.Duration
	HandshakeTimeout time.Duration
	InitialWindow    int64
}

func (l Limits) normalize() Limits {
	if l.MaxPayload == 0 {
		l.MaxPayload = DefaultMaxPayload
	}
	if l.MaxFrameRate == 0 {
		l.MaxFrameRate = 1000
	}
	if l.RateWindow == 0 {
		l.RateWindow = time.Second
	}
	if l.IdleTimeout == 0 {
		l.IdleTimeout = 2 * time.Minute
	}
	if l.HandshakeTimeout == 0 {
		l.HandshakeTimeout = 10 * time.Second
	}
	if l.InitialWindow <= 0 || l.InitialWindow > MaxWindow {
		l.InitialWindow = DefaultInitialWindow
	}
	return l
}

type rateLimiter struct {
	mu     sync.Mutex
	start  time.Time
	count  uint32
	limit  uint32
	window time.Duration
}

func newRateLimiter(limit uint32, window time.Duration) *rateLimiter {
	return &rateLimiter{start: time.Now(), limit: limit, window: window}
}

func (r *rateLimiter) allow() error {
	r.mu.Lock()
	defer r.mu.Unlock()
	now := time.Now()
	if now.Sub(r.start) >= r.window {
		r.start = now
		r.count = 0
	}
	if r.count >= r.limit {
		return fmt.Errorf("%w: frame rate exceeded", ErrRateLimited)
	}
	r.count++
	return nil
}
