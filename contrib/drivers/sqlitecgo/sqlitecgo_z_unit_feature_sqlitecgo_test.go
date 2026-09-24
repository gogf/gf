// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package sqlitecgo_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/text/gstr"
)

func featureTable(name string) string {
	return fmt.Sprintf(`feature_%s_%d`, name, gtime.TimestampNano())
}

func featureFile(name string) string {
	return gfile.Join(dbDir, fmt.Sprintf(`feature_%s_%d.db`, name, gtime.TimestampNano()))
}

// =============================================================================
// Link and connection options
// =============================================================================

// Test_SQLite_Link_MemoryDatabase tests the `sqlite::@file(:memory:)` link form.
func Test_SQLiteCgo_Link_MemoryDatabase(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		memDB, err := gdb.New(gdb.ConfigNode{
			Type:             "sqlite",
			Link:             `sqlite::@file(:memory:)`,
			MaxOpenConnCount: 1,
		})
		t.AssertNil(err)

		table := featureTable("mem")
		_, err = memDB.Exec(ctx, fmt.Sprintf(
			`CREATE TABLE %s (id INTEGER PRIMARY KEY, name VARCHAR(45))`, table,
		))
		t.AssertNil(err)

		_, err = memDB.Model(table).Data(g.Map{"id": 1, "name": "in_memory"}).Insert()
		t.AssertNil(err)

		one, err := memDB.Model(table).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["name"].String(), "in_memory")

		all, err := memDB.GetAll(ctx, `PRAGMA database_list`)
		t.AssertNil(err)
		t.Assert(len(all), 1)
		t.Assert(all[0]["name"].String(), "main")
		t.Assert(all[0]["file"].String(), "")

		t.Assert(gfile.Exists(":memory:"), false)
	})
}

// Test_SQLite_Extra_Pragmas tests that ConfigNode.Extra entries reach the database as PRAGMAs.
func Test_SQLiteCgo_Extra_Pragmas(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		tunedDB, err := gdb.New(gdb.ConfigNode{
			Type:  "sqlite",
			Link:  fmt.Sprintf(`sqlite::@file(%s)`, featureFile("pragma_wal")),
			Extra: "journal_mode=WAL&busy_timeout=7531",
		})
		t.AssertNil(err)

		mode, err := tunedDB.GetValue(ctx, `PRAGMA journal_mode`)
		t.AssertNil(err)
		t.Assert(mode.String(), "wal")

		timeout, err := tunedDB.GetValue(ctx, `PRAGMA busy_timeout`)
		t.AssertNil(err)
		t.Assert(timeout.Int(), 7531)

		plainDB, err := gdb.New(gdb.ConfigNode{
			Type: "sqlite",
			Link: fmt.Sprintf(`sqlite::@file(%s)`, featureFile("pragma_plain")),
		})
		t.AssertNil(err)

		mode, err = plainDB.GetValue(ctx, `PRAGMA journal_mode`)
		t.AssertNil(err)
		t.Assert(mode.String(), "delete")
	})
}

// Test_SQLite_New_TwoHandlesOnSameFileWAL tests two gdb handles reading and writing one WAL file.
func Test_SQLiteCgo_New_TwoHandlesOnSameFileWAL(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		node := gdb.ConfigNode{
			Type:  "sqlite",
			Link:  fmt.Sprintf(`sqlite::@file(%s)`, featureFile("two_handles")),
			Extra: "journal_mode=WAL&busy_timeout=5000",
		}
		first, err := gdb.New(node)
		t.AssertNil(err)
		second, err := gdb.New(node)
		t.AssertNil(err)

		table := featureTable("shared")
		_, err = first.Exec(ctx, fmt.Sprintf(
			`CREATE TABLE %s (id INTEGER PRIMARY KEY, writer VARCHAR(45))`, table,
		))
		t.AssertNil(err)
		defer dropTableWithDb(first, table)

		_, err = first.Model(table).Data(g.Map{"id": 1, "writer": "first"}).Insert()
		t.AssertNil(err)

		one, err := second.Model(table).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["writer"].String(), "first")

		_, err = second.Model(table).Data(g.Map{"id": 2, "writer": "second"}).Insert()
		t.AssertNil(err)

		one, err = first.Model(table).Where("id", 2).One()
		t.AssertNil(err)
		t.Assert(one["writer"].String(), "second")

		for _, handle := range []gdb.DB{first, second} {
			mode, err := handle.GetValue(ctx, `PRAGMA journal_mode`)
			t.AssertNil(err)
			t.Assert(mode.String(), "wal")
		}
	})
}

// =============================================================================
// Conflict resolution: INSERT OR IGNORE / INSERT OR REPLACE / ON CONFLICT
// =============================================================================

func createConflictTable(t *gtest.T, name string) string {
	table := featureTable(name)
	_, err := db.Exec(ctx, fmt.Sprintf(`
		CREATE TABLE %s (
			id    INTEGER PRIMARY KEY AUTOINCREMENT,
			email VARCHAR(45) NOT NULL UNIQUE,
			name  VARCHAR(45),
			score INTEGER
		)`, table,
	))
	t.AssertNil(err)
	return table
}

// Test_SQLite_InsertIgnore_OrIgnore tests InsertIgnore emitting INSERT OR IGNORE.
func Test_SQLiteCgo_InsertIgnore_OrIgnore(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		table := createConflictTable(t, "ins_ignore")
		defer dropTable(table)

		_, err := db.Model(table).Data(g.Map{
			"email": "a@example.com", "name": "alice", "score": 10,
		}).Insert()
		t.AssertNil(err)

		sqlArray, err := gdb.CatchSQL(ctx, func(ctx context.Context) error {
			_, e := db.Ctx(ctx).Model(table).Data(g.Map{
				"email": "a@example.com", "name": "duplicate", "score": 99,
			}).InsertIgnore()
			return e
		})
		t.AssertNil(err)
		t.AssertGT(len(sqlArray), 0)
		t.Assert(gstr.HasPrefix(sqlArray[len(sqlArray)-1], "INSERT OR IGNORE INTO"), true)

		count, err := db.Model(table).Count()
		t.AssertNil(err)
		t.Assert(count, 1)

		one, err := db.Model(table).Where("email", "a@example.com").One()
		t.AssertNil(err)
		t.Assert(one["name"].String(), "alice")
		t.Assert(one["score"].Int(), 10)
	})
}

// Test_SQLite_Replace_OrReplace tests Replace emitting INSERT OR REPLACE.
func Test_SQLiteCgo_Replace_OrReplace(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		table := createConflictTable(t, "replace")
		defer dropTable(table)

		id, err := db.Model(table).Data(g.Map{
			"email": "a@example.com", "name": "alice", "score": 10,
		}).InsertAndGetId()
		t.AssertNil(err)

		sqlArray, err := gdb.CatchSQL(ctx, func(ctx context.Context) error {
			_, e := db.Ctx(ctx).Model(table).Data(g.Map{
				"email": "a@example.com", "name": "replaced", "score": 20,
			}).Replace()
			return e
		})
		t.AssertNil(err)
		t.AssertGT(len(sqlArray), 0)
		t.Assert(gstr.HasPrefix(sqlArray[len(sqlArray)-1], "INSERT OR REPLACE INTO"), true)

		count, err := db.Model(table).Count()
		t.AssertNil(err)
		t.Assert(count, 1)

		one, err := db.Model(table).Where("email", "a@example.com").One()
		t.AssertNil(err)
		t.Assert(one["name"].String(), "replaced")
		t.Assert(one["score"].Int(), 20)
		t.AssertNE(one["id"].Int64(), id)
	})
}

// Test_SQLite_Save_OnConflictDoUpdate tests Save emitting ON CONFLICT (...) DO UPDATE SET.
func Test_SQLiteCgo_Save_OnConflictDoUpdate(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		table := createConflictTable(t, "on_conflict")
		defer dropTable(table)

		id, err := db.Model(table).Data(g.Map{
			"email": "a@example.com", "name": "alice", "score": 10,
		}).InsertAndGetId()
		t.AssertNil(err)

		sqlArray, err := gdb.CatchSQL(ctx, func(ctx context.Context) error {
			_, e := db.Ctx(ctx).Model(table).OnConflict("email").Data(g.Map{
				"email": "a@example.com", "name": "upserted", "score": 42,
			}).Save()
			return e
		})
		t.AssertNil(err)
		t.AssertGT(len(sqlArray), 0)
		lastSql := sqlArray[len(sqlArray)-1]
		t.Assert(gstr.Contains(lastSql, "ON CONFLICT (email) DO UPDATE SET"), true)
		t.Assert(gstr.Contains(lastSql, "EXCLUDED."), true)

		count, err := db.Model(table).Count()
		t.AssertNil(err)
		t.Assert(count, 1)

		one, err := db.Model(table).Where("email", "a@example.com").One()
		t.AssertNil(err)
		t.Assert(one["id"].Int64(), id)
		t.Assert(one["name"].String(), "upserted")
		t.Assert(one["score"].Int(), 42)
	})
}

// Test_SQLite_OnConflict_MultiColumnTarget tests a conflict target made of two columns.
func Test_SQLiteCgo_OnConflict_MultiColumnTarget(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		table := featureTable("multi_unique")
		_, err := db.Exec(ctx, fmt.Sprintf(`
			CREATE TABLE %s (
				tenant VARCHAR(45) NOT NULL,
				email  VARCHAR(45) NOT NULL,
				score  INTEGER,
				UNIQUE(tenant, email)
			)`, table,
		))
		t.AssertNil(err)
		defer dropTable(table)

		_, err = db.Model(table).Data(g.List{
			{"tenant": "t1", "email": "a@example.com", "score": 1},
			{"tenant": "t2", "email": "a@example.com", "score": 2},
		}).Insert()
		t.AssertNil(err)

		sqlArray, err := gdb.CatchSQL(ctx, func(ctx context.Context) error {
			_, e := db.Ctx(ctx).Model(table).OnConflict("tenant", "email").Data(g.Map{
				"tenant": "t1", "email": "a@example.com", "score": 100,
			}).Save()
			return e
		})
		t.AssertNil(err)
		t.AssertGT(len(sqlArray), 0)
		t.Assert(gstr.Contains(sqlArray[len(sqlArray)-1], "ON CONFLICT (tenant,email) DO UPDATE SET"), true)

		count, err := db.Model(table).Count()
		t.AssertNil(err)
		t.Assert(count, 2)

		score, err := db.Model(table).Where("tenant", "t1").Where("email", "a@example.com").Value("score")
		t.AssertNil(err)
		t.Assert(score.Int(), 100)

		score, err = db.Model(table).Where("tenant", "t2").Where("email", "a@example.com").Value("score")
		t.AssertNil(err)
		t.Assert(score.Int(), 2)
	})
}

// =============================================================================
// RETURNING clause
// =============================================================================

// Test_SQLite_Returning_Insert tests INSERT ... RETURNING through db.GetAll.
func Test_SQLiteCgo_Returning_Insert(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		table := createInitTable()
		defer dropTable(table)

		all, err := db.GetAll(ctx, fmt.Sprintf(
			`INSERT INTO %s (passport, password, nickname) VALUES ('ret_1','pass','nick_1'),('ret_2','pass','nick_2')
			 RETURNING id, passport`, table,
		))
		t.AssertNil(err)
		t.Assert(len(all), 2)
		t.Assert(all[0]["id"].Int(), TableSize+1)
		t.Assert(all[0]["passport"].String(), "ret_1")
		t.Assert(all[1]["id"].Int(), TableSize+2)
		t.Assert(all[1]["passport"].String(), "ret_2")

		count, err := db.Model(table).Count()
		t.AssertNil(err)
		t.Assert(count, TableSize+2)
	})
}

// Test_SQLite_Returning_UpdateAndDelete tests UPDATE ... RETURNING and DELETE ... RETURNING.
func Test_SQLiteCgo_Returning_UpdateAndDelete(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		table := createInitTable()
		defer dropTable(table)

		all, err := db.GetAll(ctx, fmt.Sprintf(
			`UPDATE %s SET nickname='renamed' WHERE id<=3 RETURNING id, nickname`, table,
		))
		t.AssertNil(err)
		t.Assert(len(all), 3)
		t.Assert(all[0]["id"].Int(), 1)
		t.Assert(all[2]["id"].Int(), 3)
		t.Assert(all[0]["nickname"].String(), "renamed")

		count, err := db.Model(table).Where("nickname", "renamed").Count()
		t.AssertNil(err)
		t.Assert(count, 3)

		all, err = db.GetAll(ctx, fmt.Sprintf(
			`DELETE FROM %s WHERE id>8 RETURNING id, passport`, table,
		))
		t.AssertNil(err)
		t.Assert(len(all), 2)
		t.Assert(all[0]["id"].Int(), 9)
		t.Assert(all[0]["passport"].String(), "user_9")
		t.Assert(all[1]["id"].Int(), 10)

		count, err = db.Model(table).Count()
		t.AssertNil(err)
		t.Assert(count, TableSize-2)
	})
}

// =============================================================================
// CTE and compound queries
// =============================================================================

// Test_SQLite_CTE_Basic tests a WITH ... AS common table expression.
func Test_SQLiteCgo_CTE_Basic(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		table := createInitTable()
		defer dropTable(table)

		all, err := db.GetAll(ctx, fmt.Sprintf(`
			WITH low AS (
				SELECT id, nickname FROM %s WHERE id<=5
			)
			SELECT * FROM low WHERE id>2 ORDER BY id`, table,
		))
		t.AssertNil(err)
		t.Assert(len(all), 3)
		t.Assert(all[0]["id"].Int(), 3)
		t.Assert(all[0]["nickname"].String(), "name_3")
		t.Assert(all[2]["id"].Int(), 5)
	})
}

// Test_SQLite_CTE_Recursive tests a WITH RECURSIVE common table expression.
func Test_SQLiteCgo_CTE_Recursive(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		table := featureTable("cte_tree")
		_, err := db.Exec(ctx, fmt.Sprintf(`
			CREATE TABLE %s (
				id        INTEGER PRIMARY KEY,
				parent_id INTEGER,
				name      VARCHAR(45)
			)`, table,
		))
		t.AssertNil(err)
		defer dropTable(table)

		_, err = db.Model(table).Data(g.List{
			{"id": 1, "parent_id": nil, "name": "root"},
			{"id": 2, "parent_id": 1, "name": "child_a"},
			{"id": 3, "parent_id": 1, "name": "child_b"},
			{"id": 4, "parent_id": 2, "name": "grandchild"},
		}).Insert()
		t.AssertNil(err)

		all, err := db.GetAll(ctx, fmt.Sprintf(`
			WITH RECURSIVE tree(id, name, depth) AS (
				SELECT id, name, 0 FROM %s WHERE id=1
				UNION ALL
				SELECT t.id, t.name, tree.depth+1 FROM %s t JOIN tree ON t.parent_id=tree.id
			)
			SELECT * FROM tree ORDER BY depth, id`, table, table,
		))
		t.AssertNil(err)
		t.Assert(len(all), 4)
		t.Assert(all[0]["name"].String(), "root")
		t.Assert(all[0]["depth"].Int(), 0)
		t.Assert(all[3]["name"].String(), "grandchild")
		t.Assert(all[3]["depth"].Int(), 2)

		total, err := db.Raw(`
			WITH RECURSIVE counter(n) AS (
				SELECT 1 UNION ALL SELECT n+1 FROM counter WHERE n<100
			)
			SELECT SUM(n) AS total FROM counter`,
		).Value()
		t.AssertNil(err)
		t.Assert(total.Int(), 5050)
	})
}

// Test_SQLite_Compound_ExceptIntersect tests the EXCEPT and INTERSECT compound operators.
func Test_SQLiteCgo_Compound_ExceptIntersect(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		table := createInitTable()
		defer dropTable(table)

		all, err := db.GetAll(ctx, fmt.Sprintf(
			`SELECT id FROM %s WHERE id<=5 EXCEPT SELECT id FROM %s WHERE id<=2 ORDER BY id`,
			table, table,
		))
		t.AssertNil(err)
		t.Assert(len(all), 3)
		t.Assert(all[0]["id"].Int(), 3)
		t.Assert(all[2]["id"].Int(), 5)

		all, err = db.GetAll(ctx, fmt.Sprintf(
			`SELECT id FROM %s WHERE id<=5 INTERSECT SELECT id FROM %s WHERE id>=4 ORDER BY id`,
			table, table,
		))
		t.AssertNil(err)
		t.Assert(len(all), 2)
		t.Assert(all[0]["id"].Int(), 4)
		t.Assert(all[1]["id"].Int(), 5)
	})
}

// Test_SQLite_Union_OperandOrderByLimit tests Union with ORDER BY and LIMIT on each operand.
func Test_SQLiteCgo_Union_OperandOrderByLimit(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		table := createInitTable()
		defer dropTable(table)

		var all gdb.Result
		sqlArray, err := gdb.CatchSQL(ctx, func(ctx context.Context) error {
			var e error
			all, e = db.Ctx(ctx).Union(
				db.Model(table).Fields("id").Where("id<=3").OrderAsc("id").Limit(2),
				db.Model(table).Fields("id").Where("id>8").OrderDesc("id").Limit(1),
			).All()
			return e
		})
		t.AssertNil(err)
		t.AssertGT(len(sqlArray), 0)

		lastSql := sqlArray[len(sqlArray)-1]
		t.Assert(gstr.HasPrefix(lastSql, "SELECT * FROM ("), true)
		t.Assert(gstr.Contains(lastSql, ") UNION SELECT * FROM ("), true)

		t.Assert(len(all), 3)
		t.Assert(all[0]["id"].Int(), 1)
		t.Assert(all[1]["id"].Int(), 2)
		t.Assert(all[2]["id"].Int(), 10)
	})
}

// Test_SQLite_UnionAll_Model tests UnionAll keeping the duplicates that Union removes.
func Test_SQLiteCgo_UnionAll_Model(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		table := createInitTable()
		defer dropTable(table)

		all, err := db.Union(
			db.Model(table).Fields("id").Where("id<=2"),
			db.Model(table).Fields("id").Where("id<=2"),
		).All()
		t.AssertNil(err)
		t.Assert(len(all), 2)

		all, err = db.UnionAll(
			db.Model(table).Fields("id").Where("id<=2"),
			db.Model(table).Fields("id").Where("id<=2"),
		).All()
		t.AssertNil(err)
		t.Assert(len(all), 4)

		var ids []int
		for _, record := range all {
			ids = append(ids, record["id"].Int())
		}
		t.Assert(ids, []int{1, 2, 1, 2})
	})
}

// =============================================================================
// Window functions
// =============================================================================

// Test_SQLite_Window_RowNumber tests ROW_NUMBER() OVER.
func Test_SQLiteCgo_Window_RowNumber(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		table := createInitTable()
		defer dropTable(table)

		all, err := db.GetAll(ctx, fmt.Sprintf(
			`SELECT id, ROW_NUMBER() OVER (ORDER BY id DESC) AS rn FROM %s ORDER BY id`, table,
		))
		t.AssertNil(err)
		t.Assert(len(all), TableSize)
		t.Assert(all[0]["id"].Int(), 1)
		t.Assert(all[0]["rn"].Int(), TableSize)
		t.Assert(all[TableSize-1]["id"].Int(), TableSize)
		t.Assert(all[TableSize-1]["rn"].Int(), 1)
	})
}

// Test_SQLite_Window_RankDenseRank tests RANK() and DENSE_RANK() over a partition.
func Test_SQLiteCgo_Window_RankDenseRank(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		table := featureTable("window_rank")
		_, err := db.Exec(ctx, fmt.Sprintf(`
			CREATE TABLE %s (
				id     INTEGER PRIMARY KEY AUTOINCREMENT,
				dept   VARCHAR(20),
				salary INTEGER
			)`, table,
		))
		t.AssertNil(err)
		defer dropTable(table)

		_, err = db.Model(table).Data(g.List{
			{"dept": "eng", "salary": 200},
			{"dept": "eng", "salary": 100},
			{"dept": "eng", "salary": 100},
			{"dept": "sales", "salary": 300},
			{"dept": "sales", "salary": 150},
		}).Insert()
		t.AssertNil(err)

		all, err := db.GetAll(ctx, fmt.Sprintf(`
			SELECT dept, salary,
				RANK() OVER (PARTITION BY dept ORDER BY salary DESC) AS rnk,
				DENSE_RANK() OVER (PARTITION BY dept ORDER BY salary DESC) AS dense_rnk
			FROM %s ORDER BY dept, salary DESC, id`, table,
		))
		t.AssertNil(err)
		t.Assert(len(all), 5)

		t.Assert(all[0]["dept"].String(), "eng")
		t.Assert(all[0]["salary"].Int(), 200)
		t.Assert(all[0]["rnk"].Int(), 1)
		t.Assert(all[1]["rnk"].Int(), 2)
		t.Assert(all[2]["rnk"].Int(), 2)
		t.Assert(all[1]["dense_rnk"].Int(), 2)
		t.Assert(all[2]["dense_rnk"].Int(), 2)

		t.Assert(all[3]["dept"].String(), "sales")
		t.Assert(all[3]["rnk"].Int(), 1)
		t.Assert(all[4]["rnk"].Int(), 2)
	})
}

// Test_SQLite_Window_LagLead tests LAG() and LEAD().
func Test_SQLiteCgo_Window_LagLead(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		table := createInitTable()
		defer dropTable(table)

		all, err := db.GetAll(ctx, fmt.Sprintf(`
			SELECT id,
				LAG(id, 1) OVER (ORDER BY id)  AS prev_id,
				LEAD(id, 1) OVER (ORDER BY id) AS next_id,
				LAG(id, 1, -1) OVER (ORDER BY id) AS prev_or_default
			FROM %s ORDER BY id`, table,
		))
		t.AssertNil(err)
		t.Assert(len(all), TableSize)

		t.Assert(all[0]["prev_id"], nil)
		t.Assert(all[0]["next_id"].Int(), 2)
		t.Assert(all[0]["prev_or_default"].Int(), -1)

		t.Assert(all[4]["prev_id"].Int(), 4)
		t.Assert(all[4]["next_id"].Int(), 6)

		t.Assert(all[TableSize-1]["prev_id"].Int(), TableSize-1)
		t.Assert(all[TableSize-1]["next_id"], nil)
	})
}

// Test_SQLite_Window_SumPartitionBy tests SUM() OVER with and without PARTITION BY.
func Test_SQLiteCgo_Window_SumPartitionBy(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		table := featureTable("window_sum")
		_, err := db.Exec(ctx, fmt.Sprintf(`
			CREATE TABLE %s (
				id     INTEGER PRIMARY KEY AUTOINCREMENT,
				dept   VARCHAR(20),
				amount INTEGER
			)`, table,
		))
		t.AssertNil(err)
		defer dropTable(table)

		_, err = db.Model(table).Data(g.List{
			{"dept": "eng", "amount": 10},
			{"dept": "eng", "amount": 20},
			{"dept": "sales", "amount": 100},
			{"dept": "sales", "amount": 200},
			{"dept": "sales", "amount": 300},
		}).Insert()
		t.AssertNil(err)

		all, err := db.GetAll(ctx, fmt.Sprintf(`
			SELECT id, dept, amount,
				SUM(amount) OVER (PARTITION BY dept)               AS dept_total,
				SUM(amount) OVER (PARTITION BY dept ORDER BY id)   AS dept_running
			FROM %s ORDER BY id`, table,
		))
		t.AssertNil(err)
		t.Assert(len(all), 5)

		t.Assert(all[0]["dept_total"].Int(), 30)
		t.Assert(all[1]["dept_total"].Int(), 30)
		t.Assert(all[2]["dept_total"].Int(), 600)
		t.Assert(all[4]["dept_total"].Int(), 600)

		t.Assert(all[0]["dept_running"].Int(), 10)
		t.Assert(all[1]["dept_running"].Int(), 30)
		t.Assert(all[2]["dept_running"].Int(), 100)
		t.Assert(all[3]["dept_running"].Int(), 300)
		t.Assert(all[4]["dept_running"].Int(), 600)

		all, err = db.Model(table).Fields(
			"id", "SUM(amount) OVER (ORDER BY id) AS running",
		).OrderAsc("id").All()
		t.AssertNil(err)
		t.Assert(all[4]["running"].Int(), 630)
	})
}

// =============================================================================
// JSON1 on a TEXT column
// =============================================================================

func createJsonTable(t *gtest.T, name string) string {
	table := featureTable(name)
	_, err := db.Exec(ctx, fmt.Sprintf(`
		CREATE TABLE %s (
			id  INTEGER PRIMARY KEY AUTOINCREMENT,
			doc TEXT
		)`, table,
	))
	t.AssertNil(err)
	return table
}

// Test_SQLite_JSON1_Extract tests json_extract on a TEXT column.
func Test_SQLiteCgo_JSON1_Extract(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		table := createJsonTable(t, "json_extract")
		defer dropTable(table)

		_, err := db.Model(table).Data(g.List{
			{"doc": `{"name":"alice","age":30,"addr":{"city":"beijing"},"tags":["go","db"]}`},
			{"doc": `{"name":"bob","age":25,"addr":{"city":"shanghai"},"tags":["sql"]}`},
		}).Insert()
		t.AssertNil(err)

		one, err := db.Model(table).Fields(
			`json_extract(doc,'$.name') AS name`,
			`json_extract(doc,'$.age') AS age`,
			`json_extract(doc,'$.addr.city') AS city`,
			`json_extract(doc,'$.tags[0]') AS first_tag`,
		).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["name"].String(), "alice")
		t.Assert(one["age"].Int(), 30)
		t.Assert(one["city"].String(), "beijing")
		t.Assert(one["first_tag"].String(), "go")

		all, err := db.Model(table).Fields("id").Where(
			`json_extract(doc,'$.addr.city') = ?`, "shanghai",
		).All()
		t.AssertNil(err)
		t.Assert(len(all), 1)
		t.Assert(all[0]["id"].Int(), 2)

		count, err := db.Model(table).Where(`json_extract(doc,'$.age') > ?`, 26).Count()
		t.AssertNil(err)
		t.Assert(count, 1)
	})
}

// Test_SQLite_JSON1_ArrowOperators tests the -> and ->> operators added in SQLite 3.38.
func Test_SQLiteCgo_JSON1_ArrowOperators(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		table := createJsonTable(t, "json_arrow")
		defer dropTable(table)

		_, err := db.Model(table).Data(g.Map{
			"doc": `{"name":"alice","age":30,"tags":["go","db"]}`,
		}).Insert()
		t.AssertNil(err)

		one, err := db.Model(table).Fields(
			`doc->'$.name'  AS name_json`,
			`doc->>'$.name' AS name_text`,
			`doc->>'$.age'  AS age_text`,
			`doc->'$.tags'->>0 AS first_tag`,
		).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["name_json"].String(), `"alice"`)
		t.Assert(one["name_text"].String(), "alice")
		t.Assert(one["age_text"].Int(), 30)
		t.Assert(one["first_tag"].String(), "go")

		count, err := db.Model(table).Where(`doc->>'$.name' = ?`, "alice").Count()
		t.AssertNil(err)
		t.Assert(count, 1)
	})
}

// Test_SQLite_JSON1_ArrayLengthAndValid tests json_array_length and json_valid.
func Test_SQLiteCgo_JSON1_ArrayLengthAndValid(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		table := createJsonTable(t, "json_valid")
		defer dropTable(table)

		_, err := db.Model(table).Data(g.List{
			{"doc": `{"tags":["a","b","c"]}`},
			{"doc": `{"tags":[]}`},
			{"doc": `not json at all`},
		}).Insert()
		t.AssertNil(err)

		all, err := db.Model(table).Fields(
			"id",
			"json_valid(doc) AS ok",
			`CASE WHEN json_valid(doc) THEN json_array_length(doc,'$.tags') ELSE -1 END AS tag_count`,
		).OrderAsc("id").All()
		t.AssertNil(err)
		t.Assert(len(all), 3)
		t.Assert(all[0]["ok"].Int(), 1)
		t.Assert(all[0]["tag_count"].Int(), 3)
		t.Assert(all[1]["ok"].Int(), 1)
		t.Assert(all[1]["tag_count"].Int(), 0)
		t.Assert(all[2]["ok"].Int(), 0)
		t.Assert(all[2]["tag_count"].Int(), -1)

		count, err := db.Model(table).Where("json_valid(doc)").Count()
		t.AssertNil(err)
		t.Assert(count, 2)
	})
}

// Test_SQLite_JSON1_EachTableValued tests json_each used as a table-valued function in FROM.
func Test_SQLiteCgo_JSON1_EachTableValued(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		table := createJsonTable(t, "json_each")
		defer dropTable(table)

		_, err := db.Model(table).Data(g.List{
			{"doc": `{"name":"alice","tags":["go","db","orm"]}`},
			{"doc": `{"name":"bob","tags":["sql"]}`},
		}).Insert()
		t.AssertNil(err)

		all, err := db.GetAll(ctx, fmt.Sprintf(`
			SELECT t.id, je.key AS idx, je.value AS tag
			FROM %s t, json_each(t.doc,'$.tags') je
			ORDER BY t.id, je.key`, table,
		))
		t.AssertNil(err)
		t.Assert(len(all), 4)
		t.Assert(all[0]["id"].Int(), 1)
		t.Assert(all[0]["idx"].Int(), 0)
		t.Assert(all[0]["tag"].String(), "go")
		t.Assert(all[2]["tag"].String(), "orm")
		t.Assert(all[3]["id"].Int(), 2)
		t.Assert(all[3]["tag"].String(), "sql")

		all, err = db.GetAll(ctx, fmt.Sprintf(`
			SELECT je.value AS tag, COUNT(*) AS n
			FROM %s t, json_each(t.doc,'$.tags') je
			GROUP BY je.value HAVING n>=1 ORDER BY tag`, table,
		))
		t.AssertNil(err)
		t.Assert(len(all), 4)
		t.Assert(all[0]["tag"].String(), "db")
	})
}

// Test_SQLite_JSON1_SetAndInsertInUpdate tests json_set and json_insert as Raw update values.
func Test_SQLiteCgo_JSON1_SetAndInsertInUpdate(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		table := createJsonTable(t, "json_update")
		defer dropTable(table)

		_, err := db.Model(table).Data(g.Map{"doc": `{"name":"alice","age":30}`}).Insert()
		t.AssertNil(err)

		_, err = db.Model(table).Where("id", 1).Data(g.Map{
			"doc": gdb.Raw(`json_set(doc,'$.name','updated','$.city','beijing')`),
		}).Update()
		t.AssertNil(err)

		one, err := db.Model(table).Fields(
			`doc->>'$.name' AS name`,
			`doc->>'$.city' AS city`,
			`doc->>'$.age'  AS age`,
		).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["name"].String(), "updated")
		t.Assert(one["city"].String(), "beijing")
		t.Assert(one["age"].Int(), 30)

		_, err = db.Model(table).Where("id", 1).Data(g.Map{
			"doc": gdb.Raw(`json_insert(doc,'$.name','ignored','$.extra',7)`),
		}).Update()
		t.AssertNil(err)

		one, err = db.Model(table).Fields(
			`doc->>'$.name'  AS name`,
			`doc->>'$.extra' AS extra`,
		).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["name"].String(), "updated")
		t.Assert(one["extra"].Int(), 7)
	})
}

// =============================================================================
// Storage classes and table options
// =============================================================================

func createStorageClassTable(t *gtest.T, name string) string {
	table := featureTable(name)
	_, err := db.Exec(ctx, fmt.Sprintf(`CREATE TABLE %s (id INTEGER PRIMARY KEY, v)`, table))
	t.AssertNil(err)
	_, err = db.Model(table).Data(g.List{
		{"id": 1, "v": 42},
		{"id": 2, "v": 1.5},
		{"id": 3, "v": "hello"},
		{"id": 4, "v": []byte{1, 2, 3}},
		{"id": 5, "v": nil},
	}).Insert()
	t.AssertNil(err)
	return table
}

// Test_SQLite_Type_StorageClass_Typeof tests that a value keeps its own storage class
// in a column that has no declared type, and that a single-row read returns it intact.
func Test_SQLiteCgo_Type_StorageClass_Typeof(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		table := createStorageClassTable(t, "storage_typeof")
		defer dropTable(table)

		all, err := db.GetAll(ctx, fmt.Sprintf(`SELECT id, typeof(v) AS tv FROM %s ORDER BY id`, table))
		t.AssertNil(err)
		t.Assert(len(all), 5)
		t.Assert(all[0]["tv"].String(), "integer")
		t.Assert(all[1]["tv"].String(), "real")
		t.Assert(all[2]["tv"].String(), "text")
		t.Assert(all[3]["tv"].String(), "blob")
		t.Assert(all[4]["tv"].String(), "null")

		one, err := db.Model(table).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["v"].Int64(), 42)

		one, err = db.Model(table).Where("id", 2).One()
		t.AssertNil(err)
		t.Assert(one["v"].Float64(), 1.5)

		one, err = db.Model(table).Where("id", 3).One()
		t.AssertNil(err)
		t.Assert(one["v"].String(), "hello")

		one, err = db.Model(table).Where("id", 4).One()
		t.AssertNil(err)
		t.Assert(one["v"].Bytes(), []byte{1, 2, 3})

		one, err = db.Model(table).Where("id", 5).One()
		t.AssertNil(err)
		t.Assert(one["v"], nil)
	})
}

// Test_SQLite_Type_StorageClass_MixedRows tests that every row of a mixed storage class
// column keeps its own value when the whole column is read in one result set.
func Test_SQLiteCgo_Type_StorageClass_MixedRows(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		table := createStorageClassTable(t, "storage_mixed")
		defer dropTable(table)

		all, err := db.Model(table).OrderAsc("id").All()
		t.AssertNil(err)
		t.Assert(len(all), 5)

		t.Assert(all[0]["v"].Int64(), 42)
		t.Assert(all[1]["v"].Float64(), 1.5)
		t.Assert(all[2]["v"].String(), "hello")
		t.Assert(all[3]["v"].Bytes(), []byte{1, 2, 3})
		t.Assert(all[4]["v"], nil)
	})
}

// Test_SQLite_Type_WithoutRowid tests a WITHOUT ROWID table.
func Test_SQLiteCgo_Type_WithoutRowid(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		table := featureTable("without_rowid")
		_, err := db.Exec(ctx, fmt.Sprintf(`
			CREATE TABLE %s (
				code VARCHAR(45) NOT NULL PRIMARY KEY,
				name VARCHAR(45)
			) WITHOUT ROWID`, table,
		))
		t.AssertNil(err)
		defer dropTable(table)

		_, err = db.Model(table).Data(g.List{
			{"code": "b", "name": "beta"},
			{"code": "a", "name": "alpha"},
		}).Insert()
		t.AssertNil(err)

		all, err := db.Model(table).OrderAsc("code").All()
		t.AssertNil(err)
		t.Assert(len(all), 2)
		t.Assert(all[0]["code"].String(), "a")
		t.Assert(all[1]["code"].String(), "b")

		_, err = db.GetValue(ctx, fmt.Sprintf(`SELECT rowid FROM %s`, table))
		t.AssertNE(err, nil)
		t.Assert(gstr.Contains(err.Error(), "no such column: rowid"), true)

		_, err = db.Model(table).Data(g.Map{"code": "a", "name": "duplicate"}).Insert()
		t.AssertNE(err, nil)
		t.Assert(gstr.Contains(err.Error(), "UNIQUE constraint failed"), true)
	})
}

// Test_SQLite_Type_StrictTable tests a STRICT table rejecting a wrongly typed value.
func Test_SQLiteCgo_Type_StrictTable(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		table := featureTable("strict")
		_, err := db.Exec(ctx, fmt.Sprintf(`
			CREATE TABLE %s (
				id   INTEGER PRIMARY KEY,
				n    INT  NOT NULL,
				name TEXT NOT NULL
			) STRICT`, table,
		))
		t.AssertNil(err)
		defer dropTable(table)

		_, err = db.Model(table).Data(g.Map{"id": 1, "n": 7, "name": "alice"}).Insert()
		t.AssertNil(err)

		_, err = db.Model(table).Data(g.Map{"id": 2, "n": "not_an_int", "name": "bob"}).Insert()
		t.AssertNE(err, nil)
		t.Assert(gstr.Contains(err.Error(), "cannot store TEXT value in INT column"), true)

		count, err := db.Model(table).Count()
		t.AssertNil(err)
		t.Assert(count, 1)

		all, err := db.GetAll(ctx, fmt.Sprintf(`SELECT typeof(n) AS tn, typeof(name) AS ts FROM %s`, table))
		t.AssertNil(err)
		t.Assert(all[0]["tn"].String(), "integer")
		t.Assert(all[0]["ts"].String(), "text")
	})
}

// =============================================================================
// Date and time functions
// =============================================================================

// Test_SQLite_DateTime_RawInDataAndWhere tests datetime('now') as a Raw value in Data and Where.
func Test_SQLiteCgo_DateTime_RawInDataAndWhere(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		table := featureTable("datetime_raw")
		_, err := db.Exec(ctx, fmt.Sprintf(`
			CREATE TABLE %s (
				id INTEGER PRIMARY KEY,
				ts TEXT
			)`, table,
		))
		t.AssertNil(err)
		defer dropTable(table)

		_, err = db.Model(table).Data(g.Map{
			"id": 1,
			"ts": gdb.Raw(`datetime('now')`),
		}).Insert()
		t.AssertNil(err)

		one, err := db.Model(table).Where("id", 1).One()
		t.AssertNil(err)
		stored := gtime.NewFromStr(one["ts"].String())
		t.AssertNE(stored, nil)

		age, err := db.GetValue(ctx, fmt.Sprintf(
			`SELECT (julianday('now')-julianday(ts))*86400 FROM %s WHERE id=1`, table,
		))
		t.AssertNil(err)
		t.AssertGE(age.Float64(), -1.0)
		t.AssertLE(age.Float64(), 60.0)

		count, err := db.Model(table).Where(
			"ts < ?", gdb.Raw(`datetime('now','+1 day')`),
		).Count()
		t.AssertNil(err)
		t.Assert(count, 1)

		count, err = db.Model(table).Where(
			"ts < ?", gdb.Raw(`datetime('now','-1 day')`),
		).Count()
		t.AssertNil(err)
		t.Assert(count, 0)

		_, err = db.Model(table).Where("id", 1).Data(g.Map{
			"ts": gdb.Raw(`date(ts,'+1 day')`),
		}).Update()
		t.AssertNil(err)

		shifted, err := db.Model(table).Where("id", 1).Value("ts")
		t.AssertNil(err)
		t.Assert(shifted.String(), gtime.NewFromStr(one["ts"].String()).AddDate(0, 0, 1).Format("Y-m-d"))
	})
}

// Test_SQLite_DateTime_StrftimeJulianday tests strftime, julianday and date modifiers.
func Test_SQLiteCgo_DateTime_StrftimeJulianday(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		one, err := db.GetOne(ctx, `
			SELECT
				date('2024-01-31','+1 day')      AS next_day,
				date('2024-01-31','+1 month')    AS next_month,
				date('2024-03-15','start of month','-1 day') AS prev_month_end,
				julianday('2024-01-01')          AS jd,
				strftime('%Y-%m','2024-07-09')   AS ym,
				strftime('%s','2024-01-01')      AS epoch,
				strftime('%w','2024-01-01')      AS weekday`,
		)
		t.AssertNil(err)
		t.Assert(one["next_day"].String(), "2024-02-01")
		t.Assert(one["next_month"].String(), "2024-03-02")
		t.Assert(one["prev_month_end"].String(), "2024-02-29")
		t.Assert(one["jd"].Float64(), 2460310.5)
		t.Assert(one["ym"].String(), "2024-07")
		t.Assert(one["epoch"].Int64(), 1704067200)
		t.Assert(one["weekday"].Int(), 1)
	})

	gtest.C(t, func(t *gtest.T) {
		table := createInitTable()
		defer dropTable(table)

		all, err := db.Model(table).Fields(
			"id",
			`strftime('%Y',create_time) AS y`,
			`strftime('%Y-%m-%d',create_time,'+1 day') AS tomorrow`,
		).Where("id", 1).All()
		t.AssertNil(err)
		t.Assert(all[0]["y"].String(), "2018")
		t.Assert(all[0]["tomorrow"].String(), "2018-10-25")

		count, err := db.Model(table).Where(
			`strftime('%Y',create_time) = ?`, "2018",
		).Count()
		t.AssertNil(err)
		t.Assert(count, TableSize)
	})
}

// =============================================================================
// Built-in functions and query options
// =============================================================================

// Test_SQLite_GroupConcat tests GROUP_CONCAT with the default and a custom separator.
func Test_SQLiteCgo_GroupConcat(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		table := createInitTable()
		defer dropTable(table)

		one, err := db.Model(table).Fields(`GROUP_CONCAT(nickname) AS names`).Where("id<=3").One()
		t.AssertNil(err)
		t.Assert(one["names"].String(), "name_1,name_2,name_3")

		one, err = db.Model(table).Fields(`GROUP_CONCAT(nickname,' | ') AS names`).Where("id<=3").One()
		t.AssertNil(err)
		t.Assert(one["names"].String(), "name_1 | name_2 | name_3")

		all, err := db.Model(table).Fields(
			`id % 3 AS bucket`,
			`GROUP_CONCAT(id) AS ids`,
		).Group("bucket").OrderAsc("bucket").All()
		t.AssertNil(err)
		t.Assert(len(all), 3)
		t.Assert(all[0]["bucket"].Int(), 0)
		t.Assert(all[0]["ids"].String(), "3,6,9")
		t.Assert(all[1]["ids"].String(), "1,4,7,10")
	})
}

// Test_SQLite_IifAndCoalesce tests IIF and COALESCE over NULL values.
func Test_SQLiteCgo_IifAndCoalesce(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		table := createInitTable()
		defer dropTable(table)

		_, err := db.Model(table).Where("id>5").Data(g.Map{"nickname": nil}).Update()
		t.AssertNil(err)

		all, err := db.Model(table).Fields(
			"id",
			`COALESCE(nickname,'anonymous') AS display_name`,
			`IIF(id>5,'big','small')        AS bucket`,
			`IIF(nickname IS NULL,1,0)      AS is_null`,
		).OrderAsc("id").All()
		t.AssertNil(err)
		t.Assert(len(all), TableSize)

		t.Assert(all[0]["display_name"].String(), "name_1")
		t.Assert(all[0]["bucket"].String(), "small")
		t.Assert(all[0]["is_null"].Int(), 0)

		t.Assert(all[TableSize-1]["display_name"].String(), "anonymous")
		t.Assert(all[TableSize-1]["bucket"].String(), "big")
		t.Assert(all[TableSize-1]["is_null"].Int(), 1)

		count, err := db.Model(table).Where(`COALESCE(nickname,'anonymous') = ?`, "anonymous").Count()
		t.AssertNil(err)
		t.Assert(count, 5)
	})
}

// Test_SQLite_LastInsertRowid tests last_insert_rowid() on the connection of a transaction.
func Test_SQLiteCgo_LastInsertRowid(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		table := createInitTable()
		defer dropTable(table)

		err := db.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
			_, e := tx.Model(table).Data(g.Map{
				"passport": "rowid_1", "password": "pass", "nickname": "nick_1",
			}).Insert()
			t.AssertNil(e)

			last, e := tx.GetValue(`SELECT last_insert_rowid()`)
			t.AssertNil(e)
			t.Assert(last.Int(), TableSize+1)

			_, e = tx.Model(table).Data(g.Map{
				"passport": "rowid_2", "password": "pass", "nickname": "nick_2",
			}).Insert()
			t.AssertNil(e)

			last, e = tx.GetValue(`SELECT last_insert_rowid()`)
			t.AssertNil(e)
			t.Assert(last.Int(), TableSize+2)

			passport, e := tx.Model(table).Where("id", last).Value("passport")
			t.AssertNil(e)
			t.Assert(passport.String(), "rowid_2")
			return nil
		})
		t.AssertNil(err)

		count, err := db.Model(table).Count()
		t.AssertNil(err)
		t.Assert(count, TableSize+2)
	})
}

// Test_SQLite_OrderRandom tests OrderRandom emitting ORDER BY RANDOM().
func Test_SQLiteCgo_OrderRandom(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		table := createInitTable()
		defer dropTable(table)

		var all gdb.Result
		sqlArray, err := gdb.CatchSQL(ctx, func(ctx context.Context) error {
			var e error
			all, e = db.Ctx(ctx).Model(table).Fields("id").OrderRandom().All()
			return e
		})
		t.AssertNil(err)
		t.AssertGT(len(sqlArray), 0)
		t.Assert(gstr.Contains(sqlArray[len(sqlArray)-1], "ORDER BY RANDOM()"), true)

		t.Assert(len(all), TableSize)
		seen := make(map[int]bool)
		for _, record := range all {
			seen[record["id"].Int()] = true
		}
		t.Assert(len(seen), TableSize)
		for i := 1; i <= TableSize; i++ {
			t.Assert(seen[i], true)
		}

		one, err := db.Model(table).Fields("id").OrderRandom().Limit(1).One()
		t.AssertNil(err)
		t.AssertGE(one["id"].Int(), 1)
		t.AssertLE(one["id"].Int(), TableSize)
	})
}

// Test_SQLite_LimitNegativeOne tests LIMIT -1, which SQLite reads as no limit.
func Test_SQLiteCgo_LimitNegativeOne(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		table := createInitTable()
		defer dropTable(table)

		all, err := db.GetAll(ctx, fmt.Sprintf(`SELECT id FROM %s ORDER BY id LIMIT -1`, table))
		t.AssertNil(err)
		t.Assert(len(all), TableSize)

		all, err = db.GetAll(ctx, fmt.Sprintf(`SELECT id FROM %s ORDER BY id LIMIT -1 OFFSET 8`, table))
		t.AssertNil(err)
		t.Assert(len(all), 2)
		t.Assert(all[0]["id"].Int(), 9)
		t.Assert(all[1]["id"].Int(), 10)

		all, err = db.GetAll(ctx, fmt.Sprintf(`SELECT id FROM %s ORDER BY id LIMIT 0`, table))
		t.AssertNil(err)
		t.Assert(len(all), 0)
	})
}

// Test_SQLite_ExplainQueryPlan tests EXPLAIN QUERY PLAN returning a readable plan.
func Test_SQLiteCgo_ExplainQueryPlan(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		table := createInitTable()
		defer dropTable(table)

		all, err := db.GetAll(ctx, fmt.Sprintf(
			`EXPLAIN QUERY PLAN SELECT * FROM %s WHERE id=1`, table,
		))
		t.AssertNil(err)
		t.AssertGT(len(all), 0)
		t.Assert(gstr.Contains(all[0]["detail"].String(), table), true)
		t.Assert(gstr.Contains(all[0]["detail"].String(), "SEARCH"), true)
		t.Assert(gstr.Contains(all[0]["detail"].String(), "PRIMARY KEY"), true)

		all, err = db.GetAll(ctx, fmt.Sprintf(
			`EXPLAIN QUERY PLAN SELECT * FROM %s WHERE nickname='name_1'`, table,
		))
		t.AssertNil(err)
		t.AssertGT(len(all), 0)
		t.Assert(gstr.Contains(all[0]["detail"].String(), "SCAN"), true)
	})
}
