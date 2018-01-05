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
* [一、框架介绍](https://www.kancloud.cn/johng/gf/494364)
* [二、Web服务开发](https://www.kancloud.cn/johng/gf/494366)
    * [1、Web Server](https://www.kancloud.cn/johng/gf/494366)
    * [2、MVC模式](https://www.kancloud.cn/johng/gf/494367)
    * [3、服务注册](https://www.kancloud.cn/johng/gf/494368)
        * [1)、控制器注册](https://www.kancloud.cn/johng/gf/494369)
        * [2)、执行对象注册](https://www.kancloud.cn/johng/gf/494370)
        * [3)、回调函数注册](https://www.kancloud.cn/johng/gf/494371)
    * [4、Cookie](https://www.kancloud.cn/johng/gf/494372)
    * [5、Session](https://www.kancloud.cn/johng/gf/494373)
    * [6、输入输出](https://www.kancloud.cn/johng/gf/494374)
    * [7、路由控制](https://www.kancloud.cn/johng/gf/494375)
    * [8、配置管理](https://www.kancloud.cn/johng/gf/494376)
    * [9、单例管理](https://www.kancloud.cn/johng/gf/494377)
    * [10、数据校验](https://www.kancloud.cn/johng/gf/494378)
    * [11、模板引擎](https://www.kancloud.cn/johng/gf/494379)
    * [12、数据库操作](https://www.kancloud.cn/johng/gf/494380)
* [三、网络服务开发](https://www.kancloud.cn/johng/gf/494382)
    * [1、TCPServer](https://www.kancloud.cn/johng/gf/494382)
    * [2、UDPServer](https://www.kancloud.cn/johng/gf/494383)
* 四、功能模块设计
    * 1)、缓存模块
    * 2)、日志模块
    * 3)、时间模块
    * 4)、JSON模块
    * 5)、命令行模块
    * 6)、环境变量模块
    * 7)、文件管理模块
    * 8)、并发安全容器
    * 9)、通用编码模块
    * 10)、其他模块介绍
