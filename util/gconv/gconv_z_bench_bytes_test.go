// Copyright 2017-2018 gf Author(https://github.com/jin502437344/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/jin502437344/gf.

// go test *.go -bench "Benchmark_Bytes_To_*" -benchmem

package gconv

import (
	"testing"
	"unsafe"

	"github.com/jin502437344/gf/encoding/gbinary"
)

var valueBytes = gbinary.Encode(123456789)

func Benchmark_Bytes_To_String_Normal(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = string(valueBytes)
	}
}

func Benchmark_Bytes_To_String_Unsafe(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = *(*string)(unsafe.Pointer(&valueBytes))
	}
}

func Benchmark_Bytes_To_String(b *testing.B) {
	for i := 0; i < b.N; i++ {
		String(valueBytes)
	}
}

func Benchmark_Bytes_To_Int(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Int(valueBytes)
	}
}

func Benchmark_Bytes_To_Int8(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Int8(valueBytes)
	}
}

func Benchmark_Bytes_To_Int16(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Int16(valueBytes)
	}
}

func Benchmark_Bytes_To_Int32(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Int32(valueBytes)
	}
}

func Benchmark_Bytes_To_Int64(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Int(valueBytes)
	}
}

func Benchmark_Bytes_To_Uint(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Uint(valueBytes)
	}
}

func Benchmark_Bytes_To_Uint8(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Uint8(valueBytes)
	}
}

func Benchmark_Bytes_To_Uint16(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Uint16(valueBytes)
	}
}

func Benchmark_Bytes_To_Uint32(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Uint32(valueBytes)
	}
}

func Benchmark_Bytes_To_Uint64(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Uint64(valueBytes)
	}
}

func Benchmark_Bytes_To_Float32(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Float32(valueBytes)
	}
}

func Benchmark_Bytes_To_Float64(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Float64(valueBytes)
	}
}

func Benchmark_Bytes_To_Time(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Time(valueBytes)
	}
}

func Benchmark_Bytes_To_TimeDuration(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Duration(valueBytes)
	}
}

func Benchmark_Bytes_To_Bytes(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Bytes(valueBytes)
	}
}

func Benchmark_Bytes_To_Strings(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Strings(valueBytes)
	}
}

func Benchmark_Bytes_To_Ints(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Ints(valueBytes)
	}
}

func Benchmark_Bytes_To_Floats(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Floats(valueBytes)
	}
}

func Benchmark_Bytes_To_Interfaces(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Interfaces(valueBytes)
	}
}
