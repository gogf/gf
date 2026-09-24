// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package oceanbase_test

import (
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/text/gstr"
)

// featureTable returns a unique table name for a feature test.
func featureTable(suffix string) string {
	return fmt.Sprintf(`ob_%s_%d`, suffix, gtime.TimestampNano())
}

// explainPlan runs EXPLAIN and joins the OceanBase plan rows into a single string.
func explainPlan(t *gtest.T, query string) string {
	all, err := db.GetAll(ctx, "EXPLAIN "+query)
	t.AssertNil(err)
	t.AssertGT(len(all), 0)
	var b []string
	for _, row := range all {
		b = append(b, row["Query Plan"].String())
	}
	return gstr.Join(b, "\n")
}

// =============================================================================
// Auto-increment identity reporting
// =============================================================================

// Test_OceanBase_LastInsertId_AutoIncrement verifies generated ids are reported back.
func Test_OceanBase_LastInsertId_AutoIncrement(t *testing.T) {
	table := createTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Data(g.Map{"passport": "u1", "nickname": "n1"}).Insert()
		t.AssertNil(err)
		id, err := result.LastInsertId()
		t.AssertNil(err)
		t.Assert(id, 1)

		id, err = db.Model(table).Data(g.Map{"passport": "u2", "nickname": "n2"}).InsertAndGetId()
		t.AssertNil(err)
		t.Assert(id, 2)

		one, err := db.Model(table).Where("id", id).One()
		t.AssertNil(err)
		t.Assert(one["passport"].String(), "u2")
	})
}

// Test_OceanBase_LastInsertId_ExplicitIdReportsZero pins that OceanBase reports 0 when the
// primary key is supplied by the caller, even though the row is written correctly.
func Test_OceanBase_LastInsertId_ExplicitIdReportsZero(t *testing.T) {
	table := createTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Data(g.Map{"id": 100, "passport": "u100", "nickname": "n100"}).Insert()
		t.AssertNil(err)
		id, err := result.LastInsertId()
		t.AssertNil(err)
		t.Assert(id, 0)

		id, err = db.Model(table).Data(g.Map{"id": 200, "passport": "u200", "nickname": "n200"}).InsertAndGetId()
		t.AssertNil(err)
		t.Assert(id, 0)

		// Both rows are persisted despite the reported id being 0.
		all, err := db.Model(table).OrderAsc("id").All()
		t.AssertNil(err)
		t.Assert(len(all), 2)
		t.Assert(all[0]["id"].Int(), 100)
		t.Assert(all[1]["id"].Int(), 200)
	})
}

// =============================================================================
// OceanBase session variables
// =============================================================================

// Test_OceanBase_SessionVar_ObTimeouts_ReadSetAndEnforce verifies the ob_*_timeout
// session variables are readable, settable on a pinned connection, and that a small
// ob_query_timeout actually aborts a long running statement.
func Test_OceanBase_SessionVar_ObTimeouts_ReadSetAndEnforce(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		tx, err := db.Begin(ctx)
		t.AssertNil(err)
		defer tx.Rollback()

		vars := g.SliceStr{"ob_query_timeout", "ob_trx_timeout", "ob_trx_idle_timeout"}
		originals := make(map[string]int64)
		for _, name := range vars {
			v, err := tx.GetValue("SELECT @@" + name)
			t.AssertNil(err)
			t.AssertGT(v.Int64(), int64(0))
			originals[name] = v.Int64()
		}
		defer func() {
			for name, value := range originals {
				tx.Exec(fmt.Sprintf("SET SESSION %s = %d", name, value))
			}
		}()

		for _, name := range vars {
			_, err = tx.Exec(fmt.Sprintf("SET SESSION %s = %d", name, 1234000))
			t.AssertNil(err)
			v, err := tx.GetValue("SELECT @@" + name)
			t.AssertNil(err)
			t.Assert(v.Int64(), int64(1234000))
			_, err = tx.Exec(fmt.Sprintf("SET SESSION %s = %d", name, originals[name]))
			t.AssertNil(err)
			v, err = tx.GetValue("SELECT @@" + name)
			t.AssertNil(err)
			t.Assert(v.Int64(), originals[name])
		}

		_, err = tx.Exec("SET SESSION ob_query_timeout = 100000")
		t.AssertNil(err)
		_, err = tx.GetValue("SELECT SLEEP(2)")
		t.AssertNE(err, nil)
		t.Assert(gstr.Contains(err.Error(), "4012"), true)
		t.Assert(gstr.Contains(err.Error(), "ob_query_timeout"), true)

		_, err = tx.Exec(fmt.Sprintf("SET SESSION ob_query_timeout = %d", originals["ob_query_timeout"]))
		t.AssertNil(err)
		_, err = tx.GetValue("SELECT SLEEP(1)")
		t.AssertNil(err)
	})
}

// Test_OceanBase_Hint_QueryTimeoutAndError4012 verifies the statement level QUERY_TIMEOUT
// hint aborts the statement without touching the session variable, and that OceanBase
// error 4012 survives gf's error wrapping.
func Test_OceanBase_Hint_QueryTimeoutAndError4012(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		before, err := db.GetValue(ctx, "SELECT @@ob_query_timeout")
		t.AssertNil(err)

		_, err = db.GetValue(ctx, "SELECT /*+ QUERY_TIMEOUT(200000) */ SLEEP(2)")
		t.AssertNE(err, nil)
		t.Assert(gstr.Contains(err.Error(), "200000(us)"), true)
		t.Assert(gerror.Code(err), gcode.CodeDbOperationError)

		// The OceanBase error number, SQLSTATE and message reach the caller unchanged.
		cause := gerror.Cause(err)
		t.AssertNE(cause, nil)
		t.Assert(gstr.HasPrefix(cause.Error(), "Error 4012 (HY000): Timeout"), true)
		t.Assert(gstr.Contains(cause.Error(), "ob_query_timeout"), true)

		after, err := db.GetValue(ctx, "SELECT @@ob_query_timeout")
		t.AssertNil(err)
		t.Assert(after.Int64(), before.Int64())

		v, err := db.GetValue(ctx, "SELECT /*+ QUERY_TIMEOUT(10000000) */ 1")
		t.AssertNil(err)
		t.Assert(v.Int(), 1)
	})
}

// =============================================================================
// Optimizer hints and EXPLAIN
// =============================================================================

// Test_OceanBase_Hint_IndexForcesRangeScan verifies /*+ INDEX(...) */ changes the plan.
func Test_OceanBase_Hint_IndexForcesRangeScan(t *testing.T) {
	table := featureTable("hint_index")
	if _, err := db.Exec(ctx, fmt.Sprintf(`
		CREATE TABLE %s (
			id int NOT NULL AUTO_INCREMENT,
			a  int,
			b  int,
			PRIMARY KEY (id),
			KEY idx_a (a)
		)`, table)); err != nil {
		gtest.Fatal(err)
	}
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		data := g.List{}
		for i := 1; i <= 20; i++ {
			data = append(data, g.Map{"a": i, "b": i * 10})
		}
		_, err := db.Model(table).Data(data).Insert()
		t.AssertNil(err)

		plan := explainPlan(t, fmt.Sprintf("SELECT /*+ INDEX(%s idx_a) */ * FROM %s WHERE a = 5", table, table))
		t.Assert(gstr.Contains(plan, "TABLE RANGE SCAN"), true)
		t.Assert(gstr.Contains(plan, "idx_a"), true)

		one, err := db.GetOne(ctx, fmt.Sprintf(
			"SELECT /*+ INDEX(%s idx_a) */ b FROM %s WHERE a = 5", table, table,
		))
		t.AssertNil(err)
		t.Assert(one["b"].Int(), 50)
	})
}

// Test_OceanBase_Hint_FullForcesFullScan verifies /*+ FULL(...) */ suppresses the index.
func Test_OceanBase_Hint_FullForcesFullScan(t *testing.T) {
	table := featureTable("hint_full")
	if _, err := db.Exec(ctx, fmt.Sprintf(`
		CREATE TABLE %s (
			id int NOT NULL AUTO_INCREMENT,
			a  int,
			b  int,
			PRIMARY KEY (id),
			KEY idx_a (a)
		)`, table)); err != nil {
		gtest.Fatal(err)
	}
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		data := g.List{}
		for i := 1; i <= 20; i++ {
			data = append(data, g.Map{"a": i, "b": i * 10})
		}
		_, err := db.Model(table).Data(data).Insert()
		t.AssertNil(err)

		hinted := explainPlan(t, fmt.Sprintf("SELECT /*+ FULL(%s) */ * FROM %s WHERE a = 5", table, table))
		t.Assert(gstr.Contains(hinted, "TABLE FULL SCAN"), true)
		t.Assert(gstr.Contains(hinted, "TABLE RANGE SCAN"), false)

		one, err := db.GetOne(ctx, fmt.Sprintf(
			"SELECT /*+ FULL(%s) */ b FROM %s WHERE a = 5", table, table,
		))
		t.AssertNil(err)
		t.Assert(one["b"].Int(), 50)
	})
}

// Test_OceanBase_Explain_ReturnsQueryPlanColumn verifies EXPLAIN returns OceanBase's
// single "Query Plan" column instead of MySQL's id/select_type/table row.
func Test_OceanBase_Explain_ReturnsQueryPlanColumn(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		all, err := db.GetAll(ctx, fmt.Sprintf("EXPLAIN SELECT * FROM %s WHERE id = 1", table))
		t.AssertNil(err)
		t.AssertGT(len(all), 1)
		t.Assert(len(all[0].Map()), 1)
		_, hasQueryPlan := all[0].Map()["Query Plan"]
		t.Assert(hasQueryPlan, true)
		_, hasSelectType := all[0].Map()["select_type"]
		t.Assert(hasSelectType, false)

		plan := explainPlan(t, fmt.Sprintf("SELECT * FROM %s WHERE id = 1", table))
		t.Assert(gstr.Contains(plan, "OPERATOR"), true)
		t.Assert(gstr.Contains(plan, "EST.ROWS"), true)
		t.Assert(gstr.Contains(plan, "Outputs & filters"), true)
		t.Assert(gstr.Contains(plan, table), true)
	})
}

// Test_OceanBase_Explain_ShowsPartitionPruning verifies the plan names the surviving partition.
func Test_OceanBase_Explain_ShowsPartitionPruning(t *testing.T) {
	table := featureTable("explain_part")
	if _, err := db.Exec(ctx, fmt.Sprintf(`
		CREATE TABLE %s (
			id   int NOT NULL,
			name varchar(50),
			PRIMARY KEY (id)
		) PARTITION BY HASH(id) PARTITIONS 4`, table)); err != nil {
		gtest.Fatal(err)
	}
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		data := g.List{}
		for i := 1; i <= 8; i++ {
			data = append(data, g.Map{"id": i, "name": fmt.Sprintf("n%d", i)})
		}
		_, err := db.Model(table).Data(data).Insert()
		t.AssertNil(err)

		pruned := explainPlan(t, fmt.Sprintf("SELECT * FROM %s WHERE id = 1", table))
		t.Assert(gstr.Contains(pruned, "partitions(p1)"), true)

		full := explainPlan(t, fmt.Sprintf("SELECT * FROM %s", table))
		t.Assert(gstr.Contains(full, "partitions(p[0-3])"), true)
	})
}

// =============================================================================
// Native partitioning
// =============================================================================

// Test_OceanBase_Partition_HashFourPartitions verifies HASH partitioning is created,
// reported by information_schema and prunable through raw partition selection.
func Test_OceanBase_Partition_HashFourPartitions(t *testing.T) {
	table := featureTable("part_hash")
	if _, err := db.Exec(ctx, fmt.Sprintf(`
		CREATE TABLE %s (
			id   int NOT NULL,
			name varchar(50),
			PRIMARY KEY (id)
		) PARTITION BY HASH(id) PARTITIONS 4`, table)); err != nil {
		gtest.Fatal(err)
	}
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		data := g.List{}
		for i := 1; i <= 8; i++ {
			data = append(data, g.Map{"id": i, "name": fmt.Sprintf("n%d", i)})
		}
		_, err := db.Model(table).Data(data).Insert()
		t.AssertNil(err)

		total, err := db.Model(table).Count()
		t.AssertNil(err)
		t.Assert(total, 8)

		parts, err := db.GetAll(ctx,
			"SELECT partition_name, partition_method FROM information_schema.partitions "+
				"WHERE table_schema = ? AND table_name = ? ORDER BY partition_name",
			TestSchema1, table,
		)
		t.AssertNil(err)
		t.Assert(len(parts), 4)
		for i, part := range parts {
			t.Assert(part["partition_name"].String(), fmt.Sprintf("p%d", i))
			t.Assert(part["partition_method"].String(), "HASH")
		}

		// Each hash partition holds a strict subset and together they cover the table.
		sum := 0
		for i := 0; i < 4; i++ {
			v, err := db.GetValue(ctx, fmt.Sprintf(
				"SELECT COUNT(*) FROM %s PARTITION(p%d)", table, i,
			))
			t.AssertNil(err)
			t.Assert(v.Int(), 2)
			sum += v.Int()
		}
		t.Assert(sum, 8)

		one, err := db.Model(table).Where("id", 5).One()
		t.AssertNil(err)
		t.Assert(one["name"].String(), "n5")
	})
}

// Test_OceanBase_Partition_RangeByValue verifies RANGE partitioning routes rows by value.
func Test_OceanBase_Partition_RangeByValue(t *testing.T) {
	table := featureTable("part_range")
	if _, err := db.Exec(ctx, fmt.Sprintf(`
		CREATE TABLE %s (
			id int NOT NULL,
			y  int NOT NULL,
			PRIMARY KEY (id, y)
		) PARTITION BY RANGE (y) (
			PARTITION p2020 VALUES LESS THAN (2021),
			PARTITION p2021 VALUES LESS THAN (2022),
			PARTITION pmax  VALUES LESS THAN MAXVALUE
		)`, table)); err != nil {
		gtest.Fatal(err)
	}
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(table).Data(g.List{
			{"id": 1, "y": 2019},
			{"id": 2, "y": 2020},
			{"id": 3, "y": 2021},
			{"id": 4, "y": 2030},
		}).Insert()
		t.AssertNil(err)

		parts, err := db.GetAll(ctx,
			"SELECT partition_name, partition_method, partition_description "+
				"FROM information_schema.partitions "+
				"WHERE table_schema = ? AND table_name = ? ORDER BY partition_ordinal_position",
			TestSchema1, table,
		)
		t.AssertNil(err)
		t.Assert(len(parts), 3)
		t.Assert(parts[0]["partition_name"].String(), "p2020")
		t.Assert(parts[0]["partition_method"].String(), "RANGE")
		t.Assert(parts[2]["partition_name"].String(), "pmax")

		expected := map[string]int{"p2020": 2, "p2021": 1, "pmax": 1}
		for name, want := range expected {
			v, err := db.GetValue(ctx, fmt.Sprintf(
				"SELECT COUNT(*) FROM %s PARTITION(%s)", table, name,
			))
			t.AssertNil(err)
			t.Assert(v.Int(), want)
		}

		ids, err := db.GetAll(ctx, fmt.Sprintf(
			"SELECT id FROM %s PARTITION(p2020) ORDER BY id", table,
		))
		t.AssertNil(err)
		t.Assert(len(ids), 2)
		t.Assert(ids[0]["id"].Int(), 1)
		t.Assert(ids[1]["id"].Int(), 2)
	})
}

// =============================================================================
// Sequences (OceanBase supports them in MySQL mode, MySQL does not)
// =============================================================================

// Test_OceanBase_Sequence_NextvalCurrvalAndStep verifies both spellings of sequence
// access and that START WITH / INCREMENT BY drive the generated values.
func Test_OceanBase_Sequence_NextvalCurrvalAndStep(t *testing.T) {
	basic := featureTable("seq_basic")
	if _, err := db.Exec(ctx, fmt.Sprintf(
		"CREATE SEQUENCE %s START WITH 1 INCREMENT BY 1", basic,
	)); err != nil {
		gtest.Fatal(err)
	}
	defer db.Exec(ctx, fmt.Sprintf("DROP SEQUENCE %s", basic))

	stepped := featureTable("seq_step")
	if _, err := db.Exec(ctx, fmt.Sprintf(
		"CREATE SEQUENCE %s START WITH 100 INCREMENT BY 5 MAXVALUE 1000 NOCYCLE", stepped,
	)); err != nil {
		gtest.Fatal(err)
	}
	defer db.Exec(ctx, fmt.Sprintf("DROP SEQUENCE %s", stepped))

	gtest.C(t, func(t *gtest.T) {
		v, err := db.GetValue(ctx, fmt.Sprintf("SELECT %s.nextval", basic))
		t.AssertNil(err)
		t.Assert(v.Int(), 1)

		v, err = db.GetValue(ctx, fmt.Sprintf("SELECT NEXTVAL(%s)", basic))
		t.AssertNil(err)
		t.Assert(v.Int(), 2)

		v, err = db.GetValue(ctx, fmt.Sprintf("SELECT %s.currval", basic))
		t.AssertNil(err)
		t.Assert(v.Int(), 2)

		v, err = db.GetValue(ctx, fmt.Sprintf("SELECT CURRVAL(%s)", basic))
		t.AssertNil(err)
		t.Assert(v.Int(), 2)

		for _, want := range []int{100, 105, 110, 115} {
			v, err = db.GetValue(ctx, fmt.Sprintf("SELECT %s.nextval", stepped))
			t.AssertNil(err)
			t.Assert(v.Int(), want)
		}
	})
}

// Test_OceanBase_Sequence_InsertUsingNextval verifies a sequence can drive a primary key.
func Test_OceanBase_Sequence_InsertUsingNextval(t *testing.T) {
	name := featureTable("seq_insert")
	if _, err := db.Exec(ctx, fmt.Sprintf(
		"CREATE SEQUENCE %s START WITH 50 INCREMENT BY 10", name,
	)); err != nil {
		gtest.Fatal(err)
	}
	defer db.Exec(ctx, fmt.Sprintf("DROP SEQUENCE %s", name))

	table := featureTable("seq_table")
	if _, err := db.Exec(ctx, fmt.Sprintf(
		"CREATE TABLE %s (id bigint NOT NULL, v varchar(20), PRIMARY KEY (id))", table,
	)); err != nil {
		gtest.Fatal(err)
	}
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		for _, v := range []string{"a", "b", "c"} {
			_, err := db.Exec(ctx, fmt.Sprintf(
				"INSERT INTO %s (id, v) VALUES (%s.nextval, '%s')", table, name, v,
			))
			t.AssertNil(err)
		}

		all, err := db.Model(table).OrderAsc("id").All()
		t.AssertNil(err)
		t.Assert(len(all), 3)
		t.Assert(all[0]["id"].Int(), 50)
		t.Assert(all[0]["v"].String(), "a")
		t.Assert(all[1]["id"].Int(), 60)
		t.Assert(all[2]["id"].Int(), 70)
	})
}

// =============================================================================
// JSON
// =============================================================================

// Test_OceanBase_JSON_ExtractAndUnquote verifies JSON path extraction operators.
func Test_OceanBase_JSON_ExtractAndUnquote(t *testing.T) {
	table := featureTable("json_extract")
	if _, err := db.Exec(ctx, fmt.Sprintf(
		"CREATE TABLE %s (id int NOT NULL, doc json, PRIMARY KEY (id))", table,
	)); err != nil {
		gtest.Fatal(err)
	}
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(table).Data(g.Map{
			"id":  1,
			"doc": `{"name":"alice","age":30,"tags":["go","oceanbase"]}`,
		}).Insert()
		t.AssertNil(err)

		one, err := db.Model(table).Fields(
			`JSON_EXTRACT(doc, '$.name') AS quoted`,
			`doc->>'$.name' AS unquoted`,
			`doc->'$.age' AS age`,
			`JSON_UNQUOTE(JSON_EXTRACT(doc, '$.tags[1]')) AS tag1`,
		).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["quoted"].String(), `"alice"`)
		t.Assert(one["unquoted"].String(), "alice")
		t.Assert(one["age"].Int(), 30)
		t.Assert(one["tag1"].String(), "oceanbase")
	})
}

// Test_OceanBase_JSON_ContainsTypeKeysLength verifies the JSON inspection functions.
func Test_OceanBase_JSON_ContainsTypeKeysLength(t *testing.T) {
	table := featureTable("json_inspect")
	if _, err := db.Exec(ctx, fmt.Sprintf(
		"CREATE TABLE %s (id int NOT NULL, doc json, PRIMARY KEY (id))", table,
	)); err != nil {
		gtest.Fatal(err)
	}
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(table).Data(g.List{
			{"id": 1, "doc": `{"role":"admin","tags":["go","ob"]}`},
			{"id": 2, "doc": `{"role":"user","tags":["go"]}`},
		}).Insert()
		t.AssertNil(err)

		one, err := db.Model(table).Fields(
			`JSON_LENGTH(doc->'$.tags') AS tag_count`,
			`JSON_TYPE(doc->'$.tags') AS tag_type`,
			`JSON_KEYS(doc) AS keys_json`,
			`JSON_CONTAINS(doc, '"go"', '$.tags') AS has_go`,
		).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["tag_count"].Int(), 2)
		t.Assert(one["tag_type"].String(), "ARRAY")
		t.Assert(one["keys_json"].Strings(), g.SliceStr{"role", "tags"})
		t.Assert(one["has_go"].Int(), 1)

		all, err := db.Model(table).
			Where(`JSON_UNQUOTE(JSON_EXTRACT(doc, '$.role')) = ?`, "admin").
			OrderAsc("id").All()
		t.AssertNil(err)
		t.Assert(len(all), 1)
		t.Assert(all[0]["id"].Int(), 1)
	})
}

// Test_OceanBase_JSON_SetAndAggregates verifies JSON mutation and JSON aggregate functions.
func Test_OceanBase_JSON_SetAndAggregates(t *testing.T) {
	table := featureTable("json_agg")
	if _, err := db.Exec(ctx, fmt.Sprintf(
		"CREATE TABLE %s (id int NOT NULL, name varchar(20), doc json, PRIMARY KEY (id))", table,
	)); err != nil {
		gtest.Fatal(err)
	}
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(table).Data(g.List{
			{"id": 1, "name": "alice", "doc": `{"age":30}`},
			{"id": 2, "name": "bob", "doc": `{"age":25}`},
		}).Insert()
		t.AssertNil(err)

		_, err = db.Exec(ctx, fmt.Sprintf(
			`UPDATE %s SET doc = JSON_SET(doc, '$.role', 'admin') WHERE id = 1`, table,
		))
		t.AssertNil(err)

		one, err := db.Model(table).Where("id", 1).One()
		t.AssertNil(err)
		m := one["doc"].Map()
		t.Assert(m["age"], 30)
		t.Assert(m["role"], "admin")

		one, err = db.Model(table).Fields(
			`JSON_ARRAYAGG(id) AS ids`,
			`JSON_OBJECTAGG(name, id) AS by_name`,
		).One()
		t.AssertNil(err)
		t.Assert(one["ids"].Ints(), g.SliceInt{1, 2})
		byName := one["by_name"].Map()
		t.Assert(byName["alice"], 1)
		t.Assert(byName["bob"], 2)
	})
}

// =============================================================================
// CTE, window functions, GROUP_CONCAT
// =============================================================================

// Test_OceanBase_CTE_Basic verifies a non-recursive WITH clause.
func Test_OceanBase_CTE_Basic(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		all, err := db.GetAll(ctx, fmt.Sprintf(`
			WITH top5 AS (
				SELECT id, nickname FROM %s ORDER BY id ASC LIMIT 5
			)
			SELECT * FROM top5 WHERE id > 2 ORDER BY id
		`, table))
		t.AssertNil(err)
		t.Assert(len(all), 3)
		t.Assert(all[0]["id"].Int(), 3)
		t.Assert(all[0]["nickname"].String(), "name_3")
		t.Assert(all[2]["id"].Int(), 5)
	})
}

// Test_OceanBase_CTE_Recursive verifies WITH RECURSIVE over hierarchical data.
func Test_OceanBase_CTE_Recursive(t *testing.T) {
	table := featureTable("cte_tree")
	if _, err := db.Exec(ctx, fmt.Sprintf(`
		CREATE TABLE %s (
			id        int NOT NULL,
			parent_id int NULL,
			name      varchar(50),
			PRIMARY KEY (id)
		)`, table)); err != nil {
		gtest.Fatal(err)
	}
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(table).Data(g.List{
			{"id": 1, "parent_id": nil, "name": "root"},
			{"id": 2, "parent_id": 1, "name": "child_a"},
			{"id": 3, "parent_id": 1, "name": "child_b"},
			{"id": 4, "parent_id": 2, "name": "grandchild"},
		}).Insert()
		t.AssertNil(err)

		all, err := db.GetAll(ctx, fmt.Sprintf(`
			WITH RECURSIVE tree (id, parent_id, name, depth) AS (
				SELECT id, parent_id, name, 0 FROM %s WHERE id = 1
				UNION ALL
				SELECT t.id, t.parent_id, t.name, tree.depth + 1
				FROM %s t JOIN tree ON t.parent_id = tree.id
			)
			SELECT * FROM tree ORDER BY depth, id
		`, table, table))
		t.AssertNil(err)
		t.Assert(len(all), 4)
		t.Assert(all[0]["name"].String(), "root")
		t.Assert(all[0]["depth"].Int(), 0)
		t.Assert(all[3]["name"].String(), "grandchild")
		t.Assert(all[3]["depth"].Int(), 2)
	})
}

// Test_OceanBase_Window_Functions verifies ROW_NUMBER/RANK/LAG/LEAD/SUM OVER.
func Test_OceanBase_Window_Functions(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		all, err := db.GetAll(ctx, fmt.Sprintf(`
			SELECT id,
				ROW_NUMBER() OVER (ORDER BY id DESC) AS rn,
				RANK()       OVER (ORDER BY id ASC)  AS rk,
				LAG(id)      OVER (ORDER BY id ASC)  AS prev_id,
				LEAD(id)     OVER (ORDER BY id ASC)  AS next_id,
				SUM(id)      OVER (ORDER BY id ASC)  AS running_total
			FROM %s ORDER BY id
		`, table))
		t.AssertNil(err)
		t.Assert(len(all), TableSize)

		t.Assert(all[0]["id"].Int(), 1)
		t.Assert(all[0]["rn"].Int(), TableSize)
		t.Assert(all[0]["rk"].Int(), 1)
		t.Assert(all[0]["prev_id"].IsNil(), true)
		t.Assert(all[0]["next_id"].Int(), 2)
		t.Assert(all[0]["running_total"].Int(), 1)

		t.Assert(all[4]["running_total"].Int(), 15)
		t.Assert(all[9]["rn"].Int(), 1)
		t.Assert(all[9]["prev_id"].Int(), 9)
		t.Assert(all[9]["next_id"].IsNil(), true)
		t.Assert(all[9]["running_total"].Int(), 55)
	})
}

// Test_OceanBase_GroupConcat_OrderAndSeparator verifies GROUP_CONCAT with ORDER BY and SEPARATOR.
func Test_OceanBase_GroupConcat_OrderAndSeparator(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		one, err := db.Model(table).
			Fields(`GROUP_CONCAT(nickname ORDER BY id ASC SEPARATOR '|') AS names`).
			Where("id <= ?", 3).One()
		t.AssertNil(err)
		t.Assert(one["names"].String(), "name_1|name_2|name_3")

		one, err = db.Model(table).
			Fields(`GROUP_CONCAT(DISTINCT LEFT(passport, 4) ORDER BY 1 ASC) AS prefixes`).
			One()
		t.AssertNil(err)
		t.Assert(one["prefixes"].String(), "user")
	})
}

// =============================================================================
// Upsert family
// =============================================================================

// Test_OceanBase_Save_OnDuplicateKeyUpdate verifies Save() emits ON DUPLICATE KEY UPDATE.
func Test_OceanBase_Save_OnDuplicateKeyUpdate(t *testing.T) {
	table := createTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(table).Data(g.Map{
			"id": 1, "passport": "p1", "password": "pw1", "nickname": "n1",
		}).Insert()
		t.AssertNil(err)

		result, err := db.Model(table).Data(g.Map{
			"id": 1, "passport": "p2", "password": "pw2", "nickname": "n2",
		}).Save()
		t.AssertNil(err)
		affected, err := result.RowsAffected()
		t.AssertNil(err)
		t.Assert(affected, 2)

		one, err := db.Model(table).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["passport"].String(), "p2")
		t.Assert(one["password"].String(), "pw2")
		t.Assert(one["nickname"].String(), "n2")

		count, err := db.Model(table).Count()
		t.AssertNil(err)
		t.Assert(count, 1)

		result, err = db.Model(table).Data(g.Map{
			"id": 2, "passport": "p3", "password": "pw3", "nickname": "n3",
		}).Save()
		t.AssertNil(err)
		affected, err = result.RowsAffected()
		t.AssertNil(err)
		t.Assert(affected, 1)
		count, err = db.Model(table).Count()
		t.AssertNil(err)
		t.Assert(count, 2)
	})
}

// Test_OceanBase_OnDuplicate_SelectedColumnsOnly verifies OnDuplicate limits the update list.
func Test_OceanBase_OnDuplicate_SelectedColumnsOnly(t *testing.T) {
	table := createTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(table).Data(g.Map{
			"id": 1, "passport": "keep", "password": "keep_pw", "nickname": "old",
		}).Insert()
		t.AssertNil(err)

		_, err = db.Model(table).Data(g.Map{
			"id": 1, "passport": "changed", "password": "changed_pw", "nickname": "new",
		}).OnDuplicate("nickname").Save()
		t.AssertNil(err)

		one, err := db.Model(table).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["nickname"].String(), "new")
		t.Assert(one["passport"].String(), "keep")
		t.Assert(one["password"].String(), "keep_pw")
	})
}

// Test_OceanBase_Replace_ResetsUnsetColumns verifies REPLACE deletes and re-inserts the row.
func Test_OceanBase_Replace_ResetsUnsetColumns(t *testing.T) {
	table := createTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(table).Data(g.Map{
			"id": 1, "passport": "p1", "password": "pw1", "nickname": "n1",
		}).Insert()
		t.AssertNil(err)

		result, err := db.Model(table).Data(g.Map{
			"id": 1, "passport": "p2", "nickname": "n2",
		}).Replace()
		t.AssertNil(err)
		affected, err := result.RowsAffected()
		t.AssertNil(err)
		t.Assert(affected, 2)

		one, err := db.Model(table).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["passport"].String(), "p2")
		t.Assert(one["nickname"].String(), "n2")
		// Columns omitted from REPLACE are reset because the old row was deleted.
		t.Assert(one["password"].String(), "")

		count, err := db.Model(table).Count()
		t.AssertNil(err)
		t.Assert(count, 1)
	})
}

// Test_OceanBase_InsertIgnore_SkipsDuplicateKey verifies INSERT IGNORE leaves the row untouched.
func Test_OceanBase_InsertIgnore_SkipsDuplicateKey(t *testing.T) {
	table := createTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(table).Data(g.Map{
			"id": 1, "passport": "original", "nickname": "n1",
		}).Insert()
		t.AssertNil(err)

		result, err := db.Model(table).Data(g.Map{
			"id": 1, "passport": "ignored", "nickname": "n2",
		}).InsertIgnore()
		t.AssertNil(err)
		affected, err := result.RowsAffected()
		t.AssertNil(err)
		t.Assert(affected, 0)

		one, err := db.Model(table).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["passport"].String(), "original")
		t.Assert(one["nickname"].String(), "n1")

		result, err = db.Model(table).Data(g.Map{
			"id": 2, "passport": "fresh", "nickname": "n3",
		}).InsertIgnore()
		t.AssertNil(err)
		affected, err = result.RowsAffected()
		t.AssertNil(err)
		t.Assert(affected, 1)
	})
}

// =============================================================================
// Tenant and system views
// =============================================================================

// Test_OceanBase_Version_IdentifiesOceanBase verifies the server identifies itself as OceanBase.
func Test_OceanBase_Version_IdentifiesOceanBase(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		version, err := db.GetValue(ctx, "SELECT VERSION()")
		t.AssertNil(err)
		t.Assert(gstr.Contains(version.String(), "OceanBase"), true)

		comment, err := db.GetValue(ctx, "SELECT @@version_comment")
		t.AssertNil(err)
		t.Assert(gstr.Contains(comment.String(), "OceanBase"), true)
	})
}

// Test_OceanBase_Tenant_QualifiedUserAndSystemView verifies the tenant-qualified login
// root@test lands in the `test` tenant and that the oceanbase system view is readable
// from a user tenant.
func Test_OceanBase_Tenant_QualifiedUserAndSystemView(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		all, err := db.GetAll(ctx, "SHOW TENANT")
		t.AssertNil(err)
		t.Assert(len(all), 1)
		t.Assert(all[0]["Current_tenant_name"].String(), "test")

		one, err := db.GetOne(ctx, "SELECT USER() AS u, DATABASE() AS d")
		t.AssertNil(err)
		t.Assert(gstr.HasPrefix(one["u"].String(), TestDbUser+"@"), true)
		t.Assert(one["d"].String(), TestSchema1)

		one, err = db.GetOne(ctx,
			"SELECT TENANT_NAME, TENANT_TYPE, COMPATIBILITY_MODE, STATUS "+
				"FROM oceanbase.DBA_OB_TENANTS WHERE TENANT_NAME = ?", "test",
		)
		t.AssertNil(err)
		t.Assert(one["TENANT_NAME"].String(), "test")
		t.Assert(one["TENANT_TYPE"].String(), "USER")
		t.Assert(one["COMPATIBILITY_MODE"].String(), "MYSQL")
		t.Assert(one["STATUS"].String(), "NORMAL")
	})
}

// =============================================================================
// Isolation levels
// =============================================================================

// Test_OceanBase_Isolation_DefaultIsReadCommitted pins that OceanBase defaults to
// READ-COMMITTED, where MySQL defaults to REPEATABLE-READ.
func Test_OceanBase_Isolation_DefaultIsReadCommitted(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		v, err := db.GetValue(ctx, "SELECT @@tx_isolation")
		t.AssertNil(err)
		t.Assert(v.String(), "READ-COMMITTED")

		tx, err := db.Begin(ctx)
		t.AssertNil(err)
		defer tx.Rollback()
		v, err = tx.GetValue("SELECT @@tx_isolation")
		t.AssertNil(err)
		t.Assert(v.String(), "READ-COMMITTED")
	})
}

// Test_OceanBase_Isolation_ReadCommittedSeesConcurrentCommit verifies a READ COMMITTED
// transaction picks up rows committed by another connection.
func Test_OceanBase_Isolation_ReadCommittedSeesConcurrentCommit(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		tx, err := db.BeginWithOptions(ctx, gdb.TxOptions{Isolation: sql.LevelReadCommitted})
		t.AssertNil(err)
		defer tx.Rollback()

		before, err := tx.Model(table).Count()
		t.AssertNil(err)
		t.Assert(before, TableSize)

		_, err = db.Model(table).Data(g.Map{"passport": "outsider", "nickname": "outsider"}).Insert()
		t.AssertNil(err)

		after, err := tx.Model(table).Count()
		t.AssertNil(err)
		t.Assert(after, TableSize+1)

		one, err := tx.Model(table).Where("passport", "outsider").One()
		t.AssertNil(err)
		t.Assert(one["nickname"].String(), "outsider")
	})
}

// Test_OceanBase_Isolation_RepeatableReadHoldsSnapshot verifies a REPEATABLE READ
// transaction keeps reading its own snapshot.
func Test_OceanBase_Isolation_RepeatableReadHoldsSnapshot(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		tx, err := db.BeginWithOptions(ctx, gdb.TxOptions{Isolation: sql.LevelRepeatableRead})
		t.AssertNil(err)
		defer tx.Rollback()

		before, err := tx.Model(table).Count()
		t.AssertNil(err)
		t.Assert(before, TableSize)

		_, err = db.Model(table).Data(g.Map{"passport": "outsider", "nickname": "outsider"}).Insert()
		t.AssertNil(err)

		after, err := tx.Model(table).Count()
		t.AssertNil(err)
		t.Assert(after, TableSize)

		one, err := tx.Model(table).Where("passport", "outsider").One()
		t.AssertNil(err)
		t.Assert(one.IsEmpty(), true)

		// The row is visible to a connection outside the snapshot.
		outside, err := db.Model(table).Count()
		t.AssertNil(err)
		t.Assert(outside, TableSize+1)
	})
}

// Test_OceanBase_Isolation_SerializableIsSnapshot pins that OceanBase implements
// SERIALIZABLE as snapshot isolation: a concurrent insert is not blocked and the
// transaction keeps reading its own snapshot.
func Test_OceanBase_Isolation_SerializableIsSnapshot(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		tx, err := db.BeginWithOptions(ctx, gdb.TxOptions{Isolation: sql.LevelSerializable})
		t.AssertNil(err)
		defer tx.Rollback()

		before, err := tx.Model(table).Count()
		t.AssertNil(err)
		t.Assert(before, TableSize)

		done := make(chan error, 1)
		go func() {
			_, e := db.Model(table).Data(g.Map{"passport": "outsider", "nickname": "outsider"}).Insert()
			done <- e
		}()

		select {
		case e := <-done:
			t.AssertNil(e)
		case <-time.After(2 * time.Second):
			t.Fatal("concurrent insert was blocked by the SERIALIZABLE transaction")
		}

		after, err := tx.Model(table).Count()
		t.AssertNil(err)
		t.Assert(after, TableSize)

		t.AssertNil(tx.Commit())

		final, err := db.Model(table).Count()
		t.AssertNil(err)
		t.Assert(final, TableSize+1)
	})
}

// =============================================================================
// Row locking
// =============================================================================

// Test_OceanBase_Lock_ForUpdateBlocksWriter verifies FOR UPDATE holds off a second
// writer until the holding transaction commits.
func Test_OceanBase_Lock_ForUpdateBlocksWriter(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		tx, err := db.Begin(ctx)
		t.AssertNil(err)

		one, err := tx.Model(table).LockUpdate().Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["nickname"].String(), "name_1")

		done := make(chan error, 1)
		go func() {
			_, e := db.Model(table).Data(g.Map{"nickname": "by_writer"}).Where("id", 1).Update()
			done <- e
		}()

		select {
		case e := <-done:
			tx.Rollback()
			t.Fatalf("second writer was not blocked by FOR UPDATE: %v", e)
		case <-time.After(800 * time.Millisecond):
		}

		t.AssertNil(tx.Commit())
		t.AssertNil(<-done)

		one, err = db.Model(table).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["nickname"].String(), "by_writer")
	})
}

// Test_OceanBase_Lock_ShareModeBlocksWriter verifies OceanBase honours
// LOCK IN SHARE MODE rather than treating it as a no-op: an unlocked read still
// goes through while a writer waits for the holding transaction to commit.
func Test_OceanBase_Lock_ShareModeBlocksWriter(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		tx, err := db.Begin(ctx)
		t.AssertNil(err)

		one, err := tx.Model(table).LockShared().Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["nickname"].String(), "name_1")

		// An unlocked read is served from the snapshot and is never blocked.
		plain, err := db.Model(table).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(plain["nickname"].String(), "name_1")

		done := make(chan error, 1)
		go func() {
			_, e := db.Model(table).Data(g.Map{"nickname": "by_writer"}).Where("id", 1).Update()
			done <- e
		}()

		select {
		case e := <-done:
			tx.Rollback()
			t.Fatalf("writer was not blocked by LOCK IN SHARE MODE: %v", e)
		case <-time.After(800 * time.Millisecond):
		}

		t.AssertNil(tx.Commit())
		t.AssertNil(<-done)

		one, err = db.Model(table).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["nickname"].String(), "by_writer")
	})
}

// Test_OceanBase_Lock_ShareModeIsExclusive pins that OceanBase's LOCK IN SHARE MODE is
// not shareable: a second session asking for the same lock blocks, where InnoDB would
// grant both readers the lock concurrently.
func Test_OceanBase_Lock_ShareModeIsExclusive(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		tx, err := db.Begin(ctx)
		t.AssertNil(err)

		one, err := tx.Model(table).LockShared().Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["nickname"].String(), "name_1")

		type readResult struct {
			row gdb.Record
			err error
		}
		done := make(chan readResult, 1)
		go func() {
			row, e := db.Model(table).LockShared().Where("id", 1).One()
			done <- readResult{row: row, err: e}
		}()

		select {
		case r := <-done:
			tx.Rollback()
			t.Fatalf("second LOCK IN SHARE MODE reader was not blocked: %v %v", r.row, r.err)
		case <-time.After(800 * time.Millisecond):
		}

		t.AssertNil(tx.Commit())
		r := <-done
		t.AssertNil(r.err)
		t.Assert(r.row["nickname"].String(), "name_1")

		// A row locked in share mode by another session is still readable without a lock.
		plain, err := db.Model(table).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(plain["nickname"].String(), "name_1")
	})
}

// =============================================================================
// Weak consistency reads
// =============================================================================

// Test_OceanBase_ReadConsistency_WeakSessionAndHint verifies ob_read_consistency can be
// switched to WEAK per session and through the READ_CONSISTENCY hint, and that the weak
// read path converges on the committed data.
func Test_OceanBase_ReadConsistency_WeakSessionAndHint(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		strong, err := db.Model(table).Count()
		t.AssertNil(err)
		t.Assert(strong, TableSize)

		tx, err := db.Begin(ctx)
		t.AssertNil(err)
		defer tx.Rollback()

		original, err := tx.GetValue("SELECT @@ob_read_consistency")
		t.AssertNil(err)
		t.Assert(original.String(), "STRONG")
		defer tx.Exec(fmt.Sprintf("SET SESSION ob_read_consistency = %s", original.String()))

		_, err = tx.Exec("SET SESSION ob_read_consistency = WEAK")
		t.AssertNil(err)
		v, err := tx.GetValue("SELECT @@ob_read_consistency")
		t.AssertNil(err)
		t.Assert(v.String(), "WEAK")

		// A weak read trails the leader, so poll until the snapshot catches up.
		weak := 0
		for i := 0; i < 20; i++ {
			weak, err = tx.Model(table).Count()
			t.AssertNil(err)
			if weak == TableSize {
				break
			}
			time.Sleep(200 * time.Millisecond)
		}
		t.Assert(weak, TableSize)

		_, err = tx.Exec(fmt.Sprintf("SET SESSION ob_read_consistency = %s", original.String()))
		t.AssertNil(err)
		v, err = tx.GetValue("SELECT @@ob_read_consistency")
		t.AssertNil(err)
		t.Assert(v.String(), "STRONG")

		hinted := 0
		for i := 0; i < 20; i++ {
			one, err := db.GetOne(ctx, fmt.Sprintf(
				"SELECT /*+ READ_CONSISTENCY(WEAK) */ COUNT(*) AS n FROM %s", table,
			))
			t.AssertNil(err)
			hinted = one["n"].Int()
			if hinted == TableSize {
				break
			}
			time.Sleep(200 * time.Millisecond)
		}
		t.Assert(hinted, TableSize)
	})
}
