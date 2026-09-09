# AOXNet

**AOXNet** is a documented, versioned application-layer protocol for reliable message transport over TCP, with optional authenticated TLS transport.

> Status: production-oriented reference implementation. Protocol security depends on TLS configuration and application authentication; AOXNet's CRC32 is only an accidental-corruption check.

## Highlights

- Deterministic binary framing with bounded payloads
- HELLO/WELCOME capability negotiation
- Multiplexed streams with explicit stream IDs
- Per-stream and connection-level flow control
- WINDOW_UPDATE and graceful GOAWAY shutdown
- Structured protocol error codes
- TLS 1.3 and mutual-TLS security primitives
- Peer certificate identity extraction
- Keepalive/liveness policy
- Runtime lifecycle state machine and protocol telemetry
- Fuzz coverage for frame decoding
- Benchmarks and conformance vectors
- Standard-library-only Go implementation
- CI checks for formatting, vetting, tests, and race detection

## Architecture

```text
Application
    |
AOXNet Streams
    |-- Stream lifecycle
    |-- Flow control / backpressure
    |-- Request correlation
    |
AOXNet Frames
    |-- HELLO / WELCOME
    |-- DATA
    |-- PING / PONG
    |-- WINDOW_UPDATE
    |-- GOAWAY / CLOSE / ERROR
    |
Transport
    |-- TCP
    `-- TLS 1.3 / mTLS
```

## Wire format

All integer fields use big-endian encoding.

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

Magic is `AOX1`. The default maximum payload is 1 MiB and can be configured.

## Stream model

Stream `0` is reserved for connection-level control. Local streams use odd IDs and remote streams use even IDs according to the endpoint role. Implementations must reject invalid stream parity and exhausted IDs.

Flow control exists at both stream and connection level. A sender cannot transmit DATA beyond the available send windows. Receivers advertise additional capacity using `WINDOW_UPDATE` with a positive increment.

`GOAWAY` prevents creation of new streams while allowing existing work to finish. `CLOSE` terminates the protocol session.

## Security

AOXNet does not implement cryptographic authentication itself. CRC32 detects accidental corruption but does **not** provide confidentiality, authenticity, replay protection, or active-attacker resistance.

For untrusted networks:

1. Use TLS 1.3.
2. Prefer mTLS when both peers must be authenticated.
3. Validate certificate identity according to the application's trust policy.
4. Do not treat CRC32 as a security boundary.
5. Keep payload limits, deadlines, and flow-control limits enabled.

## Reliability

The runtime exposes explicit lifecycle states and telemetry counters so applications can observe protocol health. Keepalive policies provide a basis for detecting dead peers; applications should configure deadlines appropriate to their environment.

The reference implementation intentionally favors bounded resources and fail-fast protocol errors over unbounded buffering.

## Compatibility

AOXNet uses explicit protocol versioning and capability negotiation. Unknown optional capabilities should be ignored; unsupported mandatory capabilities must fail the handshake. New frame types must not silently change the meaning of existing frame types.

## Development

Requirements:

- Go 1.22+

Run the complete local validation suite:

```bash
gofmt -w .
go vet ./...
go test ./... -count=1
go test -race ./... -count=1
go test -fuzz=FuzzDecodeNeverPanics ./aoxnet -fuzztime=30s
go test ./aoxnet -run '^$' -bench . -benchmem
```

## Example

```go
ln, err := aoxnet.Listen("127.0.0.1:9090", aoxnet.ServerConfig{})
if err != nil {
    panic(err)
}

for {
    conn, err := ln.Accept()
    if err != nil {
        continue
    }
    go conn.Serve()
}
```

## Project structure

```text
.
├── aoxnet/                 # protocol, framing, streams, runtime and security
├── cmd/
│   ├── aoxnet-server/      # reference server
│   └── aoxnet-client/      # reference client
├── .github/workflows/      # CI
├── go.mod
└── README.md
```

## License and scope

AOXNet is intended for legitimate software engineering, interoperability experiments, and authorized deployments. It is not designed as a covert channel or as a substitute for security-reviewed standards where those standards are appropriate.
