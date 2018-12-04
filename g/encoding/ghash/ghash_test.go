// Copyright 2017 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

// go test *.go -bench=".*"

package ghash_test

import (
    "gitee.com/johng/gf/g/encoding/ghash"
    "testing"
)

var (
    str = []byte("This is the test string for hash.")
)

func BenchmarkBKDRHash(b *testing.B) {
    for i := 0; i < b.N; i++ {
        ghash.BKDRHash(str)
    }
}

func BenchmarkBKDRHash64(b *testing.B) {
    for i := 0; i < b.N; i++ {
        ghash.BKDRHash64(str)
    }
}

func BenchmarkSDBMHash(b *testing.B) {
    for i := 0; i < b.N; i++ {
        ghash.SDBMHash(str)
    }
}

func BenchmarkSDBMHash64(b *testing.B) {
    for i := 0; i < b.N; i++ {
        ghash.SDBMHash64(str)
    }
}

func BenchmarkRSHash(b *testing.B) {
    for i := 0; i < b.N; i++ {
        ghash.RSHash(str)
    }
}

func BenchmarkSRSHash64(b *testing.B) {
    for i := 0; i < b.N; i++ {
        ghash.RSHash64(str)
    }
}

func BenchmarkJSHash(b *testing.B) {
    for i := 0; i < b.N; i++ {
        ghash.JSHash(str)
    }
}

func BenchmarkJSHash64(b *testing.B) {
    for i := 0; i < b.N; i++ {
        ghash.JSHash64(str)
    }
}

func BenchmarkPJWHash(b *testing.B) {
    for i := 0; i < b.N; i++ {
        ghash.PJWHash(str)
    }
}

func BenchmarkPJWHash64(b *testing.B) {
    for i := 0; i < b.N; i++ {
        ghash.PJWHash64(str)
    }
}

func BenchmarkELFHash(b *testing.B) {
    for i := 0; i < b.N; i++ {
        ghash.ELFHash(str)
    }
}

func BenchmarkELFHash64(b *testing.B) {
    for i := 0; i < b.N; i++ {
        ghash.ELFHash64(str)
    }
}

func BenchmarkDJBHash(b *testing.B) {
    for i := 0; i < b.N; i++ {
        ghash.DJBHash(str)
    }
}

func BenchmarkDJBHash64(b *testing.B) {
    for i := 0; i < b.N; i++ {
        ghash.DJBHash64(str)
    }
}

func BenchmarkAPHash(b *testing.B) {
    for i := 0; i < b.N; i++ {
        ghash.APHash(str)
    }
}

func BenchmarkAPHash64(b *testing.B) {
    for i := 0; i < b.N; i++ {
        ghash.APHash64(str)
    }
}
