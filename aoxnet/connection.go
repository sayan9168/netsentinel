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

type ServerConfig struct { MaxPayload uint32; ReadTimeout, WriteTimeout time.Duration; Name string }
type ClientConfig struct { MaxPayload uint32; ReadTimeout, WriteTimeout time.Duration; Name string; Timeout time.Duration }

type Conn struct {
    net.Conn
    r *bufio.Reader
    maxPayload uint32
    writeMu sync.Mutex
    nextRequest uint64
    ReadTimeout, WriteTimeout time.Duration
}

func NewConn(c net.Conn, maxPayload uint32) *Conn {
    if maxPayload == 0 { maxPayload = DefaultMaxPayload }
    return &Conn{Conn:c, r:bufio.NewReader(c), maxPayload:maxPayload}
}

func (c *Conn) send(f Frame) error {
    c.writeMu.Lock(); defer c.writeMu.Unlock()
    if c.WriteTimeout > 0 { _ = c.SetWriteDeadline(time.Now().Add(c.WriteTimeout)) }
    return f.Encode(c.Conn, c.maxPayload)
}

func (c *Conn) receive() (Frame, error) {
    if c.ReadTimeout > 0 { _ = c.SetReadDeadline(time.Now().Add(c.ReadTimeout)) }
    return Decode(c.r, c.maxPayload)
}

func (c *Conn) nextID() uint64 { return atomic.AddUint64(&c.nextRequest, 1) }

func (c *Conn) HandshakeClient(cfg ClientConfig) error {
    hello, _ := json.Marshal(Hello{Name:cfg.Name, Version:Version, Features:[]string{"request-id","crc32","ping"}})
    if err := c.send(Frame{Type:TypeHello, RequestID:c.nextID(), Payload:hello}); err != nil { return err }
    f, err := c.receive(); if err != nil { return err }
    if f.Type != TypeWelcome { return errors.New("expected WELCOME") }
    var w Welcome; if err := json.Unmarshal(f.Payload, &w); err != nil { return err }
    if w.Version != Version { return fmt.Errorf("server selected unsupported version %d", w.Version) }
    return nil
}

func (c *Conn) HandshakeServer(name string) error {
    f, err := c.receive(); if err != nil { return err }
    if f.Type != TypeHello { return errors.New("expected HELLO") }
    var h Hello; if err := json.Unmarshal(f.Payload, &h); err != nil { return err }
    if h.Version != Version { return fmt.Errorf("unsupported client version %d", h.Version) }
    payload, _ := json.Marshal(Welcome{Name:name, Version:Version, Features:[]string{"request-id","crc32","ping"}})
    return c.send(Frame{Type:TypeWelcome, RequestID:f.RequestID, Payload:payload})
}

func (c *Conn) Send(streamID uint32, payload []byte) (uint64, error) {
    id := c.nextID()
    return id, c.send(Frame{Type:TypeData, StreamID:streamID, RequestID:id, Payload:payload})
}

func (c *Conn) Receive() (Frame, error) { return c.receive() }

func (c *Conn) Ping() error { id := c.nextID(); return c.send(Frame{Type:TypePing, RequestID:id}) }
func (c *Conn) CloseProtocol() error { return c.send(Frame{Type:TypeClose, RequestID:c.nextID()}) }

func Listen(addr string, cfg ServerConfig) (net.Listener, error) { return net.Listen("tcp", addr) }

func Dial(addr string, cfg ClientConfig) (*Conn, error) {
    timeout := cfg.Timeout; if timeout == 0 { timeout = 10*time.Second }
    nc, err := net.DialTimeout("tcp", addr, timeout); if err != nil { return nil, err }
    c := NewConn(nc, cfg.MaxPayload); c.ReadTimeout=cfg.ReadTimeout; c.WriteTimeout=cfg.WriteTimeout
    if err := c.HandshakeClient(cfg); err != nil { _ = nc.Close(); return nil, err }
    return c, nil
}

func (c *Conn) Serve() error {
    if err := c.HandshakeServer("AOXNet Server"); err != nil { return err }
    for {
        f, err := c.receive(); if err != nil { if errors.Is(err, io.EOF) { return nil }; return err }
        switch f.Type {
        case TypePing:
            if err := c.send(Frame{Type:TypePong, RequestID:f.RequestID}); err != nil { return err }
        case TypeClose:
            return nil
        case TypeData:
            // Echo service for the reference implementation. Applications can build
            // request routing on top of StreamID and RequestID.
            if err := c.send(Frame{Type:TypeData, StreamID:f.StreamID, RequestID:f.RequestID, Payload:f.Payload}); err != nil { return err }
        default:
            return fmt.Errorf("unexpected frame type %d", f.Type)
        }
    }
}
