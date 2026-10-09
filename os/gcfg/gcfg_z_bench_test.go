// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gcfg

import "testing"

const benchJsonFileName = "gcfg-z-bench-cache.json"

func BenchmarkInstance(b *testing.B) {
	Instance("bench-instance")
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Instance("bench-instance")
	}
}

func BenchmarkInstanceParallel(b *testing.B) {
	Instance("bench-instance-parallel")
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = Instance("bench-instance-parallel")
		}
	})
}

func BenchmarkGetJson(b *testing.B) {
	adapter := newBenchJsonAdapter(b)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := adapter.getJson(benchJsonFileName); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkGetJsonParallel(b *testing.B) {
	adapter := newBenchJsonAdapter(b)
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _ = adapter.getJson(benchJsonFileName)
		}
	})
}

func newBenchJsonAdapter(b *testing.B) *AdapterFile {
	adapter, err := NewAdapterFile()
	if err != nil {
		b.Fatal(err)
	}
	adapter.SetContent(`{"v":1}`, benchJsonFileName)
	b.Cleanup(func() {
		adapter.RemoveContent(benchJsonFileName)
	})
	if _, err = adapter.getJson(benchJsonFileName); err != nil {
		b.Fatal(err)
	}
	return adapter
}
