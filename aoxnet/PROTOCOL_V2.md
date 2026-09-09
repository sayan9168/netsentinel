# AOXNet v2 Protocol

## Scope

AOXNet v2 extends the v1 binary transport with multiplexed streams, explicit stream lifecycle, connection and stream flow control, structured errors, and graceful shutdown.

## Stream model

- Stream ID `0` is reserved for connection-level control frames.
- Local streams use odd IDs.
- Remote streams use even IDs.
- A stream progresses through `IDLE -> OPEN -> HALF_CLOSED_LOCAL/HALF_CLOSED_REMOTE -> CLOSED`.
- DATA after a stream reaches `CLOSED` is a protocol error.

## Flow control

Each connection and stream has an independent send and receive window.
The default initial window is 256 KiB.
WINDOW_UPDATE carries a positive increment and windows MUST NOT exceed 2^31-1.
An endpoint MUST NOT send DATA exceeding either the stream or connection send window.

## Graceful shutdown

GOAWAY is a connection-level signal that prevents creation of new streams while allowing existing streams to finish.
CLOSE terminates the protocol connection.
Implementations should prefer GOAWAY before transport close when shutting down normally.

## Errors

Protocol errors use an ErrorCode and human-readable message. Implementations should close the affected stream for stream-local errors and terminate the connection for connection-fatal errors.

## Security

The framing checksum is an integrity guard against accidental corruption, not authentication. TLS is required when confidentiality, peer authentication, or protection against active network attackers is needed.
