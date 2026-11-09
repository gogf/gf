// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package ghttp_test

import (
	"net/http"
	"strings"
	"testing"

	"github.com/gogf/gf/v2/net/ghttp"
)

// The two benchmarks below compare the http.Header lookup cost of the canonical header name
// spelling used by the ghttp.HeaderXxx constants with a non-canonical one: the non-canonical
// name has to be rebuilt for every lookup by net/http.
func Benchmark_HeaderGetWithCanonicalName(b *testing.B) {
	header := make(http.Header)
	header.Set(ghttp.HeaderWwwAuthenticate, "Basic realm=x")
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = header.Get(ghttp.HeaderWwwAuthenticate)
	}
}

func Benchmark_HeaderGetWithNonCanonicalName(b *testing.B) {
	header := make(http.Header)
	header.Set(ghttp.HeaderWwwAuthenticate, "Basic realm=x")
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = header.Get("WWW-Authenticate")
	}
}

func Benchmark_TrimRightCharWithStrings(b *testing.B) {
	for i := 0; i < b.N; i++ {
		path := "//////////"
		strings.TrimRight(path, "/")
	}
}

func Benchmark_TrimRightCharWithSlice1(b *testing.B) {
	for i := 0; i < b.N; i++ {
		path := "//////////"
		for len(path) > 0 && path[len(path)-1] == '/' {
			path = path[:len(path)-1]
		}
	}
}

func Benchmark_TrimRightCharWithSlice2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		path := "//////////"
		for {
			if length := len(path); length > 0 && path[length-1] == '/' {
				path = path[:length-1]
			} else {
				break
			}
		}
	}
}
