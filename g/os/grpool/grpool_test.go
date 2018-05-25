// Copyright 2017 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

package grpool_test

import (
    "fmt"
    "runtime"
    "testing"
    "gitee.com/johng/gf/g/os/grpool"
)

func increment() {
    for i := 0; i < 1000000; i++ {}
}

func Test_GrpoolMemUsage(t *testing.T) {
    for i := 0; i < n; i++ {
        grpool.Add(increment)
    }
    mem := runtime.MemStats{}
    runtime.ReadMemStats(&mem)
    fmt.Println("mem usage:", mem.TotalAlloc/1024)
}

//func Test_GroroutineMemUsage(t *testing.T) {
//    for i := 0; i < n; i++ {
//        go increment()
//    }
//    mem := runtime.MemStats{}
//    runtime.ReadMemStats(&mem)
//    fmt.Println("mem usage:", mem.TotalAlloc/1024)
//}