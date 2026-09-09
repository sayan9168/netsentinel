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

    if err := client.HandshakeClient(ClientConfig{Name:"test-client"}); err != nil { t.Fatal(err) }
    if _, err := client.Send(9, []byte("hello")); err != nil { t.Fatal(err) }
    _ = client.SetReadDeadline(time.Now().Add(time.Second))
    frame, err := client.Receive()
    if err != nil { t.Fatal(err) }
    if string(frame.Payload) != "hello" || frame.StreamID != 9 { t.Fatalf("unexpected echo: %#v", frame) }
    if err := client.CloseProtocol(); err != nil { t.Fatal(err) }
    _ = client.Close(); _ = server.Close()
    if err := <-done; err != nil { t.Fatal(err) }
}
