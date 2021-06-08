// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gtimer

import (
	"testing"
	"time"
)

var (
	timer = New()
)

func Benchmark_Add(b *testing.B) {
	for i := 0; i < b.N; i++ {
		timer.Add(time.Hour, func() {

		})
	}
}

func Benchmark_PriorityQueue_Pop(b *testing.B) {
	for i := 0; i < b.N; i++ {
		timer.queue.Pop()
	}
}

func Benchmark_StartStop(b *testing.B) {
	for i := 0; i < b.N; i++ {
		timer.Start()
		timer.Stop()
	}
}
