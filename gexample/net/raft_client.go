package main

import (
    "net"
    "log"
    "g/net/graft"
    "fmt"
    "g/encoding/gjson"
)

func rpcLogSet() {
    conn, err := net.Dial("tcp", "192.168.2.124:4167")
    if err != nil {
        log.Println(err)
        return
    }

    entry      := graft.LogRequest{}
    entry.Key   = "name3"
    entry.Value = "john3"
    fmt.Println(*gjson.Encode(entry))
    e := graft.SendMsg(conn, 100, *gjson.Encode(entry))
    fmt.Println(e)
    conn.Close()
}

func main() {
    rpcLogSet()
}