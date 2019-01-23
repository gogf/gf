# GoFrame 
<img align="right" height="150px" src="https://gfer.me/cover.png">

[![Go Doc](https://godoc.org/github.com/gogf/gf?status.svg)](https://godoc.org/github.com/gogf/gf) 
[![Build Status](https://travis-ci.org/gogf/gf.svg?branch=master)](https://travis-ci.org/gogf/gf) 
[![Go Report](https://goreportcard.com/badge/github.com/gogf/gf)](https://goreportcard.com/report/github.com/gogf/gf) 
[![Documents](https://img.shields.io/badge/docs-100%25-green.svg)](https://gfer.me) 
[![License](https://img.shields.io/github/license/gogf/gf.svg?style=flat)](https://github.com/gogf/gf)
[![Language](https://img.shields.io/badge/language-go-blue.svg)](https://github.com/gogf/gf)
[![Release](https://img.shields.io/github/release/gogf/gf.svg?style=flat)](https://github.com/gogf/gf/releases)

<!--
[![Code Coverage](https://codecov.io/gh/gogf/gf/branch/master/graph/badge.svg)](https://codecov.io/gh/gogf/gf)
[![Code Helper](https://www.codetriage.com/gogf/gf/badges/users.svg)](https://www.codetriage.com/gogf/gf) 
--> 

`GF(GoFrame)` is a modular, lightweight, loosely coupled, high performance application development framework written in Go. Supporting graceful server, hot updates, multi-domain, multi-port, multi-service, HTTP/HTTPS, dynamic/hook routing and many more features. Providing a series of core components and dozens of practical modules.

# Installation
```
go get -u gitee.com/johng/gf
```
or use `go.mod`
```
require gitee.com/johng/gf latest
```

# Limitation
```
golang version >= 1.9.2
```

# Documentation

* [中文文档](https://goframe.org)

# Architecture
<div align=center>
<img src="https://gfer.me/images/arch.png"/>
</div>

# Quick Start

```go
package main

import (
    "gitee.com/johng/gf/g"
    "gitee.com/johng/gf/g/net/ghttp"
)

func main() {
    s := g.Server()
    s.BindHandler("/", func(r *ghttp.Request) {
        r.Response.Write("Hello World")
    })
    s.Run()
}
```

[View More..](https://gfer.me/start/index)


# License

`GF` is licensed under the [MIT License](LICENSE), 100% free and open-source, forever.

# Contributors(TOP 10)

<a href="https://gitee.com/johng" target="_blank" title="John"><img src="https://gitee.com/uploads/27/1309327_johng.png?1530630243" width="60" align="left"></a>

<a href="https://gitee.com/wenzi1" target="_blank" title="蚊子"><img src="https://images.gitee.com/uploads/22/1923122_wenzi1.png" width="60" align="left"></a>

<a href="https://gitee.com/zseeker" target="_blank" title="zseeker"><img src="https://gfer.me/images/contributors/zseeker.png" width="60" align="left"></a>

<a href="https://gitee.com/ymrjqyy" target="_blank" title="一墨染尽青衣颜"><img src="https://images.gitee.com/uploads/27/876827_ymrjqyy.png" width="60" align="left"></a>

<a href="https://github.com/chenyang351" target="_blank" title="chenyang351"><img src="https://avatars1.githubusercontent.com/u/30063958?s=60&v=4" width="60" align="left"></a>

<a href="https://gitee.com/wxkj" target="_blank" title="wxkj"><img src="https://gitee.com/uploads/56/91356_wxkj.png" width="60" align="left"></a>

<a href="https://github.com/wxkj001" target="_blank" title="3wxkj001
"><img src="https://avatars0.githubusercontent.com/u/7794279?s=60&v=4" width="60" align="left"></a>

<a href="https://gitee.com/zhangjinfu" target="_blank" title="张金富"><img src="https://images.gitee.com/uploads/63/356163_zhangjinfu.png" width="60" align="left"></a>

<a href="https://gitee.com/garfieldkwong" target="_blank" title="GarfieldKwong"><img src="https://gfer.me/images/contributors/garfieldkwong.png" width="60" align="left"></a>

<a href="https://gitee.com/qq1054000800" target="_blank" title="hello"><img src="https://gitee.com/uploads/9/2209_qq1054000800.jpg" width="60" align="left"></a>

<br /><br /><br />

# Donators

<a href="https://gitee.com/zhuhuan12" target="_blank" title="zhuhuan12"><img src="https://gitee.com/uploads/39/751839_zhuhuan12.png" width="60" align="left"></a>

<a href="https://gitee.com/zfan_codes" target="_blank" title="范钟"><img src="https://images.gitee.com/uploads/32/2044832_zfan_codes.png" width="60" align="left"></a>

<a href="https://gitee.com/hailaz" target="_blank" title="HaiLaz"><img src="https://gitee.com/uploads/87/1273187_hailaz.png" width="60" align="left"></a>

<a href="https://gitee.com/mg91" target="_blank" title="mg91"><img src="https://images.gitee.com/uploads/30/1410930_mg91.png" width="60" align="left"></a>





