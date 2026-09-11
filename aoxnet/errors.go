package aoxnet

import "fmt"

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

// Error implements error so protocol codes can participate in errors.Is/errors.As and %w wrapping.
func (e ErrorCode) Error() string {
	switch e {
	case ErrInvalidFrame:
		return "invalid frame"
	case ErrUnsupportedVersion:
		return "unsupported version"
	case ErrProtocolState:
		return "protocol state error"
	case ErrPayloadTooLarge:
		return "payload too large"
	case ErrTimeout:
		return "timeout"
	case ErrUnauthorized:
		return "unauthorized"
	case ErrRateLimited:
		return "rate limited"
	case ErrInternal:
		return "internal error"
	default:
		return fmt.Sprintf("unknown error code %d", uint16(e))
	}
}

// ProtocolError is a structured AOXNet protocol error.
type ProtocolError struct {
	Code    ErrorCode `json:"code"`
	Message string    `json:"message"`
}

func (e ProtocolError) Error() string { return e.Message }
