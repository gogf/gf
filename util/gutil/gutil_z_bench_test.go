// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// go test *.go -bench=".*" -benchmem

package gutil

import (
	"testing"
)

var (
	m1 = map[string]interface{}{
		"k1": "v1",
	}
	m2 = map[string]interface{}{
		"k2": "v2",
	}
)

func Benchmark_TryCatch(b *testing.B) {
	for i := 0; i < b.N; i++ {
		TryCatch(func() {

		}, func(err error) {

		})
	}
}

func Benchmark_MapMergeCopy(b *testing.B) {
	for i := 0; i < b.N; i++ {
		MapMergeCopy(m1, m2)
	}
}
