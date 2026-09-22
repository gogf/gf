// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package dm

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/test/gtest"
)

// declaredTypes are the column types of DM8, as defined by its SQL reference. The driver
// reports no list of its own, since the type name of a column comes from the server, so the
// coverage test below declares a column of each type and reads the name back instead.
var declaredTypes = []string{
	"CHAR(10)", "CHARACTER(10)", "VARCHAR(10)", "VARCHAR2(10)",
	"NCHAR(10)", "NVARCHAR(10)", "NVARCHAR2(10)",
	"TEXT", "LONG", "LONGVARCHAR", "CLOB",
	"NUMERIC(10,2)", "DECIMAL(10,2)", "DEC(10,2)", "NUMBER(10,2)",
	"INT", "INTEGER", "PLS_INTEGER", "BIGINT", "TINYINT", "SMALLINT", "BYTE",
	"FLOAT", "DOUBLE", "DOUBLE PRECISION", "REAL",
	"BIT",
	"BINARY(4)", "VARBINARY(4)", "RAW(4)", "BLOB", "IMAGE", "LONGVARBINARY", "BFILE",
	"DATE", "TIME", "TIMESTAMP", "DATETIME",
	"TIME WITH TIME ZONE", "TIMESTAMP WITH TIME ZONE", "DATETIME WITH TIME ZONE",
	"TIMESTAMP WITH LOCAL TIME ZONE",
	"INTERVAL YEAR", "INTERVAL YEAR TO MONTH", "INTERVAL MONTH",
	"INTERVAL DAY", "INTERVAL DAY TO HOUR", "INTERVAL DAY TO MINUTE",
	"INTERVAL DAY TO SECOND", "INTERVAL HOUR", "INTERVAL HOUR TO MINUTE",
	"INTERVAL HOUR TO SECOND", "INTERVAL MINUTE", "INTERVAL MINUTE TO SECOND",
	"INTERVAL SECOND",
}

// Test_LocalTypeCoverage declares a column of every DM8 type and asserts that the type name
// the driver reports for it has an explicit local type. A name missing from localTypeMap
// reaches the keyword matching of the core, which infers the type from substrings of the
// name and is wrong for any name that merely embeds a keyword, such as the INTERVAL types
// embedding "int".
//
// When this test fails, DM reports a name the map does not know. Add it to localTypeMap
// rather than relying on the core to guess it.
func Test_LocalTypeCoverage(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		conn, err := gdb.New(gdb.ConfigNode{
			Host: "127.0.0.1", Port: "5236", User: "SYSDBA", Pass: "SYSDBA001",
			Name: "SYSDBA", Type: "dm", Charset: "utf8",
		})
		t.AssertNil(err)

		var (
			ctx     = context.Background()
			table   = "test_local_type_coverage"
			columns = make([]string, 0, len(declaredTypes))
		)
		for i, declaredType := range declaredTypes {
			columns = append(columns, fmt.Sprintf("c%d %s", i, declaredType))
		}
		if _, err = conn.Exec(ctx, "DROP TABLE "+table); err != nil {
			// The table is absent on the first run, which is not an error here.
			err = nil
		}
		_, err = conn.Exec(ctx, fmt.Sprintf("CREATE TABLE %s (%s)", table, strings.Join(columns, ", ")))
		t.AssertNil(err)
		defer func() {
			_, err := conn.Exec(ctx, "DROP TABLE "+table)
			t.AssertNil(err)
		}()

		fields, err := conn.TableFields(ctx, table)
		t.AssertNil(err)
		t.Assert(len(fields), len(declaredTypes))

		var missing []string
		for _, field := range fields {
			typeName, _ := conn.GetCore().GetFormattedDBTypeNameForField(field.Type)
			if _, ok := localTypeMap[typeName]; !ok {
				missing = append(missing, fmt.Sprintf("%s (declared %s)", typeName, field.Type))
			}
		}
		t.Assert(missing, nil)
	})
}
