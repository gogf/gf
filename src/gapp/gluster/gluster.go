package main

import (
    "g/net/gip"
    "log"
    "g/net/graft"
    "g/os/gconsole"
)

func main() {
    role    := gconsole.Option.GetIndex("role")
    bindip  := gconsole.Option.GetIndex("bindip")
    monitor := gconsole.Option.GetIndex("monitor")
    if bindip != "" {
        server := graft.NewServerByIp(bindip)
        server.SetMonitor(monitor)
        server.Run()
    } else {
        ips, err := gip.IntranetIP()
        if err != nil {
            log.Println(err)
            return
        }
        for _, ip := range ips {
            server := graft.NewServerByIp(ip)
            server.SetMonitor(monitor)
            server.Run()
        }
    }

    select { }
}