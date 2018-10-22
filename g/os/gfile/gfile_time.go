// Copyright 2017-2018 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

package gfile

import (
    "os"
)

// 文件修改时间(时间戳，秒)
func MTime(path string) int64 {
    s, e := os.Stat(path)
    if e != nil {
        return 0
    }
    return s.ModTime().Unix()
}

// 文件修改时间(时间戳，毫秒)
func MTimeMillisecond(path string) int64 {
    s, e := os.Stat(path)
    if e != nil {
        return 0
    }
    return int64(s.ModTime().Nanosecond()/1000000)
}
