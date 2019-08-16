// Copyright 2018 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// go test *.go -bench=".*" -benchmem

package gvar

import "testing"

var vn = New(nil)

func Benchmark_Set(b *testing.B) {
	for i := 0; i < b.N; i++ {
		vn.Set(i)
	}
}

func Benchmark_Val(b *testing.B) {
	for i := 0; i < b.N; i++ {
		vn.Val()
	}
}

func Benchmark_IsNil(b *testing.B) {
	for i := 0; i < b.N; i++ {
		vn.IsNil()
	}
}

func Benchmark_Bytes(b *testing.B) {
	for i := 0; i < b.N; i++ {
		vn.Bytes()
	}
}

func Benchmark_String(b *testing.B) {
	for i := 0; i < b.N; i++ {
		vn.String()
	}
}

func Benchmark_Bool(b *testing.B) {
	for i := 0; i < b.N; i++ {
		vn.Bool()
	}
}

func Benchmark_Int(b *testing.B) {
	for i := 0; i < b.N; i++ {
		vn.Int()
	}
}

func Benchmark_Int8(b *testing.B) {
	for i := 0; i < b.N; i++ {
		vn.Int8()
	}
}

func Benchmark_Int16(b *testing.B) {
	for i := 0; i < b.N; i++ {
		vn.Int16()
	}
}

func Benchmark_Int32(b *testing.B) {
	for i := 0; i < b.N; i++ {
		vn.Int32()
	}
}

func Benchmark_Int64(b *testing.B) {
	for i := 0; i < b.N; i++ {
		vn.Int64()
	}
}

func Benchmark_Uint(b *testing.B) {
	for i := 0; i < b.N; i++ {
		vn.Uint()
	}
}

func Benchmark_Uint8(b *testing.B) {
	for i := 0; i < b.N; i++ {
		vn.Uint8()
	}
}

func Benchmark_Uint16(b *testing.B) {
	for i := 0; i < b.N; i++ {
		vn.Uint16()
	}
}

func Benchmark_Uint32(b *testing.B) {
	for i := 0; i < b.N; i++ {
		vn.Uint32()
	}
}

func Benchmark_Uint64(b *testing.B) {
	for i := 0; i < b.N; i++ {
		vn.Uint64()
	}
}

func Benchmark_Float32(b *testing.B) {
	for i := 0; i < b.N; i++ {
		vn.Float32()
	}
}

func Benchmark_Float64(b *testing.B) {
	for i := 0; i < b.N; i++ {
		vn.Float64()
	}
}

func Benchmark_Ints(b *testing.B) {
	for i := 0; i < b.N; i++ {
		vn.Ints()
	}
}

func Benchmark_Strings(b *testing.B) {
	for i := 0; i < b.N; i++ {
		vn.Strings()
	}
}

func Benchmark_Floats(b *testing.B) {
	for i := 0; i < b.N; i++ {
		vn.Floats()
	}
}

func Benchmark_Interfaces(b *testing.B) {
	for i := 0; i < b.N; i++ {
		vn.Interfaces()
	}
}
