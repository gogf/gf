// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/wangyougui/gf.

package gcron_test

import (
	"context"
	"time"

	"github.com/wangyougui/gf/v2/os/gcron"
	"github.com/wangyougui/gf/v2/os/glog"
)

func ExampleCronAddSingleton() {
	gcron.AddSingleton(ctx, "* * * * * *", func(ctx context.Context) {
		glog.Print(context.TODO(), "doing")
		time.Sleep(2 * time.Second)
	})
	select {}
}
