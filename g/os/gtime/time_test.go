// Copyright 2017 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

package gtime_test

import (
    "testing"
    "gitee.com/johng/gf/g/os/gtime"
)

func BenchmarkNanosecond(b *testing.B) {
    for i := 0; i < b.N; i++ {
        gtime.Nanosecond()
    }
}