// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package clickhouse_test

import (
	"context"
	"fmt"
	"math"
	"strings"
	"testing"
	"time"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/test/gtest"
)

// chDCreateReplacingTable creates a table like the default test table on the ReplacingMergeTree engine,
// which keeps only the last inserted record of the same id once its parts are merged.
func chDCreateReplacingTable() string {
	name := fmt.Sprintf(`%s_%d`, TableName, gtime.TimestampNano())
	if _, err := db.Exec(ctx, fmt.Sprintf(`
		CREATE TABLE %s (
		   id UInt64 NOT NULL,
		   passport String,
		   password String,
		   nickname String,
		   create_time DateTime
		) ENGINE = ReplacingMergeTree()
		ORDER BY id`, name,
	)); err != nil {
		gtest.Fatal(err)
	}
	return name
}

// chDOptimizeFinal merges all parts of `table`, which applies the deduplication of a ReplacingMergeTree table.
func chDOptimizeFinal(table string) error {
	_, err := db.Exec(ctx, fmt.Sprintf("OPTIMIZE TABLE %s FINAL", table))
	return err
}

// Test_Model_Batch_Insert tests batch insert with different batch sizes
func Test_Model_Batch_Insert(t *testing.T) {
	table := createTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		// Prepare data for batch insert
		data := g.Slice{}
		for i := 1; i <= 10; i++ {
			data = append(data, g.Map{
				"id":       i,
				"passport": fmt.Sprintf("batch_user_%d", i),
				"password": fmt.Sprintf("batch_pass_%d", i),
				"nickname": fmt.Sprintf("batch_name_%d", i),
			})
		}

		// Batch insert with batch size 3
		// The inserted rows are verified by reading them back.
		_, err := db.Model(table).Batch(3).Data(data).Insert()
		t.AssertNil(err)

		// Verify all records were inserted
		count, err := db.Model(table).Count()
		t.AssertNil(err)
		t.Assert(count, 10)

		// Verify specific records
		one, err := db.Model(table).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["passport"], "batch_user_1")

		one, err = db.Model(table).Where("id", 10).One()
		t.AssertNil(err)
		t.Assert(one["passport"], "batch_user_10")
	})
}

// Test_Model_Batch_Replace tests batch replace operation
func Test_Model_Batch_Replace(t *testing.T) {
	// Note: Replace through a model inserts the records, so the replacement is done by a
	// ReplacingMergeTree table once its parts are merged.
	table := chDCreateReplacingTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		// Initial insert
		data := g.Slice{}
		for i := 1; i <= 5; i++ {
			data = append(data, g.Map{
				"id":       i,
				"passport": fmt.Sprintf("original_%d", i),
			})
		}
		_, err := db.Model(table).Data(data).Insert()
		t.AssertNil(err)

		// Batch replace with overlapping ids
		replaceData := g.Slice{}
		for i := 3; i <= 8; i++ {
			replaceData = append(replaceData, g.Map{
				"id":       i,
				"passport": fmt.Sprintf("replaced_%d", i),
				"nickname": fmt.Sprintf("new_name_%d", i),
			})
		}
		_, err = db.Model(table).Batch(2).Data(replaceData).Replace()
		t.AssertNil(err)

		// Before the merge both versions of the overlapping ids exist.
		count, err := db.Model(table).Count()
		t.AssertNil(err)
		t.Assert(count, 11)
		t.AssertNil(chDOptimizeFinal(table))

		// Verify replaced records
		one, err := db.Model(table).Where("id", 3).One()
		t.AssertNil(err)
		t.Assert(one["passport"], "replaced_3")
		t.Assert(one["nickname"], "new_name_3")

		// Verify new records
		one, err = db.Model(table).Where("id", 8).One()
		t.AssertNil(err)
		t.Assert(one["passport"], "replaced_8")

		// Verify total count
		count, err = db.Model(table).Count()
		t.AssertNil(err)
		t.Assert(count, 8) // ids 1-8
	})
}

// Test_Model_Batch_Save tests batch save operation
func Test_Model_Batch_Save(t *testing.T) {
	// Note: Save through a model inserts the records, so the update is done by a
	// ReplacingMergeTree table once its parts are merged.
	table := chDCreateReplacingTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		// Initial data
		data := g.Slice{}
		for i := 1; i <= 5; i++ {
			data = append(data, g.Map{
				"id":       i,
				"passport": fmt.Sprintf("save_user_%d", i),
			})
		}
		_, err := db.Model(table).Data(data).Insert()
		t.AssertNil(err)

		// Batch save with overlapping and new ids
		saveData := g.Slice{}
		for i := 3; i <= 8; i++ {
			saveData = append(saveData, g.Map{
				"id":       i,
				"passport": fmt.Sprintf("saved_%d", i),
				"nickname": fmt.Sprintf("save_name_%d", i),
			})
		}
		_, err = db.Model(table).Batch(3).Data(saveData).Save()
		t.AssertNil(err)

		// Before the merge both versions of the overlapping ids exist.
		count, err := db.Model(table).Count()
		t.AssertNil(err)
		t.Assert(count, 11)
		t.AssertNil(chDOptimizeFinal(table))

		// Verify updated records
		one, err := db.Model(table).Where("id", 3).One()
		t.AssertNil(err)
		t.Assert(one["passport"], "saved_3")

		// Verify inserted records
		one, err = db.Model(table).Where("id", 8).One()
		t.AssertNil(err)
		t.Assert(one["passport"], "saved_8")

		// Verify total count
		count, err = db.Model(table).Count()
		t.AssertNil(err)
		t.Assert(count, 8)
	})
}

// Test_Model_Batch_LargeBatch tests batch operation with large dataset
func Test_Model_Batch_LargeBatch(t *testing.T) {
	table := createTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		// Prepare 1000+ records
		data := g.Slice{}
		totalRecords := 1500
		for i := 1; i <= totalRecords; i++ {
			data = append(data, g.Map{
				"id":       i,
				"passport": fmt.Sprintf("large_user_%d", i),
				"nickname": fmt.Sprintf("large_name_%d", i),
			})
		}

		// Batch insert with batch size 100
		_, err := db.Model(table).Batch(100).Data(data).Insert()
		t.AssertNil(err)

		// Verify count
		count, err := db.Model(table).Count()
		t.AssertNil(err)
		t.Assert(count, totalRecords)

		// Verify first and last records
		one, err := db.Model(table).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["passport"], "large_user_1")

		one, err = db.Model(table).Where("id", totalRecords).One()
		t.AssertNil(err)
		t.Assert(one["passport"], fmt.Sprintf("large_user_%d", totalRecords))
	})
}

// Test_Model_Batch_PrepareTimeout tests batch insert on a database configured with a prepare timeout.
func Test_Model_Batch_PrepareTimeout(t *testing.T) {
	table := createTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		node := *db.GetConfig()
		node.PrepareTimeout = time.Second
		timeoutDb, err := gdb.New(node)
		t.AssertNil(err)
		defer timeoutDb.Close(ctx)

		data := g.Slice{}
		for i := 1; i <= 3; i++ {
			data = append(data, g.Map{
				"id":          i,
				"passport":    fmt.Sprintf("user_%d", i),
				"password":    fmt.Sprintf("pass_%d", i),
				"nickname":    fmt.Sprintf("name_%d", i),
				"create_time": gtime.Now(),
			})
		}
		result, err := timeoutDb.Model(table).Data(data).Insert()
		t.AssertNil(err)
		n, err := result.RowsAffected()
		t.AssertNil(err)
		t.Assert(n, 3)

		count, err := db.Model(table).Count()
		t.AssertNil(err)
		t.Assert(count, 3)
	})
}

// Test_Model_Batch_InsertOrder tests that records with raw values or different fields are inserted
// in their order, which decides the record a ReplacingMergeTree table keeps.
func Test_Model_Batch_InsertOrder(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		table := chDCreateReplacingTable()
		defer dropTable(table)

		_, err := db.Model(table).Data(g.Slice{
			g.Map{"id": 1, "nickname": "first"},
			g.Map{"id": 1, "nickname": gdb.Raw("'second'")},
		}).Insert()
		t.AssertNil(err)
		t.AssertNil(chDOptimizeFinal(table))

		value, err := db.Model(table).Where("id", 1).Value("nickname")
		t.AssertNil(err)
		t.Assert(value, "second")
	})
	gtest.C(t, func(t *gtest.T) {
		table := chDCreateReplacingTable()
		defer dropTable(table)

		_, err := db.Model(table).Data(g.Slice{
			g.Map{"id": 1, "passport": "user_1", "nickname": "first"},
			g.Map{"id": 1, "nickname": "second"},
			g.Map{"id": 1, "passport": "user_1", "nickname": "third"},
		}).Insert()
		t.AssertNil(err)
		t.AssertNil(chDOptimizeFinal(table))

		value, err := db.Model(table).Where("id", 1).Value("nickname")
		t.AssertNil(err)
		t.Assert(value, "third")
	})
}

// Test_Model_Batch_EmptyBatch tests batch operation with empty data
func Test_Model_Batch_EmptyBatch(t *testing.T) {
	table := createTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		// Empty slice
		data := g.Slice{}

		// Batch insert with empty data should return error
		_, err := db.Model(table).Batch(10).Data(data).Insert()
		t.AssertNE(err, nil)
		t.AssertIN(err.Error(), "data list cannot be empty")

		// Verify no records inserted
		count, err := db.Model(table).Count()
		t.AssertNil(err)
		t.Assert(count, 0)
	})
}

// Test_Model_Batch_SingleRecord tests batch operation with single record
func Test_Model_Batch_SingleRecord(t *testing.T) {
	table := createTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		// Single record batch insert
		data := g.Slice{
			g.Map{
				"id":       1,
				"passport": "single_user",
				"nickname": "single_name",
			},
		}

		_, err := db.Model(table).Batch(10).Data(data).Insert()
		t.AssertNil(err)
		count, err := db.Model(table).Count()
		t.AssertNil(err)
		t.Assert(count, 1)

		// Verify the record
		one, err := db.Model(table).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["passport"], "single_user")
		t.Assert(one["nickname"], "single_name")
	})
}

// Test_Model_Batch_VsBatch tests performance comparison between different batch sizes
func Test_Model_Batch_VsBatch(t *testing.T) {
	table := createTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		// Prepare data
		data := g.Slice{}
		for i := 1; i <= 100; i++ {
			data = append(data, g.Map{
				"id":       i,
				"passport": fmt.Sprintf("perf_user_%d", i),
			})
		}

		// Test with batch size 1
		_, err := db.Model(table).Batch(1).Data(data).Insert()
		t.AssertNil(err)
		count, err := db.Model(table).Count()
		t.AssertNil(err)
		t.Assert(count, 100)

		// Clean up
		_, err = db.Model(table).Where("1=1").Delete()
		t.AssertNil(err)

		// Test with batch size 10
		_, err = db.Model(table).Batch(10).Data(data).Insert()
		t.AssertNil(err)
		count, err = db.Model(table).Count()
		t.AssertNil(err)
		t.Assert(count, 100)

		// Clean up
		_, err = db.Model(table).Where("1=1").Delete()
		t.AssertNil(err)

		// Test with batch size 50
		_, err = db.Model(table).Batch(50).Data(data).Insert()
		t.AssertNil(err)

		// All batch sizes should produce same result
		count, err = db.Model(table).Count()
		t.AssertNil(err)
		t.Assert(count, 100)
	})
}

// Test_Model_Batch_WithTransaction tests that a batch operation within a transaction fails as unsupported.
func Test_Model_Batch_WithTransaction(t *testing.T) {
	table := createTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		data := g.Slice{}
		for i := 1; i <= 50; i++ {
			data = append(data, g.Map{
				"id":       i,
				"passport": fmt.Sprintf("tx_batch_%d", i),
			})
		}

		// Note: ClickHouse has no transactions, so the driver rejects Transaction before calling the function.
		var called bool
		err := db.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
			called = true
			_, err := tx.Model(table).Batch(10).Data(data).Insert()
			return err
		})
		t.AssertNE(err, nil)
		t.Assert(err.Error(), "unsupported method:Transaction")
		t.Assert(called, false)

		count, err := db.Model(table).Count()
		t.AssertNil(err)
		t.Assert(count, 0)
	})
}

// Test_ClickHouse_Insert_RawRecordValues tests the other values of a record with a raw value: a time is
// written exactly in a DateTime, DateTime64 or Date column, binary data as it is, and an Array or a Map
// value fails as unsupported.
func Test_ClickHouse_Insert_RawRecordValues(t *testing.T) {
	table := chCCreateTypeTable(
		"test_raw_values",
		"name String, dt DateTime, ms DateTime64(3), us DateTime64(6), ns DateTime64(9), far DateTime64(3), "+
			"utc DateTime64(3, 'UTC'), day Date, bin String, tags Array(String), attrs Map(String, UInt8)",
	)
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		var (
			ms  = time.Date(2024, 1, 2, 3, 4, 5, 678000000, time.Local)
			us  = time.Date(1969, 12, 31, 23, 59, 59, 123456000, time.Local)
			ns  = time.Date(2024, 1, 2, 3, 4, 5, 678901234, time.Local)
			far = time.Date(2200, 1, 2, 3, 4, 5, 678000000, time.Local)
			utc = time.Date(2024, 1, 2, 3, 4, 5, 678000000, time.UTC)
		)
		_, err := db.Model(table).Data(g.Map{
			"id":   1,
			"name": gdb.Raw("concat('r', 'aw')"),
			"dt":   "2024-01-02 03:04:05",
			"ms":   "2024-01-02 03:04:05.678",
			"us":   us,
			"ns":   gtime.NewFromTime(ns),
			"far":  &far,
			"utc":  utc,
			"day":  gtime.NewFromStr("2024-01-02 12:00:00"),
			"bin":  []byte{0xff, 0x00, 0x27},
		}).Insert()
		t.AssertNil(err)

		one, err := db.Model(table).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["name"].String(), "raw")
		t.Assert(one["dt"].Time().Equal(time.Date(2024, 1, 2, 3, 4, 5, 0, time.Local)), true)
		t.Assert(one["ms"].Time().Equal(ms), true)
		t.Assert(one["us"].Time().Equal(us), true)
		t.Assert(one["ns"].Time().Equal(ns), true)
		t.Assert(one["far"].Time().Equal(far), true)
		t.Assert(one["utc"].Time().Equal(utc), true)
		t.Assert(one["day"].GTime().Format("Y-m-d"), "2024-01-02")
		t.Assert(one["bin"].Bytes(), []byte{0xff, 0x00, 0x27})
	})
	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(table).Data(g.Map{"id": 2, "name": gdb.Raw("'x'"), "tags": []string{"a", "b"}}).Insert()
		t.AssertNE(err, nil)
		t.Assert(gerror.Code(err), gcode.CodeNotSupported)

		_, err = db.Model(table).Data(g.Map{"id": 2, "name": gdb.Raw("'x'"), "attrs": map[string]uint8{"a": 1}}).Insert()
		t.AssertNE(err, nil)
		t.Assert(gerror.Code(err), gcode.CodeNotSupported)

		count, err := db.Model(table).Where("id", 2).Count()
		t.AssertNil(err)
		t.Assert(count, 0)
	})
}

// Test_ClickHouse_Insert_DryRun tests that an insert in dry run sends nothing to the server, on a
// read-only connection that would refuse it, and goes through CatchSQL as a statement a batch.
func Test_ClickHouse_Insert_DryRun(t *testing.T) {
	table := createTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		dryRunDb, err := gdb.New(chANode(func(node *gdb.ConfigNode) {
			node.Extra = "readonly=1"
			node.DryRun = true
		}))
		t.AssertNil(err)
		defer dryRunDb.Close(ctx)

		var n int64
		sqls, err := gdb.CatchSQL(ctx, func(ctx context.Context) error {
			result, err := dryRunDb.Model(table).Ctx(ctx).Data(g.List{
				g.Map{"id": 1, "passport": "user_1", "password": "pass_1", "nickname": "name_1", "create_time": gtime.Now()},
				g.Map{"id": 2, "passport": "user_2", "password": "pass_2", "nickname": "name_2", "create_time": gtime.Now()},
				g.Map{"id": 3, "passport": "user_3", "password": "pass_3", "nickname": "name_3", "create_time": gtime.Now()},
				g.Map{"id": 4, "passport": "user_4", "password": "pass_4", "nickname": gdb.Raw("'raw'"), "create_time": gtime.Now()},
			}).Batch(2).Insert()
			if err != nil {
				return err
			}
			n, err = result.RowsAffected()
			return err
		})
		t.AssertNil(err)
		t.Assert(n, 4)

		var inserts []string
		for _, sql := range sqls {
			if strings.HasPrefix(sql, "INSERT INTO "+table+"(") {
				inserts = append(inserts, sql)
			}
		}
		t.Assert(len(inserts), 3)

		count, err := db.Model(table).Count()
		t.AssertNil(err)
		t.Assert(count, 0)
	})
}

// Test_ClickHouse_Insert_BatchCount tests that the records are inserted in batches of at most the batch
// count, each of them one insert on the server, with raw values or not.
func Test_ClickHouse_Insert_BatchCount(t *testing.T) {
	for _, raw := range []bool{false, true} {
		gtest.C(t, func(t *gtest.T) {
			table := createTable()
			defer dropTable(table)

			list := g.List{}
			for i := 1; i <= 5; i++ {
				record := g.Map{
					"id":          i,
					"passport":    fmt.Sprintf("user_%d", i),
					"password":    fmt.Sprintf("pass_%d", i),
					"nickname":    fmt.Sprintf("name_%d", i),
					"create_time": gtime.Now(),
				}
				if raw {
					record["nickname"] = gdb.Raw(fmt.Sprintf("'name_%d'", i))
				}
				list = append(list, record)
			}
			_, err := db.Model(table).Data(list).Batch(2).Insert()
			t.AssertNil(err)

			count, err := db.Model(table).Count()
			t.AssertNil(err)
			t.Assert(count, 5)

			_, err = db.Exec(ctx, "SYSTEM FLUSH LOGS")
			t.AssertNil(err)
			inserts, err := db.GetValue(ctx,
				"SELECT count() FROM system.query_log WHERE type = 'QueryFinish' AND query_kind = 'Insert' "+
					"AND has(tables, currentDatabase() || '.' || ?)", table,
			)
			t.AssertNil(err)
			t.Assert(inserts.Int(), 3)
		})
	}
}

// Test_ClickHouse_Insert_TimeZones tests times of a named, a loaded and a fixed time zone written to
// DateTime, DateTime64, Date and Date32 columns, with a raw value in the record or not: the instant
// of a time is kept, and its date is the date in its own time zone.
func Test_ClickHouse_Insert_TimeZones(t *testing.T) {
	table := chCCreateTypeTable(
		"test_time_zones",
		"name String, dt DateTime, dt64 DateTime64(3), utc DateTime64(3, 'UTC'), day Date, day32 Date32",
	)
	defer dropTable(table)

	newYork, err := time.LoadLocation("America/New_York")
	gtest.AssertNil(err)
	locations := []*time.Location{time.UTC, newYork, time.FixedZone("CST", -6*3600)}
	for i, location := range locations {
		gtest.C(t, func(t *gtest.T) {
			value := time.Date(2024, 1, 2, 20, 0, 0, 0, location)
			for j, name := range []any{gdb.Raw("'raw'"), "batch"} {
				_, err := db.Model(table).Data(g.Map{
					"id":    i*2 + j + 1,
					"name":  name,
					"dt":    value,
					"dt64":  value,
					"utc":   value,
					"day":   value,
					"day32": value,
				}).Insert()
				t.AssertNil(err)

				one, err := db.Model(table).Where("id", i*2+j+1).One()
				t.AssertNil(err)
				t.Assert(one["dt"].Time().Equal(value), true)
				t.Assert(one["dt64"].Time().Equal(value), true)
				t.Assert(one["utc"].Time().Equal(value), true)
				t.Assert(one["day"].GTime().Format("Y-m-d"), "2024-01-02")
				t.Assert(one["day32"].GTime().Format("Y-m-d"), "2024-01-02")
			}
		})
	}
}

// Test_ClickHouse_Insert_RawRecordTimeRange tests that a time of a record with a raw value is written to
// a DateTime, DateTime64, Date or Date32 column up to the bounds of the column, and fails beyond them.
func Test_ClickHouse_Insert_RawRecordTimeRange(t *testing.T) {
	table := chCCreateTypeTable("test_raw_time_range", "name String, dt DateTime, dt64 DateTime64(9), d Date, d32 Date32")
	defer dropTable(table)

	maxDateTime64 := time.Unix(0, math.MaxInt64).UTC()
	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(table).Data(g.List{
			g.Map{"id": 1, "name": gdb.Raw("'min'"), "dt": time.Unix(0, 0).UTC(), "dt64": time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC), "d": "1970-01-01", "d32": "1900-01-01"},
			g.Map{"id": 2, "name": gdb.Raw("'max'"), "dt": time.Unix(math.MaxUint32, 0).UTC(), "dt64": maxDateTime64, "d": "2149-06-06", "d32": "2299-12-31"},
		}).Insert()
		t.AssertNil(err)

		all, err := db.Model(table).OrderAsc("id").All()
		t.AssertNil(err)
		t.Assert(len(all), 2)
		t.Assert(all[0]["dt"].Time().Equal(time.Unix(0, 0)), true)
		t.Assert(all[0]["dt64"].Time().Equal(time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC)), true)
		t.Assert(all[0]["d"].GTime().Format("Y-m-d"), "1970-01-01")
		t.Assert(all[0]["d32"].GTime().Format("Y-m-d"), "1900-01-01")
		t.Assert(all[1]["dt"].Time().Equal(time.Unix(math.MaxUint32, 0)), true)
		t.Assert(all[1]["dt64"].Time().Equal(maxDateTime64), true)
		t.Assert(all[1]["d"].GTime().Format("Y-m-d"), "2149-06-06")
		t.Assert(all[1]["d32"].GTime().Format("Y-m-d"), "2299-12-31")
	})
	gtest.C(t, func(t *gtest.T) {
		outOfRange := []g.Map{
			{"dt": time.Date(2200, 1, 1, 0, 0, 0, 0, time.UTC)},
			{"dt": time.Unix(-1, 0).UTC()},
			{"dt64": time.Date(1800, 1, 1, 0, 0, 0, 0, time.UTC)},
			{"dt64": maxDateTime64.Add(time.Nanosecond)},
			{"d": "2149-06-07"},
			{"d32": "1899-12-31"},
		}
		for i, data := range outOfRange {
			data["id"] = i + 3
			data["name"] = gdb.Raw("'out'")
			_, err := db.Model(table).Data(data).Insert()
			t.AssertNE(err, nil)
			t.Assert(gerror.Code(err), gcode.CodeInvalidParameter)
		}
		count, err := db.Model(table).Count()
		t.AssertNil(err)
		t.Assert(count, 2)
	})
}

// Test_ClickHouse_Insert_ToSQL tests that ToSQL returns the insert of records without executing it, in
// one batch or several, and that CatchSQL catches the insert of every batch while executing it.
func Test_ClickHouse_Insert_ToSQL(t *testing.T) {
	newList := func() g.List {
		list := g.List{}
		for i := 1; i <= 3; i++ {
			list = append(list, g.Map{
				"id":          i,
				"passport":    fmt.Sprintf("user_%d", i),
				"password":    fmt.Sprintf("pass_%d", i),
				"nickname":    fmt.Sprintf("name_%d", i),
				"create_time": gtime.Now(),
			})
		}
		return list
	}
	for _, batch := range []int{0, 2} {
		gtest.C(t, func(t *gtest.T) {
			table := createTable()
			defer dropTable(table)

			sql, err := gdb.ToSQL(ctx, func(ctx context.Context) error {
				_, err := db.Model(table).Ctx(ctx).Data(newList()).Batch(batch).Insert()
				return err
			})
			t.AssertNil(err)

			count, err := db.Model(table).Count()
			t.AssertNil(err)
			t.Assert(count, 0)
			t.Assert(strings.HasPrefix(sql, "INSERT INTO "+table+"("), true)
		})
	}
	gtest.C(t, func(t *gtest.T) {
		table := createTable()
		defer dropTable(table)

		sqls, err := gdb.CatchSQL(ctx, func(ctx context.Context) error {
			_, err := db.Model(table).Ctx(ctx).Data(newList()).Batch(2).Insert()
			return err
		})
		t.AssertNil(err)

		var inserts []string
		for _, sql := range sqls {
			if strings.HasPrefix(sql, "INSERT INTO "+table+"(") {
				inserts = append(inserts, sql)
			}
		}
		t.Assert(len(inserts), 2)
		t.Assert(strings.Contains(inserts[0], "'user_1'"), true)

		count, err := db.Model(table).Count()
		t.AssertNil(err)
		t.Assert(count, 3)
	})
}

// Test_ClickHouse_Insert_CommitError tests that a batch the server refuses when it is committed fails with
// an error of code gcode.CodeDbOperationError carrying the INSERT statement, as a statement does.
func Test_ClickHouse_Insert_CommitError(t *testing.T) {
	table := chCCreateTypeTable("test_commit_error", "name String, CONSTRAINT small_id CHECK id < 100")
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(table).Data(g.List{
			g.Map{"id": 1, "name": "a"},
			g.Map{"id": 200, "name": "b"},
		}).Insert()
		t.AssertNE(err, nil)
		t.Assert(gerror.Code(err), gcode.CodeDbOperationError)
		t.Assert(strings.Contains(err.Error(), "INSERT INTO "+table+"(id,name) VALUES (1,'a'),(200,'b')"), true)

		count, err := db.Model(table).Count()
		t.AssertNil(err)
		t.Assert(count, 0)
	})
}

// Test_ClickHouse_Insert_LogPerBatch tests that the insert of records logs one statement a batch, as a
// statement does, rather than one a record.
func Test_ClickHouse_Insert_LogPerBatch(t *testing.T) {
	table := createTable()
	defer dropTable(table)

	logDb, err := gdb.New(chANode(func(node *gdb.ConfigNode) {
		node.Debug = true
	}))
	gtest.AssertNil(err)
	defer logDb.Close(ctx)

	var (
		logs   []string
		logger = glog.New()
	)
	logger.SetStdoutPrint(false)
	logger.SetHandlers(func(ctx context.Context, in *glog.HandlerInput) {
		if content := in.Content + in.ValuesContent(); strings.Contains(content, "INSERT INTO "+table+"(") {
			logs = append(logs, content)
		}
		in.Next(ctx)
	})
	logDb.SetLogger(logger)

	for i, batch := range []int{0, 2} {
		gtest.C(t, func(t *gtest.T) {
			logs = logs[:0]
			list := g.List{}
			for j := 1; j <= 3; j++ {
				id := i*10 + j
				list = append(list, g.Map{
					"id":          id,
					"passport":    fmt.Sprintf("user_%d", id),
					"password":    fmt.Sprintf("pass_%d", id),
					"nickname":    fmt.Sprintf("name_%d", id),
					"create_time": gtime.Now(),
				})
			}
			_, err := logDb.Model(table).Data(list).Batch(batch).Insert()
			t.AssertNil(err)
			if batch == 0 {
				t.Assert(len(logs), 1)
			} else {
				t.Assert(len(logs), 2)
			}
		})
	}
}

// Test_ClickHouse_Insert_ExecTimeout tests that the ExecTimeout of the configuration limits the insert of
// a batch.
func Test_ClickHouse_Insert_ExecTimeout(t *testing.T) {
	table := createTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		timeoutDb, err := gdb.New(chANode(func(node *gdb.ConfigNode) {
			node.ExecTimeout = time.Nanosecond
		}))
		t.AssertNil(err)
		defer timeoutDb.Close(ctx)

		_, err = timeoutDb.Model(table).Data(g.Map{
			"id":          1,
			"passport":    "user_1",
			"password":    "pass_1",
			"nickname":    "name_1",
			"create_time": gtime.Now(),
		}).Insert()
		t.AssertNE(err, nil)

		count, err := db.Model(table).Count()
		t.AssertNil(err)
		t.Assert(count, 0)
	})
}

// Test_ClickHouse_Insert_ToSQLArray tests that ToSQL returns the insert of a record with an Array value, empty
// or not, as CatchSQL catches it, without executing it.
func Test_ClickHouse_Insert_ToSQLArray(t *testing.T) {
	table := chCCreateTypeTable("test_tosql_array", "tags Array(String)")
	defer dropTable(table)

	for i, tags := range [][]string{{"q", "r"}, {}} {
		gtest.C(t, func(t *gtest.T) {
			data := g.Map{"id": i + 1, "tags": tags}
			sql, err := gdb.ToSQL(ctx, func(ctx context.Context) error {
				_, err := db.Model(table).Ctx(ctx).Data(data).Insert()
				return err
			})
			t.AssertNil(err)

			count, err := db.Model(table).Where("id", i+1).Count()
			t.AssertNil(err)
			t.Assert(count, 0)

			sqls, err := gdb.CatchSQL(ctx, func(ctx context.Context) error {
				_, err := db.Model(table).Ctx(ctx).Data(data).Insert()
				return err
			})
			t.AssertNil(err)
			t.Assert(sql, sqls[len(sqls)-1])
			t.Assert(strings.HasPrefix(sql, "INSERT INTO "+table+"(id,tags) VALUES ("), true)

			count, err = db.Model(table).Where("id", i+1).Count()
			t.AssertNil(err)
			t.Assert(count, 1)
		})
	}
}
