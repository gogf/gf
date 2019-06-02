# GoFrame
<div align=center>
<img src="https://goframe.org/logo.png" width="100"/>

[![Go Doc](https://godoc.org/github.com/gogf/gf?status.svg)](https://godoc.org/github.com/gogf/gf) 
[![Build Status](https://travis-ci.org/gogf/gf.svg?branch=master)](https://travis-ci.org/gogf/gf) 
[![Go Report](https://goreportcard.com/badge/github.com/gogf/gf)](https://goreportcard.com/report/github.com/gogf/gf) 
[![Code Coverage](https://codecov.io/gh/gogf/gf/branch/master/graph/badge.svg)](https://codecov.io/gh/gogf/gf/branch/master)
[![Production Ready](https://img.shields.io/badge/production-ready-blue.svg)](https://github.com/gogf/gf)
[![License](https://img.shields.io/github/license/gogf/gf.svg?style=flat)](https://github.com/gogf/gf)

</div>

`GF(Go Frame)`是一款模块化、松耦合、生产级Go应用开发框架。提供了常用的核心开发组件，如：缓存、日志、文件、时间、队列、数组、集合、字符串、定时器、命令行、文件锁、内存锁、对象池、连接池、数据校验、数据编码、文件监控、定时任务、数据库ORM、TCP/UDP组件、进程管理/通信、
并发安全容器等等。并提供了Web服务开发的系列核心组件，如：Router、Cookie、Session、服务注册、配置管理、模板引擎等等，支持热重启、热更新、多域名、多端口、多服务、HTTPS、Rewrite等特性。


# 特点
* 模块化、松耦合设计；
* 模块丰富，开箱即用；
* 详尽的开发文档及示例；
* 完善的本地中文化支持；
* 致力于项目的通用方案；
* 更适合企业及团队使用；
* 更多请查阅文档及源码；

# 安装
```html
go get -u github.com/gogf/gf
```
或者
`go.mod`:
```
require github.com/gogf/gf latest
```

# 限制
```shell
golang版本 >= 1.10
```

# 架构
<div align=center>
<img src="https://goframe.org/images/arch.png"/>
</div>



# 文档

开发文档：[https://goframe.org](https://goframe.org)

接口文档：[https://godoc.org/github.com/gogf/gf](https://godoc.org/github.com/gogf/gf)

# 使用
```go
package main

import (
    "github.com/gogf/gf/g"
    "github.com/gogf/gf/g/net/ghttp"
)

func main() {
    s := g.Server()
    s.BindHandler("/", func(r *ghttp.Request) {
        r.Response.Write("Hello World")
    })
    s.Run()
}
```

[更多..](https://goframe.org/start/index)


# 协议

`GF` 使用非常友好的 [MIT](LICENSE) 开源协议进行发布，永久`100%`开源免费。

# 捐赠

如果您喜欢`GF`，要不[给开发者来杯咖啡吧](DONATOR.MD)！
请在捐赠时备注您的`github`/`gitee`账号名称。

# 感谢
<a href="https://www.jetbrains.com/?from=GoFrame"><img src="https://goframe.org/images/jetbrains.png" width="100" alt="JetBrains"/></a>