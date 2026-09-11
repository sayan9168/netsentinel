package aoxnet

import (
	"bufio"
	"bytes"
	"testing"
)

func BenchmarkFrameEncode(b *testing.B) {
	f := Frame{Type: TypeData, StreamID: 1, RequestID: 1, Payload: bytes.Repeat([]byte("x"), 1024)}
	for i := 0; i < b.N; i++ {
		var out bytes.Buffer
		if err := f.Encode(&out, DefaultMaxPayload); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkFrameDecode(b *testing.B) {
	f := Frame{Type: TypeData, StreamID: 1, RequestID: 1, Payload: bytes.Repeat([]byte("x"), 1024)}
	var encoded bytes.Buffer
	if err := f.Encode(&encoded, DefaultMaxPayload); err != nil {
		b.Fatal(err)
	}
	data := encoded.Bytes()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := Decode(bufio.NewReader(bytes.NewReader(data)), DefaultMaxPayload)
		if err != nil {
			b.Fatal(err)
		}
	}
}
