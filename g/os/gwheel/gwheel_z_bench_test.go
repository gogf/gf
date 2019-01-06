// Copyright 2019 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

package gwheel_test

import (
    "gitee.com/johng/gf/g/container/gtype"
    "gitee.com/johng/gf/g/os/gwheel"
    "testing"
    "time"
)

var (
    nowNanoseconds   = time.Now().UnixNano()
    entryUpdate      = gtype.NewInt64()
    entryStatus      = gtype.NewInt(gwheel.STATUS_RUNNING)
    entryTimes       = gtype.NewInt(-1)
    entryInterval    = int64(0)
    entryIsSingleton = gtype.NewBool()
)
func Benchmark_Add(b *testing.B) {
    for i := 0; i < b.N; i++ {
        // 基准测试的时候不能设置为1秒，否则大量的任务会崩掉系统
        gwheel.Add(time.Hour, func() {

        })
    }
}

// 测试最坏情况的任务检测开销
func Benchmark_RunnableCheck(b *testing.B) {
    for i := 0; i < b.N; i++ {
        if nowNanoseconds - entryUpdate.Val() >= entryInterval {
            // 是否关闭
            if entryStatus.Val() == gwheel.STATUS_CLOSED {
                continue
            }
            // 是否单例
            if entryIsSingleton.Val() {
                if entryStatus.Set(gwheel.STATUS_RUNNING) == gwheel.STATUS_RUNNING {
                    continue
                }
            }
            // 次数限制
            if entryTimes.Add(-1) == 0 {
                if  entryStatus.Set(gwheel.STATUS_CLOSED) == gwheel.STATUS_CLOSED {
                    continue
                }
            }
            entryUpdate.Set(nowNanoseconds)
        }
    }
}

