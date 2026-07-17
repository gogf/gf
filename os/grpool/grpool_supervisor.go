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
	changed := p.limitChanged.Load()
	if p.limitChanger != nil {
		changed = p.limitChanger(ctx, &p.limit)
		if v := p.limit.Load(); v <= 0 && v != -1 {
			p.limit.Store(-1)
			changed = true
		}
	}
	if !changed && p.count.Val() == 0 {
		changed = true
	}

	if changed {
		if p.list.Size() > 0 {
			limit := p.limit.Load()
			if limit == -1 {
				size := p.list.Size()
				for i := 0; i < size; i++ {
					p.checkAndForkNewGoroutineWorker()
				}
				return
			}
			n := limit - int64(p.count.Val())
			if n <= 0 {
				return
			}
			number := p.list.Size()
			if n < int64(number) {
				number = int(n)
			}
			for i := 0; i < number; i++ {
				p.checkAndForkNewGoroutineWorker()
			}
		}
		p.limitChanged.Store(false)
	}
}
