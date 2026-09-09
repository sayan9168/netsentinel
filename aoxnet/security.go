package aoxnet

import (
    "crypto/tls"
    "crypto/x509"
    "errors"
    "fmt"
    "time"
)

type SecurityConfig struct {
    TLSConfig *tls.Config
    RequireTLS bool
    RequireClientCertificate bool
}

func ValidateTLSConfig(cfg SecurityConfig) error {
    if !cfg.RequireTLS { return nil }
    if cfg.TLSConfig == nil { return errors.New("TLS is required but no TLS config was provided") }
    if cfg.RequireClientCertificate && cfg.TLSConfig.ClientAuth == tls.NoClientCert { return errors.New("client certificate authentication is required") }
    return nil
}

func NewServerTLSConfig(cert tls.Certificate, roots *x509.CertPool, requireClientCert bool) *tls.Config {
    cfg := &tls.Config{MinVersion: tls.VersionTLS13, Certificates: []tls.Certificate{cert}, CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256}}
    if requireClientCert { cfg.ClientCAs = roots; cfg.ClientAuth = tls.RequireAndVerifyClientCert }
    return cfg
}

func NewClientTLSConfig(serverName string, roots *x509.CertPool, clientCert *tls.Certificate) *tls.Config {
    cfg := &tls.Config{MinVersion: tls.VersionTLS13, ServerName: serverName, RootCAs: roots, CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256}}
    if clientCert != nil { cfg.Certificates = []tls.Certificate{*clientCert} }
    return cfg
}

type PeerIdentity struct { Subject string; Issuer string; VerifiedChains int }

func PeerIdentityFromState(state tls.ConnectionState) (PeerIdentity, error) {
    if len(state.PeerCertificates) == 0 { return PeerIdentity{}, errors.New("peer presented no certificate") }
    cert := state.PeerCertificates[0]
    return PeerIdentity{Subject: cert.Subject.String(), Issuer: cert.Issuer.String(), VerifiedChains: len(state.VerifiedChains)}, nil
}

func HandshakeDeadline(d time.Duration) error { if d < 0 { return fmt.Errorf("invalid handshake deadline: %s", d) }; return nil }
