package main

import (
    "fmt"
    "time"
    "net"
)

func main() {
    addr := "127.0.0.1:8999"

    // Server
    go func() {
        tcpaddr, err := net.ResolveTCPAddr("tcp4", addr)
        if err != nil {
            panic(err)
        }
        listen, err := net.ListenTCP("tcp", tcpaddr)
        if err != nil {
            panic(err)
        }
        for  {
            if conn, err := listen.Accept(); err != nil {
                panic(err)
            } else if conn != nil {
                go func(conn net.Conn) {
                    buffer := make([]byte, 1024)
                    n, err := conn.Read(buffer)
                    if err != nil {
                        fmt.Println(err)
                    } else {
                        fmt.Println(">", string(buffer[0 : n]))
                    }
                    conn.Close()
                }(conn)
            }
        }
    }()

    time.Sleep(time.Second)

    // Client
    if conn, err := net.Dial("tcp", addr); err == nil {
        for i := 0; i < 2; i++ {
            _, err := conn.Write([]byte("hello"))
            if err != nil {
                fmt.Println(err)
                conn.Close()
                break
            } else {
                fmt.Println("ok")
            }
            // sleep 10 seconds and re-send
            time.Sleep(10*time.Second)
        }
    } else {
        panic(err)
    }

}