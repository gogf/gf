// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package sqlitecgo

import (
	"path/filepath"
	"testing"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/test/gtest"
)

func Test_Open_Extra(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		db, err := (&Driver{}).Open(&gdb.ConfigNode{
			Name:  filepath.Join(t.TempDir(), "extra.db"),
			Extra: "busy_timeout=5000&journal_mode=WAL",
		})
		t.AssertNil(err)
		defer db.Close()

		var mode string
		t.AssertNil(db.QueryRow(`PRAGMA journal_mode`).Scan(&mode))
		t.Assert(mode, "wal")

		var timeout int
		t.AssertNil(db.QueryRow(`PRAGMA busy_timeout`).Scan(&timeout))
		t.Assert(timeout, 5000)
	})
}
