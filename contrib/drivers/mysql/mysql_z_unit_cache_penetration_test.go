// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package mysql_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/guid"
)

// Test_Model_Cache_EmptyResult_ForceTrue tests that empty results are cached when Force=true
func Test_Model_Cache_EmptyResult_ForceTrue(t *testing.T) {
	table := createTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		// First query: should query database and cache empty result
		result, err := db.Model(table).Where("id", 999).Cache(gdb.CacheOption{
			Duration: time.Second * 5,
			Name:     "test_empty_force_true",
			Force:    true,
		}).All()
		t.AssertNil(err)
		t.Assert(len(result), 0)

		// Insert data after caching empty result
		_, err = db.Model(table).Data(g.MapStrAny{
			"id":       999,
			"passport": "passport_999",
			"password": "password_999",
			"nickname": "nickname_999",
		}).Insert()
		t.AssertNil(err)

		// Second query: should return cached empty result (not the newly inserted data)
		result, err = db.Model(table).Where("id", 999).Cache(gdb.CacheOption{
			Duration: time.Second * 5,
			Name:     "test_empty_force_true",
			Force:    true,
		}).All()
		t.AssertNil(err)
		t.Assert(len(result), 0) // Still returns empty because it's cached

		// Wait for cache to expire
		time.Sleep(time.Second * 6)

		// Third query: cache expired, should return new data
		result, err = db.Model(table).Where("id", 999).Cache(gdb.CacheOption{
			Duration: time.Second * 5,
			Name:     "test_empty_force_true_2",
			Force:    true,
		}).All()
		t.AssertNil(err)
		t.Assert(len(result), 1)
		t.Assert(result[0]["id"], 999)
	})
}

// Test_Model_Cache_EmptyResult_ForceFalse tests that empty results are NOT cached when Force=false
func Test_Model_Cache_EmptyResult_ForceFalse(t *testing.T) {
	table := createTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		// First query: should query database, get empty result but NOT cache it
		result, err := db.Model(table).Where("id", 888).Cache(gdb.CacheOption{
			Duration: time.Second * 5,
			Name:     "test_empty_force_false",
			Force:    false,
		}).All()
		t.AssertNil(err)
		t.Assert(len(result), 0)

		// Insert data after querying empty result
		_, err = db.Model(table).Data(g.MapStrAny{
			"id":       888,
			"passport": "passport_888",
			"password": "password_888",
			"nickname": "nickname_888",
		}).Insert()
		t.AssertNil(err)

		// Second query: should return new data (because empty result was not cached)
		result, err = db.Model(table).Where("id", 888).Cache(gdb.CacheOption{
			Duration: time.Second * 5,
			Name:     "test_empty_force_false",
			Force:    false,
		}).All()
		t.AssertNil(err)
		t.Assert(len(result), 1)
		t.Assert(result[0]["id"], 888)
	})
}

// Test_Model_Cache_Count_Zero_ForceTrue tests Count=0 caching with Force=true
func Test_Model_Cache_Count_Zero_ForceTrue(t *testing.T) {
	table := createTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		// First query: Count should be 0 and cached
		count, err := db.Model(table).Where("id", 777).Cache(gdb.CacheOption{
			Duration: time.Second * 5,
			Name:     "test_count_zero_force_true",
			Force:    true,
		}).Count()
		t.AssertNil(err)
		t.Assert(count, int64(0))

		// Insert data
		_, err = db.Model(table).Data(g.MapStrAny{
			"id":       777,
			"passport": "passport_777",
			"password": "password_777",
			"nickname": "nickname_777",
		}).Insert()
		t.AssertNil(err)

		// Second query: should still return 0 from cache
		count, err = db.Model(table).Where("id", 777).Cache(gdb.CacheOption{
			Duration: time.Second * 5,
			Name:     "test_count_zero_force_true",
			Force:    true,
		}).Count()
		t.AssertNil(err)
		t.Assert(count, int64(0)) // Still 0 because cached
	})
}

// Test_Model_Cache_Count_Zero_ForceFalse tests Count=0 caching with Force=false
func Test_Model_Cache_Count_Zero_ForceFalse(t *testing.T) {
	table := createTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		// First query: Count should be 0 but not cached
		count, err := db.Model(table).Where("id", 666).Cache(gdb.CacheOption{
			Duration: time.Second * 5,
			Name:     "test_count_zero_force_false",
			Force:    false,
		}).Count()
		t.AssertNil(err)
		t.Assert(count, int64(0))

		// Insert data
		_, err = db.Model(table).Data(g.MapStrAny{
			"id":       666,
			"passport": "passport_666",
			"password": "password_666",
			"nickname": "nickname_666",
		}).Insert()
		t.AssertNil(err)

		// Second query: should return new count
		count, err = db.Model(table).Where("id", 666).Cache(gdb.CacheOption{
			Duration: time.Second * 5,
			Name:     "test_count_zero_force_false",
			Force:    false,
		}).Count()
		t.AssertNil(err)
		t.Assert(count, int64(1)) // Returns 1 because empty result was not cached
	})
}

// Test_Model_Cache_Value_Empty_ForceTrue tests Value="" caching with Force=true
func Test_Model_Cache_Value_Empty_ForceTrue(t *testing.T) {
	table := createTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		// First query: Value should be empty and cached
		value, err := db.Model(table).Where("id", 555).Cache(gdb.CacheOption{
			Duration: time.Second * 5,
			Name:     "test_value_empty_force_true",
			Force:    true,
		}).Value("passport")
		t.AssertNil(err)
		t.Assert(value.Int(), 0)

		// Insert data
		_, err = db.Model(table).Data(g.MapStrAny{
			"id":       555,
			"passport": "passport_555",
			"password": "password_555",
			"nickname": "nickname_555",
		}).Insert()
		t.AssertNil(err)

		// Second query: should still return empty from cache
		value, err = db.Model(table).Where("id", 555).Cache(gdb.CacheOption{
			Duration: time.Second * 5,
			Name:     "test_value_empty_force_true",
			Force:    true,
		}).Value("passport")
		t.AssertNil(err)
		t.Assert(value.String(), "") // Still empty because cached
	})
}

// Test_Model_Cache_Value_Empty_ForceFalse tests Value="" caching with Force=false
func Test_Model_Cache_Value_Empty_ForceFalse(t *testing.T) {
	table := createTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		// First query: Value should be empty but not cached
		value, err := db.Model(table).Where("id", 444).Cache(gdb.CacheOption{
			Duration: time.Second * 5,
			Name:     "test_value_empty_force_false",
			Force:    false,
		}).Value("passport")
		t.AssertNil(err)
		t.Assert(value.String(), "")

		// Insert data
		_, err = db.Model(table).Data(g.MapStrAny{
			"id":       444,
			"passport": "passport_444",
			"password": "password_444",
			"nickname": "nickname_444",
		}).Insert()
		t.AssertNil(err)

		// Second query: should return new value
		value, err = db.Model(table).Where("id", 444).Cache(gdb.CacheOption{
			Duration: time.Second * 5,
			Name:     "test_value_empty_force_false",
			Force:    false,
		}).Value("passport")
		t.AssertNil(err)
		t.Assert(value.String(), "passport_444") // Returns value because empty was not cached
	})
}

// Test_Model_Cache_One_Empty_ForceTrue tests One() with empty result and Force=true
func Test_Model_Cache_One_Empty_ForceTrue(t *testing.T) {
	table := createTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		// First query: should return nil and cache it
		one, err := db.Model(table).Where("id", 333).Cache(gdb.CacheOption{
			Duration: time.Second * 5,
			Name:     "test_one_empty_force_true",
			Force:    true,
		}).One()
		t.AssertNil(err)
		t.Assert(one.IsEmpty(), true)

		// Insert data
		_, err = db.Model(table).Data(g.MapStrAny{
			"id":       333,
			"passport": "passport_333",
			"password": "password_333",
			"nickname": "nickname_333",
		}).Insert()
		t.AssertNil(err)

		// Second query: should still return empty from cache
		one, err = db.Model(table).Where("id", 333).Cache(gdb.CacheOption{
			Duration: time.Second * 5,
			Name:     "test_one_empty_force_true",
			Force:    true,
		}).One()
		t.AssertNil(err)
		t.Assert(one.IsEmpty(), true) // Still empty because cached
	})
}

// Test_Model_Cache_One_Empty_ForceFalse tests One() with empty result and Force=false
func Test_Model_Cache_One_Empty_ForceFalse(t *testing.T) {
	table := createTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		// First query: should return nil but not cache it
		one, err := db.Model(table).Where("id", 222).Cache(gdb.CacheOption{
			Duration: time.Second * 5,
			Name:     "test_one_empty_force_false",
			Force:    false,
		}).One()
		t.AssertNil(err)
		t.Assert(one.IsEmpty(), true)

		// Insert data
		_, err = db.Model(table).Data(g.MapStrAny{
			"id":       222,
			"passport": "passport_222",
			"password": "password_222",
			"nickname": "nickname_222",
		}).Insert()
		t.AssertNil(err)

		// Second query: should return new data
		one, err = db.Model(table).Where("id", 222).Cache(gdb.CacheOption{
			Duration: time.Second * 5,
			Name:     "test_one_empty_force_false",
			Force:    false,
		}).One()
		t.AssertNil(err)
		t.Assert(one.IsEmpty(), false)
		t.Assert(one["id"], 222)
	})
}

// Test_Model_Cache_NonEmpty_Result tests caching of non-empty results (should work the same way)
func Test_Model_Cache_NonEmpty_Result(t *testing.T) {
	table := createTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		// Insert data first
		_, err := db.Model(table).Data(g.MapStrAny{
			"id":       111,
			"passport": "passport_111",
			"password": "password_111",
			"nickname": "nickname_111",
		}).Insert()
		t.AssertNil(err)

		// First query with Force=false: should cache non-empty result
		result, err := db.Model(table).Where("id", 111).Cache(gdb.CacheOption{
			Duration: time.Second * 5,
			Name:     "test_nonempty_force_false",
			Force:    false,
		}).All()
		t.AssertNil(err)
		t.Assert(len(result), 1)
		t.Assert(result[0]["passport"], "passport_111")

		// Update data
		_, err = db.Model(table).Data("passport", "passport_111_updated").Where("id", 111).Update()
		t.AssertNil(err)

		// Second query: should still return cached result
		result, err = db.Model(table).Where("id", 111).Cache(gdb.CacheOption{
			Duration: time.Second * 5,
			Name:     "test_nonempty_force_false",
			Force:    false,
		}).All()
		t.AssertNil(err)
		t.Assert(len(result), 1)
		t.Assert(result[0]["passport"], "passport_111") // Still old value from cache
	})

	gtest.C(t, func(t *gtest.T) {
		// Test with Force=true as well
		result, err := db.Model(table).Where("id", 111).Cache(gdb.CacheOption{
			Duration: time.Second * 5,
			Name:     "test_nonempty_force_true",
			Force:    true,
		}).All()
		t.AssertNil(err)
		t.Assert(len(result), 1)
		t.Assert(result[0]["passport"], "passport_111_updated") // New query, gets updated value
	})
}

// Test_Model_Cache_Penetration_Prevention tests cache penetration scenario
func Test_Model_Cache_Penetration_Prevention(t *testing.T) {
	table := createTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		// Simulate cache penetration attack: multiple queries for non-existent data
		queryCount := 0
		cacheName := guid.S()

		// With Force=false: each query hits the database
		for i := 0; i < 5; i++ {
			_, err := db.Model(table).Where("id", 9999).Cache(gdb.CacheOption{
				Duration: time.Second * 5,
				Name:     cacheName + "_no_force",
				Force:    false,
			}).All()
			t.AssertNil(err)
			queryCount++
		}
		// All 5 queries should hit database because empty results are not cached

		// With Force=true: only first query hits the database
		cacheName2 := guid.S()
		for i := 0; i < 5; i++ {
			_, err := db.Model(table).Where("id", 8888).Cache(gdb.CacheOption{
				Duration: time.Second * 5,
				Name:     cacheName2 + "_with_force",
				Force:    true,
			}).All()
			t.AssertNil(err)
		}
		// Only first query should hit database; subsequent queries use cache
	})
}

// Test_Model_Cache_CountColumn_Zero_ForceTrue tests CountColumn with zero result and Force=true
func Test_Model_Cache_CountColumn_Zero_ForceTrue(t *testing.T) {
	table := createTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		// First query: CountColumn should be 0 and cached
		count, err := db.Model(table).Where("id > ?", 10000).Cache(gdb.CacheOption{
			Duration: time.Second * 5,
			Name:     "test_count_column_zero_force_true",
			Force:    true,
		}).CountColumn("id")
		t.AssertNil(err)
		t.Assert(count, int64(0))

		// Insert data
		_, err = db.Model(table).Data(g.MapStrAny{
			"id":       10001,
			"passport": "passport_10001",
			"password": "password_10001",
			"nickname": "nickname_10001",
		}).Insert()
		t.AssertNil(err)

		// Second query: should still return 0 from cache
		count, err = db.Model(table).Where("id > ?", 10000).Cache(gdb.CacheOption{
			Duration: time.Second * 5,
			Name:     "test_count_column_zero_force_true",
			Force:    true,
		}).CountColumn("id")
		t.AssertNil(err)
		t.Assert(count, int64(0)) // Still 0 because cached
	})
}

// Test_Model_Cache_Multiple_Empty_Queries tests multiple different empty queries
func Test_Model_Cache_Multiple_Empty_Queries(t *testing.T) {
	table := createTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		// Query 1: id=1000
		result1, err := db.Model(table).Where("id", 1000).Cache(gdb.CacheOption{
			Duration: time.Second * 5,
			Name:     "test_multiple_1",
			Force:    true,
		}).All()
		t.AssertNil(err)
		t.Assert(len(result1), 0)

		// Query 2: id=2000
		result2, err := db.Model(table).Where("id", 2000).Cache(gdb.CacheOption{
			Duration: time.Second * 5,
			Name:     "test_multiple_2",
			Force:    true,
		}).All()
		t.AssertNil(err)
		t.Assert(len(result2), 0)

		// Insert data for id=1000 only
		_, err = db.Model(table).Data(g.MapStrAny{
			"id":       1000,
			"passport": "passport_1000",
			"password": "password_1000",
			"nickname": "nickname_1000",
		}).Insert()
		t.AssertNil(err)

		// Query both again: id=1000 should return empty (cached), id=2000 should also return empty (cached)
		result1, err = db.Model(table).Where("id", 1000).Cache(gdb.CacheOption{
			Duration: time.Second * 5,
			Name:     "test_multiple_1",
			Force:    true,
		}).All()
		t.AssertNil(err)
		t.Assert(len(result1), 0) // Still empty from cache

		result2, err = db.Model(table).Where("id", 2000).Cache(gdb.CacheOption{
			Duration: time.Second * 5,
			Name:     "test_multiple_2",
			Force:    true,
		}).All()
		t.AssertNil(err)
		t.Assert(len(result2), 0) // Still empty from cache
	})
}

// Test_Model_Cache_Clear_Empty_Result tests clearing cached empty results
func Test_Model_Cache_Clear_Empty_Result(t *testing.T) {
	table := createTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		// Cache empty result
		result, err := db.Model(table).Where("id", 5000).Cache(gdb.CacheOption{
			Duration: time.Second * 5,
			Name:     "test_clear_empty",
			Force:    true,
		}).All()
		t.AssertNil(err)
		t.Assert(len(result), 0)

		// Insert data
		_, err = db.Model(table).Data(g.MapStrAny{
			"id":       5000,
			"passport": "passport_5000",
			"password": "password_5000",
			"nickname": "nickname_5000",
		}).Insert()
		t.AssertNil(err)

		// Clear cache using negative duration
		_, err = db.Model(table).Where("id", 5000).Cache(gdb.CacheOption{
			Duration: -1,
			Name:     "test_clear_empty",
			Force:    true,
		}).All()
		t.AssertNil(err)

		// Query again: should get new data after cache clear
		result, err = db.Model(table).Where("id", 5000).Cache(gdb.CacheOption{
			Duration: time.Second * 5,
			Name:     "test_clear_empty_new",
			Force:    true,
		}).All()
		t.AssertNil(err)
		t.Assert(len(result), 1)
		t.Assert(result[0]["id"], 5000)
	})
}

// Test_Model_Cache_Transaction_Disabled tests that cache is disabled in transactions
func Test_Model_Cache_Transaction_Disabled(t *testing.T) {
	table := createTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		err := db.Transaction(context.TODO(), func(ctx context.Context, tx gdb.TX) error {
			// Query in transaction: cache should be disabled
			result, err := tx.Model(table).Where("id", 6000).Cache(gdb.CacheOption{
				Duration: time.Second * 5,
				Name:     "test_transaction_cache",
				Force:    true,
			}).All()
			t.AssertNil(err)
			t.Assert(len(result), 0)

			// Insert data in transaction
			_, err = tx.Model(table).Data(g.MapStrAny{
				"id":       6000,
				"passport": "passport_6000",
				"password": "password_6000",
				"nickname": "nickname_6000",
			}).Insert()
			t.AssertNil(err)

			// Query again in transaction: should see new data (cache disabled)
			result, err = tx.Model(table).Where("id", 6000).Cache(gdb.CacheOption{
				Duration: time.Second * 5,
				Name:     "test_transaction_cache",
				Force:    true,
			}).All()
			t.AssertNil(err)
			t.Assert(len(result), 1)
			t.Assert(result[0]["id"], 6000)

			return nil
		})
		t.AssertNil(err)
	})
}

// Test_Model_Cache_EmptyLogic_ResultNil tests empty result detection for nil result
func Test_Model_Cache_EmptyLogic_ResultNil(t *testing.T) {
	table := createTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		// Query non-existent data: result should be empty (nil or empty array)
		// With Force=true, should cache the empty result
		result, err := db.Model(table).Where("id", 7001).Cache(gdb.CacheOption{
			Duration: time.Second * 5,
			Name:     "test_empty_nil_force_true",
			Force:    true,
		}).All()
		t.AssertNil(err)
		t.Assert(len(result), 0)

		// Insert data
		_, err = db.Model(table).Data(g.MapStrAny{
			"id":       7001,
			"passport": "passport_7001",
			"password": "password_7001",
			"nickname": "nickname_7001",
		}).Insert()
		t.AssertNil(err)

		// Query again: should return cached empty result (verifying isEmpty logic)
		result, err = db.Model(table).Where("id", 7001).Cache(gdb.CacheOption{
			Duration: time.Second * 5,
			Name:     "test_empty_nil_force_true",
			Force:    true,
		}).All()
		t.AssertNil(err)
		t.Assert(len(result), 0) // Still empty because cached
	})
}

// Test_Model_Cache_EmptyLogic_CountZero tests empty result detection for Count=0
func Test_Model_Cache_EmptyLogic_CountZero(t *testing.T) {
	table := createTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		// Count with no matching records: should be 0 (empty)
		// With Force=true, should cache the zero count
		count, err := db.Model(table).Where("id > ?", 8000).Cache(gdb.CacheOption{
			Duration: time.Second * 5,
			Name:     "test_empty_count_zero",
			Force:    true,
		}).Count()
		t.AssertNil(err)
		t.Assert(count, int64(0))

		// Insert data with id > 8000
		_, err = db.Model(table).Data(g.MapStrAny{
			"id":       8001,
			"passport": "passport_8001",
			"password": "password_8001",
			"nickname": "nickname_8001",
		}).Insert()
		t.AssertNil(err)

		// Query again: should return cached 0 (verifying Count=0 empty logic)
		count, err = db.Model(table).Where("id > ?", 8000).Cache(gdb.CacheOption{
			Duration: time.Second * 5,
			Name:     "test_empty_count_zero",
			Force:    true,
		}).Count()
		t.AssertNil(err)
		t.Assert(count, int64(0)) // Still 0 because cached
	})
}

// Test_Model_Cache_EmptyLogic_ValueEmpty tests empty result detection for Value=""
func Test_Model_Cache_EmptyLogic_ValueEmpty(t *testing.T) {
	table := createTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		// Value query with no matching records: should be empty
		// With Force=true, should cache the empty value
		value, err := db.Model(table).Where("id", 8100).Cache(gdb.CacheOption{
			Duration: time.Second * 5,
			Name:     "test_empty_value",
			Force:    true,
		}).Value("nickname")
		t.AssertNil(err)
		t.Assert(value.String(), "")

		// Insert data
		_, err = db.Model(table).Data(g.MapStrAny{
			"id":       8100,
			"passport": "passport_8100",
			"password": "password_8100",
			"nickname": "nickname_8100",
		}).Insert()
		t.AssertNil(err)

		// Query again: should return cached empty value (verifying Value="" empty logic)
		value, err = db.Model(table).Where("id", 8100).Cache(gdb.CacheOption{
			Duration: time.Second * 5,
			Name:     "test_empty_value",
			Force:    true,
		}).Value("nickname")
		t.AssertNil(err)
		t.Assert(value.String(), "") // Still empty because cached
	})
}

// Test_Model_Cache_IsCached_Flag_WithForceTrue tests IsCached flag behavior with Force=true
func Test_Model_Cache_IsCached_Flag_WithForceTrue(t *testing.T) {
	table := createTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		// First query: empty result should be cached with IsCached=true
		result1, err := db.Model(table).Where("id", 8200).Cache(gdb.CacheOption{
			Duration: time.Second * 5,
			Name:     "test_iscached_true",
			Force:    true,
		}).All()
		t.AssertNil(err)
		t.Assert(len(result1), 0)

		// Second query: should return cached result (proving IsCached=true works)
		// If IsCached flag didn't work, this would query database again
		result2, err := db.Model(table).Where("id", 8200).Cache(gdb.CacheOption{
			Duration: time.Second * 5,
			Name:     "test_iscached_true",
			Force:    true,
		}).All()
		t.AssertNil(err)
		t.Assert(len(result2), 0)

		// Insert data after cache
		_, err = db.Model(table).Data(g.MapStrAny{
			"id":       8200,
			"passport": "passport_8200",
			"password": "password_8200",
			"nickname": "nickname_8200",
		}).Insert()
		t.AssertNil(err)

		// Third query: still returns empty (IsCached=true prevents re-query)
		result3, err := db.Model(table).Where("id", 8200).Cache(gdb.CacheOption{
			Duration: time.Second * 5,
			Name:     "test_iscached_true",
			Force:    true,
		}).All()
		t.AssertNil(err)
		t.Assert(len(result3), 0) // IsCached=true flag is working
	})
}

// Test_Model_Cache_IsCached_Flag_WithForceFalse tests IsCached flag behavior with Force=false
func Test_Model_Cache_IsCached_Flag_WithForceFalse(t *testing.T) {
	table := createTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		// First query: empty result should NOT be cached (Force=false)
		result1, err := db.Model(table).Where("id", 8300).Cache(gdb.CacheOption{
			Duration: time.Second * 5,
			Name:     "test_iscached_false",
			Force:    false,
		}).All()
		t.AssertNil(err)
		t.Assert(len(result1), 0)

		// Insert data immediately
		_, err = db.Model(table).Data(g.MapStrAny{
			"id":       8300,
			"passport": "passport_8300",
			"password": "password_8300",
			"nickname": "nickname_8300",
		}).Insert()
		t.AssertNil(err)

		// Second query: should return new data (empty result was not cached)
		result2, err := db.Model(table).Where("id", 8300).Cache(gdb.CacheOption{
			Duration: time.Second * 5,
			Name:     "test_iscached_false",
			Force:    false,
		}).All()
		t.AssertNil(err)
		t.Assert(len(result2), 1)
		t.Assert(result2[0]["id"], 8300) // Got new data, proving empty was not cached
	})
}

// Test_Model_Cache_Force_And_IsCached_Combination tests the combination of Force and IsCached
func Test_Model_Cache_Force_And_IsCached_Combination(t *testing.T) {
	table := createTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		// Scenario 1: Force=true with empty result
		// Expected: isEmpty=true, Force=true, should cache with IsCached=true
		result, err := db.Model(table).Where("id", 8400).Cache(gdb.CacheOption{
			Duration: time.Second * 5,
			Name:     "test_combination_1",
			Force:    true,
		}).All()
		t.AssertNil(err)
		t.Assert(len(result), 0)

		// Verify it's cached by querying again
		_, err = db.Model(table).Data(g.MapStrAny{
			"id":       8400,
			"passport": "passport_8400",
			"password": "password_8400",
			"nickname": "nickname_8400",
		}).Insert()
		t.AssertNil(err)

		result, err = db.Model(table).Where("id", 8400).Cache(gdb.CacheOption{
			Duration: time.Second * 5,
			Name:     "test_combination_1",
			Force:    true,
		}).All()
		t.AssertNil(err)
		t.Assert(len(result), 0) // Still cached

		// Scenario 2: Force=false with empty result
		// Expected: isEmpty=true, Force=false, should NOT cache (no IsCached set)
		count, err := db.Model(table).Where("id", 8500).Cache(gdb.CacheOption{
			Duration: time.Second * 5,
			Name:     "test_combination_2",
			Force:    false,
		}).Count()
		t.AssertNil(err)
		t.Assert(count, int64(0))

		// Insert and query again
		_, err = db.Model(table).Data(g.MapStrAny{
			"id":       8500,
			"passport": "passport_8500",
			"password": "password_8500",
			"nickname": "nickname_8500",
		}).Insert()
		t.AssertNil(err)

		count, err = db.Model(table).Where("id", 8500).Cache(gdb.CacheOption{
			Duration: time.Second * 5,
			Name:     "test_combination_2",
			Force:    false,
		}).Count()
		t.AssertNil(err)
		t.Assert(count, int64(1)) // Got new data, not cached

		// Scenario 3: Force=true with non-empty result
		// Expected: isEmpty=false, should cache normally with IsCached=true
		result, err = db.Model(table).Where("id", 8500).Cache(gdb.CacheOption{
			Duration: time.Second * 5,
			Name:     "test_combination_3",
			Force:    true,
		}).All()
		t.AssertNil(err)
		t.Assert(len(result), 1)

		// Update data
		_, err = db.Model(table).Where("id", 8500).Data("nickname", "updated").Update()
		t.AssertNil(err)

		// Query again: should return cached old data
		result, err = db.Model(table).Where("id", 8500).Cache(gdb.CacheOption{
			Duration: time.Second * 5,
			Name:     "test_combination_3",
			Force:    true,
		}).All()
		t.AssertNil(err)
		t.Assert(result[0]["nickname"], "nickname_8500") // Old cached value
	})
}

// Test_Model_Cache_EmptyLogic_AllTypes tests empty detection for all query types
func Test_Model_Cache_EmptyLogic_AllTypes(t *testing.T) {
	table := createTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		// Test All() with empty result
		result, err := db.Model(table).Where("id", 9001).Cache(gdb.CacheOption{
			Duration: time.Second * 5,
			Name:     "test_all_empty",
			Force:    true,
		}).All()
		t.AssertNil(err)
		t.Assert(len(result), 0)

		// Test One() with empty result
		one, err := db.Model(table).Where("id", 9002).Cache(gdb.CacheOption{
			Duration: time.Second * 5,
			Name:     "test_one_empty",
			Force:    true,
		}).One()
		t.AssertNil(err)
		t.Assert(one.IsEmpty(), true)

		// Test Count() with 0 result
		count, err := db.Model(table).Where("id", 9003).Cache(gdb.CacheOption{
			Duration: time.Second * 5,
			Name:     "test_count_empty",
			Force:    true,
		}).Count()
		t.AssertNil(err)
		t.Assert(count, int64(0))

		// Test Value() with empty result
		value, err := db.Model(table).Where("id", 9004).Cache(gdb.CacheOption{
			Duration: time.Second * 5,
			Name:     "test_value_empty",
			Force:    true,
		}).Value("nickname")
		t.AssertNil(err)
		t.Assert(value.IsEmpty(), true)

		// Test CountColumn() with 0 result
		countCol, err := db.Model(table).Where("id", 9005).Cache(gdb.CacheOption{
			Duration: time.Second * 5,
			Name:     "test_countcol_empty",
			Force:    true,
		}).CountColumn("id")
		t.AssertNil(err)
		t.Assert(countCol, int64(0))

		// Insert data for all above queries
		for i := 9001; i <= 9005; i++ {
			_, err = db.Model(table).Data(g.MapStrAny{
				"id":       i,
				"passport": fmt.Sprintf("passport_%d", i),
				"password": fmt.Sprintf("password_%d", i),
				"nickname": fmt.Sprintf("nickname_%d", i),
			}).Insert()
			t.AssertNil(err)
		}

		// All queries should still return empty (cached)
		result, _ = db.Model(table).Where("id", 9001).Cache(gdb.CacheOption{
			Duration: time.Second * 5,
			Name:     "test_all_empty",
			Force:    true,
		}).All()
		t.Assert(len(result), 0)

		one, _ = db.Model(table).Where("id", 9002).Cache(gdb.CacheOption{
			Duration: time.Second * 5,
			Name:     "test_one_empty",
			Force:    true,
		}).One()
		t.Assert(one.IsEmpty(), true)

		count, _ = db.Model(table).Where("id", 9003).Cache(gdb.CacheOption{
			Duration: time.Second * 5,
			Name:     "test_count_empty",
			Force:    true,
		}).Count()
		t.Assert(count, int64(0))

		value, _ = db.Model(table).Where("id", 9004).Cache(gdb.CacheOption{
			Duration: time.Second * 5,
			Name:     "test_value_empty",
			Force:    true,
		}).Value("nickname")
		t.Assert(value.IsEmpty(), true)

		countCol, _ = db.Model(table).Where("id", 9005).Cache(gdb.CacheOption{
			Duration: time.Second * 5,
			Name:     "test_countcol_empty",
			Force:    true,
		}).CountColumn("id")
		t.Assert(countCol, int64(0))
	})
}
