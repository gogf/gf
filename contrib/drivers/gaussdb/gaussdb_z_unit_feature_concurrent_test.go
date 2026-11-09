// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gaussdb_test

import (
	"fmt"
	"sync"
	"testing"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/test/gtest"
)

// Note: the worker goroutines below only record their outcome; every assertion runs
// on the test goroutine after Wait(). gtest assertions panic on failure, and a panic
// raised in a spawned goroutine is not recovered by gtest.C — it would abort the whole
// test binary and skip every remaining test rather than failing just this one.

// Test_Concurrent_Insert tests concurrent Insert operations.
func Test_Concurrent_Insert(t *testing.T) {
	table := createTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		var (
			wg          sync.WaitGroup
			concurrency = 10
			errs        = make([]error, concurrency)
		)
		wg.Add(concurrency)
		for i := 0; i < concurrency; i++ {
			go func(idx, id int) {
				defer wg.Done()
				_, errs[idx] = db.Model(table).Insert(g.Map{
					"passport":    fmt.Sprintf("user_%d", id),
					"password":    fmt.Sprintf("pass_%d", id),
					"nickname":    fmt.Sprintf("name_%d", id),
					"create_time": gtime.Now().String(),
				})
			}(i, i+1)
		}
		wg.Wait()

		for _, err := range errs {
			t.AssertNil(err)
		}

		// Verify all records inserted
		count, err := db.Model(table).Count()
		t.AssertNil(err)
		t.Assert(count, concurrency)
	})
}

// Test_Concurrent_Update tests concurrent Update operations.
func Test_Concurrent_Update(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		var (
			wg          sync.WaitGroup
			concurrency = 5
			errs        = make([]error, concurrency)
		)
		wg.Add(concurrency)
		for i := 0; i < concurrency; i++ {
			go func(idx, id int) {
				defer wg.Done()
				_, errs[idx] = db.Model(table).Data(g.Map{
					"nickname": fmt.Sprintf("updated_%d", id),
				}).Where("id", id+1).Update()
			}(i, i)
		}
		wg.Wait()

		for _, err := range errs {
			t.AssertNil(err)
		}

		// Verify updates
		for i := 0; i < concurrency; i++ {
			one, err := db.Model(table).Where("id", i+1).One()
			t.AssertNil(err)
			t.Assert(one["nickname"].String(), fmt.Sprintf("updated_%d", i))
		}
	})
}

// Test_Concurrent_Delete tests concurrent Delete operations.
func Test_Concurrent_Delete(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		var (
			wg          sync.WaitGroup
			concurrency = 5
			errs        = make([]error, concurrency)
		)
		wg.Add(concurrency)
		for i := 0; i < concurrency; i++ {
			go func(idx, id int) {
				defer wg.Done()
				_, errs[idx] = db.Model(table).Where("id", id+1).Delete()
			}(i, i)
		}
		wg.Wait()

		for _, err := range errs {
			t.AssertNil(err)
		}

		// Verify deletions
		count, err := db.Model(table).Count()
		t.AssertNil(err)
		t.Assert(count, TableSize-concurrency)
	})
}

// Test_Concurrent_Query tests concurrent Query operations.
func Test_Concurrent_Query(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		var (
			wg          sync.WaitGroup
			concurrency = 20
			errs        = make([]error, concurrency)
			results     = make([]gdb.Record, concurrency)
		)
		wg.Add(concurrency)
		for i := 0; i < concurrency; i++ {
			go func(idx, id int) {
				defer wg.Done()
				results[idx], errs[idx] = db.Model(table).Where("id", (id%TableSize)+1).One()
			}(i, i)
		}
		wg.Wait()

		for i := 0; i < concurrency; i++ {
			t.AssertNil(errs[i])
			t.AssertNE(results[i], nil)
			t.Assert(results[i]["id"].Int(), (i%TableSize)+1)
		}
	})
}

// Test_Concurrent_Transaction tests concurrent transaction operations.
func Test_Concurrent_Transaction(t *testing.T) {
	table := createTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		var (
			wg          sync.WaitGroup
			concurrency = 10
			errs        = make([]error, concurrency)
		)
		wg.Add(concurrency)
		for i := 0; i < concurrency; i++ {
			go func(idx, id int) {
				defer wg.Done()
				errs[idx] = db.Transaction(ctx, func(ctx g.Ctx, tx gdb.TX) error {
					_, err := tx.Model(table).Insert(g.Map{
						"passport":    fmt.Sprintf("user_%d", id),
						"password":    fmt.Sprintf("pass_%d", id),
						"nickname":    fmt.Sprintf("name_%d", id),
						"create_time": gtime.Now().String(),
					})
					return err
				})
			}(i, i+1)
		}
		wg.Wait()

		for _, err := range errs {
			t.AssertNil(err)
		}

		// Verify all transactions committed
		count, err := db.Model(table).Count()
		t.AssertNil(err)
		t.Assert(count, concurrency)
	})
}

// Test_Concurrent_Mixed_Operations tests mixed concurrent operations.
func Test_Concurrent_Mixed_Operations(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		// createInitTable seeds its rows with explicit ids, which leaves the bigserial
		// sequence at its initial value. Advance it past the seeded rows so the inserts
		// below cannot collide with them and every operation is expected to succeed.
		_, err := db.Exec(ctx, fmt.Sprintf(
			`SELECT setval(pg_get_serial_sequence('%s', 'id'), %d)`, table, TableSize,
		))
		t.AssertNil(err)

		var (
			wg         sync.WaitGroup
			operations = 30
			inserts    int
			errs       = make([]error, operations)
			queried    = make([]gdb.Record, operations)
			// Row id -> the nickname its updater writes. The schedule below gives every
			// seeded row exactly one updater, so each write is deterministically visible.
			updated = make(map[int]string)
		)
		wg.Add(operations)
		for i := 0; i < operations; i++ {
			switch i % 3 {
			case 0: // Insert
				inserts++
				go func(idx, id int) {
					defer wg.Done()
					_, errs[idx] = db.Model(table).Insert(g.Map{
						"passport":    fmt.Sprintf("new_user_%d", id),
						"password":    fmt.Sprintf("new_pass_%d", id),
						"nickname":    fmt.Sprintf("new_name_%d", id),
						"create_time": gtime.Now().String(),
					})
				}(i, i)
			case 1: // Update
				updated[(i%TableSize)+1] = fmt.Sprintf("concurrent_%d", i)
				go func(idx, id int) {
					defer wg.Done()
					_, errs[idx] = db.Model(table).Data(g.Map{
						"nickname": fmt.Sprintf("concurrent_%d", id),
					}).Where("id", (id%TableSize)+1).Update()
				}(i, i)
			case 2: // Query
				go func(idx, id int) {
					defer wg.Done()
					queried[idx], errs[idx] = db.Model(table).Where("id", (id%TableSize)+1).One()
				}(i, i)
			}
		}
		wg.Wait()

		// Every operation succeeded, and each query returned the row it asked for.
		for i, err := range errs {
			t.AssertNil(err)
			if i%3 == 2 {
				t.AssertNE(queried[i], nil)
				t.Assert(queried[i]["id"].Int(), (i%TableSize)+1)
			}
		}

		// Every insert landed, so the row count is exact rather than a lower bound.
		count, err := db.Model(table).Count()
		t.AssertNil(err)
		t.Assert(count, TableSize+inserts)

		// The pre-existing rows are all still readable, and each carries the nickname
		// written by its concurrent updater.
		all, err := db.Model(table).Where("id<=?", TableSize).Order("id").All()
		t.AssertNil(err)
		t.Assert(len(all), TableSize)
		t.Assert(len(updated), TableSize)
		for _, v := range all {
			t.Assert(v["nickname"].String(), updated[v["id"].Int()])
		}
	})
}

// Test_Concurrent_Connection_Pool tests connection pool under load.
func Test_Concurrent_Connection_Pool(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		var (
			wg              sync.WaitGroup
			concurrency     = 50
			queriesPerRound = 5
			errs            = make([]error, concurrency)
		)
		wg.Add(concurrency)
		for i := 0; i < concurrency; i++ {
			go func(idx, id int) {
				defer wg.Done()
				// Each goroutine performs multiple operations
				for j := 0; j < queriesPerRound; j++ {
					if _, err := db.Model(table).Where("id", (id%TableSize)+1).One(); err != nil {
						errs[idx] = err
						return
					}
				}
			}(i, i)
		}
		wg.Wait()

		for _, err := range errs {
			t.AssertNil(err)
		}
	})
}

// Test_Concurrent_Schema_Switch tests concurrent operations across two tables.
// GaussDB-adapted from MariaDB: the MariaDB version uses db2 (a separate schema).
// GaussDB init_test.go only configures a single DB/schema, so we use two distinct
// tables within the same schema to exercise the same concurrency pathway.
func Test_Concurrent_Schema_Switch(t *testing.T) {
	table1 := createTableWithDb(db, TablePrefix+"schema_1")
	table2 := createTableWithDb(db, TablePrefix+"schema_2")
	defer dropTableWithDb(db, table1)
	defer dropTableWithDb(db, table2)

	gtest.C(t, func(t *gtest.T) {
		var (
			wg          sync.WaitGroup
			concurrency = 10
			errs1       = make([]error, concurrency)
			errs2       = make([]error, concurrency)
		)
		wg.Add(concurrency * 2)
		for i := 0; i < concurrency; i++ {
			// Insert to table1
			go func(idx int) {
				defer wg.Done()
				_, errs1[idx] = db.Model(table1).Insert(g.Map{
					"passport":    fmt.Sprintf("user_s1_%d", idx),
					"password":    fmt.Sprintf("pass_%d", idx),
					"nickname":    fmt.Sprintf("name_%d", idx),
					"create_time": gtime.Now().String(),
				})
			}(i)

			// Insert to table2
			go func(idx int) {
				defer wg.Done()
				_, errs2[idx] = db.Model(table2).Insert(g.Map{
					"passport":    fmt.Sprintf("user_s2_%d", idx),
					"password":    fmt.Sprintf("pass_%d", idx),
					"nickname":    fmt.Sprintf("name_%d", idx),
					"create_time": gtime.Now().String(),
				})
			}(i)
		}
		wg.Wait()

		for i := 0; i < concurrency; i++ {
			t.AssertNil(errs1[i])
			t.AssertNil(errs2[i])
		}

		// Verify both tables
		count1, err := db.Model(table1).Count()
		t.AssertNil(err)
		t.Assert(count1, concurrency)

		count2, err := db.Model(table2).Count()
		t.AssertNil(err)
		t.Assert(count2, concurrency)
	})
}

// Test_Concurrent_Model_Clone tests concurrent model cloning.
func Test_Concurrent_Model_Clone(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		var (
			baseModel   = db.Model(table).Where("id>", 0)
			wg          sync.WaitGroup
			concurrency = 20
			errs        = make([]error, concurrency)
			counts      = make([]int, concurrency)
		)
		wg.Add(concurrency)
		for i := 0; i < concurrency; i++ {
			go func(idx int) {
				defer wg.Done()
				// Clone model for each goroutine
				result, err := baseModel.Clone().Where("id<=", TableSize/2).All()
				errs[idx] = err
				counts[idx] = len(result)
			}(i)
		}
		wg.Wait()

		for i := 0; i < concurrency; i++ {
			t.AssertNil(errs[i])
			// Cloning must not let the goroutines leak conditions into each other,
			// so every clone sees exactly the same row set.
			t.Assert(counts[i], TableSize/2)
		}
	})
}

// Test_Concurrent_Batch_Insert tests concurrent batch insert operations.
func Test_Concurrent_Batch_Insert(t *testing.T) {
	table := createTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		var (
			wg          sync.WaitGroup
			concurrency = 5
			batchSize   = 10
			errs        = make([]error, concurrency)
		)
		wg.Add(concurrency)
		for i := 0; i < concurrency; i++ {
			go func(idx, batchId int) {
				defer wg.Done()
				batch := make([]g.Map, 0, batchSize)
				for j := 0; j < batchSize; j++ {
					id := batchId*batchSize + j
					batch = append(batch, g.Map{
						"passport":    fmt.Sprintf("batch_user_%d", id),
						"password":    fmt.Sprintf("pass_%d", id),
						"nickname":    fmt.Sprintf("name_%d", id),
						"create_time": gtime.Now().String(),
					})
				}
				_, errs[idx] = db.Model(table).Data(batch).Insert()
			}(i, i)
		}
		wg.Wait()

		for _, err := range errs {
			t.AssertNil(err)
		}

		// Verify all batch inserts
		count, err := db.Model(table).Count()
		t.AssertNil(err)
		t.Assert(count, concurrency*batchSize)
	})
}
