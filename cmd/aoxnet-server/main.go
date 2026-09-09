package main

import (
    "log"
    "os"

    "github.com/sayan9168/netsentinel/aoxnet"
)

func main() {
    addr := "127.0.0.1:9090"
    if len(os.Args) > 1 { addr = os.Args[1] }
    ln, err := aoxnet.Listen(addr, aoxnet.ServerConfig{})
    if err != nil { log.Fatal(err) }
    defer ln.Close()
    log.Printf("AOXNet server listening on %s", addr)
    for {
        nc, err := ln.Accept()
        if err != nil { log.Printf("accept: %v", err); continue }
        go func() {
            c := aoxnet.NewConn(nc, aoxnet.DefaultMaxPayload)
            defer c.Close()
            if err := c.Serve(); err != nil { log.Printf("connection: %v", err) }
        }()
    }
}
