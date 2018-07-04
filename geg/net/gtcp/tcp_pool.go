package main

import (
    "gitee.com/johng/gf/g/net/gtcp"
    "time"
    "net"
    "gitee.com/johng/gf/g/os/gtime"
    "fmt"
    "gitee.com/johng/gf/g/os/glog"
)

func main() {
    go gtcp.NewServer("127.0.0.1:8999", func(conn net.Conn) {
        for {
            buffer := make([]byte, 1024)
            if length, err := conn.Read(buffer); err == nil {
                conn.Write(append([]byte("> "), buffer[0 : length]...))
            }
            //conn.Close()
        }
    }).Run()

    time.Sleep(time.Second)

    for {
       if conn, err := gtcp.NewConn("127.0.0.1", 8999); err == nil {
           if b, err := conn.SendReceive([]byte(gtime.Datetime())); err == nil {
               fmt.Println(string(b), conn.LocalAddr(), conn.RemoteAddr())
               conn.Close()
           } else {
               glog.Error(err)
           }
       } else {
           glog.Error(err)
       }
       time.Sleep(time.Second)
    }
}