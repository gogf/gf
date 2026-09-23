// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package clickhouse

import (
	"context"
	"strings"
	"testing"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/test/gtest"
)

// Test_LocalTypeCoverage asserts that every type family the server knows has an explicit
// local type. A family missing from localTypeMap reaches the keyword matching of the core,
// which infers the type from substrings of the name and is wrong for any name that merely
// embeds a keyword, such as `Point` embedding "int".
//
// The families are read from the server rather than from a list compiled into the driver,
// so that a server upgrade adding a type is caught as well.
//
// When this test fails, add the reported families to localTypeMap rather than relying on
// the core to guess them.
func Test_LocalTypeCoverage(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		conn, err := gdb.New(gdb.ConfigNode{
			Host: "127.0.0.1", Port: "9000", User: "default", Name: "default", Type: "clickhouse",
		})
		t.AssertNil(err)

		var ctx = context.Background()
		all, err := conn.GetAll(ctx, `SELECT name FROM system.data_type_families`)
		t.AssertNil(err)
		t.AssertGT(len(all), 0)

		var missing []string
		for _, record := range all {
			name := strings.ToLower(record["name"].String())
			if _, ok := localTypeMap[name]; !ok {
				missing = append(missing, name)
			}
		}
		t.Assert(missing, nil)
	})
}
