package main

import (
    "testing"
    "strings"
)

var s = "arwerwerwerwerwerwrwerwerwerwersefsdgsdfgsddsfgsdfsd timeout"

func Benchmark_Contains(b *testing.B) {
    for i := 0; i < b.N; i ++ {
        strings.Contains(s, "timeout")
    }
}

func Benchmark_EqualFold(b *testing.B) {
    for i := 0; i < b.N; i ++ {
        strings.EqualFold(s[len(s) - 7:], "timeout")
    }
}


