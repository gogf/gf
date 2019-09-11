# GoFrame

[![Go Doc](https://godoc.org/github.com/gogf/gf?status.svg)](https://godoc.org/github.com/gogf/gf/g#pkg-subdirectories)
[![Build Status](https://travis-ci.org/gogf/gf.svg?branch=master)](https://travis-ci.org/gogf/gf)
[![Go Report](https://goreportcard.com/badge/github.com/gogf/gf?v=1)](https://goreportcard.com/report/github.com/gogf/gf)
[![Code Coverage](https://codecov.io/gh/gogf/gf/branch/master/graph/badge.svg)](https://codecov.io/gh/gogf/gf/branch/master)
[![Production Ready](https://img.shields.io/badge/production-ready-blue.svg)](https://github.com/gogf/gf)
[![License](https://img.shields.io/github/license/gogf/gf.svg?style=flat)](https://github.com/gogf/gf)

English | [简体中文](README_ZH.MD)

`GF(GoFrame)` is a modular, full-featured and production-ready application development framework of golang. Providing a series of core components and dozens of practical modules, such as: memcache, configure, validator, logging, array/queue/set/map containers, timer/timing tasks, file/memory lock, object pool, database ORM, etc. Supporting web server integrated with router, cookie, session, middleware, logger, template, https, hooks, rewrites and many more features.


# Installation
```
go get -u github.com/gogf/gf
```
suggested using `go.mod`:
```
require github.com/gogf/gf latest
```

# Limitation
```
golang version >= 1.10
```

# Documentation

* [APIDoc](https://godoc.org/github.com/gogf/gf)
* [中文文档](https://goframe.org)

# Architecture
<div align=center>
<img src="https://goframe.org/images/arch.png?v=10"/>
</div>

# Quick Start

## Hello World
```go
package main

import (
    "github.com/gogf/gf/frame/g"
    "github.com/gogf/gf/net/ghttp"
)

func main() {
    s := g.Server()
    s.BindHandler("/", func(r *ghttp.Request) {
        r.Response.Write("Hello World")
    })
    s.Run()
}
```
## Router & Middleware
```go
package main

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
)

func main() {
    s := g.Server()
    s.Group("/api.v2", func(g *ghttp.RouterGroup) {
        g.Middleware(func(r *ghttp.Request) {
            r.Response.Write("start")
            r.Middleware.Next()
            r.Response.Write("end")
        })
        g.Group("/order", func(g *ghttp.RouterGroup) {
            g.GET("/list", func(r *ghttp.Request) {
                r.Response.Write("list")
            })
        })
        g.Group("/user", func(g *ghttp.RouterGroup) {
            g.GET("/info", func(r *ghttp.Request) {
                r.Response.Write("info")
            })
            g.POST("/edit", func(r *ghttp.Request) {
                r.Response.Write("edit")
            })
        })
        g.Group("/hook", func(g *ghttp.RouterGroup) {
            g.Hook("/*", ghttp.HOOK_BEFORE_SERVE, func(r *ghttp.Request) {
                r.Response.Write("hook any")
            })
            g.Hook("/:name", ghttp.HOOK_BEFORE_SERVE, func(r *ghttp.Request) {
                r.Response.Write("hook name")
            })
        })
    })
    s.SetPort(8199)
    s.Run()
}
```

## Multi ports & domains
```go
package main

import (
    "github.com/gogf/gf/frame/g"
    "github.com/gogf/gf/net/ghttp"
)

func Hello1(r *ghttp.Request) {
    r.Response.Write("127.0.0.1: Hello1!")
}

func Hello2(r *ghttp.Request) {
    r.Response.Write("localhost: Hello2!")
}

func main() {
    s := g.Server()
    s.Domain("127.0.0.1").BindHandler("/", Hello1)
    s.Domain("localhost").BindHandler("/", Hello2)
    s.SetPort(8100, 8200, 8300)
    s.Run()
}
```
## Template Engine
```go
package main

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
)

func main() {
    s := g.Server()
    s.BindHandler("/template", func(r *ghttp.Request) {
        r.Response.WriteTpl("index.tpl", g.Map{
            "id":   123,
            "name": "john",
        })
    })
    s.SetPort(8199)
    s.Run()
}
```
## File Uploading
```go
func Upload(r *ghttp.Request) {
    if f, h, e := r.FormFile("upload-file"); e == nil {
        defer f.Close()
        name   := gfile.Basename(h.Filename)
        buffer := make([]byte, h.Size)
        f.Read(buffer)
        gfile.PutBytes("/tmp/" + name, buffer)
        r.Response.Write(name + " uploaded successly")
    } else {
        r.Response.Write(e.Error())
    }
}
```

## ORM Operations

### 1. Retrieving instance
```go
db := g.DB()
db := g.DB("user-center")
```
### 2. Chaining Operations

`Where + string`
```go
// SELECT * FROM user WHERE uid>1 LIMIT 0,10
r, err := db.Table("user").Where("uid > ?", 1).Limit(0, 10).Select()

// SELECT uid,name FROM user WHERE uid>1 LIMIT 0,10
r, err := db.Table("user").Fileds("uid,name").Where("uid > ?", 1).Limit(0, 10).Select()

// SELECT * FROM user WHERE uid=1
r, err := db.Table("user").Where("u.uid=1",).One()
r, err := db.Table("user").Where("u.uid", 1).One()
r, err := db.Table("user").Where("u.uid=?", 1).One()
// SELECT * FROM user WHERE (uid=1) AND (name='john')
r, err := db.Table("user").Where("uid", 1).Where("name", "john").One()
r, err := db.Table("user").Where("uid=?", 1).And("name=?", "john").One()
// SELECT * FROM user WHERE (uid=1) OR (name='john')
r, err := db.Table("user").Where("uid=?", 1).Or("name=?", "john").One()
```
`Where + map`
```go
// SELECT * FROM user WHERE uid=1 AND name='john'
r, err := db.Table("user").Where(g.Map{"uid" : 1, "name" : "john"}).One()
// SELECT * FROM user WHERE uid=1 AND age>18
r, err := db.Table("user").Where(g.Map{"uid" : 1, "age>" : 18}).One()
```
`Where + struct/*struct`
```go
type User struct {
    Id       int    `json:"uid"`
    UserName string `gconv:"name"`
}
// SELECT * FROM user WHERE uid =1 AND name='john'
r, err := db.Table("user").Where(User{ Id : 1, UserName : "john"}).One()
// SELECT * FROM user WHERE uid =1
r, err := db.Table("user").Where(&User{ Id : 1}).One()
```
### 3. Update & Delete
```go
// UPDATE user SET name='john guo' WHERE name='john'
r, err := db.Table("user").Data(gdb.Map{"name" : "john guo"}).Where("name=?", "john").Update()
r, err := db.Table("user").Data("name='john guo'").Where("name=?", "john").Update()
// UPDATE user SET status=1 ORDER BY login_time asc LIMIT 10
r, err := db.Table("user").Data("status", 1).OrderBy("login_time asc").Limit(10).Update

// DELETE FROM user WHERE uid=10
r, err := db.Table("user").Where("uid=?", 10).Delete()
// DELETE FROM user ORDER BY login_time asc LIMIT 10
r, err := db.Table("user").OrderBy("login_time asc").Limit(10).Delete()
```
### 4. Insert & Replace & Save
```go
r, err := db.Table("user").Data(g.Map{"name": "john"}).Insert()
r, err := db.Table("user").Data(g.Map{"uid": 10000, "name": "john"}).Replace()
r, err := db.Table("user").Data(g.Map{"uid": 10001, "name": "john"}).Save()
```
### 5. Transaction
```go
if tx, err := db.Begin(); err == nil {
    r, err := tx.Save("user", g.Map{
        "uid"  :  1,
        "name" : "john",
    })
    tx.Commit()
}
```
### 6. Error Handling
```go
func GetOrderInfo(id int) (order *Order, err error) {
    err = g.DB().Table("order").Where("id", id).Struct(&order)
    if err != nil && err == sql.ErrNoRows {
        err = nil
    }
    return
}
```

[More Features...](https://goframe.org/start/index)


# License

`GF` is licensed under the [MIT License](LICENSE), 100% free and open-source, forever.

# Donators

We currently accept donation by Alipay/WechatPay, please note your github/gitee account in your payment bill. If you like `GF`, why not [buy developer a cup of coffee](DONATOR.MD)?

# Thanks
<a href="https://www.jetbrains.com/?from=GoFrame"><img src="https://goframe.org/images/jetbrains.png" width="100" alt="JetBrains"/></a>


<!--
# Sponsor
We appreciate any kind of sponsorship for `GF` development. If you've got some interested, please contact john@goframe.org.
-->




