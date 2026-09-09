package aoxnet

import "testing"

func FuzzDecodeNeverPanics(f *testing.F) {
    f.Add([]byte("AOX1"))
    f.Add(make([]byte, HeaderSize+TrailerSize))
    f.Fuzz(func(t *testing.T, data []byte) {
        _, _ = DecodeBytes(data, DefaultMaxPayload)
    })
}
