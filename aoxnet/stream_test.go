package aoxnet

import "testing"

func TestStreamFlowControl(t *testing.T) {
	s := NewStream(1, 100)
	if err := s.Open(); err != nil {
		t.Fatal(err)
	}
	if err := s.ConsumeSendWindow(60); err != nil {
		t.Fatal(err)
	}
	if err := s.ConsumeSendWindow(50); err == nil {
		t.Fatal("expected flow-control error")
	}
	if err := s.AddSendWindow(20); err != nil {
		t.Fatal(err)
	}
	if err := s.ConsumeSendWindow(60); err != nil {
		t.Fatal(err)
	}
}

func TestStreamRejectsLocalWritesAfterClose(t *testing.T) {
	s := NewStream(1, 100)
	if err := s.Open(); err != nil {
		t.Fatal(err)
	}
	if err := s.CloseLocal(); err != nil {
		t.Fatal(err)
	}
	if err := s.ConsumeSendWindow(1); err == nil {
		t.Fatal("expected send rejection after local close")
	}
}

func TestStreamCloseTransitions(t *testing.T) {
	s := NewStream(1, 100)
	if err := s.Open(); err != nil {
		t.Fatal(err)
	}
	if err := s.CloseLocal(); err != nil {
		t.Fatal(err)
	}
	if s.State != StreamHalfClosedLocal {
		t.Fatalf("unexpected state: %v", s.State)
	}
	if err := s.CloseRemote(); err != nil {
		t.Fatal(err)
	}
	if s.State != StreamClosed {
		t.Fatalf("unexpected state: %v", s.State)
	}
}

func TestStreamManager(t *testing.T) {
	m := NewStreamManager(100)
	s, err := m.OpenLocal()
	if err != nil {
		t.Fatal(err)
	}
	if s.ID != 1 {
		t.Fatalf("expected first local stream 1, got %d", s.ID)
	}
	m.BeginGoAway()
	if _, err := m.OpenLocal(); err == nil {
		t.Fatal("expected GOAWAY rejection")
	}
}

func TestStreamManagerServerParity(t *testing.T) {
	m := NewStreamManagerWithRole(100, false)
	local, err := m.OpenLocal()
	if err != nil {
		t.Fatal(err)
	}
	if local.ID != 2 {
		t.Fatalf("expected first server-local stream 2, got %d", local.ID)
	}
	if _, err := m.RegisterRemote(1); err != nil {
		t.Fatal(err)
	}
	if _, err := m.RegisterRemote(2); err == nil {
		t.Fatal("expected invalid remote parity")
	}
}
