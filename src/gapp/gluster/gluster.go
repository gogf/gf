// 高性能的分布式集群管理工具
// 1、分布式数据管理(使用KV方式)
// 2、服务注册与发现
// 3、服务健康检查
// @todo 命令行管理功能完善，本地配置与集群配置冲突解决方案
package main

import (
    "gapp/gluster/gluster"
)

func main() {
    server := gluster.NewServer()
    server.Run()

    select { }
}