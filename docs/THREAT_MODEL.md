# AOXNet Threat Model

## Security goals

- Confidentiality and integrity through TLS 1.3.
- Optional mutual TLS for authenticated peers.
- Bounded frames and flow-control windows to reduce resource exhaustion.
- Explicit protocol state transitions and structured errors.

## Out of scope

CRC32 is not a cryptographic authenticator. AOXNet must not treat it as protection against an active attacker. Applications requiring peer authorization should use mTLS and validate certificate identity.

## Abuse resistance

Implementations should enforce maximum payloads, frame rates, idle timeouts, stream counts, and connection lifetimes. Backpressure must be propagated instead of buffering unbounded application data.
