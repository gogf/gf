package main

import (
    "gitee.com/johng/gf/g/net/gscanner"
    "net"
    "fmt"
    "time"
)

func main() {
    //gscanner.New().SetTimeout(3*time.Second).ScanIp("192.168.2.1", "192.168.2.255", 80, func(conn net.Conn){
    //    fmt.Println(conn.RemoteAddr())
    //})

    gscanner.New().SetTimeout(3*time.Second).ScanIp("120.76.249.1", "120.76.249.255", 80, func(conn net.Conn){
        fmt.Println(conn.RemoteAddr())
    })

    //gscanner.New().SetTimeout(6*time.Second).ScanPort("120.76.249.2", func(conn net.Conn){
    //    //fmt.Println("yes")
    //    fmt.Println(conn.RemoteAddr())
    //})
}