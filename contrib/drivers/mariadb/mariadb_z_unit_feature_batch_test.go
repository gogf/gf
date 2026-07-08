// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package mariadb_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/test/gtest"
)

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
		result, err := db.Model(table).Batch(3).Data(data).Insert()
		t.AssertNil(err)
		n, _ := result.RowsAffected()
		t.Assert(n, 10)

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
	table := createTable()
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
		result, err := db.Model(table).Batch(2).Data(replaceData).Replace()
		t.AssertNil(err)
		n, _ := result.RowsAffected()
		t.AssertGT(n, 0)

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
		count, err := db.Model(table).Count()
		t.AssertNil(err)
		t.Assert(count, 8) // ids 1-8
	})
}

// Test_Model_Batch_Save tests batch save operation
func Test_Model_Batch_Save(t *testing.T) {
	table := createTable()
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
		result, err := db.Model(table).Batch(3).Data(saveData).Save()
		t.AssertNil(err)
		n, _ := result.RowsAffected()
		t.AssertGT(n, 0)

		// Verify updated records
		one, err := db.Model(table).Where("id", 3).One()
		t.AssertNil(err)
		t.Assert(one["passport"], "saved_3")

		// Verify inserted records
		one, err = db.Model(table).Where("id", 8).One()
		t.AssertNil(err)
		t.Assert(one["passport"], "saved_8")

		// Verify total count
		count, err := db.Model(table).Count()
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
		result, err := db.Model(table).Batch(100).Data(data).Insert()
		t.AssertNil(err)
		n, _ := result.RowsAffected()
		t.Assert(n, totalRecords)

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

		result, err := db.Model(table).Batch(10).Data(data).Insert()
		t.AssertNil(err)
		n, _ := result.RowsAffected()
		t.Assert(n, 1)

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
		result, err := db.Model(table).Batch(1).Data(data).Insert()
		t.AssertNil(err)
		n, _ := result.RowsAffected()
		t.Assert(n, 100)

		// Clean up
		_, err = db.Model(table).Where("1=1").Delete()
		t.AssertNil(err)

		// Test with batch size 10
		result, err = db.Model(table).Batch(10).Data(data).Insert()
		t.AssertNil(err)
		n, _ = result.RowsAffected()
		t.Assert(n, 100)

		// Clean up
		_, err = db.Model(table).Where("1=1").Delete()
		t.AssertNil(err)

		// Test with batch size 50
		result, err = db.Model(table).Batch(50).Data(data).Insert()
		t.AssertNil(err)
		n, _ = result.RowsAffected()
		t.Assert(n, 100)

		// All batch sizes should produce same result
		count, err := db.Model(table).Count()
		t.AssertNil(err)
		t.Assert(count, 100)
	})
}

// Test_Model_Batch_WithTransaction tests batch operation within transaction
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

		// Test commit
		err := db.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
			result, err := tx.Model(table).Batch(10).Data(data).Insert()
			t.AssertNil(err)
			n, _ := result.RowsAffected()
			t.Assert(n, 50)
			return nil
		})
		t.AssertNil(err)

		// Verify commit
		count, err := db.Model(table).Count()
		t.AssertNil(err)
		t.Assert(count, 50)

		// Clean up
		_, err = db.Model(table).Where("1=1").Delete()
		t.AssertNil(err)

		// Test rollback
		err = db.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
			_, err := tx.Model(table).Batch(10).Data(data).Insert()
			t.AssertNil(err)
			return fmt.Errorf("rollback test")
		})
		t.AssertNE(err, nil)

		// Verify rollback - no records should exist
		count, err = db.Model(table).Count()
		t.AssertNil(err)
		t.Assert(count, 0)
	})
}
