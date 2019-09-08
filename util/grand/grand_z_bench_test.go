// Copyright 2018 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// go test *.go -bench=".*"

package grand_test

import (
	"testing"

	"github.com/gogf/gf/util/grand"
)

var buffer = make([]byte, 8)

func Benchmark_Rand(b *testing.B) {
	for i := 0; i < b.N; i++ {
		grand.N(0, 999999999)
	}
}

//func Benchmark_Buffer(b *testing.B) {
//    for i := 0; i < b.N; i++ {
//        if _, err := rand.Read(buffer); err == nil {
//            binary.LittleEndian.Uint64(buffer)
//        }
//    }
//}
