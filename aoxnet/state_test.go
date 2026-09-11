package aoxnet

import "testing"

func TestConnectionStateMachine(t *testing.T) {
	s := NewStateMachine()
	if s.State() != StateNew {
		t.Fatal("unexpected initial state")
	}
	if err := s.BeginHandshake(); err != nil {
		t.Fatal(err)
	}
	if err := s.Open(); err != nil {
		t.Fatal(err)
	}
	if err := s.GoAway(); err != nil {
		t.Fatal(err)
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	if s.State() != StateClosed {
		t.Fatal("connection did not close")
	}
}

func TestInvalidConnectionTransition(t *testing.T) {
	s := NewStateMachine()
	if err := s.Open(); err == nil {
		t.Fatal("opening before handshake should fail")
	}
}
