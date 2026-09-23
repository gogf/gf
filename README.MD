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

`GoFrame` requires `Go 1.23` or later.

In an existing module, add the framework:

```bash
go get -u github.com/gogf/gf/v2
```

Install the `gf` CLI:

```bash
go install github.com/gogf/gf/cmd/gf/v2@latest
```

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

## Documentation

- Official Site: [https://goframe.org](https://goframe.org)
- Official Site(en): [https://goframe.org/en](https://goframe.org/en)
- 国内镜像：[https://goframe.org.cn](https://goframe.org.cn)
- Mirror Site: [https://pages.goframe.org](https://pages.goframe.org)
- Mirror Site: [Offline Docs](https://github.com/gogf/goframe.org-pdf?tab=readme-ov-file#%E6%9C%80%E6%96%B0%E7%89%88%E6%9C%AC)
- GoDoc API: [https://pkg.go.dev/github.com/gogf/gf/v2](https://pkg.go.dev/github.com/gogf/gf/v2)
- Doc Source: [https://github.com/gogf/gf-site](https://github.com/gogf/gf-site)

## Contributors

💖 [Thanks to all the contributors who made GoFrame possible](https://github.com/gogf/gf/graphs/contributors) 💖

<a href="https://github.com/gogf/gf/graphs/contributors">
<img src="https://goframe.org/img/contributors.svg?version=v2.10.3" alt="goframe contributors"/>
</a>

## License

`GoFrame` is licensed under the [MIT License](LICENSE), 100% free and open-source, forever.
