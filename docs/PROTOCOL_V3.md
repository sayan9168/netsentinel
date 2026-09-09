# AOXNet v3 Production Profile

AOXNet v3 defines a production-oriented transport profile on top of the v2 framing and multiplexing model.

## Transport

TCP is the baseline transport. TLS 1.3 is the recommended secure transport and should be mandatory for untrusted networks. Mutual TLS may be enabled when peers require cryptographic identity.

## Reliability

Connections use bounded frame sizes, frame-rate limits, idle deadlines, explicit stream lifecycle, flow-control windows, and graceful GOAWAY shutdown.

## Keepalive

Implementations SHOULD send periodic PING frames and terminate peers that repeatedly fail to respond within the configured timeout.

## Compatibility

Version negotiation occurs during HELLO/WELCOME. Unknown optional features must be ignored; unsupported mandatory features must fail the handshake with a structured protocol error.

## Security

CRC32 provides corruption detection only. Authentication, confidentiality, and active-attacker resistance come from TLS 1.3. Certificate identity must be validated by the application or mTLS policy.
