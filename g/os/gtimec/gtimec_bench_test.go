// Copyright 2019 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

package gtimec_test

import (
    "gitee.com/johng/gf/g/os/gtimec"
    "testing"
)

func Benchmark_Add(b *testing.B) {
    for i := 0; i < b.N; i++ {
        gtimec.Add(1, func() {

        })
    }
}
