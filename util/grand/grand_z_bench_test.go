// Copyright 2018 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// go test *.go -bench=".*"

package grand_test

import (
	"crypto/rand"
	"encoding/binary"
	"testing"

	"github.com/gogf/gf/util/grand"
)

var buffer = make([]byte, 8)

func Benchmark_Rand1(b *testing.B) {
	for i := 0; i < b.N; i++ {
		grand.N(0, 99)
	}
}

func Benchmark_Rand2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		grand.N(0, 999999999)
	}
}

func Benchmark_Str(b *testing.B) {
	for i := 0; i < b.N; i++ {
		grand.S(16)
	}
}

func Benchmark_StrSymbols(b *testing.B) {
	for i := 0; i < b.N; i++ {
		grand.S(16, true)
	}
}

func Benchmark_Buffer(b *testing.B) {
	for i := 0; i < b.N; i++ {
		if _, err := rand.Read(buffer); err == nil {
			binary.LittleEndian.Uint64(buffer)
		}
	}
}
