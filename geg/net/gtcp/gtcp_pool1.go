package main

import (
    "fmt"
    "time"
    "net"
    "gitee.com/johng/gf/g/net/gtcp"
    "gitee.com/johng/gf/g/os/glog"
    "gitee.com/johng/gf/g/os/gtime"
)

func main() {
    // Server
    go gtcp.NewServer("127.0.0.1:8999", func(conn net.Conn) {
        c := gtcp.NewConnByNetConn(conn)
        defer c.Close()
        for {
            data, err := c.Receive(-1)
            if len(data) > 0 {
                if err := c.Send(append([]byte("> "), data...)); err != nil {
                    fmt.Println(err)
                }
            }
            if err != nil {
                break
            }
        }
    }).Run()

    time.Sleep(time.Second)

    // Client
    for {
       if conn, err := gtcp.NewPoolConn("127.0.0.1:8999"); err == nil {
           if b, err := conn.SendReceive([]byte(gtime.Datetime()), -1); err == nil {
               fmt.Println(string(b), conn.LocalAddr(), conn.RemoteAddr())
           } else {
               fmt.Println(err)
           }
           conn.Close()
       } else {
           glog.Error(err)
       }
       time.Sleep(time.Second)
    }
}