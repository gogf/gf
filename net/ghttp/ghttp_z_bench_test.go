// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package ghttp_test

import (
	"strings"
	"testing"
)

func Benchmark_TrimRightCharWithStrings(b *testing.B) {
	for i := 0; i < b.N; i++ {
		path := "//////////"
		path = strings.TrimRight(path, "/")
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
