// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package oracle_test

import (
	"fmt"
	"sync"
	"testing"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/test/gtest"
)

// Test_Concurrent_Insert tests concurrent Insert operations
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
			go func(id int) {
				defer wg.Done()
				_, errs[id-1] = db.Model(table).Insert(g.Map{
					"passport": fmt.Sprintf("user_%d", id),
					"password": fmt.Sprintf("pass_%d", id),
					"nickname": fmt.Sprintf("name_%d", id),
				})
			}(i + 1)
		}
		wg.Wait()
		for _, err := range errs {
			t.AssertNil(err)
		}

		count, err := db.Model(table).Count()
		t.AssertNil(err)
		t.Assert(count, concurrency)
	})
}

// Test_Concurrent_Update tests concurrent Update operations
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
			go func(id int) {
				defer wg.Done()
				_, errs[id] = db.Model(table).Data(g.Map{
					"nickname": fmt.Sprintf("updated_%d", id),
				}).Where("id", id+1).Update()
			}(i)
		}
		wg.Wait()
		for _, err := range errs {
			t.AssertNil(err)
		}

		for i := 0; i < concurrency; i++ {
			one, err := db.Model(table).Where("id", i+1).One()
			t.AssertNil(err)
			t.Assert(one["NICKNAME"].String(), fmt.Sprintf("updated_%d", i))
		}
	})
}

// Test_Concurrent_Delete tests concurrent Delete operations
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
			go func(id int) {
				defer wg.Done()
				_, errs[id] = db.Model(table).Where("id", id+1).Delete()
			}(i)
		}
		wg.Wait()
		for _, err := range errs {
			t.AssertNil(err)
		}

		count, err := db.Model(table).Count()
		t.AssertNil(err)
		t.Assert(count, TableSize-concurrency)
	})
}

// Test_Concurrent_Query tests concurrent Query operations
func Test_Concurrent_Query(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		var (
			wg          sync.WaitGroup
			concurrency = 20
			results     = make([]gdb.Record, concurrency)
			errs        = make([]error, concurrency)
		)
		wg.Add(concurrency)
		for i := 0; i < concurrency; i++ {
			go func(id int) {
				defer wg.Done()
				results[id], errs[id] = db.Model(table).Where("id", (id%TableSize)+1).One()
			}(i)
		}
		wg.Wait()
		for i := 0; i < concurrency; i++ {
			t.AssertNil(errs[i])
			t.Assert(results[i]["ID"], (i%TableSize)+1)
		}
	})
}

// Test_Concurrent_Transaction tests concurrent transaction operations
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
			go func(id int) {
				defer wg.Done()
				errs[id-1] = db.Transaction(ctx, func(ctx g.Ctx, tx gdb.TX) error {
					_, err := tx.Model(table).Insert(g.Map{
						"passport": fmt.Sprintf("user_%d", id),
						"password": fmt.Sprintf("pass_%d", id),
						"nickname": fmt.Sprintf("name_%d", id),
					})
					return err
				})
			}(i + 1)
		}
		wg.Wait()
		for _, err := range errs {
			t.AssertNil(err)
		}

		count, err := db.Model(table).Count()
		t.AssertNil(err)
		t.Assert(count, concurrency)
	})
}

// Test_Concurrent_Mixed_Operations tests mixed concurrent operations
func Test_Concurrent_Mixed_Operations(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		var (
			wg         sync.WaitGroup
			operations = 30
			errs       = make([]error, operations)
		)
		wg.Add(operations)
		for i := 0; i < operations; i++ {
			op := i % 3
			switch op {
			case 0:
				go func(id int) {
					defer wg.Done()
					_, errs[id] = db.Model(table).Insert(g.Map{
						"id":       TableSize + id + 1,
						"passport": fmt.Sprintf("new_user_%d", id),
						"password": fmt.Sprintf("new_pass_%d", id),
						"nickname": fmt.Sprintf("new_name_%d", id),
					})
				}(i)
			case 1:
				go func(id int) {
					defer wg.Done()
					targetId := (id % TableSize) + 1
					_, errs[id] = db.Model(table).Data(g.Map{
						"nickname": fmt.Sprintf("concurrent_%d", id),
					}).Where("id", targetId).Update()
				}(i)
			case 2:
				go func(id int) {
					defer wg.Done()
					targetId := (id % TableSize) + 1
					_, errs[id] = db.Model(table).Where("id", targetId).One()
				}(i)
			}
		}
		wg.Wait()
		for _, err := range errs {
			t.AssertNil(err)
		}

		count, err := db.Model(table).Count()
		t.AssertNil(err)
		t.Assert(count, TableSize+operations/3)
	})
}

// Test_Concurrent_Connection_Pool tests connection pool under load
func Test_Concurrent_Connection_Pool(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		var (
			wg          sync.WaitGroup
			concurrency = 50
			errs        = make([]error, concurrency)
		)
		wg.Add(concurrency)
		for i := 0; i < concurrency; i++ {
			go func(id int) {
				defer wg.Done()
				for j := 0; j < 5 && errs[id] == nil; j++ {
					_, errs[id] = db.Model(table).Where("id", (id%TableSize)+1).One()
				}
			}(i)
		}
		wg.Wait()
		for _, err := range errs {
			t.AssertNil(err)
		}
	})
}

// Test_Concurrent_Schema_Switch tests concurrent schema switching
func Test_Concurrent_Schema_Switch(t *testing.T) {
	table1 := createTableWithDb(db, "test_schema_1")
	table2 := createTableWithDb(dblink, "test_schema_2")
	defer dropTableWithDb(db, table1)
	defer dropTableWithDb(dblink, table2)

	gtest.C(t, func(t *gtest.T) {
		var (
			wg          sync.WaitGroup
			concurrency = 10
			errs1       = make([]error, concurrency)
			errs2       = make([]error, concurrency)
		)
		wg.Add(concurrency * 2)
		for i := 0; i < concurrency; i++ {
			go func(id int) {
				defer wg.Done()
				_, errs1[id] = db.Model(table1).Insert(g.Map{
					"passport": fmt.Sprintf("user_s1_%d", id),
					"password": fmt.Sprintf("pass_%d", id),
					"nickname": fmt.Sprintf("name_%d", id),
				})
			}(i)

			go func(id int) {
				defer wg.Done()
				_, errs2[id] = dblink.Model(table2).Insert(g.Map{
					"passport": fmt.Sprintf("user_s2_%d", id),
					"password": fmt.Sprintf("pass_%d", id),
					"nickname": fmt.Sprintf("name_%d", id),
				})
			}(i)
		}
		wg.Wait()
		for i := 0; i < concurrency; i++ {
			t.AssertNil(errs1[i])
			t.AssertNil(errs2[i])
		}

		count1, err := db.Model(table1).Count()
		t.AssertNil(err)
		t.Assert(count1, concurrency)

		count2, err := dblink.Model(table2).Count()
		t.AssertNil(err)
		t.Assert(count2, concurrency)
	})
}

// Test_Concurrent_Model_Clone tests concurrent model cloning
func Test_Concurrent_Model_Clone(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		var (
			baseModel   = db.Model(table).Where("id>", 0)
			wg          sync.WaitGroup
			concurrency = 20
			results     = make([]gdb.Result, concurrency)
			errs        = make([]error, concurrency)
		)
		wg.Add(concurrency)
		for i := 0; i < concurrency; i++ {
			go func(id int) {
				defer wg.Done()
				m := baseModel.Clone()
				results[id], errs[id] = m.Where("id<=", TableSize/2).All()
			}(i)
		}
		wg.Wait()
		for i := 0; i < concurrency; i++ {
			t.AssertNil(errs[i])
			t.Assert(len(results[i]), TableSize/2)
		}
	})
}

// Test_Concurrent_Batch_Insert tests concurrent batch insert operations
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
			go func(batchId int) {
				defer wg.Done()
				batch := make([]g.Map, 0, batchSize)
				for j := 0; j < batchSize; j++ {
					id := batchId*batchSize + j
					batch = append(batch, g.Map{
						"passport": fmt.Sprintf("batch_user_%d", id),
						"password": fmt.Sprintf("pass_%d", id),
						"nickname": fmt.Sprintf("name_%d", id),
					})
				}
				_, errs[batchId] = db.Model(table).Data(batch).Insert()
			}(i)
		}
		wg.Wait()
		for _, err := range errs {
			t.AssertNil(err)
		}

		count, err := db.Model(table).Count()
		t.AssertNil(err)
		t.Assert(count, concurrency*batchSize)
	})
}
