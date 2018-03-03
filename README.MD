<div align=center>
<img src="http://cover.kancloud.cn/johng/gf" width="200"/>
</div>

<div align=center>
    <a href="https://godoc.org/github.com/johng-cn/gf" target="_blank">
        <img src="https://godoc.org/github.com/johng-cn/gf?status.svg"/>
    </a>
    <!--
    <a href="https://www.codetriage.com/johng-cn/gf" target="_blank">
        <img src="https://www.codetriage.com/johng-cn/gf/badges/users.svg"/>
    </a>
    -->
</div>

## 介绍
gf(Go Frame)是一款为Web服务及网络服务开发设计的，模块化、低耦合、轻量级、高性能的Go语言开发框架。
实现了通用的HTTP/TCP/UDP Server，并提供了Web服务开发的系列核心组件，
包括：MVC、Cookie、Session、模板引擎、路由控制、配置管理、数据校验、数据库操作(ORM)等等，
并且提供了数十个实用开发模块，包括：缓存模块、日志模块、JSON模块、命令行模块、二进制模块、环境变量模块、并发安全容器、Goroutine池等等。

gf是开源的，免费的，基于MIT协议进行分发，开源项目地址(gitee与github仓库保持实时同步)：
**Gitee**( https://gitee.com/johng/gf )，**Github**( https://github.com/johng-cn/gf )
	
使用中有任何问题/建议，欢迎加入技术QQ群交流：116707870。如有优秀的gf框架使用案例，欢迎联系作者将地址展示到项目库中，您的牛逼将被世人所瞻仰。

## 特点
1. 纯Go语言实现，开源、免费、跨平台、高性能；
1. 充分利用Go语言特点，采用模块化、低耦合设计；
1. 为Web服务及网络服务开发设计，轻量级简便易用；
1. 精心设计的核心组件及数十个丰富实用的开发模块库；
1. 功能强大的HttpServer及服务注册功能，提高Web开发效率；
1. 丰富详尽的框架文档及专业的技术支持/讨论群，易于使用及维护；
1. 更多特点请查阅框架手册( http://gf.johng.cn )；


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

    
## 贡献
gf是开源的、免费的软件，这意味着任何人都可以为其开发和进步贡献力量。gf的项目源代码目前同时托管在 Gitee 和 Github 平台上，您可以选择您喜欢的平台来 fork 项目和合并你的贡献，两个平台的仓库将会保持即时的同步。我们非常欢迎有更多的朋友加入到gf框架的开发中来，您为gf所做出的任何贡献都将会被记录到gf的史册中。