# AOXNet

AOXNet is a small, documented, application-layer transport protocol designed for reliable message framing over TCP.

## Goals

- Deterministic binary framing
- Request/response correlation
- Explicit protocol versioning
- Lightweight handshake and ping/pong
- CRC32 integrity check for accidental corruption
- Standard-library-only Go implementation
- Easy extension without changing the transport layer

AOXNet does **not** attempt to replace TCP, TLS, or established cryptography. For untrusted networks, run it over TLS or another authenticated secure transport.

## Wire format

Every frame is encoded in network byte order (big endian):

| Field | Size |
|---|---:|
| Magic | 4 bytes |
| Version | 1 byte |
| Type | 1 byte |
| Flags | 2 bytes |
| Stream ID | 4 bytes |
| Request ID | 8 bytes |
| Payload Length | 4 bytes |
| Payload | variable |
| CRC32 | 4 bytes |

Magic is `AOX1`. The maximum payload is configurable and defaults to 1 MiB.

## Message types

- `HELLO` / `WELCOME`: capability negotiation
- `DATA`: application payload
- `PING` / `PONG`: liveness
- `ERROR`: protocol/application error
- `CLOSE`: graceful shutdown

## Example

```go
ln, _ := aoxnet.Listen("127.0.0.1:9090", aoxnet.ServerConfig{})
for {
    conn, err := ln.Accept()
    if err != nil { continue }
    go conn.Serve()
}
```

This project is intended for learning, experimentation, and authorized software development. It is not a covert channel or a replacement for security-reviewed production protocols.
