package main

import (
    "net"
    "gitee.com/johng/gf/g/os/glog"
    "gitee.com/johng/gf/g/net/gtcp"
)

func main() {
    gtcp.NewServer(":8999", func(conn net.Conn) {
        c := gtcp.NewConnByNetConn(conn)
        defer c.Close()
        for {
            if data, err := c.Receive(); err == nil {
                glog.Println(string(data))
            }
        }
    }).Run()
}