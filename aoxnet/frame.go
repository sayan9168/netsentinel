package aoxnet

import (
	"bufio"
	"encoding/binary"
	"errors"
	"fmt"
	"hash/crc32"
	"io"
)

const (
	Magic              uint32 = 0x414F5831 // AOX1
	Version            byte   = 1
	HeaderSize                = 24
	TrailerSize               = 4
	DefaultMaxPayload         = 1 << 20
)

type Type byte

const (
	TypeHello Type = iota + 1
	TypeWelcome
	TypeData
	TypePing
	TypePong
	TypeError
	TypeClose
	TypeWindowUpdate
	TypeGoAway
)

type Frame struct {
	Version   byte
	Type      Type
	Flags     uint16
	StreamID  uint32
	RequestID uint64
	Payload   []byte
}

func (f Frame) Encode(w io.Writer, maxPayload uint32) error {
	if maxPayload == 0 { maxPayload = DefaultMaxPayload }
	if f.Version == 0 { f.Version = Version }
	if len(f.Payload) > int(maxPayload) { return fmt.Errorf("payload too large: %d", len(f.Payload)) }
	h := make([]byte, HeaderSize)
	binary.BigEndian.PutUint32(h[0:4], Magic)
	h[4] = f.Version
	h[5] = byte(f.Type)
	binary.BigEndian.PutUint16(h[6:8], f.Flags)
	binary.BigEndian.PutUint32(h[8:12], f.StreamID)
	binary.BigEndian.PutUint64(h[12:20], f.RequestID)
	binary.BigEndian.PutUint32(h[20:24], uint32(len(f.Payload)))
	crc := crc32.NewIEEE(); _, _ = crc.Write(h); _, _ = crc.Write(f.Payload)
	if _, err := w.Write(h); err != nil { return err }
	if _, err := w.Write(f.Payload); err != nil { return err }
	var trailer [4]byte; binary.BigEndian.PutUint32(trailer[:], crc.Sum32())
	_, err := w.Write(trailer[:]); return err
}

func Decode(r *bufio.Reader, maxPayload uint32) (Frame, error) {
	if maxPayload == 0 { maxPayload = DefaultMaxPayload }
	h := make([]byte, HeaderSize)
	if _, err := io.ReadFull(r, h); err != nil { return Frame{}, err }
	if binary.BigEndian.Uint32(h[0:4]) != Magic { return Frame{}, errors.New("invalid AOXNet magic") }
	if h[4] != Version { return Frame{}, fmt.Errorf("unsupported protocol version: %d", h[4]) }
	n := binary.BigEndian.Uint32(h[20:24])
	if n > maxPayload { return Frame{}, fmt.Errorf("payload exceeds limit: %d", n) }
	payload := make([]byte, n)
	if _, err := io.ReadFull(r, payload); err != nil { return Frame{}, err }
	var trailer [4]byte
	if _, err := io.ReadFull(r, trailer[:]); err != nil { return Frame{}, err }
	crc := crc32.NewIEEE(); _, _ = crc.Write(h); _, _ = crc.Write(payload)
	if binary.BigEndian.Uint32(trailer[:]) != crc.Sum32() { return Frame{}, errors.New("CRC32 mismatch") }
	return Frame{Version:h[4], Type:Type(h[5]), Flags:binary.BigEndian.Uint16(h[6:8]), StreamID:binary.BigEndian.Uint32(h[8:12]), RequestID:binary.BigEndian.Uint64(h[12:20]), Payload:payload}, nil
}
