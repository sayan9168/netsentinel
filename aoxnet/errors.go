package aoxnet

// ErrorCode identifies a protocol-level failure.
type ErrorCode uint16

const (
	ErrUnknown ErrorCode = iota
	ErrInvalidFrame
	ErrUnsupportedVersion
	ErrProtocolState
	ErrPayloadTooLarge
	ErrTimeout
	ErrUnauthorized
	ErrRateLimited
	ErrInternal
)

// ProtocolError is a structured AOXNet protocol error.
type ProtocolError struct {
	Code    ErrorCode `json:"code"`
	Message string    `json:"message"`
}

func (e ProtocolError) Error() string { return e.Message }
