package main

import (
    "g/net/graft"
    "g/net/gip"
    "log"
)



func main() {
    ips, err := gip.IntranetIP()
    if err != nil {
        log.Println(err)
        return
    }

    for _, ip := range ips {
        //fmt.Println(ip)
        graft.NewServerByIp(ip).Run()
    }
    select {

    }
}