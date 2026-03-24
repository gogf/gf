// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package mariadb_test

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
		var wg sync.WaitGroup
		concurrency := 10

		wg.Add(concurrency)
		for i := 0; i < concurrency; i++ {
			go func(id int) {
				defer wg.Done()
				_, err := db.Model(table).Insert(g.Map{
					"passport": fmt.Sprintf("user_%d", id),
					"password": fmt.Sprintf("pass_%d", id),
					"nickname": fmt.Sprintf("name_%d", id),
				})
				t.AssertNil(err)
			}(i + 1)
		}
		wg.Wait()

		// Verify all records inserted
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
		var wg sync.WaitGroup
		concurrency := 5

		wg.Add(concurrency)
		for i := 0; i < concurrency; i++ {
			go func(id int) {
				defer wg.Done()
				_, err := db.Model(table).Data(g.Map{
					"nickname": fmt.Sprintf("updated_%d", id),
				}).Where("id", id+1).Update()
				t.AssertNil(err)
			}(i)
		}
		wg.Wait()

		// Verify updates
		for i := 0; i < concurrency; i++ {
			one, err := db.Model(table).Where("id", i+1).One()
			t.AssertNil(err)
			t.Assert(one["nickname"].String(), fmt.Sprintf("updated_%d", i))
		}
	})
}

// Test_Concurrent_Delete tests concurrent Delete operations
func Test_Concurrent_Delete(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		var wg sync.WaitGroup
		concurrency := 5

		wg.Add(concurrency)
		for i := 0; i < concurrency; i++ {
			go func(id int) {
				defer wg.Done()
				_, err := db.Model(table).Where("id", id+1).Delete()
				t.AssertNil(err)
			}(i)
		}
		wg.Wait()

		// Verify deletions
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
		var wg sync.WaitGroup
		concurrency := 20

		wg.Add(concurrency)
		for i := 0; i < concurrency; i++ {
			go func(id int) {
				defer wg.Done()
				result, err := db.Model(table).Where("id", (id%TableSize)+1).One()
				t.AssertNil(err)
				t.AssertNE(result, nil)
			}(i)
		}
		wg.Wait()
	})
}

// Test_Concurrent_Transaction tests concurrent transaction operations
func Test_Concurrent_Transaction(t *testing.T) {
	table := createTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		var wg sync.WaitGroup
		concurrency := 10

		wg.Add(concurrency)
		for i := 0; i < concurrency; i++ {
			go func(id int) {
				defer wg.Done()
				err := db.Transaction(ctx, func(ctx g.Ctx, tx gdb.TX) error {
					_, err := tx.Model(table).Insert(g.Map{
						"passport": fmt.Sprintf("user_%d", id),
						"password": fmt.Sprintf("pass_%d", id),
						"nickname": fmt.Sprintf("name_%d", id),
					})
					return err
				})
				t.AssertNil(err)
			}(i + 1)
		}
		wg.Wait()

		// Verify all transactions committed
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
		var wg sync.WaitGroup
		operations := 30

		wg.Add(operations)
		for i := 0; i < operations; i++ {
			op := i % 3
			switch op {
			case 0: // Insert
				go func(id int) {
					defer wg.Done()
					_, _ = db.Model(table).Insert(g.Map{
						"passport": fmt.Sprintf("new_user_%d", id),
						"password": fmt.Sprintf("new_pass_%d", id),
						"nickname": fmt.Sprintf("new_name_%d", id),
					})
				}(i)
			case 1: // Update
				go func(id int) {
					defer wg.Done()
					targetId := (id % TableSize) + 1
					_, _ = db.Model(table).Data(g.Map{
						"nickname": fmt.Sprintf("concurrent_%d", id),
					}).Where("id", targetId).Update()
				}(i)
			case 2: // Query
				go func(id int) {
					defer wg.Done()
					targetId := (id % TableSize) + 1
					_, _ = db.Model(table).Where("id", targetId).One()
				}(i)
			}
		}
		wg.Wait()

		// Verify database is still consistent
		count, err := db.Model(table).Count()
		t.AssertNil(err)
		t.AssertGT(count, TableSize)
	})
}

// Test_Concurrent_Connection_Pool tests connection pool under load
func Test_Concurrent_Connection_Pool(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		var wg sync.WaitGroup
		concurrency := 50

		wg.Add(concurrency)
		for i := 0; i < concurrency; i++ {
			go func(id int) {
				defer wg.Done()
				// Each goroutine performs multiple operations
				for j := 0; j < 5; j++ {
					_, err := db.Model(table).Where("id", (id%TableSize)+1).One()
					t.AssertNil(err)
				}
			}(i)
		}
		wg.Wait()
	})
}

// Test_Concurrent_Schema_Switch tests concurrent schema switching
func Test_Concurrent_Schema_Switch(t *testing.T) {
	table1 := createTableWithDb(db, "test_schema_1")
	table2 := createTableWithDb(db2, "test_schema_2")
	defer dropTableWithDb(db, table1)
	defer dropTableWithDb(db2, table2)

	gtest.C(t, func(t *gtest.T) {
		var wg sync.WaitGroup
		concurrency := 10

		wg.Add(concurrency * 2)
		for i := 0; i < concurrency; i++ {
			// Insert to schema1
			go func(id int) {
				defer wg.Done()
				_, err := db.Model(table1).Insert(g.Map{
					"passport": fmt.Sprintf("user_s1_%d", id),
					"password": fmt.Sprintf("pass_%d", id),
					"nickname": fmt.Sprintf("name_%d", id),
				})
				t.AssertNil(err)
			}(i)

			// Insert to schema2
			go func(id int) {
				defer wg.Done()
				_, err := db2.Model(table2).Insert(g.Map{
					"passport": fmt.Sprintf("user_s2_%d", id),
					"password": fmt.Sprintf("pass_%d", id),
					"nickname": fmt.Sprintf("name_%d", id),
				})
				t.AssertNil(err)
			}(i)
		}
		wg.Wait()

		// Verify both schemas
		count1, err := db.Model(table1).Count()
		t.AssertNil(err)
		t.Assert(count1, concurrency)

		count2, err := db2.Model(table2).Count()
		t.AssertNil(err)
		t.Assert(count2, concurrency)
	})
}

// Test_Concurrent_Model_Clone tests concurrent model cloning
func Test_Concurrent_Model_Clone(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		baseModel := db.Model(table).Where("id>", 0)
		var wg sync.WaitGroup
		concurrency := 20

		wg.Add(concurrency)
		for i := 0; i < concurrency; i++ {
			go func(id int) {
				defer wg.Done()
				// Clone model for each goroutine
				m := baseModel.Clone()
				result, err := m.Where("id<=", TableSize/2).All()
				t.AssertNil(err)
				t.AssertGT(len(result), 0)
			}(i)
		}
		wg.Wait()
	})
}

// Test_Concurrent_Batch_Insert tests concurrent batch insert operations
func Test_Concurrent_Batch_Insert(t *testing.T) {
	table := createTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		var wg sync.WaitGroup
		concurrency := 5
		batchSize := 10

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
				_, err := db.Model(table).Data(batch).Insert()
				t.AssertNil(err)
			}(i)
		}
		wg.Wait()

		// Verify all batch inserts
		count, err := db.Model(table).Count()
		t.AssertNil(err)
		t.Assert(count, concurrency*batchSize)
	})
}
