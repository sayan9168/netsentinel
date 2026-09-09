package aoxnet

import (
    "errors"
    "sync"
)

type ConnState uint8

const (
    StateNew ConnState = iota
    StateHandshaking
    StateOpen
    StateGoingAway
    StateClosed
)

var ErrInvalidTransition = errors.New("invalid connection state transition")

type StateMachine struct {
    mu sync.RWMutex
    state ConnState
}

func NewStateMachine() *StateMachine { return &StateMachine{state: StateNew} }
func (s *StateMachine) State() ConnState { s.mu.RLock(); defer s.mu.RUnlock(); return s.state }

func (s *StateMachine) BeginHandshake() error {
    s.mu.Lock(); defer s.mu.Unlock()
    if s.state != StateNew { return ErrInvalidTransition }
    s.state = StateHandshaking
    return nil
}
func (s *StateMachine) Open() error {
    s.mu.Lock(); defer s.mu.Unlock()
    if s.state != StateHandshaking { return ErrInvalidTransition }
    s.state = StateOpen
    return nil
}
func (s *StateMachine) GoAway() error {
    s.mu.Lock(); defer s.mu.Unlock()
    if s.state != StateOpen { return ErrInvalidTransition }
    s.state = StateGoingAway
    return nil
}
func (s *StateMachine) Close() error {
    s.mu.Lock(); defer s.mu.Unlock()
    if s.state == StateClosed { return nil }
    s.state = StateClosed
    return nil
}
