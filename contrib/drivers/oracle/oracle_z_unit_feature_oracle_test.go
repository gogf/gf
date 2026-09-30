// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package oracle_test

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gconv"
)

// oracleFeatureTableName returns a table name starting with `prefix` that is unique within the test run.
func oracleFeatureTableName(prefix string) string {
	return fmt.Sprintf("%s_%d", prefix, gtime.TimestampMicro()%1e9)
}

// oracleFeatureExec executes `statements` in order and stops the test run if one of them fails.
func oracleFeatureExec(statements ...string) {
	for _, statement := range statements {
		if _, err := db.Exec(ctx, statement); err != nil {
			gtest.Fatal(err)
		}
	}
}

// oracleFeatureOpen returns a database object for the test database with `extra` as its Extra option.
func oracleFeatureOpen(extra string, maxOpen int) gdb.DB {
	node := gdb.ConfigNode{
		Host:             TestDbIP,
		Port:             TestDbPort,
		User:             TestDbUser,
		Pass:             TestDbPass,
		Name:             TestDbName,
		Type:             TestDbType,
		Extra:            extra,
		MaxOpenConnCount: maxOpen,
	}
	newDb, err := gdb.New(node)
	if err != nil {
		gtest.Fatal(err)
	}
	return newDb
}

// oracleFeatureAssertOra asserts that `err` carries the Oracle error `code`.
func oracleFeatureAssertOra(t *gtest.T, err error, code string) {
	t.AssertNE(err, nil)
	if !gstr.Contains(err.Error(), code) {
		t.Error(fmt.Sprintf("expected %s, got: %v", code, err))
	}
}

// oracleFeatureRoundTrips returns the number of round trips the session of a database object opened
// with `extra` takes to fetch 200 rows.
func oracleFeatureRoundTrips(t *gtest.T, extra string) int {
	const statSql = "SELECT s.VALUE FROM V$MYSTAT s JOIN V$STATNAME n ON s.STATISTIC# = n.STATISTIC# " +
		"WHERE n.NAME = 'SQL*Net roundtrips to/from client'"
	extraDb := oracleFeatureOpen(extra, 1)
	defer extraDb.Close(ctx)

	before, err := extraDb.GetValue(ctx, statSql)
	t.AssertNil(err)
	all, err := extraDb.GetAll(ctx, "SELECT LEVEL LV FROM dual CONNECT BY LEVEL <= 200")
	t.AssertNil(err)
	t.Assert(len(all), 200)
	after, err := extraDb.GetValue(ctx, statSql)
	t.AssertNil(err)
	return after.Int() - before.Int()
}

// Test_Oracle_Date_KeepsTimeOfDay tests that a DATE column is read back with its time of day.
func Test_Oracle_Date_KeepsTimeOfDay(t *testing.T) {
	table := fmt.Sprintf("t_date_%d", gtime.TimestampMicro()%1e9)
	if _, err := db.Exec(ctx, fmt.Sprintf(`CREATE TABLE %s (ID NUMBER(10) NOT NULL, D DATE, PRIMARY KEY (ID))`, table)); err != nil {
		gtest.Fatal(err)
	}
	defer dropTable(table)

	var (
		withTime = gtime.NewFromStr("2026-03-04 15:26:37")
		midnight = gtime.NewFromStr("2026-03-05 00:00:00")
	)
	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(table).Data(g.List{
			{"ID": 1, "D": withTime},
			{"ID": 2, "D": midnight},
		}).Insert()
		t.AssertNil(err)

		stored, err := db.GetValue(ctx, fmt.Sprintf(`SELECT TO_CHAR(D, 'YYYY-MM-DD HH24:MI:SS') FROM %s WHERE ID = 1`, table))
		t.AssertNil(err)
		t.Assert(stored.String(), "2026-03-04 15:26:37")
	})

	gtest.C(t, func(t *gtest.T) {
		one, err := db.Model(table).Where("ID", 1).One()
		t.AssertNil(err)
		t.Assert(one["D"].String(), "2026-03-04 15:26:37")

		one, err = db.Model(table).Where("ID", 2).One()
		t.AssertNil(err)
		t.Assert(one["D"].String(), "2026-03-05 00:00:00")
	})

	gtest.C(t, func(t *gtest.T) {
		value, err := db.Model(table).Where("ID", 1).Value("D")
		t.AssertNil(err)
		t.Assert(value.GTime().String(), "2026-03-04 15:26:37")
	})

	gtest.C(t, func(t *gtest.T) {
		type Row struct {
			Id int
			D  *gtime.Time
		}
		var rows []Row
		err := db.Model(table).OrderAsc("ID").Scan(&rows)
		t.AssertNil(err)
		t.Assert(len(rows), 2)
		t.Assert(rows[0].D.String(), "2026-03-04 15:26:37")
		t.Assert(rows[1].D.String(), "2026-03-05 00:00:00")
	})

	gtest.C(t, func(t *gtest.T) {
		count, err := db.Model(table).Where("D", withTime).Count()
		t.AssertNil(err)
		t.Assert(count, 1)

		count, err = db.Model(table).WhereGT("D", gtime.NewFromStr("2026-03-04 15:00:00")).Count()
		t.AssertNil(err)
		t.Assert(count, 2)
	})
}

// Test_Oracle_Date_LocalTimeZone tests that a DATE column keeps the instant of a time value
// written in any time zone, reading it back in the local time zone.
func Test_Oracle_Date_LocalTimeZone(t *testing.T) {
	table := fmt.Sprintf("t_date_%d", gtime.TimestampMicro()%1e9)
	if _, err := db.Exec(ctx, fmt.Sprintf(`CREATE TABLE %s (ID NUMBER(10) NOT NULL, D DATE, PRIMARY KEY (ID))`, table)); err != nil {
		gtest.Fatal(err)
	}
	defer dropTable(table)

	var (
		local = time.Date(2026, 3, 4, 15, 26, 37, 0, time.Local)
		utc   = time.Date(2026, 3, 4, 15, 26, 37, 0, time.UTC)
	)
	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(table).Data(g.List{
			{"ID": 1, "D": local},
			{"ID": 2, "D": utc},
		}).Insert()
		t.AssertNil(err)

		one, err := db.Model(table).Where("ID", 1).One()
		t.AssertNil(err)
		t.Assert(one["D"].Time().Equal(local), true)
		t.Assert(one["D"].String(), "2026-03-04 15:26:37")

		one, err = db.Model(table).Where("ID", 2).One()
		t.AssertNil(err)
		t.Assert(one["D"].Time().Equal(utc), true)
		t.Assert(one["D"].String(), utc.In(time.Local).Format("2006-01-02 15:04:05"))
	})
}

// Test_Oracle_Rownum_RawSqlShapes tests the LIMIT rewrite of raw sql whatever the letter case and
// the whitespace around its tokens, and of statements starting with WITH or a parenthesized operand.
func Test_Oracle_Rownum_RawSqlShapes(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		all, err := db.GetAll(ctx, fmt.Sprintf("select id from %s order by id limit 2, 3", table))
		t.AssertNil(err)
		t.Assert(gconv.Ints(all.Array("ID")), []int{3, 4, 5})

		all, err = db.GetAll(ctx, fmt.Sprintf("select id from %s order by id desc Limit 3", table))
		t.AssertNil(err)
		t.Assert(gconv.Ints(all.Array("ID")), []int{10, 9, 8})
	})

	gtest.C(t, func(t *gtest.T) {
		all, err := db.GetAll(ctx, fmt.Sprintf(
			"SELECT id\n\tFROM %s\n\tWHERE id > ?\n\tORDER BY id\n\tLIMIT 2\n\tOFFSET 4", table,
		), 1)
		t.AssertNil(err)
		t.Assert(gconv.Ints(all.Array("ID")), []int{6, 7})

		all, err = db.GetAll(ctx, fmt.Sprintf("SELECT id FROM %s ORDER BY id LIMIT\r\n3,\r\n2", table))
		t.AssertNil(err)
		t.Assert(gconv.Ints(all.Array("ID")), []int{4, 5})
	})

	gtest.C(t, func(t *gtest.T) {
		all, err := db.GetAll(ctx, fmt.Sprintf(
			"WITH x AS (SELECT id FROM %s WHERE id > 3) SELECT id FROM x ORDER BY id LIMIT 1, 2", table,
		))
		t.AssertNil(err)
		t.Assert(gconv.Ints(all.Array("ID")), []int{5, 6})
	})

	gtest.C(t, func(t *gtest.T) {
		all, err := db.GetAll(ctx, fmt.Sprintf(
			"WITH x AS (SELECT id FROM %s ORDER BY id DESC LIMIT 3) SELECT id FROM x ORDER BY id", table,
		))
		t.AssertNil(err)
		t.Assert(gconv.Ints(all.Array("ID")), []int{8, 9, 10})
	})

	gtest.C(t, func(t *gtest.T) {
		all, err := db.GetAll(ctx, fmt.Sprintf(
			"(SELECT id FROM %s WHERE id < 3) UNION ALL (SELECT id FROM %s WHERE id > 8) ORDER BY id LIMIT 1, 2",
			table, table,
		))
		t.AssertNil(err)
		t.Assert(gconv.Ints(all.Array("ID")), []int{2, 9})
	})
}

// Test_Oracle_Rownum_HelperColumnHidden tests that the ROW_NUMBER__ column of the ROWNUM rewrite
// never reaches the results of paged queries.
func Test_Oracle_Rownum_HelperColumnHidden(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		all, err := db.Model(table).Order("id").Page(2, 3).All()
		t.AssertNil(err)
		t.Assert(gconv.Ints(all.Array("ID")), []int{4, 5, 6})
		for _, record := range all {
			t.Assert(len(record), 6)
			_, ok := record["ROW_NUMBER__"]
			t.Assert(ok, false)
		}

		all, err = db.Model(table).Fields("id, nickname").Order("id").Page(3, 2).All()
		t.AssertNil(err)
		t.Assert(len(all), 2)
		t.Assert(len(all[0]), 2)
		t.Assert(all[0]["ID"], 5)
		t.Assert(all[1]["NICKNAME"], "name_6")
	})

	gtest.C(t, func(t *gtest.T) {
		one, err := db.Model(table).Order("id").Page(2, 3).One()
		t.AssertNil(err)
		t.Assert(len(one), 6)
		t.Assert(one["ID"], 4)

		value, err := db.Model(table).Fields("nickname").Order("id").Page(2, 3).Value()
		t.AssertNil(err)
		t.Assert(value, "name_4")

		array, err := db.Model(table).Fields("nickname").Order("id").Page(2, 3).Array()
		t.AssertNil(err)
		t.Assert(array, g.Slice{"name_4", "name_5", "name_6"})
	})

	gtest.C(t, func(t *gtest.T) {
		type User struct {
			Id        int
			Nickname  string
			RowNumber int
		}
		var users []User
		err := db.Model(table).Order("id").Page(4, 3).Scan(&users)
		t.AssertNil(err)
		t.Assert(len(users), 1)
		t.Assert(users[0].Id, 10)
		t.Assert(users[0].Nickname, "name_10")
		t.Assert(users[0].RowNumber, 0)
	})

	gtest.C(t, func(t *gtest.T) {
		all, err := db.GetAll(ctx, fmt.Sprintf("SELECT id FROM %s ORDER BY id LIMIT 2, 2", table))
		t.AssertNil(err)
		t.Assert(len(all), 2)
		t.Assert(len(all[0]), 1)
		t.Assert(all[0]["ID"], 3)

		value, err := db.GetValue(ctx, fmt.Sprintf("SELECT nickname FROM %s ORDER BY id LIMIT 2, 1", table))
		t.AssertNil(err)
		t.Assert(value, "name_3")

		array, err := db.GetArray(ctx, fmt.Sprintf("SELECT nickname FROM %s ORDER BY id LIMIT 2, 3", table))
		t.AssertNil(err)
		t.Assert(array, g.Slice{"name_3", "name_4", "name_5"})
	})
}

// Test_Oracle_Rownum_JoinSelectStar tests paging a join whose tables share column names.
func Test_Oracle_Rownum_JoinSelectStar(t *testing.T) {
	table1 := createInitTable()
	defer dropTable(table1)
	table2 := createInitTable()
	defer dropTable(table2)

	gtest.C(t, func(t *gtest.T) {
		all, err := db.Model(table1+" a").
			InnerJoin(table2+" b", "a.id = b.id").
			Fields("*").
			Order("a.id").
			Page(2, 3).
			All()
		t.AssertNil(err)
		t.Assert(gconv.Ints(all.Array("ID")), []int{4, 5, 6})
		t.Assert(len(all[0]), 6)
		t.Assert(all[0]["NICKNAME"], "name_4")
	})

	gtest.C(t, func(t *gtest.T) {
		all, err := db.Model(table1+" a").
			InnerJoin(table2+" b", "a.id = b.id").
			Fields("a.*, b.*").
			Where("a.id", 4).
			All()
		t.AssertNil(err)
		t.Assert(len(all), 1)

		_, err = db.Model(table1+" a").
			InnerJoin(table2+" b", "a.id = b.id").
			Fields("a.*, b.*").
			Order("a.id").
			Page(2, 3).
			All()
		oracleFeatureAssertOra(t, err, "ORA-00918")
	})
}

// Test_Oracle_LastInsertId_TriggerSequence tests that the id a trigger takes from a sequence is
// returned as the last insert id, for single and batch inserts.
func Test_Oracle_LastInsertId_TriggerSequence(t *testing.T) {
	table := createTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		for i := 1; i <= 2; i++ {
			id, err := db.Model(table).Data(g.Map{
				"passport": fmt.Sprintf("seq_%d", i),
				"password": "p",
				"nickname": "n",
			}).InsertAndGetId()
			t.AssertNil(err)
			t.Assert(id, i)
		}

		result, err := db.Model(table).Data(g.Map{"passport": "seq_3", "password": "p", "nickname": "n"}).Insert()
		t.AssertNil(err)
		id, err := result.LastInsertId()
		t.AssertNil(err)
		t.Assert(id, 3)

		value, err := db.Model(table).Where("id", 3).Value("passport")
		t.AssertNil(err)
		t.Assert(value, "seq_3")
	})

	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Data(g.List{
			{"passport": "batch_1", "password": "p", "nickname": "n"},
			{"passport": "batch_2", "password": "p", "nickname": "n"},
			{"passport": "batch_3", "password": "p", "nickname": "n"},
		}).Batch(2).Insert()
		t.AssertNil(err)
		n, err := result.RowsAffected()
		t.AssertNil(err)
		t.Assert(n, 3)
		id, err := result.LastInsertId()
		t.AssertNil(err)
		t.Assert(id, 6)

		array, err := db.Model(table).WhereLike("passport", "batch_%").Order("id").Array("id")
		t.AssertNil(err)
		t.Assert(gconv.Ints(array), []int{4, 5, 6})
	})
}

// Test_Oracle_LastInsertId_NonIntegerPrimaryKey tests inserting without the primary key into
// tables whose primary key is not an integer and is filled by the database.
func Test_Oracle_LastInsertId_NonIntegerPrimaryKey(t *testing.T) {
	var (
		guidTable = oracleFeatureTableName("t_guid")
		rawTable  = oracleFeatureTableName("t_raw")
	)
	oracleFeatureExec(
		fmt.Sprintf(`CREATE TABLE %s (ID VARCHAR2(32) NOT NULL, NAME VARCHAR2(20), PRIMARY KEY (ID))`, guidTable),
		fmt.Sprintf(`CREATE OR REPLACE TRIGGER %s_TRG BEFORE INSERT ON %s FOR EACH ROW
BEGIN
    IF :NEW.ID IS NULL THEN
        :NEW.ID := RAWTOHEX(SYS_GUID());
    END IF;
END;`, guidTable, guidTable),
		fmt.Sprintf(`CREATE TABLE %s (ID RAW(16) DEFAULT SYS_GUID() NOT NULL, NAME VARCHAR2(20), PRIMARY KEY (ID))`, rawTable),
	)
	defer dropTable(guidTable)
	defer dropTable(rawTable)

	gtest.C(t, func(t *gtest.T) {
		_, err := db.Exec(ctx, fmt.Sprintf("INSERT INTO %s(NAME) VALUES(?)", guidTable), "raw_sql")
		t.AssertNil(err)

		_, err = db.Model(guidTable).Data(g.Map{"name": "model"}).Insert()
		t.AssertNil(err)

		value, err := db.Model(guidTable).Where("name", "model").Value("id")
		t.AssertNil(err)
		t.Assert(len(value.String()), 32)
	})

	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(rawTable).Data(g.Map{"name": "model"}).Insert()
		t.AssertNil(err)

		count, err := db.Model(rawTable).Where("name", "model").Count()
		t.AssertNil(err)
		t.Assert(count, 1)
	})

	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(guidTable).Data(g.Map{"name": "last_insert_id"}).Insert()
		t.AssertNil(err)
		affected, err := result.RowsAffected()
		t.AssertNil(err)
		t.Assert(affected, 1)
		_, err = result.LastInsertId()
		t.Assert(gerror.Code(err), gcode.CodeNotSupported)
	})
}

// Test_Oracle_LastInsertId_NoPrimaryKey tests InsertAndGetId on a table without primary key.
func Test_Oracle_LastInsertId_NoPrimaryKey(t *testing.T) {
	table := oracleFeatureTableName("t_nopk")
	oracleFeatureExec(fmt.Sprintf(`CREATE TABLE %s (ID NUMBER(10), NAME VARCHAR2(20))`, table))
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		id, err := db.Model(table).Data(g.Map{"id": 7, "name": "a"}).InsertAndGetId()
		t.AssertNil(err)
		t.Assert(id, 0)

		id, err = db.Model(table).Data(g.Map{"name": "b"}).InsertAndGetId()
		t.AssertNil(err)
		t.Assert(id, 0)

		all, err := db.Model(table).OrderAsc("name").All()
		t.AssertNil(err)
		t.Assert(len(all), 2)
		t.Assert(all[0]["ID"], 7)
		t.Assert(all[1]["ID"].IsNil(), true)
	})
}

// Test_Oracle_LastInsertId_Transaction tests InsertAndGetId inside a transaction, and that a rolled
// back insert still consumes its sequence value.
func Test_Oracle_LastInsertId_Transaction(t *testing.T) {
	table := createTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		var id int64
		err := db.Transaction(ctx, func(ctx context.Context, tx gdb.TX) (err error) {
			id, err = tx.Model(table).Data(g.Map{"passport": "commit", "password": "p", "nickname": "n"}).InsertAndGetId()
			return err
		})
		t.AssertNil(err)
		t.Assert(id, 1)

		value, err := db.Model(table).Where("id", id).Value("passport")
		t.AssertNil(err)
		t.Assert(value, "commit")
	})

	gtest.C(t, func(t *gtest.T) {
		var id int64
		err := db.Transaction(ctx, func(ctx context.Context, tx gdb.TX) (err error) {
			id, err = tx.Model(table).Data(g.Map{"passport": "rollback", "password": "p", "nickname": "n"}).InsertAndGetId()
			if err != nil {
				return err
			}
			return gerror.New("rollback")
		})
		t.AssertNE(err, nil)
		t.Assert(id, 2)

		count, err := db.Model(table).Where("passport", "rollback").Count()
		t.AssertNil(err)
		t.Assert(count, 0)

		id, err = db.Model(table).Data(g.Map{"passport": "after", "password": "p", "nickname": "n"}).InsertAndGetId()
		t.AssertNil(err)
		t.Assert(id, 3)
	})
}

// Test_Oracle_Sequence_Nextval tests reading sequence values through raw sql, the DUAL table and
// an insert that takes the primary key from a sequence.
func Test_Oracle_Sequence_Nextval(t *testing.T) {
	var (
		table    = createTable()
		sequence = table + "_ID_SEQ"
		bigSeq   = oracleFeatureTableName("s_big")
	)
	defer dropTable(table)
	oracleFeatureExec(fmt.Sprintf(`CREATE SEQUENCE %s START WITH 9007199254740993`, bigSeq))
	defer db.Exec(ctx, "DROP SEQUENCE "+bigSeq)

	gtest.C(t, func(t *gtest.T) {
		value, err := db.GetValue(ctx, fmt.Sprintf("SELECT %s.NEXTVAL FROM dual", sequence))
		t.AssertNil(err)
		t.Assert(value.Int64(), 1)

		value, err = db.GetValue(ctx, fmt.Sprintf("SELECT %s.NEXTVAL FROM dual", sequence))
		t.AssertNil(err)
		t.Assert(value.Int64(), 2)

		all, err := db.Model("dual").Fields(sequence + ".NEXTVAL").All()
		t.AssertNil(err)
		t.Assert(len(all), 1)
		t.Assert(all[0]["NEXTVAL"].Int64(), 3)
	})

	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(table).Data(g.Map{
			"id":       gdb.Raw(sequence + ".NEXTVAL"),
			"passport": "from_seq",
			"password": "p",
			"nickname": "n",
		}).Insert()
		t.AssertNil(err)

		value, err := db.Model(table).Where("passport", "from_seq").Value("id")
		t.AssertNil(err)
		t.Assert(value.Int64(), 4)
	})

	gtest.C(t, func(t *gtest.T) {
		value, err := db.GetValue(ctx, fmt.Sprintf("SELECT %s.NEXTVAL FROM dual", bigSeq))
		t.AssertNil(err)
		t.Assert(value.String(), "9007199254740993")
		t.Assert(value.Int64(), int64(9007199254740993))

		value, err = db.GetValue(ctx, fmt.Sprintf("SELECT %s.NEXTVAL FROM dual", bigSeq))
		t.AssertNil(err)
		t.Assert(value.Int64(), int64(9007199254740994))
	})
}

// Test_Oracle_Sequence_NextvalSingleRow tests that reading a sequence value from DUAL with Value and One
// fails with ORA-02287, as they limit the query to one row with a ROWNUM sub-query, where Oracle does
// not allow a sequence. Test_Oracle_Sequence_Nextval reads it with All and raw sql instead.
func Test_Oracle_Sequence_NextvalSingleRow(t *testing.T) {
	var (
		table    = createTable()
		sequence = table + "_ID_SEQ"
	)
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model("dual").Fields(sequence + ".NEXTVAL").Value()
		oracleFeatureAssertOra(t, err, "ORA-02287")
	})

	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model("dual").Fields(sequence + ".NEXTVAL N").One()
		oracleFeatureAssertOra(t, err, "ORA-02287")
	})
}

// Test_Oracle_Sequence_NextvalSave tests that a sequence value given as Raw to Save fails with
// ORA-02287, as its MERGE takes the values from a sub-query, where Oracle does not allow a sequence,
// while Insert writes it into the VALUES clause.
func Test_Oracle_Sequence_NextvalSave(t *testing.T) {
	var (
		table    = createTable()
		sequence = table + "_ID_SEQ"
	)
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(table).Data(g.Map{
			"id": 1, "passport": gdb.Raw(sequence + ".NEXTVAL"), "password": "p", "nickname": "n",
		}).Save()
		oracleFeatureAssertOra(t, err, "ORA-02287")
	})

	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(table).Data(g.Map{
			"id": gdb.Raw(sequence + ".NEXTVAL"), "passport": "p", "password": "p", "nickname": "n",
		}).Insert()
		t.AssertNil(err)

		count, err := db.Model(table).Count()
		t.AssertNil(err)
		t.Assert(count, 1)
	})
}

// Test_Oracle_Placeholder_InStringLiteral tests that a question mark inside a string literal is
// kept as text instead of being turned into a bind placeholder.
func Test_Oracle_Placeholder_InStringLiteral(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		value, err := db.GetValue(ctx, "SELECT '?' FROM dual")
		t.AssertNil(err)
		t.Assert(value.String(), "?")
	})

	gtest.C(t, func(t *gtest.T) {
		count, err := db.Model(table).Where("REGEXP_LIKE(passport, ?)", "^user_10?$").Count()
		t.AssertNil(err)
		t.Assert(count, 2)

		count, err = db.Model(table).Where("REGEXP_LIKE(passport, '^user_10?$')").Count()
		t.AssertNil(err)
		t.Assert(count, 2)
	})
}

// Test_Oracle_Placeholder_InComment tests that question marks inside comments do not disturb the
// binding of the real placeholders.
func Test_Oracle_Placeholder_InComment(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		value, err := db.GetValue(ctx, "SELECT 1 FROM dual -- ?")
		t.AssertNil(err)
		t.Assert(value.Int(), 1)

		value, err = db.GetValue(ctx, "SELECT /* ? */ ? FROM dual", 5)
		t.AssertNil(err)
		t.Assert(value.Int(), 5)

		value, err = db.GetValue(ctx, "SELECT ? FROM dual -- ?", 6)
		t.AssertNil(err)
		t.Assert(value.Int(), 6)
	})

	gtest.C(t, func(t *gtest.T) {
		count, err := db.Model(table).Where("id > ? /* ? */", 8).Count()
		t.AssertNil(err)
		t.Assert(count, 2)
	})
}

// Test_Oracle_Placeholder_DoubleQuoteInLiteral tests that a double quote inside a string literal is kept.
func Test_Oracle_Placeholder_DoubleQuoteInLiteral(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		value, err := db.GetValue(ctx, "SELECT ? FROM dual", `say "hi"`)
		t.AssertNil(err)
		t.Assert(value, `say "hi"`)

		value, err = db.GetValue(ctx, `SELECT 'say "hi"' FROM dual`)
		t.AssertNil(err)
		t.Assert(value.String(), `say "hi"`)
	})

	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(table).Data("nickname", `say "hi"`).Where("id", 1).Update()
		t.AssertNil(err)

		count, err := db.Model(table).Where("nickname", `say "hi"`).Count()
		t.AssertNil(err)
		t.Assert(count, 1)

		count, err = db.Model(table).Where(`nickname = 'say "hi"'`).Count()
		t.AssertNil(err)
		t.Assert(count, 1)
	})
}

// Test_Oracle_Merge_OnlyConflictKey tests Save with data holding only the primary key, which leaves
// nothing to update when the row exists.
func Test_Oracle_Merge_OnlyConflictKey(t *testing.T) {
	table := oracleFeatureTableName("t_mkey")
	oracleFeatureExec(
		fmt.Sprintf(`CREATE TABLE %s (ID NUMBER(10) NOT NULL, NAME VARCHAR2(20), SCORE NUMBER(10), PRIMARY KEY (ID))`, table),
		fmt.Sprintf(`INSERT INTO %s (ID, NAME, SCORE) VALUES (1, 'a', 10)`, table),
	)
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Data(g.Map{"id": 1}).Save()
		t.AssertNil(err)
		n, err := result.RowsAffected()
		t.AssertNil(err)
		t.Assert(n, 0)

		result, err = db.Model(table).Data(g.Map{"id": 1}).Replace()
		t.AssertNil(err)
		n, err = result.RowsAffected()
		t.AssertNil(err)
		t.Assert(n, 0)

		one, err := db.Model(table).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["NAME"], "a")
		t.Assert(one["SCORE"], 10)
	})

	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Data(g.Map{"id": 2}).Save()
		t.AssertNil(err)
		n, err := result.RowsAffected()
		t.AssertNil(err)
		t.Assert(n, 1)

		one, err := db.Model(table).Where("id", 2).One()
		t.AssertNil(err)
		t.Assert(one["ID"], 2)
		t.Assert(one["NAME"].IsNil(), true)
		t.Assert(one["SCORE"].IsNil(), true)
	})
}

// Test_Oracle_Merge_OnDuplicateCounter tests that a gdb.Counter in OnDuplicate adds to the value
// stored in the matched row, as MySQL does, not to the value being inserted.
func Test_Oracle_Merge_OnDuplicateCounter(t *testing.T) {
	table := oracleFeatureTableName("t_mcnt")
	oracleFeatureExec(
		fmt.Sprintf(`CREATE TABLE %s (ID NUMBER(10) NOT NULL, NAME VARCHAR2(20), SCORE NUMBER(10), PRIMARY KEY (ID))`, table),
		fmt.Sprintf(`INSERT INTO %s (ID, NAME, SCORE) VALUES (1, 'a', 10)`, table),
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

// Test_Oracle_Merge_NullConflictKey tests Save with a NULL primary key, which never matches, so the
// trigger fills a new id for every call.
func Test_Oracle_Merge_NullConflictKey(t *testing.T) {
	table := createTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		for i := 0; i < 2; i++ {
			result, err := db.Model(table).Data(g.Map{
				"id":       nil,
				"passport": "null_key",
				"password": "p",
				"nickname": "n",
			}).Save()
			t.AssertNil(err)
			n, err := result.RowsAffected()
			t.AssertNil(err)
			t.Assert(n, 1)
		}

		array, err := db.Model(table).Where("passport", "null_key").Order("id").Array("id")
		t.AssertNil(err)
		t.Assert(gconv.Ints(array), []int{1, 2})
	})
}

// Test_Oracle_Merge_NoConflictKey tests the upsert operations on data without primary key, with and
// without OnConflict.
func Test_Oracle_Merge_NoConflictKey(t *testing.T) {
	table := oracleFeatureTableName("t_mnok")
	oracleFeatureExec(fmt.Sprintf(`CREATE TABLE %s (ID NUMBER(10), NAME VARCHAR2(20))`, table))
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		data := g.Map{"id": 5, "name": "a"}
		_, err := db.Model(table).Data(data).Save()
		t.Assert(gerror.Code(err), gcode.CodeMissingParameter)
		_, err = db.Model(table).Data(data).Replace()
		t.Assert(gerror.Code(err), gcode.CodeMissingParameter)
		_, err = db.Model(table).Data(data).InsertIgnore()
		t.Assert(gerror.Code(err), gcode.CodeMissingParameter)

		count, err := db.Model(table).Count()
		t.AssertNil(err)
		t.Assert(count, 0)
	})

	gtest.C(t, func(t *gtest.T) {
		userTable := createTable()
		defer dropTable(userTable)

		_, err := db.Model(userTable).Data(g.Map{"passport": "p", "password": "p", "nickname": "n"}).Save()
		t.Assert(gerror.Code(err), gcode.CodeMissingParameter)
	})

	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Data(g.Map{"id": 5, "name": "a"}).OnConflict("id").Save()
		t.AssertNil(err)
		n, _ := result.RowsAffected()
		t.Assert(n, 1)

		result, err = db.Model(table).Data(g.Map{"id": 5, "name": "b"}).OnConflict("id").Save()
		t.AssertNil(err)
		n, _ = result.RowsAffected()
		t.Assert(n, 1)

		all, err := db.Model(table).All()
		t.AssertNil(err)
		t.Assert(len(all), 1)
		t.Assert(all[0]["NAME"], "b")
	})
}

// Test_Oracle_Merge_OnConflictReplaceAndIgnore tests that Replace and InsertIgnore detect conflicts
// on the columns given by OnConflict, as the driver builds their MERGE from them.
func Test_Oracle_Merge_OnConflictReplaceAndIgnore(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		table := oracleFeatureTableName("t_mocf")
		oracleFeatureExec(
			fmt.Sprintf(`CREATE TABLE %s (ID NUMBER(10), NAME VARCHAR2(20))`, table),
			fmt.Sprintf(`INSERT INTO %s (ID, NAME) VALUES (5, 'a')`, table),
		)
		defer dropTable(table)

		_, err := db.Model(table).Data(g.Map{"id": 5, "name": "b"}).OnConflict("id").InsertIgnore()
		t.AssertNil(err)

		value, err := db.Model(table).Where("id", 5).Value("name")
		t.AssertNil(err)
		t.Assert(value, "a")
	})

	gtest.C(t, func(t *gtest.T) {
		table := oracleFeatureTableName("t_mocf")
		oracleFeatureExec(
			fmt.Sprintf(`CREATE TABLE %s (ID NUMBER(10), NAME VARCHAR2(20))`, table),
			fmt.Sprintf(`INSERT INTO %s (ID, NAME) VALUES (5, 'a')`, table),
		)
		defer dropTable(table)

		_, err := db.Model(table).Data(g.Map{"id": 5, "name": "b"}).OnConflict("id").Replace()
		t.AssertNil(err)

		value, err := db.Model(table).Where("id", 5).Value("name")
		t.AssertNil(err)
		t.Assert(value, "b")
	})

	gtest.C(t, func(t *gtest.T) {
		table := createInitTable()
		defer dropTable(table)

		result, err := db.Model(table).Data(g.Map{
			"id": 11, "passport": "user_3", "password": "p", "nickname": "ignored",
		}).OnConflict("passport").InsertIgnore()
		t.AssertNil(err)
		n, err := result.RowsAffected()
		t.AssertNil(err)
		t.Assert(n, 0)

		count, err := db.Model(table).Where("passport", "user_3").Count()
		t.AssertNil(err)
		t.Assert(count, 1)
	})

	// A plain Insert does not use OnConflict, so that a value of a type OnConflict does not accept
	// is ignored as before.
	gtest.C(t, func(t *gtest.T) {
		table := createTable()
		defer dropTable(table)

		_, err := db.Model(table).Data(g.Map{
			"id": 1, "passport": "p", "password": "p", "nickname": "n",
		}).OnConflict(g.Map{"id": 1}).Insert()
		t.AssertNil(err)

		count, err := db.Model(table).Count()
		t.AssertNil(err)
		t.Assert(count, 1)
	})
}

// Test_Oracle_Merge_InsertIgnoreConflict tests that InsertIgnore reports no affected row and keeps
// the existing row on a conflict.
func Test_Oracle_Merge_InsertIgnoreConflict(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Data(g.Map{
			"id": 1, "passport": "ignored", "password": "p", "nickname": "ignored",
		}).InsertIgnore()
		t.AssertNil(err)
		n, err := result.RowsAffected()
		t.AssertNil(err)
		t.Assert(n, 0)

		one, err := db.Model(table).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["PASSPORT"], "user_1")
		t.Assert(one["NICKNAME"], "name_1")
	})

	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Data(g.List{
			{"id": 2, "passport": "ignored", "password": "p", "nickname": "ignored"},
			{"id": 11, "passport": "user_11", "password": "p", "nickname": "name_11"},
		}).InsertIgnore()
		t.AssertNil(err)
		n, err := result.RowsAffected()
		t.AssertNil(err)
		t.Assert(n, 1)

		count, err := db.Model(table).Count()
		t.AssertNil(err)
		t.Assert(count, TableSize+1)
		value, err := db.Model(table).Where("id", 2).Value("nickname")
		t.AssertNil(err)
		t.Assert(value, "name_2")
	})
}

// Test_Oracle_Merge_RawExec tests a hand-written MERGE statement with bind arguments through Exec.
func Test_Oracle_Merge_RawExec(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	mergeSql := fmt.Sprintf(`MERGE INTO %s t
USING (SELECT ? ID, ? NICK FROM dual) s ON (t.ID = s.ID)
WHEN MATCHED THEN UPDATE SET t.NICKNAME = s.NICK DELETE WHERE t.NICKNAME = 'drop_me'
WHEN NOT MATCHED THEN INSERT (ID, PASSPORT, PASSWORD, NICKNAME) VALUES (s.ID, 'merged', 'p', s.NICK)`, table)

	gtest.C(t, func(t *gtest.T) {
		result, err := db.Exec(ctx, mergeSql, 2, "merged_2")
		t.AssertNil(err)
		n, err := result.RowsAffected()
		t.AssertNil(err)
		t.Assert(n, 1)
		value, err := db.Model(table).Where("id", 2).Value("nickname")
		t.AssertNil(err)
		t.Assert(value, "merged_2")
	})

	gtest.C(t, func(t *gtest.T) {
		result, err := db.Exec(ctx, mergeSql, 11, "merged_11")
		t.AssertNil(err)
		n, err := result.RowsAffected()
		t.AssertNil(err)
		t.Assert(n, 1)
		one, err := db.Model(table).Where("id", 11).One()
		t.AssertNil(err)
		t.Assert(one["PASSPORT"], "merged")
		t.Assert(one["NICKNAME"], "merged_11")
	})

	gtest.C(t, func(t *gtest.T) {
		_, err := db.Exec(ctx, mergeSql, 3, "drop_me")
		t.AssertNil(err)
		count, err := db.Model(table).Where("id", 3).Count()
		t.AssertNil(err)
		t.Assert(count, 0)
		count, err = db.Model(table).Count()
		t.AssertNil(err)
		t.Assert(count, TableSize)
	})
}

// Test_Oracle_TableFields_TypeStrings tests the field types that TableFields reports for Oracle
// types without a MySQL counterpart, and that a NUMBER without precision keeps its decimals.
func Test_Oracle_TableFields_TypeStrings(t *testing.T) {
	table := oracleFeatureTableName("t_types")
	oracleFeatureExec(fmt.Sprintf(`CREATE TABLE %s (
    ID NUMBER(10) NOT NULL,
    N NUMBER,
    I INTEGER,
    F FLOAT,
    BD BINARY_DOUBLE,
    TZ TIMESTAMP WITH TIME ZONE,
    TL TIMESTAMP WITH LOCAL TIME ZONE,
    IDS INTERVAL DAY TO SECOND,
    IYM INTERVAL YEAR TO MONTH,
    NV NVARCHAR2(10),
    NC NCHAR(4),
    VB VARCHAR2(6 BYTE),
    VC VARCHAR2(3 CHAR),
    R RAW(8),
    PRIMARY KEY (ID)
)`, table))
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		fields, err := db.TableFields(ctx, table)
		t.AssertNil(err)
		t.Assert(len(fields), 14)
		expect := map[string]string{
			"ID":  "INT(10,0)",
			"N":   "INT(,)",
			"I":   "INT(,0)",
			"F":   "FLOAT(126,)",
			"BD":  "BINARY_DOUBLE(8)",
			"TZ":  "TIMESTAMP(6) WITH TIME ZONE(13)",
			"TL":  "TIMESTAMP(6) WITH LOCAL TIME ZONE(11)",
			"IDS": "INTERVAL DAY(2) TO SECOND(6)(11)",
			"IYM": "INTERVAL YEAR(2) TO MONTH(5)",
			"NV":  "NVARCHAR2(20)",
			"NC":  "NCHAR(8)",
			"VB":  "VARCHAR2(6)",
			"VC":  "VARCHAR2(12)",
			"R":   "RAW(8)",
		}
		for name, fieldType := range expect {
			t.Assert(fields[name].Type, fieldType)
		}
		t.Assert(fields["ID"].Key, "PRI")
		t.Assert(fields["ID"].Null, false)
		t.Assert(fields["N"].Null, true)
	})

	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(table).Data(g.List{
			{"id": 1, "n": 3.14159},
			{"id": 2, "n": "12345678901234567890"},
			{"id": 3, "n": "0.000000000000000000000000000001"},
		}).Insert()
		t.AssertNil(err)

		array, err := db.Model(table).Order("id").Array("n")
		t.AssertNil(err)
		t.Assert(array[0].Float64(), 3.14159)
		t.Assert(array[1].String(), "12345678901234567890")
		t.Assert(array[2].Float64(), 1e-30)
	})
}

// Test_Oracle_TableFields_CommentAndDefault tests what TableFields reports for a column that has a
// comment and a default value.
func Test_Oracle_TableFields_CommentAndDefault(t *testing.T) {
	table := oracleFeatureTableName("t_cmt")
	oracleFeatureExec(
		fmt.Sprintf(`CREATE TABLE %s (ID NUMBER(10) NOT NULL, NAME VARCHAR2(20) DEFAULT 'dflt', PRIMARY KEY (ID))`, table),
		fmt.Sprintf(`COMMENT ON COLUMN %s.NAME IS 'the name'`, table),
	)
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		comment, err := db.GetValue(ctx,
			"SELECT COMMENTS FROM USER_COL_COMMENTS WHERE TABLE_NAME = ? AND COLUMN_NAME = 'NAME'",
			strings.ToUpper(table),
		)
		t.AssertNil(err)
		t.Assert(comment, "the name")

		_, err = db.Model(table).Data(g.Map{"id": 1}).Insert()
		t.AssertNil(err)
		value, err := db.Model(table).Where("id", 1).Value("name")
		t.AssertNil(err)
		t.Assert(value, "dflt")
	})

	gtest.C(t, func(t *gtest.T) {
		fields, err := db.TableFields(ctx, table)
		t.AssertNil(err)
		t.Assert(fields["NAME"].Type, "VARCHAR2(20)")
		t.Assert(fields["NAME"].Comment, "")
		t.Assert(fields["NAME"].Default, nil)
		t.Assert(fields["NAME"].Extra, "")
	})
}

// Test_Oracle_TableFields_View tests TableFields and queries on a view, and that Tables lists no view.
func Test_Oracle_TableFields_View(t *testing.T) {
	var (
		table = createInitTable()
		view  = oracleFeatureTableName("v_user")
	)
	defer dropTable(table)
	oracleFeatureExec(fmt.Sprintf(`CREATE VIEW %s AS SELECT ID, NICKNAME FROM %s WHERE ID > 8`, view, table))
	defer db.Exec(ctx, "DROP VIEW "+view)

	gtest.C(t, func(t *gtest.T) {
		fields, err := db.TableFields(ctx, view)
		t.AssertNil(err)
		t.Assert(len(fields), 2)
		t.Assert(fields["ID"].Type, "INT(10,0)")
		t.Assert(fields["ID"].Key, "")
		t.Assert(fields["ID"].Null, false)
		t.Assert(fields["NICKNAME"].Type, "VARCHAR2(45)")
	})

	gtest.C(t, func(t *gtest.T) {
		all, err := db.Model(view).Order("id").All()
		t.AssertNil(err)
		t.Assert(gconv.Ints(all.Array("ID")), []int{9, 10})
		t.Assert(all[1]["NICKNAME"], "name_10")

		value, err := db.Model(view).Where("nickname", "name_9").Value("id")
		t.AssertNil(err)
		t.Assert(value, 9)
	})

	gtest.C(t, func(t *gtest.T) {
		tables, err := db.Tables(ctx)
		t.AssertNil(err)
		t.AssertIN(strings.ToUpper(table), tables)
		t.AssertNI(strings.ToUpper(view), tables)
	})
}

// Test_Oracle_Model_SchemaQualifiedTable tests querying and changing a table named with its schema.
func Test_Oracle_Model_SchemaQualifiedTable(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		one, err := db.Model("SYSTEM."+table).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["PASSPORT"], "user_1")

		one, err = db.Model("system."+table).Fields("id, nickname").Where("id", 2).One()
		t.AssertNil(err)
		t.Assert(len(one), 2)
		t.Assert(one["NICKNAME"], "name_2")

		all, err := db.Model("SYSTEM."+table+" u").Where("u.id > ?", 8).Order("u.id").All()
		t.AssertNil(err)
		t.Assert(gconv.Ints(all.Array("ID")), []int{9, 10})
	})

	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model("SYSTEM."+table).Where("id", 4).Delete()
		t.AssertNil(err)
		n, _ := result.RowsAffected()
		t.Assert(n, 1)

		_, err = db.Model("SYSTEM." + table).Data(g.Map{
			"id": 11, "passport": "user_11", "password": "p", "nickname": "name_11",
		}).Insert()
		t.AssertNil(err)
		count, err := db.Model(table).Count()
		t.AssertNil(err)
		t.Assert(count, TableSize)
	})
}

// Test_Oracle_Model_SchemaQualifiedFieldTypes tests the operations that need the field types of a
// table named with its schema.
func Test_Oracle_Model_SchemaQualifiedFieldTypes(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		table := createTable()
		defer dropTable(table)

		fields, err := db.TableFields(ctx, "SYSTEM."+table)
		t.AssertNil(err)
		t.Assert(len(fields), 6)
	})

	gtest.C(t, func(t *gtest.T) {
		table := createTable()
		defer dropTable(table)

		id, err := db.Model("SYSTEM." + table).Data(g.Map{"passport": "p", "password": "p", "nickname": "n"}).InsertAndGetId()
		t.AssertNil(err)
		value, err := db.Model(table).Where("passport", "p").Value("id")
		t.AssertNil(err)
		t.Assert(value, 1)
		t.Assert(id, 1)
	})

	gtest.C(t, func(t *gtest.T) {
		table := createInitTable()
		defer dropTable(table)

		_, err := db.Model("SYSTEM."+table).Data("nickname", "schema").Where("id", 3).Update()
		t.AssertNil(err)

		value, err := db.Model(table).Where("id", 3).Value("nickname")
		t.AssertNil(err)
		t.Assert(value, "schema")
	})

	gtest.C(t, func(t *gtest.T) {
		table := oracleFeatureTableName("t_sclob")
		oracleFeatureExec(fmt.Sprintf(`CREATE TABLE %s (ID NUMBER(10) NOT NULL, C CLOB, PRIMARY KEY (ID))`, table))
		defer dropTable(table)

		text := strings.Repeat("x", 5000)
		_, err := db.Model(table).Data(g.Map{"id": 1, "c": text}).Insert()
		t.AssertNil(err)

		_, err = db.Model("SYSTEM." + table).Data(g.Map{"id": 2, "c": text}).Insert()
		t.AssertNil(err)

		value, err := db.Model(table).Where("id", 2).Value("c")
		t.AssertNil(err)
		t.Assert(len(value.String()), 5000)
	})
}

// Test_Oracle_Model_OtherOwnerTable tests the model operations that need the fields of a table
// qualified with an owner other than the connected user, while the connected user has a table of
// the same name with other columns.
func Test_Oracle_Model_OtherOwnerTable(t *testing.T) {
	var (
		owner      = "GF_OWNER"
		table      = oracleFeatureTableName("t_owner")
		ownerTable = owner + "." + table
	)
	_, _ = db.Exec(ctx, "DROP USER "+owner+" CASCADE")
	oracleFeatureExec(
		fmt.Sprintf(`CREATE USER %s IDENTIFIED BY gf_owner DEFAULT TABLESPACE USERS QUOTA UNLIMITED ON USERS`, owner),
		fmt.Sprintf(`CREATE TABLE %s (ID NUMBER(10) NOT NULL, PASSPORT VARCHAR2(45), NICKNAME VARCHAR2(45), PRIMARY KEY (ID))`, ownerTable),
		fmt.Sprintf(`CREATE TABLE %s (ID NUMBER(10) NOT NULL, NAME VARCHAR2(45), PRIMARY KEY (ID))`, table),
	)
	defer func() {
		dropTable(table)
		_, _ = db.Exec(ctx, "DROP USER "+owner+" CASCADE")
	}()

	gtest.C(t, func(t *gtest.T) {
		t.Assert(db.GetTableNameForFields(ownerTable), ownerTable)
		t.Assert(db.GetTableNameForFields(ownerTable+" u"), ownerTable)
		t.Assert(db.GetTableNameForFields(ownerTable+" AS u, "+table+" t"), ownerTable)
		t.Assert(db.GetTableNameForFields(table+" u"), table)
		t.Assert(db.GetTableNameForFields(","+ownerTable), ownerTable)

		fields, err := db.Model(ownerTable + " u").TableFields(ownerTable + " u")
		t.AssertNil(err)
		t.Assert(len(fields), 3)
		t.AssertNE(fields["PASSPORT"], nil)

		has, err := db.GetCore().HasField(ctx, ownerTable, "nickname")
		t.AssertNil(err)
		t.Assert(has, true)
	})

	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(ownerTable).Data(g.Map{"id": 1, "passport": "user_1", "nickname": "name_1"}).Insert()
		t.AssertNil(err)

		result, err := db.Model(ownerTable).Data("nickname", "owner").Where("id", 1).Update()
		t.AssertNil(err)
		n, _ := result.RowsAffected()
		t.Assert(n, 1)

		one, err := db.GetOne(ctx, fmt.Sprintf(`SELECT PASSPORT, NICKNAME FROM %s WHERE ID = 1`, ownerTable))
		t.AssertNil(err)
		t.Assert(one["PASSPORT"], "user_1")
		t.Assert(one["NICKNAME"], "owner")

		one, err = db.Model(ownerTable).FieldsEx("passport").Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(len(one), 2)
		t.Assert(one["NICKNAME"], "owner")

		count, err := db.Model(table).Count()
		t.AssertNil(err)
		t.Assert(count, 0)
	})
}

// Test_Oracle_Lob_LargeRoundTrip tests CLOB and BLOB values beyond the 4000 and 32767 byte bind
// limits through insert, update and select.
func Test_Oracle_Lob_LargeRoundTrip(t *testing.T) {
	table := oracleFeatureTableName("t_lob")
	oracleFeatureExec(fmt.Sprintf(`CREATE TABLE %s (ID NUMBER(10) NOT NULL, C CLOB, B BLOB, PRIMARY KEY (ID))`, table))
	defer dropTable(table)

	var (
		text4k    = strings.Repeat("a", 4001)
		text32k   = strings.Repeat("0123456789", 3300)
		multibyte = strings.Repeat("中", 1500)
		binary    = make([]byte, 40000)
	)
	for i := range binary {
		binary[i] = byte(i % 256)
	}

	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(table).Data(g.List{
			{"id": 1, "c": text4k, "b": binary[:5000]},
			{"id": 2, "c": text32k, "b": binary},
			{"id": 3, "c": multibyte},
		}).Insert()
		t.AssertNil(err)

		all, err := db.Model(table).Order("id").All()
		t.AssertNil(err)
		t.Assert(len(all), 3)
		t.Assert(all[0]["C"].String() == text4k, true)
		t.Assert(all[0]["B"].Bytes(), binary[:5000])
		t.Assert(all[1]["C"].String() == text32k, true)
		t.Assert(all[1]["B"].Bytes(), binary)
		t.Assert(all[2]["C"].String() == multibyte, true)

		lengths, err := db.Model(table).Fields("DBMS_LOB.GETLENGTH(c) L").Order("id").Array()
		t.AssertNil(err)
		t.Assert(gconv.Ints(lengths), []int{4001, 33000, 1500})
	})

	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(table).Data(g.Map{"c": text32k + text32k, "b": binary[:33000]}).Where("id", 1).Update()
		t.AssertNil(err)

		one, err := db.Model(table).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(len(one["C"].String()), 66000)
		t.Assert(one["C"].String() == text32k+text32k, true)
		t.Assert(one["B"].Bytes(), binary[:33000])
	})
}

// Test_Oracle_Lob_ComparisonRejected tests that Oracle rejects a CLOB in comparisons, ORDER BY,
// DISTINCT and GROUP BY, while LIKE and the DBMS_LOB functions work.
func Test_Oracle_Lob_ComparisonRejected(t *testing.T) {
	table := oracleFeatureTableName("t_lobcmp")
	oracleFeatureExec(fmt.Sprintf(`CREATE TABLE %s (ID NUMBER(10) NOT NULL, C CLOB, PRIMARY KEY (ID))`, table))
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(table).Data(g.List{
			{"id": 1, "c": "abc"},
			{"id": 2, "c": strings.Repeat("x", 5000) + "needle"},
		}).Insert()
		t.AssertNil(err)

		_, err = db.Model(table).Where("c", "abc").Count()
		oracleFeatureAssertOra(t, err, "ORA-00932")
		_, err = db.Model(table).Order("c").All()
		oracleFeatureAssertOra(t, err, "ORA-00932")
		_, err = db.Model(table).Fields("DISTINCT c").All()
		oracleFeatureAssertOra(t, err, "ORA-00932")
		_, err = db.Model(table).Fields("c").Group("c").All()
		oracleFeatureAssertOra(t, err, "ORA-00932")
	})

	gtest.C(t, func(t *gtest.T) {
		array, err := db.Model(table).WhereLike("c", "%needle").Array("id")
		t.AssertNil(err)
		t.Assert(gconv.Ints(array), []int{2})

		array, err = db.Model(table).Where("DBMS_LOB.INSTR(c, ?) > 0", "needle").Array("id")
		t.AssertNil(err)
		t.Assert(gconv.Ints(array), []int{2})

		array, err = db.Model(table).Where("DBMS_LOB.COMPARE(c, TO_CLOB(?)) = 0", "abc").Array("id")
		t.AssertNil(err)
		t.Assert(gconv.Ints(array), []int{1})
	})
}

// Test_Oracle_Varchar2_ByteAndCharSemantics tests the length limits of VARCHAR2 in bytes and in
// characters, and NVARCHAR2 with characters beyond the basic plane.
func Test_Oracle_Varchar2_ByteAndCharSemantics(t *testing.T) {
	table := oracleFeatureTableName("t_vchar")
	oracleFeatureExec(fmt.Sprintf(
		`CREATE TABLE %s (ID NUMBER(10) NOT NULL, VB VARCHAR2(6 BYTE), VC VARCHAR2(3 CHAR), NV NVARCHAR2(10), PRIMARY KEY (ID))`,
		table,
	))
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(table).Data(g.Map{"id": 1, "vb": "中文字"}).Insert()
		oracleFeatureAssertOra(t, err, "ORA-12899")

		_, err = db.Model(table).Data(g.Map{"id": 2, "vc": "abcd"}).Insert()
		oracleFeatureAssertOra(t, err, "ORA-12899")

		count, err := db.Model(table).Count()
		t.AssertNil(err)
		t.Assert(count, 0)
	})

	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(table).Data(g.Map{"id": 3, "vb": "中文", "vc": "中文字", "nv": "中文字符串😀"}).Insert()
		t.AssertNil(err)

		one, err := db.Model(table).Where("id", 3).One()
		t.AssertNil(err)
		t.Assert(one["VB"], "中文")
		t.Assert(one["VC"], "中文字")
		t.Assert(one["NV"], "中文字符串😀")

		value, err := db.Model(table).Where("nv", "中文字符串😀").Value("id")
		t.AssertNil(err)
		t.Assert(value, 3)
	})
}

// Test_Oracle_WhereIn_Over1000 tests that an IN list is limited to 1000 values, which gf does not split.
func Test_Oracle_WhereIn_Over1000(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	ids := make([]int, 1001)
	for i := range ids {
		ids[i] = i + 1
	}

	gtest.C(t, func(t *gtest.T) {
		count, err := db.Model(table).WhereIn("id", ids[:1000]).Count()
		t.AssertNil(err)
		t.Assert(count, TableSize)

		count, err = db.Model(table).WhereNotIn("id", ids[2:1000]).Count()
		t.AssertNil(err)
		t.Assert(count, 2)
	})

	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(table).WhereIn("id", ids).Count()
		oracleFeatureAssertOra(t, err, "ORA-01795")

		_, err = db.Model(table).Where("id", ids).Count()
		oracleFeatureAssertOra(t, err, "ORA-01795")

		_, err = db.Model(table).WhereNotIn("id", ids).Count()
		oracleFeatureAssertOra(t, err, "ORA-01795")
	})
}

// Test_Oracle_EmptyString_IsNull tests that Oracle stores an empty string as NULL and how conditions
// on an empty string behave.
func Test_Oracle_EmptyString_IsNull(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(table).Data(g.Map{"id": 11, "passport": "", "password": "p", "nickname": "n"}).Insert()
		oracleFeatureAssertOra(t, err, "ORA-01400")
	})

	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(table).Data("create_time", "").Where("id", 1).Update()
		t.AssertNil(err)

		value, err := db.Model(table).Where("id", 1).Value("create_time")
		t.AssertNil(err)
		t.Assert(value.IsNil(), true)

		count, err := db.Model(table).WhereNull("create_time").Count()
		t.AssertNil(err)
		t.Assert(count, 1)

		count, err = db.Model(table).Where("create_time", "").Count()
		t.AssertNil(err)
		t.Assert(count, 0)

		count, err = db.Model(table).WhereNot("create_time", "").Count()
		t.AssertNil(err)
		t.Assert(count, 0)

		count, err = db.Model(table).OmitEmptyWhere().Where("create_time", "").Count()
		t.AssertNil(err)
		t.Assert(count, TableSize)
	})
}

// Test_Oracle_Extra_OptionParsing tests how the Extra option reaches the underlying driver:
// malformed options are dropped, the others are passed on and validated by the driver.
func Test_Oracle_Extra_OptionParsing(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		for _, extra := range []string{"PREFETCH_ROWS", "a=b=c&PREFETCH_ROWS", "PREFETCH_ROWS=abc"} {
			extraDb := oracleFeatureOpen(extra, 0)
			value, err := extraDb.GetValue(ctx, "SELECT 1 FROM dual")
			extraDb.Close(ctx)
			t.AssertNil(err)
			t.Assert(value, 1)
		}
	})

	gtest.C(t, func(t *gtest.T) {
		expect := map[string]string{
			"FOO=bar":                "unknown URL option: FOO",
			"connection timeout=abc": "CONNECTION TIMEOUT value must be an integer",
			"LOB FETCH=bad":          "LOB FETCH value should be",
		}
		for extra, message := range expect {
			extraDb := oracleFeatureOpen(extra, 0)
			_, err := extraDb.GetValue(ctx, "SELECT 1 FROM dual")
			extraDb.Close(ctx)
			t.AssertNE(err, nil)
			t.Assert(gstr.Contains(err.Error(), message), true)
		}
	})

	gtest.C(t, func(t *gtest.T) {
		table := oracleFeatureTableName("t_lobpre")
		oracleFeatureExec(fmt.Sprintf(`CREATE TABLE %s (ID NUMBER(10) NOT NULL, C CLOB, B BLOB, PRIMARY KEY (ID))`, table))
		defer dropTable(table)

		text := strings.Repeat("中文", 3000)
		_, err := db.Model(table).Data(g.Map{"id": 1, "c": text, "b": []byte(text)}).Insert()
		t.AssertNil(err)

		extraDb := oracleFeatureOpen("LOB FETCH=PRE&CONNECTION TIMEOUT=5", 0)
		defer extraDb.Close(ctx)
		one, err := extraDb.Model(table).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["C"].String() == text, true)
		t.Assert(one["B"].Bytes(), []byte(text))
	})
}

// Test_Oracle_Extra_PrefetchRows tests that PREFETCH_ROWS given in Extra replaces the default and
// changes the number of round trips a fetch takes.
func Test_Oracle_Extra_PrefetchRows(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			oneRow      = oracleFeatureRoundTrips(t, "PREFETCH_ROWS=1")
			defaultRows = oracleFeatureRoundTrips(t, "")
			manyRows    = oracleFeatureRoundTrips(t, "prefetch_rows=1000")
		)
		t.AssertGT(oneRow, 100)
		t.AssertLT(defaultRows, oneRow)
		t.AssertLT(manyRows, defaultRows)
		t.AssertLE(manyRows, 10)
	})
}

// Test_Oracle_Extra_NlsTerritory tests that LANGUAGE and TERRITORY given in Extra set the session
// NLS parameters, and that numbers are written correctly under a comma decimal separator.
func Test_Oracle_Extra_NlsTerritory(t *testing.T) {
	table := oracleFeatureTableName("t_nls")
	oracleFeatureExec(fmt.Sprintf(
		`CREATE TABLE %s (ID NUMBER(10) NOT NULL, AMOUNT NUMBER(18,2), RATIO BINARY_DOUBLE, PRIMARY KEY (ID))`, table,
	))
	defer dropTable(table)

	extraDb := oracleFeatureOpen("LANGUAGE=GERMAN&TERRITORY=GERMANY", 0)
	defer extraDb.Close(ctx)

	gtest.C(t, func(t *gtest.T) {
		one, err := extraDb.GetOne(ctx, "SELECT TO_CHAR(DATE '2026-03-04', 'fmDay') D, TO_CHAR(1234.5, 'FM9G999D0') N FROM dual")
		t.AssertNil(err)
		t.Assert(one["D"], "Mittwoch")
		t.Assert(one["N"], "1.234,5")

		one, err = db.GetOne(ctx, "SELECT TO_CHAR(DATE '2026-03-04', 'fmDay') D, TO_CHAR(1234.5, 'FM9G999D0') N FROM dual")
		t.AssertNil(err)
		t.Assert(one["D"], "Wednesday")
		t.Assert(one["N"], "1,234.5")
	})

	gtest.C(t, func(t *gtest.T) {
		_, err := extraDb.Model(table).Data(g.Map{"id": 1, "amount": 2.5, "ratio": 1.25}).Insert()
		t.AssertNil(err)

		one, err := db.Model(table).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["AMOUNT"].Float64(), 2.5)
		t.Assert(one["RATIO"].Float64(), 1.25)

		value, err := extraDb.Model(table).Where("amount", 2.5).Value("ratio")
		t.AssertNil(err)
		t.Assert(value.Float64(), 1.25)
	})
}

// Test_Oracle_ConnectBy_Hierarchy tests hierarchical queries with CONNECT BY, including paging
// that keeps ORDER SIBLINGS BY.
func Test_Oracle_ConnectBy_Hierarchy(t *testing.T) {
	table := oracleFeatureTableName("t_tree")
	oracleFeatureExec(fmt.Sprintf(
		`CREATE TABLE %s (ID NUMBER(10) NOT NULL, PARENT_ID NUMBER(10), NAME VARCHAR2(20), PRIMARY KEY (ID))`, table,
	))
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(table).Data(g.List{
			{"id": 1, "parent_id": nil, "name": "root"},
			{"id": 2, "parent_id": 1, "name": "b"},
			{"id": 3, "parent_id": 1, "name": "a"},
			{"id": 4, "parent_id": 2, "name": "b1"},
			{"id": 5, "parent_id": 3, "name": "a1"},
			{"id": 6, "parent_id": 3, "name": "a2"},
		}).Insert()
		t.AssertNil(err)

		all, err := db.GetAll(ctx, fmt.Sprintf(`SELECT ID, LEVEL LV, SYS_CONNECT_BY_PATH(NAME, '/') PATH, CONNECT_BY_ISLEAF LEAF
FROM %s START WITH ID = ? CONNECT BY PRIOR ID = PARENT_ID ORDER SIBLINGS BY NAME`, table), 1)
		t.AssertNil(err)
		t.Assert(gconv.Ints(all.Array("ID")), []int{1, 3, 5, 6, 2, 4})
		t.Assert(gconv.Ints(all.Array("LV")), []int{1, 2, 3, 3, 2, 3})
		t.Assert(all[2]["PATH"], "/root/a/a1")
		t.Assert(gconv.Ints(all.Array("LEAF")), []int{0, 0, 1, 1, 0, 1})
	})

	gtest.C(t, func(t *gtest.T) {
		all, err := db.GetAll(ctx, fmt.Sprintf(
			`SELECT ID FROM %s START WITH PARENT_ID IS NULL CONNECT BY PRIOR ID = PARENT_ID ORDER SIBLINGS BY NAME LIMIT 1, 3`,
			table,
		))
		t.AssertNil(err)
		t.Assert(gconv.Ints(all.Array("ID")), []int{3, 5, 6})

		array, err := db.GetArray(ctx, "SELECT LEVEL FROM dual CONNECT BY LEVEL <= ?", 4)
		t.AssertNil(err)
		t.Assert(gconv.Ints(array), []int{1, 2, 3, 4})
	})
}

// Test_Oracle_Listagg_WithinGroup tests LISTAGG ... WITHIN GROUP through the Model with GROUP BY and HAVING.
func Test_Oracle_Listagg_WithinGroup(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		all, err := db.Model(table).
			Fields("MOD(id, 2) M, LISTAGG(id, ',') WITHIN GROUP (ORDER BY id) IDS").
			Group("MOD(id, 2)").
			Order("M").
			All()
		t.AssertNil(err)
		t.Assert(len(all), 2)
		t.Assert(all[0]["IDS"], "2,4,6,8,10")
		t.Assert(all[1]["IDS"], "1,3,5,7,9")
	})

	gtest.C(t, func(t *gtest.T) {
		value, err := db.Model(table).
			Fields("LISTAGG(nickname, '|') WITHIN GROUP (ORDER BY id DESC) NAMES").
			Where("id > ?", 7).
			Value()
		t.AssertNil(err)
		t.Assert(value, "name_10|name_9|name_8")

		all, err := db.Model(table).
			Fields("MOD(id, 3) M, LISTAGG(id, ',') WITHIN GROUP (ORDER BY id) IDS").
			Group("MOD(id, 3)").
			Having("COUNT(*) > ?", 3).
			All()
		t.AssertNil(err)
		t.Assert(len(all), 1)
		t.Assert(all[0]["IDS"], "1,4,7,10")
	})
}

// Test_Oracle_WindowFunction_Page tests window functions combined with paging, where ROW_NUMBER is
// computed over the whole table before the page is taken.
func Test_Oracle_WindowFunction_Page(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		all, err := db.Model(table).
			Fields("id, ROW_NUMBER() OVER (ORDER BY id DESC) RN, RANK() OVER (ORDER BY MOD(id, 3)) RK, "+
				"LAG(id) OVER (ORDER BY id) PREV_ID, LEAD(id) OVER (ORDER BY id) NEXT_ID").
			Order("id").
			Page(2, 3).
			All()
		t.AssertNil(err)
		t.Assert(len(all), 3)
		t.Assert(len(all[0]), 5)
		t.Assert(gconv.Ints(all.Array("ID")), []int{4, 5, 6})
		t.Assert(gconv.Ints(all.Array("RN")), []int{7, 6, 5})
		t.Assert(gconv.Ints(all.Array("RK")), []int{4, 8, 1})
		t.Assert(gconv.Ints(all.Array("PREV_ID")), []int{3, 4, 5})
		t.Assert(gconv.Ints(all.Array("NEXT_ID")), []int{5, 6, 7})
	})

	gtest.C(t, func(t *gtest.T) {
		one, err := db.Model(table).
			Fields("id, LAG(id) OVER (ORDER BY id) PREV_ID").
			Order("id").
			One()
		t.AssertNil(err)
		t.Assert(one["ID"], 1)
		t.Assert(one["PREV_ID"].IsNil(), true)
	})
}

// Test_Oracle_PlSql_OutBind tests anonymous PL/SQL blocks through Exec, returning values into
// positional and named out binds.
func Test_Oracle_PlSql_OutBind(t *testing.T) {
	table := createTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		var id int64
		_, err := db.Exec(ctx, fmt.Sprintf(
			"BEGIN INSERT INTO %s (PASSPORT, PASSWORD, NICKNAME) VALUES (?, 'p', 'n') RETURNING ID INTO ?; END;", table,
		), "positional", sql.Out{Dest: &id})
		t.AssertNil(err)
		t.Assert(id, 1)

		value, err := db.Model(table).Where("id", id).Value("passport")
		t.AssertNil(err)
		t.Assert(value, "positional")
	})

	gtest.C(t, func(t *gtest.T) {
		var id int64
		_, err := db.Exec(ctx, fmt.Sprintf(
			"BEGIN INSERT INTO %s (PASSPORT, PASSWORD, NICKNAME) VALUES (:passport, 'p', 'n') RETURNING ID INTO :id; END;", table,
		), sql.Named("passport", "named"), sql.Named("id", sql.Out{Dest: &id}))
		t.AssertNil(err)
		t.Assert(id, 2)

		value, err := db.Model(table).Where("id", id).Value("passport")
		t.AssertNil(err)
		t.Assert(value, "named")
	})

	gtest.C(t, func(t *gtest.T) {
		_, err := db.Exec(ctx, fmt.Sprintf(
			"BEGIN FOR i IN 1..3 LOOP INSERT INTO %s (PASSPORT, PASSWORD, NICKNAME) VALUES ('loop_' || i, 'p', ?); END LOOP; END;",
			table,
		), "loop")
		t.AssertNil(err)

		array, err := db.Model(table).Where("nickname", "loop").Order("id").Array("passport")
		t.AssertNil(err)
		t.Assert(array, g.Slice{"loop_1", "loop_2", "loop_3"})
	})
}

// Test_Oracle_LastInsertId_NonIntegerPrimaryKeyId tests that InsertAndGetId inserts into a table whose
// primary key is not an integer and returns 0, as there is no integer id to return.
func Test_Oracle_LastInsertId_NonIntegerPrimaryKeyId(t *testing.T) {
	table := oracleFeatureTableName("t_rawid")
	oracleFeatureExec(fmt.Sprintf(`CREATE TABLE %s (ID RAW(16) DEFAULT SYS_GUID() NOT NULL, NAME VARCHAR2(20), PRIMARY KEY (ID))`, table))
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		id, err := db.Model(table).Data(g.Map{"name": "a"}).InsertAndGetId()
		t.Assert(gerror.Code(err), gcode.CodeNotSupported)
		t.Assert(id, 0)

		count, err := db.Model(table).Where("name", "a").Count()
		t.AssertNil(err)
		t.Assert(count, 1)
	})
}

func Test_Oracle_TableAlias_AsKeyword(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		all, err := db.Model(fmt.Sprintf("(SELECT ID FROM %s) AS u", table)).All()
		t.AssertNil(err)
		t.Assert(len(all), TableSize)
	})

	gtest.C(t, func(t *gtest.T) {
		sub := db.Model(table).Fields("ID").Where("ID", 1)
		all, err := db.Model("? AS a", sub).All()
		t.AssertNil(err)
		t.Assert(len(all), 1)

		value, err := db.Model("(?) AS sub", sub).Value("sub.ID")
		t.AssertNil(err)
		t.Assert(value.Int(), 1)

		value, err = db.Model("? AS a, ? AS b", sub, sub).Where("a.ID=b.ID").Value("a.ID")
		t.AssertNil(err)
		t.Assert(value.Int(), 1)
	})

	gtest.C(t, func(t *gtest.T) {
		count, err := db.Model(table, "a").
			LeftJoin(fmt.Sprintf("(SELECT ID FROM %s) AS b", table), "a.ID=b.ID").
			Where("b.ID IS NOT NULL").
			Count()
		t.AssertNil(err)
		t.Assert(count, TableSize)
	})

	gtest.C(t, func(t *gtest.T) {
		value, err := db.Model(fmt.Sprintf("(SELECT ID AS USER_ID, NICKNAME FROM %s) AS u", table)).
			Where("u.USER_ID", 2).
			Value("u.NICKNAME")
		t.AssertNil(err)
		t.Assert(value, "name_2")

		sub := db.Raw(fmt.Sprintf("SELECT 'a,  b AS c' AS V FROM %s WHERE ID=1", table))
		value, err = db.Model("(?) AS u", sub).Value("u.V")
		t.AssertNil(err)
		t.Assert(value, "a,  b AS c")
	})
}
