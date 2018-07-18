package main

import (
    "fmt"
    "log"
    "net"
    "time"
)

func main() {
    addr := "127.0.0.1:8999"

    tcpaddr, err := net.ResolveTCPAddr("tcp4", addr)
    if err != nil {
        log.Fatal(err)
    }
    listener, err := net.ListenTCP("tcp", tcpaddr)
    if err != nil {
        log.Fatal(err)
    }

    // Server
    done := make(chan error)
    go func(listener net.Listener, done chan<- error) {
        for {
            conn, err := listener.Accept()
            if err != nil {
                done <- err
                return
            }
            go func(conn net.Conn) {
                var buffer [1024]byte
                n, err := conn.Read(buffer[:])
                if err != nil {
                    log.Println(err)
                } else {
                    log.Println(">", string(buffer[0:n]))
                }
                if err := conn.Close(); err != nil {
                    log.Println("error closing server conn:", err)
                }
            }(conn)
        }
    }(listener, done)

    // Client
    conn, err := net.Dial("tcp", addr)
    if err != nil {
        log.Fatal(err)
    }
    for i := 0; i < 2; i++ {
        _, err := conn.Write([]byte("hello"))
        if err != nil {
            log.Println(err)
            err = conn.Close()
            if err != nil {
                log.Println("error closing client conn:", err)
            }
            break
        }
        fmt.Println("ok")
        time.Sleep(2 * time.Second)
    }

    // Shut the server down and wait for it to report back
    err = listener.Close()
    if err != nil {
        log.Fatal("error closing listener:", err)
    }
    err = <-done
    if err != nil {
        log.Println("server returned:", err)
    }
}