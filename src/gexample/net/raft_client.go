package main

import (
    "net"
    "os"
    "fmt"
)



func main() {
    conn, err := net.Dial("tcp", "192.168.2.102:4166")
    defer conn.Close()
    if err != nil {
        os.Exit(1)
    }

    conn.Write([]byte(""))
    var msg [20]byte
    n, err := conn.Read(msg[0:])

    fmt.Println(string(msg[0:n]))
}