// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package dm_test

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"math"
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

// dmFeatureTableName returns a table name starting with `prefix` that is unique within the test run.
func dmFeatureTableName(prefix string) string {
	return fmt.Sprintf("%s_%d", prefix, gtime.TimestampMicro()%1e9)
}

// dmFeatureExec executes `statements` in order and stops the test run if one of them fails.
func dmFeatureExec(statements ...string) {
	for _, statement := range statements {
		if _, err := db.Exec(ctx, statement); err != nil {
			gtest.Fatal(err)
		}
	}
}

// dmFeatureOpen returns a database object for the test database with `extra` as its Extra option.
func dmFeatureOpen(extra string, maxOpen int) gdb.DB {
	node := gdb.ConfigNode{
		Host:             TestDBHost,
		Port:             TestDBPort,
		User:             TestDBUser,
		Pass:             TestDBPass,
		Name:             TestDBName,
		Type:             TestDBType,
		Extra:            extra,
		MaxOpenConnCount: maxOpen,
	}
	newDb, err := gdb.New(node)
	if err != nil {
		gtest.Fatal(err)
	}
	return newDb
}

// The codes of the DM server errors that the tests of this file expect besides the dmErr constants.
const (
	dmErrParamIncompatible = -5403
	dmErrColumnTooLong     = -6169
	dmErrLobCompare        = -2685
	dmErrDataOutOfRange    = 6092
	dmErrInvalidSchema     = -2103
	dmErrIdentityInsertOff = -2723
)

// Test_DM_Time_HistoricalOffsets tests that a time of a zone whose offset is not the current one,
// like the local mean time of Asia/Shanghai before 1901 or its daylight saving time of 1986 to
// 1991, is written as its wall clock and read back as the same time, through Data, Where and Exec.
// The underlying driver converts a time from the offset of its zone, in whole minutes, to the time
// zone of the session, which is the current offset of the time zone, so that it wrote
// 1880-10-24 00:00 as 1880-10-23 23:55 and 1988-07-01 12:00 as 11:00. The CI runs the tests in
// the local time zone Asia/Shanghai, which the pool of this test configures as its Timezone.
func Test_DM_Time_HistoricalOffsets(t *testing.T) {
	location, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		gtest.Fatal(err)
	}
	pool, err := gdb.New(gdb.ConfigNode{
		Host:     TestDBHost,
		Port:     TestDBPort,
		User:     TestDBUser,
		Pass:     TestDBPass,
		Name:     TestDBName,
		Type:     TestDBType,
		Timezone: location.String(),
	})
	if err != nil {
		gtest.Fatal(err)
	}
	defer pool.Close(ctx)

	table := dmFeatureTableName("t_time_offset")
	if _, err = pool.Exec(ctx, fmt.Sprintf(
		`CREATE TABLE %s (ID NUMBER(10) NOT NULL, D DATE, TS TIMESTAMP, PRIMARY KEY (ID))`, table,
	)); err != nil {
		gtest.Fatal(err)
	}
	defer dropTable(table)

	var values = []string{"1880-10-24 00:00:00", "1988-07-01 12:00:00", "1000-01-01 00:00:00", "2022-02-14 12:00:01"}
	gtest.C(t, func(t *gtest.T) {
		for i, value := range values {
			at, err := time.ParseInLocation(time.DateTime, value, location)
			t.AssertNil(err)
			_, err = pool.Model(table).Data(g.Map{"ID": i, "D": at, "TS": at}).Insert()
			t.AssertNil(err)

			stored, err := pool.GetOne(ctx, fmt.Sprintf(
				`SELECT TO_CHAR(D, 'YYYY-MM-DD') D, TO_CHAR(TS, 'YYYY-MM-DD HH24:MI:SS') TS FROM %s WHERE ID = ?`, table,
			), i)
			t.AssertNil(err)
			t.Assert(stored["D"], value[:10])
			t.Assert(stored["TS"], value)

			one, err := pool.Model(table).Where("TS", at).One()
			t.AssertNil(err)
			t.Assert(one["ID"].Int(), i)
			t.Assert(one["D"].Time().In(location).Format(time.DateOnly), value[:10])
			t.Assert(one["TS"].Time().In(location).Format(time.DateTime), value)
			t.Assert(one["TS"].Time().Equal(at), true)
		}
	})
	gtest.C(t, func(t *gtest.T) {
		at, err := time.ParseInLocation(time.DateTime, values[1], location)
		t.AssertNil(err)
		_, err = pool.Exec(ctx, fmt.Sprintf(`UPDATE %s SET TS = ? WHERE ID = 0`, table), at)
		t.AssertNil(err)
		value, err := pool.GetValue(ctx, fmt.Sprintf(`SELECT TO_CHAR(TS, 'YYYY-MM-DD HH24:MI:SS') FROM %s WHERE ID = 0`, table))
		t.AssertNil(err)
		t.Assert(value, values[1])
	})
}

// Test_DM_Time_ConfiguredTimezone tests that a time is written as its wall clock in the configured
// Timezone, also in its daylight saving time, and read back as the same time.
func Test_DM_Time_ConfiguredTimezone(t *testing.T) {
	pool, err := gdb.New(gdb.ConfigNode{
		Host:     TestDBHost,
		Port:     TestDBPort,
		User:     TestDBUser,
		Pass:     TestDBPass,
		Name:     TestDBName,
		Type:     TestDBType,
		Timezone: "Europe/Moscow",
	})
	if err != nil {
		gtest.Fatal(err)
	}
	defer pool.Close(ctx)

	table := dmFeatureTableName("t_time_tz")
	if _, err = pool.Exec(ctx, fmt.Sprintf(`CREATE TABLE %s (ID NUMBER(10) NOT NULL, TS TIMESTAMP, PRIMARY KEY (ID))`, table)); err != nil {
		gtest.Fatal(err)
	}
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		for i, c := range []struct {
			at     time.Time
			stored string
		}{
			{time.Date(2022, 2, 14, 12, 0, 0, 0, time.UTC), "2022-02-14 15:00:00"},
			{time.Date(1988, 7, 1, 3, 0, 0, 0, time.UTC), "1988-07-01 07:00:00"},
		} {
			_, err := pool.Model(table).Data(g.Map{"ID": i, "TS": c.at}).Insert()
			t.AssertNil(err)

			stored, err := pool.GetValue(ctx, fmt.Sprintf(`SELECT TO_CHAR(TS, 'YYYY-MM-DD HH24:MI:SS') FROM %s WHERE ID = ?`, table), i)
			t.AssertNil(err)
			t.Assert(stored, c.stored)

			value, err := pool.Model(table).Where("TS", c.at).Value("TS")
			t.AssertNil(err)
			t.Assert(value.Time().Equal(c.at), true)

			count, err := pool.Model(table).Where("TS", sql.NullTime{Time: c.at, Valid: true}).Count()
			t.AssertNil(err)
			t.Assert(count, 1)
		}
	})
}

// Test_DM_Time_ColumnTypes tests the time written to the columns of each kind: a character column
// keeps its text, a column with a time zone keeps its instant also for another offset than the one
// of the session, and a column with the local time zone keeps its instant also in a session whose
// time zone is not the one of the database, through Data, and a DATE read in a time zone whose
// daylight saving time begins at midnight keeps its date.
func Test_DM_Time_ColumnTypes(t *testing.T) {
	newPool := func(timezone string) gdb.DB {
		pool, err := gdb.New(gdb.ConfigNode{
			Host:     TestDBHost,
			Port:     TestDBPort,
			User:     TestDBUser,
			Pass:     TestDBPass,
			Name:     TestDBName,
			Type:     TestDBType,
			Timezone: timezone,
		})
		if err != nil {
			gtest.Fatal(err)
		}
		return pool
	}
	table := dmFeatureTableName("t_time_types")
	dmFeatureExec(fmt.Sprintf(`CREATE TABLE %s (
		ID NUMBER(10) NOT NULL, V VARCHAR(19), C CHAR(19), D DATE,
		TZ TIMESTAMP WITH TIME ZONE, LTZ TIMESTAMP WITH LOCAL TIME ZONE, PRIMARY KEY (ID)
	)`, table))
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		at := time.Date(2026, 7, 15, 12, 0, 0, 0, time.Local)
		_, err := db.Model(table).Data(g.Map{"ID": 1, "V": at, "C": at}).Insert()
		t.AssertNil(err)
		one, err := db.Model(table).Where("ID", 1).One()
		t.AssertNil(err)
		t.Assert(one["V"].String(), "2026-07-15 12:00:00")
		t.Assert(gstr.Trim(one["C"].String()), "2026-07-15 12:00:00")
	})
	gtest.C(t, func(t *gtest.T) {
		pool := newPool("America/New_York")
		defer pool.Close(ctx)
		for i, at := range []time.Time{
			time.Date(2026, 1, 15, 12, 0, 0, 0, time.UTC),
			time.Date(2026, 7, 15, 12, 0, 0, 0, time.FixedZone("", 9*3600)),
		} {
			_, err := pool.Model(table).Data(g.Map{"ID": 10 + i, "TZ": at}).Insert()
			t.AssertNil(err)
			value, err := pool.Model(table).Where("ID", 10+i).Value("TZ")
			t.AssertNil(err)
			t.Assert(value.Time().Equal(at), true)
		}
	})
	gtest.C(t, func(t *gtest.T) {
		pool := newPool("Europe/London")
		defer pool.Close(ctx)
		at := time.Date(2026, 7, 15, 11, 0, 0, 0, time.UTC)
		_, err := pool.Model(table).Data(g.Map{"ID": 20, "LTZ": at}).Insert()
		t.AssertNil(err)
		value, err := pool.Model(table).Where("ID", 20).Value("LTZ")
		t.AssertNil(err)
		t.Assert(value.Time().Equal(at), true)
		node := gdb.ConfigNode{
			Host:             TestDBHost,
			Port:             TestDBPort,
			User:             TestDBUser,
			Pass:             TestDBPass,
			Name:             TestDBName,
			Type:             TestDBType,
			Extra:            "timeZone=60",
			MaxOpenConnCount: 1,
		}
		fixedPool, err := gdb.New(node)
		t.AssertNil(err)
		defer fixedPool.Close(ctx)
		stored, err := fixedPool.GetValue(ctx, fmt.Sprintf(`SELECT TO_CHAR(LTZ, 'YYYY-MM-DD HH24:MI:SS') FROM %s WHERE ID = 20`, table))
		t.AssertNil(err)
		t.Assert(stored, "2026-07-15 12:00:00")
	})
	gtest.C(t, func(t *gtest.T) {
		dmFeatureExec(fmt.Sprintf(`INSERT INTO %s (ID, D) VALUES (30, DATE '2026-03-08')`, table))
		pool := newPool("America/Havana")
		defer pool.Close(ctx)
		one, err := pool.Model(table).Where("ID", 30).One()
		t.AssertNil(err)
		t.Assert(one["D"].String(), "2026-03-08")
	})
	gtest.C(t, func(t *gtest.T) {
		pool := newPool("Asia/Shanghai")
		defer pool.Close(ctx)
		at := time.Date(2026, 7, 15, 12, 0, 0, 0, time.Local)
		st, err := pool.Prepare(ctx, fmt.Sprintf("INSERT INTO %s (ID, D) VALUES (?, ?)", table))
		t.AssertNil(err)
		defer st.Close()
		_, err = st.Exec(40, at)
		t.AssertNil(err)
		one, err := pool.Model(table).Where("ID", 40).One()
		t.AssertNil(err)
		t.Assert(one["D"].String(), at.Format("2006-01-02"))
	})
}

// Test_DM_Time_RawArguments tests that a time passed as a raw argument, of any type and time zone,
// is bound as its wall clock in the configured Timezone: written by Exec and by a prepared
// statement, and compared by Where.
func Test_DM_Time_RawArguments(t *testing.T) {
	location, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		gtest.Fatal(err)
	}
	pool, err := gdb.New(gdb.ConfigNode{
		Host:     TestDBHost,
		Port:     TestDBPort,
		User:     TestDBUser,
		Pass:     TestDBPass,
		Name:     TestDBName,
		Type:     TestDBType,
		Timezone: location.String(),
	})
	if err != nil {
		gtest.Fatal(err)
	}
	defer pool.Close(ctx)

	table := dmFeatureTableName("t_time_raw")
	dmFeatureExec(fmt.Sprintf(`CREATE TABLE %s (ID NUMBER(10) NOT NULL, TS TIMESTAMP, PRIMARY KEY (ID))`, table))
	defer dropTable(table)

	// 1988-07-01 03:00 UTC is 12:00 of the daylight saving time of Asia/Shanghai, +09:00.
	at := time.Date(1988, 7, 1, 3, 0, 0, 0, time.UTC)
	stored := "1988-07-01 12:00:00"
	values := []any{
		at,
		at.In(time.FixedZone("", 0)),
		at.In(time.FixedZone("", 8*3600)),
		at.In(location),
		&at,
		gtime.NewFromTime(at.In(time.FixedZone("", 0))),
		*gtime.NewFromTime(at),
		sql.NullTime{Time: at.In(time.FixedZone("", 0)), Valid: true},
	}
	gtest.C(t, func(t *gtest.T) {
		for i, value := range values {
			_, err := pool.Exec(ctx, fmt.Sprintf(`INSERT INTO %s (ID, TS) VALUES (?, ?)`, table), i, value)
			t.AssertNil(err)
			s, err := pool.GetValue(ctx, fmt.Sprintf(`SELECT TO_CHAR(TS, 'YYYY-MM-DD HH24:MI:SS') FROM %s WHERE ID = ?`, table), i)
			t.AssertNil(err)
			t.Assert(s, stored)
		}
		for _, value := range values {
			count, err := pool.Model(table).Where("TS", value).Count()
			t.AssertNil(err)
			t.Assert(count, len(values))
		}
	})
	gtest.C(t, func(t *gtest.T) {
		st, err := pool.Prepare(ctx, fmt.Sprintf(`INSERT INTO %s (ID, TS) VALUES (?, ?)`, table))
		t.AssertNil(err)
		defer st.Close()
		_, err = st.Exec(100, at.In(time.FixedZone("", 0)))
		t.AssertNil(err)
		s, err := pool.GetValue(ctx, fmt.Sprintf(`SELECT TO_CHAR(TS, 'YYYY-MM-DD HH24:MI:SS') FROM %s WHERE ID = 100`, table))
		t.AssertNil(err)
		t.Assert(s, stored)
	})
}

// Test_DM_Time_SkippedWallClock tests that a wall clock that the daylight saving time of the
// configured Timezone skips is read as it is stored.
func Test_DM_Time_SkippedWallClock(t *testing.T) {
	pool, err := gdb.New(gdb.ConfigNode{
		Host:     TestDBHost,
		Port:     TestDBPort,
		User:     TestDBUser,
		Pass:     TestDBPass,
		Name:     TestDBName,
		Type:     TestDBType,
		Timezone: "Europe/London",
	})
	if err != nil {
		gtest.Fatal(err)
	}
	defer pool.Close(ctx)

	table := dmFeatureTableName("t_time_skipped")
	dmFeatureExec(fmt.Sprintf(`CREATE TABLE %s (ID NUMBER(10) NOT NULL, TS TIMESTAMP, PRIMARY KEY (ID))`, table))
	defer dropTable(table)

	// Europe/London skips 2026-03-29 01:00 to 02:00.
	dmFeatureExec(fmt.Sprintf(`INSERT INTO %s (ID, TS) VALUES (1, TIMESTAMP '2026-03-29 01:30:00')`, table))
	gtest.C(t, func(t *gtest.T) {
		value, err := pool.Model(table).Where("ID", 1).Value("TS")
		t.AssertNil(err)
		t.Assert(value.Time().Format(time.DateTime), "2026-03-29 01:30:00")
	})
}

// Test_DM_Date_DropsTimeOfDay tests that a DATE column keeps only the date of a time value, while
// a DATETIME column keeps its time of day.
// Note: Oracle DATE keeps the time of day; DM DATE is date-only like MySQL DATE, and DM has
// DATETIME and TIMESTAMP for a time of day.
func Test_DM_Date_DropsTimeOfDay(t *testing.T) {
	table := fmt.Sprintf("t_date_%d", gtime.TimestampMicro()%1e9)
	if _, err := db.Exec(ctx, fmt.Sprintf(`CREATE TABLE %s (ID NUMBER(10) NOT NULL, D DATE, DT DATETIME, PRIMARY KEY (ID))`, table)); err != nil {
		gtest.Fatal(err)
	}
	defer dropTable(table)

	var (
		withTime = gtime.NewFromStr("2026-03-04 15:26:37")
		midnight = gtime.NewFromStr("2026-03-05 00:00:00")
	)
	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(table).Data(g.List{
			{"ID": 1, "D": withTime, "DT": withTime},
			{"ID": 2, "D": midnight, "DT": midnight},
		}).Insert()
		t.AssertNil(err)

		stored, err := db.GetOne(ctx, fmt.Sprintf(
			`SELECT TO_CHAR(D, 'YYYY-MM-DD HH24:MI:SS') D, TO_CHAR(DT, 'YYYY-MM-DD HH24:MI:SS') DT FROM %s WHERE ID = 1`, table,
		))
		t.AssertNil(err)
		t.Assert(stored["D"], "2026-03-04 00:00:00")
		t.Assert(stored["DT"], "2026-03-04 15:26:37")
	})

	gtest.C(t, func(t *gtest.T) {
		one, err := db.Model(table).Where("ID", 1).One()
		t.AssertNil(err)
		t.Assert(one["D"].String(), "2026-03-04")
		t.Assert(one["DT"].String(), "2026-03-04 15:26:37")

		one, err = db.Model(table).Where("ID", 2).One()
		t.AssertNil(err)
		t.Assert(one["D"].String(), "2026-03-05")
		t.Assert(one["DT"].String(), "2026-03-05 00:00:00")
	})

	gtest.C(t, func(t *gtest.T) {
		value, err := db.Model(table).Where("ID", 1).Value("D")
		t.AssertNil(err)
		t.Assert(value.GTime().String(), "2026-03-04 00:00:00")

		value, err = db.Model(table).Where("ID", 1).Value("DT")
		t.AssertNil(err)
		t.Assert(value.GTime().String(), "2026-03-04 15:26:37")
	})

	gtest.C(t, func(t *gtest.T) {
		type Row struct {
			Id int
			D  *gtime.Time
			Dt *gtime.Time
		}
		var rows []Row
		err := db.Model(table).OrderAsc("ID").Scan(&rows)
		t.AssertNil(err)
		t.Assert(len(rows), 2)
		t.Assert(rows[0].D.String(), "2026-03-04 00:00:00")
		t.Assert(rows[0].Dt.String(), "2026-03-04 15:26:37")
		t.Assert(rows[1].D.String(), "2026-03-05 00:00:00")
		t.Assert(rows[1].Dt.String(), "2026-03-05 00:00:00")
	})

	gtest.C(t, func(t *gtest.T) {
		count, err := db.Model(table).Where("D", withTime).Count()
		t.AssertNil(err)
		t.Assert(count, 0)

		count, err = db.Model(table).Where("D", "2026-03-04").Count()
		t.AssertNil(err)
		t.Assert(count, 1)

		count, err = db.Model(table).Where("DT", withTime).Count()
		t.AssertNil(err)
		t.Assert(count, 1)

		count, err = db.Model(table).WhereGT("D", gtime.NewFromStr("2026-03-04 15:00:00")).Count()
		t.AssertNil(err)
		t.Assert(count, 1)

		count, err = db.Model(table).WhereGT("DT", gtime.NewFromStr("2026-03-04 15:00:00")).Count()
		t.AssertNil(err)
		t.Assert(count, 2)
	})
}

// Test_DM_Date_LocalTimeZone tests that a DATETIME column keeps the instant of a time value
// written in any time zone, reading it back in the local time zone.
// Note: DM DATE is date-only, so the time of day is kept in a DATETIME column.
func Test_DM_Date_LocalTimeZone(t *testing.T) {
	table := fmt.Sprintf("t_date_%d", gtime.TimestampMicro()%1e9)
	if _, err := db.Exec(ctx, fmt.Sprintf(`CREATE TABLE %s (ID NUMBER(10) NOT NULL, D DATETIME, PRIMARY KEY (ID))`, table)); err != nil {
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

// Test_DM_Limit_RawSqlShapes tests LIMIT in raw sql whatever the letter case and the whitespace
// around its tokens, and in statements starting with WITH or a parenthesized operand.
// Note: the driver passes LIMIT to DM, which supports it, instead of rewriting it to ROWNUM.
func Test_DM_Limit_RawSqlShapes(t *testing.T) {
	table := createInitUserTable()
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

// Test_DM_Page_NoExtraColumn tests that paged queries return the selected columns only, with no
// helper column like the ROW_NUMBER__ of a ROWNUM rewrite.
func Test_DM_Page_NoExtraColumn(t *testing.T) {
	table := createInitUserTable()
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

// Test_DM_Page_JoinSelectStar tests paging a join whose tables share column names.
func Test_DM_Page_JoinSelectStar(t *testing.T) {
	table1 := createInitUserTable()
	defer dropTable(table1)
	table2 := createInitUserTable()
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

		// Note: Oracle rejects the ROWNUM sub-query of the page for its duplicate column names
		// (ORA-00918); DM pages the join with LIMIT, as MySQL does.
		all, err = db.Model(table1+" a").
			InnerJoin(table2+" b", "a.id = b.id").
			Fields("a.*, b.*").
			Order("a.id").
			Page(2, 3).
			All()
		t.AssertNil(err)
		t.Assert(gconv.Ints(all.Array("ID")), []int{4, 5, 6})
		t.Assert(len(all[0]), 6)
		t.Assert(all[0]["NICKNAME"], "name_4")
	})
}

// Test_DM_Limit_Overflow tests paging with a count that reaches the largest int, whose sum with
// the offset does not fit in an int.
func Test_DM_Limit_Overflow(t *testing.T) {
	table := createInitUserTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		all, err := db.Model(table).Order("id").Limit(1, math.MaxInt).All()
		t.AssertNil(err)
		t.Assert(gconv.Ints(all.Array("ID")), []int{2, 3, 4, 5, 6, 7, 8, 9, 10})

		all, err = db.Model(table).Order("id").Limit(math.MaxInt, math.MaxInt).All()
		t.AssertNil(err)
		t.Assert(len(all), 0)

		all, err = db.Model(table).Order("id").Limit(math.MaxInt).All()
		t.AssertNil(err)
		t.Assert(len(all), TableSize)
	})
}

// Test_DM_Limit_OffsetInSubQuery tests a query with an offset used as an operand of a union and
// as a sub-query, which select the columns of the query only.
func Test_DM_Limit_OffsetInSubQuery(t *testing.T) {
	table := createInitUserTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		all, err := db.UnionAll(
			db.Model(table).Order("id").Limit(2),
			db.Model(table).Order("id").Limit(1),
		).All()
		t.AssertNil(err)
		t.Assert(gconv.Ints(all.Array("ID")), []int{1, 2, 1})
	})

	// Note: Oracle 11g fails these (ORA-01789, ORA-00913), as its ROWNUM rewrite of an offset adds
	// a column; DM runs LIMIT with an offset as it is.
	gtest.C(t, func(t *gtest.T) {
		all, err := db.UnionAll(
			db.Model(table).Order("id").Limit(2, 3),
			db.Model(table).Order("id").Limit(1),
		).All()
		t.AssertNil(err)
		t.Assert(gconv.Ints(all.Array("ID")), []int{3, 4, 5, 1})
	})

	gtest.C(t, func(t *gtest.T) {
		all, err := db.Model(table).
			Where("id IN(?)", db.Model(table).Fields("id").Order("id").Limit(2, 3)).
			Order("id").
			All()
		t.AssertNil(err)
		t.Assert(gconv.Ints(all.Array("ID")), []int{3, 4, 5})
	})
}

// Test_DM_LastInsertId_TriggerSequence tests that the id a trigger takes from a sequence is
// returned as the last insert id, for single and batch inserts.
func Test_DM_LastInsertId_TriggerSequence(t *testing.T) {
	table := createUserTable()
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

	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Data(g.List{
			{"passport": "multi_1", "password": "p", "nickname": "n"},
			{"passport": "multi_2", "password": "p", "nickname": "n"},
		}).Insert()
		t.AssertNil(err)
		n, err := result.RowsAffected()
		t.AssertNil(err)
		t.Assert(n, 2)
		id, err := result.LastInsertId()
		t.AssertNil(err)
		t.Assert(id, 8)

		value, err := db.Model(table).Where("id", id).Value("passport")
		t.AssertNil(err)
		t.Assert(value, "multi_2")
	})
}

// Test_DM_LastInsertId_IdentityBatch tests the last insert id of several records inserted into a
// table with an identity column: the identity that the server generates for the last record, which
// leaves the column out, or else the value that the last record gives the column, even 0, which DM
// keeps as it is with IDENTITY_INSERT on and rejects without it.
func Test_DM_LastInsertId_IdentityBatch(t *testing.T) {
	table := createTableWithIdentity(dmFeatureTableName("t_identity_batch"))
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Data(g.List{{"ACCOUNT_NAME": "a1"}, {"ACCOUNT_NAME": "a2"}}).Insert()
		t.AssertNil(err)
		n, err := result.RowsAffected()
		t.AssertNil(err)
		t.Assert(n, 2)
		id, err := result.LastInsertId()
		t.AssertNil(err)
		name, err := db.Model(table).Where("ID", id).Value("ACCOUNT_NAME")
		t.AssertNil(err)
		t.Assert(name, "a2")
	})
	gtest.C(t, func(t *gtest.T) {
		list := g.List{{"ID": 100, "ACCOUNT_NAME": "b1"}, {"ID": 0, "ACCOUNT_NAME": "b2"}}
		_, err := db.Model(table).Data(list).Insert()
		assertDmError(t, err, dmErrIdentityInsertOff)

		err = db.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
			if _, err := tx.Exec(fmt.Sprintf(`SET IDENTITY_INSERT %s ON`, table)); err != nil {
				return err
			}
			result, err := tx.Model(table).Data(list).Insert()
			t.AssertNil(err)
			id, err := result.LastInsertId()
			t.AssertNil(err)
			t.Assert(id, 0)
			_, err = tx.Exec(fmt.Sprintf(`SET IDENTITY_INSERT %s OFF`, table))
			return err
		})
		t.AssertNil(err)
		name, err := db.Model(table).Where("ID", 0).Value("ACCOUNT_NAME")
		t.AssertNil(err)
		t.Assert(name, "b2")
	})
}

// Test_DM_Records_AllOrNone tests that an insert or a save of several records writes none of them
// when the last one fails: on its own, and in the transaction of the caller, which goes on. DM
// inserts the records by one statement, and merges them by one statement each, which run in one
// transaction, or back to a savepoint in the transaction of the caller.
func Test_DM_Records_AllOrNone(t *testing.T) {
	table := createUserTable()
	defer dropTable(table)

	insertList := g.List{
		{"passport": "insert_1", "password": "p", "nickname": "n"},
		{"passport": nil, "password": "p", "nickname": "n"},
	}
	saveList := g.List{
		{"id": 101, "passport": "save_1", "password": "p", "nickname": "n"},
		{"id": 102, "passport": nil, "password": "p", "nickname": "n"},
	}
	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(table).Data(insertList).Insert()
		assertDmError(t, err, dmErrNotNullViolation)
		_, err = db.Model(table).Data(saveList).Save()
		assertDmError(t, err, dmErrNotNullViolation)

		count, err := db.Model(table).Count()
		t.AssertNil(err)
		t.Assert(count, 0)
	})
	gtest.C(t, func(t *gtest.T) {
		err := db.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
			if _, err := tx.Model(table).Data(g.Map{"passport": "kept", "password": "p", "nickname": "n"}).Insert(); err != nil {
				return err
			}
			_, err := tx.Model(table).Data(insertList).Insert()
			assertDmError(t, err, dmErrNotNullViolation)
			_, err = tx.Model(table).Data(saveList).Save()
			assertDmError(t, err, dmErrNotNullViolation)
			return nil
		})
		t.AssertNil(err)

		array, err := db.Model(table).Array("passport")
		t.AssertNil(err)
		t.Assert(gconv.Strings(array), []string{"kept"})
	})
}

// Test_DM_Records_Schema tests that a save of several records, which DM merges in a transaction of
// its own, writes the table of the schema of the model, while the configured schema has a table of
// the same name.
func Test_DM_Records_Schema(t *testing.T) {
	table := createUserTable()
	defer dropTable(table)
	schema := strings.ToUpper(dmFeatureTableName("gf_records"))
	dmFeatureExec(
		fmt.Sprintf(`CREATE SCHEMA %s`, schema),
		fmt.Sprintf(
			`CREATE TABLE %s.%s (ID NUMBER(10) NOT NULL, PASSPORT VARCHAR(45) NOT NULL, PASSWORD CHAR(32) NOT NULL, NICKNAME VARCHAR(45) NOT NULL, PRIMARY KEY (ID))`,
			schema, table,
		),
	)
	defer func() {
		_, _ = db.Exec(ctx, fmt.Sprintf(`DROP SCHEMA %s CASCADE`, schema))
	}()

	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(table).Schema(schema).Data(g.List{
			{"id": 1, "passport": "s_1", "password": "p", "nickname": "n"},
			{"id": 2, "passport": "s_2", "password": "p", "nickname": "n"},
		}).Save()
		t.AssertNil(err)

		count, err := db.Model(table).Schema(schema).Count()
		t.AssertNil(err)
		t.Assert(count, 2)
		count, err = db.Model(table).Count()
		t.AssertNil(err)
		t.Assert(count, 0)
	})
}

// Test_DM_LastInsertId_NonIntegerPrimaryKey tests inserting without the primary key into
// tables whose primary key is not an integer and is filled by the database.
func Test_DM_LastInsertId_NonIntegerPrimaryKey(t *testing.T) {
	var (
		guidTable = dmFeatureTableName("t_guid")
		rawTable  = dmFeatureTableName("t_raw")
	)
	dmFeatureExec(
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

// Test_DM_LastInsertId_NoPrimaryKey tests InsertAndGetId on a table without primary key.
func Test_DM_LastInsertId_NoPrimaryKey(t *testing.T) {
	table := dmFeatureTableName("t_nopk")
	dmFeatureExec(fmt.Sprintf(`CREATE TABLE %s (ID NUMBER(10), NAME VARCHAR2(20))`, table))
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

// Test_DM_LastInsertId_IdentityGivenBySome tests a batch of which only some records give the
// identity column, which DM inserts with IDENTITY_INSERT on: the insert succeeds, and LastInsertId
// fails with CodeNotSupported, as the ROWID of the last record does not tell which value it took.
// Without IDENTITY_INSERT, DM rejects the record giving the identity column.
// Note: the core inserts the records of different fields by separate statements in the iteration
// order of a map, so that the record without the identity column may be inserted before the other
// one is rejected.
func Test_DM_LastInsertId_IdentityGivenBySome(t *testing.T) {
	table := dmFeatureTableName("t_identity_some")
	dmFeatureExec(fmt.Sprintf(`CREATE TABLE %s (ID INT IDENTITY(1, 1) PRIMARY KEY, NAME VARCHAR(20))`, table))
	defer dropTable(table)
	var list = g.List{{"ID": 100, "NAME": "a"}, {"NAME": "b"}}

	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(table).Data(list).Insert()
		assertDmError(t, err, dmErrIdentityInsertOff)
		_, err = db.Model(table).Delete("1=1")
		t.AssertNil(err)
	})
	gtest.C(t, func(t *gtest.T) {
		err := db.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
			if _, err := tx.Exec(fmt.Sprintf(`SET IDENTITY_INSERT %s ON`, table)); err != nil {
				return err
			}
			result, err := tx.Model(table).Data(list).Insert()
			t.AssertNil(err)
			affected, err := result.RowsAffected()
			t.AssertNil(err)
			t.Assert(affected, 2)
			_, err = result.LastInsertId()
			t.Assert(gerror.Code(err), gcode.CodeNotSupported)
			_, err = tx.Exec(fmt.Sprintf(`SET IDENTITY_INSERT %s OFF`, table))
			return err
		})
		t.AssertNil(err)
		count, err := db.Model(table).Count()
		t.AssertNil(err)
		t.Assert(count, 2)
	})
}

// Test_DM_LastInsertId_Transaction tests InsertAndGetId inside a transaction, and that a rolled
// back insert still consumes its sequence value.
func Test_DM_LastInsertId_Transaction(t *testing.T) {
	table := createUserTable()
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

// Test_DM_Sequence_Nextval tests reading sequence values through raw sql, the DUAL table and
// an insert that takes the primary key from a sequence.
func Test_DM_Sequence_Nextval(t *testing.T) {
	var (
		table    = createUserTable()
		sequence = table + "_ID_SEQ"
		bigSeq   = dmFeatureTableName("s_big")
	)
	defer dropTable(table)
	dmFeatureExec(fmt.Sprintf(`CREATE SEQUENCE %s START WITH 9007199254740993`, bigSeq))
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

	// Note: Oracle reads NEXTVAL, a NUMBER without precision, as float64, losing the integers
	// beyond 2^53; DM reports NEXTVAL as BIGINT, read as int64 without loss. A cast to NUMBER(19)
	// is reported as DEC(19,0), read as its exact decimal string.
	gtest.C(t, func(t *gtest.T) {
		value, err := db.GetValue(ctx, fmt.Sprintf("SELECT CAST(%s.NEXTVAL AS NUMBER(19)) AS ID FROM dual", bigSeq))
		t.AssertNil(err)
		t.Assert(value.String(), "9007199254740993")
		t.AssertEQ(value.Val(), "9007199254740993")

		value, err = db.GetValue(ctx, fmt.Sprintf("SELECT %s.NEXTVAL FROM dual", bigSeq))
		t.AssertNil(err)
		t.AssertEQ(value.Val(), int64(9007199254740994))
	})
}

// Test_DM_Sequence_NextvalSingleRow tests reading a sequence value from DUAL with Value and One,
// which limit the query to one row.
// Note: Oracle fails these with ORA-02287, as its ROWNUM sub-query does not allow a sequence; DM
// limits the query with LIMIT, which allows it.
func Test_DM_Sequence_NextvalSingleRow(t *testing.T) {
	var (
		table    = createUserTable()
		sequence = table + "_ID_SEQ"
	)
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		value, err := db.Model("dual").Fields(sequence + ".NEXTVAL").Value()
		t.AssertNil(err)
		t.AssertEQ(value.Val(), int64(1))
	})

	gtest.C(t, func(t *gtest.T) {
		one, err := db.Model("dual").Fields(sequence + ".NEXTVAL N").One()
		t.AssertNil(err)
		t.Assert(len(one), 1)
		t.AssertEQ(one["N"].Val(), int64(2))
	})
}

// Test_DM_Sequence_NextvalSave tests a sequence value given as Raw to Save, whose MERGE takes the
// values from a sub-query, and to Insert, which writes it into the VALUES clause.
// Note: Oracle fails the Save with ORA-02287, as it does not allow a sequence in the sub-query of
// MERGE; DM allows it.
func Test_DM_Sequence_NextvalSave(t *testing.T) {
	var (
		table    = createUserTable()
		sequence = table + "_ID_SEQ"
	)
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		var previous int64
		for i := 0; i < 2; i++ {
			result, err := db.Model(table).Data(g.Map{
				"id": 1, "passport": gdb.Raw(sequence + ".NEXTVAL"), "password": "p", "nickname": "n",
			}).Save()
			t.AssertNil(err)
			n, err := result.RowsAffected()
			t.AssertNil(err)
			t.Assert(n, 1)

			passport, err := db.Model(table).Where("id", 1).Value("passport")
			t.AssertNil(err)
			current, err := db.GetValue(ctx, fmt.Sprintf("SELECT %s.CURRVAL FROM dual", sequence))
			t.AssertNil(err)
			t.AssertGT(passport.Int64(), previous)
			t.AssertLE(passport.Int64(), current.Int64())
			previous = passport.Int64()
		}

		count, err := db.Model(table).Count()
		t.AssertNil(err)
		t.Assert(count, 1)
	})

	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(table).Data(g.Map{
			"id": gdb.Raw(sequence + ".NEXTVAL"), "passport": "p", "password": "p", "nickname": "n",
		}).Insert()
		t.AssertNil(err)

		count, err := db.Model(table).Count()
		t.AssertNil(err)
		t.Assert(count, 2)
	})
}

// Test_DM_Placeholder_InStringLiteral tests that a question mark inside a string literal is
// kept as text instead of being turned into a bind placeholder.
func Test_DM_Placeholder_InStringLiteral(t *testing.T) {
	table := createInitUserTable()
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

// Test_DM_Placeholder_InComment tests that question marks inside comments do not disturb the
// binding of the real placeholders.
func Test_DM_Placeholder_InComment(t *testing.T) {
	table := createInitUserTable()
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

// Test_DM_Placeholder_DoubleQuoteInLiteral tests that a double quote inside a string literal is kept.
func Test_DM_Placeholder_DoubleQuoteInLiteral(t *testing.T) {
	table := createInitUserTable()
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

// Test_DM_Merge_OnlyConflictKey tests Save with data holding only the primary key, which leaves
// nothing to update when the row exists.
func Test_DM_Merge_OnlyConflictKey(t *testing.T) {
	table := dmFeatureTableName("t_mkey")
	dmFeatureExec(
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

// Test_DM_Merge_OnDuplicateCounter tests that a gdb.Counter in OnDuplicate adds to the value
// stored in the matched row, as MySQL does, not to the value being inserted.
func Test_DM_Merge_OnDuplicateCounter(t *testing.T) {
	table := dmFeatureTableName("t_mcnt")
	dmFeatureExec(
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

// Test_DM_Merge_NullConflictKey tests Save with a NULL primary key, which never matches, so the
// trigger fills a new id for every call.
func Test_DM_Merge_NullConflictKey(t *testing.T) {
	table := createUserTable()
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

// Test_DM_Merge_NoConflictKey tests the upsert operations on data without primary key, with and
// without OnConflict.
func Test_DM_Merge_NoConflictKey(t *testing.T) {
	table := dmFeatureTableName("t_mnok")
	dmFeatureExec(fmt.Sprintf(`CREATE TABLE %s (ID NUMBER(10), NAME VARCHAR2(20))`, table))
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
		userTable := createUserTable()
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

// Test_DM_Merge_OnConflictReplaceAndIgnore tests that Replace and InsertIgnore detect conflicts
// on the columns given by OnConflict, as the driver builds their MERGE from them.
func Test_DM_Merge_OnConflictReplaceAndIgnore(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		table := dmFeatureTableName("t_mocf")
		dmFeatureExec(
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
		table := dmFeatureTableName("t_mocf")
		dmFeatureExec(
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
		table := createInitUserTable()
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
		table := createUserTable()
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

// Test_DM_Merge_InsertIgnoreConflict tests that InsertIgnore reports no affected row and keeps
// the existing row on a conflict.
func Test_DM_Merge_InsertIgnoreConflict(t *testing.T) {
	table := createInitUserTable()
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

// Test_DM_Merge_RawExec tests a hand-written MERGE statement with bind arguments through Exec.
func Test_DM_Merge_RawExec(t *testing.T) {
	table := createInitUserTable()
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

// Test_DM_Merge_OnDuplicateExFolded tests that OnDuplicateEx leaves out a column named in
// another letter case than in OnDuplicate, as DM resolves both to the same column.
func Test_DM_Merge_OnDuplicateExFolded(t *testing.T) {
	table := createInitUserTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(table).OnDuplicate("nickname,create_time").OnDuplicateEx("NICKNAME").Data(g.Map{
			"id": 1, "passport": "user_1", "password": "pass_1", "nickname": "name_new", "create_time": CreateTime,
		}).Save()
		t.AssertNil(err)

		one, err := db.Model(table).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["CREATE_TIME"], CreateTime)
		t.Assert(one["NICKNAME"], "name_1")
	})

	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(table).OnDuplicate(g.Map{"nickname": "nickname", "create_time": "create_time"}).
			OnDuplicateEx("Nickname").Data(g.Map{
			"id": 2, "passport": "user_2", "password": "pass_2", "nickname": "name_new", "create_time": CreateTime,
		}).Save()
		t.AssertNil(err)

		one, err := db.Model(table).Where("id", 2).One()
		t.AssertNil(err)
		t.Assert(one["CREATE_TIME"], CreateTime)
		t.Assert(one["NICKNAME"], "name_2")
	})
}

// Test_DM_Merge_CompositeKeyOrder tests that the MERGE statement of a table with a composite
// primary key matches the key columns in the same order every time.
func Test_DM_Merge_CompositeKeyOrder(t *testing.T) {
	table := dmFeatureTableName("t_mck")
	dmFeatureExec(fmt.Sprintf(
		`CREATE TABLE %s (K1 NUMBER(10), K2 NUMBER(10), K3 NUMBER(10), NAME VARCHAR2(20), PRIMARY KEY (K3, K1, K2))`,
		table,
	))
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		for i := 0; i < 20; i++ {
			sqlArray, err := gdb.CatchSQL(ctx, func(ctx context.Context) error {
				_, err := db.Model(table).Ctx(ctx).Data(g.Map{"k1": 1, "k2": 2, "k3": 3, "name": "a"}).Save()
				return err
			})
			t.AssertNil(err)
			mergeSql := sqlArray[len(sqlArray)-1]
			t.Assert(gstr.HasPrefix(mergeSql, "MERGE INTO"), true)
			t.Assert(gstr.Contains(mergeSql, "ON (T1.K1 = T2.K1 AND T1.K2 = T2.K2 AND T1.K3 = T2.K3)"), true)
		}
	})

	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(table).Data(g.Map{"k1": 1, "k2": 2, "k3": 3, "name": "a"}).Save()
		t.AssertNil(err)
		_, err = db.Model(table).Data(g.Map{"k1": 1, "k2": 2, "k3": 3, "name": "b"}).Save()
		t.AssertNil(err)

		all, err := db.Model(table).All()
		t.AssertNil(err)
		t.Assert(len(all), 1)
		t.Assert(all[0]["NAME"], "b")
	})
}

// Test_DM_TableFields_TypeStrings tests the field types that TableFields reports for the Oracle
// compatible types without a MySQL counterpart, as DM reports them, and that a NUMBER without
// precision keeps its decimals.
func Test_DM_TableFields_TypeStrings(t *testing.T) {
	table := dmFeatureTableName("t_types")
	dmFeatureExec(fmt.Sprintf(`CREATE TABLE %s (
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
		// Note: DM reports BINARY_DOUBLE as DOUBLE, TIMESTAMP WITH TIME ZONE as DATETIME WITH
		// TIME ZONE and RAW as VARBINARY, the length in bytes of the character types, which is
		// 6 for VARCHAR2(3 CHAR) with the GB18030 charset of the test server, and the storage
		// size of the other types.
		expect := map[string]string{
			"ID":  "NUMBER(10,0)",
			"N":   "NUMBER(22)",
			"I":   "INTEGER(4)",
			"F":   "FLOAT(8)",
			"BD":  "DOUBLE(8)",
			"TZ":  "DATETIME WITH TIME ZONE(10)",
			"TL":  "TIMESTAMP WITH LOCAL TIME ZONE(8)",
			"IDS": "INTERVAL DAY TO SECOND(2,6)",
			"IYM": "INTERVAL YEAR TO MONTH(2,6)",
			"NV":  "NVARCHAR2(20)",
			"NC":  "NCHAR(4)",
			"VB":  "VARCHAR2(6)",
			"VC":  "VARCHAR2(6)",
			"R":   "VARBINARY(8)",
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

		// Note: DM writes the leading zero of a fraction that Oracle leaves out.
		array, err := db.Model(table).Order("id").Array("TO_CHAR(n)")
		t.AssertNil(err)
		t.Assert(array, g.Slice{"3.14159", "12345678901234567890", "0.000000000000000000000000000001"})

		// Note: Oracle reads the value of a NUMBER column without precision as float64; on DM it
		// is the exact decimal string, like a DECIMAL value on MySQL.
		array, err = db.Model(table).Order("id").Array("n")
		t.AssertNil(err)
		t.Assert(array[0].String(), "3.14159")
		t.Assert(array[1].String(), "12345678901234567890")
		t.Assert(array[2].String(), "0.000000000000000000000000000001")
	})
}

// Test_DM_TableFields_CommentAndDefault tests what TableFields reports for a column that has a
// comment and a default value.
func Test_DM_TableFields_CommentAndDefault(t *testing.T) {
	table := dmFeatureTableName("t_cmt")
	dmFeatureExec(
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

	// Note: the Oracle driver reports no comment and no default; the DM driver reports the
	// comment, and the default as the expression DM stores, which quotes a string.
	gtest.C(t, func(t *gtest.T) {
		fields, err := db.TableFields(ctx, table)
		t.AssertNil(err)
		t.Assert(fields["NAME"].Type, "VARCHAR2(20)")
		t.Assert(fields["NAME"].Comment, "the name")
		t.Assert(fields["NAME"].Default, "'dflt'")
		t.Assert(fields["NAME"].Extra, "")
		t.Assert(fields["ID"].Comment, "")
		t.Assert(fields["ID"].Default, nil)
	})
}

// Test_DM_TableFields_View tests TableFields and queries on a view, and that Tables lists no view.
func Test_DM_TableFields_View(t *testing.T) {
	var (
		table = createInitUserTable()
		view  = dmFeatureTableName("v_user")
	)
	defer dropTable(table)
	dmFeatureExec(fmt.Sprintf(`CREATE VIEW %s AS SELECT ID, NICKNAME FROM %s WHERE ID > 8`, view, table))
	defer db.Exec(ctx, "DROP VIEW "+view)

	gtest.C(t, func(t *gtest.T) {
		fields, err := db.TableFields(ctx, view)
		t.AssertNil(err)
		t.Assert(len(fields), 2)
		t.Assert(fields["ID"].Type, "NUMBER(10,0)")
		t.Assert(fields["ID"].Key, "")
		t.Assert(fields["ID"].Null, false)
		t.Assert(fields["NICKNAME"].Type, "VARCHAR(45)")
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

// Test_DM_Model_SchemaQualifiedTable tests querying and changing a table named with its schema,
// the schema of the connection.
func Test_DM_Model_SchemaQualifiedTable(t *testing.T) {
	var (
		table  = createInitUserTable()
		schema = TestDBName
	)
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		one, err := db.Model(schema+"."+table).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["PASSPORT"], "user_1")

		one, err = db.Model(strings.ToLower(schema)+"."+table).Fields("id, nickname").Where("id", 2).One()
		t.AssertNil(err)
		t.Assert(len(one), 2)
		t.Assert(one["NICKNAME"], "name_2")

		all, err := db.Model(schema+"."+table+" u").Where("u.id > ?", 8).Order("u.id").All()
		t.AssertNil(err)
		t.Assert(gconv.Ints(all.Array("ID")), []int{9, 10})
	})

	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(schema+"."+table).Where("id", 4).Delete()
		t.AssertNil(err)
		n, _ := result.RowsAffected()
		t.Assert(n, 1)

		_, err = db.Model(schema + "." + table).Data(g.Map{
			"id": 11, "passport": "user_11", "password": "p", "nickname": "name_11",
		}).Insert()
		t.AssertNil(err)
		count, err := db.Model(table).Count()
		t.AssertNil(err)
		t.Assert(count, TableSize)
	})
}

// Test_DM_Model_SchemaQualifiedFieldTypes tests the operations that need the field types of a
// table named with its schema.
func Test_DM_Model_SchemaQualifiedFieldTypes(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		table := createUserTable()
		defer dropTable(table)

		fields, err := db.TableFields(ctx, TestDBName+"."+table)
		t.AssertNil(err)
		t.Assert(len(fields), 6)
	})

	gtest.C(t, func(t *gtest.T) {
		table := createUserTable()
		defer dropTable(table)

		id, err := db.Model(TestDBName + "." + table).Data(g.Map{"passport": "p", "password": "p", "nickname": "n"}).InsertAndGetId()
		t.AssertNil(err)
		value, err := db.Model(table).Where("passport", "p").Value("id")
		t.AssertNil(err)
		t.Assert(value, 1)
		t.Assert(id, 1)
	})

	gtest.C(t, func(t *gtest.T) {
		table := createInitUserTable()
		defer dropTable(table)

		_, err := db.Model(TestDBName+"."+table).Data("nickname", "schema").Where("id", 3).Update()
		t.AssertNil(err)

		value, err := db.Model(table).Where("id", 3).Value("nickname")
		t.AssertNil(err)
		t.Assert(value, "schema")
	})

	gtest.C(t, func(t *gtest.T) {
		table := dmFeatureTableName("t_sclob")
		dmFeatureExec(fmt.Sprintf(`CREATE TABLE %s (ID NUMBER(10) NOT NULL, C CLOB, PRIMARY KEY (ID))`, table))
		defer dropTable(table)

		text := strings.Repeat("x", 5000)
		_, err := db.Model(table).Data(g.Map{"id": 1, "c": text}).Insert()
		t.AssertNil(err)

		_, err = db.Model(TestDBName + "." + table).Data(g.Map{"id": 2, "c": text}).Insert()
		t.AssertNil(err)

		value, err := db.Model(table).Where("id", 2).Value("c")
		t.AssertNil(err)
		t.Assert(len(value.String()), 5000)
	})
}

// Test_DM_Model_OtherSchemaTable tests the model operations that need the fields of a table
// qualified with a schema other than the schema of the connection, while the schema of the
// connection has a table of the same name with other columns.
// Note: Oracle has a schema per user; DM creates a schema of the connected user with CREATE SCHEMA.
func Test_DM_Model_OtherSchemaTable(t *testing.T) {
	var (
		schema      = strings.ToUpper(dmFeatureTableName("gf_schema"))
		table       = dmFeatureTableName("t_schema")
		schemaTable = schema + "." + table
	)
	dmFeatureExec(
		fmt.Sprintf(`CREATE SCHEMA %s`, schema),
		fmt.Sprintf(`CREATE TABLE %s (ID NUMBER(10) NOT NULL, PASSPORT VARCHAR2(45), NICKNAME VARCHAR2(45), PRIMARY KEY (ID))`, schemaTable),
		fmt.Sprintf(`CREATE TABLE %s (ID NUMBER(10) NOT NULL, NAME VARCHAR2(45), PRIMARY KEY (ID))`, table),
	)
	defer func() {
		dropTable(table)
		dmFeatureExec(fmt.Sprintf(`DROP SCHEMA %s CASCADE`, schema))
	}()

	gtest.C(t, func(t *gtest.T) {
		t.Assert(db.GetTableNameForFields(schemaTable), schemaTable)
		t.Assert(db.GetTableNameForFields(schemaTable+" u"), schemaTable)
		t.Assert(db.GetTableNameForFields(schemaTable+" AS u, "+table+" t"), schemaTable)
		t.Assert(db.GetTableNameForFields(table+" u"), table)
		t.Assert(db.GetTableNameForFields(","+schemaTable), schemaTable)

		fields, err := db.Model(schemaTable + " u").TableFields(schemaTable + " u")
		t.AssertNil(err)
		t.Assert(len(fields), 3)
		t.AssertNE(fields["PASSPORT"], nil)

		has, err := db.GetCore().HasField(ctx, schemaTable, "nickname")
		t.AssertNil(err)
		t.Assert(has, true)
	})

	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(schemaTable).Data(g.Map{"id": 1, "passport": "user_1", "nickname": "name_1"}).Insert()
		t.AssertNil(err)

		result, err := db.Model(schemaTable).Data("nickname", "owner").Where("id", 1).Update()
		t.AssertNil(err)
		n, _ := result.RowsAffected()
		t.Assert(n, 1)

		one, err := db.GetOne(ctx, fmt.Sprintf(`SELECT PASSPORT, NICKNAME FROM %s WHERE ID = 1`, schemaTable))
		t.AssertNil(err)
		t.Assert(one["PASSPORT"], "user_1")
		t.Assert(one["NICKNAME"], "owner")

		one, err = db.Model(schemaTable).FieldsEx("passport").Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(len(one), 2)
		t.Assert(one["NICKNAME"], "owner")

		count, err := db.Model(table).Count()
		t.AssertNil(err)
		t.Assert(count, 0)
	})
}

// Test_DM_Model_QuotedSchemaTable tests a table whose name is quoted in lower case and qualified
// with its schema, which DM looks up as it is written.
func Test_DM_Model_QuotedSchemaTable(t *testing.T) {
	var (
		name  = dmFeatureTableName("t_quoted")
		table = fmt.Sprintf(`"%s"."%s"`, TestDBName, name)
	)
	dmFeatureExec(fmt.Sprintf(`CREATE TABLE %s (ID NUMBER(10) NOT NULL, NICKNAME VARCHAR2(45), PRIMARY KEY (ID))`, table))
	defer dmFeatureExec(fmt.Sprintf(`DROP TABLE %s`, table))

	gtest.C(t, func(t *gtest.T) {
		t.Assert(db.GetTableNameForFields(table+" u"), table)

		fields, err := db.TableFields(ctx, table)
		t.AssertNil(err)
		t.Assert(len(fields), 2)
		t.Assert(fields["ID"].Key, "PRI")
		t.Assert(fields["NICKNAME"].Type, "VARCHAR2(45)")
	})

	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(table).Data(g.Map{"id": 1, "nickname": "name_1"}).Insert()
		t.AssertNil(err)

		one, err := db.Model(table).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["NICKNAME"], "name_1")
	})

	gtest.C(t, func(t *gtest.T) {
		quoted := fmt.Sprintf(`"%s"`, name)
		t.Assert(db.GetTableNameForFields(quoted+" u"), quoted)

		fields, err := db.TableFields(ctx, quoted)
		t.AssertNil(err)
		t.Assert(len(fields), 2)
		t.Assert(fields["ID"].Key, "PRI")

		_, err = db.Model(quoted).Data(g.Map{"id": 2, "nickname": "name_2"}).Insert()
		t.AssertNil(err)

		one, err := db.Model(quoted).Where("id", 2).One()
		t.AssertNil(err)
		t.Assert(one["NICKNAME"], "name_2")
	})
}

// Test_DM_Model_AliasedWrite tests writing through a model whose table has an alias, given by
// As, by Model or in the table name, which is the target of INSERT and MERGE without its alias and
// whose field types still convert the values.
func Test_DM_Model_AliasedWrite(t *testing.T) {
	table := dmFeatureTableName("t_alias")
	// Note: DM DATE is date-only, so the time of day is kept in a DATETIME column.
	dmFeatureExec(fmt.Sprintf(
		`CREATE TABLE %s (ID NUMBER(10) NOT NULL, NICKNAME VARCHAR2(45), INFO VARCHAR2(200), DT DATETIME, PRIMARY KEY (ID))`,
		table,
	))
	defer dropTable(table)

	var (
		models = []func() *gdb.Model{
			func() *gdb.Model { return db.Model(table).As("a") },
			func() *gdb.Model { return db.Model(table, "a") },
			func() *gdb.Model { return db.Model(table + " AS a") },
			func() *gdb.Model { return db.Model(table + " a") },
		}
		dt = time.Date(2024, 1, 15, 4, 30, 45, 0, time.UTC)
	)
	assertRow := func(t *gtest.T, id int, nickname string, info string, dt time.Time) {
		one, err := db.GetOne(ctx, fmt.Sprintf(
			"SELECT NICKNAME, INFO, TO_CHAR(DT, 'YYYY-MM-DD HH24:MI:SS') AS DT FROM %s WHERE ID = ?", table,
		), id)
		t.AssertNil(err)
		t.Assert(one["NICKNAME"], nickname)
		t.Assert(one["INFO"], info)
		t.Assert(one["DT"], dt.In(time.Local).Format("2006-01-02 15:04:05"))
	}

	for i, model := range models {
		id := (i + 1) * 10

		gtest.C(t, func(t *gtest.T) {
			_, err := model().Data(g.Map{"id": id, "nickname": "insert", "info": g.Map{"k": "insert"}, "dt": dt}).Insert()
			t.AssertNil(err)
			assertRow(t, id, "insert", `{"k":"insert"}`, dt)
		})

		gtest.C(t, func(t *gtest.T) {
			_, err := model().Data(g.Map{
				"nickname": "update", "info": g.Map{"k": "update"}, "dt": dt.Add(time.Hour),
			}).Where("id", id).Update()
			t.AssertNil(err)
			assertRow(t, id, "update", `{"k":"update"}`, dt.Add(time.Hour))
		})

		gtest.C(t, func(t *gtest.T) {
			_, err := model().Data(g.Map{"id": id, "nickname": "save", "info": g.Map{"k": "save"}, "dt": dt}).Save()
			t.AssertNil(err)
			assertRow(t, id, "save", `{"k":"save"}`, dt)

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
	}

	gtest.C(t, func(t *gtest.T) {
		count, err := db.Model(table).Count()
		t.AssertNil(err)
		t.Assert(count, len(models)*3)
	})
}

// Test_DM_Lob_LargeRoundTrip tests CLOB and BLOB values of more than 4000 and 32767 bytes, and of
// multibyte characters, through insert, update and select.
func Test_DM_Lob_LargeRoundTrip(t *testing.T) {
	table := dmFeatureTableName("t_lob")
	dmFeatureExec(fmt.Sprintf(`CREATE TABLE %s (ID NUMBER(10) NOT NULL, C CLOB, B BLOB, PRIMARY KEY (ID))`, table))
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

// Test_DM_Lob_ComparisonRejected tests that DM rejects a CLOB in comparisons (-6105), ORDER BY,
// DISTINCT and GROUP BY (-2685), while LIKE and the DBMS_LOB functions work.
func Test_DM_Lob_ComparisonRejected(t *testing.T) {
	table := dmFeatureTableName("t_lobcmp")
	dmFeatureExec(fmt.Sprintf(`CREATE TABLE %s (ID NUMBER(10) NOT NULL, C CLOB, PRIMARY KEY (ID))`, table))
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(table).Data(g.List{
			{"id": 1, "c": "abc"},
			{"id": 2, "c": strings.Repeat("x", 5000) + "needle"},
		}).Insert()
		t.AssertNil(err)

		_, err = db.Model(table).Where("c", "abc").Count()
		assertDmError(t, err, dmErrTypeMismatch)
		_, err = db.Model(table).Order("c").All()
		assertDmError(t, err, dmErrLobCompare)
		_, err = db.Model(table).Fields("DISTINCT c").All()
		assertDmError(t, err, dmErrLobCompare)
		_, err = db.Model(table).Fields("c").Group("c").All()
		assertDmError(t, err, dmErrLobCompare)
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

// Test_DM_Varchar2_ByteAndCharSemantics tests the length limits of VARCHAR2 in bytes and in
// characters, and NVARCHAR2 with characters beyond the basic plane.
// Note: the GB18030 charset of the test server stores a Chinese character in 2 bytes, where the
// AL32UTF8 charset of Oracle takes 3, and DM rejects a bound value beyond the byte length of the
// column with -5403 and a value beyond the character length with -6169, where Oracle has ORA-12899.
func Test_DM_Varchar2_ByteAndCharSemantics(t *testing.T) {
	table := dmFeatureTableName("t_vchar")
	dmFeatureExec(fmt.Sprintf(
		`CREATE TABLE %s (ID NUMBER(10) NOT NULL, VB VARCHAR2(6 BYTE), VC VARCHAR2(3 CHAR), NV NVARCHAR2(10), PRIMARY KEY (ID))`,
		table,
	))
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(table).Data(g.Map{"id": 1, "vb": "中文字中"}).Insert()
		assertDmError(t, err, dmErrParamIncompatible)

		_, err = db.Model(table).Data(g.Map{"id": 2, "vc": "abcd"}).Insert()
		assertDmError(t, err, dmErrColumnTooLong)

		count, err := db.Model(table).Count()
		t.AssertNil(err)
		t.Assert(count, 0)
	})

	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(table).Data(g.Map{"id": 3, "vb": "中文字", "vc": "中文字", "nv": "中文字符串😀"}).Insert()
		t.AssertNil(err)

		one, err := db.Model(table).Where("id", 3).One()
		t.AssertNil(err)
		t.Assert(one["VB"], "中文字")
		t.Assert(one["VC"], "中文字")
		t.Assert(one["NV"], "中文字符串😀")

		value, err := db.Model(table).Where("nv", "中文字符串😀").Value("id")
		t.AssertNil(err)
		t.Assert(value, 3)

		lengths, err := db.GetOne(ctx, fmt.Sprintf(
			"SELECT LENGTHB(VB) VB, LENGTHB(VC) VC, LENGTHB(NV) NV FROM %s WHERE ID = 3", table,
		))
		t.AssertNil(err)
		t.Assert(lengths, g.Map{"VB": 6, "VC": 6, "NV": 14})
	})
}

// Test_DM_WhereIn_Over1000 tests an IN list of more than 1000 values, which gf does not split.
// Note: Oracle limits an IN list to 1000 values (ORA-01795); DM does not.
func Test_DM_WhereIn_Over1000(t *testing.T) {
	table := createInitUserTable()
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
		count, err := db.Model(table).WhereIn("id", ids).Count()
		t.AssertNil(err)
		t.Assert(count, TableSize)

		count, err = db.Model(table).Where("id", ids).Count()
		t.AssertNil(err)
		t.Assert(count, TableSize)

		count, err = db.Model(table).WhereNotIn("id", ids).Count()
		t.AssertNil(err)
		t.Assert(count, 0)
	})
}

// Test_DM_EmptyString_NotNull tests that DM stores an empty string as it is, not as NULL, and how
// conditions on an empty string behave.
// Note: Oracle stores an empty string as NULL, so that it fails a NOT NULL column (ORA-01400);
// DM keeps it like MySQL.
func Test_DM_EmptyString_NotNull(t *testing.T) {
	table := createInitUserTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(table).Data(g.Map{
			"id": 11, "passport": "", "password": "p", "nickname": "n", "create_time": CreateTime,
		}).Insert()
		t.AssertNil(err)

		value, err := db.Model(table).Where("id", 11).Value("passport")
		t.AssertNil(err)
		t.Assert(value.IsNil(), false)
		t.Assert(value.String(), "")
	})

	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(table).Data("create_time", "").Where("id", 1).Update()
		t.AssertNil(err)

		value, err := db.Model(table).Where("id", 1).Value("create_time")
		t.AssertNil(err)
		t.Assert(value.IsNil(), false)
		t.Assert(value.String(), "")

		count, err := db.Model(table).WhereNull("create_time").Count()
		t.AssertNil(err)
		t.Assert(count, 0)

		count, err = db.Model(table).Where("create_time", "").Count()
		t.AssertNil(err)
		t.Assert(count, 1)

		count, err = db.Model(table).WhereNot("create_time", "").Count()
		t.AssertNil(err)
		t.Assert(count, TableSize)

		count, err = db.Model(table).OmitEmptyWhere().Where("create_time", "").Count()
		t.AssertNil(err)
		t.Assert(count, TableSize+1)
	})
}

// dmFeatureValuerString is a string type implementing driver.Valuer.
type dmFeatureValuerString string

// Value implements driver.Valuer, returning the passport of the user numbered by the string.
func (s dmFeatureValuerString) Value() (driver.Value, error) {
	return "user_" + string(s), nil
}

// Test_DM_NamedString_Arguments tests values of a defined string type and of a driver.Valuer as
// arguments of Where, of Data and of a raw statement, including named arguments.
func Test_DM_NamedString_Arguments(t *testing.T) {
	type rvStr string
	table := createInitUserTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		one, err := db.Model(table).Where("passport", rvStr("user_2")).One()
		t.AssertNil(err)
		t.Assert(one["ID"], 2)

		passport := rvStr("user_3")
		one, err = db.Model(table).Where("passport", &passport).One()
		t.AssertNil(err)
		t.Assert(one["ID"], 3)

		count, err := db.Model(table).Where("passport IN(?)", g.Slice{rvStr("user_1"), rvStr("user_4")}).Count()
		t.AssertNil(err)
		t.Assert(count, 2)
	})

	gtest.C(t, func(t *gtest.T) {
		_, err := db.Exec(ctx, fmt.Sprintf("UPDATE %s SET nickname = ? WHERE id = ?", table), rvStr("y"), 5)
		t.AssertNil(err)

		_, err = db.Model(table).Data("nickname = ?", rvStr("z")).Where("id", 6).Update()
		t.AssertNil(err)

		_, err = db.Model(table).Data(g.Map{"nickname": rvStr("w")}).Where("id", 7).Update()
		t.AssertNil(err)

		_, err = db.Exec(
			ctx, fmt.Sprintf("UPDATE %s SET nickname = :nickname WHERE id = :id", table),
			sql.Named("nickname", rvStr("v")), sql.Named("id", 9),
		)
		t.AssertNil(err)

		array, err := db.Model(table).WhereIn("id", g.Slice{5, 6, 7, 9}).Order("id").Array("nickname")
		t.AssertNil(err)
		t.Assert(array, g.Slice{"y", "z", "w", "v"})
	})

	gtest.C(t, func(t *gtest.T) {
		one, err := db.Model(table).Where("passport", dmFeatureValuerString("8")).One()
		t.AssertNil(err)
		t.Assert(one["ID"], 8)
	})
}

// Test_DM_Double_BeyondNumberRange tests floats beyond the range of NUMBER, from 1e-130 to 1e126
// exclusive, in DOUBLE and REAL columns, as arguments of Where, of Data in string form and of a
// raw statement.
// Note: DM has DOUBLE and REAL for the BINARY_DOUBLE and BINARY_FLOAT of Oracle, whose driver
// encodes a float as NUMBER, and the underlying driver binds a float as DOUBLE.
func Test_DM_Double_BeyondNumberRange(t *testing.T) {
	table := dmFeatureTableName("t_bdrange")
	dmFeatureExec(
		fmt.Sprintf(`CREATE TABLE %s (ID NUMBER(10) NOT NULL, D DOUBLE, F REAL, PRIMARY KEY (ID))`, table),
		fmt.Sprintf(`INSERT INTO %s (ID, D, F) VALUES (1, 0, 0)`, table),
	)
	defer dropTable(table)

	values := []float64{
		3e200, -3e200, 1e126, -1e126, math.Nextafter(1e126, 0), 1e-130, -1e-130, math.Nextafter(1e-130, 0),
		-5e-324, 0.1, math.Pi, math.MaxFloat64, math.SmallestNonzeroFloat64,
	}

	gtest.C(t, func(t *gtest.T) {
		for _, value := range values {
			_, err := db.Model(table).Data("d = ?", value).Where("id", 1).Update()
			t.AssertNil(err)

			stored, err := db.Model(table).Where("id", 1).Value("d")
			t.AssertNil(err)
			t.Assert(stored.Float64(), value)
		}
	})

	gtest.C(t, func(t *gtest.T) {
		for _, value := range values {
			_, err := db.Exec(ctx, fmt.Sprintf("UPDATE %s SET d = ? WHERE id = ?", table), value, 1)
			t.AssertNil(err)

			stored, err := db.Model(table).Where("id", 1).Value("d")
			t.AssertNil(err)
			t.Assert(stored.Float64(), value)
		}
	})

	// Note: DM has no BINARY_DOUBLE literal like 3e+200d, so the literal is a plain number.
	gtest.C(t, func(t *gtest.T) {
		for _, value := range values {
			dmFeatureExec(fmt.Sprintf("UPDATE %s SET d = %v WHERE id = 1", table, value))

			count, err := db.Model(table).Where("d", value).Count()
			t.AssertNil(err)
			t.Assert(count, 1)
		}
	})

	// Note: Oracle stores the infinities in BINARY_DOUBLE; DM rejects them as an overflow, as
	// MySQL rejects them in DOUBLE.
	gtest.C(t, func(t *gtest.T) {
		for _, value := range []float64{math.Inf(1), math.Inf(-1)} {
			_, err := db.Model(table).Data(g.Map{"d": value}).Where("id", 1).Update()
			assertDmError(t, err, dmErrNumericOverflow)

			_, err = db.Model(table).Data(g.Map{"f": float32(value)}).Where("id", 1).Update()
			assertDmError(t, err, dmErrNumericOverflow)
		}
	})

	// Note: DM stores NaN, which equals no value in a condition, NaN included, where Oracle
	// matches NaN with NaN.
	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(table).Data(g.Map{"d": math.NaN(), "f": float32(math.NaN())}).Where("id", 1).Update()
		t.AssertNil(err)

		one, err := db.Model(table).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["D"].String(), "NaN")
		t.Assert(one["F"].String(), "NaN")

		count, err := db.Model(table).Where("d", math.NaN()).Count()
		t.AssertNil(err)
		t.Assert(count, 0)
		count, err = db.Model(table).Where("f", float32(math.NaN())).Count()
		t.AssertNil(err)
		t.Assert(count, 0)
	})

	// Note: Oracle fails both with ORA-01426; the underlying driver rejects the value for the
	// NUMBER(10) column (6092), and the condition matches no row.
	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(table).Data("id = ?", 1e200).Where("id", 1).Update()
		assertDmError(t, err, dmErrDataOutOfRange)

		count, err := db.Model(table).Where("id", 1e200).Count()
		t.AssertNil(err)
		t.Assert(count, 0)
	})

	gtest.C(t, func(t *gtest.T) {
		value := sql.NullFloat64{Float64: 2e200, Valid: true}
		_, err := db.Exec(ctx, fmt.Sprintf("UPDATE %s SET d = ? WHERE id = ?", table), value, 1)
		t.AssertNil(err)
		stored, err := db.Model(table).Where("id", 1).Value("d")
		t.AssertNil(err)
		t.Assert(stored.Float64(), 2e200)

		count, err := db.Model(table).Where("d", value).Count()
		t.AssertNil(err)
		t.Assert(count, 1)

		value.Float64 = -2e200
		_, err = db.Model(table).Data("d = ?", value).Where("id", 1).Update()
		t.AssertNil(err)
		stored, err = db.Model(table).Where("id", 1).Value("d")
		t.AssertNil(err)
		t.Assert(stored.Float64(), -2e200)

		value.Float64 = 2e-200
		_, err = db.Model(table).Data(g.Map{"d": value}).Where("id", 1).Update()
		t.AssertNil(err)
		stored, err = db.Model(table).Where("id", 1).Value("d")
		t.AssertNil(err)
		t.Assert(stored.Float64(), 2e-200)

		var null *sql.NullFloat64
		_, err = db.Exec(ctx, fmt.Sprintf("UPDATE %s SET d = ? WHERE id = ?", table), null, 1)
		t.AssertNil(err)
		stored, err = db.Model(table).Where("id", 1).Value("d")
		t.AssertNil(err)
		t.Assert(stored.IsNil(), true)
	})

	gtest.C(t, func(t *gtest.T) {
		_, err := db.Exec(
			ctx, fmt.Sprintf("UPDATE %s SET d = :d WHERE id = :id", table), sql.Named("d", 4e200), sql.Named("id", 1),
		)
		t.AssertNil(err)
		stored, err := db.Model(table).Where("id", 1).Value("d")
		t.AssertNil(err)
		t.Assert(stored.Float64(), 4e200)

		count, err := db.GetCount(ctx, fmt.Sprintf("SELECT COUNT(1) FROM %s WHERE d = :d", table), sql.Named("d", 4e200))
		t.AssertNil(err)
		t.Assert(count, 1)
	})
}

// Test_DM_Number_ResultType tests the Go type of the values of NUMBER columns, which follows the
// column and not the value: int64 for a column without decimal places, or uint64 and then the
// decimal string for a value beyond the range of int64, and the exact decimal string for a column
// with decimal places, like a DECIMAL column on MySQL.
func Test_DM_Number_ResultType(t *testing.T) {
	table := dmFeatureTableName("t_numtype")
	dmFeatureExec(fmt.Sprintf(`CREATE TABLE %s (
    ID NUMBER(10) NOT NULL,
    AMOUNT NUMBER(20,4),
    BIG NUMBER(38,0),
    ANY_NUM NUMBER,
    PRIMARY KEY (ID)
)`, table))
	defer dropTable(table)

	_, err := db.Model(table).Data(g.List{
		{"id": 1, "amount": "1", "big": "9223372036854775807", "any_num": "1"},
		{"id": 2, "amount": "1.5", "big": "9223372036854775808", "any_num": "1.5"},
		{"id": 3, "amount": "1234567890123456.1234", "big": "99999999999999999999999999999999999999", "any_num": "12345678901234567890.5"},
		{"id": 4, "amount": "-2", "big": "-9223372036854775809", "any_num": "-2"},
	}).Insert()
	gtest.AssertNil(err)

	gtest.C(t, func(t *gtest.T) {
		all, err := db.Model(table).Order("id").All()
		t.AssertNil(err)
		t.Assert(len(all), 4)
		// Note: Oracle reads NUMBER(p,s>0) and NUMBER without precision as float64; on DM they
		// are exact. DM reports a NUMBER without precision to the underlying driver as NUMBER(38,0),
		// so only its decimal text is asserted.
		for _, record := range all {
			t.Assert(fmt.Sprintf("%T", record["ID"].Val()), "int64")
			t.Assert(fmt.Sprintf("%T", record["AMOUNT"].Val()), "string")
		}
		t.Assert(all.Array("AMOUNT"), g.Slice{"1.0000", "1.5000", "1234567890123456.1234", "-2.0000"})
		t.Assert(gconv.Strings(all.Array("ANY_NUM")), []string{"1", "1.5", "12345678901234567890.5", "-2"})
	})

	gtest.C(t, func(t *gtest.T) {
		array, err := db.Model(table).Order("id").Array("big")
		t.AssertNil(err)
		t.AssertEQ(array[0].Val(), int64(9223372036854775807))
		t.AssertEQ(array[1].Val(), uint64(9223372036854775808))
		t.AssertEQ(array[2].Val(), "99999999999999999999999999999999999999")
		t.AssertEQ(array[3].Val(), "-9223372036854775809")
	})

	gtest.C(t, func(t *gtest.T) {
		// Note: Oracle reports SUM of a NUMBER(10) as a NUMBER without precision, read as float64;
		// DM reports it as DEC(38,0), an exact number read as its text as DECIMAL is on MySQL.
		value, err := db.Model(table).Where("id", 2).Value("SUM(id)")
		t.AssertNil(err)
		t.AssertEQ(value.Val(), "2")

		count, err := db.Model(table).Count()
		t.AssertNil(err)
		t.Assert(count, 4)
	})
}

// Test_DM_Extra_OptionParsing tests how the Extra option reaches the underlying driver: an
// option is passed on and applied, while a malformed option, an unknown option and an invalid
// value are ignored, the default value of the option being kept.
// Note: go-ora rejects an unknown option or an invalid value when it connects; the DM driver
// ignores them.
func Test_DM_Extra_OptionParsing(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		for _, extra := range []string{
			"columnNameCase", "a=b=c&columnNameCase", "columnNameCase=bad",
			"FOO=bar", "connectTimeout=abc", "LobMode=bad",
		} {
			extraDb := dmFeatureOpen(extra, 0)
			one, err := extraDb.GetOne(ctx, "SELECT 1 AS ID FROM dual")
			extraDb.Close(ctx)
			t.AssertNil(err)
			t.Assert(one, g.Map{"ID": 1})
		}
	})

	gtest.C(t, func(t *gtest.T) {
		table := createInitUserTable()
		defer dropTable(table)

		extraDb := dmFeatureOpen("columnNameCase=lower", 0)
		defer extraDb.Close(ctx)
		one, err := extraDb.GetOne(ctx, "SELECT 1 AS ID FROM dual")
		t.AssertNil(err)
		t.Assert(one, g.Map{"id": 1})

		one, err = extraDb.Model(table).Fields("id, nickname").Where("id", 2).One()
		t.AssertNil(err)
		t.Assert(one, g.Map{"id": 2, "nickname": "name_2"})
	})

	gtest.C(t, func(t *gtest.T) {
		table := dmFeatureTableName("t_lobpre")
		dmFeatureExec(fmt.Sprintf(`CREATE TABLE %s (ID NUMBER(10) NOT NULL, C CLOB, B BLOB, PRIMARY KEY (ID))`, table))
		defer dropTable(table)

		text := strings.Repeat("中文", 3000)
		_, err := db.Model(table).Data(g.Map{"id": 1, "c": text, "b": []byte(text)}).Insert()
		t.AssertNil(err)

		extraDb := dmFeatureOpen("LobMode=2&connectTimeout=5000", 0)
		defer extraDb.Close(ctx)
		one, err := extraDb.Model(table).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["C"].String() == text, true)
		t.Assert(one["B"].Bytes(), []byte(text))
	})
}

// Test_DM_Nls_LanguageAndDateFormat tests that the NLS parameters a session sets apply to that session
// only, and that values are written and matched correctly under a date format other than the
// default.
// Note: DM has no LANGUAGE and TERRITORY connection options like go-ora, nor NLS_TERRITORY and
// NLS_NUMERIC_CHARACTERS, so the session sets NLS_DATE_LANGUAGE and the date formats.
func Test_DM_Nls_LanguageAndDateFormat(t *testing.T) {
	table := dmFeatureTableName("t_nls")
	dmFeatureExec(fmt.Sprintf(
		`CREATE TABLE %s (ID NUMBER(10) NOT NULL, D DATE, TM TIME, DT DATETIME, AMOUNT NUMBER(18,2), RATIO DOUBLE, PRIMARY KEY (ID))`, table,
	))
	defer dropTable(table)

	// Note: the pool has one connection, so that every statement runs in the session that sets
	// the NLS parameters.
	sessionDb := dmFeatureOpen("", 1)
	defer sessionDb.Close(ctx)
	for _, statement := range []string{
		"ALTER SESSION SET NLS_DATE_LANGUAGE = 'AMERICAN'",
		"ALTER SESSION SET NLS_DATE_FORMAT = 'DD.MM.YYYY'",
		"ALTER SESSION SET NLS_TIME_FORMAT = 'HH24.MI.SS'",
		"ALTER SESSION SET NLS_TIMESTAMP_FORMAT = 'DD.MM.YYYY HH24:MI:SS'",
	} {
		if _, err := sessionDb.Exec(ctx, statement); err != nil {
			gtest.Fatal(err)
		}
	}

	const nlsSql = "SELECT TO_CHAR(DATE '2026-03-04', 'fmDay') W, TO_CHAR(DATE '2026-03-04') D, " +
		"TO_CHAR(TIMESTAMP '2026-03-04 15:26:37') T FROM dual"
	gtest.C(t, func(t *gtest.T) {
		one, err := sessionDb.GetOne(ctx, nlsSql)
		t.AssertNil(err)
		t.Assert(one, g.Map{"W": "WEDNESDAY", "D": "04.03.2026", "T": "04.03.2026 15:26:37"})

		one, err = db.GetOne(ctx, nlsSql)
		t.AssertNil(err)
		t.Assert(one, g.Map{"W": "星期三", "D": "2026-03-04", "T": "2026-03-04 15:26:37.000000"})
	})

	withTime := gtime.NewFromStr("2026-03-04 15:26:37")
	gtest.C(t, func(t *gtest.T) {
		_, err := sessionDb.Model(table).Data(g.Map{"id": 1, "dt": withTime, "amount": 2.5, "ratio": 1.25}).Insert()
		t.AssertNil(err)

		one, err := db.GetOne(ctx, fmt.Sprintf(
			"SELECT TO_CHAR(DT, 'YYYY-MM-DD HH24:MI:SS') DT, AMOUNT, RATIO FROM %s WHERE ID = 1", table,
		))
		t.AssertNil(err)
		t.Assert(one["DT"], "2026-03-04 15:26:37")
		t.Assert(one["AMOUNT"].Float64(), 2.5)
		t.Assert(one["RATIO"].Float64(), 1.25)

		count, err := sessionDb.Model(table).Where("dt", withTime).Where("amount", 2.5).Count()
		t.AssertNil(err)
		t.Assert(count, 1)
	})

	// A time is written to a DATE and a TIME column whatever the NLS_DATE_FORMAT and NLS_TIME_FORMAT
	// of the session are.
	gtest.C(t, func(t *gtest.T) {
		_, err := sessionDb.Model(table).Data(g.Map{"id": 2, "d": withTime, "tm": withTime}).Insert()
		t.AssertNil(err)

		one, err := db.GetOne(ctx, fmt.Sprintf(
			"SELECT TO_CHAR(D, 'YYYY-MM-DD HH24:MI:SS') D, TO_CHAR(TM, 'HH24:MI:SS') TM FROM %s WHERE ID = 2", table,
		))
		t.AssertNil(err)
		t.Assert(one, g.Map{"D": "2026-03-04 00:00:00", "TM": "15:26:37"})

		count, err := sessionDb.Model(table).Where("d", gtime.NewFromStr("2026-03-04 00:00:00")).Count()
		t.AssertNil(err)
		t.Assert(count, 1)
	})
}

// Test_DM_ConnectBy_Hierarchy tests hierarchical queries with CONNECT BY, including paging
// that keeps ORDER SIBLINGS BY.
func Test_DM_ConnectBy_Hierarchy(t *testing.T) {
	table := dmFeatureTableName("t_tree")
	dmFeatureExec(fmt.Sprintf(
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

// Test_DM_Listagg_WithinGroup tests LISTAGG ... WITHIN GROUP through the Model with GROUP BY and HAVING.
func Test_DM_Listagg_WithinGroup(t *testing.T) {
	table := createInitUserTable()
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

// Test_DM_WindowFunction_Page tests window functions combined with paging, where ROW_NUMBER is
// computed over the whole table before the page is taken.
func Test_DM_WindowFunction_Page(t *testing.T) {
	table := createInitUserTable()
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

// Test_DM_PlSql_OutBind tests anonymous PL/SQL blocks through Exec, returning values into
// positional and named out binds.
// Note: the underlying driver fails to convert a NUMBER out bind into *int64 (6007), so the
// destination is a sql.NullInt64.
func Test_DM_PlSql_OutBind(t *testing.T) {
	table := createUserTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		var id sql.NullInt64
		_, err := db.Exec(ctx, fmt.Sprintf(
			"BEGIN INSERT INTO %s (PASSPORT, PASSWORD, NICKNAME) VALUES (?, 'p', 'n') RETURNING ID INTO ?; END;", table,
		), "positional", sql.Out{Dest: &id})
		t.AssertNil(err)
		t.Assert(id, sql.NullInt64{Int64: 1, Valid: true})

		value, err := db.Model(table).Where("id", id.Int64).Value("passport")
		t.AssertNil(err)
		t.Assert(value, "positional")
	})

	gtest.C(t, func(t *gtest.T) {
		var id sql.NullInt64
		_, err := db.Exec(ctx, fmt.Sprintf(
			"BEGIN INSERT INTO %s (PASSPORT, PASSWORD, NICKNAME) VALUES (:passport, 'p', 'n') RETURNING ID INTO :id; END;", table,
		), sql.Named("passport", "named"), sql.Named("id", sql.Out{Dest: &id}))
		t.AssertNil(err)
		t.Assert(id, sql.NullInt64{Int64: 2, Valid: true})

		value, err := db.Model(table).Where("id", id.Int64).Value("passport")
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

// Test_DM_LastInsertId_NonIntegerPrimaryKeyId tests that InsertAndGetId inserts into a table whose
// primary key is not an integer and returns 0 with a CodeNotSupported error, as there is no integer
// id to return.
func Test_DM_LastInsertId_NonIntegerPrimaryKeyId(t *testing.T) {
	table := dmFeatureTableName("t_rawid")
	dmFeatureExec(fmt.Sprintf(`CREATE TABLE %s (ID RAW(16) DEFAULT SYS_GUID() NOT NULL, NAME VARCHAR2(20), PRIMARY KEY (ID))`, table))
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

// Test_DM_TableAlias_AsKeyword tests table expressions giving the alias of a table or a sub-query
// with the AS keyword.
func Test_DM_TableAlias_AsKeyword(t *testing.T) {
	table := createInitUserTable()
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

// Test_Model_Insert_EmptyIdentity tests an insert giving the AUTO_INCREMENT column an empty
// string, which DM takes as NULL, so that it generates the value.
func Test_Model_Insert_EmptyIdentity(t *testing.T) {
	table := dmFeatureTableName("t_identity")
	dmFeatureExec(fmt.Sprintf(`CREATE TABLE %s (ID INT AUTO_INCREMENT NOT NULL, PASSPORT VARCHAR(45), PRIMARY KEY (ID))`, table))
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		for i := 1; i <= 2; i++ {
			id, err := db.Model(table).Data(g.Map{"id": "", "passport": fmt.Sprintf("user_%d", i)}).InsertAndGetId()
			t.AssertNil(err)
			t.Assert(id, i)
		}

		one, err := db.Model(table).Where("passport", "user_2").One()
		t.AssertNil(err)
		t.Assert(one["ID"], 2)
	})
}

// Test_Transaction_Options_EndWithTransaction tests that the isolation level and the read-only
// mode of a transaction end with it on its pooled connection, also when the transaction executes
// no statement, returns an error before executing one, or only queries.
func Test_Transaction_Options_EndWithTransaction(t *testing.T) {
	table := createInitUserTable()
	defer dropTable(table)

	// Note: the pool has one connection, so that the statements following a transaction run on
	// the connection of the transaction.
	pool, err := gdb.New(gdb.ConfigNode{
		Host:             TestDBHost,
		Port:             TestDBPort,
		User:             TestDBUser,
		Pass:             TestDBPass,
		Name:             TestDBName,
		Type:             TestDBType,
		Charset:          TestCharset,
		MaxOpenConnCount: 1,
	})
	gtest.AssertNil(err)
	defer pool.Close(ctx)

	gtest.C(t, func(t *gtest.T) {
		const transactionSql = "SELECT ISOLATION, READ_ONLY FROM V$TRX WHERE SESS_ID = SESSID()"
		var (
			errRollback = errors.New("rollback")
			// Note: V$TRX reports READ COMMITTED, the default isolation level of DM, as 1 and
			// SERIALIZABLE as 3.
			defaultTransaction = g.Map{"ISOLATION": 1, "READ_ONLY": "N"}
			// The statements following a transaction, in a transaction with the default options
			// and in autocommit mode.
			followers = []func() error{
				func() error {
					return pool.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
						one, err := tx.GetOne(transactionSql)
						t.AssertNil(err)
						t.Assert(one, defaultTransaction)
						_, err = tx.Update(table, g.Map{"nickname": "transaction"}, "id=1")
						return err
					})
				},
				func() error {
					one, err := pool.GetOne(ctx, transactionSql)
					t.AssertNil(err)
					t.Assert(one, defaultTransaction)
					_, err = pool.Update(ctx, table, g.Map{"nickname": "autocommit"}, "id=2")
					return err
				},
			}
		)
		for _, readOnly := range []bool{false, true} {
			var (
				options = gdb.TxOptions{
					Propagation: gdb.PropagationRequiresNew,
					Isolation:   sql.LevelSerializable,
					ReadOnly:    readOnly,
				}
				transaction = g.Map{"ISOLATION": 3, "READ_ONLY": "N"}
			)
			if readOnly {
				transaction["READ_ONLY"] = "Y"
			}
			for _, item := range []struct {
				body func(ctx context.Context, tx gdb.TX) error
				err  error
			}{
				// A transaction executing no statement.
				{func(ctx context.Context, tx gdb.TX) error { return nil }, nil},
				// A transaction returning an error before executing a statement.
				{func(ctx context.Context, tx gdb.TX) error { return errRollback }, errRollback},
				// A transaction only querying, which runs with the options.
				{func(ctx context.Context, tx gdb.TX) error {
					one, err := tx.GetOne(transactionSql)
					t.AssertNil(err)
					t.Assert(one, transaction)
					return err
				}, nil},
			} {
				for _, follow := range followers {
					t.Assert(pool.TransactionWithOptions(ctx, options, item.body), item.err)
					t.AssertNil(follow())
				}
			}
		}
	})
}

// Test_Transaction_Isolation_RepeatableRead_WriteConflict tests that a REPEATABLE READ transaction,
// which DM runs as SERIALIZABLE, fails to update a record that another transaction has updated
// since it began (-7066).
func Test_Transaction_Isolation_RepeatableRead_WriteConflict(t *testing.T) {
	table := createInitUserTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		err := db.TransactionWithOptions(ctx, gdb.TxOptions{
			Propagation: gdb.PropagationRequiresNew,
			Isolation:   sql.LevelRepeatableRead,
		}, func(ctx context.Context, tx gdb.TX) error {
			isolation, err := tx.GetValue("SELECT ISOLATION FROM V$TRX WHERE SESS_ID = SESSID()")
			t.AssertNil(err)
			t.Assert(isolation, 3)

			_, err = db.Model(table).Data("passport", "other").Where("id", 1).Update()
			t.AssertNil(err)
			_, err = tx.Model(table).Data("passport", "mine").Where("id", 1).Update()
			return err
		})
		t.Assert(isDmError(err, -7066), true)

		value, err := db.Model(table).Where("id", 1).Value("passport")
		t.AssertNil(err)
		t.Assert(value, "other")
	})
}

// Test_DataType_VarBinary_Where tests a condition on a VARBINARY column.
func Test_DataType_VarBinary_Where(t *testing.T) {
	table := "test_varbinary_" + gtime.TimestampMicroStr()
	_, err := db.Exec(ctx, "CREATE TABLE "+table+" (id INT PRIMARY KEY, data VARBINARY(32))")
	if err != nil {
		t.Fatal(err)
	}
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(table).Data(g.Map{"id": 1, "data": "abc"}).Insert()
		t.AssertNil(err)

		// Note: a string condition is bound as text, which DM converts to binary data as
		// hexadecimal digits, so the value is matched by its bytes.
		one, err := db.Model(table).Where("data", []byte("abc")).One()
		t.AssertNil(err)
		t.Assert(one["ID"], 1)
		t.Assert(one["DATA"], "abc")
	})
}

// Test_DM_SoftTime_NumberColumns tests soft time fields holding a unix timestamp in NUMBER(10)
// columns, the integer columns of tables written for Oracle.
func Test_DM_SoftTime_NumberColumns(t *testing.T) {
	table := dmFeatureTableName("t_soft_number")
	dmFeatureExec(fmt.Sprintf(`CREATE TABLE %s (
    ID NUMBER(10) NOT NULL,
    NAME VARCHAR(45),
    CREATE_AT NUMBER(10),
    UPDATE_AT NUMBER(10),
    DELETE_AT NUMBER(10),
    PRIMARY KEY (ID)
)`, table))
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(table).Data(g.Map{"id": 1, "name": "name_1"}).Insert()
		t.AssertNil(err)

		one, err := db.Model(table).WherePri(1).One()
		t.AssertNil(err)
		t.AssertGT(one["CREATE_AT"].Int64(), 0)
		t.AssertGT(one["UPDATE_AT"].Int64(), 0)
		t.Assert(one["DELETE_AT"].Int64(), 0)

		_, err = db.Model(table).WherePri(1).Delete()
		t.AssertNil(err)

		count, err := db.Model(table).Count()
		t.AssertNil(err)
		t.Assert(count, 0)

		one, err = db.Model(table).Unscoped().WherePri(1).One()
		t.AssertNil(err)
		t.Assert(one["NAME"], "name_1")
		t.AssertGT(one["DELETE_AT"].Int64(), 0)
	})
}

// Test_DM_Schema_LetterCase tests the schema that a database object opened with a configured name
// works in: the name in upper case, as DM folds an unquoted name, or the name inside double quotes
// as it is, so that an unquoted name does not reach a schema created in lower case.
func Test_DM_Schema_LetterCase(t *testing.T) {
	var (
		upper = strings.ToUpper(dmFeatureTableName("gf_upper"))
		lower = dmFeatureTableName("gf_lower")
	)
	dmFeatureExec(
		fmt.Sprintf(`CREATE SCHEMA %s`, upper),
		fmt.Sprintf(`CREATE SCHEMA "%s"`, lower),
	)
	defer func() {
		_, _ = db.Exec(ctx, fmt.Sprintf(`DROP SCHEMA %s CASCADE`, upper))
		_, _ = db.Exec(ctx, fmt.Sprintf(`DROP SCHEMA "%s" CASCADE`, lower))
	}()

	currentSchema := func(name string) (string, error) {
		schemaDb, err := gdb.New(gdb.ConfigNode{
			Host:    TestDBHost,
			Port:    TestDBPort,
			User:    TestDBUser,
			Pass:    TestDBPass,
			Name:    name,
			Type:    TestDBType,
			Charset: TestCharset,
		})
		if err != nil {
			return "", err
		}
		defer schemaDb.Close(ctx)
		value, err := schemaDb.GetValue(ctx, "SELECT SF_GET_SCHEMA_NAME_BY_ID(CURRENT_SCHID()) FROM DUAL")
		return value.String(), err
	}
	gtest.C(t, func(t *gtest.T) {
		for name, schema := range map[string]string{
			strings.ToLower(upper): upper,
			upper:                  upper,
			`"` + lower + `"`:      lower,
		} {
			current, err := currentSchema(name)
			t.AssertNil(err)
			t.Assert(current, schema)
		}

		_, err := currentSchema(lower)
		assertDmError(t, err, dmErrInvalidSchema)
	})
	gtest.C(t, func(t *gtest.T) {
		dmFeatureExec(fmt.Sprintf(`CREATE TABLE "%s".T_ITEM (ID INT IDENTITY(1, 1) PRIMARY KEY, NAME VARCHAR(20))`, lower))
		schemaDb, err := gdb.New(gdb.ConfigNode{
			Host:    TestDBHost,
			Port:    TestDBPort,
			User:    TestDBUser,
			Pass:    TestDBPass,
			Name:    `"` + lower + `"`,
			Type:    TestDBType,
			Charset: TestCharset,
		})
		t.AssertNil(err)
		defer schemaDb.Close(ctx)

		fields, err := schemaDb.TableFields(ctx, "T_ITEM")
		t.AssertNil(err)
		t.Assert(len(fields), 2)
		t.Assert(fields["ID"].Extra, "auto_increment")

		id, err := schemaDb.Model("T_ITEM").Data(g.Map{"NAME": "a"}).InsertAndGetId()
		t.AssertNil(err)
		t.Assert(id, 1)
	})
}
