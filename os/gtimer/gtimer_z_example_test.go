// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gtimer_test

import (
	"context"
	"fmt"
	"time"

	"github.com/gogf/gf/v2/os/gtimer"
)

func Example_add() {
	var (
		ctx      = context.Background()
		now      = time.Now()
		interval = 1400 * time.Millisecond
	)
	gtimer.Add(ctx, interval, func(ctx context.Context) {
		fmt.Println(time.Now(), time.Duration(time.Now().UnixNano()-now.UnixNano()))
		now = time.Now()
	})

	select {}
}
