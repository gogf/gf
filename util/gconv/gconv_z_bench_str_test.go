// Copyright 2017-2018 gf Author(https://github.com/jin502437344/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/jin502437344/gf.

// go test *.go -bench=".*" -benchmem

package gconv

import (
	"testing"
)

var valueStr = "123456789"

func Benchmark_Str_To_String(b *testing.B) {
	for i := 0; i < b.N; i++ {
		String(valueStr)
	}
}

func Benchmark_Str_To_Int(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Int(valueStr)
	}
}

func Benchmark_Str_To_Int8(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Int8(valueStr)
	}
}

func Benchmark_Str_To_Int16(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Int16(valueStr)
	}
}

func Benchmark_Str_To_Int32(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Int32(valueStr)
	}
}

func Benchmark_Str_To_Int64(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Int(valueStr)
	}
}

func Benchmark_Str_To_Uint(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Uint(valueStr)
	}
}

func Benchmark_Str_To_Uint8(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Uint8(valueStr)
	}
}

func Benchmark_Str_To_Uint16(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Uint16(valueStr)
	}
}

func Benchmark_Str_To_Uint32(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Uint32(valueStr)
	}
}

func Benchmark_Str_To_Uint64(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Uint64(valueStr)
	}
}

func Benchmark_Str_To_Float32(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Float32(valueStr)
	}
}

func Benchmark_Str_To_Float64(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Float64(valueStr)
	}
}

func Benchmark_Str_To_Time(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Time(valueStr)
	}
}

func Benchmark_Str_To_TimeDuration(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Duration(valueStr)
	}
}

func Benchmark_Str_To_Bytes(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Bytes(valueStr)
	}
}

func Benchmark_Str_To_Strings(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Strings(valueStr)
	}
}

func Benchmark_Str_To_Ints(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Ints(valueStr)
	}
}

func Benchmark_Str_To_Floats(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Floats(valueStr)
	}
}

func Benchmark_Str_To_Interfaces(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Interfaces(valueStr)
	}
}
