# AOXNet Protocol Specification

## Status

AOXNet v1 is an experimental binary application protocol. The wire format is versioned and designed for deterministic parsing, bounded resource usage, and encrypted deployment through TLS.

## Transport

- Default transport: TCP.
- Recommended production transport: TLS over TCP.
- Application protocol data is never treated as cryptographically authenticated by CRC32.
- TLS certificate validation must be enabled for untrusted networks.

## Frame Layout

All integer fields are unsigned and encoded in network byte order (big-endian).

| Offset | Size | Field |
|---:|---:|---|
| 0 | 4 | Magic (`AOX1`) |
| 4 | 1 | Version |
| 5 | 1 | Type |
| 6 | 2 | Flags |
| 8 | 4 | Stream ID |
| 12 | 8 | Request ID |
| 20 | 4 | Payload length |
| 24 | N | Payload |
| 24+N | 4 | CRC32 |

The implementation rejects frames whose payload exceeds the configured maximum before allocating the payload buffer.

## Frame Types

1. HELLO
2. WELCOME
3. DATA
4. PING
5. PONG
6. ERROR
7. CLOSE

Unknown types are protocol errors and must not be silently accepted.

## Session State

A session starts in `HANDSHAKE`.

1. Client sends HELLO.
2. Server validates the protocol version and sends WELCOME.
3. Session enters `OPEN`.
4. DATA, PING/PONG, ERROR, and CLOSE are processed according to their type.
5. CLOSE terminates the session gracefully.

Handshake and idle timeouts are mandatory defensive controls in the reference implementation.

## Resource Controls

Implementations should enforce:

- maximum payload size;
- maximum frame rate;
- bounded handshake duration;
- idle connection timeout;
- write deadlines where appropriate.

These controls reduce accidental resource exhaustion and basic denial-of-service exposure.

## Integrity and Security

CRC32 detects accidental corruption but is not a security primitive. It does not provide authentication, confidentiality, replay protection, or resistance to deliberate tampering. Production deployments should use TLS with certificate verification.

## Compatibility Rules

- Version changes require explicit negotiation or a new protocol version.
- New optional features should be advertised through the HELLO/WELCOME feature list.
- Existing frame fields must not be reinterpreted incompatibly within a version.
- Payload schemas must define their own versioning when application messages evolve.
