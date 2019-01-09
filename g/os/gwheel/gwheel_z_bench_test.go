// Copyright 2019 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

package gwheel_test

import (
    "gitee.com/johng/gf/g/os/gwheel"
    "testing"
    "time"
)


func Benchmark_Add(b *testing.B) {
    for i := 0; i < b.N; i++ {
        // 基准测试的时候不能设置为1秒，否则大量的任务会崩掉系统
        gwheel.Add(time.Hour, func() {

        })
    }
}
