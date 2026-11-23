English | [简体中文](README.zh_CN.MD)

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

`GoFrame` is a powerful framework for faster, easier, and more efficient project development.
Use it as a full application framework, or pick individual components as needed.

## Installation

`GoFrame` requires `Go 1.25` or later.

In an existing module, add the framework:

```bash
go get -u github.com/gogf/gf/v2
```

Install the `gf` CLI:

```bash
go install github.com/gogf/gf/cmd/gf/v2@latest
```

## Documentation

- Official Site: [https://goframe.org](https://goframe.org)
- Official Site(en): [https://goframe.org/en](https://goframe.org/en)
- 国内镜像：[https://goframe.org.cn](https://goframe.org.cn)
- Mirror Site: [https://pages.goframe.org](https://pages.goframe.org)
- Mirror Site: [Offline Docs](https://github.com/gogf/goframe.org-pdf?tab=readme-ov-file#%E6%9C%80%E6%96%B0%E7%89%88%E6%9C%AC)
- GoDoc API: [https://pkg.go.dev/github.com/gogf/gf/v2](https://pkg.go.dev/github.com/gogf/gf/v2)
- Doc Source: [https://github.com/gogf/gf-site](https://github.com/gogf/gf-site)

## Usage

Create `main.go`:

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

Run it:

```bash
go mod init hello
go mod tidy
go run main.go
```

Open [http://127.0.0.1:8000](http://127.0.0.1:8000).

To scaffold a project with the `gf` CLI:

```bash
gf init hello
cd hello && gf run main.go
```

## Feature Map

`GoFrame` is a collection of independent modules. Use it as a full application framework, or import only the packages you need.

### Core components

| Top-level package | Sub-packages | Description |
| --- | --- | --- |
| `frame/` | `g`, `gins` | Application facade and singleton container |
| `net/` | `ghttp`, `gclient`, `gtcp`, `gudp`, `goai`, `gsvc`, `gsel`, `gtrace`, `gipv4`, `gipv6` | HTTP server and client, TCP/UDP, OpenAPI, service registry, load balancing, tracing, IP utilities |
| `database/` | `gdb`, `gredis` | ORM and Redis client core |
| `os/` (services) | `gcfg`, `glog`, `gcache`, `gsession`, `gview`, `gcron`, `gtimer`, `gcmd`, `gmetric`, `gmlock`, `gmutex` | Config, logging, cache, session, template, cron, CLI, metrics, memory lock, mutex |

### Utility libraries

| Top-level package | Sub-packages | Description |
| --- | --- | --- |
| `os/` (base) | `gfile`, `gfsnotify`, `gproc`, `grpool`, `gfpool`, `gctx`, `genv`, `gtime`, `gres`, `gspath`, `gstructs`, `gbuild` | Files, file watching, processes, goroutine pool, file pointer pool, context, environment, time, resources, search path, structs, build info |
| `container/` | `garray`, `gmap`, `gset`, `glist`, `gqueue`, `gring`, `gtree`, `gpool`, `gtype`, `gvar` | Concurrent-safe array, map, set, list, queue, ring, tree, object pool, type wrapper, generic variable |
| `encoding/` | `gjson`, `gyaml`, `gtoml`, `gini`, `gxml`, `gproperties`, `gbase64`, `gbinary`, `gcharset`, `gcompress`, `ghash`, `ghtml`, `gurl` | Codecs, charset conversion, compression, hashing, HTML and URL handling |
| `crypto/` | `gaes`, `gdes`, `grsa`, `gmd5`, `gsha1`, `gsha256`, `gsha512`, `gcrc32` | Symmetric and asymmetric encryption, digest algorithms |
| `util/` | `gconv`, `gvalid`, `gmeta`, `gtag`, `gmode`, `gpage`, `grand`, `guid`, `gutil` | Type conversion, validation, metadata, struct tags, mode bits, pagination, random, unique ID, general utilities |
| `text/` | `gstr`, `gregex` | String handling and regex wrapper |
| `errors/` | `gerror`, `gcode` | Stack-aware errors and error codes |
| `i18n/` | `gi18n` | Internationalization |
| `debug/` | `gdebug` | Debug information |
| `test/` | `gtest` | Unit test helpers |

### Ecosystem modules (`contrib/`)

| Area | Modules |
| --- | --- |
| [Database drivers](contrib/drivers/README.MD) | [`mysql`](contrib/drivers/mysql), [`pgsql`](contrib/drivers/pgsql), [`mssql`](contrib/drivers/mssql), [`oracle`](contrib/drivers/oracle), [`sqlite`](contrib/drivers/sqlite), [`sqlitecgo`](contrib/drivers/sqlitecgo), [`clickhouse`](contrib/drivers/clickhouse), [`dm`](contrib/drivers/dm), [`gaussdb`](contrib/drivers/gaussdb), [`mariadb`](contrib/drivers/mariadb), [`oceanbase`](contrib/drivers/oceanbase), [`tidb`](contrib/drivers/tidb) |
| Config adapters | `apollo`, `consul`, `nacos`, `polaris`, `kubecm` |
| Service registry | `etcd`, `consul`, `nacos`, `polaris`, `zookeeper`, `file` |
| RPC | `grpcx` |
| Tracing / metrics | `otlpgrpc`, `otlphttp`, `otelmetric` |
| Redis adapter | `redis` |

## Contributors

💖 [Thanks to all the contributors who made GoFrame possible](https://github.com/gogf/gf/graphs/contributors) 💖

<a href="https://github.com/gogf/gf/graphs/contributors">
<img src="https://goframe.org/img/contributors.svg?version=v2.10.5" alt="goframe contributors"/>
</a>

## License

`GoFrame` is licensed under the [MIT License](LICENSE), 100% free and open-source, forever.
