<div align=center>
    <img src="http://cover.kancloud.cn/johng/gf" width="200"/>
</div>

<div align=center>
    <a href="https://godoc.org/github.com/johng-cn/gf" target="_blank">
        <img src="https://godoc.org/github.com/johng-cn/gf?status.svg"/>
    </a>
</div>

# 介绍
gf(Go Frame)是一款为Web服务及网络服务开发设计的，模块化、低耦合、轻量级、高性能的Go语言开发框架。
实现了通用的HTTP/TCP/UDP Server，并提供了Web服务开发的系列核心组件，
包括：MVC、Cookie、Session、服务注册、路由控制、配置管理、模板引擎、数据校验、数据库操作(ORM)等等，
并且提供了数十个实用开发模块，包括：缓存模块、日志模块、JSON模块、命令行模块、二进制模块、环境变量模块、并发安全容器、Goroutine池等等。

gf是开源的，免费的，基于MIT协议进行分发，开源项目地址(gitee与github仓库保持实时同步)：
**Gitee**( https://gitee.com/johng/gf )，**Github**( https://github.com/johng-cn/gf )
	
使用中有任何问题/建议，欢迎加入技术QQ群交流：116707870。如有优秀的gf框架使用案例，欢迎联系作者将地址展示到项目库中，您的牛逼将被世人所瞻仰。

# 特点
1. 纯Go语言实现，开源、免费、跨平台、高性能；
1. 充分利用Go语言特点，采用模块化、低耦合设计；
1. 为Web服务及网络服务开发设计，轻量级简便易用；
1. 精心设计的核心组件及数十个丰富实用的开发模块库；
1. 功能强大的HttpServer及服务注册功能，提高Web开发效率；
1. 丰富详尽的框架文档及专业的技术支持/讨论群，易于使用及维护；
1. 更多特点请查阅框架手册( http://gf.johng.cn )；

# 文档
* [框架介绍](http://gf.johng.cn/494364)
* [加入团队](http://gf.johng.cn/512841)
* [Web服务开发](http://gf.johng.cn/494647)
    * [开始使用](http://gf.johng.cn/494366)
    * [开发模式](http://gf.johng.cn/494367)
    * [服务注册](http://gf.johng.cn/494368)
        * [基本介绍](http://gf.johng.cn/590880)
        * [控制器注册](http://gf.johng.cn/590881)
        * [执行对象注册](http://gf.johng.cn/590882)
        * [回调方法注册](http://gf.johng.cn/590883)
        * [事件回调注册](http://gf.johng.cn/590884)
        * [服务性能分析](http://gf.johng.cn/592298)
    * [路由控制](http://gf.johng.cn/49437)
    * [Cookie](http://gf.johng.cn/494372)
    * [Session](http://gf.johng.cn/494373)
    * [输入输出](http://gf.johng.cn/494374)
        * [请求输入](http://gf.johng.cn/590900)
        * [请求输出](http://gf.johng.cn/590901)
        * [文件上传](http://gf.johng.cn/590902)
        * [HTTP客户端](http://gf.johng.cn/590903)
    * [配置管理](http://gf.johng.cn/494376)
    * [单例管理](http://gf.johng.cn/494377)
    * [数据校验](http://gf.johng.cn/494378)
    * [模板引擎](http://gf.johng.cn/494379)
        * [使用方法](http://gf.johng.cn/591642)
        * [基本语法](http://gf.johng.cn/591643)
        * [基本函数](http://gf.johng.cn/591644)
    * [数据库ORM](http://gf.johng.cn/494380)
* [网络服务开发](http://gf.johng.cn/494648)
    * [TCPServer](http://gf.johng.cn/494382)
    * [UDPServer](http://gf.johng.cn/494383)
* [功能模块设计](http://gf.johng.cn/494384)
    * [缓存模块](http://gf.johng.cn/494385)
    * [日志模块](http://gf.johng.cn/494386)
    * [时间模块](http://gf.johng.cn/494387)
    * [类型转换](http://gf.johng.cn/588061)
    * [JSON模块](http://gf.johng.cn/494388)
    * [命令行模块](http://gf.johng.cn/494389)
    * [二进制模块](http://gf.johng.cn/500342)
    * [Goroutine池](http://gf.johng.cn/504458)
    * [数据编码解析](http://gf.johng.cn/511393)
    * [环境变量模块](http://gf.johng.cn/494390)
    * [文件管理模块](http://gf.johng.cn/494391)
    * [文件监控模块](http://gf.johng.cn/593848)
    * [并发安全容器](http://gf.johng.cn/494392)
        * [gmap](http://gf.johng.cn/590858)
        * [gtype](http://gf.johng.cn/590859)
        * [glist](http://gf.johng.cn/590860)
        * [gset](http://gf.johng.cn/590861)
        * [garray](http://gf.johng.cn/590862)
        * [gqueue](http://gf.johng.cn/5908603)
        * [gchan](http://gf.johng.cn/590864)
    * [通用编码模块](http://gf.johng.cn/494393)
    * [其他模块介绍](http://gf.johng.cn/494394)

    
# 贡献
gf是开源的、免费的软件，这意味着任何人都可以为其开发和进步贡献力量。
gf的项目源代码目前同时托管在 Gitee 和 Github 平台上，您可以选择您喜欢的平台来 fork 项目和合并你的贡献，
两个平台的仓库将会保持即时的同步。我们非常欢迎有更多的朋友加入到gf框架的开发中来，
您为gf所做出的任何贡献都将会被记录到gf的史册中。

# 安装
```
go get -u gitee.com/johng/gf
```

# 使用
1. **哈喽世界！**
    ```go
    package main
    
    import "gitee.com/johng/gf/g/net/ghttp"
    
    func main() {
        s := ghttp.GetServer()
        s.BindHandler("/", func(r *ghttp.Request){
            r.Response.Write("哈喽世界！")
        })
        s.Run()
    }
    ```
1. **域名 & 多域名支持**

    1. 自定义域名支持：
        ```go
        package main
        
        import "gitee.com/johng/gf/g/net/ghttp"
    
        func main() {
            s := ghttp.GetServer()
            s.Domain("127.0.0.1").BindHandler("/", func(r *ghttp.Request) {
                r.Response.Write("127.0.0.1")   
            })
            s.Domain("localhost").BindHandler("/", func(r *ghttp.Request) {
                r.Response.Write("localhost")
            })
            s.Run()
        }
        ```
    1. 多域名绑定支持：
        ```go
        package main
        
        import "gitee.com/johng/gf/g/net/ghttp"
    
        func main() {
            s := ghttp.GetServer()
            s.Domain("localhost").BindHandler("/", func(r *ghttp.Request) {
                r.Response.Write("localhost")
            })
            s.Domain("localhost1，localhost2，localhost3").BindHandler("/", func(r *ghttp.Request) {
                r.Response.Write("localhostx")
            })
            s.Run()
        }
        ```

1. **服务注册**
    ```go
    package main
    
    import (
        "gitee.com/johng/gf/g/net/ghttp"
        "gitee.com/johng/gf/g/frame/gmvc"
    )
    
    // 定义业务相关的控制器对象
    type ControllerUser struct {
        gmvc.Controller
    }
    // 定义操作逻辑 - 展示姓名
    func (c *ControllerUser) Name() {
        c.Response.Write("John")
    }
    
    // 定义操作逻辑 - 展示年龄
    func (c *ControllerUser) Age() {
        c.Response.Write("18")
    }
     
    func main() {
        // 绑定控制器到指定URI，所有控制器的公开方法将会映射到指定URI末尾
        // 例如该方法执行后，查看效果可访问：
        // http://127.0.0.1:8199/user/name
        // http://127.0.0.1:8199/user/age
        s := ghttp.GetServer()
        s.BindController("/user", &ControllerUser{})
        s.Run()
    }
    ```

1. **路由控制**

    ```go
    package main
    
    import "gitee.com/johng/gf/g/net/ghttp"
    
    func main() {
        s := ghttp.GetServer()
        s.BindHandler("/:name", func(r *ghttp.Request){
            r.Response.Write(r.Get("name"))
        })
        s.BindHandler("/:name/:action", func(r *ghttp.Request){
            r.Response.Write(r.Get("name"))
            r.Response.Write(r.Get("action"))
        })
        s.BindHandler("/:name/*any", func(r *ghttp.Request){
            r.Response.Write(r.Get("name"))
            r.Response.Write(r.Get("any"))
        })
        s.SetPort(8199)
        s.Run()
    }
    ```

1. **模板引擎**
   
    ```go
    package main
    
    import (
        "gitee.com/johng/gf/g/net/ghttp"
        "gitee.com/johng/gf/g/frame/gmvc"
    )
    
    type ControllerTemplate struct {
        gmvc.Controller
    }
    
    func (c *ControllerTemplate) Info() {
        c.View.Assign("name", "john")
        c.View.Assigns(map[string]interface{}{
            "age"   : 18,
            "score" : 100,
        })
        c.View.Display("user/index.tpl")
    }
    
    func main() {
        s := ghttp.GetServer()
        s.BindControllerMethod("/template", &ControllerTemplate{}, "Info")
        s.Run()
    }
    ```

1. **数据校验**
    ```go
    rule := "url|min-length:11"
    msgs := map[string]string{
        "url"       : "请输入正确的URL地址",
        "minlength" : "地址长度至少为:min位"
    }
    fmt.Println(gvalid.Check("http://gfjohngcn", rule, msgs))
    ```
    
1. **数据库ORM**

    1. **方法操作**
        1. **获取ORM单例对象**
            ```go
            // 获取默认配置的数据库对象(配置名称为"default")
            db, err := gdb.Instance()
            // 获取配置分组名称为"user-center"的数据库对象
            db, err := gdb.Instance("user-center")
            ```
        
        2. **数据写入**
            ```go
            r, err := db.Insert("user", gdb.Map {
                "name": "john",
            })
            ```
        
        3. **数据查询(列表)**
            ```go
            list, err := db.GetAll("select * from user limit 2")
            ```
        
        4. **数据查询(单条)**
            ```go
            one, err := db.GetOne("select * from user limit 2")
            // 或者
            one, err := db.GetOne("select * from user where uid=1000")
            ```
        
        5. **数据保存**
            ```go
            r, err := db.Save("user", gdb.Map {
                "uid"  :  1,
                "name" : "john",
            })
            ```
        
        7. **数据更新/删除**
            ```go
            // db.Update/db.Delete 同理
            r, err := db.Update("user", gdb.Map {"name": "john"}, "uid=?", 10000)
            r, err := db.Update("user", "name='john'", "uid=10000")
            r, err := db.Update("user", "name=?", "uid=?", "john", 10000)
            ```
    
    1. **链式操作**
    
        1. **链式查询**
            ```go
            // 查询多条记录并使用Limit分页
            r, err := db.Table("user u").LeftJoin("user_detail ud", "u.uid=ud.uid").Fields("u.*, ud.site").Where("u.uid > ?", 1).Limit(0, 10).Select()
            // 查询符合条件的单条记录(第一条)
            r, err := db.Table("user u").LeftJoin("user_detail ud", "u.uid=ud.uid").Fields("u.*,ud.site").Where("u.uid=?", 1).One()
            // 查询字段值
            r, err := db.Table("user u").LeftJoin("user_detail ud", "u.uid=ud.uid").Fields("ud.site").Where("u.uid=?", 1).Value()
            // 分组及排序
            r, err := db.Table("user u").LeftJoin("user_detail ud", "u.uid=ud.uid").Fields("u.*,ud.city").GroupBy("city").OrderBy("register_time asc").Select()
            ```
        
        2. **链式更新/删除**
            ```go
            // 更新
            r, err := db.Table("user").Data(gdb.Map{"name" : "john2"}).Where("name=?", "john").Update()
            r, err := db.Table("user").Data("name='john3'").Where("name=?", "john2").Update()
            // 删除
            r, err := db.Table("user").Where("uid=?", 10).Delete()
            ```
        
        3. **链式写入/保存**
            ```go
            r, err := db.Table("user").Data(gdb.Map{"name": "john"}).Insert()
            r, err := db.Table("user").Data(gdb.Map{"uid": 10000, "name": "john"}).Replace()
            r, err := db.Table("user").Data(gdb.Map{"uid": 10001, "name": "john"}).Save()
            ```
        
    1. **事务操作**
    
        1. **开启事务操作**
            ```go
            if tx, err := db.Begin(); err == nil {
                fmt.Println("开启事务操作")
            }
            ```
        2. **事务回滚操作**
            ```go
            if tx, err := db.Begin(); err == nil {
                r, err := tx.Save("user", gdb.Map{
                    "uid"  :  1,
                    "name" : "john",
                })
                tx.Rollback()
                fmt.Println(r, err)
            }
            ```
        3. **事务提交操作**
            ```go
            if tx, err := db.Begin(); err == nil {
                r, err := tx.Save("user", gdb.Map{
                    "uid"  :  1,
                    "name" : "john",
                })
                tx.Commit()
                fmt.Println(r, err)
            }
            ```
        4. **事务链式操作**
            ```go
            if tx, err := db.Begin(); err == nil {
                r, err := tx.Table("user").Data(gdb.Map{"uid":1, "name": "john_1"}).Save()
                tx.Commit()
                fmt.Println(r, err)
            }
            ```

1. **数据编码解析**

    ```go
    // JSON
    data :=
        `{
            "users" : {
                    "count" : 100,
                    "list"  : [
                        {"name" : "Ming", "score" : 60},
                        {"name" : "John", "score" : 99.5}
                    ]
            }
        }`

    if p, e := gparser.LoadContent([]byte(data), "json"); e != nil {
        glog.Error(e)
    } else {
        fmt.Println("John Score:", p.GetFloat32("users.list.1.score"))
    }
    ```

    ```go
    // XML
    data :=
        `<?xml version="1.0" encoding="UTF-8"?>
         <note>
             <to>Tove</to>
             <from>Jani</from>
             <heading>Reminder</heading>
             <body>Don't forget me this weekend!</body>
         </note>`

    if p, e := gparser.LoadContent([]byte(data), "xml"); e != nil {
        glog.Error(e)
    } else {
        fmt.Println("Heading:", p.GetString("note.heading"))
    }
    ```


...


更多特性及示例请查看官方开发文档：[gf.johng.cn](http://gf.johng.cn)

