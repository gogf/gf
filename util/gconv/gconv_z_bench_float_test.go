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

var valueFloat = float64(1.23456789)

func Benchmark_Float_To_String(b *testing.B) {
	for i := 0; i < b.N; i++ {
		String(valueFloat)
	}
}

func Benchmark_Float_To_Int(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Int(valueFloat)
	}
}

func Benchmark_Float_To_Int8(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Int8(valueFloat)
	}
}

func Benchmark_Float_To_Int16(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Int16(valueFloat)
	}
}

func Benchmark_Float_To_Int32(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Int32(valueFloat)
	}
}

func Benchmark_Float_To_Int64(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Int(valueFloat)
	}
}

func Benchmark_Float_To_Uint(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Uint(valueFloat)
	}
}

func Benchmark_Float_To_Uint8(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Uint8(valueFloat)
	}
}

func Benchmark_Float_To_Uint16(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Uint16(valueFloat)
	}
}

func Benchmark_Float_To_Uint32(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Uint32(valueFloat)
	}
}

func Benchmark_Float_To_Uint64(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Uint64(valueFloat)
	}
}

func Benchmark_Float_To_Float32(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Float32(valueFloat)
	}
}

func Benchmark_Float_To_Float64(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Float64(valueFloat)
	}
}

func Benchmark_Float_To_Time(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Time(valueFloat)
	}
}

func Benchmark_Float_To_TimeDuration(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Duration(valueFloat)
	}
}

func Benchmark_Float_To_Bytes(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Bytes(valueFloat)
	}
}

func Benchmark_Float_To_Strings(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Strings(valueFloat)
	}
}

func Benchmark_Float_To_Ints(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Ints(valueFloat)
	}
}

func Benchmark_Float_To_Floats(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Floats(valueFloat)
	}
}

func Benchmark_Float_To_Interfaces(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Interfaces(valueFloat)
	}
}
