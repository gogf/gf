## 安装
```
go get -u gitee.com/johng/gf
```

## 使用
```go
package main

import "gitee.com/johng/gf/g/net/ghttp"

func main() {
    ghttp.GetServer().BindHandler("/", func(r *ghttp.Request){
        r.Response.WriteString("Hello World!")
    })
    ghttp.GetServer().Run()
}
```
## 文档
* [框架介绍](https://www.kancloud.cn/johng/gf/框架介绍.md)
* [Web服务开发](https://www.kancloud.cn/johng/gf/Web服务开发.md)
    * [WebServer](https://www.kancloud.cn/johng/gf/HTTPServer.md)
    * [MVC模式](https://www.kancloud.cn/johng/gf/MVC模式.md)
    * [服务注册](https://www.kancloud.cn/johng/gf/服务注册.md)
        * [控制器注册](https://www.kancloud.cn/johng/gf/控制器注册.md)
        * [执行对象注册](https://www.kancloud.cn/johng/gf/执行对象注册.md)
        * [回调函数注册](https://www.kancloud.cn/johng/gf/回调函数注册.md)
    * [Cookie](https://www.kancloud.cn/johng/gf/Cookie.md)
    * [Session](https://www.kancloud.cn/johng/gf/Session.md)
    * [输入输出](https://www.kancloud.cn/johng/gf/输入输出.md)
    * [路由控制](https://www.kancloud.cn/johng/gf/路由控制.md)
    * [配置管理](https://www.kancloud.cn/johng/gf/配置管理.md)
    * [单例管理](https://www.kancloud.cn/johng/gf/单例管理.md)
    * [数据校验](https://www.kancloud.cn/johng/gf/数据校验.md)
    * [模板引擎](https://www.kancloud.cn/johng/gf/模板引擎.md)
    * [数据库操作](https://www.kancloud.cn/johng/gf/数据库操作.md)
* [网络服务开发](https://www.kancloud.cn/johng/gf/网络服务开发.md)
    * [TCPServer](https://www.kancloud.cn/johng/gf/TCPServer.md)
    * [UDPServer](https://www.kancloud.cn/johng/gf/UDPServer.md)
* [功能模块设计](https://www.kancloud.cn/johng/gf/功能模块设计.md)
    * [缓存模块](https://www.kancloud.cn/johng/gf/缓存模块.md)
    * [日志模块](https://www.kancloud.cn/johng/gf/日志模块.md)
    * [时间模块](https://www.kancloud.cn/johng/gf/时间模块.md)
    * [JSON模块](https://www.kancloud.cn/johng/gf/JSON模块.md)
    * [命令行模块](https://www.kancloud.cn/johng/gf/命令行模块.md)
    * [环境变量模块](https://www.kancloud.cn/johng/gf/环境变量模块.md)
    * [文件管理模块](https://www.kancloud.cn/johng/gf/文件管理模块.md)
    * [并发安全容器](https://www.kancloud.cn/johng/gf/并发安全容器.md)
    * [通用编码模块](https://www.kancloud.cn/johng/gf/通用编码模块.md)
    * [其他模块介绍](https://www.kancloud.cn/johng/gf/其他模块介绍.md)
