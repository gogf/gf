package test

import (
    "strings"
    "testing"
)

var s = "/name/john///."
var c = "./"

func t1(s string) string {
    if len(s) == 0 {
        return s
    }
    for _, cut := range c {
        for s[len(s) - 1] == uint8(cut) {
            s = s[:len(s) - 1]
            if len(s) == 0 {
                return s
            }
        }
    }
    return s
}

func t2(s string) string {
    return strings.TrimRight(s, c)
}

func Benchmark_t1(b *testing.B) {
    for i := 0; i < b.N; i++ {
        t1(s)
    }
}

func Benchmark_t2(b *testing.B) {
    for i := 0; i < b.N; i++ {
        t2(s)
    }
}

