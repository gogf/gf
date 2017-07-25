package main

import (
    "g/net/gscanner"
    "net"
    "fmt"
)

func main() {
    //gscanner.TcpScan("192.168.2.1", "192.168.2.255", 80, func(conn net.Conn){
    //    fmt.Println(conn.RemoteAddr())
    //})

    gscanner.TcpScan("120.76.249.1", "120.76.249.255", 80, func(conn net.Conn){
        fmt.Println(conn.RemoteAddr())
    })
}