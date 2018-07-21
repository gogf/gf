package main

import (
    "gitee.com/johng/gf/g/net/gudp"
    "fmt"
)

func main() {
    gudp.NewServer("127.0.0.1:8999", func(conn *gudp.Conn) {
        defer conn.Close()
        for {
            if data, _ := conn.Recv(-1); len(data) > 0 {
                fmt.Println(string(data))
            }
        }
    }).Run()
}