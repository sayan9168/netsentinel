package aoxnet

import (
	"errors"
	"fmt"
	"sync"
)

const (
	DefaultInitialWindow int64 = 256 * 1024
	MaxWindow            int64 = 1<<31 - 1
)

type StreamManager struct {
	mu             sync.RWMutex
	streams        map[uint32]*Stream
	nextLocalID    uint32
	initialWindow  int64
	connectionRecv int64
	connectionSend int64
	localIsClient  bool
	closed         bool
	goAway         bool
}

// NewStreamManager preserves the client-side odd/local, even/remote default.
func NewStreamManager(initialWindow int64) *StreamManager {
	return NewStreamManagerWithRole(initialWindow, true)
}

// NewStreamManagerWithRole configures stream ID parity for the local endpoint.
// Clients use odd local IDs; servers use even local IDs.
func NewStreamManagerWithRole(initialWindow int64, localIsClient bool) *StreamManager {
	if initialWindow <= 0 || initialWindow > MaxWindow {
		initialWindow = DefaultInitialWindow
	}
	next := uint32(2)
	if localIsClient {
		next = 1
	}
	return &StreamManager{
		streams:        make(map[uint32]*Stream),
		nextLocalID:    next,
		initialWindow:  initialWindow,
		connectionRecv: initialWindow,
		connectionSend: initialWindow,
		localIsClient:  localIsClient,
	}
}

// SetRole is intended for an endpoint before it opens or accepts streams.
func (m *StreamManager) SetRole(localIsClient bool) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.closed || len(m.streams) != 0 {
		return errors.New("cannot change stream role after streams are active")
	}
	m.localIsClient = localIsClient
	if localIsClient {
		m.nextLocalID = 1
	} else {
		m.nextLocalID = 2
	}
	return nil
}

func (m *StreamManager) OpenLocal() (*Stream, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.closed || m.goAway {
		return nil, errors.New("connection is shutting down")
	}
	id := m.nextLocalID
	m.nextLocalID += 2
	if id == 0 {
		return nil, errors.New("stream ID exhausted")
	}
	s := NewStream(id, m.initialWindow)
	if err := s.Open(); err != nil {
		return nil, err
	}
	m.streams[id] = s
	return s, nil
}

func (m *StreamManager) Get(id uint32) (*Stream, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	s, ok := m.streams[id]
	return s, ok
}

// RegisterRemote accepts the parity opposite to the local endpoint.
func (m *StreamManager) RegisterRemote(id uint32) (*Stream, error) {
	if id == 0 {
		return nil, errors.New("stream ID zero is reserved")
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	expectedEven := m.localIsClient
	if (id%2 == 0) != expectedEven {
		return nil, errors.New("remote stream ID has invalid parity")
	}
	if m.closed || m.goAway {
		return nil, errors.New("connection is shutting down")
	}
	if _, exists := m.streams[id]; exists {
		return nil, errors.New("stream already exists")
	}
	s := NewStream(id, m.initialWindow)
	if err := s.Open(); err != nil {
		return nil, err
	}
	m.streams[id] = s
	return s, nil
}

func (m *StreamManager) Remove(id uint32) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.streams, id)
}

func (m *StreamManager) ConsumeConnectionSend(n int64) error {
	if n < 0 {
		return errors.New("negative window consumption")
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	if n > m.connectionSend {
		return fmt.Errorf("connection flow-control window exceeded: need=%d available=%d", n, m.connectionSend)
	}
	m.connectionSend -= n
	return nil
}

func (m *StreamManager) AddConnectionSend(n int64) error {
	if n <= 0 {
		return errors.New("window increment must be positive")
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	if n > MaxWindow || m.connectionSend > MaxWindow-n {
		return errors.New("connection window overflow")
	}
	m.connectionSend += n
	return nil
}

func (m *StreamManager) ConsumeConnectionRecv(n int64) error {
	if n < 0 {
		return errors.New("negative window consumption")
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	if n > m.connectionRecv {
		return fmt.Errorf("connection receive window exceeded: need=%d available=%d", n, m.connectionRecv)
	}
	m.connectionRecv -= n
	return nil
}

func (m *StreamManager) AddConnectionRecv(n int64) error {
	if n <= 0 {
		return errors.New("window increment must be positive")
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	if n > MaxWindow || m.connectionRecv > MaxWindow-n {
		return errors.New("connection window overflow")
	}
	m.connectionRecv += n
	return nil
}

func (m *StreamManager) BeginGoAway() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.goAway = true
}

func (m *StreamManager) Close() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.closed = true
	for _, s := range m.streams {
		s.mu.Lock()
		s.State = StreamClosed
		s.mu.Unlock()
	}
}
