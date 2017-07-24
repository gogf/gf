package main

import (
    "g/net/gscanner"
    "net"
    "fmt"
)

func main() {
    gscanner.TcpScan("192.168.2.100", "192.168.2.130", 80, func(conn net.Conn){
        fmt.Println(conn.RemoteAddr())
    })
}