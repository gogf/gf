package main

import (
    "net"
    "gitee.com/johng/gf/g/net/gtcp"
)

func main() {
    gtcp.NewServer(":8999", func(conn net.Conn) {
        for {
            buffer := make([]byte, 1024)
            if length, err := conn.Read(buffer); err == nil {
                conn.Write(append([]byte("> "), buffer[0 : length]...))
            }
        }
    }).Run()
}