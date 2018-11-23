// Copyright 2018 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

// go test *.go -bench=".*"

package grand_test

import (
    "gitee.com/johng/gf/g/util/grand"
    "testing"
)

var buffer = make([]byte, 8)

func Benchmark_Rand(b *testing.B) {
    for i := 0; i < b.N; i++ {
        grand.Rand(0, 999999999)
    }
}

//func Benchmark_Buffer(b *testing.B) {
//    for i := 0; i < b.N; i++ {
//        if _, err := rand.Read(buffer); err == nil {
//            binary.LittleEndian.Uint64(buffer)
//        }
//    }
//}
