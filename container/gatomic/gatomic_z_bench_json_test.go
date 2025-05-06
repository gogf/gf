// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// go test *.go -bench=".+\_Json" -benchmem

package gatomic_test

import (
	"testing"

	"github.com/gogf/gf/v3/container/gatomic"
	"github.com/gogf/gf/v3/internal/json"
)

var (
	vBool      = gatomic.NewBool()
	vByte      = gatomic.NewByte()
	vBytes     = gatomic.NewBytes()
	vFloat32   = gatomic.NewFloat32()
	vFloat64   = gatomic.NewFloat64()
	vInt       = gatomic.NewInt()
	vInt32     = gatomic.NewInt32()
	vInt64     = gatomic.NewInt64()
	vInterface = gatomic.NewInterface()
	vString    = gatomic.NewString()
	vUint      = gatomic.NewUint()
	vUint32    = gatomic.NewUint32()
	vUint64    = gatomic.NewUint64()
)

func Benchmark_Bool_Json(b *testing.B) {
	for i := 0; i < b.N; i++ {
		json.Marshal(vBool)
	}
}

func Benchmark_Byte_Json(b *testing.B) {
	for i := 0; i < b.N; i++ {
		json.Marshal(vByte)
	}
}

func Benchmark_Bytes_Json(b *testing.B) {
	for i := 0; i < b.N; i++ {
		json.Marshal(vBytes)
	}
}

func Benchmark_Float32_Json(b *testing.B) {
	for i := 0; i < b.N; i++ {
		json.Marshal(vFloat32)
	}
}

func Benchmark_Float64_Json(b *testing.B) {
	for i := 0; i < b.N; i++ {
		json.Marshal(vFloat64)
	}
}

func Benchmark_Int_Json(b *testing.B) {
	for i := 0; i < b.N; i++ {
		json.Marshal(vInt)
	}
}

func Benchmark_Int32_Json(b *testing.B) {
	for i := 0; i < b.N; i++ {
		json.Marshal(vInt32)
	}
}

func Benchmark_Int64_Json(b *testing.B) {
	for i := 0; i < b.N; i++ {
		json.Marshal(vInt64)
	}
}

func Benchmark_Interface_Json(b *testing.B) {
	for i := 0; i < b.N; i++ {
		json.Marshal(vInterface)
	}
}

func Benchmark_String_Json(b *testing.B) {
	for i := 0; i < b.N; i++ {
		json.Marshal(vString)
	}
}

func Benchmark_Uint_Json(b *testing.B) {
	for i := 0; i < b.N; i++ {
		json.Marshal(vUint)
	}
}

func Benchmark_Uint32_Json(b *testing.B) {
	for i := 0; i < b.N; i++ {
		json.Marshal(vUint64)
	}
}
