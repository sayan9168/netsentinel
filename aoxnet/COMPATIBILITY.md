# AOXNet Compatibility Contract

AOXNet implementations MUST reject unsupported protocol versions rather than silently downgrade.

Feature negotiation is additive: peers advertise capabilities during HELLO/WELCOME and implementations must tolerate unknown feature names.

A frame type is valid only when its stream and payload shape match the protocol specification. Invalid control frames are protocol errors, not application data.

Version changes that alter framing, stream semantics, security requirements, or error interpretation require a new protocol version.

The production compatibility suite should include golden vectors for valid frames, malformed lengths, invalid CRCs, unsupported versions, control-frame shape violations, flow-control exhaustion, and GOAWAY behavior.
