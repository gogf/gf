<div align=center>
    <img src="http://cover.kancloud.cn/johng/gf" width="200"/>
</div>
<div align=center>
    <a href="https://godoc.org/github.com/johng-cn/gf" target="_blank">
        <img src="https://godoc.org/github.com/johng-cn/gf?status.svg"/>
    </a>
</div>

# 安装
```shell
go get -u gitee.com/johng/gf
```

# 介绍
gf是一款模块化、松耦合、轻量级、高性能的Web开发框架。开源项目地址(仓库保持实时同步)：
[Gitee](https://gitee.com/johng/gf)，[Github](https://github.com/johng-cn/gf)

使用中有任何问题/建议，欢迎加入技术QQ群交流：116707870。如有优秀的框架使用案例，欢迎联系作者将地址展示到项目库中，您的牛逼将被世人所瞻仰。

# 特点
1. 轻量级、高性能，模块化、松耦合设计，丰富的开发模块；
1. 热重启、热更新特性，并支持Web界面及命令行管理接口；
1. 专业的技术交流群，完善的开发文档，良好的中文化支持；
1. 支持多种形式的服务注册特性，灵活高效的路由控制管理；
1. 支持服务事件回调注册功能，可供选择的pprof性能分析模块；
1. 支持配置文件及模板文件的自动检测更新机制，即修改即生效；
1. 强大的数据/表单校验模块，支持常用的40种及自定义校验规则；
1. 提供了对基本数据类型的并发安全封装，提供了常用的并发安全容器；
1. 支持Go变量/Json/Xml/Yml/Toml任意数据格式之间的相互转换及创建；
1. 强大的数据库管理模块，支持应用层级的集群管理、读写分离、负载均衡，方法及链式ORM操作；
1. 更多特点请查阅框架手册( http://gf.johng.cn )；

# 文档
* [框架介绍](http://gf.johng.cn/494364)
* [加入团队](http://gf.johng.cn/512841)
* [源码目录](http://gf.johng.cn/597940)
* [性能测试](http://gf.johng.cn/637443)
* [Web服务开发](http://gf.johng.cn/494647)
    * [开始使用](http://gf.johng.cn/494366)
    * [开发模式](http://gf.johng.cn/494367)
    * [服务注册](http://gf.johng.cn/494368)
        * [基本介绍](http://gf.johng.cn/590880)
        * [控制器注册](http://gf.johng.cn/590881)
        * [执行对象注册](http://gf.johng.cn/590882)
        * [回调方法注册](http://gf.johng.cn/590883)
        * [事件回调注册](http://gf.johng.cn/590884)
    * [路由控制](http://gf.johng.cn/49437)
    * [Cookie](http://gf.johng.cn/494372)
    * [Session](http://gf.johng.cn/494373)
    * [输入输出](http://gf.johng.cn/494374)
        * [请求输入](http://gf.johng.cn/590900)
        * [请求输出](http://gf.johng.cn/590901)
        * [文件上传](http://gf.johng.cn/590902)
        * [HTTP客户端](http://gf.johng.cn/590903)
    * [对象管理](http://gf.johng.cn/494377)
    * [配置管理](http://gf.johng.cn/494376)
        * [基本用法](http://gf.johng.cn/608519)
        * [其他用法](http://gf.johng.cn/608520)
    * [数据校验](http://gf.johng.cn/494378)
    * [分页管理](http://gf.johng.cn/597381)
        * [基本介绍](http://gf.johng.cn/597431)
        * [动态分页](http://gf.johng.cn/597432)
        * [静态分页](http://gf.johng.cn/597433)
        * [Ajax分页](http://gf.johng.cn/597434)
        * [自定义分页](http://gf.johng.cn/597435)
    * [模板引擎](http://gf.johng.cn/494379)
        * [使用方法](http://gf.johng.cn/591642)
        * [基本语法](http://gf.johng.cn/591643)
        * [基本函数](http://gf.johng.cn/591644)
    * [数据库ORM](http://gf.johng.cn/494380)
        * [ORM基本介绍](http://gf.johng.cn/606753)
        * [ORM使用配置](http://gf.johng.cn/596222)
        * [ORM链式操作](http://gf.johng.cn/596224)
        * [ORM方法操作](http://gf.johng.cn/596223)
        * [ORM事务操作](http://gf.johng.cn/596225)
        * [ORM高级特性](http://gf.johng.cn/606435)
    * [框架开发进阶](http://gf.johng.cn/598801)
        * [HTTPS服务](http://gf.johng.cn/598802)
        * [平滑重启特性](http://gf.johng.cn/625833)
        * [服务性能分析](http://gf.johng.cn/592298)
        * [服务日志管理](http://gf.johng.cn/596202)
        * [WebSocket服务](http://gf.johng.cn/643193)
        * [自定义状态码处理](http://gf.johng.cn/609401)
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
    * [文件锁模块](http://gf.johng.cn/626062)
    * [Goroutine池](http://gf.johng.cn/504458)
    * [Redis客户端](http://gf.johng.cn/596344)
    * [数据编码解析](http://gf.johng.cn/511393)
    * [环境变量模块](http://gf.johng.cn/494390)
    * [文件管理模块](http://gf.johng.cn/494391)
    * [文件监控模块](http://gf.johng.cn/593848)
    * [进程管理通信](http://gf.johng.cn/626063)
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

# 使用
## 哈喽世界！
```go
package main

import (
    "gitee.com/johng/gf/g"
    "gitee.com/johng/gf/g/net/ghttp"
)

func main() {
    s := g.Server()
    s.BindHandler("/", func(r *ghttp.Request) {
        r.Response.Write("哈喽世界！")
    })
    s.Run()
}
```
## 多域名支持
```go
package main

import (
    "gitee.com/johng/gf/g"
    "gitee.com/johng/gf/g/net/ghttp"
)

func main() {
    s := g.Server()
    s.Domain("localhost1，localhost2，localhost3").BindHandler("/", func(r *ghttp.Request) {
        r.Response.Write("localhostx")
    })
    s.Run()
}
```
## 多端口支持
```go
package main

import (
    "gitee.com/johng/gf/g"
    "gitee.com/johng/gf/g/net/ghttp"
)

func main() {
    s := g.Server()
    s.BindHandler("/", func(r *ghttp.Request){
        r.Response.Writeln("go frame!")
    })
    s.SetPort(8080, 8081, 8082)
    s.Run()
}
```

## 服务注册
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

func main() {
    // 绑定控制器到指定URI，所有控制器的公开方法将会映射到指定URI末尾
    // 例如该方法执行后，查看效果可访问：
    // http://127.0.0.1:8199/user/name
    s := ghttp.GetServer()
    s.BindController("/user", &ControllerUser{})
    s.Run()
}
```

## 路由控制
```go
package main

import "gitee.com/johng/gf/g/net/ghttp"

func main() {
    s := ghttp.GetServer()
    s.BindHandler("/:name/:action", func(r *ghttp.Request){
        r.Response.Write(r.Get("name"))
        r.Response.Write(r.Get("action"))
    })
    s.SetPort(8199)
    s.Run()
}
```

## 模板引擎
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
    c.View.DisplayContent(`
        <html>
            <head>
                <title>gf template engine</title>
            </head>
            <body>
                <p>Name: {{.name}}</p>
                <p>Age:  {{.age}}</p>
                <p>Score:{{.score}}</p>
            </body>
        </html>
    `)
}

func main() {
    s := ghttp.GetServer()
    s.BindController("/template", &ControllerTemplate{})
    s.SetPort(8199)
    s.Run()
}
```

## 数据校验
```go
rule := "url|min-length:11"
msgs := map[string]string{
    "url"       : "请输入正确的URL地址",
    "minlength" : "地址长度至少为:min位"
}
fmt.Println(gvalid.Check("http://gfjohngcn", rule, msgs))
```

## 数据库ORM

### 链式操作

#### ORM单例对象
```go
// 获取默认配置的数据库对象(配置名称为"default")
db, err := gdb.New()
// 获取配置分组名称为"user-center"的数据库对象
db, err := gdb.New("user-center")
```
#### 单表/联表查询
```go
// 查询多条记录并使用Limit分页
r, err := db.Table("user").Where("u.uid > ?", 1).Limit(0, 10).Select()
// 查询符合条件的单条记录(第一条)
r, err := db.Table("user u").LeftJoin("user_detail ud", "u.uid=ud.uid").Fields("u.*,ud.site").Where("u.uid=?", 1).One()
// 查询指定字段值
r, err := db.Table("user u").RightJoin("user_detail ud", "u.uid=ud.uid").Fields("ud.site").Where("u.uid=?", 1).Value()
// 分组及排序
r, err := db.Table("user u").InnerJoin("user_detail ud", "u.uid=ud.uid").Fields("u.*,ud.city").GroupBy("city").OrderBy("register_time asc").Select()
// 不使用john的联表查询
r, err := db.Table("user u,user_detail ud").Where("u.uid=ud.uid").Fields("u.*,ud.city").All()
// 不使用Fields方法指定查询字段时，默认查询为"*"
r, err := db.Table("user").Where("u.uid=1",).One()
```

#### like查询
```go
r, err := db.Table("user").Where("name like ?", "%john%").Select()
```
    
#### 链式更新/删除
```go
// 更新
r, err := db.Table("user").Data(gdb.Map{"name" : "john2"}).Where("name=?", "john").Update()
r, err := db.Table("user").Data("name='john3'").Where("name=?", "john2").Update()
// 删除
r, err := db.Table("user").Where("uid=?", 10).Delete()
// Data数值方法的参数形式比较灵活
r, err := db.Table("user").Data(`name="john"`).Update()
r, err := db.Table("user").Data("name", "john").Update()
r, err := db.Table("user").Data(g.Map{"name" : "john"}).Update()
```
#### 链式写入/保存
```go
r, err := db.Table("user").Data(gdb.Map{"name": "john"}).Insert()
r, err := db.Table("user").Data(gdb.Map{"uid": 10000, "name": "john"}).Replace()
r, err := db.Table("user").Data(gdb.Map{"uid": 10001, "name": "john"}).Save()
```

#### 链式批量写入
```go
r, err := db.Table("user").Data(gdb.List{
    {"name": "john_1"},
    {"name": "john_2"},
    {"name": "john_3"},
    {"name": "john_4"},
}).Insert()
```

#### 链式批量保存
```go
r, err := db.Table("user").Data(gdb.List{
    {"uid":10000, "name": "john_1"},
    {"uid":10001, "name": "john_2"},
    {"uid":10002, "name": "john_3"},
    {"uid":10003, "name": "john_4"},
}).Save()
```

### 事务操作

#### 开启事务操作
```go
if tx, err := db.Begin(); err == nil {
    fmt.Println("开启事务操作")
}
```

#### 事务回滚操作
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
        
#### 事务提交操作
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
        
#### 事务链式操作
```go
if tx, err := db.Begin(); err == nil {
    r, err := tx.Table("user").Data(gdb.Map{"uid":1, "name": "john"}).Save()
    tx.Commit()
    fmt.Println(r, err)
}
```

## 数据编码解析
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
