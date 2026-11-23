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
| `frame/` | [`g`](frame/g), [`gins`](frame/gins) | Application facade and singleton container |
| `net/` | [`ghttp`](net/ghttp), [`gclient`](net/gclient), [`gtcp`](net/gtcp), [`gudp`](net/gudp), [`goai`](net/goai), [`gsvc`](net/gsvc), [`gsel`](net/gsel), [`gtrace`](net/gtrace), [`gipv4`](net/gipv4), [`gipv6`](net/gipv6) | HTTP server and client, TCP/UDP, OpenAPI, service registry, load balancing, tracing, IP utilities |
| `database/` | [`gdb`](database/gdb), [`gredis`](database/gredis) | ORM and Redis client core |
| `os/` (services) | [`gcfg`](os/gcfg), [`glog`](os/glog), [`gcache`](os/gcache), [`gsession`](os/gsession), [`gview`](os/gview), [`gcron`](os/gcron), [`gtimer`](os/gtimer), [`gcmd`](os/gcmd), [`gmetric`](os/gmetric), [`gmlock`](os/gmlock), [`gmutex`](os/gmutex) | Config, logging, cache, session, template, cron, CLI, metrics, memory lock, mutex |

### Utility libraries

| Top-level package | Sub-packages | Description |
| --- | --- | --- |
| `os/` (base) | [`gfile`](os/gfile), [`gfsnotify`](os/gfsnotify), [`gproc`](os/gproc), [`grpool`](os/grpool), [`gfpool`](os/gfpool), [`gctx`](os/gctx), [`genv`](os/genv), [`gtime`](os/gtime), [`gres`](os/gres), [`gspath`](os/gspath), [`gstructs`](os/gstructs), [`gbuild`](os/gbuild) | Files, file watching, processes, goroutine pool, file pointer pool, context, environment, time, resources, search path, structs, build info |
| `container/` | [`garray`](container/garray), [`gmap`](container/gmap), [`gset`](container/gset), [`glist`](container/glist), [`gqueue`](container/gqueue), [`gring`](container/gring), [`gtree`](container/gtree), [`gpool`](container/gpool), [`gtype`](container/gtype), [`gvar`](container/gvar) | Concurrent-safe array, map, set, list, queue, ring, tree, object pool, type wrapper, generic variable |
| `encoding/` | [`gjson`](encoding/gjson), [`gyaml`](encoding/gyaml), [`gtoml`](encoding/gtoml), [`gini`](encoding/gini), [`gxml`](encoding/gxml), [`gproperties`](encoding/gproperties), [`gbase64`](encoding/gbase64), [`gbinary`](encoding/gbinary), [`gcharset`](encoding/gcharset), [`gcompress`](encoding/gcompress), [`ghash`](encoding/ghash), [`ghtml`](encoding/ghtml), [`gurl`](encoding/gurl) | Codecs, charset conversion, compression, hashing, HTML and URL handling |
| `crypto/` | [`gaes`](crypto/gaes), [`gdes`](crypto/gdes), [`grsa`](crypto/grsa), [`gmd5`](crypto/gmd5), [`gsha1`](crypto/gsha1), [`gsha256`](crypto/gsha256), [`gsha512`](crypto/gsha512), [`gcrc32`](crypto/gcrc32) | Symmetric and asymmetric encryption, digest algorithms |
| `util/` | [`gconv`](util/gconv), [`gvalid`](util/gvalid), [`gmeta`](util/gmeta), [`gtag`](util/gtag), [`gmode`](util/gmode), [`gpage`](util/gpage), [`grand`](util/grand), [`guid`](util/guid), [`gutil`](util/gutil) | Type conversion, validation, metadata, struct tags, mode bits, pagination, random, unique ID, general utilities |
| `text/` | [`gstr`](text/gstr), [`gregex`](text/gregex) | String handling and regex wrapper |
| `errors/` | [`gerror`](errors/gerror), [`gcode`](errors/gcode) | Stack-aware errors and error codes |
| `i18n/` | [`gi18n`](i18n/gi18n) | Internationalization |
| `debug/` | [`gdebug`](debug/gdebug) | Debug information |
| `test/` | [`gtest`](test/gtest) | Unit test helpers |

### Ecosystem modules (`contrib/`)

| Area | Modules |
| --- | --- |
| [Database drivers](contrib/drivers/README.MD) | [`mysql`](contrib/drivers/mysql), [`pgsql`](contrib/drivers/pgsql), [`mssql`](contrib/drivers/mssql), [`oracle`](contrib/drivers/oracle), [`sqlite`](contrib/drivers/sqlite), [`sqlitecgo`](contrib/drivers/sqlitecgo), [`clickhouse`](contrib/drivers/clickhouse), [`dm`](contrib/drivers/dm), [`gaussdb`](contrib/drivers/gaussdb), [`mariadb`](contrib/drivers/mariadb), [`oceanbase`](contrib/drivers/oceanbase), [`tidb`](contrib/drivers/tidb) |
| Config adapters | [`apollo`](contrib/config/apollo/README.MD), [`consul`](contrib/config/consul), [`nacos`](contrib/config/nacos/README.MD), [`polaris`](contrib/config/polaris/README.MD), [`kubecm`](contrib/config/kubecm/README.MD) |
| Service registry | [`etcd`](contrib/registry/etcd/README.MD), [`consul`](contrib/registry/consul/README.MD), [`nacos`](contrib/registry/nacos/README.MD), [`polaris`](contrib/registry/polaris/README.MD), [`zookeeper`](contrib/registry/zookeeper/README.MD), [`file`](contrib/registry/file/README.MD) |
| RPC | [`grpcx`](contrib/rpc/grpcx) |
| Tracing / metrics | [`otlpgrpc`](contrib/trace/otlpgrpc), [`otlphttp`](contrib/trace/otlphttp), [`otelmetric`](contrib/metric/otelmetric/README.MD) |
| Redis adapter | [`redis`](contrib/nosql/redis/README.MD) |

## Contributors

💖 [Thanks to all the contributors who made GoFrame possible](https://github.com/gogf/gf/graphs/contributors) 💖

<a href="https://github.com/gogf/gf/graphs/contributors">
<img src="https://goframe.org/img/contributors.svg?version=v2.10.5" alt="goframe contributors"/>
</a>

## License

`GoFrame` is licensed under the [MIT License](LICENSE), 100% free and open-source, forever.
