// Copyright 2017 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

// go test *.go -bench=".*"

package ghash_test

import (
    "testing"
    "gitee.com/johng/gf/g/encoding/ghash"
    "gitee.com/johng/gf/g/encoding/gbinary"
)

func BenchmarkBKDRHash(b *testing.B) {
    for i := 0; i < b.N; i++ {
        ghash.BKDRHash(gbinary.EncodeInt(i))
    }
}

func BenchmarkBKDRHash64(b *testing.B) {
    for i := 0; i < b.N; i++ {
        ghash.BKDRHash64(gbinary.EncodeInt(i))
    }
}

func BenchmarkSDBMHash(b *testing.B) {
    for i := 0; i < b.N; i++ {
        ghash.SDBMHash(gbinary.EncodeInt(i))
    }
}

func BenchmarkSDBMHash64(b *testing.B) {
    for i := 0; i < b.N; i++ {
        ghash.SDBMHash64(gbinary.EncodeInt(i))
    }
}

func BenchmarkRSHash(b *testing.B) {
    for i := 0; i < b.N; i++ {
        ghash.RSHash(gbinary.EncodeInt(i))
    }
}

func BenchmarkSRSHash64(b *testing.B) {
    for i := 0; i < b.N; i++ {
        ghash.RSHash64(gbinary.EncodeInt(i))
    }
}

func BenchmarkJSHash(b *testing.B) {
    for i := 0; i < b.N; i++ {
        ghash.JSHash(gbinary.EncodeInt(i))
    }
}

func BenchmarkJSHash64(b *testing.B) {
    for i := 0; i < b.N; i++ {
        ghash.JSHash64(gbinary.EncodeInt(i))
    }
}

func BenchmarkPJWHash(b *testing.B) {
    for i := 0; i < b.N; i++ {
        ghash.PJWHash(gbinary.EncodeInt(i))
    }
}

func BenchmarkPJWHash64(b *testing.B) {
    for i := 0; i < b.N; i++ {
        ghash.PJWHash64(gbinary.EncodeInt(i))
    }
}

func BenchmarkELFHash(b *testing.B) {
    for i := 0; i < b.N; i++ {
        ghash.ELFHash(gbinary.EncodeInt(i))
    }
}

func BenchmarkELFHash64(b *testing.B) {
    for i := 0; i < b.N; i++ {
        ghash.ELFHash64(gbinary.EncodeInt(i))
    }
}

func BenchmarkDJBHash(b *testing.B) {
    for i := 0; i < b.N; i++ {
        ghash.DJBHash(gbinary.EncodeInt(i))
    }
}

func BenchmarkDJBHash64(b *testing.B) {
    for i := 0; i < b.N; i++ {
        ghash.DJBHash64(gbinary.EncodeInt(i))
    }
}

func BenchmarkAPHash(b *testing.B) {
    for i := 0; i < b.N; i++ {
        ghash.APHash(gbinary.EncodeInt(i))
    }
}

func BenchmarkAPHash64(b *testing.B) {
    for i := 0; i < b.N; i++ {
        ghash.APHash64(gbinary.EncodeInt(i))
    }
}
