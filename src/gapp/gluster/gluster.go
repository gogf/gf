// 高性能的分布式集群管理工具
// 1、分布式KV数据管理
// 2、服务注册与发现
// 3、服务健康检查
// 4、服务负载均衡
package main

import (
    "gapp/gluster/gluster"
)

func main() {
    server := gluster.NewServer()
    server.Run()

    select { }
}