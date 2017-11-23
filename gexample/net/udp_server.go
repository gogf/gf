package main

import (
    "net"
    "fmt"
    "gf/g/net/gudp"
)

func main() {
    gudp.NewServer(":8999", func(conn *net.UDPConn) {
        var buf [1024]byte
        count, raddr, err := conn.ReadFromUDP(buf[0:])
        if err != nil {
            return
        }
        fmt.Println(raddr.String() + ":", string(buf[0:count]))
        _, err = conn.WriteToUDP([]byte("hi"), raddr)
    }).Run()
}