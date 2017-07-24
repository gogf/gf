package main

import (
    "g/net/gscanner"
    "net"
    "fmt"
)

func main() {
    gscanner.UdpScan("192.168.2.100", "192.168.2.103", 8999, func(conn net.Conn){
        //conn.Write([]byte("1"))
        //var msg [20]byte
        //n, err := conn.Read(msg[0:])
        //if err != nil {
        //    fmt.Println(err)
        //}
        //fmt.Println(string(msg[0:n]))
        fmt.Println(conn.RemoteAddr())
    })
    //
    //gscanner.UdpScan("192.168.2.100", "192.168.2.103", 80, func(conn net.Conn){
    //    fmt.Println(conn.RemoteAddr())
    //})
}