package aoxnet

import (
    "context"
    "sync"
    "time"
)

// KeepaliveConfig defines bounded peer liveness checks.
type KeepaliveConfig struct {
    Interval  time.Duration
    Timeout   time.Duration
    MaxMissed int
}

func (k KeepaliveConfig) normalize() KeepaliveConfig {
    if k.Interval <= 0 { k.Interval = 30 * time.Second }
    if k.Timeout <= 0 { k.Timeout = 10 * time.Second }
    if k.MaxMissed <= 0 { k.MaxMissed = 3 }
    return k
}

type KeepaliveState struct {
    mu sync.Mutex
    pending map[uint64]time.Time
    missed int
}

func NewKeepaliveState() *KeepaliveState { return &KeepaliveState{pending: make(map[uint64]time.Time)} }
func (s *KeepaliveState) Sent(id uint64, now time.Time) { s.mu.Lock(); s.pending[id] = now; s.mu.Unlock() }
func (s *KeepaliveState) Ack(id uint64) { s.mu.Lock(); delete(s.pending, id); s.missed = 0; s.mu.Unlock() }
func (s *KeepaliveState) Expired(now time.Time, timeout time.Duration) int {
    s.mu.Lock(); defer s.mu.Unlock()
    n := 0
    for id, sent := range s.pending { if now.Sub(sent) >= timeout { delete(s.pending, id); n++ } }
    s.missed += n
    return s.missed
}

// RunKeepalive provides transport-independent peer liveness orchestration.
func RunKeepalive(ctx context.Context, cfg KeepaliveConfig, probe func() (uint64, error), state *KeepaliveState, onDead func()) {
    cfg = cfg.normalize()
    ticker := time.NewTicker(cfg.Interval)
    defer ticker.Stop()
    for {
        select {
        case <-ctx.Done(): return
        case now := <-ticker.C:
            if state.Expired(now, cfg.Timeout) >= cfg.MaxMissed { onDead(); return }
            if id, err := probe(); err == nil { state.Sent(id, now) }
        }
    }
}
