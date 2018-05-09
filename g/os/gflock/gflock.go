// Copyright 2018 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

// 文件锁.
package gflock

import (
    "github.com/theckman/go-flock"
    "gitee.com/johng/gf/g/os/gfile"
)

// 文件锁
type Locker struct {
    flock *flock.Flock
}

// 创建文件锁
func New(file string) *Locker {
    path := gfile.TempDir() + gfile.Separator + file
    lock := flock.NewFlock(path)
    return &Locker{
        lock,
    }
}


