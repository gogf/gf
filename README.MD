<div align=center>
<img src="http://cover.kancloud.cn/johng/gf" width="200"/>
</div>

## 介绍
gf(Go Frame)是一款为Web服务及网络服务开发设计的，模块化、低耦合、轻量级的Go语言开发框架。
实现了通用的HTTP/TCP/UDP Server，并提供了Web服务开发的系列核心组件，
包括：MVC、Cookie、Session、模板引擎、路由控制、配置管理、数据校验、数据库操作(ORM)等等，
并且提供了数十个实用开发模块，包括：缓存模块、日志模块、JSON模块、命令行模块、二进制模块、环境变量模块、并发安全容器、Goroutine池等等。

gf是开源的，免费的，基于MIT协议进行分发，开源项目地址(gitee与github仓库保持实时同步)：
1. Gitee

	https://gitee.com/johng/gf
	
2. Github

	https://github.com/johng-cn/gf
	
使用中有任何问题/建议，欢迎加入技术QQ群交流：116707870

## 安装
```
go get -u gitee.com/johng/gf
```

## 使用
```go
package main

import "gitee.com/johng/gf/g/net/ghttp"

func main() {
    s := ghttp.GetServer()
    s.BindHandler("/", func(r *ghttp.Request){
        r.Response.WriteString("Hello World!")
    })
    s.Run()
}
```
## 文档
* [框架介绍](http://gf.johng.cn/494364)
* [加入团队](http://gf.johng.cn/512841)
* [Web服务开发](http://gf.johng.cn/494647)
    * [WebServer](http://gf.johng.cn/494366)
    * [MVC模式](http://gf.johng.cn/494367)
    * [服务注册](http://gf.johng.cn/494368)
    * [Cookie](http://gf.johng.cn/494372)
    * [Session](http://gf.johng.cn/494373)
    * [输入输出](http://gf.johng.cn/494374)
    * [路由控制](http://gf.johng.cn/49437)
    * [配置管理](http://gf.johng.cn/494376)
    * [单例管理](http://gf.johng.cn/494377)
    * [数据校验](http://gf.johng.cn/494378)
    * [模板引擎](http://gf.johng.cn/494379)
    * [数据库操作](http://gf.johng.cn/494380)
* [网络服务开发](http://gf.johng.cn/494648)
    * [TCPServer](http://gf.johng.cn/494382)
    * [UDPServer](http://gf.johng.cn/494383)
* [功能模块设计](http://gf.johng.cn/494384)
    * [缓存模块](http://gf.johng.cn/494385)
    * [日志模块](http://gf.johng.cn/494386)
    * [时间模块](http://gf.johng.cn/494387)
    * [JSON模块](http://gf.johng.cn/494388)
    * [命令行模块](http://gf.johng.cn/494389)
    * [二进制模块](http://gf.johng.cn/500342)
    * [HTTP客户端](http://gf.johng.cn/499674)
    * [Goroutine池](http://gf.johng.cn/504458)
    * [数据编码解析](http://gf.johng.cn/511393)
    * [环境变量模块](http://gf.johng.cn/494390)
    * [文件管理模块](http://gf.johng.cn/494391)
    * [并发安全容器](http://gf.johng.cn/494392)
    * [通用编码模块](http://gf.johng.cn/494393)
    * [其他模块介绍](http://gf.johng.cn/494394)
