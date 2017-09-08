package main

import (
    "fmt"
    "net"
)

type ST struct {
    I int64
}

func main() {

    interfaces, err :=  net.Interfaces()
    if err != nil {
        panic("Poor soul, here is what you got: " + err.Error())
    }
    for _, inter := range interfaces {
        fmt.Println(inter.Name, inter.HardwareAddr)
    }

}