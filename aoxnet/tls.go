package aoxnet

import (
	"crypto/tls"
	"fmt"
	"net"
	"time"
)

// TLSConfig wraps crypto/tls configuration for authenticated, encrypted AOXNet transport.
type TLSConfig struct {
	Config  *tls.Config
	Timeout time.Duration
}

// ListenTLS creates an AOXNet listener protected by TLS.
func ListenTLS(addr string, cfg ServerConfig, tlsCfg *tls.Config) (net.Listener, error) {
	if tlsCfg == nil {
		return nil, fmt.Errorf("TLS configuration is required")
	}
	if len(tlsCfg.Certificates) == 0 {
		return nil, fmt.Errorf("TLS server certificate is required")
	}
	return tls.Listen("tcp", addr, tlsCfg)
}

// DialTLS establishes a TLS-protected AOXNet connection and performs the AOXNet handshake.
func DialTLS(addr string, cfg ClientConfig, tlsCfg *tls.Config) (*Conn, error) {
	if tlsCfg == nil {
		return nil, fmt.Errorf("TLS configuration is required")
	}
	timeout := cfg.Timeout
	if timeout == 0 {
		timeout = 10 * time.Second
	}
	nc, err := tls.DialWithDialer(&net.Dialer{Timeout: timeout}, "tcp", addr, tlsCfg)
	if err != nil {
		return nil, err
	}
	c := NewConn(nc, cfg.MaxPayload)
	c.ReadTimeout = cfg.ReadTimeout
	c.WriteTimeout = cfg.WriteTimeout
	if err := c.HandshakeClient(cfg); err != nil {
		_ = nc.Close()
		return nil, err
	}
	return c, nil
}
