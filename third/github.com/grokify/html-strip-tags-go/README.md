HTML StripTags for Go
=====================

[![Used By][used-by-svg]][used-by-link]
[![Build Status][build-status-svg]][build-status-link]
[![Go Report Card][goreport-svg]][goreport-link]
[![Docs][docs-godoc-svg]][docs-godoc-link]
[![License][license-svg]][license-link]

This is a Go package containing an extracted version of the unexported `stripTags` function in `html/template/html.go`.

:warning: This package does not protect against untrusted input. Please use [bluemonday](https://github.com/microcosm-cc/bluemonday) if you have untrusted data :warning:

## Background

* The `stripTags` function in `html/template/html.go` is very useful, however, it is not exported.
* Requests were made [on GitHub](https://github.com/golang/go/issues/5884) without success.
* This package is a repo for work done by [Christopher Hesse](https://github.com/christopherhesse) provided in this [Gist](https://gist.github.com/christopherhesse/d422447a086d373a967f).

## Installation

```bash
$ go get github.com/grokify/html-strip-tags-go
```

## Usage

```go
import(
    "github.com/gogf/gf/third/github.com/grokify/html-strip-tags-go" // => strip
)

func main() {
    original := "<h1>Hello World</h1>"
    stripped := strip.StripTags(original) // => "Hello World"
}
```

 [used-by-svg]: https://sourcegraph.com/github.com/grokify/html-strip-tags-go/-/badge.svg
 [used-by-link]: https://sourcegraph.com/github.com/grokify/html-strip-tags-go?badge
 [goreport-svg]: https://goreportcard.com/badge/github.com/grokify/html-strip-tags-go
 [goreport-link]: https://goreportcard.com/report/github.com/grokify/html-strip-tags-go
 [build-status-svg]: https://api.travis-ci.org/grokify/html-strip-tags-go.svg?branch=master
 [build-status-link]: https://travis-ci.org/grokify/html-strip-tags-go
 [coverage-status-svg]: https://coveralls.io/repos/grokify/html-strip-tags-go/badge.svg?branch=master
 [coverage-status-link]: https://coveralls.io/r/grokify/html-strip-tags-go?branch=master
 [codeclimate-status-svg]: https://codeclimate.com/github/grokify/html-strip-tags-go/badges/gpa.svg
 [codeclimate-status-link]: https://codeclimate.com/github/grokify/html-strip-tags-go
 [docs-godoc-svg]: https://img.shields.io/badge/docs-godoc-blue.svg
 [docs-godoc-link]: https://godoc.org/github.com/grokify/html-strip-tags-go
 [license-svg]: https://img.shields.io/badge/license-BSD--style+patent--grant-blue.svg
 [license-link]: https://github.com/grokify/html-strip-tags-go/blob/master/LICENSE
