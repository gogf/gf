// Copyright 2018 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// go test *.go -bench=".*" -benchmem

package gvar

import "testing"

var varObj = Create(nil)

func Benchmark_Obj_Set(b *testing.B) {
	for i := 0; i < b.N; i++ {
		varObj.Set(i)
	}
}

func Benchmark_Obj_Val(b *testing.B) {
	for i := 0; i < b.N; i++ {
		varObj.Val()
	}
}

func Benchmark_Obj_IsNil(b *testing.B) {
	for i := 0; i < b.N; i++ {
		varObj.IsNil()
	}
}

func Benchmark_Obj_Bytes(b *testing.B) {
	for i := 0; i < b.N; i++ {
		varObj.Bytes()
	}
}

func Benchmark_Obj_String(b *testing.B) {
	for i := 0; i < b.N; i++ {
		varObj.String()
	}
}

func Benchmark_Obj_Bool(b *testing.B) {
	for i := 0; i < b.N; i++ {
		varObj.Bool()
	}
}

func Benchmark_Obj_Int(b *testing.B) {
	for i := 0; i < b.N; i++ {
		varObj.Int()
	}
}

func Benchmark_Obj_Int8(b *testing.B) {
	for i := 0; i < b.N; i++ {
		varObj.Int8()
	}
}

func Benchmark_Obj_Int16(b *testing.B) {
	for i := 0; i < b.N; i++ {
		varObj.Int16()
	}
}

func Benchmark_Obj_Int32(b *testing.B) {
	for i := 0; i < b.N; i++ {
		varObj.Int32()
	}
}

func Benchmark_Obj_Int64(b *testing.B) {
	for i := 0; i < b.N; i++ {
		varObj.Int64()
	}
}

func Benchmark_Obj_Uint(b *testing.B) {
	for i := 0; i < b.N; i++ {
		varObj.Uint()
	}
}

func Benchmark_Obj_Uint8(b *testing.B) {
	for i := 0; i < b.N; i++ {
		varObj.Uint8()
	}
}

func Benchmark_Obj_Uint16(b *testing.B) {
	for i := 0; i < b.N; i++ {
		varObj.Uint16()
	}
}

func Benchmark_Obj_Uint32(b *testing.B) {
	for i := 0; i < b.N; i++ {
		varObj.Uint32()
	}
}

func Benchmark_Obj_Uint64(b *testing.B) {
	for i := 0; i < b.N; i++ {
		varObj.Uint64()
	}
}

func Benchmark_Obj_Float32(b *testing.B) {
	for i := 0; i < b.N; i++ {
		varObj.Float32()
	}
}

func Benchmark_Obj_Float64(b *testing.B) {
	for i := 0; i < b.N; i++ {
		varObj.Float64()
	}
}

func Benchmark_Obj_Ints(b *testing.B) {
	for i := 0; i < b.N; i++ {
		varObj.Ints()
	}
}

func Benchmark_Obj_Strings(b *testing.B) {
	for i := 0; i < b.N; i++ {
		varObj.Strings()
	}
}

func Benchmark_Obj_Floats(b *testing.B) {
	for i := 0; i < b.N; i++ {
		varObj.Floats()
	}
}

func Benchmark_Obj_Interfaces(b *testing.B) {
	for i := 0; i < b.N; i++ {
		varObj.Interfaces()
	}
}
