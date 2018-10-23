// Copyright 2017-2018 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

package gfile

import (
    "fmt"
    "os"
)

// 文件大小(bytes)
func Size(path string) int64 {
    s, e := os.Stat(path)
    if e != nil {
        return 0
    }
    return s.Size()
}

// 格式化文件大小
func ReadableSize(path string) string {
    return FormatSize(float64(Size(path)))
}

// 格式化文件大小
func FormatSize(raw float64) string {
    var t float64 = 1024
    var d float64 = 1

    if raw < t {
        return fmt.Sprintf("%.2fB", raw/d)
    }

    d *= 1024
    t *= 1024

    if raw < t {
        return fmt.Sprintf("%.2fK", raw/d)
    }

    d *= 1024
    t *= 1024

    if raw < t {
        return fmt.Sprintf("%.2fM", raw/d)
    }

    d *= 1024
    t *= 1024

    if raw < t {
        return fmt.Sprintf("%.2fG", raw/d)
    }

    d *= 1024
    t *= 1024

    if raw < t {
        return fmt.Sprintf("%.2fT", raw/d)
    }

    d *= 1024
    t *= 1024

    if raw < t {
        return fmt.Sprintf("%.2fP", raw/d)
    }

    return "TooLarge"
}