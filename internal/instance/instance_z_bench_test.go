// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package instance_test

import (
	"testing"

	"github.com/gogf/gf/v2/internal/instance"
)

func BenchmarkGetOrSetFuncLock(b *testing.B) {
	var (
		key = "bench-get-or-set-func-lock"
		f   = func() any { return &struct{}{} }
	)
	instance.GetOrSetFuncLock(key, f)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = instance.GetOrSetFuncLock(key, f)
	}
}

func BenchmarkGetOrSetFuncLockParallel(b *testing.B) {
	var (
		key = "bench-get-or-set-func-lock-parallel"
		f   = func() any { return &struct{}{} }
	)
	instance.GetOrSetFuncLock(key, f)
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = instance.GetOrSetFuncLock(key, f)
		}
	})
}
