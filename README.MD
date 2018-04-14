<div align=center>
    <img src="http://cover.kancloud.cn/johng/gf" width="200"/>
</div>

<div align=center>
    <a href="https://godoc.org/github.com/johng-cn/gf" target="_blank">
        <img src="https://godoc.org/github.com/johng-cn/gf?status.svg"/>
    </a>
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
    
    gf框架提供了非常强大的服务注册方式，这也是相对于其他框架如：gin/beego/httprouter最突出的特点之一。
    
    gf框架提供了三种服务注册方式：控制器注册、执行对象注册、回调函数注册，具体介绍请查看官方文档([http://gf.johng.cn/494368](http://gf.johng.cn/494368))。服务注册方式比较：
    
    |  注册方式        |  使用难度  |  安全系数  |  执行性能  | 内存消耗  |
    | ---                   |     ---       | --- | --- | ---|
    |  控制器注册     |  低  |   高 |  低  |  高 |
    |  执行对象注册  |  中  |  中  |  中  | 中 |
    |  回调函数注册  |  高  |  低  |  高  | 低 |
    
    1. **控制器注册**
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
            s.SetPort(8199)
            s.Run()
        }
        ```
    1. **RESTful控制器注册**
        
        支持的HTTP Method: ```GET,PUT,POST,DELETE,PATCH,HEAD,CONNECT,OPTIONS,TRACE```
        ```go
        package main
        
        import (
            "gitee.com/johng/gf/g/net/ghttp"
            "gitee.com/johng/gf/g/frame/gmvc"
        )
        
        // 测试控制器
        type ControllerRest struct {
            gmvc.Controller
        }   
        
        // RESTFul - GET
        func (c *ControllerRest) Get() {
            c.Response.Write("RESTful HTTP Method GET")
        }

        // RESTFul - DELETE
        func (c *ControllerRest) Delete() {
            c.Response.Write("RESTful HTTP Method DELETE")
        }
        
        // 该方法无法映射，将会无法访问到
        func (c *ControllerRest) Hello() {
            c.Response.Write("Hello")
        }
     
        func main() {
            // 控制器公开方法中与HTTP Method方法同名的方法将会自动绑定映射
            s := ghttp.GetServer()
            s.BindControllerRest("/user.rest", &ControllerRest{})
            s.SetPort(8199)
            s.Run()
        }
        ```
    1. **其他注册方式**
    
        由于服务注册方式功能相当丰富，其他注册方式请查看官方开发文档 - 服务注册章节（[http://gf.johng.cn/494368](http://gf.johng.cn/494368)）。
    
1. **路由控制**

    gf框架提供了自建的非常强大的路由控制功能，支持流行的命名匹配规则及模糊匹配规则，并提供了优秀的优先级管理机制。
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
        s.BindHandler("/:name/:any", func(r *ghttp.Request){
            r.Response.Write(r.Get("name"))
            r.Response.Write(r.Get("any"))
        })
        s.SetPort(8199)
        s.Run()
    }
    ```
    1. **命名匹配规则**
    
        使用```:name```方式进行匹配(```name```为自定义的匹配名称)，匹配成功后对应匹配参数会被解析为GET参数并传递给注册的服务使用。
        
        匹配示例1：
        ```go
        rule: /user/:user
        
        /user/john/profile        no match
        /user/                    no match
        /user/john                match
        /user/you                 match
        ```
        匹配示例2：
        ```go
        rule: /:name/action
        
        /john/name                no match
        /john/action              match
        /smith/info               no match
        /smith/info/age           no match
        /smith/action             match
        ```
        匹配示例3：
        ```go
        rule: /:name/:action
        
        /john/name                match
        /john/info                match
        /smith/info               match
        /smith/info/age           no match
        /smith/action/del         no match
        ```
    
    1. **模糊匹配规则**

        使用```*any```方式进行匹配(```any```为自定义的匹配名称)，一般常用语末尾匹配，将会匹配URI随后所有的参数，并将匹配参数解析为GET参数并传递给注册的服务使用。
        
        匹配示例1：
        ```go
        rule: /src/*path
        
        /src/                     match
        /src/somefile.go          match
        /src/subdir/somefile.go   match
        /user/                    no match
        /user/john                no match
        ```
        匹配示例2：
        ```go
        rule: /src/*path/:action
        
        /src/                     no match
        /src/somefile.go          no match
        /src/somefile.go/del      match
        /src/subdir/file.go/del   match
        ```
        匹配示例3：
        ```go
        rule: /src/*path/show
        
        /src/                     no match
        /src/somefile.go          no match
        /src/somefile.go/del      no match
        /src/somefile.go/show     match
        /src/subdir/file.go/show  match
        ```
    
    1. **路由优先级控制**
    
        优先级控制最主要的是两点因素，1、层级越深的规则优先级越高；2、命名匹配比模糊匹配优先级高。
        我们来看示例（左边的规则优先级比右边高）：
        ```go
        /user/name               >            /user/:action
        /:name/info              >            /:name/:action
        /:name/:action           >            /:name/*action
        /src/path/del            >            /src/path
        /src/path/del            >            /src/path/:action
        /src/path/*any           >            /src/path
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
        
        6. **批量操作**
            ```go
            // BatchInsert/BatchReplace/BatchSave 同理
            _, err := db.BatchInsert("user", gdb.List {
                {"name": "john_1"},
                {"name": "john_2"},
                {"name": "john_3"},
                {"name": "john_4"},
            }, 10)
            ```
        
        7. **数据更新/删除**
            ```go
            // db.Update/db.Delete 同理
            r, err := db.Update("user", gdb.Map {"name": "john"}, "uid=?", 10000)
            r, err := db.Update("user", "name='john'", "uid=10000")
            r, err := db.Update("user", "name=?", "uid=?", "john", 10000)
            ```
            注意，参数域支持并建议使用预处理模式进行输入，避免SQL注入风险。
    
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
        
        4. **链式批量写入**
            ```go
            r, err := db.Table("user").Data(gdb.List{
                {"name": "john_1"},
                {"name": "john_2"},
                {"name": "john_3"},
                {"name": "john_4"},
            }).Insert()
            ```
            可以指定批量操作中分批写入数据库的每批次写入条数数量：
            ```go
            r, err := db.Table("user").Data(gdb.List{
                {"name": "john_1"},
                {"name": "john_2"},
                {"name": "john_3"},
                {"name": "john_4"},
            }).Batch(2).Insert()
            ```
        
        5. **链式批量保存**
            ```go
            r, err := db.Table("user").Data(gdb.List{
                {"uid":10000, "name": "john_1"},
                {"uid":10001, "name": "john_2"},
                {"uid":10002, "name": "john_3"},
                {"uid":10003, "name": "john_4"},
            }).Save()
            ```
        
    1. **事务操作**
    
        开启事务操作可以通过执行```db.Begin```方法，该方法返回事务的操作对象，类型为```*gdb.Tx```，通过该对象执行后续的数据库操作，并可通过```tx.Commit```提交修改，或者通过```tx.Rollback```回滚修改。
        
        1. **开启事务操作**
            ```go
            if tx, err := db.Begin(); err == nil {
                fmt.Println("开启事务操作")
            }
            ```
            事务操作对象可以执行所有db对象的方法，具体请参考[API文档](https://godoc.org/github.com/johng-cn/gf/g/database/gdb)。
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
            事务操作对象仍然可以通过```tx.Table```或者```tx.From```方法返回一个链式操作的对象，该对象与```db.Table```或者```db.From```方法返回值相同，只不过数据库操作在事务上执行，可提交或回滚。
            ```go
            if tx, err := db.Begin(); err == nil {
                r, err := tx.Table("user").Data(gdb.Map{"uid":1, "name": "john_1"}).Save()
                tx.Commit()
                fmt.Println(r, err)
            }
            ```
            其他链式操作请参考上述链式操作章节。


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
    * [数据库ORM](http://gf.johng.cn/494380)
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