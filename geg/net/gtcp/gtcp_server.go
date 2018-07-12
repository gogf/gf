package main

import (
    "net"
    "gitee.com/johng/gf/g/net/gtcp"
)

func main() {
    gtcp.NewServer(":8999", func(conn net.Conn) {
        defer conn.Close()
        for {
            if data, err := gtcp.Receive(conn); err == nil {
                gtcp.Send(conn, append([]byte("> "), data...))
            }
        }
    }).Run()
}