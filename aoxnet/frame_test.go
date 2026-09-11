package aoxnet

import (
	"bufio"
	"bytes"
	"testing"
)

func TestFrameRoundTrip(t *testing.T) {
	original := Frame{Version: Version, Type: TypeData, Flags: 3, StreamID: 7, RequestID: 42, Payload: []byte("AOXNet test")}
	var b bytes.Buffer
	if err := original.Encode(&b, DefaultMaxPayload); err != nil {
		t.Fatal(err)
	}
	got, err := Decode(bufio.NewReader(&b), DefaultMaxPayload)
	if err != nil {
		t.Fatal(err)
	}
	if got.Type != original.Type || got.Flags != original.Flags || got.StreamID != original.StreamID || got.RequestID != original.RequestID || !bytes.Equal(got.Payload, original.Payload) {
		t.Fatalf("round trip mismatch: %#v", got)
	}
}

func TestRejectsOversizedPayload(t *testing.T) {
	var b bytes.Buffer
	f := Frame{Type: TypeData, Payload: []byte("12345")}
	if err := f.Encode(&b, 4); err == nil {
		t.Fatal("expected oversized payload error")
	}
}
