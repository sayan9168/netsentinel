package aoxnet

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

type Hello struct { Name string `json:"name"`; Version byte `json:"version"`; Features []string `json:"features"` }
type Welcome struct { Name string `json:"name"`; Version byte `json:"version"`; Features []string `json:"features"` }

type ServerConfig struct {
	MaxPayload uint32
	ReadTimeout, WriteTimeout time.Duration
	Name string
	Limits Limits
}
type ClientConfig struct {
	MaxPayload uint32
	ReadTimeout, WriteTimeout time.Duration
	Name string
	Timeout time.Duration
	Limits Limits
}

type Conn struct {
	net.Conn
	r *bufio.Reader
	maxPayload uint32
	limits Limits
	rate *rateLimiter
	writeMu sync.Mutex
	nextRequest uint64
	ReadTimeout, WriteTimeout time.Duration
}

func NewConn(c net.Conn, maxPayload uint32) *Conn {
	limits := Limits{MaxPayload: maxPayload}.normalize()
	return newConnWithLimits(c, limits)
}

func newConnWithLimits(c net.Conn, limits Limits) *Conn {
	limits = limits.normalize()
	return &Conn{Conn: c, r: bufio.NewReader(c), maxPayload: limits.MaxPayload, limits: limits, rate: newRateLimiter(limits.MaxFrameRate, limits.RateWindow)}
}

func (c *Conn) send(f Frame) error {
	c.writeMu.Lock()
	defer c.writeMu.Unlock()
	if c.WriteTimeout > 0 { _ = c.SetWriteDeadline(time.Now().Add(c.WriteTimeout)) }
	return f.Encode(c.Conn, c.maxPayload)
}

func (c *Conn) receive() (Frame, error) {
	if err := c.rate.allow(); err != nil { return Frame{}, err }
	timeout := c.ReadTimeout
	if timeout == 0 { timeout = c.limits.IdleTimeout }
	if timeout > 0 { _ = c.SetReadDeadline(time.Now().Add(timeout)) }
	return Decode(c.r, c.maxPayload)
}

func (c *Conn) nextID() uint64 { return atomic.AddUint64(&c.nextRequest, 1) }

func (c *Conn) sendError(code ErrorCode, message string, requestID uint64) error {
	payload, err := json.Marshal(ProtocolError{Code: code, Message: message})
	if err != nil { return err }
	return c.send(Frame{Type: TypeError, RequestID: requestID, Payload: payload})
}

func (c *Conn) HandshakeClient(cfg ClientConfig) error {
	if c.limits.HandshakeTimeout > 0 { _ = c.SetDeadline(time.Now().Add(c.limits.HandshakeTimeout)) }
	defer c.SetDeadline(time.Time{})
	hello, err := json.Marshal(Hello{Name: cfg.Name, Version: Version, Features: []string{"request-id", "crc32", "ping", "tls-transport"}})
	if err != nil { return err }
	id := c.nextID()
	if err := c.send(Frame{Type: TypeHello, RequestID: id, Payload: hello}); err != nil { return err }
	f, err := c.receive()
	if err != nil { return err }
	if f.Type == TypeError { return decodeProtocolError(f.Payload) }
	if f.Type != TypeWelcome { return fmt.Errorf("%w: expected WELCOME", ErrProtocolState) }
	var w Welcome
	if err := json.Unmarshal(f.Payload, &w); err != nil { return fmt.Errorf("%w: invalid WELCOME payload: %v", ErrInvalidFrame, err) }
	if w.Version != Version { return fmt.Errorf("%w: server selected unsupported version %d", ErrUnsupportedVersion, w.Version) }
	return nil
}

func (c *Conn) HandshakeServer(name string) error {
	if c.limits.HandshakeTimeout > 0 { _ = c.SetDeadline(time.Now().Add(c.limits.HandshakeTimeout)) }
	defer c.SetDeadline(time.Time{})
	f, err := c.receive()
	if err != nil { return err }
	if f.Type != TypeHello {
		_ = c.sendError(ErrProtocolState, "expected HELLO", f.RequestID)
		return fmt.Errorf("%w: expected HELLO", ErrProtocolState)
	}
	var h Hello
	if err := json.Unmarshal(f.Payload, &h); err != nil {
		_ = c.sendError(ErrInvalidFrame, "invalid HELLO payload", f.RequestID)
		return err
	}
	if h.Version != Version {
		_ = c.sendError(ErrUnsupportedVersion, fmt.Sprintf("unsupported client version %d", h.Version), f.RequestID)
		return fmt.Errorf("%w: %d", ErrUnsupportedVersion, h.Version)
	}
	payload, err := json.Marshal(Welcome{Name: name, Version: Version, Features: []string{"request-id", "crc32", "ping", "tls-transport"}})
	if err != nil { return err }
	return c.send(Frame{Type: TypeWelcome, RequestID: f.RequestID, Payload: payload})
}

func decodeProtocolError(payload []byte) error {
	var pe ProtocolError
	if err := json.Unmarshal(payload, &pe); err != nil { return fmt.Errorf("%w: malformed ERROR frame", ErrInvalidFrame) }
	if pe.Message == "" { pe.Message = "remote protocol error" }
	return pe
}

func (c *Conn) Send(streamID uint32, payload []byte) (uint64, error) {
	id := c.nextID()
	return id, c.send(Frame{Type: TypeData, StreamID: streamID, RequestID: id, Payload: payload})
}

func (c *Conn) Receive() (Frame, error) { return c.receive() }

func (c *Conn) Ping() error { id := c.nextID(); return c.send(Frame{Type: TypePing, RequestID: id}) }
func (c *Conn) CloseProtocol() error { return c.send(Frame{Type: TypeClose, RequestID: c.nextID()}) }

func Listen(addr string, cfg ServerConfig) (net.Listener, error) {
	cfg.Limits = cfg.Limits.normalize()
	return net.Listen("tcp", addr)
}

func Dial(addr string, cfg ClientConfig) (*Conn, error) {
	timeout := cfg.Timeout
	if timeout == 0 { timeout = 10 * time.Second }
	nc, err := net.DialTimeout("tcp", addr, timeout)
	if err != nil { return nil, err }
	limits := cfg.Limits
	if cfg.MaxPayload != 0 { limits.MaxPayload = cfg.MaxPayload }
	c := newConnWithLimits(nc, limits)
	c.ReadTimeout = cfg.ReadTimeout
	c.WriteTimeout = cfg.WriteTimeout
	if err := c.HandshakeClient(cfg); err != nil { _ = nc.Close(); return nil, err }
	return c, nil
}

func (c *Conn) Serve() error {
	if err := c.HandshakeServer("AOXNet Server"); err != nil { return err }
	for {
		f, err := c.receive()
		if err != nil {
			if errors.Is(err, io.EOF) { return nil }
			return err
		}
		switch f.Type {
		case TypePing:
			if err := c.send(Frame{Type: TypePong, RequestID: f.RequestID}); err != nil { return err }
		case TypeClose:
			return nil
		case TypeData:
			if err := c.send(Frame{Type: TypeData, StreamID: f.StreamID, RequestID: f.RequestID, Payload: f.Payload}); err != nil { return err }
		case TypeError:
			return decodeProtocolError(f.Payload)
		default:
			_ = c.sendError(ErrProtocolState, fmt.Sprintf("unexpected frame type %d", f.Type), f.RequestID)
			return fmt.Errorf("%w: unexpected frame type %d", ErrProtocolState, f.Type)
		}
	}
}
