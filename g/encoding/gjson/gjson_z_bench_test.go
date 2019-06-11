// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gjson_test

import (
    "github.com/gogf/gf/g/encoding/gjson"
    "testing"
)

func Benchmark_Set1(b *testing.B) {
    for i := 0; i < b.N; i++ {
        p := gjson.New(map[string]string{
            "k1" : "v1",
            "k2" : "v2",
        })
        p.Set("k1.k11", []int{1,2,3})
    }
}

func Benchmark_Set2(b *testing.B) {
    for i := 0; i < b.N; i++ {
        p := gjson.New([]string{"a"})
        p.Set("0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0", []int{1,2,3})
    }
}

