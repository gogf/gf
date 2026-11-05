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
	"testing"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gconv"
)

// insertTableName returns a table name starting with `prefix` that is unique within the test run.
func insertTableName(prefix string) string {
	return fmt.Sprintf("%s_%d", prefix, gtime.TimestampNano())
}

// insertExec executes `statements` in order and stops the test run if one of them fails.
func insertExec(statements ...string) {
	for _, statement := range statements {
		if _, err := db.Exec(ctx, statement); err != nil {
			gtest.Fatal(err)
		}
	}
}

// insertCreateIdentityTable creates a table whose primary key ID is an IDENTITY column and returns
// its name.
func insertCreateIdentityTable() string {
	table := insertTableName("insert_identity")
	insertExec(fmt.Sprintf(
		`CREATE TABLE [%s] (ID int IDENTITY(1,1) NOT NULL, NAME varchar(45) NULL, PRIMARY KEY (ID))`, table,
	))
	return table
}

// insertCreatePlainTable creates a table whose primary key ID is filled by the caller and returns
// its name.
func insertCreatePlainTable() string {
	table := insertTableName("insert_plain")
	insertExec(fmt.Sprintf(
		`CREATE TABLE [%s] (ID int NOT NULL, NAME varchar(45) NULL, PRIMARY KEY (ID))`, table,
	))
	return table
}

// insertUnquote returns `sql` without double quotes and surrounding spaces, so that an assertion on
// it holds whether or not the driver keeps the quotes gf puts around identifiers.
func insertUnquote(sql string) string {
	return gstr.Trim(gstr.Replace(sql, `"`, ""))
}

// insertAssertResult asserts that `result` reports `affected` rows and 0 as the last insert id,
// both without error.
func insertAssertResult(t *gtest.T, result sql.Result, affected int64) {
	n, err := result.RowsAffected()
	t.AssertNil(err)
	t.Assert(n, affected)
	id, err := result.LastInsertId()
	t.AssertNil(err)
	t.Assert(id, 0)
}

// Test_MSSQL_Insert_LastInsertIdIdentity tests the last insert id and the affected rows of inserts
// into a table whose primary key is an IDENTITY column, for single and batch inserts.
func Test_MSSQL_Insert_LastInsertIdIdentity(t *testing.T) {
	table := insertCreateIdentityTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		id, err := db.Model(table).Data(g.Map{"name": "a"}).InsertAndGetId()
		t.AssertNil(err)
		t.Assert(id, 1)

		result, err := db.Model(table).Data(g.Map{"name": "b"}).Insert()
		t.AssertNil(err)
		n, err := result.RowsAffected()
		t.AssertNil(err)
		t.Assert(n, 1)
		id, err = result.LastInsertId()
		t.AssertNil(err)
		t.Assert(id, 2)
	})

	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Data(g.List{{"name": "c"}, {"name": "d"}, {"name": "e"}}).Insert()
		t.AssertNil(err)
		n, err := result.RowsAffected()
		t.AssertNil(err)
		t.Assert(n, 3)
		id, err := result.LastInsertId()
		t.AssertNil(err)
		t.Assert(id, 3)
	})

	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Data(g.List{{"name": "f"}, {"name": "g"}, {"name": "h"}}).Batch(2).Insert()
		t.AssertNil(err)
		n, err := result.RowsAffected()
		t.AssertNil(err)
		t.Assert(n, 3)
		id, err := result.LastInsertId()
		t.AssertNil(err)
		t.Assert(id, 8)

		array, err := db.Model(table).OrderAsc("id").Array("name")
		t.AssertNil(err)
		t.Assert(gconv.Strings(array), []string{"a", "b", "c", "d", "e", "f", "g", "h"})
	})
}

// Test_MSSQL_Insert_LastInsertIdNoIdentity tests that inserts into a table without IDENTITY column
// report 0 as the last insert id without error.
func Test_MSSQL_Insert_LastInsertIdNoIdentity(t *testing.T) {
	table := insertCreatePlainTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		id, err := db.Model(table).Data(g.Map{"id": 1, "name": "a"}).InsertAndGetId()
		t.AssertNil(err)
		t.Assert(id, 0)

		result, err := db.Model(table).Data(g.List{{"id": 2, "name": "b"}, {"id": 3, "name": "c"}}).Insert()
		t.AssertNil(err)
		insertAssertResult(t, result, 2)

		count, err := db.Model(table).Count()
		t.AssertNil(err)
		t.Assert(count, 3)
	})
}

// Test_MSSQL_Insert_LastInsertIdIdentityNotPrimaryKey tests that the value of an IDENTITY column
// that is not the primary key is reported as the last insert id.
func Test_MSSQL_Insert_LastInsertIdIdentityNotPrimaryKey(t *testing.T) {
	table := insertTableName("insert_seq")
	insertExec(fmt.Sprintf(
		`CREATE TABLE [%s] (CODE varchar(20) NOT NULL, SEQ int IDENTITY(100,5) NOT NULL, PRIMARY KEY (CODE))`, table,
	))
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		id, err := db.Model(table).Data(g.Map{"code": "a"}).InsertAndGetId()
		t.AssertNil(err)
		t.Assert(id, 100)

		id, err = db.Model(table).Data(g.Map{"code": "b"}).InsertAndGetId()
		t.AssertNil(err)
		t.Assert(id, 105)
	})
}

// Test_MSSQL_Insert_LastInsertIdTransaction tests InsertAndGetId through a transaction object and
// through a context carrying the transaction, and that the inserts are rolled back with it.
func Test_MSSQL_Insert_LastInsertIdTransaction(t *testing.T) {
	table := insertCreateIdentityTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		var ids []int64
		err := db.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
			id, err := tx.Model(table).Data(g.Map{"name": "tx"}).InsertAndGetId()
			if err != nil {
				return err
			}
			ids = append(ids, id)
			id, err = db.Model(table).Ctx(ctx).Data(g.Map{"name": "ctx"}).InsertAndGetId()
			if err != nil {
				return err
			}
			ids = append(ids, id)
			return gerror.New("rollback")
		})
		t.AssertNE(err, nil)
		t.Assert(ids, []int64{1, 2})

		count, err := db.Model(table).Count()
		t.AssertNil(err)
		t.Assert(count, 0)
	})

	gtest.C(t, func(t *gtest.T) {
		var id int64
		err := db.Transaction(ctx, func(ctx context.Context, tx gdb.TX) (err error) {
			id, err = db.Model(table).Ctx(ctx).Data(g.Map{"name": "commit"}).InsertAndGetId()
			return err
		})
		t.AssertNil(err)
		t.Assert(id, 3)

		value, err := db.Model(table).Where("id", id).Value("name")
		t.AssertNil(err)
		t.Assert(value, "commit")
	})
}

// Test_MSSQL_Insert_ToSQL tests that ToSQL returns the statements of Insert and Save without
// executing them.
func Test_MSSQL_Insert_ToSQL(t *testing.T) {
	var (
		identityTable = insertCreateIdentityTable()
		plainTable    = insertCreatePlainTable()
	)
	defer dropTable(identityTable)
	defer dropTable(plainTable)

	gtest.C(t, func(t *gtest.T) {
		var id int64
		sqlStr, err := gdb.ToSQL(ctx, func(ctx context.Context) (err error) {
			id, err = db.Model(identityTable).Ctx(ctx).Data(g.Map{"name": "a"}).InsertAndGetId()
			return err
		})
		t.AssertNil(err)
		t.Assert(id, 0)
		t.Assert(insertUnquote(sqlStr), fmt.Sprintf(
			"INSERT INTO %s(NAME) OUTPUT INSERTED.ID VALUES('a')", identityTable,
		))

		count, err := db.Model(identityTable).Count()
		t.AssertNil(err)
		t.Assert(count, 0)
	})

	gtest.C(t, func(t *gtest.T) {
		sqlStr, err := gdb.ToSQL(ctx, func(ctx context.Context) error {
			_, err := db.Model(plainTable).Ctx(ctx).Data(g.Map{"name": "a", "id": 1}).Insert()
			return err
		})
		t.AssertNil(err)
		t.Assert(insertUnquote(sqlStr), fmt.Sprintf("INSERT INTO %s(ID,NAME) VALUES(1,'a')", plainTable))

		sqlStr, err = gdb.ToSQL(ctx, func(ctx context.Context) error {
			_, err := db.Model(plainTable).Ctx(ctx).Data(g.Map{"name": "b", "id": 1}).Save()
			return err
		})
		t.AssertNil(err)
		t.Assert(insertUnquote(sqlStr), fmt.Sprintf(
			"MERGE INTO %s WITH (HOLDLOCK) AS T1 USING (VALUES(1,'b')) T2 (ID,NAME) ON (T1.ID = T2.ID) "+
				"WHEN NOT MATCHED THEN INSERT(ID,NAME) VALUES (T2.ID,T2.NAME) "+
				"WHEN MATCHED THEN UPDATE SET T1.NAME = T2.NAME;",
			plainTable,
		))

		count, err := db.Model(plainTable).Count()
		t.AssertNil(err)
		t.Assert(count, 0)
	})
}

// Test_MSSQL_Insert_CatchSQL tests that CatchSQL records the statements of Insert, which are
// executed.
func Test_MSSQL_Insert_CatchSQL(t *testing.T) {
	table := insertCreateIdentityTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		var id int64
		sqlArray, err := gdb.CatchSQL(ctx, func(ctx context.Context) (err error) {
			id, err = db.Model(table).Ctx(ctx).Data(g.Map{"name": "a"}).InsertAndGetId()
			return err
		})
		t.AssertNil(err)
		t.Assert(id, 1)
		t.AssertGT(len(sqlArray), 0)
		t.Assert(insertUnquote(sqlArray[len(sqlArray)-1]), fmt.Sprintf(
			"INSERT INTO %s(NAME) OUTPUT INSERTED.ID VALUES('a')", table,
		))

		value, err := db.Model(table).Where("id", 1).Value("name")
		t.AssertNil(err)
		t.Assert(value, "a")
	})

	gtest.C(t, func(t *gtest.T) {
		sqlArray, err := gdb.CatchSQL(ctx, func(ctx context.Context) error {
			_, err := db.Model(table).Ctx(ctx).Data(g.List{{"name": "b"}, {"name": "c"}}).Batch(1).Insert()
			return err
		})
		t.AssertNil(err)
		t.AssertGE(len(sqlArray), 2)
		t.Assert(insertUnquote(sqlArray[len(sqlArray)-2]), fmt.Sprintf(
			"INSERT INTO %s(NAME) OUTPUT INSERTED.ID VALUES('b')", table,
		))
		t.Assert(insertUnquote(sqlArray[len(sqlArray)-1]), fmt.Sprintf(
			"INSERT INTO %s(NAME) OUTPUT INSERTED.ID VALUES('c')", table,
		))

		count, err := db.Model(table).Count()
		t.AssertNil(err)
		t.Assert(count, 3)
	})
}

// Test_MSSQL_Insert_RawExec tests hand-written INSERT statements through Exec, with and without
// column list, which are executed as they are written.
func Test_MSSQL_Insert_RawExec(t *testing.T) {
	var (
		identityTable = insertCreateIdentityTable()
		plainTable    = insertCreatePlainTable()
	)
	defer dropTable(identityTable)
	defer dropTable(plainTable)

	gtest.C(t, func(t *gtest.T) {
		result, err := db.Exec(ctx, fmt.Sprintf("INSERT INTO %s VALUES(?, ?)", plainTable), 1, "a")
		t.AssertNil(err)
		n, err := result.RowsAffected()
		t.AssertNil(err)
		t.Assert(n, 1)

		result, err = db.Exec(ctx, fmt.Sprintf("insert into %s values (?, ?), (?, ?)", plainTable), 2, "b", 3, "c")
		t.AssertNil(err)
		n, err = result.RowsAffected()
		t.AssertNil(err)
		t.Assert(n, 2)

		array, err := db.Model(plainTable).OrderAsc("id").Array("name")
		t.AssertNil(err)
		t.Assert(gconv.Strings(array), []string{"a", "b", "c"})
	})

	gtest.C(t, func(t *gtest.T) {
		result, err := db.Exec(ctx, fmt.Sprintf("INSERT INTO %s VALUES(?)", identityTable), "a")
		t.AssertNil(err)
		n, err := result.RowsAffected()
		t.AssertNil(err)
		t.Assert(n, 1)

		sqlArray, err := gdb.CatchSQL(ctx, func(ctx context.Context) error {
			_, err := db.Exec(ctx, fmt.Sprintf("INSERT INTO %s(NAME) VALUES(?)", identityTable), "b")
			return err
		})
		t.AssertNil(err)
		t.Assert(sqlArray, []string{fmt.Sprintf("INSERT INTO %s(NAME) VALUES('b')", identityTable)})

		array, err := db.Model(identityTable).OrderAsc("id").Array("name")
		t.AssertNil(err)
		t.Assert(gconv.Strings(array), []string{"a", "b"})
	})
}

// Test_MSSQL_Merge_LastInsertId tests that Save, Replace and InsertIgnore report the affected rows,
// and 0 as the last insert id without error.
func Test_MSSQL_Merge_LastInsertId(t *testing.T) {
	table := insertCreatePlainTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Data(g.Map{"id": 1, "name": "save"}).Save()
		t.AssertNil(err)
		insertAssertResult(t, result, 1)

		result, err = db.Model(table).Data(g.Map{"id": 1, "name": "replace"}).Replace()
		t.AssertNil(err)
		insertAssertResult(t, result, 1)

		result, err = db.Model(table).Data(g.Map{"id": 2, "name": "ignore"}).InsertIgnore()
		t.AssertNil(err)
		insertAssertResult(t, result, 1)

		result, err = db.Model(table).Data(g.Map{"id": 2, "name": "ignored"}).InsertIgnore()
		t.AssertNil(err)
		insertAssertResult(t, result, 0)

		array, err := db.Model(table).OrderAsc("id").Array("name")
		t.AssertNil(err)
		t.Assert(gconv.Strings(array), []string{"replace", "ignore"})
	})
}

// Test_MSSQL_Merge_Batch tests that Save, Replace and InsertIgnore write every record of a list and
// sum the affected rows.
func Test_MSSQL_Merge_Batch(t *testing.T) {
	table := insertCreatePlainTable()
	defer dropTable(table)
	insertExec(fmt.Sprintf(`INSERT INTO [%s] (ID, NAME) VALUES (1, 'a'), (2, 'b')`, table))

	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Data(g.List{
			{"id": 2, "name": "ignored"},
			{"id": 3, "name": "c"},
			{"id": 4, "name": "d"},
		}).InsertIgnore()
		t.AssertNil(err)
		insertAssertResult(t, result, 2)

		result, err = db.Model(table).Data(g.List{
			{"id": 1, "name": "a1"},
			{"id": 5, "name": "e"},
		}).Batch(1).Save()
		t.AssertNil(err)
		insertAssertResult(t, result, 2)

		result, err = db.Model(table).Data(g.List{
			{"id": 2, "name": "b2"},
			{"id": 6, "name": "f"},
		}).Replace()
		t.AssertNil(err)
		insertAssertResult(t, result, 2)

		array, err := db.Model(table).OrderAsc("id").Array("name")
		t.AssertNil(err)
		t.Assert(gconv.Strings(array), []string{"a1", "b2", "c", "d", "e", "f"})
	})
}

// Test_MSSQL_Merge_RawValue tests that gdb.Raw values of Save and InsertIgnore are written as SQL.
func Test_MSSQL_Merge_RawValue(t *testing.T) {
	table := insertCreatePlainTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(table).Data(g.Map{"id": 1, "name": gdb.Raw("UPPER('save')")}).Save()
		t.AssertNil(err)
		value, err := db.Model(table).Where("id", 1).Value("name")
		t.AssertNil(err)
		t.Assert(value, "SAVE")

		_, err = db.Model(table).Data(g.Map{"id": 1, "name": gdb.Raw("CONCAT('save', '_2')")}).Save()
		t.AssertNil(err)
		value, err = db.Model(table).Where("id", 1).Value("name")
		t.AssertNil(err)
		t.Assert(value, "save_2")

		_, err = db.Model(table).Data(g.Map{"id": gdb.Raw("1+1"), "name": "ignore"}).InsertIgnore()
		t.AssertNil(err)
		value, err = db.Model(table).Where("id", 2).Value("name")
		t.AssertNil(err)
		t.Assert(value, "ignore")
	})
}

// Test_MSSQL_Merge_OnConflictReplaceAndIgnore tests that Replace and InsertIgnore detect conflicts
// on the columns given by OnConflict, as the driver builds their MERGE from them.
func Test_MSSQL_Merge_OnConflictReplaceAndIgnore(t *testing.T) {
	table := insertTableName("insert_nopk")
	insertExec(
		fmt.Sprintf(`CREATE TABLE [%s] (ID int NULL, NAME varchar(45) NULL)`, table),
		fmt.Sprintf(`INSERT INTO [%s] (ID, NAME) VALUES (5, 'a')`, table),
	)
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(table).Data(g.Map{"id": 5, "name": "b"}).InsertIgnore()
		t.Assert(gerror.Code(err), gcode.CodeMissingParameter)

		_, err = db.Model(table).Data(g.Map{"id": 5, "name": "b"}).OnConflict("id").InsertIgnore()
		t.AssertNil(err)
		value, err := db.Model(table).Where("id", 5).Value("name")
		t.AssertNil(err)
		t.Assert(value, "a")

		_, err = db.Model(table).Data(g.Map{"id": 5, "name": "c"}).OnConflict("id").Replace()
		t.AssertNil(err)
		value, err = db.Model(table).Where("id", 5).Value("name")
		t.AssertNil(err)
		t.Assert(value, "c")

		count, err := db.Model(table).Count()
		t.AssertNil(err)
		t.Assert(count, 1)
	})
}

// Test_MSSQL_Merge_ConflictKeyCase tests that the conflict columns are left out of the MERGE update
// whatever their letter case, which matters for an IDENTITY primary key that cannot be updated.
func Test_MSSQL_Merge_ConflictKeyCase(t *testing.T) {
	table := insertTableName("insert_case")
	insertExec(
		fmt.Sprintf(`CREATE TABLE [%s] (id int IDENTITY(1,1) NOT NULL, name varchar(45) NULL, PRIMARY KEY (id))`, table),
		fmt.Sprintf(`INSERT INTO [%s] (name) VALUES ('a')`, table),
	)
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		err := db.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
			if _, err := tx.Exec(fmt.Sprintf("SET IDENTITY_INSERT [%s] ON", table)); err != nil {
				return err
			}
			if _, err := tx.Model(table).Data(g.Map{"id": 1, "name": "b"}).Save(); err != nil {
				return err
			}
			if _, err := tx.Model(table).Data(g.Map{"ID": 1, "NAME": "c"}).OnConflict("ID").Save(); err != nil {
				return err
			}
			_, err := tx.Exec(fmt.Sprintf("SET IDENTITY_INSERT [%s] OFF", table))
			return err
		})
		t.AssertNil(err)

		all, err := db.Model(table).All()
		t.AssertNil(err)
		t.Assert(len(all), 1)
		t.Assert(all[0]["id"], 1)
		t.Assert(all[0]["name"], "c")
	})
}

// Test_MSSQL_Merge_ReplaceCreatedField tests that Replace overwrites the soft created field of an
// existing row while Save keeps it.
func Test_MSSQL_Merge_ReplaceCreatedField(t *testing.T) {
	table := insertTableName("insert_soft")
	insertExec(fmt.Sprintf(
		`CREATE TABLE [%s] (ID int NOT NULL, NAME varchar(45) NULL, CREATED_AT datetime2(0) NULL, PRIMARY KEY (ID))`,
		table,
	))
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(table).Data(g.Map{
			"id": 1, "name": "a", "created_at": gtime.NewFromStr("2024-05-30 20:00:00"),
		}).Insert()
		t.AssertNil(err)

		_, err = db.Model(table).Data(g.Map{
			"id": 1, "name": "b", "created_at": gtime.NewFromStr("2024-05-30 21:00:00"),
		}).Save()
		t.AssertNil(err)
		one, err := db.Model(table).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["NAME"], "b")
		t.Assert(one["CREATED_AT"].String(), "2024-05-30 20:00:00")

		_, err = db.Model(table).Data(g.Map{
			"id": 1, "name": "c", "created_at": gtime.NewFromStr("2024-05-30 22:00:00"),
		}).Replace()
		t.AssertNil(err)
		one, err = db.Model(table).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["NAME"], "c")
		t.Assert(one["CREATED_AT"].String(), "2024-05-30 22:00:00")
	})
}

// Test_MSSQL_Merge_OnDuplicateCounter tests that a gdb.Counter in OnDuplicate adds to the value
// stored in the matched row, as MySQL does, not to the value being inserted.
func Test_MSSQL_Merge_OnDuplicateCounter(t *testing.T) {
	table := insertTableName("insert_counter")
	insertExec(
		fmt.Sprintf(`CREATE TABLE [%s] (ID int NOT NULL, NAME varchar(45) NULL, SCORE int NULL, PRIMARY KEY (ID))`, table),
		fmt.Sprintf(`INSERT INTO [%s] (ID, NAME, SCORE) VALUES (1, 'a', 10)`, table),
	)
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(table).Data(g.Map{"id": 1, "name": "b", "score": 0}).
			OnDuplicate(g.Map{"score": gdb.Counter{Field: "score", Value: 1}}).
			Save()
		t.AssertNil(err)

		one, err := db.Model(table).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["NAME"], "a")
		t.Assert(one["SCORE"], 11)
	})

	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(table).Data(g.Map{"id": 1, "score": 0}).
			OnDuplicate(g.Map{"score": &gdb.Counter{Field: "score", Value: -5}}).
			Save()
		t.AssertNil(err)

		value, err := db.Model(table).Where("id", 1).Value("score")
		t.AssertNil(err)
		t.Assert(value, 6)
	})

	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(table).Data(g.Map{"id": 2, "name": "c", "score": 7}).
			OnDuplicate(g.Map{"score": gdb.Counter{Field: "score", Value: 1}}).
			Save()
		t.AssertNil(err)

		value, err := db.Model(table).Where("id", 2).Value("score")
		t.AssertNil(err)
		t.Assert(value, 7)
	})
}
