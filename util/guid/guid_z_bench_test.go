// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// go test *.go -bench=".*"

package guid_test

import (
	"testing"

	"github.com/gogf/gf/v2/util/guid"
)

func Benchmark_S(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			guid.S()
		}
	})
}

func Benchmark_S_Data_1(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			guid.S([]byte("123"))
		}
	})
}

func Benchmark_S_Data_2(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			guid.S([]byte("123"), []byte("456"))
		}
	})
}
