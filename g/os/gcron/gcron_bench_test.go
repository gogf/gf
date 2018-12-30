// Copyright 2018 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

package gcron_test

import (
    "gitee.com/johng/gf/g/os/gcron"
    "testing"
)

func Benchmark_Add(b *testing.B) {
    for i := 0; i < b.N; i++ {
        gcron.Add("1 1 1 1 1 1", func() {

        })
    }
}
