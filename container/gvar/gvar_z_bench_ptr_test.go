// Copyright 2018 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// go test *.go -bench=".*" -benchmem

package gvar

import "testing"

var varPtr = New(nil)

func Benchmark_Ptr_Set(b *testing.B) {
	for i := 0; i < b.N; i++ {
		varPtr.Set(i)
	}
}

func Benchmark_Ptr_Val(b *testing.B) {
	for i := 0; i < b.N; i++ {
		varPtr.Val()
	}
}

func Benchmark_Ptr_IsNil(b *testing.B) {
	for i := 0; i < b.N; i++ {
		varPtr.IsNil()
	}
}

func Benchmark_Ptr_Bytes(b *testing.B) {
	for i := 0; i < b.N; i++ {
		varPtr.Bytes()
	}
}

func Benchmark_Ptr_String(b *testing.B) {
	for i := 0; i < b.N; i++ {
		varPtr.String()
	}
}

func Benchmark_Ptr_Bool(b *testing.B) {
	for i := 0; i < b.N; i++ {
		varPtr.Bool()
	}
}

func Benchmark_Ptr_Int(b *testing.B) {
	for i := 0; i < b.N; i++ {
		varPtr.Int()
	}
}

func Benchmark_Ptr_Int8(b *testing.B) {
	for i := 0; i < b.N; i++ {
		varPtr.Int8()
	}
}

func Benchmark_Ptr_Int16(b *testing.B) {
	for i := 0; i < b.N; i++ {
		varPtr.Int16()
	}
}

func Benchmark_Ptr_Int32(b *testing.B) {
	for i := 0; i < b.N; i++ {
		varPtr.Int32()
	}
}

func Benchmark_Ptr_Int64(b *testing.B) {
	for i := 0; i < b.N; i++ {
		varPtr.Int64()
	}
}

func Benchmark_Ptr_Uint(b *testing.B) {
	for i := 0; i < b.N; i++ {
		varPtr.Uint()
	}
}

func Benchmark_Ptr_Uint8(b *testing.B) {
	for i := 0; i < b.N; i++ {
		varPtr.Uint8()
	}
}

func Benchmark_Ptr_Uint16(b *testing.B) {
	for i := 0; i < b.N; i++ {
		varPtr.Uint16()
	}
}

func Benchmark_Ptr_Uint32(b *testing.B) {
	for i := 0; i < b.N; i++ {
		varPtr.Uint32()
	}
}

func Benchmark_Ptr_Uint64(b *testing.B) {
	for i := 0; i < b.N; i++ {
		varPtr.Uint64()
	}
}

func Benchmark_Ptr_Float32(b *testing.B) {
	for i := 0; i < b.N; i++ {
		varPtr.Float32()
	}
}

func Benchmark_Ptr_Float64(b *testing.B) {
	for i := 0; i < b.N; i++ {
		varPtr.Float64()
	}
}

func Benchmark_Ptr_Ints(b *testing.B) {
	for i := 0; i < b.N; i++ {
		varPtr.Ints()
	}
}

func Benchmark_Ptr_Strings(b *testing.B) {
	for i := 0; i < b.N; i++ {
		varPtr.Strings()
	}
}

func Benchmark_Ptr_Floats(b *testing.B) {
	for i := 0; i < b.N; i++ {
		varPtr.Floats()
	}
}

func Benchmark_Ptr_Interfaces(b *testing.B) {
	for i := 0; i < b.N; i++ {
		varPtr.Interfaces()
	}
}
