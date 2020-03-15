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
