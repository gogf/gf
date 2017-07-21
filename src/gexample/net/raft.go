package main

import (
    "net"
    "fmt"
    "g/net/gtcp"
)

func Handler(conn net.Conn){
    client := conn.RemoteAddr().String()
    fmt.Println("connected from: ", client)
    buf := make([]byte,1024)
    for {
        lenght, err := conn.Read(buf)
        if (err != nil) {
            fmt.Println("closed from:", client)
            conn.Close()
            break
        }
        if lenght > 0 {
            buf[lenght] = 0
        }
        fmt.Println(client, string(buf[0:lenght]))
    }
}


func main() {
    gtcp.NewTCPServer(":22200", Handler).Run()
    gtcp.NewTCPServer(":22201", Handler).Run()
    gtcp.ServerWaitGroup.Wait()
}