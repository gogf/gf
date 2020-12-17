# GoFrame

[![Go Doc](https://godoc.org/github.com/gogf/gf?status.svg)](https://godoc.org/github.com/gogf/gf)
[![Build Status](https://travis-ci.org/gogf/gf.svg?branch=master)](https://travis-ci.org/gogf/gf)
[![Go Report](https://goreportcard.com/badge/github.com/gogf/gf?v=1)](https://goreportcard.com/report/github.com/gogf/gf)
[![Code Coverage](https://codecov.io/gh/gogf/gf/branch/master/graph/badge.svg)](https://codecov.io/gh/gogf/gf/branch/master)
[![Production Ready](https://img.shields.io/badge/production-ready-blue.svg)](https://github.com/gogf/gf)
[![License](https://img.shields.io/github/license/gogf/gf.svg?style=flat)](https://github.com/gogf/gf)

English | [简体中文](README_ZH.MD)

`GF(GoFrame)` is a modular, powerful, high-performance and production-ready application development framework 
of golang. Providing a series of core components and dozens of practical modules, such as: 
cache, logging, containers, timer, resource, validator, database orm, etc. 
Supporting web server integrated with router, cookie, session, middleware, logger, configure, 
template, https, hooks, rewrites and many more features. 

> If you're a newbie to `Go`, you may consider `GoFrame` easy and great as `Laravel` in `PHP`, `SpringBoot` in `Java` or `Django` in `Python`.

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
golang version >= 1.11
```

# Architecture
<div align=center>
<img src="https://goframe.org/images/arch.png?v=12"/>
</div>

# Packages
1. **Primary Package**

    The `gf` repository maintains some basic and most commonly used packages, keeping it as lightweight and simple as possible. 

1. **Community Package**

    The community packages are contributed and maintained by community members, which are hosted in `gogf` organization. Some of the community packages are separated from the `gf` repository, which are not of common usage or are with heavy dependencies. 

# Performance

Here's the most popular Golang frameworks and libraries performance testing result in `WEB Server`. Performance testing cases source codes are hosted at: https://github.com/gogf/gf-performance

## Environment

    OS   : Ubuntu 18.04 amd64
    CPU  : AMD A8-6600K x 4
    MEM  : 32GB
    GO   : v1.13.4

## Testing Tool

`ab`: Apache HTTP server benchmarking tool.

Command:
```
ab -t 10 -c 100 http://127.0.0.1:3000/hello
ab -t 10 -c 100 http://127.0.0.1:3000/query?id=10000
ab -t 10 -c 100 http://127.0.0.1:3000/json
```
The concurrency starts from `100` to `10000`.

> Run `5` times for each case of each project and pick up the best testing result.

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

## 2. Json Response
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

# Documentation

* 中文官网: https://goframe.org
* GoDoc API: https://godoc.org/github.com/gogf/gf


# Discussion
- QQ Group：[116707870](//shang.qq.com/wpa/qunwpa?idkey=195f91eceeb5d7fa76009b7cd5a4641f70bf4897b7f5a520635eb26ff17adfe7)
- WX Group：Add friend`389961817` in WeChat, commenting `GF`
- Issues：https://github.com/gogf/gf/issues

> It's recommended learning `GoFrame` through its awesome source codes and API reference.

# License

`GF` is licensed under the [MIT License](LICENSE), 100% free and open-source, forever.

# Part Of Users

- [Tencent](https://www.tencent.com/)
- [ZTE](https://www.zte.com.cn/china/)
- [Ant Financial Services](https://www.antfin.com/)
- [MedLinker](https://www.medlinker.com/)
- [KuCoin](https://www.kucoin.io/)
- [LeYouJia](https://www.leyoujia.com/)
- [IGG](https://igg.com)
- [XiMaLaYa](https://www.ximalaya.com)
- [ZYBang](https://www.zybang.com/)

> We list part of the users here, if your company or products are using `GoFrame`, please let us know [here](https://github.com/gogf/gf/issues/168).


# Contributors
This project exists thanks to all the people who contribute. [[Contributors](https://github.com/gogf/gf/graphs/contributors)].
<a href="https://github.com/gogf/gf/graphs/contributors"><img src="https://opencollective.com/goframe/contributors.svg?width=890&button=false" /></a>


# Donators

We currently accept donation by Alipay/WechatPay, please note your github/gitee account in your payment bill. If you like `GF`, why not [buy developer a cup of coffee](DONATOR.MD)?

# Sponsors
We appreciate any kind of sponsorship for `GF` development. If you've got some interesting, please contact WeChat `389961817` / Email `john@goframe.org`.



# Thanks
<a href="https://www.jetbrains.com/?from=GoFrame"><img src="https://goframe.org/images/jetbrains.png" width="100" alt="JetBrains"/></a>












