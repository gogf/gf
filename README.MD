<div align=center>
    <img src="https://gfer.me/cover.png" width="150"/>
</div>


`GF(Go Frame)`是一款模块化、松耦合、轻量级、高性能的Go应用开发框架。支持热重启、热更新、多域名、多端口、多服务、HTTP/HTTPS、动态路由等特性
，并提供了Web服务开发的系列核心组件，如：Router、Cookie、Session、服务注册、配置管理、模板引擎、数据校验、分页管理、数据库ORM等等等等，
并且提供了数十个内置核心开发模块集，如：缓存、日志、时间、命令行、二进制、文件锁、内存锁、对象池、连接池、数据编码、进程管理、进程通信、文件监控、定时任务、TCP/UDP组件、
并发安全容器等等等等等等。

开源项目地址(仓库保持实时同步)：
[Gitee](https://gitee.com/johng/gf)，[Github](https://github.com/johng-cn/gf)。
使用中有任何问题/建议，欢迎加入技术QQ群交流：**116707870**。
如有优秀的框架使用案例，欢迎联系作者将地址展示到项目库中，您的牛逼将被世人所瞻仰。

# 安装
```html
go get -u gitee.com/johng/gf
```

# 限制
```shell
golang版本 >= 1.9.2
```

# 架构
<div align=center>
<img src="https://gfer.me/images/arch.png"/>
</div>


# 特点
1. 轻量级、高性能，模块化、松耦合设计，丰富的开发模块；
1. 热重启、热更新特性，并支持Web界面及命令行管理接口；
1. 专业的技术交流群，完善的开发文档及示例代码，良好的中文化支持；
1. 支持多种形式的服务注册特性，灵活高效的路由控制管理；
1. 支持服务事件回调注册功能，可供选择的pprof性能分析模块；
1. 支持配置文件及模板文件的自动检测更新机制，即修改即生效；
1. 支持自定义日期时间格式的时间模块，类似PHP日期时间格式化；
1. 强大的数据/表单校验模块，支持常用的40种及自定义校验规则；
1. 强大的网络通信TCP/UDP组件，并提供TCP连接池特性，简便高效；
1. 提供了对基本数据类型的并发安全封装，提供了常用的数据结构容器；
1. 支持Go变量/Json/Xml/Yml/Toml任意数据格式之间的相互转换及创建；
1. 强大的数据库ORM，支持应用层级的集群管理、读写分离、负载均衡，查询缓存、方法及链式ORM操作；
1. 更多特点请查阅框架[手册](https://gfer.me)和[源码](https://godoc.org/github.com/johng-cn/gf)；

# 文档
[https://gfer.me](https://gfer.me)


# 使用
```go
package main

import (
    "gitee.com/johng/gf/g"
    "gitee.com/johng/gf/g/net/ghttp"
)

func main() {
    s := g.Server()
    s.BindHandler("/", func(r *ghttp.Request) {
        r.Response.Write("Hello World")
    })
    s.Run()
}
```