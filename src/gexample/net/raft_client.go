package main

import (
    "net"
    "os"
    "fmt"
    "log"
    "g/net/graft"
    "time"
)

func rpcLogSet() {
    conn, err := net.Dial("tcp", "192.168.2.102:4167")
    if err != nil {
        log.Println(err)
        return
    }

    log      := graft.LogEntry{}
    log.Id    = time.Now().UnixNano()
    log.Act   = "set"
    log.Key   = "name"
    log.Value = "john"

    conn.Write([]byte(""))
    var msg [20]byte
    n, err := conn.Read(msg[0:])

    fmt.Println(string(msg[0:n]))
    conn.Close()
}

func main() {

}