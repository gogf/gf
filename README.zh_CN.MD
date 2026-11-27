[English](README.MD) | 简体中文

<div align=center>
<img src="https://goframe.org/img/logo_full.png" width="300" alt="goframe logo"/>

[![Go Reference](https://pkg.go.dev/badge/github.com/gogf/gf/v2.svg)](https://pkg.go.dev/github.com/gogf/gf/v2)
[![GoFrame CI](https://github.com/gogf/gf/actions/workflows/ci-main.yml/badge.svg)](https://github.com/gogf/gf/actions/workflows/ci-main.yml)
[![Code Coverage](https://codecov.io/gh/gogf/gf/branch/master/graph/badge.svg)](https://codecov.io/gh/gogf/gf)
[![Production Ready](https://img.shields.io/badge/production-ready-blue.svg?style=flat)](https://github.com/gogf/gf)
[![License](https://img.shields.io/github/license/gogf/gf.svg?style=flat)](https://github.com/gogf/gf)
[![Ask DeepWiki](https://deepwiki.com/badge.svg)](https://deepwiki.com/gogf/gf)

[![Release](https://img.shields.io/github/v/release/gogf/gf?style=flat)](https://github.com/gogf/gf/releases)
[![GitHub pull requests](https://img.shields.io/github/issues-pr/gogf/gf?style=flat)](https://github.com/gogf/gf/pulls)
[![GitHub closed pull requests](https://img.shields.io/github/issues-pr-closed/gogf/gf?style=flat)](https://github.com/gogf/gf/pulls?q=is%3Apr+is%3Aclosed)
[![GitHub issues](https://img.shields.io/github/issues/gogf/gf?style=flat)](https://github.com/gogf/gf/issues)
[![GitHub closed issues](https://img.shields.io/github/issues-closed/gogf/gf?style=flat)](https://github.com/gogf/gf/issues?q=is%3Aissue+is%3Aclosed)
![Stars](https://img.shields.io/github/stars/gogf/gf?style=flat)
![Forks](https://img.shields.io/github/forks/gogf/gf?style=flat)

</div>

`GoFrame`是一款强大的框架，为了更快、更轻松、更高效的项目开发。既可作为完整的应用框架，也可按需选用独立组件。

## 安装

`GoFrame`要求`Go 1.25`或更高版本。

在已有模块中添加框架：

```bash
go get -v github.com/gogf/gf/v2@latest
```

安装`gf`命令行工具：

```bash
go install github.com/gogf/gf/cmd/gf/v2@latest
```

## 文档

- 官方网站：[https://goframe.org](https://goframe.org)
- 官方网站（en）：[https://goframe.org/en](https://goframe.org/en)
- 国内镜像：[https://goframe.org.cn](https://goframe.org.cn)
- 镜像网站：[https://pages.goframe.org](https://pages.goframe.org)
- 镜像网站：[离线文档](https://github.com/gogf/goframe.org-pdf?tab=readme-ov-file#%E6%9C%80%E6%96%B0%E7%89%88%E6%9C%AC)
- `Go`包文档：[https://pkg.go.dev/github.com/gogf/gf/v2](https://pkg.go.dev/github.com/gogf/gf/v2)
- 文档源码：[https://github.com/gogf/gf-site](https://github.com/gogf/gf-site)

## 使用

创建`main.go`：

```go
package main

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

func main() {
	s := g.Server()
	s.BindHandler("/", func(r *ghttp.Request) {
		r.Response.Write("Hello World!")
	})
	s.SetPort(8000)
	s.Run()
}
```

运行：

```bash
go mod init hello
go mod tidy
go run main.go
```

打开[http://127.0.0.1:8000](http://127.0.0.1:8000)。

使用`gf`命令行工具创建项目脚手架：

```bash
gf init hello
cd hello && gf run main.go
```

## 功能地图

`GoFrame`由一组相互独立的模块组成，既可作为完整的应用框架使用，也可只引入需要的包。

### 应用框架

构建一个完整应用所需的运行时组件，涵盖入口门面、网络通信、数据访问与常用服务。

| 顶层包 | 子包 | 说明 |
| --- | --- | --- |
| `frame/` | [`g`](frame/g)、[`gins`](frame/gins) | 应用门面与单例容器。`g`提供`g.Server()`、`g.DB()`、`g.Cfg()`等全局入口，`gins`负责单例注册与核心组件实例管理 |
| `net/` | [`ghttp`](net/ghttp)、[`gclient`](net/gclient)、[`gtcp`](net/gtcp)、[`gudp`](net/gudp)、[`gipv4`](net/gipv4)、[`gipv6`](net/gipv6) | 网络通信。`ghttp`提供路由、中间件、静态文件、WebSocket、平滑重启；`gclient`提供链式 HTTP 客户端；`gtcp`/`gudp`提供 TCP/UDP 服务端与客户端；`gipv4`/`gipv6`提供 IP 地址解析与校验 |
| `database/` | [`gdb`](database/gdb)、[`gredis`](database/gredis) | 数据访问。`gdb`是驱动无关的 ORM，支持模型、事务、软删除、分表分库、钩子；`gredis`定义 Redis 客户端抽象，支持单机、集群、哨兵模式 |
| `os/` | [`gcfg`](os/gcfg)、[`glog`](os/glog)、[`gcache`](os/gcache)、[`gsession`](os/gsession)、[`gview`](os/gview)、[`gcron`](os/gcron)、[`gtimer`](os/gtimer) | 应用服务。配置支持多格式与热更新；日志支持分级输出、文件切割与链路上下文；缓存提供带过期时间的内存 LRU 与可插拔适配；会话支持文件、内存、Redis 存储；模板支持内置函数与 i18n；定时任务支持秒级 cron、间隔与一次性定时器 |

### 服务治理

面向微服务场景的抽象层，只定义接口与通用能力，具体实现由`contrib/`提供。

| 顶层包 | 子包 | 说明 |
| --- | --- | --- |
| `net/` | [`gsvc`](net/gsvc)、[`gsel`](net/gsel)、[`gtrace`](net/gtrace)、[`goai`](net/goai) | `gsvc`定义服务注册与发现接口；`gsel`定义负载均衡选择器；`gtrace`封装 OpenTelemetry 链路追踪；`goai`生成 OpenAPI 文档 |
| `os/` | [`gmetric`](os/gmetric) | 指标抽象层，定义`Counter`、`UpDownCounter`、`Histogram`以及`ObservableCounter`、`ObservableUpDownCounter`、`ObservableGauge`等指标接口，上报实现见`contrib/metric/otelmetric` |

### 工具库

可独立引入的通用能力，不依赖框架运行时，单独使用也成立。

| 顶层包 | 子包 | 说明 |
| --- | --- | --- |
| `os/` | [`gfile`](os/gfile)、[`gfsnotify`](os/gfsnotify)、[`gproc`](os/gproc)、[`grpool`](os/grpool)、[`gfpool`](os/gfpool)、[`gctx`](os/gctx)、[`genv`](os/genv)、[`gtime`](os/gtime)、[`gres`](os/gres)、[`gspath`](os/gspath)、[`gstructs`](os/gstructs)、[`gbuild`](os/gbuild)、[`gcmd`](os/gcmd)、[`gmlock`](os/gmlock)、[`gmutex`](os/gmutex) | 系统与运行时工具。文件读写、文件监听、进程管理、协程池、文件指针池、上下文、环境变量、时间处理、资源管理、搜索路径、结构体反射、构建信息、命令行解析、内存锁、互斥锁 |
| `container/` | [`garray`](container/garray)、[`gmap`](container/gmap)、[`gset`](container/gset)、[`glist`](container/glist)、[`gqueue`](container/gqueue)、[`gring`](container/gring)、[`gtree`](container/gtree)、[`gpool`](container/gpool)、[`gtype`](container/gtype)、[`gvar`](container/gvar) | 并发安全容器。数组、字典、集合、列表、队列、环、树、对象池、类型包装、通用变量，均支持并发访问 |
| `encoding/` | [`gjson`](encoding/gjson)、[`gyaml`](encoding/gyaml)、[`gtoml`](encoding/gtoml)、[`gini`](encoding/gini)、[`gxml`](encoding/gxml)、[`gproperties`](encoding/gproperties)、[`gbase64`](encoding/gbase64)、[`gbinary`](encoding/gbinary)、[`gcharset`](encoding/gcharset)、[`gcompress`](encoding/gcompress)、[`ghash`](encoding/ghash)、[`ghtml`](encoding/ghtml)、[`gurl`](encoding/gurl) | 编解码与格式转换。JSON、YAML、TOML、INI、XML、Properties 解析，Base64、二进制、字符集转换、压缩、哈希、HTML 转义、URL 编解码 |
| `crypto/` | [`gaes`](crypto/gaes)、[`gdes`](crypto/gdes)、[`grsa`](crypto/grsa)、[`gmd5`](crypto/gmd5)、[`gsha1`](crypto/gsha1)、[`gsha256`](crypto/gsha256)、[`gsha512`](crypto/gsha512)、[`gcrc32`](crypto/gcrc32) | 加密与摘要。AES、DES 对称加密，RSA 非对称加密，MD5、SHA1、SHA256、SHA512 摘要，CRC32 校验 |
| `util/` | [`gconv`](util/gconv)、[`gvalid`](util/gvalid)、[`gmeta`](util/gmeta)、[`gtag`](util/gtag)、[`gmode`](util/gmode)、[`gpage`](util/gpage)、[`grand`](util/grand)、[`guid`](util/guid)、[`gutil`](util/gutil) | 通用工具。类型转换、数据校验、元数据、结构体标签、运行模式（`develop`、`testing`、`staging`、`product`）、分页、随机数、唯一 ID |
| `text/` | [`gstr`](text/gstr)、[`gregex`](text/gregex) | 文本处理。字符串查找、替换、分割、格式化，正则表达式封装 |
| `errors/` | [`gerror`](errors/gerror)、[`gcode`](errors/gcode) | 错误处理。`gerror`提供带堆栈的错误包装，`gcode`定义标准错误码 |
| `i18n/` | [`gi18n`](i18n/gi18n) | 国际化，支持多语言资源加载与翻译 |
| `debug/` | [`gdebug`](debug/gdebug) | 调试辅助，打印调用栈与变量信息 |
| `test/` | [`gtest`](test/gtest) | 单元测试辅助，提供`gtest.C`用例块与链式断言 |

### 生态模块（`contrib/`）

需要单独`go get`的独立模块，按领域划分。

| 领域 | 模块 | 说明 |
| --- | --- | --- |
| [数据库驱动](contrib/drivers/README.MD) | [`mysql`](contrib/drivers/mysql)、[`pgsql`](contrib/drivers/pgsql)、[`mssql`](contrib/drivers/mssql)、[`oracle`](contrib/drivers/oracle)、[`sqlite`](contrib/drivers/sqlite)、[`sqlitecgo`](contrib/drivers/sqlitecgo)、[`clickhouse`](contrib/drivers/clickhouse)、[`dm`](contrib/drivers/dm)、[`gaussdb`](contrib/drivers/gaussdb)、[`mariadb`](contrib/drivers/mariadb)、[`oceanbase`](contrib/drivers/oceanbase)、[`tidb`](contrib/drivers/tidb) | `gdb`的数据库驱动实现，通过匿名导入注册到框架 |
| 配置适配 | [`apollo`](contrib/config/apollo/README.MD)、[`consul`](contrib/config/consul)、[`nacos`](contrib/config/nacos/README.MD)、[`polaris`](contrib/config/polaris/README.MD)、[`kubecm`](contrib/config/kubecm/README.MD) | `gcfg`的远程配置源适配，支持配置监听与热更新 |
| 服务注册 | [`etcd`](contrib/registry/etcd/README.MD)、[`consul`](contrib/registry/consul/README.MD)、[`nacos`](contrib/registry/nacos/README.MD)、[`polaris`](contrib/registry/polaris/README.MD)、[`zookeeper`](contrib/registry/zookeeper/README.MD)、[`file`](contrib/registry/file/README.MD) | `gsvc`的服务注册与发现实现，`file`为本地文件实现，适合开发调试 |
| RPC | [`grpcx`](contrib/rpc/grpcx) | 基于`gRPC`的 RPC 组件，集成服务注册、负载均衡与链路追踪 |
| 链路追踪 / 指标 | [`otlpgrpc`](contrib/trace/otlpgrpc)、[`otlphttp`](contrib/trace/otlphttp)、[`otelmetric`](contrib/metric/otelmetric/README.MD) | `gtrace`与`gmetric`的 OpenTelemetry 实现，支持 gRPC 与 HTTP 两种上报协议 |
| Redis 适配 | [`redis`](contrib/nosql/redis/README.MD) | `gredis`的适配实现，基于`go-redis`，同时可作为`gcache`与`gsession`的存储后端 |
| SDK 工具 | [`httpclient`](contrib/sdk/httpclient) | 面向 SDK 场景的 HTTP 客户端封装，支持自定义响应处理与请求响应日志 |

## 命令行工具

`gf`是`GoFrame`提供的项目开发命令行工具，覆盖项目创建、代码生成、构建部署等日常开发环节。

常用命令：

| 命令 | 说明 |
| --- | --- |
| `gf init` | 创建并初始化一个空的`GoFrame`项目 |
| `gf run` | 以热编译方式运行`Go`代码 |
| `gf gen dao` | 根据数据库表结构生成`dao`、`do`、`entity`代码 |
| `gf gen ctrl` | 解析`api`定义，生成控制器与`SDK`代码 |
| `gf build` | 交叉编译项目，支持多平台与多架构 |

完整命令列表与使用说明见[`cmd/gf`](cmd/gf/README.zh_CN.MD)。

## 贡献者

💖 [感谢所有使 GoFrame 成为可能的贡献者](https://github.com/gogf/gf/graphs/contributors) 💖

<a href="https://github.com/gogf/gf/graphs/contributors">
<img src="https://goframe.org/img/contributors.svg?version=v2.10.6" alt="goframe contributors"/>
</a>

## 许可证

`GoFrame`采用 [MIT License](LICENSE) 许可，100%开源和免费。
