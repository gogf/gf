// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gdb

import (
	"database/sql"
	"sync/atomic"
	"testing"

	"github.com/gogf/gf/v2/container/gmap"
)

// benchLinkDB counts the connection pools opened by the benchmarks.
type benchLinkDB struct {
	*Core
	sqlDB *sql.DB
	opens int32
}

func (d *benchLinkDB) Open(_ *ConfigNode) (*sql.DB, error) {
	atomic.AddInt32(&d.opens, 1)
	return d.sqlDB, nil
}

func newBenchLinkCore(b *testing.B, node *ConfigNode) *Core {
	sqlDB, err := sql.Open(mockRowsErrDriverName, "0:eof")
	if err != nil {
		b.Fatal(err)
	}
	b.Cleanup(func() {
		_ = sqlDB.Close()
	})
	db := &benchLinkDB{sqlDB: sqlDB}
	core := &Core{
		links:  gmap.NewKVMapWithChecker[ConfigNode, *sql.DB](linksChecker, true),
		config: node,
	}
	db.Core = core
	core.db = db
	if _, err = core.getSqlDb(true); err != nil {
		b.Fatal(err)
	}
	if opens := atomic.LoadInt32(&db.opens); opens != 1 {
		b.Fatalf("unexpected pool opens: %d", opens)
	}
	return core
}

func BenchmarkGetSqlDb(b *testing.B) {
	core := newBenchLinkCore(b, &ConfigNode{Host: "127.0.0.1", Port: "3306", Name: "bench", Type: "mock"})
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = core.getSqlDb(true)
	}
}

func BenchmarkGetSqlDbParallel(b *testing.B) {
	core := newBenchLinkCore(b, &ConfigNode{Host: "127.0.0.1", Port: "3306", Name: "bench-parallel", Type: "mock"})
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _ = core.getSqlDb(true)
		}
	})
}
