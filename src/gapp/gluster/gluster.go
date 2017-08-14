package main

import (
    "g/net/gip"
    "log"
    "g/os/gconsole"
    "gapp/gluster/gluster"
)

func main() {
    //role    := gconsole.Option.GetIndex("role")
    bindip  := gconsole.Option.GetIndex("bindip")
    monitor := gconsole.Option.GetIndex("monitor")
    if bindip != "" {
        server := gluster.NewServerByIp(bindip)
        server.SetFileName("gluster.db")
        server.SetMonitor(monitor)
        server.Run()
    } else {
        ips, err := gip.IntranetIP()
        if err != nil {
            log.Println(err)
            return
        }
        for _, ip := range ips {
            server := gluster.NewServerByIp(ip)
            server.SetFileName("gluster.db")
            server.SetMonitor(monitor)
            server.Run()
        }
    }
    select { }
}