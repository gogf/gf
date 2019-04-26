# GoFrame 
<img align="right" height="150px" src="https://goframe.org/cover.png">

[![Go Doc](https://godoc.org/github.com/gogf/gf?status.svg)](https://godoc.org/github.com/gogf/gf) 
[![Build Status](https://travis-ci.org/gogf/gf.svg?branch=master)](https://travis-ci.org/gogf/gf) 
[![Go Report](https://goreportcard.com/badge/github.com/gogf/gf)](https://goreportcard.com/report/github.com/gogf/gf) 
[![Code Coverage](https://codecov.io/gh/gogf/gf/branch/master/graph/badge.svg)](https://codecov.io/gh/gogf/gf/branch/master)
[![Production Ready](https://img.shields.io/badge/production-ready-blue.svg)](https://github.com/gogf/gf)
[![License](https://img.shields.io/github/license/gogf/gf.svg?style=flat)](https://github.com/gogf/gf)

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

* [APIDoc](https://godoc.org/github.com/gogf/gf)
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

- [aloncn](https://github.com/aloncn)
- [chenyang351](https://github.com/chenyang351)
- [garfieldkwong](https://gitee.com/garfieldkwong)
- [hailaz](https://gitee.com/hailaz)
- [johng](https://johng.cn)
- [jroam](https://github.com/jroam)
- [pibigstar](https://github.com/pibigstar)
- [qq1054000800](https://gitee.com/qq1054000800)
- [qq976739120](https://github.com/qq976739120)
- [touzijiao](https://github.com/touzijiao)
- [wenzi1](https://gitee.com/wenzi1)
- [wxkj001](https://github.com/wxkj001)
- [ymrjqyy](https://gitee.com/ymrjqyy)
- [youyixiao](https://github.com/youyixiao)
- [zhangjinfu](https://gitee.com/zhangjinfu)
- [zhaopengme](https://github.com/zhaopengme)
- [zseeker](https://gitee.com/zseeker)

# Donators

- [flyke-xu](https://gitee.com/flyke-xu)
- [hailaz](https://gitee.com/hailaz)
- [ireadx](https://github.com/ireadx)
- [mg91](https://gitee.com/mg91)
- [pibigstar](https://github.com/pibigstar)
- [tiangenglan](https://gitee.com/tiangenglan)
- [wxkj](https://gitee.com/wxkj)
- [zhuhuan12](https://gitee.com/zhuhuan12)
- [zfan_codes](https://gitee.com/zfan_codes)






