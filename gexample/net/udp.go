package main

import (
    "net"
    "fmt"
    "os"
)

func main() {
    conn, err := net.Dial("udp", "127.0.0.1:8999")
    defer conn.Close()
    if err != nil {
        os.Exit(1)
    }

    conn.Write([]byte(""))
    var msg [20]byte
    n, err := conn.Read(msg[0:])

    fmt.Println(string(msg[0:n]))
}