// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package dm

import (
	"testing"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/test/gtest"
)

func Test_SplitTableName(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var cases = []struct {
			table  string
			schema string
			name   string
			ok     bool
		}{
			{"user", "", "USER", true},
			{`"user"`, "", "user", true},
			{"user_$#1", "", "USER_$#1", true},
			{"SYSDBA.user", "SYSDBA", "USER", true},
			{"sysdba.user", "SYSDBA", "USER", true},
			{`"SYSDBA"."user"`, "SYSDBA", "user", true},
			{`sysdba."user"`, "SYSDBA", "user", true},
			{`"gf_schema".user`, "gf_schema", "USER", true},
			{`"SYSDBA"."my.user"`, "SYSDBA", "my.user", true},
			{`SYSDBA"."user`, "", "", false},
			{`"SYSDBA.user`, "", "", false},
			{`"".user`, "", "", false},
			{"SYSDBA.", "", "", false},
			{"a.b.c", "", "", false},
			{"user u", "", "", false},
			{"(SELECT", "", "", false},
			{"", "", "", false},
		}
		for _, c := range cases {
			schema, name, ok := splitTableName(c.table)
			t.Assert(schema, c.schema)
			t.Assert(name, c.name)
			t.Assert(ok, c.ok)
		}
	})
}

func Test_GetTableNameForFields(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		db, err := gdb.New(gdb.ConfigNode{Type: "dm", Host: "127.0.0.1", Port: "5236", Name: "SYSDBA"})
		t.AssertNil(err)
		var cases = []struct {
			table  string
			expect string
		}{
			{"t", "t"},
			{"t u", "t"},
			{"t AS u, x y", "t"},
			{"SYSDBA.t u", "SYSDBA.t"},
			{`"t_lower" u`, `"t_lower"`},
			{`"SYSDBA"."t_lower" AS u`, `"SYSDBA"."t_lower"`},
			{`"a""b".t`, ""},
			{"a.b.c", ""},
			{"SYSDBA..t", ""},
			{`"SYSDBA.t`, ""},
			{"(SELECT 1) s", ""},
		}
		for _, c := range cases {
			t.Assert(db.GetTableNameForFields(c.table), c.expect)
		}
	})
}
