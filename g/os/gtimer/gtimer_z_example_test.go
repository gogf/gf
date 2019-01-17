// Copyright 2019 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

package gtimer_test

import (
    "fmt"
    "gitee.com/johng/gf/g/os/gtimer"
    "time"
)

func ExampleAdd() {
    now      := time.Now()
    interval := 1400*time.Millisecond
    gtimer.Add(interval, func() {
        fmt.Println(time.Now(), time.Duration(time.Now().UnixNano() - now.UnixNano()))
        now = time.Now()
    })

    select { }
}
