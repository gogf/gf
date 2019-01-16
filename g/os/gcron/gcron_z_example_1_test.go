// Copyright 2019 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.


package gcron_test

import (
    "gitee.com/johng/gf/g/os/gcron"
    "gitee.com/johng/gf/g/os/glog"
    "time"
)

func ExampleCron_AddSingleton() {
    gcron.AddSingleton("* * * * * *", func() {
        glog.Println("doing")
        time.Sleep(2*time.Second)
    })
    select { }
}
