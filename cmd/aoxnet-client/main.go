package main

import (
	"fmt"
	"log"
	"os"

	"github.com/sayan9168/netsentinel/aoxnet"
)

func main() {
	addr := "127.0.0.1:9090"
	if len(os.Args) > 1 {
		addr = os.Args[1]
	}
	c, err := aoxnet.Dial(addr, aoxnet.ClientConfig{Name: "AOXNet CLI"})
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	id, err := c.Send(1, []byte("hello from AOXNet"))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("sent request=%d\n", id)

	f, err := c.Receive()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("received type=%d request=%d stream=%d payload=%q\n", f.Type, f.RequestID, f.StreamID, f.Payload)

	if err := c.CloseProtocol(); err != nil {
		log.Fatal(err)
	}
}
