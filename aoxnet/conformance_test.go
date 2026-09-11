package aoxnet

import (
	"bytes"
	"testing"
)

func TestFrameConformance(t *testing.T) {
	cases := []Frame{
		{Type: TypePing, RequestID: 1},
		{Type: TypeData, StreamID: 1, RequestID: 2, Payload: []byte("AOXNet")},
		{Type: TypeGoAway, RequestID: 3},
	}
	for _, in := range cases {
		var b bytes.Buffer
		if err := in.Encode(&b, DefaultMaxPayload); err != nil {
			t.Fatal(err)
		}
		out, err := Decode(&b, DefaultMaxPayload)
		if err != nil {
			t.Fatal(err)
		}
		if out.Type != in.Type || out.StreamID != in.StreamID || out.RequestID != in.RequestID || !bytes.Equal(out.Payload, in.Payload) {
			t.Fatalf("round trip mismatch: %#v != %#v", out, in)
		}
	}
}
