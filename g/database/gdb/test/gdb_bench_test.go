// Copyright 2018 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

// go test *.go -bench=".*"

package test

import (
    "testing"
    "gitee.com/johng/gf/g"
)

func init() {
    // 这里需要修改为本地配置文件的目录地址
    g.Config().SetPath("/home/john/Workspace/Go/GOPATH/src/gitee.com/johng/gf/geg/frame")
}

func Benchmark_New(b *testing.B) {
    for i := 0; i < b.N; i++ {
        g.Database()
    }
}

func Benchmark_NewAndClose(b *testing.B) {
    for i := 0; i < b.N; i++ {
        g.Database().Close()
    }
}
