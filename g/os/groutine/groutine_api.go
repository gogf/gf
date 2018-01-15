// Copyright 2017 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

package groutine

import (
    "gitee.com/johng/gf/g/container/gset"
    "gitee.com/johng/gf/g/container/glist"
)

// 创建goroutine池管理对象
func New() *Pool {
    p := &Pool {
        jobs  : gset.NewInterfaceSet(),
        queue : glist.NewSafeList(),
        funcs : make(chan func(), 1000000),
    }
    p.loop()
    return p
}

// 添加异步任务
func (p *Pool) Add(f func()) {
    p.funcs <- f
}

// 关闭池，所有的任务将会停止，此后继续添加的任务将不会被执行
func (p *Pool) Close() {
    p.funcs <- nil
    p.jobs.Iterator(func(v interface{}){
        v.(*PoolJob).stop()
    })
}