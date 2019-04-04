# GoFrame 
<img align="right" height="150px" src="https://goframe.org/cover.png">

[![Go Doc](https://godoc.org/github.com/gogf/gf?status.svg)](https://godoc.org/github.com/gogf/gf) 
[![Build Status](https://travis-ci.org/gogf/gf.svg?branch=master)](https://travis-ci.org/gogf/gf) 
[![Go Report](https://goreportcard.com/badge/github.com/gogf/gf)](https://goreportcard.com/report/github.com/gogf/gf) 
[![Code Coverage](https://codecov.io/gh/gogf/gf/branch/master/graph/badge.svg)](https://codecov.io/gh/gogf/gf/branch/master)
[![License](https://img.shields.io/github/license/gogf/gf.svg?style=flat)](https://github.com/gogf/gf)
[![Language](https://img.shields.io/badge/language-go-blue.svg)](https://github.com/gogf/gf)
[![Release](https://img.shields.io/github/release/gogf/gf.svg?style=flat)](https://github.com/gogf/gf/releases)

<!--
GoFrame is a modular, loose-coupled, production-ready and most-powerful application development framework of golang. Providing a series of core components and dozens of practical modules, such as: cache, logging, containers, timer, validator, database orm, etc. Supporting web server integrated with router, cookie, session, logger, configure, template, https, hooks, rewrites and many more features. 
-->

`GF(GoFrame)` is a modular, loose-coupled, production-ready and most-powerful application development framework of golang. Providing a series of core components and dozens of practical modules, such as: memcache, configure, validator, logging, array/queue/set/map containers, timer/timing tasks, file/memory lock, object pool, database ORM, etc. Supporting web server integrated with router, cookie, session, logger, template, https, hooks, rewrites and many more features. 

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
golang version >= 1.9.2
```

# Documentation

* [GoDoc](https://godoc.org/github.com/gogf/gf)
* [中文文档](https://goframe.org)

# Architecture
<div align=center>
<img src="https://goframe.org/images/arch.png"/>
</div>

# Quick Start

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

[View More..](https://goframe.org/start/index)


# License

`GF` is licensed under the [MIT License](LICENSE), 100% free and open-source, forever.

# Contributors

- [johng](https://gitee.com/johng)
- [zhaopengme](https://github.com/zhaopengme)
- [wenzi1](https://gitee.com/wenzi1)
- [zseeker](https://gitee.com/zseeker)
- [ymrjqyy](https://gitee.com/ymrjqyy)
- [chenyang351](https://github.com/chenyang351)
- [wxkj](https://gitee.com/wxkj)
- [wxkj001](https://github.com/wxkj001)
- [zhangjinfu](https://gitee.com/zhangjinfu)
- [garfieldkwong](https://gitee.com/garfieldkwong)
- [qq1054000800](https://gitee.com/qq1054000800)

# Donators

- [tiangenglan](https://gitee.com/tiangenglan)
- [zhuhuan12](https://gitee.com/zhuhuan12)
- [zfan_codes](https://gitee.com/zfan_codes)
- [hailaz](https://gitee.com/hailaz)
- [mg91](https://gitee.com/mg91)
- [wxkj](https://gitee.com/wxkj)
- [pibigstar](https://github.com/pibigstar)
- [flyke-xu](https://gitee.com/flyke-xu)





