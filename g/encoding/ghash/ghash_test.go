<<<<<<< HEAD
// Copyright 2017 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.
=======
// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.
>>>>>>> upstream/master

// go test *.go -bench=".*"

package ghash_test

import (
<<<<<<< HEAD
    "testing"
    "gitee.com/johng/gf/g/encoding/ghash"
    "gitee.com/johng/gf/g/encoding/gbinary"
=======
    "github.com/gogf/gf/g/encoding/ghash"
    "testing"
)

var (
    str = []byte("This is the test string for hash.")
>>>>>>> upstream/master
)

func BenchmarkBKDRHash(b *testing.B) {
    for i := 0; i < b.N; i++ {
<<<<<<< HEAD
        ghash.BKDRHash(gbinary.EncodeInt(i))
=======
        ghash.BKDRHash(str)
>>>>>>> upstream/master
    }
}

func BenchmarkBKDRHash64(b *testing.B) {
    for i := 0; i < b.N; i++ {
<<<<<<< HEAD
        ghash.BKDRHash64(gbinary.EncodeInt(i))
=======
        ghash.BKDRHash64(str)
>>>>>>> upstream/master
    }
}

func BenchmarkSDBMHash(b *testing.B) {
    for i := 0; i < b.N; i++ {
<<<<<<< HEAD
        ghash.SDBMHash(gbinary.EncodeInt(i))
=======
        ghash.SDBMHash(str)
>>>>>>> upstream/master
    }
}

func BenchmarkSDBMHash64(b *testing.B) {
    for i := 0; i < b.N; i++ {
<<<<<<< HEAD
        ghash.SDBMHash64(gbinary.EncodeInt(i))
=======
        ghash.SDBMHash64(str)
>>>>>>> upstream/master
    }
}

func BenchmarkRSHash(b *testing.B) {
    for i := 0; i < b.N; i++ {
<<<<<<< HEAD
        ghash.RSHash(gbinary.EncodeInt(i))
=======
        ghash.RSHash(str)
>>>>>>> upstream/master
    }
}

func BenchmarkSRSHash64(b *testing.B) {
    for i := 0; i < b.N; i++ {
<<<<<<< HEAD
        ghash.RSHash64(gbinary.EncodeInt(i))
=======
        ghash.RSHash64(str)
>>>>>>> upstream/master
    }
}

func BenchmarkJSHash(b *testing.B) {
    for i := 0; i < b.N; i++ {
<<<<<<< HEAD
        ghash.JSHash(gbinary.EncodeInt(i))
=======
        ghash.JSHash(str)
>>>>>>> upstream/master
    }
}

func BenchmarkJSHash64(b *testing.B) {
    for i := 0; i < b.N; i++ {
<<<<<<< HEAD
        ghash.JSHash64(gbinary.EncodeInt(i))
=======
        ghash.JSHash64(str)
>>>>>>> upstream/master
    }
}

func BenchmarkPJWHash(b *testing.B) {
    for i := 0; i < b.N; i++ {
<<<<<<< HEAD
        ghash.PJWHash(gbinary.EncodeInt(i))
=======
        ghash.PJWHash(str)
>>>>>>> upstream/master
    }
}

func BenchmarkPJWHash64(b *testing.B) {
    for i := 0; i < b.N; i++ {
<<<<<<< HEAD
        ghash.PJWHash64(gbinary.EncodeInt(i))
=======
        ghash.PJWHash64(str)
>>>>>>> upstream/master
    }
}

func BenchmarkELFHash(b *testing.B) {
    for i := 0; i < b.N; i++ {
<<<<<<< HEAD
        ghash.ELFHash(gbinary.EncodeInt(i))
=======
        ghash.ELFHash(str)
>>>>>>> upstream/master
    }
}

func BenchmarkELFHash64(b *testing.B) {
    for i := 0; i < b.N; i++ {
<<<<<<< HEAD
        ghash.ELFHash64(gbinary.EncodeInt(i))
=======
        ghash.ELFHash64(str)
>>>>>>> upstream/master
    }
}

func BenchmarkDJBHash(b *testing.B) {
    for i := 0; i < b.N; i++ {
<<<<<<< HEAD
        ghash.DJBHash(gbinary.EncodeInt(i))
=======
        ghash.DJBHash(str)
>>>>>>> upstream/master
    }
}

func BenchmarkDJBHash64(b *testing.B) {
    for i := 0; i < b.N; i++ {
<<<<<<< HEAD
        ghash.DJBHash64(gbinary.EncodeInt(i))
=======
        ghash.DJBHash64(str)
>>>>>>> upstream/master
    }
}

func BenchmarkAPHash(b *testing.B) {
    for i := 0; i < b.N; i++ {
<<<<<<< HEAD
        ghash.APHash(gbinary.EncodeInt(i))
=======
        ghash.APHash(str)
>>>>>>> upstream/master
    }
}

func BenchmarkAPHash64(b *testing.B) {
    for i := 0; i < b.N; i++ {
<<<<<<< HEAD
        ghash.APHash64(gbinary.EncodeInt(i))
=======
        ghash.APHash64(str)
>>>>>>> upstream/master
    }
}
