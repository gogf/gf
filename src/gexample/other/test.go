package main

import (
    "fmt"
    "net"
)

func main() {
    conn, err := net.Dial("udp", "1.0.0.111:1234")
    if err != nil {
        fmt.Println(err)
    }
    fmt.Println(conn.RemoteAddr())
    conn.Close()
}