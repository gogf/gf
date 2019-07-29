package gfpool

import (
	"os"
	"testing"
)

func Benchmark_os_Open_Close_ALLFlags(b *testing.B) {
	for i := 0; i < b.N; i++ {
		f, _ := os.OpenFile("/tmp/bench-test", os.O_RDWR|os.O_CREATE|os.O_TRUNC|os.O_APPEND, 0666)
		f.Close()
	}
}

func Benchmark_gfpool_Open_Close_ALLFlags(b *testing.B) {
	for i := 0; i < b.N; i++ {
		f, _ := Open("/tmp/bench-test", os.O_RDWR|os.O_CREATE|os.O_TRUNC|os.O_APPEND, 0666)
		f.Close()
	}
}

func Benchmark_os_Open_Close_RDWR(b *testing.B) {
	for i := 0; i < b.N; i++ {
		f, _ := os.OpenFile("/tmp/bench-test", os.O_RDWR, 0666)
		f.Close()
	}
}

func Benchmark_gfpool_Open_Close_RDWR(b *testing.B) {
	for i := 0; i < b.N; i++ {
		f, _ := Open("/tmp/bench-test", os.O_RDWR, 0666)
		f.Close()
	}
}

func Benchmark_os_Open_Close_RDONLY(b *testing.B) {
	for i := 0; i < b.N; i++ {
		f, _ := os.OpenFile("/tmp/bench-test", os.O_RDONLY, 0666)
		f.Close()
	}
}

func Benchmark_gfpool_Open_Close_RDONLY(b *testing.B) {
	for i := 0; i < b.N; i++ {
		f, _ := Open("/tmp/bench-test", os.O_RDONLY, 0666)
		f.Close()
	}
}
