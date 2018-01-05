package main

import (
    "fmt"
    "net"
    "gitee.com/johng/gf/g/net/gudp"
)

func main() {
    gudp.NewServer(":8999", func(conn *net.UDPConn) {
        buffer := make([]byte, 1024)
        if length, addr, err := conn.ReadFromUDP(buffer); err == nil {
            fmt.Println(string(buffer[0 : length]), "from", addr.String())
        }
    }).Run()
}