package aoxnet

import "sync/atomic"

// Metrics contains lock-free protocol counters suitable for embedding in telemetry.
type Metrics struct {
	FramesIn        uint64
	FramesOut       uint64
	BytesIn         uint64
	BytesOut        uint64
	ProtocolErrors  uint64
	ActiveStreams   int64
	KeepaliveMisses uint64
}

func (m *Metrics) AddFrameIn(bytes int) {
	atomic.AddUint64(&m.FramesIn, 1)
	atomic.AddUint64(&m.BytesIn, uint64(bytes))
}
func (m *Metrics) AddFrameOut(bytes int) {
	atomic.AddUint64(&m.FramesOut, 1)
	atomic.AddUint64(&m.BytesOut, uint64(bytes))
}
func (m *Metrics) AddProtocolError() { atomic.AddUint64(&m.ProtocolErrors, 1) }
func (m *Metrics) StreamOpened()     { atomic.AddInt64(&m.ActiveStreams, 1) }
func (m *Metrics) StreamClosed() {
	if atomic.AddInt64(&m.ActiveStreams, -1) < 0 {
		atomic.StoreInt64(&m.ActiveStreams, 0)
	}
}
func (m *Metrics) AddKeepaliveMiss() { atomic.AddUint64(&m.KeepaliveMisses, 1) }

// MetricsSnapshot is an immutable point-in-time view of protocol counters.
type MetricsSnapshot struct {
	FramesIn        uint64
	FramesOut       uint64
	BytesIn         uint64
	BytesOut        uint64
	ProtocolErrors  uint64
	ActiveStreams   int64
	KeepaliveMisses uint64
}

// Snapshot returns atomically loaded metrics suitable for logs, dashboards, or exports.
func (m *Metrics) Snapshot() MetricsSnapshot {
	return MetricsSnapshot{
		FramesIn:        atomic.LoadUint64(&m.FramesIn),
		FramesOut:       atomic.LoadUint64(&m.FramesOut),
		BytesIn:         atomic.LoadUint64(&m.BytesIn),
		BytesOut:        atomic.LoadUint64(&m.BytesOut),
		ProtocolErrors:  atomic.LoadUint64(&m.ProtocolErrors),
		ActiveStreams:   atomic.LoadInt64(&m.ActiveStreams),
		KeepaliveMisses: atomic.LoadUint64(&m.KeepaliveMisses),
	}
}
