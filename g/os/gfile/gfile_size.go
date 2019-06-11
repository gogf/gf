// Copyright 2017-2018 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gfile

import (
    "fmt"
    "os"
)

// Size returns the size of file specified by <path> in byte.
func Size(path string) int64 {
    s, e := os.Stat(path)
    if e != nil {
        return 0
    }
    return s.Size()
}

// ReadableSize formats size of file given by <path>, for more human readable.
func ReadableSize(path string) string {
    return FormatSize(float64(Size(path)))
}

// FormatSize formats size <raw> for more human readable.
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