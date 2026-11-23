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
go get -u github.com/gogf/gf/v2
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

### 核心组件

| 顶层包 | 子包 | 说明 |
| --- | --- | --- |
| `frame/` | [`g`](frame/g)、[`gins`](frame/gins) | 应用门面与单例容器 |
| `net/` | [`ghttp`](net/ghttp)、[`gclient`](net/gclient)、[`gtcp`](net/gtcp)、[`gudp`](net/gudp)、[`goai`](net/goai)、[`gsvc`](net/gsvc)、[`gsel`](net/gsel)、[`gtrace`](net/gtrace)、[`gipv4`](net/gipv4)、[`gipv6`](net/gipv6) | HTTP 服务与客户端、TCP/UDP、OpenAPI、服务注册、负载均衡、链路追踪、IP 工具 |
| `database/` | [`gdb`](database/gdb)、[`gredis`](database/gredis) | ORM 与 Redis 客户端核心 |
| `os/`（服务） | [`gcfg`](os/gcfg)、[`glog`](os/glog)、[`gcache`](os/gcache)、[`gsession`](os/gsession)、[`gview`](os/gview)、[`gcron`](os/gcron)、[`gtimer`](os/gtimer)、[`gcmd`](os/gcmd)、[`gmetric`](os/gmetric)、[`gmlock`](os/gmlock)、[`gmutex`](os/gmutex) | 配置、日志、缓存、会话、模板、定时任务、命令行、指标、内存锁、互斥锁 |

### 工具库

| 顶层包 | 子包 | 说明 |
| --- | --- | --- |
| `os/`（基础） | [`gfile`](os/gfile)、[`gfsnotify`](os/gfsnotify)、[`gproc`](os/gproc)、[`grpool`](os/grpool)、[`gfpool`](os/gfpool)、[`gctx`](os/gctx)、[`genv`](os/genv)、[`gtime`](os/gtime)、[`gres`](os/gres)、[`gspath`](os/gspath)、[`gstructs`](os/gstructs)、[`gbuild`](os/gbuild) | 文件、文件监听、进程、协程池、文件指针池、上下文、环境变量、时间、资源、搜索路径、结构体、构建信息 |
| `container/` | [`garray`](container/garray)、[`gmap`](container/gmap)、[`gset`](container/gset)、[`glist`](container/glist)、[`gqueue`](container/gqueue)、[`gring`](container/gring)、[`gtree`](container/gtree)、[`gpool`](container/gpool)、[`gtype`](container/gtype)、[`gvar`](container/gvar) | 并发安全的数组、字典、集合、列表、队列、环、树、对象池、类型包装、通用变量 |
| `encoding/` | [`gjson`](encoding/gjson)、[`gyaml`](encoding/gyaml)、[`gtoml`](encoding/gtoml)、[`gini`](encoding/gini)、[`gxml`](encoding/gxml)、[`gproperties`](encoding/gproperties)、[`gbase64`](encoding/gbase64)、[`gbinary`](encoding/gbinary)、[`gcharset`](encoding/gcharset)、[`gcompress`](encoding/gcompress)、[`ghash`](encoding/ghash)、[`ghtml`](encoding/ghtml)、[`gurl`](encoding/gurl) | 各类编解码、字符集转换、压缩、哈希、HTML 与 URL 处理 |
| `crypto/` | [`gaes`](crypto/gaes)、[`gdes`](crypto/gdes)、[`grsa`](crypto/grsa)、[`gmd5`](crypto/gmd5)、[`gsha1`](crypto/gsha1)、[`gsha256`](crypto/gsha256)、[`gsha512`](crypto/gsha512)、[`gcrc32`](crypto/gcrc32) | 对称与非对称加密、摘要算法 |
| `util/` | [`gconv`](util/gconv)、[`gvalid`](util/gvalid)、[`gmeta`](util/gmeta)、[`gtag`](util/gtag)、[`gmode`](util/gmode)、[`gpage`](util/gpage)、[`grand`](util/grand)、[`guid`](util/guid)、[`gutil`](util/gutil) | 类型转换、数据校验、元数据、结构体标签、权限位、分页、随机数、唯一 ID、通用工具 |
| `text/` | [`gstr`](text/gstr)、[`gregex`](text/gregex) | 字符串处理与正则封装 |
| `errors/` | [`gerror`](errors/gerror)、[`gcode`](errors/gcode) | 带堆栈的错误与错误码 |
| `i18n/` | [`gi18n`](i18n/gi18n) | 国际化 |
| `debug/` | [`gdebug`](debug/gdebug) | 调试信息 |
| `test/` | [`gtest`](test/gtest) | 单元测试辅助 |

### 生态模块（`contrib/`）

| 领域 | 模块 |
| --- | --- |
| [数据库驱动](contrib/drivers/README.MD) | [`mysql`](contrib/drivers/mysql)、[`pgsql`](contrib/drivers/pgsql)、[`mssql`](contrib/drivers/mssql)、[`oracle`](contrib/drivers/oracle)、[`sqlite`](contrib/drivers/sqlite)、[`sqlitecgo`](contrib/drivers/sqlitecgo)、[`clickhouse`](contrib/drivers/clickhouse)、[`dm`](contrib/drivers/dm)、[`gaussdb`](contrib/drivers/gaussdb)、[`mariadb`](contrib/drivers/mariadb)、[`oceanbase`](contrib/drivers/oceanbase)、[`tidb`](contrib/drivers/tidb) |
| 配置适配 | [`apollo`](contrib/config/apollo/README.MD)、[`consul`](contrib/config/consul)、[`nacos`](contrib/config/nacos/README.MD)、[`polaris`](contrib/config/polaris/README.MD)、[`kubecm`](contrib/config/kubecm/README.MD) |
| 服务注册 | [`etcd`](contrib/registry/etcd/README.MD)、[`consul`](contrib/registry/consul/README.MD)、[`nacos`](contrib/registry/nacos/README.MD)、[`polaris`](contrib/registry/polaris/README.MD)、[`zookeeper`](contrib/registry/zookeeper/README.MD)、[`file`](contrib/registry/file/README.MD) |
| RPC | [`grpcx`](contrib/rpc/grpcx) |
| 链路追踪 / 指标 | [`otlpgrpc`](contrib/trace/otlpgrpc)、[`otlphttp`](contrib/trace/otlphttp)、[`otelmetric`](contrib/metric/otelmetric/README.MD) |
| Redis 适配 | [`redis`](contrib/nosql/redis/README.MD) |

## 贡献者

💖 [感谢所有使 GoFrame 成为可能的贡献者](https://github.com/gogf/gf/graphs/contributors) 💖

<a href="https://github.com/gogf/gf/graphs/contributors">
<img src="https://goframe.org/img/contributors.svg?version=v2.10.5" alt="goframe contributors"/>
</a>

## 许可证

`GoFrame`采用 [MIT License](LICENSE) 许可，100%开源和免费。
