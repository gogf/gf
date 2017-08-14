package main

import (
    "g/net/gip"
    "log"
    "gapp/gluster/gluster"
)

func main() {
    //role    := gconsole.Option.GetIndex("role")
    //bindip   := gconsole.Option.Get("bindip")
    bindip   := "192.168.2.102"
    ips, err := gip.IntranetIP()
    if err != nil {
        log.Println(err)
        return
    }
    if len(ips) > 1 && bindip == "" {
        log.Println("You have serveral local ips, please specify one to bind")
        return
    }
    if bindip == "" {
        bindip = ips[0]
    }

    server := gluster.NewServerByIp(bindip)
    server.Run()

    select { }
}