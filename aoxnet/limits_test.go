package aoxnet

import "testing"

func TestProtocolLimitsRejectInvalidInitialWindow(t *testing.T) {
    m := NewStreamManager(MaxWindow + 1)
    if m.initialWindow != DefaultInitialWindow { t.Fatalf("invalid initial window was not normalized") }
}

func TestRemoteStreamParity(t *testing.T) {
    m := NewStreamManager(DefaultInitialWindow)
    if _, err := m.RegisterRemote(1); err == nil { t.Fatal("odd remote stream ID should be rejected by the current role-independent manager") }
    if _, err := m.RegisterRemote(2); err != nil { t.Fatalf("even remote stream ID rejected: %v", err) }
}
