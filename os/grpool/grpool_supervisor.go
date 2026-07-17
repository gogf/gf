// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package grpool

import (
	"context"

	"github.com/gogf/gf/v2/os/gtimer"
)

// supervisor checks the job list and fork new worker goroutine to handle the job
// if there are jobs but no workers in pool.
func (p *Pool) supervisor(ctx context.Context) {
	if p.IsClosed() {
		gtimer.Exit()
	}
	if p.IsPaused() {
		return
	}
	var changed = false
	if p.limitChanger != nil {
		changed = p.limitChanger(ctx, &p.limit)
	}
	if p.count.Val() == 0 {
		changed = true
	}

	if p.list.Size() > 0 && changed {
		limit := int(p.limit.Load())
		n := limit - p.count.Val()
		if limit <= 0 || n > 0 {
			var number = p.list.Size()
			if n > 0 && n < number {
				number = n
			}
			for i := 0; i < number; i++ {
				p.checkAndForkNewGoroutineWorker()
			}
		}
	}
}
