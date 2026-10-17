// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package mssql_test

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/test/gtest"
)

func Test_Mssql_TableFields_View(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)
	view := table + "_view"
	_, err := db.Exec(ctx, fmt.Sprintf("CREATE VIEW [%s] AS SELECT ID, PASSPORT, NICKNAME FROM [%s]", view, table))
	gtest.AssertNil(err)
	defer db.Exec(ctx, fmt.Sprintf("DROP VIEW [%s]", view))

	gtest.C(t, func(t *gtest.T) {
		fields, err := db.TableFields(ctx, view)
		t.AssertNil(err)
		t.Assert(len(fields), 3)
		t.Assert(fields["ID"].Index, 0)
		t.Assert(fields["NICKNAME"].Type, "varchar(45)")
	})

	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(view).Data(g.Map{"nickname": "view_name"}).Where("id", 3).Update()
		t.AssertNil(err)
		n, err := result.RowsAffected()
		t.AssertNil(err)
		t.Assert(n, 1)

		value, err := db.Model(table).Where("id", 3).Value("NICKNAME")
		t.AssertNil(err)
		t.Assert(value, "view_name")
	})
}

func Test_Mssql_TableFields_NotExist(t *testing.T) {
	table := fmt.Sprintf("user_%d", gtime.TimestampNano())
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		fields, err := db.TableFields(ctx, table)
		t.AssertNE(err, nil)
		t.Assert(len(fields), 0)

		_, err = db.Model(table).Data(g.Map{"id": 1, "passport": "p1"}).Insert()
		t.AssertNE(err, nil)
	})

	gtest.C(t, func(t *gtest.T) {
		createTable(table)

		fields, err := db.TableFields(ctx, table)
		t.AssertNil(err)
		t.Assert(len(fields), 7)

		_, err = db.Model(table).Data(g.Map{"id": 1, "passport": "p1"}).Insert()
		t.AssertNil(err)
		value, err := db.Model(table).Where("id", 1).Value("PASSPORT")
		t.AssertNil(err)
		t.Assert(value, "p1")
	})
}

// Test_Mssql_Timezone_Extra tests the time zone in which the times of columns without time zone are
// read and written: the local one by default, or the one that the timezone configuration or the
// timezone parameter of Extra names, the latter if both are set.
func Test_Mssql_Timezone_Extra(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		value, err := db.Model(table).Where("id", 1).Value("CREATE_TIME")
		t.AssertNil(err)
		t.Assert(value.GTime().Location().String(), time.Local.String())
		t.Assert(value.GTime().Format("Y-m-d H:i:s"), CreateTime)
	})

	gtest.C(t, func(t *gtest.T) {
		node := *db.GetConfig()
		node.Extra = "timezone=Asia/Tokyo"
		tokyoDb, err := gdb.New(node)
		t.AssertNil(err)
		defer tokyoDb.Close(ctx)

		value, err := tokyoDb.Model(table).Where("id", 1).Value("CREATE_TIME")
		t.AssertNil(err)
		t.Assert(value.GTime().Location().String(), "Asia/Tokyo")
		t.Assert(value.GTime().Format("Y-m-d H:i:s"), CreateTime)
	})

	gtest.C(t, func(t *gtest.T) {
		for _, c := range []struct {
			timezone string
			extra    string
			location string
		}{
			{"", "timezone=Asia/Tokyo", "Asia/Tokyo"},
			{"Asia/Tokyo", "", "Asia/Tokyo"},
			{"America/New_York", "timezone=Asia/Tokyo", "Asia/Tokyo"},
			{"America/New_York", "", "America/New_York"},
		} {
			node := *db.GetConfig()
			node.Timezone, node.Extra = c.timezone, c.extra
			zoneDb, err := gdb.New(node)
			t.AssertNil(err)
			defer zoneDb.Close(ctx)
			location, err := time.LoadLocation(c.location)
			t.AssertNil(err)

			value, err := zoneDb.Model(table).Where("id", 3).Value("CREATE_TIME")
			t.AssertNil(err)
			t.Assert(value.GTime().Location().String(), c.location)
			t.Assert(value.GTime().Format("Y-m-d H:i:s"), CreateTime)

			written := time.Date(2024, 1, 2, 3, 4, 5, 0, time.UTC)
			_, err = zoneDb.Model(table).Data(g.Map{"create_time": written}).Where("id", 1).Update()
			t.AssertNil(err)
			_, err = zoneDb.Exec(
				ctx, fmt.Sprintf("UPDATE [%s] SET CREATE_TIME = ? WHERE ID = ?", table), written, 2,
			)
			t.AssertNil(err)
			for _, id := range []int{1, 2} {
				raw, err := db.GetValue(ctx, fmt.Sprintf(
					"SELECT CONVERT(varchar(19), CREATE_TIME, 120) FROM [%s] WHERE ID = ?", table,
				), id)
				t.AssertNil(err)
				t.Assert(raw, written.In(location).Format("2006-01-02 15:04:05"))
				value, err := zoneDb.Model(table).Where("id", id).Value("CREATE_TIME")
				t.AssertNil(err)
				t.Assert(value.Time().Equal(written), true)
			}
		}
	})
}

func Test_Mssql_Lock_TableHint(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		for lock, hint := range map[string]string{
			gdb.LockForUpdate:           "WITH (UPDLOCK, ROWLOCK)",
			gdb.LockForUpdateNowait:     "WITH (UPDLOCK, ROWLOCK, NOWAIT)",
			gdb.LockForUpdateSkipLocked: "WITH (UPDLOCK, ROWLOCK, READPAST)",
			gdb.LockInShareMode:         "WITH (REPEATABLEREAD, ROWLOCK)",
			gdb.LockForShare:            "WITH (REPEATABLEREAD, ROWLOCK)",
			gdb.LockForShareNowait:      "WITH (REPEATABLEREAD, ROWLOCK, NOWAIT)",
			gdb.LockForShareSkipLocked:  "WITH (REPEATABLEREAD, ROWLOCK, READPAST)",
			gdb.LockWithUpdLockHoldLock: gdb.LockWithUpdLockHoldLock,
			"for update":                "WITH (UPDLOCK, ROWLOCK)",
		} {
			err := db.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
				sql, err := gdb.CatchSQL(ctx, func(ctx context.Context) error {
					one, err := tx.Model(table).Ctx(ctx).Lock(lock).Where("id", 1).One()
					t.Assert(one["PASSPORT"], "user_1")
					return err
				})
				t.Assert(len(sql), 1)
				unquoted := strings.ReplaceAll(sql[0], `"`, "")
				t.Assert(strings.Contains(unquoted, table+" "+hint+" WHERE"), true)
				return err
			})
			t.AssertNil(err)
		}
	})

	gtest.C(t, func(t *gtest.T) {
		err := db.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
			sql, err := gdb.CatchSQL(ctx, func(ctx context.Context) error {
				all, err := tx.Model(table, "a").Ctx(ctx).
					LeftJoin(table, "b", "a.ID=b.ID").
					Fields("a.ID", "b.PASSPORT").
					LockUpdate().
					Where("a.ID", g.Slice{1, 2}).
					Order("a.ID").
					All()
				t.Assert(len(all), 2)
				t.Assert(all[1]["PASSPORT"], "user_2")
				return err
			})
			t.AssertGT(len(sql), 0)
			hintIndex := strings.Index(sql[len(sql)-1], "WITH (UPDLOCK, ROWLOCK)")
			joinIndex := strings.Index(sql[len(sql)-1], "LEFT JOIN")
			t.AssertGT(hintIndex, 0)
			t.AssertLT(hintIndex, joinIndex)
			return err
		})
		t.AssertNil(err)
	})
}

// dialectSortedFieldNames returns the names of `fields` in the order of their index.
func dialectSortedFieldNames(fields map[string]*gdb.TableField) []string {
	names := make([]string, len(fields))
	for name, field := range fields {
		if field.Index >= 0 && field.Index < len(names) {
			names[field.Index] = name
		}
	}
	return names
}

// Test_Mssql_TableFields_Schema tests that the fields of a table are looked up in its schema, the
// default schema of the user if the name has none, so that a table of the same name in another
// schema is never merged in, and that a model writes to the table of a qualified name.
func Test_Mssql_TableFields_Schema(t *testing.T) {
	var (
		table  = insertTableName("t_schema")
		schema = "s" + table
	)
	insertExec(
		fmt.Sprintf(`CREATE TABLE [%s] (ID int IDENTITY(1,1) NOT NULL PRIMARY KEY, PASSPORT varchar(45), NICKNAME varchar(45))`, table),
		fmt.Sprintf(`CREATE SCHEMA [%s]`, schema),
		fmt.Sprintf(`CREATE TABLE [%s].[%s] (ID int NOT NULL PRIMARY KEY, OTHER_COL varchar(45), SECRET varchar(45), NICKNAME varchar(45))`, schema, table),
	)
	defer dropTable(table)
	defer db.Exec(ctx, fmt.Sprintf(`DROP SCHEMA [%s]`, schema))
	defer db.Exec(ctx, fmt.Sprintf(`DROP TABLE [%s].[%s]`, schema, table))

	gtest.C(t, func(t *gtest.T) {
		for _, name := range []string{table, `"` + table + `"`, "dbo." + table, `"dbo"."` + table + `"`, "[dbo].[" + table + "]", "test.dbo." + table} {
			fields, err := db.TableFields(ctx, name)
			t.AssertNil(err)
			t.Assert(dialectSortedFieldNames(fields), []string{"ID", "PASSPORT", "NICKNAME"})
			t.Assert(fields["ID"].Extra, "IDENTITY")
		}
		for _, name := range []string{schema + "." + table, `"` + schema + `"."` + table + `"`, "[" + schema + "].[" + table + "]", "test." + schema + "." + table} {
			fields, err := db.TableFields(ctx, name)
			t.AssertNil(err)
			t.Assert(dialectSortedFieldNames(fields), []string{"ID", "OTHER_COL", "SECRET", "NICKNAME"})
			t.Assert(fields["ID"].Key, "PRI")
			t.Assert(fields["ID"].Extra, "")
		}
	})

	gtest.C(t, func(t *gtest.T) {
		has, err := db.GetCore().HasField(ctx, table, "secret")
		t.AssertNil(err)
		t.Assert(has, false)
		has, err = db.GetCore().HasField(ctx, schema+"."+table, "secret")
		t.AssertNil(err)
		t.Assert(has, true)
	})

	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Data(g.Map{"passport": "p1", "nickname": "n1", "secret": "s1"}).Insert()
		t.AssertNil(err)
		id, err := result.LastInsertId()
		t.AssertNil(err)
		t.Assert(id, 1)
		one, err := db.Model(table).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one, g.Map{"ID": 1, "PASSPORT": "p1", "NICKNAME": "n1"})
	})

	gtest.C(t, func(t *gtest.T) {
		qualified := schema + "." + table
		_, err := db.Model(qualified).Data(g.Map{"id": 1, "other_col": "o1", "secret": "s1", "passport": "p1"}).Insert()
		t.AssertNil(err)
		_, err = db.Model(qualified).Data(g.Map{"nickname": "n1"}).Where("id", 1).Update()
		t.AssertNil(err)
		_, err = db.Model(qualified).Data(g.Map{"id": 2, "secret": "s2"}).Save()
		t.AssertNil(err)
		_, err = db.Model("test." + qualified).Data(g.Map{"id": 3, "secret": "s3"}).Insert()
		t.AssertNil(err)
		all, err := db.Model(qualified).Order("id").All()
		t.AssertNil(err)
		t.Assert(all, g.List{
			{"ID": 1, "OTHER_COL": "o1", "SECRET": "s1", "NICKNAME": "n1"},
			{"ID": 2, "OTHER_COL": nil, "SECRET": "s2", "NICKNAME": nil},
			{"ID": 3, "OTHER_COL": nil, "SECRET": "s3", "NICKNAME": nil},
		})
		count, err := db.Model(table).Count()
		t.AssertNil(err)
		t.Assert(count, 1)
	})
}

// Test_Mssql_TableFields_QuotedName tests the fields of a table whose name holds a single quote
// and a closing bracket, given quoted with double quotes or brackets, and of a table in a schema
// whose name holds a closing bracket, a dot and a single quote, given with its schema and with its
// database.
func Test_Mssql_TableFields_QuotedName(t *testing.T) {
	var (
		table   = insertTableName("t'quote]")
		escaped = strings.ReplaceAll(table, "]", "]]")
		schema  = insertTableName("s]x.y'z")
		inner   = strings.ReplaceAll(schema, "]", "]]")
	)
	insertExec(
		fmt.Sprintf(`CREATE TABLE [%s] (ID int NOT NULL PRIMARY KEY, NAME varchar(45))`, escaped),
		fmt.Sprintf(`CREATE SCHEMA [%s]`, inner),
		fmt.Sprintf(`CREATE TABLE [%s].[%s] (ID int NOT NULL PRIMARY KEY, OTHER varchar(45))`, inner, escaped),
	)
	defer db.Exec(ctx, fmt.Sprintf(`DROP TABLE [%s]`, escaped))
	defer db.Exec(ctx, fmt.Sprintf(`DROP SCHEMA [%s]`, inner))
	defer db.Exec(ctx, fmt.Sprintf(`DROP TABLE [%s].[%s]`, inner, escaped))

	gtest.C(t, func(t *gtest.T) {
		for _, name := range []string{`"` + table + `"`, "[" + escaped + "]", `"dbo"."` + table + `"`, "dbo.[" + escaped + "]"} {
			fields, err := db.TableFields(ctx, name)
			t.AssertNil(err)
			t.Assert(dialectSortedFieldNames(fields), []string{"ID", "NAME"})
		}
	})

	gtest.C(t, func(t *gtest.T) {
		for _, name := range []string{
			`"` + schema + `"."` + table + `"`,
			"[" + inner + "].[" + escaped + "]",
			"test.[" + inner + "].[" + escaped + "]",
		} {
			fields, err := db.TableFields(ctx, name)
			t.AssertNil(err)
			t.Assert(name+" "+strings.Join(dialectSortedFieldNames(fields), ","), name+" ID,OTHER")
		}
	})
}

// Test_Mssql_HasField_UpperCase tests that HasField matches a field name, quoted or not, that is
// the stored name of a column or, in upper case, the stored name, as it does on Oracle.
func Test_Mssql_HasField_UpperCase(t *testing.T) {
	table := insertTableName("t_has_field")
	insertExec(fmt.Sprintf(
		`CREATE TABLE [%s] (Id int NOT NULL PRIMARY KEY, MixedCase varchar(45), lower_name varchar(45), UPPER_NAME varchar(45))`,
		table,
	))
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		for field, expect := range map[string]bool{
			"Id":           true,
			"ID":           false,
			"MixedCase":    true,
			`"MixedCase"`:  true,
			"mixedcase":    false,
			"MIXEDCASE":    false,
			"lower_name":   true,
			"LOWER_NAME":   false,
			`"Lower_Name"`: false,
			"UPPER_NAME":   true,
			"upper_name":   true,
			"Upper_Name":   true,
			`"upper_name"`: true,
			"mixed_case":   false,
			"nickname":     false,
		} {
			has, err := db.GetCore().HasField(ctx, table, field)
			t.AssertNil(err)
			t.Assert(fmt.Sprintf("%s=%v", field, has), fmt.Sprintf("%s=%v", field, expect))
		}
	})

	gtest.C(t, func(t *gtest.T) {
		has, err := db.Model(table + " AS a").HasField("upper_name")
		t.AssertNil(err)
		t.Assert(has, true)
		has, err = db.Model("dbo." + table).HasField("MixedCase")
		t.AssertNil(err)
		t.Assert(has, true)
		has, err = db.Model("dbo." + table).HasField("mixedcase")
		t.AssertNil(err)
		t.Assert(has, false)
	})
}

// Test_Mssql_TableFields_Synonym tests that the fields of a synonym are those of the table it
// stands for and are cached like those of a table, while a synonym of a missing table has none.
func Test_Mssql_TableFields_Synonym(t *testing.T) {
	var (
		table   = insertTableName("t_synonym")
		synonym = "syn_" + table
		dangled = "syn_missing_" + table
	)
	insertExec(
		fmt.Sprintf(`CREATE TABLE [%s] (ID int NOT NULL PRIMARY KEY, NAME varchar(45))`, table),
		fmt.Sprintf(`CREATE SYNONYM [%s] FOR [dbo].[%s]`, synonym, table),
		fmt.Sprintf(`CREATE SYNONYM [%s] FOR [dbo].[%s_missing]`, dangled, table),
	)
	defer dropTable(table)
	defer db.Exec(ctx, fmt.Sprintf(`DROP SYNONYM [%s]`, synonym))
	defer db.Exec(ctx, fmt.Sprintf(`DROP SYNONYM [%s]`, dangled))

	gtest.C(t, func(t *gtest.T) {
		var counts []int
		for i := 0; i < 3; i++ {
			sqls, err := gdb.CatchSQL(ctx, func(ctx context.Context) error {
				_, err := db.Model(synonym).Ctx(ctx).Where("id", 1).All()
				return err
			})
			t.AssertNil(err)
			counts = append(counts, len(sqls))
		}
		t.Assert(counts, []int{2, 1, 1})
	})

	gtest.C(t, func(t *gtest.T) {
		fields, err := db.TableFields(ctx, synonym)
		t.AssertNil(err)
		t.Assert(dialectSortedFieldNames(fields), []string{"ID", "NAME"})
		t.Assert(fields["ID"].Key, "PRI")

		_, err = db.Model(synonym).Data(g.Map{"id": 1, "name": "n1"}).Insert()
		t.AssertNil(err)
		_, err = db.Model(synonym).Data(g.Map{"id": 1, "name": "n2"}).Save()
		t.AssertNil(err)
		value, err := db.Model(table).Where("id", 1).Value("NAME")
		t.AssertNil(err)
		t.Assert(value, "n2")
	})

	gtest.C(t, func(t *gtest.T) {
		fields, err := db.TableFields(ctx, dangled)
		t.AssertNE(err, nil)
		t.Assert(len(fields), 0)
	})
}

// Test_Mssql_Model_JoinThreePartName tests that the conditions of a joined model whose table is
// named with its database and schema are prefixed with the last part of the name, which SQL Server
// exposes the table by.
func Test_Mssql_Model_JoinThreePartName(t *testing.T) {
	var (
		table1 = insertTableName("t_join_left")
		table2 = insertTableName("t_join_right")
	)
	insertExec(
		fmt.Sprintf(`CREATE TABLE [%s] (ID int NOT NULL PRIMARY KEY, NAME varchar(45))`, table1),
		fmt.Sprintf(`CREATE TABLE [%s] (ID int NOT NULL PRIMARY KEY, NICKNAME varchar(45))`, table2),
		fmt.Sprintf(`INSERT INTO [%s] (ID, NAME) VALUES (1, 'a'), (2, 'b')`, table1),
		fmt.Sprintf(`INSERT INTO [%s] (ID, NICKNAME) VALUES (1, 'x'), (2, 'y')`, table2),
	)
	defer dropTable(table1)
	defer dropTable(table2)

	gtest.C(t, func(t *gtest.T) {
		model := func(ctx context.Context) *gdb.Model {
			return db.Model(fmt.Sprintf("%s.dbo.%s", TestSchema, table1)).Ctx(ctx).
				LeftJoin(fmt.Sprintf("%s.dbo.%s", TestSchema, table2), fmt.Sprintf("%s.ID=%s.ID", table1, table2)).
				Fields(table1+".NAME", table2+".NICKNAME").Where("id", 2)
		}
		sqlStr, err := gdb.ToSQL(ctx, func(ctx context.Context) error {
			_, err := model(ctx).All()
			return err
		})
		t.AssertNil(err)
		t.Assert(insertUnquote(sqlStr), fmt.Sprintf(
			"SELECT %s.NAME,%s.NICKNAME FROM %s.dbo.%s LEFT JOIN %s.dbo.%s ON (%s.ID=%s.ID) WHERE %s.id=2",
			table1, table2, TestSchema, table1, TestSchema, table2, table1, table2, table1,
		))

		all, err := model(ctx).All()
		t.AssertNil(err)
		t.Assert(all, g.List{{"NAME": "b", "NICKNAME": "y"}})
	})
}

// Test_Mssql_Model_AliasedWrite tests writing through a model whose table has an alias, which SQL
// Server rejects on the target of INSERT, MERGE, UPDATE and DELETE: the values are converted by
// the types of the columns of the table, keeping the time zone of a time written to a
// datetimeoffset column, the rows are written to the table, and an insert returns the value of
// its IDENTITY column.
func Test_Mssql_Model_AliasedWrite(t *testing.T) {
	table := insertTableName("t_alias")
	insertExec(fmt.Sprintf(
		`CREATE TABLE [%s] (ID int NOT NULL PRIMARY KEY, NICKNAME varchar(45), INFO nvarchar(200), DT datetime2(0), DTO datetimeoffset(0))`, table,
	))
	defer dropTable(table)

	var (
		models = []func() *gdb.Model{
			func() *gdb.Model { return db.Model(table).As("a") },
			func() *gdb.Model { return db.Model(table, "a") },
			func() *gdb.Model { return db.Model(table + " AS a") },
			func() *gdb.Model { return db.Model(table + " a") },
			func() *gdb.Model { return db.Model("dbo." + table + " a") },
			func() *gdb.Model { return db.Model(`"dbo"."` + table + `" AS "a"`) },
			func() *gdb.Model { return db.Model("[" + table + "] a") },
			func() *gdb.Model { return db.Model("[dbo].[" + table + "] a") },
			func() *gdb.Model { return db.Model("[dbo].[" + table + "] AS a") },
			func() *gdb.Model { return db.Model("[" + table + "]").As("a") },
		}
		dt  = time.Date(2024, 1, 15, 4, 30, 45, 0, time.UTC)
		dto = time.Date(2024, 1, 15, 4, 30, 45, 0, time.FixedZone("UTC-5", -5*3600))
	)
	assertRow := func(t *gtest.T, id int, nickname string, info string, dt time.Time) {
		one, err := db.GetOne(ctx, fmt.Sprintf(
			"SELECT NICKNAME, INFO, CONVERT(varchar(19), DT, 120) AS DT FROM [%s] WHERE ID = ?", table,
		), id)
		t.AssertNil(err)
		t.Assert(one["NICKNAME"], nickname)
		t.Assert(one["INFO"], info)
		t.Assert(one["DT"], dt.In(time.Local).Format("2006-01-02 15:04:05"))
	}
	assertOffset := func(t *gtest.T, id int) {
		one, err := db.GetOne(ctx, fmt.Sprintf(
			"SELECT DTO, DATEPART(TZOFFSET, DTO) AS TZ FROM [%s] WHERE ID = ?", table,
		), id)
		t.AssertNil(err)
		t.Assert(one["TZ"], -300)
		t.Assert(one["DTO"].Time().Equal(dto), true)
	}

	for i, model := range models {
		id := (i + 1) * 10

		gtest.C(t, func(t *gtest.T) {
			_, err := model().Data(g.Map{
				"id": id, "nickname": "insert", "info": g.Map{"k": "insert"}, "dt": dt, "dto": dto,
			}).Insert()
			t.AssertNil(err)
			assertRow(t, id, "insert", `{"k":"insert"}`, dt)
			assertOffset(t, id)
		})

		gtest.C(t, func(t *gtest.T) {
			_, err := model().Data(g.Map{
				"nickname": "update", "info": g.Map{"k": "update"}, "dt": dt.Add(time.Hour),
				"dto": dto,
			}).Where("id", id).Update()
			t.AssertNil(err)
			assertRow(t, id, "update", `{"k":"update"}`, dt.Add(time.Hour))
			assertOffset(t, id)

			_, err = model().Data(g.Map{"nickname": "update_a"}).Where("a.id", id).Update()
			t.AssertNil(err)
			assertRow(t, id, "update_a", `{"k":"update"}`, dt.Add(time.Hour))
		})

		gtest.C(t, func(t *gtest.T) {
			_, err := model().Data(g.Map{
				"id": id, "nickname": "save", "info": g.Map{"k": "save"}, "dt": dt, "dto": dto,
			}).Save()
			t.AssertNil(err)
			assertRow(t, id, "save", `{"k":"save"}`, dt)
			assertOffset(t, id)

			_, err = model().Data(g.Map{"id": id + 1, "nickname": "save", "info": g.Map{"k": "new"}, "dt": dt}).Save()
			t.AssertNil(err)
			assertRow(t, id+1, "save", `{"k":"new"}`, dt)
		})

		gtest.C(t, func(t *gtest.T) {
			_, err := model().OnConflict("id").Data(g.Map{
				"id": id, "nickname": "conflict", "info": g.Map{"k": "conflict"}, "dt": dt.Add(2 * time.Hour),
			}).Save()
			t.AssertNil(err)
			assertRow(t, id, "conflict", `{"k":"conflict"}`, dt.Add(2*time.Hour))

			_, err = model().OnDuplicate("nickname").Data(g.Map{
				"id": id, "nickname": "duplicate", "info": g.Map{"k": "duplicate"}, "dt": dt,
			}).Save()
			t.AssertNil(err)
			assertRow(t, id, "duplicate", `{"k":"conflict"}`, dt.Add(2*time.Hour))
		})

		gtest.C(t, func(t *gtest.T) {
			_, err := model().Data(g.Map{"id": id, "nickname": "replace", "info": g.Map{"k": "replace"}, "dt": dt}).Replace()
			t.AssertNil(err)
			assertRow(t, id, "replace", `{"k":"replace"}`, dt)
		})

		gtest.C(t, func(t *gtest.T) {
			result, err := model().Data(g.Map{"id": id, "nickname": "ignored", "info": g.Map{"k": "ignored"}, "dt": dt}).InsertIgnore()
			t.AssertNil(err)
			n, err := result.RowsAffected()
			t.AssertNil(err)
			t.Assert(n, 0)
			assertRow(t, id, "replace", `{"k":"replace"}`, dt)

			result, err = model().Data(g.Map{"id": id + 2, "nickname": "ignore", "info": g.Map{"k": "ignore"}, "dt": dt}).InsertIgnore()
			t.AssertNil(err)
			n, err = result.RowsAffected()
			t.AssertNil(err)
			t.Assert(n, 1)
			assertRow(t, id+2, "ignore", `{"k":"ignore"}`, dt)
		})

		gtest.C(t, func(t *gtest.T) {
			result, err := model().Where("a.id", id+2).Delete()
			t.AssertNil(err)
			n, err := result.RowsAffected()
			t.AssertNil(err)
			t.Assert(n, 1)
			count, err := model().Where("a.id BETWEEN ? AND ?", id, id+2).Count()
			t.AssertNil(err)
			t.Assert(count, 2)
		})
	}

	gtest.C(t, func(t *gtest.T) {
		count, err := db.Model(table).Count()
		t.AssertNil(err)
		t.Assert(count, len(models)*2)
	})

	gtest.C(t, func(t *gtest.T) {
		identityTable := insertCreateIdentityTable()
		defer dropTable(identityTable)

		for i, model := range []*gdb.Model{
			db.Model(identityTable).As("a"),
			db.Model(identityTable + " a"),
			db.Model(`"dbo"."` + identityTable + `" AS "a"`),
			db.Model("[dbo].[" + identityTable + "] AS a"),
		} {
			id, err := model.Data(g.Map{"name": "identity"}).InsertAndGetId()
			t.AssertNil(err)
			t.Assert(id, i+1)
		}
	})
}

// Test_Mssql_Time_ZonelessWrite tests that a time in any time zone written to a date and time column
// without time zone is stored as its wall clock in the local time zone, which the driver reads it
// in, so that it comes back as the same instant, whether it is written through the data of a model
// or as an argument, as a time, a gtime or a sql.NullTime; that a date or a time column stores the
// date or time of day of a time of the data of a model, and of an argument in the local time zone;
// that an invalid sql.NullTime is written as NULL; and that a column with time zone keeps the
// instant, and the time zone of a time of the data of a model, while an argument comes in the local
// time zone.
func Test_Mssql_Time_ZonelessWrite(t *testing.T) {
	table := insertTableName("t_zoneless")
	insertExec(fmt.Sprintf(
		`CREATE TABLE [%s] (ID int NOT NULL PRIMARY KEY, D datetime2(0), DT datetime, SDT smalldatetime, DA date, TI time(0), DTO datetimeoffset(0))`,
		table,
	))
	defer dropTable(table)

	var (
		local  = time.Date(2024, 1, 2, 3, 4, 0, 0, time.Local)
		values = []time.Time{
			local,
			local.UTC(),
			local.In(time.FixedZone("UTC-5", -5*3600)),
			local.In(time.FixedZone("UTC+14", 14*3600)),
		}
		columns = []string{"D", "DT", "SDT", "DTO"}
	)
	assertOffset := func(t *gtest.T, id int, expect time.Time) {
		value, err := db.Model(table).Where("id", id).Value("DTO")
		t.AssertNil(err)
		_, offset := value.Time().Zone()
		_, expectOffset := expect.Zone()
		t.Assert(offset, expectOffset)
	}
	assertRow := func(t *gtest.T, id int, day time.Time) {
		one, err := db.Model(table).Where("id", id).One()
		t.AssertNil(err)
		for _, column := range columns {
			t.Assert(fmt.Sprintf("%d %s %v", id, column, one[column].Time().Equal(local)), fmt.Sprintf("%d %s true", id, column))
		}
		t.Assert(one["DA"].Time().Format("2006-01-02"), day.Format("2006-01-02"))
		raw, err := db.GetOne(ctx, fmt.Sprintf(
			"SELECT CONVERT(varchar(19), D, 120) AS D, CONVERT(varchar(8), TI, 108) AS TI FROM [%s] WHERE ID = ?", table,
		), id)
		t.AssertNil(err)
		t.Assert(raw["D"], local.Format("2006-01-02 15:04:05"))
		t.Assert(raw["TI"], day.Format("15:04:05"))
	}

	for i, value := range values {
		gtest.C(t, func(t *gtest.T) {
			id := i*10 + 1
			_, err := db.Model(table).Data(g.Map{
				"id": id, "d": value, "dt": value, "sdt": value, "da": value, "ti": value, "dto": value,
			}).Insert()
			t.AssertNil(err)
			assertRow(t, id, value)
			assertOffset(t, id, value)

			id++
			_, err = db.Model(table).Data(g.Map{
				"id": id, "d": gtime.NewFromTime(value), "dt": &value, "sdt": *gtime.NewFromTime(value), "da": &value,
				"ti": gtime.NewFromTime(value), "dto": value,
			}).Insert()
			t.AssertNil(err)
			assertRow(t, id, value)

			id++
			_, err = db.Exec(ctx, fmt.Sprintf(
				"INSERT INTO [%s] (ID, D, DT, SDT, DA, TI, DTO) VALUES (?, ?, ?, ?, ?, ?, ?)", table,
			), id, value, value, &value, gtime.NewFromTime(value), value, value)
			t.AssertNil(err)
			assertRow(t, id, local)
			assertOffset(t, id, local)

			id++
			_, err = db.Model(table).Data(g.Map{"id": id}).Insert()
			t.AssertNil(err)
			_, err = db.Model(table).Data(
				"D = ?, DT = ?, SDT = ?, DA = ?, TI = ?, DTO = ?", value, value, value, value, value, value,
			).Where("id", id).Update()
			t.AssertNil(err)
			assertRow(t, id, local)

			id++
			_, err = db.Model(table).Data(g.Map{"id": id}).Insert()
			t.AssertNil(err)
			_, err = db.Model(table).Data(g.Map{
				"d": value, "dt": value, "sdt": value, "da": value, "ti": value, "dto": value,
			}).Where("id", id).Update()
			t.AssertNil(err)
			assertRow(t, id, value)

			id++
			nullTime := sql.NullTime{Time: value, Valid: true}
			_, err = db.Exec(ctx, fmt.Sprintf(
				"INSERT INTO [%s] (ID, D, DT, SDT, DA, TI, DTO) VALUES (?, ?, ?, ?, ?, ?, ?)", table,
			), id, nullTime, &nullTime, nullTime, &nullTime, nullTime, nullTime)
			t.AssertNil(err)
			assertRow(t, id, local)

			id++
			_, err = db.Model(table).Data(g.Map{"id": id}).Insert()
			t.AssertNil(err)
			_, err = db.Model(table).Data(
				"D = ?, DT = ?, SDT = ?, DA = ?, TI = ?, DTO = ?",
				nullTime, &nullTime, nullTime, &nullTime, nullTime, nullTime,
			).Where("id", id).Update()
			t.AssertNil(err)
			assertRow(t, id, local)
		})
	}

	gtest.C(t, func(t *gtest.T) {
		for _, value := range values {
			count, err := db.Model(table).Where("DTO", value).Count()
			t.AssertNil(err)
			t.Assert(count, len(values)*7)
		}
	})

	gtest.C(t, func(t *gtest.T) {
		_, err := db.Exec(ctx, fmt.Sprintf(
			"INSERT INTO [%s] (ID, D, DA, DTO) VALUES (?, ?, ?, ?)", table,
		), 1000, sql.NullTime{}, &sql.NullTime{}, sql.NullTime{})
		t.AssertNil(err)
		one, err := db.Model(table).Where("id", 1000).One()
		t.AssertNil(err)
		t.Assert(one["D"].IsNil(), true)
		t.Assert(one["DA"].IsNil(), true)
		t.Assert(one["DTO"].IsNil(), true)
	})
}

// Test_Mssql_Time_DateAndTimeOfDay tests that a time written through the data of a model to a date
// or a time column keeps its own date or time of day, whatever the time zone of the driver, as on
// MySQL.
func Test_Mssql_Time_DateAndTimeOfDay(t *testing.T) {
	table := insertTableName("t_date_time")
	insertExec(fmt.Sprintf(`CREATE TABLE [%s] (ID int NOT NULL PRIMARY KEY, DA date, TI time(0))`, table))
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		node := *db.GetConfig()
		node.Timezone = "America/New_York"
		zoneDb, err := gdb.New(node)
		t.AssertNil(err)
		defer zoneDb.Close(ctx)

		date, err := time.Parse("2006-01-02", "2024-01-02")
		t.AssertNil(err)
		timeOfDay, err := time.Parse("15:04:05", "03:04:05")
		t.AssertNil(err)
		_, err = zoneDb.Model(table).Data(g.Map{"id": 1, "da": date, "ti": timeOfDay}).Insert()
		t.AssertNil(err)
		_, err = zoneDb.Model(table).Data(g.Map{"id": 2}).Insert()
		t.AssertNil(err)
		_, err = zoneDb.Model(table).Data(g.Map{"da": date, "ti": timeOfDay}).Where("id", 2).Update()
		t.AssertNil(err)

		for _, id := range []int{1, 2} {
			raw, err := db.GetOne(ctx, fmt.Sprintf(
				"SELECT CONVERT(varchar(10), DA, 120) AS DA, CONVERT(varchar(8), TI, 108) AS TI FROM [%s] WHERE ID = ?", table,
			), id)
			t.AssertNil(err)
			t.Assert(raw, g.Map{"DA": "2024-01-02", "TI": "03:04:05"})
		}
	})
}

// Test_Mssql_Merge_BatchAtomic tests that Save, Replace and InsertIgnore of several records write
// none of them if one fails, as the MERGE statements of the records run in one transaction, or
// back to a savepoint of the transaction of the caller, which keeps its other writes and commits.
func Test_Mssql_Merge_BatchAtomic(t *testing.T) {
	table := insertTableName("t_merge_atomic")
	insertExec(fmt.Sprintf(
		`CREATE TABLE [%s] (ID int NOT NULL PRIMARY KEY, N int NOT NULL CHECK (N < 100))`, table,
	))
	defer dropTable(table)

	var (
		failing = g.List{{"id": 1, "n": 1}, {"id": 2, "n": 2}, {"id": 3, "n": 300}}
		valid   = g.List{{"id": 1, "n": 1}, {"id": 2, "n": 2}, {"id": 3, "n": 3}}
		writes  = map[string]func(m *gdb.Model) error{
			"Save": func(m *gdb.Model) error {
				_, err := m.Save()
				return err
			},
			"Replace": func(m *gdb.Model) error {
				_, err := m.Replace()
				return err
			},
			"InsertIgnore": func(m *gdb.Model) error {
				_, err := m.InsertIgnore()
				return err
			},
		}
	)
	for name, write := range writes {
		gtest.C(t, func(t *gtest.T) {
			err := write(db.Model(table).Data(failing))
			t.AssertNE(err, nil)
			count, err := db.Model(table).Count()
			t.AssertNil(err)
			t.Assert(fmt.Sprintf("%s %d", name, count), name+" 0")
		})

		gtest.C(t, func(t *gtest.T) {
			err := db.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
				if err := write(db.Model(table).Ctx(ctx).Data(valid)); err != nil {
					return err
				}
				return fmt.Errorf("rollback %s", name)
			})
			t.AssertNE(err, nil)
			count, err := db.Model(table).Count()
			t.AssertNil(err)
			t.Assert(fmt.Sprintf("%s %d", name, count), name+" 0")
		})

		for entry, model := range map[string]func(ctx context.Context, tx gdb.TX) *gdb.Model{
			"tx.Model": func(ctx context.Context, tx gdb.TX) *gdb.Model { return tx.Model(table) },
			"Ctx(ctx)": func(ctx context.Context, tx gdb.TX) *gdb.Model { return db.Model(table).Ctx(ctx) },
		} {
			gtest.C(t, func(t *gtest.T) {
				defer db.Model(table).Where("1=1").Delete()
				label := name + " " + entry
				err := db.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
					if _, err := tx.Insert(table, g.Map{"id": 10, "n": 10}); err != nil {
						return err
					}
					t.AssertNE(write(model(ctx, tx).Data(failing)), nil)
					count, err := tx.Model(table).Where("id < ?", 10).Count()
					t.AssertNil(err)
					t.Assert(fmt.Sprintf("%s %d", label, count), label+" 0")
					if err = write(model(ctx, tx).Data(valid)); err != nil {
						return err
					}
					count, err = tx.Model(table).Count()
					t.AssertNil(err)
					t.Assert(fmt.Sprintf("%s %d", label, count), label+" 4")
					return nil
				})
				t.AssertNil(err)
				count, err := db.Model(table).Count()
				t.AssertNil(err)
				t.Assert(fmt.Sprintf("%s %d", label, count), label+" 4")
			})
		}

		gtest.C(t, func(t *gtest.T) {
			defer db.Model(table).Where("1=1").Delete()
			t.AssertNil(write(db.Model(table).Data(valid)))
			count, err := db.Model(table).Count()
			t.AssertNil(err)
			t.Assert(fmt.Sprintf("%s %d", name, count), name+" 3")
		})
	}
}

// Test_Mssql_Merge_Concurrent tests that Save and InsertIgnore of the same new record by several
// goroutines at once write it once without error, as the MERGE of the driver holds the lock on the
// key it looks up until its transaction ends.
func Test_Mssql_Merge_Concurrent(t *testing.T) {
	table := insertTableName("t_merge_concurrent")
	insertExec(fmt.Sprintf(`CREATE TABLE [%s] (ID int NOT NULL PRIMARY KEY, NAME varchar(45))`, table))
	defer dropTable(table)

	const (
		concurrency = 10
		rounds      = 50
	)
	for _, c := range []struct {
		name  string
		write func(m *gdb.Model, data ...any) (sql.Result, error)
	}{
		{"Save", (*gdb.Model).Save},
		{"InsertIgnore", (*gdb.Model).InsertIgnore},
	} {
		gtest.C(t, func(t *gtest.T) {
			defer db.Model(table).Where("1=1").Delete()
			for id := 1; id <= rounds; id++ {
				var (
					wg      sync.WaitGroup
					barrier = make(chan struct{})
					errs    = make([]error, concurrency)
				)
				wg.Add(concurrency)
				for i := 0; i < concurrency; i++ {
					go func() {
						defer wg.Done()
						<-barrier
						_, errs[i] = c.write(db.Model(table), g.Map{"id": id, "name": fmt.Sprintf("name_%d", i)})
					}()
				}
				close(barrier)
				wg.Wait()
				for _, err := range errs {
					t.AssertNil(err)
				}
				count, err := db.Model(table).Where("id", id).Count()
				t.AssertNil(err)
				t.Assert(fmt.Sprintf("%s %d: %d", c.name, id, count), fmt.Sprintf("%s %d: 1", c.name, id))
			}
		})
	}
}
