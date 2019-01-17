// Copyright 2019 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

package gtimer_test

import (
    "gitee.com/johng/gf/g/os/gtimer"
    "testing"
    "time"
)

var (
    timer = gtimer.New(5, 30*time.Millisecond)
)
func Benchmark_Add(b *testing.B) {
    for i := 0; i < b.N; i++ {
        timer.Add(time.Hour, func() {

        })
    }
}

func Benchmark_StartStop(b *testing.B) {
    for i := 0; i < b.N; i++ {
        timer.Start()
        timer.Stop()
    }
}
