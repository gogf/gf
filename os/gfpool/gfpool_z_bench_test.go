// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gfpool

import (
	"os"
	"testing"
)

func Benchmark_OS_Open_Close_ALLFlags(b *testing.B) {
	for i := 0; i < b.N; i++ {
		f, err := os.OpenFile("/tmp/bench-test", os.O_RDWR|os.O_CREATE|os.O_TRUNC|os.O_APPEND, 0666)
		if err != nil {
			panic(err)
		}
		f.Close()
	}
}

func Benchmark_GFPool_Open_Close_ALLFlags(b *testing.B) {
	for i := 0; i < b.N; i++ {
		f, err := Open("/tmp/bench-test", os.O_RDWR|os.O_CREATE|os.O_TRUNC|os.O_APPEND, 0666)
		if err != nil {
			panic(err)
		}
		f.Close()
	}
}

func Benchmark_OS_Open_Close_RDWR(b *testing.B) {
	for i := 0; i < b.N; i++ {
		f, err := os.OpenFile("/tmp/bench-test", os.O_RDWR, 0666)
		if err != nil {
			panic(err)
		}
		f.Close()
	}
}

func Benchmark_GFPool_Open_Close_RDWR(b *testing.B) {
	for i := 0; i < b.N; i++ {
		f, err := Open("/tmp/bench-test", os.O_RDWR, 0666)
		if err != nil {
			panic(err)
		}
		f.Close()
	}
}

func Benchmark_OS_Open_Close_RDONLY(b *testing.B) {
	for i := 0; i < b.N; i++ {
		f, err := os.OpenFile("/tmp/bench-test", os.O_RDONLY, 0666)
		if err != nil {
			panic(err)
		}
		f.Close()
	}
}

func Benchmark_GFPool_Open_Close_RDONLY(b *testing.B) {
	for i := 0; i < b.N; i++ {
		f, err := Open("/tmp/bench-test", os.O_RDONLY, 0666)
		if err != nil {
			panic(err)
		}
		f.Close()
	}
}

func Benchmark_Stat(b *testing.B) {
	f, err := os.Create("/tmp/bench-test-stat")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	for i := 0; i < b.N; i++ {
		f.Stat()
	}
}
