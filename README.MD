# GoFrame 
<img align="right" height="150px" src="https://gfer.me/cover.png">

[![Go Doc](https://godoc.org/github.com/johng-cn/gf?status.svg)](https://godoc.org/github.com/johng-cn/gf) 
[![Build Status](https://travis-ci.org/johng-cn/gf.svg?branch=master)](https://travis-ci.org/johng-cn/gf) 
[![Go Report](https://goreportcard.com/badge/github.com/johng-cn/gf)](https://goreportcard.com/report/github.com/johng-cn/gf) 
[![Documents](https://img.shields.io/badge/docs-100%25-green.svg)](https://gfer.me) 
[![License](https://img.shields.io/github/license/johng-cn/gf.svg?style=flat)](https://github.com/johng-cn/gf)
[![Language](https://img.shields.io/badge/language-go-blue.svg)](https://github.com/johng-cn/gf)
[![Release](https://img.shields.io/github/release/johng-cn/gf.svg?style=flat)](https://github.com/johng-cn/gf/releases)

<!--
[![Code Coverage](https://codecov.io/gh/johng-cn/gf/branch/master/graph/badge.svg)](https://codecov.io/gh/johng-cn/gf)
[![Code Helper](https://www.codetriage.com/johng-cn/gf/badges/users.svg)](https://www.codetriage.com/johng-cn/gf) 
--> 

`GF(GoFrame)` is a modular, lightweight, loosely coupled, high performance application development framework written in Go. Supporting graceful server, hot updates, multi-domain, multi-port, multi-service, HTTP/HTTPS, dynamic/hook routing and many more features. Providing a series of core components and dozens of practical modules.

# Installation
```
go get -u gitee.com/johng/gf
```
# Limitation
```
golang version >= 1.9.2
```

# Documentation

* [中文文档](https://gfer.me/)

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

# License

GF is licensed under the [MIT License](LICENSE), 100% free and open-source.
