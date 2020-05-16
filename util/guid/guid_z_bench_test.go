// Copyright 2018 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// go test *.go -bench=".*"

package guid_test

import (
	"github.com/gogf/gf/util/guid"
	"testing"
)

func Benchmark_S(b *testing.B) {
	for i := 0; i < b.N; i++ {
		guid.S()
	}
}

func Benchmark_New1(b *testing.B) {
	for i := 0; i < b.N; i++ {
		guid.New([]byte("123"))
	}
}

func Benchmark_New2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		guid.New([]byte("123"), []byte("456"))
	}
}

func Benchmark_New3(b *testing.B) {
	for i := 0; i < b.N; i++ {
		guid.New([]byte("123"), []byte("456"), []byte("789"))
	}
}
