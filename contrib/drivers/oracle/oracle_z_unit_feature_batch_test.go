// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package oracle_test

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
		data := g.Slice{}
		for i := 1; i <= 10; i++ {
			data = append(data, g.Map{
				"id":       i,
				"passport": fmt.Sprintf("batch_user_%d", i),
				"password": fmt.Sprintf("batch_pass_%d", i),
				"nickname": fmt.Sprintf("batch_name_%d", i),
			})
		}

		result, err := db.Model(table).Batch(3).Data(data).Insert()
		t.AssertNil(err)
		n, _ := result.RowsAffected()
		t.Assert(n, 10)

		count, err := db.Model(table).Count()
		t.AssertNil(err)
		t.Assert(count, 10)

		one, err := db.Model(table).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["PASSPORT"], "batch_user_1")

		one, err = db.Model(table).Where("id", 10).One()
		t.AssertNil(err)
		t.Assert(one["PASSPORT"], "batch_user_10")
	})
}

// Test_Model_Batch_Replace tests batch replace operation
func Test_Model_Batch_Replace(t *testing.T) {
	table := createTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		data := g.Slice{}
		for i := 1; i <= 5; i++ {
			data = append(data, g.Map{
				"id":       i,
				"passport": fmt.Sprintf("original_%d", i),
				"password": fmt.Sprintf("original_pass_%d", i),
				"nickname": fmt.Sprintf("original_name_%d", i),
			})
		}
		_, err := db.Model(table).Data(data).Insert()
		t.AssertNil(err)

		replaceData := g.Slice{}
		for i := 3; i <= 8; i++ {
			replaceData = append(replaceData, g.Map{
				"id":       i,
				"passport": fmt.Sprintf("replaced_%d", i),
				"password": fmt.Sprintf("replaced_pass_%d", i),
				"nickname": fmt.Sprintf("new_name_%d", i),
			})
		}
		result, err := db.Model(table).Batch(2).Data(replaceData).Replace()
		t.AssertNil(err)
		n, _ := result.RowsAffected()
		t.Assert(n, 6)

		one, err := db.Model(table).Where("id", 3).One()
		t.AssertNil(err)
		t.Assert(one["PASSPORT"], "replaced_3")
		t.Assert(one["NICKNAME"], "new_name_3")

		one, err = db.Model(table).Where("id", 8).One()
		t.AssertNil(err)
		t.Assert(one["PASSPORT"], "replaced_8")

		count, err := db.Model(table).Count()
		t.AssertNil(err)
		t.Assert(count, 8)
	})
}

// Test_Model_Batch_Save tests batch save operation
func Test_Model_Batch_Save(t *testing.T) {
	table := createTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		data := g.Slice{}
		for i := 1; i <= 5; i++ {
			data = append(data, g.Map{
				"id":       i,
				"passport": fmt.Sprintf("save_user_%d", i),
				"password": fmt.Sprintf("save_pass_%d", i),
				"nickname": fmt.Sprintf("save_name_%d", i),
			})
		}
		_, err := db.Model(table).Data(data).Insert()
		t.AssertNil(err)

		saveData := g.Slice{}
		for i := 3; i <= 8; i++ {
			saveData = append(saveData, g.Map{
				"id":       i,
				"passport": fmt.Sprintf("saved_%d", i),
				"password": fmt.Sprintf("saved_pass_%d", i),
				"nickname": fmt.Sprintf("save_name_%d", i),
			})
		}
		result, err := db.Model(table).Batch(3).Data(saveData).Save()
		t.AssertNil(err)
		n, _ := result.RowsAffected()
		t.Assert(n, 6)

		one, err := db.Model(table).Where("id", 3).One()
		t.AssertNil(err)
		t.Assert(one["PASSPORT"], "saved_3")

		one, err = db.Model(table).Where("id", 8).One()
		t.AssertNil(err)
		t.Assert(one["PASSPORT"], "saved_8")

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
		data := g.Slice{}
		totalRecords := 250
		for i := 1; i <= totalRecords; i++ {
			data = append(data, g.Map{
				"id":       i,
				"passport": fmt.Sprintf("large_user_%d", i),
				"password": fmt.Sprintf("large_pass_%d", i),
				"nickname": fmt.Sprintf("large_name_%d", i),
			})
		}

		result, err := db.Model(table).Batch(100).Data(data).Insert()
		t.AssertNil(err)
		n, _ := result.RowsAffected()
		t.Assert(n, totalRecords)

		count, err := db.Model(table).Count()
		t.AssertNil(err)
		t.Assert(count, totalRecords)

		one, err := db.Model(table).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["PASSPORT"], "large_user_1")

		one, err = db.Model(table).Where("id", totalRecords).One()
		t.AssertNil(err)
		t.Assert(one["PASSPORT"], fmt.Sprintf("large_user_%d", totalRecords))
	})
}

// Test_Model_Batch_EmptyBatch tests batch operation with empty data
func Test_Model_Batch_EmptyBatch(t *testing.T) {
	table := createTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		data := g.Slice{}

		_, err := db.Model(table).Batch(10).Data(data).Insert()
		t.AssertNE(err, nil)
		t.AssertIN(err.Error(), "data list cannot be empty")

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
		data := g.Slice{
			g.Map{
				"id":       1,
				"passport": "single_user",
				"password": "single_pass",
				"nickname": "single_name",
			},
		}

		result, err := db.Model(table).Batch(10).Data(data).Insert()
		t.AssertNil(err)
		n, _ := result.RowsAffected()
		t.Assert(n, 1)

		one, err := db.Model(table).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["PASSPORT"], "single_user")
		t.Assert(one["NICKNAME"], "single_name")
	})
}

// Test_Model_Batch_VsBatch tests performance comparison between different batch sizes
func Test_Model_Batch_VsBatch(t *testing.T) {
	table := createTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		data := g.Slice{}
		for i := 1; i <= 100; i++ {
			data = append(data, g.Map{
				"id":       i,
				"passport": fmt.Sprintf("perf_user_%d", i),
				"password": fmt.Sprintf("perf_pass_%d", i),
				"nickname": fmt.Sprintf("perf_name_%d", i),
			})
		}

		result, err := db.Model(table).Batch(1).Data(data).Insert()
		t.AssertNil(err)
		n, _ := result.RowsAffected()
		t.Assert(n, 100)

		_, err = db.Model(table).Where("1=1").Delete()
		t.AssertNil(err)

		result, err = db.Model(table).Batch(10).Data(data).Insert()
		t.AssertNil(err)
		n, _ = result.RowsAffected()
		t.Assert(n, 100)

		_, err = db.Model(table).Where("1=1").Delete()
		t.AssertNil(err)

		result, err = db.Model(table).Batch(50).Data(data).Insert()
		t.AssertNil(err)
		n, _ = result.RowsAffected()
		t.Assert(n, 100)

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
				"password": fmt.Sprintf("tx_pass_%d", i),
				"nickname": fmt.Sprintf("tx_name_%d", i),
			})
		}

		err := db.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
			result, err := tx.Model(table).Batch(10).Data(data).Insert()
			t.AssertNil(err)
			n, _ := result.RowsAffected()
			t.Assert(n, 50)
			return nil
		})
		t.AssertNil(err)

		count, err := db.Model(table).Count()
		t.AssertNil(err)
		t.Assert(count, 50)

		_, err = db.Model(table).Where("1=1").Delete()
		t.AssertNil(err)

		err = db.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
			_, err := tx.Model(table).Batch(10).Data(data).Insert()
			t.AssertNil(err)
			return fmt.Errorf("rollback test")
		})
		t.AssertNE(err, nil)
		t.Assert(err.Error(), "rollback test")

		count, err = db.Model(table).Count()
		t.AssertNil(err)
		t.Assert(count, 0)
	})
}

// Test_Model_Batch_Save_DifferentFields tests batch save and insert ignore with records carrying different fields
func Test_Model_Batch_Save_DifferentFields(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Data(g.List{
			g.Map{"id": 1, "nickname": "saved_1"},
			g.Map{"id": 11, "passport": "user_11", "password": "pass_11", "nickname": "name_11"},
		}).Save()
		t.AssertNil(err)
		n, _ := result.RowsAffected()
		t.Assert(n, 2)

		one, err := db.Model(table).WherePri(1).One()
		t.AssertNil(err)
		t.Assert(one["PASSPORT"], "user_1")
		t.Assert(one["NICKNAME"], "saved_1")

		one, err = db.Model(table).WherePri(11).One()
		t.AssertNil(err)
		t.Assert(one["PASSPORT"], "user_11")
		t.Assert(one["NICKNAME"], "name_11")
	})

	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Data(g.List{
			g.Map{"id": 2, "passport": "ignored_2", "password": "pass_2", "nickname": "ignored_2"},
			g.Map{"id": 12, "passport": "user_12", "password": "pass_12", "nickname": "name_12"},
		}).InsertIgnore()
		t.AssertNil(err)
		n, _ := result.RowsAffected()
		t.Assert(n, 1)

		one, err := db.Model(table).WherePri(2).One()
		t.AssertNil(err)
		t.Assert(one["PASSPORT"], "user_2")

		one, err = db.Model(table).WherePri(12).One()
		t.AssertNil(err)
		t.Assert(one["PASSPORT"], "user_12")

		count, err := db.Model(table).Count()
		t.AssertNil(err)
		t.Assert(count, TableSize+2)
	})
}
