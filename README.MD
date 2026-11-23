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

### Application framework

Runtime components required to build a complete application: entry facade, network communication, data access and common services.

| Top-level package | Sub-packages | Description |
| --- | --- | --- |
| `frame/` | [`g`](frame/g), [`gins`](frame/gins) | Application facade and singleton container. `g` exposes global entry points such as `g.Server()`, `g.DB()` and `g.Cfg()`; `gins` handles singleton registration and dependency injection |
| `net/` | [`ghttp`](net/ghttp), [`gclient`](net/gclient), [`gtcp`](net/gtcp), [`gudp`](net/gudp), [`gipv4`](net/gipv4), [`gipv6`](net/gipv6) | Network communication. `ghttp` provides routing, middleware, static files, WebSocket and graceful restart; `gclient` provides a chainable HTTP client; `gtcp`/`gudp` provide TCP/UDP server and client; `gipv4`/`gipv6` provide IP address parsing and validation |
| `database/` | [`gdb`](database/gdb), [`gredis`](database/gredis) | Data access. `gdb` is a driver-agnostic ORM supporting models, transactions, soft delete, sharding and hooks; `gredis` defines the Redis client abstraction with single, cluster and sentinel modes |
| `os/` | [`gcfg`](os/gcfg), [`glog`](os/glog), [`gcache`](os/gcache), [`gsession`](os/gsession), [`gview`](os/gview), [`gcron`](os/gcron), [`gtimer`](os/gtimer) | Application services. Config supports multiple formats and hot reload; logging supports leveled output, file rotation and trace context; cache provides in-memory LRU with expiry and pluggable adapters; session supports file, memory and Redis storage; template supports built-in functions and i18n; cron supports second-level schedules, intervals and one-shot timers |

### Service governance

Abstraction layer for microservice scenarios. It defines interfaces and common capabilities only; concrete implementations live in `contrib/`.

| Top-level package | Sub-packages | Description |
| --- | --- | --- |
| `net/` | [`gsvc`](net/gsvc), [`gsel`](net/gsel), [`gtrace`](net/gtrace), [`goai`](net/goai) | `gsvc` defines service registry and discovery interfaces; `gsel` defines load-balancing selectors; `gtrace` wraps OpenTelemetry tracing; `goai` generates OpenAPI documents |
| `os/` | [`gmetric`](os/gmetric) | Metric abstraction layer defining Counter, Gauge and Histogram interfaces. See `contrib/metric/otelmetric` for the reporting implementation |

### Utility libraries

General-purpose capabilities that can be imported independently without the framework runtime.

| Top-level package | Sub-packages | Description |
| --- | --- | --- |
| `os/` | [`gfile`](os/gfile), [`gfsnotify`](os/gfsnotify), [`gproc`](os/gproc), [`grpool`](os/grpool), [`gfpool`](os/gfpool), [`gctx`](os/gctx), [`genv`](os/genv), [`gtime`](os/gtime), [`gres`](os/gres), [`gspath`](os/gspath), [`gstructs`](os/gstructs), [`gbuild`](os/gbuild), [`gcmd`](os/gcmd), [`gmlock`](os/gmlock), [`gmutex`](os/gmutex) | System and runtime utilities. File I/O, file watching, process management, goroutine pool, file pointer pool, context, environment variables, time handling, resource management, search path, struct reflection, build info, CLI parsing, memory lock, mutex |
| `container/` | [`garray`](container/garray), [`gmap`](container/gmap), [`gset`](container/gset), [`glist`](container/glist), [`gqueue`](container/gqueue), [`gring`](container/gring), [`gtree`](container/gtree), [`gpool`](container/gpool), [`gtype`](container/gtype), [`gvar`](container/gvar) | Concurrent-safe containers. Array, map, set, list, queue, ring, tree, object pool, type wrapper and generic variable, all safe for concurrent access |
| `encoding/` | [`gjson`](encoding/gjson), [`gyaml`](encoding/gyaml), [`gtoml`](encoding/gtoml), [`gini`](encoding/gini), [`gxml`](encoding/gxml), [`gproperties`](encoding/gproperties), [`gbase64`](encoding/gbase64), [`gbinary`](encoding/gbinary), [`gcharset`](encoding/gcharset), [`gcompress`](encoding/gcompress), [`ghash`](encoding/ghash), [`ghtml`](encoding/ghtml), [`gurl`](encoding/gurl) | Codecs and format conversion. JSON, YAML, TOML, INI, XML and Properties parsing; Base64, binary, charset conversion, compression, hashing, HTML escaping and URL encoding |
| `crypto/` | [`gaes`](crypto/gaes), [`gdes`](crypto/gdes), [`grsa`](crypto/grsa), [`gmd5`](crypto/gmd5), [`gsha1`](crypto/gsha1), [`gsha256`](crypto/gsha256), [`gsha512`](crypto/gsha512), [`gcrc32`](crypto/gcrc32) | Encryption and digests. AES and DES symmetric encryption, RSA asymmetric encryption, MD5, SHA1, SHA256 and SHA512 digests, CRC32 checksums |
| `util/` | [`gconv`](util/gconv), [`gvalid`](util/gvalid), [`gmeta`](util/gmeta), [`gtag`](util/gtag), [`gmode`](util/gmode), [`gpage`](util/gpage), [`grand`](util/grand), [`guid`](util/guid), [`gutil`](util/gutil) | General utilities. Type conversion, validation, metadata, struct tags, mode bits, pagination, random numbers and unique IDs |
| `text/` | [`gstr`](text/gstr), [`gregex`](text/gregex) | Text processing. String search, replace, split and format; regex wrapper |
| `errors/` | [`gerror`](errors/gerror), [`gcode`](errors/gcode) | Error handling. `gerror` provides stack-aware error wrapping; `gcode` defines standard error codes |
| `i18n/` | [`gi18n`](i18n/gi18n) | Internationalization with multi-language resource loading and translation |
| `debug/` | [`gdebug`](debug/gdebug) | Debug helpers for printing call stacks and variable information |
| `test/` | [`gtest`](test/gtest) | Unit test helpers providing `gtest.C` blocks and fluent assertions |

### Ecosystem modules (`contrib/`)

Standalone modules that require a separate `go get`, grouped by domain.

| Area | Modules |
| --- | --- |
| [Database drivers](contrib/drivers/README.MD) | [`mysql`](contrib/drivers/mysql), [`pgsql`](contrib/drivers/pgsql), [`mssql`](contrib/drivers/mssql), [`oracle`](contrib/drivers/oracle), [`sqlite`](contrib/drivers/sqlite), [`sqlitecgo`](contrib/drivers/sqlitecgo), [`clickhouse`](contrib/drivers/clickhouse), [`dm`](contrib/drivers/dm), [`gaussdb`](contrib/drivers/gaussdb), [`mariadb`](contrib/drivers/mariadb), [`oceanbase`](contrib/drivers/oceanbase), [`tidb`](contrib/drivers/tidb) |
| Config adapters | [`apollo`](contrib/config/apollo/README.MD), [`consul`](contrib/config/consul), [`nacos`](contrib/config/nacos/README.MD), [`polaris`](contrib/config/polaris/README.MD), [`kubecm`](contrib/config/kubecm/README.MD) |
| Service registry | [`etcd`](contrib/registry/etcd/README.MD), [`consul`](contrib/registry/consul/README.MD), [`nacos`](contrib/registry/nacos/README.MD), [`polaris`](contrib/registry/polaris/README.MD), [`zookeeper`](contrib/registry/zookeeper/README.MD), [`file`](contrib/registry/file/README.MD) |
| RPC | [`grpcx`](contrib/rpc/grpcx) |
| Tracing / metrics | [`otlpgrpc`](contrib/trace/otlpgrpc), [`otlphttp`](contrib/trace/otlphttp), [`otelmetric`](contrib/metric/otelmetric/README.MD) |
| Redis adapter | [`redis`](contrib/nosql/redis/README.MD) |
| SDK utilities | [`httpclient`](contrib/sdk/httpclient) |

## Command-line tool

The `gf` CLI ships with `GoFrame` and covers daily development tasks such as project scaffolding, code generation, building and deployment.

### Project and runtime

| Command | Description |
| --- | --- |
| `gf init` | Create and initialize an empty `GoFrame` project |
| `gf run` | Run Go code with hot-compiled-like behavior, rebuilding and restarting on change |
| `gf build` | Cross-build the project for multiple platforms and architectures |
| `gf pack` | Pack any file or directory into a resource file or Go source file |
| `gf docker` | Build a Docker image for the current project |
| `gf up` | Upgrade the `GoFrame` version or tool version in the current project |

### Code generation

| Command | Description |
| --- | --- |
| `gf gen dao` | Generate `dao`, `do` and `entity` code from database tables |
| `gf gen ctrl` | Parse `api` definitions and generate controller and SDK code |
| `gf gen service` | Parse structs and methods in packages and generate `service` interface files |
| `gf gen pb` | Parse `proto` files and generate protobuf code |
| `gf gen pbentity` | Generate `entity` message files in protobuf3 format |
| `gf gen enums` | Parse Go files in the project and generate enum definition files |

### Environment and helpers

| Command | Description |
| --- | --- |
| `gf env` | Show current Golang environment variables |
| `gf fix` | Auto-fix code after upgrading to a new `GoFrame` version |
| `gf tpl` | Template parsing and building |
| `gf doc` | Download `https://pages.goframe.org/` to run locally |
| `gf install` | Install the `gf` binary to the system (may require admin permission) |
| `gf version` | Show version information of the current binary |
| `gf agents` | Link local AI coding-agent paths to `.agents` and `AGENTS.md` |

## Contributors

💖 [Thanks to all the contributors who made GoFrame possible](https://github.com/gogf/gf/graphs/contributors) 💖

<a href="https://github.com/gogf/gf/graphs/contributors">
<img src="https://goframe.org/img/contributors.svg?version=v2.10.5" alt="goframe contributors"/>
</a>

## License

`GoFrame` is licensed under the [MIT License](LICENSE), 100% free and open-source, forever.
