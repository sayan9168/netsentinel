# AOXNet Production Profile

## Reliability

- Bounded frame and payload sizes are mandatory.
- Stream and connection flow-control windows prevent unbounded sender pressure.
- GOAWAY is the graceful shutdown primitive; new streams are rejected after shutdown begins.
- Keepalive is transport-independent and must be paired with a real PING/PONG transport loop.
- Read and write deadlines must be configured for Internet-facing deployments.

## Security

- CRC32 detects accidental corruption only; it is not authentication.
- TLS 1.3 is the baseline for confidentiality and peer authentication on untrusted networks.
- mTLS is recommended for service-to-service deployments requiring explicit client identity.
- Certificates must be validated using normal hostname/PKI policy; disabling verification is not a production configuration.

## Multiplexing

- Stream 0 is reserved for connection control.
- The current v2 stream manager uses odd local IDs and even remote IDs. A production symmetric implementation must make parity role-aware rather than hard-coding one endpoint's perspective.
- Flow-control increments are positive and bounded by the maximum window.

## Operational guidance

- Expose counters for frames, protocol errors, bytes, active streams, flow-control stalls, keepalive failures, and connection closes.
- Fuzz frame decoding and control-frame parsing before declaring the wire implementation stable.
- Benchmark encode/decode and concurrent stream workloads.
- Treat protocol version and feature negotiation as compatibility boundaries.
