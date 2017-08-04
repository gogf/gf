package main

import (
    "g/net/gip"
    "log"
    "g/net/graft"
)

func main() {
    ips, err := gip.IntranetIP()
    if err != nil {
        log.Println(err)
        return
    }

    for _, ip := range ips {
        graft.NewServerByIp(ip).Run()
    }

    select { }
}