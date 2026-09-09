package aoxnet

import (
	"errors"
	"fmt"
	"sync"
)

const (
	FlagEndStream uint16 = 1 << iota
	FlagEndMessage
	FlagWindowUpdate
	FlagProtocolError
)

type StreamState uint8

const (
	StreamIdle StreamState = iota
	StreamOpen
	StreamHalfClosedLocal
	StreamHalfClosedRemote
	StreamClosed
)

type Stream struct {
	ID       uint32
	State    StreamState
	RecvWin  int64
	SendWin  int64
	mu       sync.Mutex
}

func NewStream(id uint32, initialWindow int64) *Stream {
	if initialWindow <= 0 { initialWindow = DefaultInitialWindow }
	return &Stream{ID: id, State: StreamIdle, RecvWin: initialWindow, SendWin: initialWindow}
}

func (s *Stream) Open() error {
	s.mu.Lock(); defer s.mu.Unlock()
	if s.State != StreamIdle { return errors.New("stream is not idle") }
	s.State = StreamOpen
	return nil
}

func (s *Stream) ConsumeSendWindow(n int64) error {
	if n < 0 { return errors.New("negative window consumption") }
	s.mu.Lock(); defer s.mu.Unlock()
	if s.State == StreamClosed { return errors.New("stream is closed") }
	if n > s.SendWin { return fmt.Errorf("flow-control window exceeded: need=%d available=%d", n, s.SendWin) }
	s.SendWin -= n
	return nil
}

func (s *Stream) AddSendWindow(n int64) error {
	if n <= 0 { return errors.New("window increment must be positive") }
	s.mu.Lock(); defer s.mu.Unlock()
	if s.SendWin > MaxWindow-n { return errors.New("flow-control window overflow") }
	s.SendWin += n
	return nil
}

func (s *Stream) ConsumeRecvWindow(n int64) error {
	if n < 0 { return errors.New("negative receive consumption") }
	s.mu.Lock(); defer s.mu.Unlock()
	if n > s.RecvWin { return fmt.Errorf("receive window exceeded: need=%d available=%d", n, s.RecvWin) }
	s.RecvWin -= n
	return nil
}

func (s *Stream) AddRecvWindow(n int64) error {
	if n <= 0 { return errors.New("window increment must be positive") }
	s.mu.Lock(); defer s.mu.Unlock()
	if s.RecvWin > MaxWindow-n { return errors.New("receive window overflow") }
	s.RecvWin += n
	return nil
}

func (s *Stream) CloseLocal() error {
	s.mu.Lock(); defer s.mu.Unlock()
	switch s.State {
	case StreamOpen: s.State = StreamHalfClosedLocal
	case StreamHalfClosedRemote: s.State = StreamClosed
	default: return errors.New("invalid local close transition")
	}
	return nil
}

func (s *Stream) CloseRemote() error {
	s.mu.Lock(); defer s.mu.Unlock()
	switch s.State {
	case StreamOpen: s.State = StreamHalfClosedRemote
	case StreamHalfClosedLocal: s.State = StreamClosed
	default: return errors.New("invalid remote close transition")
	}
	return nil
}
