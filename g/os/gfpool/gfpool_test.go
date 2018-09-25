package gfpool

import (
    "testing"
    "os"
)

func Benchmark_os_Open_Close(b *testing.B) {
    for i := 0; i < b.N; i++ {
        f, _ := os.OpenFile("/tmp/bench-test", os.O_RDWR|os.O_CREATE|os.O_TRUNC|os.O_APPEND, 0766)
        f.Close()
    }
}

func Benchmark_gfpool_Open_Close(b *testing.B) {
    for i := 0; i < b.N; i++ {
        f, _ := Open("/tmp/bench-test", os.O_RDWR|os.O_CREATE|os.O_TRUNC|os.O_APPEND, 0766)
        f.Close()
    }
}

