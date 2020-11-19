# GoFrame
[![Go Doc](https://godoc.org/github.com/gogf/gf?status.svg)](https://godoc.org/github.com/gogf/gf) 
[![Build Status](https://travis-ci.org/gogf/gf.svg?branch=master)](https://travis-ci.org/gogf/gf) 
[![Go Report](https://goreportcard.com/badge/github.com/gogf/gf?v=1)](https://goreportcard.com/report/github.com/gogf/gf)
[![Code Coverage](https://codecov.io/gh/gogf/gf/branch/master/graph/badge.svg)](https://codecov.io/gh/gogf/gf/branch/master)
[![Production Ready](https://img.shields.io/badge/production-ready-blue.svg)](https://github.com/gogf/gf)
[![License](https://img.shields.io/github/license/gogf/gf.svg?style=flat)](https://github.com/gogf/gf)

[English](README.MD) | 简体中文

`GF(Go Frame)`是一款模块化、高性能、生产级的Go基础开发框架。
实现了比较完善的基础设施建设以及开发工具链，提供了常用的基础开发模块，
如：缓存、日志、队列、数组、集合、容器、定时器、命令行、内存锁、对象池、
配置管理、资源管理、数据校验、数据编码、定时任务、数据库ORM、TCP/UDP组件、进程管理/通信等等。
并提供了Web服务开发的系列核心组件，如：Router、Cookie、Session、Middleware、服务注册、模板引擎等等，
支持热重启、热更新、域名绑定、TLS/HTTPS、Rewrite等特性。

> 如果您初识`Go`语言，您可以将`GoFrame`类似于`PHP`中的`Laravel`, `Java`中的`SpringBoot`或者`Python`中的`Django`。

# 特点
* 模块化、松耦合设计；
* 模块丰富、开箱即用；
* 简便易用、易于维护；
* 高代码质量、高单元测试覆盖率；
* 社区活跃，大牛谦逊低调脾气好；
* 详尽的开发文档及示例；
* 完善的本地中文化支持；
* 设计为团队及企业使用；

# 地址
- **主库**：https://github.com/gogf/gf 
- **码云**：https://gitee.com/johng/gf 

# 安装
```html
go get -u -v github.com/gogf/gf
```
推荐使用 `go.mod`:
```
require github.com/gogf/gf latest
```

# 限制
```shell
golang版本 >= 1.11
```



# 架构
<div align=center>
<img src="https://goframe.org/images/arch.png?v=12"/>
</div>

# 模块

1. **核心模块**

    `GoFrame`提供了一些基础的、常用的模块，简单、易用和轻量级，并保持极少的外部依赖，这些模块由`gf`主仓库细致维护。

1. **社区模块**

    社区模块主要由社区贡献并维护，大部分也是由`gf`主仓库的贡献者提供及维护，存放于`gogf`空间下，与`gf`主仓库处于同一级别。有的社区模块是从`gf`主仓库中剥离出来单独维护的模块，这些模块并不是特别常用，或者对外部依赖较重。


# 性能

以下是目前最流行的`WEB Server` Golang框架/类库性能测试结果。
性能测试用例源代码仓库: https://github.com/gogf/gf-performance

## 环境:

    OS   : Ubuntu 18.04 amd64
    CPU  : AMD A8-6600K x 4
    MEM  : 32GB
    GO   : v1.13.4

## 工具

`ab`: Apache HTTP server benchmarking tool.

测试命令:
```
ab -t 10 -c 100 http://127.0.0.1:3000/hello
ab -t 10 -c 100 http://127.0.0.1:3000/query?id=10000
ab -t 10 -c 100 http://127.0.0.1:3000/json
```
并发客户端数量从 `100` 递增到 `10000`。

> 每个项目的每个用例均运行`5`次，取最优的结果展示。

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


# 文档

开发文档：https://goframe.org

接口文档：https://godoc.org/github.com/gogf/gf

# 帮助
- QQ交流群：[116707870](//shang.qq.com/wpa/qunwpa?idkey=195f91eceeb5d7fa76009b7cd5a4641f70bf4897b7f5a520635eb26ff17adfe7)
- WX交流群：微信添加`389961817`备注`GF`
- 主库ISSUE：https://github.com/gogf/gf/issues

> 建议通过阅读`GoFrame`的源码以及API文档深度学习`GoFrame`，了解更多的精妙设计。

# 协议

`GF` 使用非常友好的 [MIT](LICENSE) 开源协议进行发布，永久`100%`开源免费。

# 用户

- [腾讯科技](https://www.tencent.com/)
- [中兴科技](https://www.zte.com.cn/china/)
- [蚂蚁金服](https://www.antfin.com/)
- [医联科技](https://www.medlinker.com/)
- [库币科技](https://www.kucoin.io/)
- [乐有家](https://www.leyoujia.com/)
- [IGG](https://igg.com)
- [喜马拉雅](https://www.ximalaya.com)
- [作业帮](https://www.zybang.com/)

> 在这里只列举了部分知名的用户，如果您的企业或者产品正在使用`GoFrame`，欢迎到 [这里](https://github.com/gogf/gf/issues/168) 留言。

# 贡献

感谢所有参与`GoFrame`开发的贡献者。 [[贡献者列表](https://github.com/gogf/gf/graphs/contributors)].
<a href="https://github.com/gogf/gf/graphs/contributors"><img src="https://opencollective.com/goframe/contributors.svg?width=890&button=false" /></a>


# 捐赠

如果您喜欢`GF`，要不[给开发者来杯咖啡吧](DONATOR.MD)！
请在捐赠时备注您的`github`/`gitee`账号名称。

# 赞助

赞助支持`GF`框架的快速研发，如果您感兴趣，请联系 微信 `389961817` / 邮件 `john@goframe.org`。

# 感谢
<a href="https://www.jetbrains.com/?from=GoFrame"><img src="https://goframe.org/images/jetbrains.png" width="100" alt="JetBrains"/></a>

