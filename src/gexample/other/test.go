package main

import (
    "g/net/gip"
    "log"
)

type ttt struct {
    Name string
    Age  int `json:"age"`
    Info struct{
        grade string
    }
}
func main() {

    ips, err := gip.IntranetIP()
    if err != nil {
        log.Println("error", err)
        return
    }
    for _, ip := range ips {
        println(ip)
    }

}