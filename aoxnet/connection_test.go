package aoxnet

import (
	"net"
	"testing"
	"time"
)

func TestHandshakeAndEcho(t *testing.T) {
	clientNet, serverNet := net.Pipe()
	server := NewConn(serverNet, DefaultMaxPayload)
	client := NewConn(clientNet, DefaultMaxPayload)
	done := make(chan error, 1)
	go func() { done <- server.Serve() }()

	if err := client.HandshakeClient(ClientConfig{Name: "test-client"}); err != nil {
		t.Fatal(err)
	}
	stream, err := client.OpenStream()
	if err != nil {
		t.Fatal(err)
	}
	if stream.ID != 1 {
		t.Fatalf("expected first local stream ID 1, got %d", stream.ID)
	}
	if _, err := client.Send(stream.ID, []byte("hello")); err != nil {
		t.Fatal(err)
	}
	_ = client.SetReadDeadline(time.Now().Add(time.Second))
	frame, err := client.Receive()
	if err != nil {
		t.Fatal(err)
	}
	if string(frame.Payload) != "hello" || frame.StreamID != stream.ID {
		t.Fatalf("unexpected echo: %#v", frame)
	}
	if err := client.CloseProtocol(); err != nil {
		t.Fatal(err)
	}
	_ = client.Close()
	_ = server.Close()
	if err := <-done; err != nil {
		t.Fatal(err)
	}
}

func TestWindowUpdateAndGoAway(t *testing.T) {
	clientNet, serverNet := net.Pipe()
	client := NewConn(clientNet, DefaultMaxPayload)
	server := NewConn(serverNet, DefaultMaxPayload)
	done := make(chan error, 1)
	go func() { done <- server.Serve() }()
	defer client.Close()
	defer server.Close()

	if err := client.HandshakeClient(ClientConfig{Name: "test-client"}); err != nil {
		t.Fatal(err)
	}
	stream, err := client.OpenStream()
	if err != nil {
		t.Fatal(err)
	}
	if err := client.SendWindowUpdate(stream.ID, 1024); err != nil {
		t.Fatal(err)
	}
	if err := client.SendGoAway(stream.ID, ErrUnknown); err != nil {
		t.Fatal(err)
	}
	if _, err := client.OpenStream(); err == nil {
		t.Fatal("expected new stream creation to fail after GOAWAY")
	}
	_ = client.Close()
	_ = server.Close()
	<-done
}
