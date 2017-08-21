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
    server := gluster.NewServer()
    server.Run()

    select { }
}