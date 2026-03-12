// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package mariadb_test

import (
	"context"
	"testing"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/test/gtest"
)

// Test_Model_Insert_NilData tests Insert with nil data
func Test_Model_Insert_NilData(t *testing.T) {
	table := createTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(table).Data(nil).Insert()
		t.AssertNE(err, nil)
	})
}

// Test_Model_Insert_EmptyMap tests Insert with empty map
func Test_Model_Insert_EmptyMap(t *testing.T) {
	table := createTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(table).Data(g.Map{}).Insert()
		t.AssertNE(err, nil)
	})
}

// Test_Model_Insert_EmptySlice tests Insert with empty slice
func Test_Model_Insert_EmptySlice(t *testing.T) {
	table := createTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(table).Data(g.Slice{}).Insert()
		t.AssertNE(err, nil)
	})
}

// Test_Model_Update_NilData tests Update with nil data
func Test_Model_Update_NilData(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(table).Data(nil).Where("id", 1).Update()
		t.AssertNE(err, nil)
	})
}

// Test_Model_Update_EmptyData tests Update with empty data
func Test_Model_Update_EmptyData(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(table).Data(g.Map{}).Where("id", 1).Update()
		t.AssertNE(err, nil)
	})
}

// Test_Model_Update_NoWhere tests Update without WHERE clause is rejected by framework
func Test_Model_Update_NoWhere(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		// Framework safety check: Update without WHERE should return error
		_, err := db.Model(table).Data(g.Map{"nickname": "updated"}).Update()
		t.AssertNE(err, nil)
	})
}

// Test_Model_Delete_NoWhere tests Delete without WHERE clause is rejected by framework
func Test_Model_Delete_NoWhere(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		// Framework safety check: Delete without WHERE should return error
		_, err := db.Model(table).Delete()
		t.AssertNE(err, nil)
	})
}

// Test_Model_Scan_NilPointer tests Scan with nil pointer
func Test_Model_Scan_NilPointer(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		err := db.Model(table).Where("id", 1).Scan(nil)
		t.AssertNE(err, nil)
	})
}

// Test_Model_Scan_InvalidPointer tests Scan with invalid pointer type
func Test_Model_Scan_InvalidPointer(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		var str string
		err := db.Model(table).Where("id", 1).Scan(&str)
		t.AssertNE(err, nil)
	})
}

// Test_Model_Scan_EmptyResult tests Scan with empty result
func Test_Model_Scan_EmptyResult(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	type User struct {
		Id int
	}

	// Scan initialized struct with empty result returns sql.ErrNoRows
	gtest.C(t, func(t *gtest.T) {
		var user User
		err := db.Model(table).Where("id > ?", 1000).Scan(&user)
		t.AssertNE(err, nil)
	})

	// Scan nil pointer with empty result returns nil error
	gtest.C(t, func(t *gtest.T) {
		var user *User
		err := db.Model(table).Where("id > ?", 1000).Scan(&user)
		t.AssertNil(err)
		t.Assert(user, nil)
	})
}

// Test_Model_Where_InvalidOperator tests Where with invalid operator
func Test_Model_Where_InvalidOperator(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		// Invalid SQL should cause error at query time
		_, err := db.Model(table).Where("id INVALID_OP ?", 1).All()
		t.AssertNE(err, nil)
	})
}

// Test_Model_Where_EmptyString tests Where with empty string
func Test_Model_Where_EmptyString(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Where("").All()
		t.AssertNil(err)
		t.Assert(len(result), TableSize) // Empty WHERE returns all
	})
}

// Test_Model_Fields_InvalidField tests Fields with non-existent field
func Test_Model_Fields_InvalidField(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(table).Fields("non_existent_field").All()
		t.AssertNE(err, nil)
	})
}

// Test_Model_Fields_Empty tests Fields with empty string
// Regression test for #4697: Fields("") should handle empty string gracefully
// https://github.com/gogf/gf/issues/4697
func Test_Model_Fields_Empty(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Fields("").Limit(1).All()
		t.AssertNil(err)
		t.AssertLE(len(result), 1)
	})
}

// Test_Model_Order_InvalidSyntax tests Order with invalid syntax
func Test_Model_Order_InvalidSyntax(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		// Invalid ORDER BY syntax
		_, err := db.Model(table).Order("id INVALID").All()
		t.AssertNE(err, nil)
	})
}

// Test_Model_Group_UnknownColumn tests Group with non-existent column
func Test_Model_Group_UnknownColumn(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(table).Group("non_existent_field").All()
		t.AssertNE(err, nil)
	})
}

// Test_Model_TableNotExist tests querying non-existent table
func Test_Model_TableNotExist(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model("non_existent_table_xyz").All()
		t.AssertNE(err, nil)
	})
}

// Test_Model_InvalidTableName tests invalid table name
func Test_Model_InvalidTableName(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// Empty table name
		_, err := db.Model("").All()
		t.AssertNE(err, nil)
	})
}

// Test_Model_SQLInjection_Where tests SQL injection prevention in Where
func Test_Model_SQLInjection_Where(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		// Attempt SQL injection through string column parameter.
		// Using string column `nickname` instead of int column `id`,
		// because MySQL coerces "1 OR 1=1" to 1 for int columns.
		maliciousInput := "1 OR 1=1"
		result, err := db.Model(table).Where("nickname = ?", maliciousInput).All()
		t.AssertNil(err)
		t.Assert(len(result), 0) // Should not return all records
	})

	gtest.C(t, func(t *gtest.T) {
		// Attempt SQL injection with quotes, using string column to avoid
		// MySQL implicit int conversion (which would coerce "1'..." to 1)
		maliciousInput := "1'; DROP TABLE " + table + "; --"
		result, err := db.Model(table).Where("nickname = ?", maliciousInput).All()
		t.AssertNil(err)
		t.Assert(len(result), 0)
		// Table should still exist
		count, err := db.Model(table).Count()
		t.AssertNil(err)
		t.Assert(count, TableSize)
	})
}

// Test_Model_SQLInjection_Insert tests SQL injection prevention in Insert
func Test_Model_SQLInjection_Insert(t *testing.T) {
	table := createTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		maliciousData := g.Map{
			"id":       1,
			"passport": "'; DROP TABLE " + table + "; --",
			"password": "pwd",
			"nickname": "test",
		}
		_, err := db.Model(table).Data(maliciousData).Insert()
		t.AssertNil(err)

		// Verify data was inserted correctly and table still exists
		one, err := db.Model(table).Where("id", 1).One()
		t.AssertNil(err)
		t.AssertNE(one, nil)
		t.Assert(one["passport"].String(), "'; DROP TABLE "+table+"; --")
	})
}

// Test_Model_SQLInjection_Update tests SQL injection prevention in Update
func Test_Model_SQLInjection_Update(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		// Use shorter malicious string to fit in nickname column
		maliciousData := g.Map{
			"nickname": "'; DELETE FROM users; --",
		}
		_, err := db.Model(table).Data(maliciousData).Where("id", 1).Update()
		t.AssertNil(err)

		// Verify only one record was updated (parameterized query prevents injection)
		one, err := db.Model(table).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["nickname"].String(), "'; DELETE FROM users; --")

		// Other records should still exist (injection was prevented)
		count, err := db.Model(table).Count()
		t.AssertNil(err)
		t.Assert(count, TableSize)
	})
}

// Test_Model_Context_Cancelled tests query with cancelled context
func Test_Model_Context_Cancelled(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel() // Cancel immediately

		_, err := db.Model(table).Ctx(ctx).All()
		t.AssertNE(err, nil)
		t.Assert(gerror.Is(err, context.Canceled), true)
	})
}

// Test_Model_Value_EmptyResult tests Value with empty result
func Test_Model_Value_EmptyResult(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		value, err := db.Model(table).Where("id > ?", 1000).Value()
		t.AssertNil(err)
		t.Assert(value.IsEmpty(), true)
	})
}

// Test_Model_Array_EmptyResult tests Array with empty result
func Test_Model_Array_EmptyResult(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		array, err := db.Model(table).Where("id > ?", 1000).Array()
		t.AssertNil(err)
		t.Assert(len(array), 0)
	})
}

// Test_Model_Count_InvalidTable tests Count on invalid table
func Test_Model_Count_InvalidTable(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model("non_existent_table").Count()
		t.AssertNE(err, nil)
	})
}

// Test_Model_Max_EmptyResult tests Max with empty result
func Test_Model_Max_EmptyResult(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		value, err := db.Model(table).Where("id > ?", 1000).Max("id")
		t.AssertNil(err)
		t.Assert(value, 0) // Returns 0 for empty result
	})
}

// Test_Model_Min_EmptyResult tests Min with empty result
func Test_Model_Min_EmptyResult(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		value, err := db.Model(table).Where("id > ?", 1000).Min("id")
		t.AssertNil(err)
		t.Assert(value, 0) // Returns 0 for empty result
	})
}

// Test_Model_Avg_EmptyResult tests Avg with empty result
func Test_Model_Avg_EmptyResult(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		value, err := db.Model(table).Where("id > ?", 1000).Avg("id")
		t.AssertNil(err)
		t.Assert(value, 0) // Returns 0 for empty result
	})
}

// Test_Model_Sum_EmptyResult tests Sum with empty result
func Test_Model_Sum_EmptyResult(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		value, err := db.Model(table).Where("id > ?", 1000).Sum("id")
		t.AssertNil(err)
		t.Assert(value, 0) // Returns 0 for empty result
	})
}

// Test_Model_One_NilResult tests One returning nil
func Test_Model_One_NilResult(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		one, err := db.Model(table).Where("id > ?", 1000).One()
		t.AssertNil(err)
		t.Assert(one, nil)
	})
}

// Test_TX_Rollback_AfterError tests transaction rollback after error
func Test_TX_Rollback_AfterError(t *testing.T) {
	table := createTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		err := db.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
			// Insert valid record
			_, err := tx.Model(table).Data(g.Map{
				"id":       1,
				"passport": "pass1",
				"password": "pwd1",
				"nickname": "name1",
			}).Insert()
			if err != nil {
				return err
			}

			// Insert duplicate id (should fail)
			_, err = tx.Model(table).Data(g.Map{
				"id":       1, // Duplicate
				"passport": "pass2",
				"password": "pwd2",
				"nickname": "name2",
			}).Insert()

			return err // Return error to trigger rollback
		})

		t.AssertNE(err, nil)

		// Verify rollback - table should be empty
		count, err := db.Model(table).Count()
		t.AssertNil(err)
		t.Assert(count, 0)
	})
}

// Test_Model_Insert_DuplicateKey tests handling of duplicate key error
func Test_Model_Insert_DuplicateKey(t *testing.T) {
	table := createTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		data := g.Map{
			"id":       1,
			"passport": "pass",
			"password": "pwd",
			"nickname": "name",
		}

		// First insert should succeed
		_, err := db.Model(table).Data(data).Insert()
		t.AssertNil(err)

		// Second insert with same id should fail
		_, err = db.Model(table).Data(data).Insert()
		t.AssertNE(err, nil)
	})
}

// Test_Model_All_InvalidConnection tests query with invalid connection
func Test_Model_All_InvalidConnection(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		if dbInvalid == nil {
			t.Skip("dbInvalid not configured")
		}
		_, err := dbInvalid.Model("test_table").All()
		t.AssertNE(err, nil)
	})
}
