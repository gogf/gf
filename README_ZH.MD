# GoFrame
<img align="right" height="150px" src="https://gfer.me/cover.png">

[![Go Doc](https://godoc.org/github.com/johng-cn/gf?status.svg)](https://godoc.org/github.com/johng-cn/gf) 
[![Build Status](https://travis-ci.org/johng-cn/gf.svg?branch=master)](https://travis-ci.org/johng-cn/gf) 
[![Go Report](https://goreportcard.com/badge/github.com/johng-cn/gf)](https://goreportcard.com/report/github.com/johng-cn/gf) 
[![Documents](https://img.shields.io/badge/docs-100%25-green.svg)](https://gfer.me) 
[![License](https://img.shields.io/github/license/johng-cn/gf.svg?style=flat)](https://github.com/johng-cn/gf)
[![Language](https://img.shields.io/badge/language-go-blue.svg)](https://github.com/johng-cn/gf)
[![Release](https://img.shields.io/github/release/johng-cn/gf.svg?style=flat)](https://github.com/johng-cn/gf/releases)

<!--
[![Code Coverage](https://codecov.io/gh/johng-cn/gf/branch/master/graph/badge.svg)](https://codecov.io/gh/johng-cn/gf)
[![Code Helper](https://www.codetriage.com/johng-cn/gf/badges/users.svg)](https://www.codetriage.com/johng-cn/gf) 
--> 

`GF(Go Frame)`是一款模块化、松耦合、轻量级、高性能的Go应用开发框架。支持热重启、热更新、多域名、多端口、多服务、HTTP/HTTPS、动态路由等特性
，并提供了Web服务开发的系列核心组件，如：Router、Cookie、Session、服务注册、配置管理、模板引擎、数据校验、分页管理、数据库ORM等等等等，
并且提供了数十个内置核心开发模块集，如：缓存、日志、时间、命令行、二进制、文件锁、内存锁、对象池、连接池、数据编码、进程管理、进程通信、文件监控、定时任务、TCP/UDP组件、
并发安全容器等等等等等等。

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