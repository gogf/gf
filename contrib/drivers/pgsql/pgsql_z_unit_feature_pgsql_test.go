// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package pgsql_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/text/gstr"
)

// =============================================================================
// Layer 3: JSONB operator integration tests
// =============================================================================

// Test_PgSQL_JSONB_Arrow_Operator tests the -> and ->> operators for JSONB field access.
func Test_PgSQL_JSONB_Arrow_Operator(t *testing.T) {
	table := fmt.Sprintf(`%s_%d`, TablePrefix+"jsonb_arrow", gtime.TimestampNano())
	if _, err := db.Exec(ctx, fmt.Sprintf(`
		CREATE TABLE %s (
			id bigserial PRIMARY KEY,
			data jsonb
		);`, table)); err != nil {
		gtest.Fatal(err)
	}
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		// Insert test data
		for i := 1; i <= 5; i++ {
			_, err := db.Model(table).Data(g.Map{
				"data": g.Map{
					"name": fmt.Sprintf("user_%d", i),
					"age":  20 + i,
					"tags": g.Slice{"go", "pgsql"},
				},
			}).Insert()
			t.AssertNil(err)
		}

		// -> returns jsonb (with quotes for strings)
		one, err := db.Model(table).Fields("data->'name' as name_json").Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["name_json"].String(), `"user_1"`)

		// ->> returns text (without quotes)
		one, err = db.Model(table).Fields("data->>'name' as name_text").Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["name_text"].String(), "user_1")

		// Nested -> for array element
		one, err = db.Model(table).Fields("data->'tags'->0 as first_tag").Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["first_tag"].String(), `"go"`)
	})
}

// Test_PgSQL_JSONB_Contains tests the @> (contains) and <@ (contained by) operators.
func Test_PgSQL_JSONB_Contains(t *testing.T) {
	table := fmt.Sprintf(`%s_%d`, TablePrefix+"jsonb_contains", gtime.TimestampNano())
	if _, err := db.Exec(ctx, fmt.Sprintf(`
		CREATE TABLE %s (
			id bigserial PRIMARY KEY,
			data jsonb
		);`, table)); err != nil {
		gtest.Fatal(err)
	}
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(table).Data(g.Map{"data": `{"name":"alice","role":"admin","level":5}`}).Insert()
		t.AssertNil(err)
		_, err = db.Model(table).Data(g.Map{"data": `{"name":"bob","role":"user","level":1}`}).Insert()
		t.AssertNil(err)
		_, err = db.Model(table).Data(g.Map{"data": `{"name":"charlie","role":"admin","level":3}`}).Insert()
		t.AssertNil(err)

		// @> contains: find admins
		all, err := db.Model(table).Where("data @> ?", `{"role":"admin"}`).OrderAsc("id").All()
		t.AssertNil(err)
		t.Assert(len(all), 2)
		t.Assert(all[0]["data"].Map()["name"], "alice")
		t.Assert(all[1]["data"].Map()["name"], "charlie")

		// <@ contained by
		count, err := db.Model(table).Where("? <@ data", `{"role":"user"}`).Count()
		t.AssertNil(err)
		t.Assert(count, 1)
	})
}

// Test_PgSQL_JSONB_Existence tests the ? (key existence) operator.
func Test_PgSQL_JSONB_Existence(t *testing.T) {
	table := fmt.Sprintf(`%s_%d`, TablePrefix+"jsonb_exist", gtime.TimestampNano())
	if _, err := db.Exec(ctx, fmt.Sprintf(`
		CREATE TABLE %s (
			id bigserial PRIMARY KEY,
			data jsonb
		);`, table)); err != nil {
		gtest.Fatal(err)
	}
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(table).Data(g.Map{"data": `{"name":"alice","email":"a@b.com"}`}).Insert()
		t.AssertNil(err)
		_, err = db.Model(table).Data(g.Map{"data": `{"name":"bob"}`}).Insert()
		t.AssertNil(err)

		// jsonb_exists is the function form of the ? key-existence operator.
		// We use it here because GoFrame's DoFilter replaces bare ? with $N placeholders.
		result, err := db.GetValue(ctx, fmt.Sprintf(
			`SELECT COUNT(*) FROM %s WHERE jsonb_exists(data, 'email')`, table,
		))
		t.AssertNil(err)
		t.Assert(result.Int(), 1)
	})
}

// Test_PgSQL_JSONB_Path tests jsonb_path_query_first (PG12+).
func Test_PgSQL_JSONB_Path(t *testing.T) {
	table := fmt.Sprintf(`%s_%d`, TablePrefix+"jsonb_path", gtime.TimestampNano())
	if _, err := db.Exec(ctx, fmt.Sprintf(`
		CREATE TABLE %s (
			id bigserial PRIMARY KEY,
			data jsonb
		);`, table)); err != nil {
		gtest.Fatal(err)
	}
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(table).Data(g.Map{
			"data": `{"users":[{"name":"alice","age":30},{"name":"bob","age":25}]}`,
		}).Insert()
		t.AssertNil(err)

		// jsonb_path_query_first
		one, err := db.Model(table).Fields(
			`jsonb_path_query_first(data, '$.users[0].name') as first_name`,
		).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["first_name"].String(), `"alice"`)
	})
}

// Test_PgSQL_JSONB_Concat tests the || (concatenation/merge) operator.
func Test_PgSQL_JSONB_Concat(t *testing.T) {
	table := fmt.Sprintf(`%s_%d`, TablePrefix+"jsonb_concat", gtime.TimestampNano())
	if _, err := db.Exec(ctx, fmt.Sprintf(`
		CREATE TABLE %s (
			id bigserial PRIMARY KEY,
			data jsonb
		);`, table)); err != nil {
		gtest.Fatal(err)
	}
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(table).Data(g.Map{"data": `{"name":"alice"}`}).Insert()
		t.AssertNil(err)

		// || merge operator via raw SQL update
		_, err = db.Exec(ctx, fmt.Sprintf(
			`UPDATE %s SET data = data || '{"role":"admin"}' WHERE id = 1`, table,
		))
		t.AssertNil(err)

		one, err := db.Model(table).Where("id", 1).One()
		t.AssertNil(err)
		m := one["data"].Map()
		t.Assert(m["name"], "alice")
		t.Assert(m["role"], "admin")
	})
}

// Test_PgSQL_JSONB_Remove tests the - (key removal) operator.
func Test_PgSQL_JSONB_Remove(t *testing.T) {
	table := fmt.Sprintf(`%s_%d`, TablePrefix+"jsonb_remove", gtime.TimestampNano())
	if _, err := db.Exec(ctx, fmt.Sprintf(`
		CREATE TABLE %s (
			id bigserial PRIMARY KEY,
			data jsonb
		);`, table)); err != nil {
		gtest.Fatal(err)
	}
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(table).Data(g.Map{"data": `{"name":"alice","age":30,"temp":"remove_me"}`}).Insert()
		t.AssertNil(err)

		// - key removal operator
		_, err = db.Exec(ctx, fmt.Sprintf(
			`UPDATE %s SET data = data - 'temp' WHERE id = 1`, table,
		))
		t.AssertNil(err)

		one, err := db.Model(table).Where("id", 1).One()
		t.AssertNil(err)
		m := one["data"].Map()
		t.Assert(m["name"], "alice")
		t.Assert(m["age"], 30)
		_, hasTmp := m["temp"]
		t.Assert(hasTmp, false)
	})
}

// Test_PgSQL_JSONB_Agg tests jsonb_agg and jsonb_object_agg aggregate functions.
func Test_PgSQL_JSONB_Agg(t *testing.T) {
	table := fmt.Sprintf(`%s_%d`, TablePrefix+"jsonb_agg", gtime.TimestampNano())
	if _, err := db.Exec(ctx, fmt.Sprintf(`
		CREATE TABLE %s (
			id bigserial PRIMARY KEY,
			name varchar(50),
			score int
		);`, table)); err != nil {
		gtest.Fatal(err)
	}
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		for i := 1; i <= 3; i++ {
			_, err := db.Model(table).Data(g.Map{
				"name":  fmt.Sprintf("user_%d", i),
				"score": i * 10,
			}).Insert()
			t.AssertNil(err)
		}

		// jsonb_agg: aggregate names into a JSON array
		one, err := db.Model(table).Fields("jsonb_agg(name) as names").One()
		t.AssertNil(err)
		names := one["names"].Strings()
		t.Assert(len(names), 3)

		// jsonb_object_agg: aggregate into key-value pairs
		one, err = db.Model(table).Fields("jsonb_object_agg(name, score) as scores").One()
		t.AssertNil(err)
		scores := one["scores"].Map()
		t.Assert(scores["user_1"], 10)
		t.Assert(scores["user_2"], 20)
		t.Assert(scores["user_3"], 30)
	})
}

// =============================================================================
// Layer 3: RETURNING clause tests
// =============================================================================

// Test_PgSQL_Returning_Insert tests INSERT ... RETURNING.
func Test_PgSQL_Returning_Insert(t *testing.T) {
	table := fmt.Sprintf(`%s_%d`, TablePrefix+"returning_ins", gtime.TimestampNano())
	if _, err := db.Exec(ctx, fmt.Sprintf(`
		CREATE TABLE %s (
			id bigserial PRIMARY KEY,
			name varchar(100),
			created_at timestamp DEFAULT now()
		);`, table)); err != nil {
		gtest.Fatal(err)
	}
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		// InsertAndGetId uses RETURNING internally for PgSQL
		id, err := db.Model(table).Data(g.Map{"name": "alice"}).InsertAndGetId()
		t.AssertNil(err)
		t.Assert(id, 1)

		id, err = db.Model(table).Data(g.Map{"name": "bob"}).InsertAndGetId()
		t.AssertNil(err)
		t.Assert(id, 2)

		// Verify data
		count, err := db.Model(table).Count()
		t.AssertNil(err)
		t.Assert(count, 2)
	})
}

// Test_PgSQL_Returning_Insert_Batch tests batch INSERT ... RETURNING.
func Test_PgSQL_Returning_Insert_Batch(t *testing.T) {
	table := fmt.Sprintf(`%s_%d`, TablePrefix+"returning_batch", gtime.TimestampNano())
	if _, err := db.Exec(ctx, fmt.Sprintf(`
		CREATE TABLE %s (
			id bigserial PRIMARY KEY,
			name varchar(100)
		);`, table)); err != nil {
		gtest.Fatal(err)
	}
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		data := g.List{
			{"name": "alice"},
			{"name": "bob"},
			{"name": "charlie"},
		}
		result, err := db.Model(table).Data(data).Insert()
		t.AssertNil(err)
		n, _ := result.RowsAffected()
		t.Assert(n, 3)

		// Verify all inserted
		all, err := db.Model(table).OrderAsc("id").All()
		t.AssertNil(err)
		t.Assert(len(all), 3)
		t.Assert(all[0]["name"].String(), "alice")
		t.Assert(all[2]["name"].String(), "charlie")
	})
}

// Test_PgSQL_Returning_Upsert tests INSERT ... ON CONFLICT ... RETURNING (Save).
func Test_PgSQL_Returning_Upsert(t *testing.T) {
	table := fmt.Sprintf(`%s_%d`, TablePrefix+"returning_upsert", gtime.TimestampNano())
	if _, err := db.Exec(ctx, fmt.Sprintf(`
		CREATE TABLE %s (
			id bigserial PRIMARY KEY,
			name varchar(100),
			value int DEFAULT 0
		);`, table)); err != nil {
		gtest.Fatal(err)
	}
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		// First insert
		_, err := db.Model(table).Data(g.Map{
			"id": 1, "name": "alice", "value": 10,
		}).Save()
		t.AssertNil(err)

		one, err := db.Model(table).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["name"].String(), "alice")
		t.Assert(one["value"].Int(), 10)

		// Upsert: update existing
		_, err = db.Model(table).Data(g.Map{
			"id": 1, "name": "alice_updated", "value": 20,
		}).Save()
		t.AssertNil(err)

		one, err = db.Model(table).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["name"].String(), "alice_updated")
		t.Assert(one["value"].Int(), 20)

		// Only 1 row total
		count, err := db.Model(table).Count()
		t.AssertNil(err)
		t.Assert(count, 1)
	})
}

// Test_PgSQL_Returning_CatchSQL verifies CatchSQL captures INSERT ... RETURNING.
// Regression for a bug where pgsql/gaussdb DoExec bypassed Core.DoQuery and
// silently dropped the SQL from CatchSQLManager on InsertAndGetId.
func Test_PgSQL_Returning_CatchSQL(t *testing.T) {
	table := fmt.Sprintf(`%s_%d`, TablePrefix+"returning_catchsql", gtime.TimestampNano())
	if _, err := db.Exec(ctx, fmt.Sprintf(`
		CREATE TABLE %s (
			id bigserial PRIMARY KEY,
			name varchar(100)
		);`, table)); err != nil {
		gtest.Fatal(err)
	}
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		sqlArray, err := gdb.CatchSQL(ctx, func(ctx context.Context) error {
			_, e := db.Ctx(ctx).Model(table).Data(g.Map{"name": "alice"}).InsertAndGetId()
			return e
		})
		t.AssertNil(err)
		t.AssertGT(len(sqlArray), 0)
		// The captured SQL must contain the RETURNING clause.
		t.Assert(gstr.Contains(sqlArray[len(sqlArray)-1], `RETURNING "id"`), true)
		// Insert must have executed (CatchSQL uses DoCommit=true).
		count, err := db.Model(table).Count()
		t.AssertNil(err)
		t.Assert(count, 1)
	})
}

// Test_PgSQL_Returning_ToSQL verifies ToSQL captures without executing.
// Regression for a bug where pgsql/gaussdb DoExec bypassed Core.DoQuery and
// executed the INSERT even when ToSQL (DoCommit=false) was active.
func Test_PgSQL_Returning_ToSQL(t *testing.T) {
	table := fmt.Sprintf(`%s_%d`, TablePrefix+"returning_tosql", gtime.TimestampNano())
	if _, err := db.Exec(ctx, fmt.Sprintf(`
		CREATE TABLE %s (
			id bigserial PRIMARY KEY,
			name varchar(100)
		);`, table)); err != nil {
		gtest.Fatal(err)
	}
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		before, err := db.Model(table).Count()
		t.AssertNil(err)

		capturedSql, err := gdb.ToSQL(ctx, func(ctx context.Context) error {
			_, e := db.Ctx(ctx).Model(table).Data(g.Map{"name": "bob"}).InsertAndGetId()
			return e
		})
		t.AssertNil(err)
		t.AssertNE(capturedSql, "")
		// Row count must be unchanged — ToSQL must NOT execute.
		after, err := db.Model(table).Count()
		t.AssertNil(err)
		t.Assert(after, before)
	})
}

// =============================================================================
// Layer 3: CTE (Common Table Expression) tests
// =============================================================================

// Test_PgSQL_CTE_Basic tests basic WITH ... AS query.
func Test_PgSQL_CTE_Basic(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		sql := fmt.Sprintf(`
			WITH top5 AS (
				SELECT id, nickname FROM %s ORDER BY id ASC LIMIT 5
			)
			SELECT * FROM top5 WHERE id > 2
		`, table)
		all, err := db.GetAll(ctx, sql)
		t.AssertNil(err)
		t.Assert(len(all), 3) // ids 3, 4, 5
		t.Assert(all[0]["id"].Int(), 3)
	})
}

// Test_PgSQL_CTE_Recursive tests recursive CTE for hierarchical data.
func Test_PgSQL_CTE_Recursive(t *testing.T) {
	table := fmt.Sprintf(`%s_%d`, TablePrefix+"cte_tree", gtime.TimestampNano())
	if _, err := db.Exec(ctx, fmt.Sprintf(`
		CREATE TABLE %s (
			id int PRIMARY KEY,
			parent_id int,
			name varchar(50)
		);`, table)); err != nil {
		gtest.Fatal(err)
	}
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		// Build a tree:  1 -> 2 -> 4
		//                1 -> 3
		data := g.List{
			{"id": 1, "parent_id": nil, "name": "root"},
			{"id": 2, "parent_id": 1, "name": "child_a"},
			{"id": 3, "parent_id": 1, "name": "child_b"},
			{"id": 4, "parent_id": 2, "name": "grandchild"},
		}
		_, err := db.Model(table).Data(data).Insert()
		t.AssertNil(err)

		// Recursive CTE to get all descendants of node 1
		sql := fmt.Sprintf(`
			WITH RECURSIVE tree AS (
				SELECT id, parent_id, name, 0 AS depth FROM %s WHERE id = 1
				UNION ALL
				SELECT t.id, t.parent_id, t.name, tree.depth + 1
				FROM %s t JOIN tree ON t.parent_id = tree.id
			)
			SELECT * FROM tree ORDER BY depth, id
		`, table, table)
		all, err := db.GetAll(ctx, sql)
		t.AssertNil(err)
		t.Assert(len(all), 4)
		t.Assert(all[0]["name"].String(), "root")
		t.Assert(all[0]["depth"].Int(), 0)
		t.Assert(all[3]["name"].String(), "grandchild")
		t.Assert(all[3]["depth"].Int(), 2)
	})
}

// Test_PgSQL_CTE_Update tests CTE used with UPDATE (writable CTE).
func Test_PgSQL_CTE_Update(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		// Use CTE to identify rows, then update them
		sql := fmt.Sprintf(`
			WITH to_update AS (
				SELECT id FROM %s WHERE id <= 3
			)
			UPDATE %s SET nickname = 'updated'
			WHERE id IN (SELECT id FROM to_update)
		`, table, table)
		_, err := db.Exec(ctx, sql)
		t.AssertNil(err)

		// Verify
		all, err := db.Model(table).Where("nickname", "updated").OrderAsc("id").All()
		t.AssertNil(err)
		t.Assert(len(all), 3)
		t.Assert(all[0]["id"].Int(), 1)
		t.Assert(all[2]["id"].Int(), 3)
	})
}

// Test_PgSQL_CTE_Delete tests CTE used with DELETE (writable CTE).
func Test_PgSQL_CTE_Delete(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		// Use CTE to identify rows, then delete them
		sql := fmt.Sprintf(`
			WITH to_delete AS (
				SELECT id FROM %s WHERE id > 8
			)
			DELETE FROM %s WHERE id IN (SELECT id FROM to_delete)
		`, table, table)
		_, err := db.Exec(ctx, sql)
		t.AssertNil(err)

		count, err := db.Model(table).Count()
		t.AssertNil(err)
		t.Assert(count, 8)
	})
}

// =============================================================================
// Layer 3: Window function tests
// =============================================================================

// Test_PgSQL_Window_RowNumber tests ROW_NUMBER() window function.
func Test_PgSQL_Window_RowNumber(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		sql := fmt.Sprintf(`
			SELECT id, nickname, ROW_NUMBER() OVER (ORDER BY id DESC) as rn
			FROM %s
		`, table)
		all, err := db.GetAll(ctx, sql)
		t.AssertNil(err)
		t.Assert(len(all), TableSize)
		// id=10 should have rn=1 (ordered desc)
		t.Assert(all[0]["id"].Int(), 10)
		t.Assert(all[0]["rn"].Int(), 1)
	})
}

// Test_PgSQL_Window_RankDenseRank tests RANK() and DENSE_RANK() window functions.
func Test_PgSQL_Window_RankDenseRank(t *testing.T) {
	table := fmt.Sprintf(`%s_%d`, TablePrefix+"window_rank", gtime.TimestampNano())
	if _, err := db.Exec(ctx, fmt.Sprintf(`
		CREATE TABLE %s (
			id bigserial PRIMARY KEY,
			dept varchar(20),
			salary int
		);`, table)); err != nil {
		gtest.Fatal(err)
	}
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		data := g.List{
			{"dept": "eng", "salary": 100},
			{"dept": "eng", "salary": 100},
			{"dept": "eng", "salary": 200},
			{"dept": "sales", "salary": 150},
			{"dept": "sales", "salary": 250},
		}
		_, err := db.Model(table).Data(data).Insert()
		t.AssertNil(err)

		sql := fmt.Sprintf(`
			SELECT dept, salary,
				RANK() OVER (PARTITION BY dept ORDER BY salary DESC) as rnk,
				DENSE_RANK() OVER (PARTITION BY dept ORDER BY salary DESC) as dense_rnk
			FROM %s ORDER BY dept, salary DESC
		`, table)
		all, err := db.GetAll(ctx, sql)
		t.AssertNil(err)
		t.Assert(len(all), 5)

		// eng: 200(rank=1), 100(rank=2), 100(rank=2)
		t.Assert(all[0]["salary"].Int(), 200)
		t.Assert(all[0]["rnk"].Int(), 1)
		t.Assert(all[1]["salary"].Int(), 100)
		t.Assert(all[1]["rnk"].Int(), 2)
		t.Assert(all[1]["dense_rnk"].Int(), 2)
	})
}

// Test_PgSQL_Window_LagLead tests LAG() and LEAD() window functions.
func Test_PgSQL_Window_LagLead(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		sql := fmt.Sprintf(`
			SELECT id,
				LAG(id, 1) OVER (ORDER BY id) as prev_id,
				LEAD(id, 1) OVER (ORDER BY id) as next_id
			FROM %s ORDER BY id
		`, table)
		all, err := db.GetAll(ctx, sql)
		t.AssertNil(err)
		t.Assert(len(all), TableSize)
		// First row: prev_id = NULL
		t.Assert(all[0]["prev_id"].IsNil() || all[0]["prev_id"].IsEmpty(), true)
		t.Assert(all[0]["next_id"].Int(), 2)
		// Last row: next_id = NULL
		t.Assert(all[9]["prev_id"].Int(), 9)
		t.Assert(all[9]["next_id"].IsNil() || all[9]["next_id"].IsEmpty(), true)
	})
}

// Test_PgSQL_Window_SumOver tests SUM() OVER (cumulative sum).
func Test_PgSQL_Window_SumOver(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		sql := fmt.Sprintf(`
			SELECT id, SUM(id) OVER (ORDER BY id) as running_total
			FROM %s ORDER BY id
		`, table)
		all, err := db.GetAll(ctx, sql)
		t.AssertNil(err)
		t.Assert(len(all), TableSize)
		// running_total: 1, 3, 6, 10, 15, 21, 28, 36, 45, 55
		t.Assert(all[0]["running_total"].Int(), 1)
		t.Assert(all[4]["running_total"].Int(), 15) // 1+2+3+4+5
		t.Assert(all[9]["running_total"].Int(), 55) // sum(1..10)
	})
}

// =============================================================================
// Layer 3: Array operations tests
// =============================================================================

// Test_PgSQL_Array_ANY tests the ANY() array operator.
func Test_PgSQL_Array_ANY(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		// Use ANY with array literal
		all, err := db.GetAll(ctx, fmt.Sprintf(
			`SELECT * FROM %s WHERE id = ANY(ARRAY[2, 5, 8]) ORDER BY id`, table,
		))
		t.AssertNil(err)
		t.Assert(len(all), 3)
		t.Assert(all[0]["id"].Int(), 2)
		t.Assert(all[1]["id"].Int(), 5)
		t.Assert(all[2]["id"].Int(), 8)
	})
}

// Test_PgSQL_Array_Contains tests array @> and <@ operators.
func Test_PgSQL_Array_Contains(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		// Insert data with array columns
		_, err := db.Exec(ctx, fmt.Sprintf(
			`UPDATE %s SET favorite_movie = ARRAY['movie_a', 'movie_b'] WHERE id = 1`, table,
		))
		t.AssertNil(err)
		_, err = db.Exec(ctx, fmt.Sprintf(
			`UPDATE %s SET favorite_movie = ARRAY['movie_b', 'movie_c'] WHERE id = 2`, table,
		))
		t.AssertNil(err)

		// @> contains: find rows where array contains 'movie_a'
		// Cast ARRAY literal to varchar[] to match the column type (varchar[]).
		all, err := db.GetAll(ctx, fmt.Sprintf(
			`SELECT id FROM %s WHERE favorite_movie @> ARRAY['movie_a']::varchar[] ORDER BY id`, table,
		))
		t.AssertNil(err)
		t.Assert(len(all), 1)
		t.Assert(all[0]["id"].Int(), 1)

		// @> contains: find rows where array contains 'movie_b' (both rows)
		all, err = db.GetAll(ctx, fmt.Sprintf(
			`SELECT id FROM %s WHERE favorite_movie @> ARRAY['movie_b']::varchar[] ORDER BY id`, table,
		))
		t.AssertNil(err)
		t.Assert(len(all), 2)
	})
}

// Test_PgSQL_Array_Unnest tests array_agg and unnest functions.
func Test_PgSQL_Array_Unnest(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		// array_agg: aggregate column values into array
		one, err := db.Model(table).Fields("array_agg(id ORDER BY id) as ids").Where("id <= 5").One()
		t.AssertNil(err)
		ids := one["ids"].Ints()
		t.Assert(len(ids), 5)
		t.Assert(ids[0], 1)
		t.Assert(ids[4], 5)
	})
}

// =============================================================================
// Layer 3: Full-text search tests
// =============================================================================

// Test_PgSQL_FullText_Search tests tsvector and tsquery full-text search.
func Test_PgSQL_FullText_Search(t *testing.T) {
	table := fmt.Sprintf(`%s_%d`, TablePrefix+"fts", gtime.TimestampNano())
	if _, err := db.Exec(ctx, fmt.Sprintf(`
		CREATE TABLE %s (
			id bigserial PRIMARY KEY,
			title varchar(200),
			body text,
			tsv tsvector
		);`, table)); err != nil {
		gtest.Fatal(err)
	}
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		// Insert documents
		docs := g.List{
			{"title": "PostgreSQL is great", "body": "PostgreSQL is an advanced open source database"},
			{"title": "GoFrame ORM", "body": "GoFrame provides powerful ORM features for Go"},
			{"title": "Database testing", "body": "Testing database drivers is important for PostgreSQL"},
		}
		for _, doc := range docs {
			_, err := db.Exec(ctx, fmt.Sprintf(
				`INSERT INTO %s (title, body, tsv) VALUES ('%s', '%s', to_tsvector('english', '%s %s'))`,
				table, doc["title"], doc["body"], doc["title"], doc["body"],
			))
			t.AssertNil(err)
		}

		// Full-text search for 'PostgreSQL'
		all, err := db.GetAll(ctx, fmt.Sprintf(
			`SELECT id, title FROM %s WHERE tsv @@ to_tsquery('english', 'PostgreSQL') ORDER BY id`, table,
		))
		t.AssertNil(err)
		t.Assert(len(all), 2) // docs 1 and 3
		t.Assert(all[0]["title"].String(), "PostgreSQL is great")

		// Full-text search with AND
		all, err = db.GetAll(ctx, fmt.Sprintf(
			`SELECT id, title FROM %s WHERE tsv @@ to_tsquery('english', 'database & testing') ORDER BY id`, table,
		))
		t.AssertNil(err)
		t.Assert(len(all), 1)
		t.Assert(all[0]["title"].String(), "Database testing")
	})
}

// =============================================================================
// Layer 3: Advanced PgSQL features
// =============================================================================

// Test_PgSQL_Generate_Series tests generate_series() table-valued function.
func Test_PgSQL_Generate_Series(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		all, err := db.GetAll(ctx, `SELECT generate_series(1, 5) as n`)
		t.AssertNil(err)
		t.Assert(len(all), 5)
		for i, row := range all {
			t.Assert(row["n"].Int(), i+1)
		}
	})
}

// Test_PgSQL_Lateral_Join tests LATERAL join.
func Test_PgSQL_Lateral_Join(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		// LATERAL join: for each row, get the next row's id
		sql := fmt.Sprintf(`
			SELECT t1.id, t2.next_id
			FROM %s t1,
			LATERAL (SELECT id as next_id FROM %s WHERE id = t1.id + 1) t2
			ORDER BY t1.id
		`, table, table)
		all, err := db.GetAll(ctx, sql)
		t.AssertNil(err)
		t.Assert(len(all), 9) // ids 1-9 (id=10 has no next)
		t.Assert(all[0]["id"].Int(), 1)
		t.Assert(all[0]["next_id"].Int(), 2)
	})
}

// Test_PgSQL_Distinct_On tests DISTINCT ON (PgSQL-specific).
func Test_PgSQL_Distinct_On(t *testing.T) {
	table := fmt.Sprintf(`%s_%d`, TablePrefix+"distinct_on", gtime.TimestampNano())
	if _, err := db.Exec(ctx, fmt.Sprintf(`
		CREATE TABLE %s (
			id bigserial PRIMARY KEY,
			category varchar(20),
			value int
		);`, table)); err != nil {
		gtest.Fatal(err)
	}
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		data := g.List{
			{"category": "A", "value": 10},
			{"category": "A", "value": 20},
			{"category": "B", "value": 30},
			{"category": "B", "value": 40},
			{"category": "C", "value": 50},
		}
		_, err := db.Model(table).Data(data).Insert()
		t.AssertNil(err)

		// DISTINCT ON: get first row per category (ordered by value desc = max)
		all, err := db.GetAll(ctx, fmt.Sprintf(`
			SELECT DISTINCT ON (category) category, value
			FROM %s ORDER BY category, value DESC
		`, table))
		t.AssertNil(err)
		t.Assert(len(all), 3)
		t.Assert(all[0]["category"].String(), "A")
		t.Assert(all[0]["value"].Int(), 20)
		t.Assert(all[1]["category"].String(), "B")
		t.Assert(all[1]["value"].Int(), 40)
		t.Assert(all[2]["category"].String(), "C")
		t.Assert(all[2]["value"].Int(), 50)
	})
}

// Test_PgSQL_Explain tests EXPLAIN (ANALYZE) for query plan inspection.
func Test_PgSQL_Explain(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		// Just verify EXPLAIN doesn't error out
		all, err := db.GetAll(ctx, fmt.Sprintf(`EXPLAIN SELECT * FROM %s WHERE id = 1`, table))
		t.AssertNil(err)
		t.AssertGT(len(all), 0)
	})
}

// Test_PgSQL_Coalesce tests COALESCE function with NULL handling.
func Test_PgSQL_Coalesce(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		// Set some nicknames to NULL
		_, err := db.Exec(ctx, fmt.Sprintf(`UPDATE %s SET nickname = NULL WHERE id > 5`, table))
		t.AssertNil(err)

		all, err := db.GetAll(ctx, fmt.Sprintf(
			`SELECT id, COALESCE(nickname, 'anonymous') as display_name FROM %s ORDER BY id`, table,
		))
		t.AssertNil(err)
		t.Assert(len(all), TableSize)
		t.Assert(all[0]["display_name"].String(), "name_1")
		t.Assert(all[9]["display_name"].String(), "anonymous")
	})
}

// Test_PgSQL_String_Agg tests string_agg aggregate function.
func Test_PgSQL_String_Agg(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		one, err := db.Model(table).Fields("string_agg(nickname, ',' ORDER BY id) as names").Where("id <= 3").One()
		t.AssertNil(err)
		t.Assert(one["names"].String(), "name_1,name_2,name_3")
	})
}

// Test_PgSQL_Filter_Clause tests the FILTER (WHERE ...) clause on aggregates.
func Test_PgSQL_Filter_Clause(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		sql := fmt.Sprintf(`
			SELECT
				COUNT(*) as total,
				COUNT(*) FILTER (WHERE id <= 5) as first_half,
				COUNT(*) FILTER (WHERE id > 5) as second_half
			FROM %s
		`, table)
		one, err := db.GetOne(ctx, sql)
		t.AssertNil(err)
		t.Assert(one["total"].Int(), 10)
		t.Assert(one["first_half"].Int(), 5)
		t.Assert(one["second_half"].Int(), 5)
	})
}
