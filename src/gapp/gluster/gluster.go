// 高性能的分布式集群管理工具
// 1、分布式数据管理(使用KV方式)
// 2、服务注册与发现
// 3、服务健康检查
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