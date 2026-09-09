# AOXNet v1

AOXNet (Application-Oriented Exchange Network) is a compact application-layer protocol over TCP.

## Design

1. TCP provides ordered, reliable byte transport.
2. AOXNet adds message boundaries and request/stream identifiers.
3. HELLO/WELCOME negotiates the protocol version.
4. CRC32 detects accidental frame corruption.
5. Applications use `StreamID` to multiplex logical conversations.
6. `RequestID` correlates requests and responses.

## Security model

AOXNet's CRC32 is an integrity check, not authentication and not cryptographic protection. For real networks, use TLS with certificate validation. Never treat AOXNet alone as secure against an active network attacker.

## Extension rules

- Never reuse an existing frame type for a different semantic meaning.
- Add new types with explicit versioning when wire compatibility changes.
- Unknown types should produce a protocol error rather than being silently executed.
- Enforce payload limits before allocation.
