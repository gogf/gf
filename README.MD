# GoFrame

[![Go Doc](https://godoc.org/github.com/gogf/gf?status.svg)](https://godoc.org/github.com/gogf/gf)
[![Build Status](https://travis-ci.org/gogf/gf.svg?branch=master)](https://travis-ci.org/gogf/gf)
[![Go Report](https://goreportcard.com/badge/github.com/gogf/gf?v=1)](https://goreportcard.com/report/github.com/gogf/gf)
[![Code Coverage](https://codecov.io/gh/gogf/gf/branch/master/graph/badge.svg)](https://codecov.io/gh/gogf/gf/branch/master)
[![Production Ready](https://img.shields.io/badge/production-ready-blue.svg)](https://github.com/gogf/gf)
[![License](https://img.shields.io/github/license/gogf/gf.svg?style=flat)](https://github.com/gogf/gf)

English | [简体中文](README_ZH.MD)

`GF(GoFrame)` is a modular, full-featured and production-ready application development framework 
of golang. Providing a series of core components and dozens of practical modules, such as: 
cache, logging, containers, timer, resource, validator, database orm, etc. 
Supporting web server integrated with router, cookie, session, middleware, logger, configure, 
template, https, hooks, rewrites and many more features. 


# Installation
```
go get -u -v github.com/gogf/gf
```
suggested using `go.mod`:
```
require github.com/gogf/gf latest
```

# Limitation
```
golang version >= 1.13
```

# Documentation

* 中文官网: https://goframe.org
* GoDoc API: https://godoc.org/github.com/gogf/gf


# Discussion
- QQ Group：[116707870](//shang.qq.com/wpa/qunwpa?idkey=195f91eceeb5d7fa76009b7cd5a4641f70bf4897b7f5a520635eb26ff17adfe7)
- WX Group：Add friend`389961817` in WeChat, commenting `GF`
- Issues：https://github.com/gogf/gf/issues

> It's recommended learning `GoFrame` through its awesome source codes and API reference.

# Architecture
<div align=center>
<img src="https://goframe.org/images/arch.png?v=11"/>
</div>

# Performance

> Run `5` times for each case of each project and pick up the best testing result.

Performance testing cases source codes are hosted at: https://github.com/gogf/gf-performance

## 1. Hello World
<table>
<tr>
<th>Throughputs</th>
<th>Mean Latency</th>
<th>P99 Latency</th>
</tr>
<tr>
<td width="30%"><img src="http://gfcdn.johng.cn/images/performance/throughputs1.jpeg"></td>
<td width="30%"><img src="http://gfcdn.johng.cn/images/performance/meanlatency1.jpeg"></td>
<td width="30%"><img src="http://gfcdn.johng.cn/images/performance/p99latency1.jpeg"></td>
</tr>
</table>

## 2. JSON Response.
<table>
<tr>
<th>Throughputs</th>
<th>Mean Latency</th>
<th>P99 Latency</th>
</tr>
<tr>
<td width="30%"><img src="http://gfcdn.johng.cn/images/performance/throughputs3.jpeg"></td>
<td width="30%"><img src="http://gfcdn.johng.cn/images/performance/meanlatency3.jpeg"></td>
<td width="30%"><img src="http://gfcdn.johng.cn/images/performance/p99latency3.jpeg"></td>
</tr>
</table>

# License

`GF` is licensed under the [MIT License](LICENSE), 100% free and open-source, forever.

# Contributors
This project exists thanks to all the people who contribute. [[Contributors](https://github.com/gogf/gf/graphs/contributors)].
<a href="https://github.com/gogf/gf/graphs/contributors"><img src="https://opencollective.com/goframe/contributors.svg?width=890&button=false" /></a>

# Donators

We currently accept donation by Alipay/WechatPay, please note your github/gitee account in your payment bill. If you like `GF`, why not [buy developer a cup of coffee](DONATOR.MD)?

# Sponsors
We appreciate any kind of sponsorship for `GF` development. If you've got some interesting, please contact WeChat `389961817` / Email `john@goframe.org`.


# Thanks
<a href="https://www.jetbrains.com/?from=GoFrame"><img src="https://goframe.org/images/jetbrains.png" width="100" alt="JetBrains"/></a>








