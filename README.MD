# GoFrame
<div align=center>
<img src="https://goframe.org/logo.png" width="100"/>

[![Go Doc](https://godoc.org/github.com/gogf/gf?status.svg)](https://godoc.org/github.com/gogf/gf/g#pkg-subdirectories)
[![Build Status](https://travis-ci.org/gogf/gf.svg?branch=master)](https://travis-ci.org/gogf/gf)
[![Go Report](https://goreportcard.com/badge/github.com/gogf/gf?v=1)](https://goreportcard.com/report/github.com/gogf/gf)
[![Code Coverage](https://codecov.io/gh/gogf/gf/branch/master/graph/badge.svg)](https://codecov.io/gh/gogf/gf/branch/master)
[![Production Ready](https://img.shields.io/badge/production-ready-blue.svg)](https://github.com/gogf/gf)
[![License](https://img.shields.io/github/license/gogf/gf.svg?style=flat)](https://github.com/gogf/gf)

</div>

`GF(GoFrame)` is a modular, full-featured and production-ready application development framework of golang. Providing a series of core components and dozens of practical modules, such as: memcache, configure, validator, logging, array/queue/set/map containers, timer/timing tasks, file/memory lock, object pool, database ORM, etc. Supporting web server integrated with router, cookie, session, logger, template, https, hooks, rewrites and many more features.


# Installation
```
go get -u github.com/gogf/gf
```
or use `go.mod`:
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
<img src="https://goframe.org/images/arch.png"/>
</div>

# Quick Start

## Hello World
```go
package main

import (
    "github.com/gogf/gf/g"
    "github.com/gogf/gf/g/net/ghttp"
)

func main() {
    s := g.Server()
    s.BindHandler("/", func(r *ghttp.Request) {
        r.Response.Write("Hello World")
    })
    s.Run()
}
```

## Rich Router
```go
package main

import (
    "github.com/gogf/gf/g/net/ghttp"
    "github.com/gogf/gf/g"
)

func main() {
    s := g.Server()
    s.BindHandler("/{class}-{course}/:name/*act", func(r *ghttp.Request) {
        r.Response.Writeln(r.Get("class"))
        r.Response.Writeln(r.Get("course"))
        r.Response.Writeln(r.Get("name"))
        r.Response.Writeln(r.Get("act"))
    })
    s.SetPort(8199)
    s.Run()
}
```
## Group Routers
```go
package main

import (
	"github.com/gogf/gf/g"
	"github.com/gogf/gf/g/net/ghttp"
)

type Object struct {}

func (o *Object) Show(r *ghttp.Request) {
	r.Response.Writeln("Object Show")
}

func (o *Object) Delete(r *ghttp.Request) {
	r.Response.Writeln("Object REST Delete")
}

func Handler(r *ghttp.Request) {
	r.Response.Writeln("Handler")
}

func HookHandler(r *ghttp.Request) {
	r.Response.Writeln("Hook Handler")
}

func main() {
	s     := g.Server()
	obj   := new(Object)
	group := s.Group("/api")
	group.ALL ("*",            HookHandler, ghttp.HOOK_BEFORE_SERVE)
	group.ALL ("/handler",     Handler)
	group.ALL ("/obj",         obj)
	group.GET ("/obj/showit",  obj, "Show")
	group.REST("/obj/rest",    obj)
	s.SetPort(8199)
	s.Run()
}
```
or
```go
func main() {
	s   := g.Server()
	obj := new(Object)
	s.Group("/api").Bind([]ghttp.GroupItem{
		{"ALL",  "*",            HookHandler, ghttp.HOOK_BEFORE_SERVE},
		{"ALL",  "/handler",     Handler},
		{"ALL",  "/obj",         obj},
		{"GET",  "/obj/showit",  obj, "Show"},
		{"REST", "/obj/rest",    obj},
	})
	s.SetPort(8199)
	s.Run()
}
```

## Multi ports & domains

```go
package main

import (
    "github.com/gogf/gf/g"
    "github.com/gogf/gf/g/net/ghttp"
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
	"github.com/gogf/gf/g"
	"github.com/gogf/gf/g/net/ghttp"
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
        gfile.PutBinContents("/tmp/" + name, buffer)
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

#### 1). Basic
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
Complex one:
```go
conditions := g.Map{
    "title like ?"         : "%九寨%",
    "online"               : 1,
    "hits between ? and ?" : g.Slice{1, 10},
    "exp > 0"              : nil,
    "category"             : g.Slice{100, 200},
}
result, err := db.Table("article").Where(conditions).All()
// SELECT * FROM article WHERE title like '%九寨%' AND online=1 AND hits between 1 and 10 AND exp > 0 AND category IN(100,200)
```

#### 2). `join`
```go
// SELECT u.*,ud.site FROM user u LEFT JOIN user_detail ud ON u.uid=ud.uid WHERE u.uid=1 LIMIT 1
r, err := db.Table("user u").LeftJoin("user_detail ud", "u.uid=ud.uid").Fields("u.*,ud.site").Where("u.uid=?", 1).One()

// SELECT ud.site FROM user u RIGHT JOIN user_detail ud ON u.uid=ud.uid WHERE u.uid=1 LIMIT 1
r, err := db.Table("user u").RightJoin("user_detail ud", "u.uid=ud.uid").Fields("ud.site").Where("u.uid=?", 1).Value()

// SELECT u.*,ud.site FROM user u INNER JOIN user_detail ud ON u.uid=ud.uid GROUP BY city ORDER BY register_time asc
r, err := db.Table("user u").InnerJoin("user_detail ud", "u.uid=ud.uid").Fields("u.*,ud.city").GroupBy("city").OrderBy("register_time asc").Select()

// SELECT u.*,ud.city FROM user u,user_detail ud WHERE u.uid=ud.uid
r, err := db.Table("user u,user_detail ud").Where("u.uid=ud.uid").Fields("u.*,ud.city").All()
```

#### 3). `select in`
```go
// SELECT * FROM user WHERE uid IN(100,10000,90000)
r, err := db.Table("user").Where("uid IN(?,?,?)", 100, 10000, 90000).All()
// SELECT * FROM user WHERE gender=1 AND uid IN(100,10000,90000)
r, err := db.Table("user").Where("gender=? AND uid IN(?)", 1, g.Slice{100, 10000, 90000}).All()
// SELECT COUNT(*) FROM user WHERE age in(18,50)
r, err := db.Table("user").Where("age IN(?,?)", 18, 50).Count()
```
With `map`:
```go
// SELECT * FROM user WHERE gender=1 AND uid IN(100,10000,90000)
r, err := db.Table("user").Where(g.Map{
    "gender" : 1,
    "uid"    : g.Slice{100,10000,90000},
}).All()
```
With `struct`:
```go
type User struct {
    Id     []int  `gconv:"uid"`
    Gender int    `gconv:"gender"`
}
// SELECT * FROM user WHERE uid IN(100,10000,90000) AND gender=1
r, err := db.Table("user").Where(User{
    "gender" : 1,
    "uid"    : []int{100, 10000, 90000},
}).All()
```

#### 4). `like`
```go
// SELECT * FROM user WHERE name like '%john%'
r, err := db.Table("user").Where("name like ?", "%john%").Select()
// SELECT * FROM user WHERE birthday like '1990-%'
r, err := db.Table("user").Where("birthday like ?", "1990-%").Select()
```

#### 5). `sum`
```go
// SELECT SUM(score) FROM user WHERE uid=1
r, err := db.Table("user").Fields("SUM(score)").Where("uid=?", 1).Value()
```

#### 6). `count`
```go
// SELECT COUNT(1) FROM user WHERE `birthday`='1990-10-01'
r, err := db.Table("user").Where("birthday=?", "1990-10-01").Count()
// SELECT COUNT(uid) FROM user WHERE `birthday`='1990-10-01'
r, err := db.Table("user").Fields("uid").Where("birthday=?", "1990-10-01").Count()
```

#### 7). `distinct`
```go
// SELECT DISTINCT uid,name FROM user 
r, err := db.Table("user").Fields("DISTINCT uid,name").Select()
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




