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
* [框架介绍](https://www.kancloud.cn/johng/gf/494364)
* [Web服务开发](https://www.kancloud.cn/johng/gf/494647)
    * [WebServer](https://www.kancloud.cn/johng/gf/494366)
    * [MVC模式](https://www.kancloud.cn/johng/gf/494367)
    * [服务注册](https://www.kancloud.cn/johng/gf/494368)
    * [Cookie](https://www.kancloud.cn/johng/gf/494372)
    * [Session](https://www.kancloud.cn/johng/gf/494373)
    * [输入输出](https://www.kancloud.cn/johng/gf/494374)
    * [路由控制](https://www.kancloud.cn/johng/gf/49437)
    * [配置管理](https://www.kancloud.cn/johng/gf/494376)
    * [单例管理](https://www.kancloud.cn/johng/gf/494377)
    * [数据校验](https://www.kancloud.cn/johng/gf/494378)
    * [模板引擎](https://www.kancloud.cn/johng/gf/494379)
    * [数据库操作](https://www.kancloud.cn/johng/gf/494380)
* [网络服务开发](https://www.kancloud.cn/johng/gf/494648)
    * [TCPServer](https://www.kancloud.cn/johng/gf/494382)
    * [UDPServer](https://www.kancloud.cn/johng/gf/494383)
* [功能模块设计](https://www.kancloud.cn/johng/gf/494384)
    * [缓存模块](https://www.kancloud.cn/johng/gf/494384)
    * [日志模块](https://www.kancloud.cn/johng/gf/494384)
    * [时间模块](https://www.kancloud.cn/johng/gf/494384)
    * [JSON模块](https://www.kancloud.cn/johng/gf/494384)
    * [命令行模块](https://www.kancloud.cn/johng/gf/494384)
    * [HTTP客户端](https://www.kancloud.cn/johng/gf/499674)
    * [环境变量模块](https://www.kancloud.cn/johng/gf/494384)
    * [文件管理模块](https://www.kancloud.cn/johng/gf/494384)
    * [并发安全容器](https://www.kancloud.cn/johng/gf/494384)
    * [通用编码模块](https://www.kancloud.cn/johng/gf/494384)
    * [其他模块介绍](https://www.kancloud.cn/johng/gf/494384)
