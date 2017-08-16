package main

import (
    "g/net/gip"
    "gapp/gluster/gluster"
    "g/os/glog"
    "g/os/gconsole"
)

func main() {
    //role    := gconsole.Option.GetIndex("role")
    bindip   := gconsole.Option.Get("bindip")
    ips, err := gip.IntranetIP()
    if err != nil {
        glog.Println(err)
        return
    }
    if len(ips) > 1 && bindip == "" {
        glog.Fatalln("You have serveral local ips, please specify one to bind")
        return
    }
    if bindip == "" {
        bindip = ips[0]
    }

    server := gluster.NewServerByIp(bindip)
    server.Run()

    select { }
}