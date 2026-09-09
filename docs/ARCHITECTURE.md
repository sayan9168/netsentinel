# AOXNet Architecture

AOXNet is layered so protocol semantics remain independent from the underlying transport.

## Layers

1. **Application** — application messages and authorization.
2. **Runtime** — lifecycle, stream management, telemetry and liveness policy.
3. **Multiplexing** — stream IDs, stream state and flow-control windows.
4. **Framing** — bounded binary frames, version checks and CRC32 corruption detection.
5. **Transport** — TCP or authenticated TLS 1.3.

## Resource safety

Implementations should enforce maximum payload sizes, bounded queues, connection/stream windows, deadlines and explicit shutdown. Avoid unbounded goroutine or buffer creation from remote input.

## Failure model

Malformed frames, invalid control messages, flow-control violations and unsupported versions are protocol errors. Implementations should fail closed rather than reinterpret invalid input.

## Security boundary

The protocol layer is not an authentication boundary. TLS provides transport confidentiality and peer authentication when correctly configured. Application-level authorization is still required after authentication.

## Observability

Expose counters and lifecycle state to the host application. At minimum, operators should be able to detect protocol errors, connection churn, active-stream growth, bytes transferred and keepalive failures.
