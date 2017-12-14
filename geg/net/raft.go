package main

import (
    "gitee.com/johng/gf/g/net/graft"
    "gitee.com/johng/gf/g/net/gip"
    "log"
)



func main() {
    ips, err := gipv4.IntranetIP()
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