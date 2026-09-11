package aoxnet

import "testing"

func TestNegotiateFeatures(t *testing.T) {
	got := NegotiateFeatures([]string{"request-id", "ping", "custom"}, []string{"ping", "request-id"})
	want := []string{"request-id", "ping"}
	if len(got) != len(want) {
		t.Fatalf("got %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("got %v, want %v", got, want)
		}
	}
}

func TestDefaultFeatures(t *testing.T) {
	features := DefaultFeatures()
	if len(features) < 7 {
		t.Fatalf("expected core feature set, got %v", features)
	}
	if !HasFeature(features, FeatureMultiplexing) || !HasFeature(features, FeatureFlowControl) {
		t.Fatalf("missing required transport features: %v", features)
	}
}

func TestMetricsSnapshot(t *testing.T) {
	var m Metrics
	m.AddFrameIn(10)
	m.AddFrameOut(20)
	m.AddProtocolError()
	m.StreamOpened()
	m.AddKeepaliveMiss()

	s := m.Snapshot()
	if s.FramesIn != 1 || s.BytesIn != 10 || s.FramesOut != 1 || s.BytesOut != 20 {
		t.Fatalf("unexpected frame metrics: %+v", s)
	}
	if s.ProtocolErrors != 1 || s.ActiveStreams != 1 || s.KeepaliveMisses != 1 {
		t.Fatalf("unexpected counters: %+v", s)
	}
}
