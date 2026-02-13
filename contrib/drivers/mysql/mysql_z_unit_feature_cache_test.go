// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package mysql_test

import (
	"context"
	"testing"
	"time"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/test/gtest"
)

// Test_Model_Cache_Basic tests basic cache functionality
func Test_Model_Cache_Basic(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		// First query - cache miss, result from DB
		one, err := db.Model(table).Cache(gdb.CacheOption{
			Duration: time.Second * 10,
			Name:     "test_cache_basic",
		}).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["id"], 1)
		t.Assert(one["passport"], "user_1")

		// Update the record in DB
		_, err = db.Model(table).Data(g.Map{"passport": "updated_user"}).Where("id", 1).Update()
		t.AssertNil(err)

		// Second query - cache hit, still returns old cached value
		one, err = db.Model(table).Cache(gdb.CacheOption{
			Duration: time.Second * 10,
			Name:     "test_cache_basic",
		}).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["passport"], "user_1") // cached value, not "updated_user"

		// Query without cache - returns updated value from DB
		one, err = db.Model(table).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["passport"], "updated_user")
	})
}

// Test_Model_Cache_TTL tests cache TTL expiration
func Test_Model_Cache_TTL(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		// Cache with short TTL
		one, err := db.Model(table).Cache(gdb.CacheOption{
			Duration: time.Millisecond * 100, // 100ms TTL
			Name:     "test_cache_ttl",
		}).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["passport"], "user_1")

		// Update record
		_, err = db.Model(table).Data(g.Map{"passport": "ttl_test"}).Where("id", 1).Update()
		t.AssertNil(err)

		// Immediate query - cache still valid
		one, err = db.Model(table).Cache(gdb.CacheOption{
			Duration: time.Millisecond * 100,
			Name:     "test_cache_ttl",
		}).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["passport"], "user_1") // cached value

		// Wait for cache to expire
		time.Sleep(time.Millisecond * 150)

		// Query after expiration - should get fresh data
		one, err = db.Model(table).Cache(gdb.CacheOption{
			Duration: time.Millisecond * 100,
			Name:     "test_cache_ttl",
		}).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["passport"], "ttl_test") // fresh value from DB
	})
}

// Test_Model_Cache_Clear tests clearing cache with negative duration
func Test_Model_Cache_Clear(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		// Set cache
		one, err := db.Model(table).Cache(gdb.CacheOption{
			Duration: time.Second * 60,
			Name:     "test_cache_clear",
		}).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["passport"], "user_1")

		// Update record
		_, err = db.Model(table).Data(g.Map{"passport": "cleared"}).Where("id", 1).Update()
		t.AssertNil(err)

		// Clear cache with negative duration
		_, err = db.Model(table).Cache(gdb.CacheOption{
			Duration: -1,
			Name:     "test_cache_clear",
		}).Where("id", 1).One()
		t.AssertNil(err)

		// Query again - should get fresh data since cache was cleared
		one, err = db.Model(table).Cache(gdb.CacheOption{
			Duration: time.Second * 60,
			Name:     "test_cache_clear",
		}).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["passport"], "cleared")
	})
}

// Test_Model_Cache_NoExpire tests cache with no expiration (Duration=0)
func Test_Model_Cache_NoExpire(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		// Cache with no expiration
		one, err := db.Model(table).Cache(gdb.CacheOption{
			Duration: 0, // never expires
			Name:     "test_cache_no_expire",
		}).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["passport"], "user_1")

		// Update record
		_, err = db.Model(table).Data(g.Map{"passport": "no_expire_test"}).Where("id", 1).Update()
		t.AssertNil(err)

		// Wait a bit
		time.Sleep(time.Millisecond * 100)

		// Query - cache should still be valid
		one, err = db.Model(table).Cache(gdb.CacheOption{
			Duration: 0,
			Name:     "test_cache_no_expire",
		}).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["passport"], "user_1") // cached value persists

		// Clear the cache
		_, err = db.Model(table).Cache(gdb.CacheOption{
			Duration: -1,
			Name:     "test_cache_no_expire",
		}).Where("id", 1).One()
		t.AssertNil(err)
	})
}

// Test_Model_Cache_Force tests Force option to cache nil results
func Test_Model_Cache_Force(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		// Query non-existent record with Force=true
		one, err := db.Model(table).Cache(gdb.CacheOption{
			Duration: time.Second * 10,
			Name:     "test_cache_force",
			Force:    true, // cache even nil results
		}).Where("id", 9999).One()
		t.AssertNil(err)
		t.Assert(one.IsEmpty(), true)

		// Insert the record
		_, err = db.Model(table).Data(g.Map{
			"id":       9999,
			"passport": "force_test",
		}).Insert()
		t.AssertNil(err)

		// Query again - should return cached nil result
		one, err = db.Model(table).Cache(gdb.CacheOption{
			Duration: time.Second * 10,
			Name:     "test_cache_force",
			Force:    true,
		}).Where("id", 9999).One()
		t.AssertNil(err)
		t.Assert(one.IsEmpty(), true) // cached nil result

		// Clear cache
		_, err = db.Model(table).Cache(gdb.CacheOption{
			Duration: -1,
			Name:     "test_cache_force",
		}).Where("id", 9999).One()
		t.AssertNil(err)

		// Query without cache - should return the inserted record
		one, err = db.Model(table).Where("id", 9999).One()
		t.AssertNil(err)
		t.Assert(one["passport"], "force_test")
	})
}

// Test_Model_Cache_DisabledInTransaction tests cache is disabled in transactions
func Test_Model_Cache_DisabledInTransaction(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		err := db.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
			// First query in transaction
			one, err := tx.Model(table).Cache(gdb.CacheOption{
				Duration: time.Second * 10,
				Name:     "test_tx_cache",
			}).Where("id", 1).One()
			t.AssertNil(err)
			t.Assert(one["passport"], "user_1")

			// Update in transaction
			_, err = tx.Model(table).Data(g.Map{"passport": "tx_update"}).Where("id", 1).Update()
			t.AssertNil(err)

			// Second query - should see updated value (cache disabled in tx)
			one, err = tx.Model(table).Cache(gdb.CacheOption{
				Duration: time.Second * 10,
				Name:     "test_tx_cache",
			}).Where("id", 1).One()
			t.AssertNil(err)
			t.Assert(one["passport"], "tx_update") // not cached, fresh from DB

			return nil
		})
		t.AssertNil(err)
	})
}

// Test_Model_PageCache tests pagination cache
func Test_Model_PageCache(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		// First page query with cache
		all, err := db.Model(table).PageCache(
			gdb.CacheOption{Duration: time.Second * 10, Name: "test_page_count"},
			gdb.CacheOption{Duration: time.Second * 10, Name: "test_page_data"},
		).Page(1, 3).All()
		t.AssertNil(err)
		t.Assert(len(all), 3)

		// Insert new record
		_, err = db.Model(table).Data(g.Map{
			"id":       11,
			"passport": "user_11",
		}).Insert()
		t.AssertNil(err)

		// Query again - should return cached results
		all, err = db.Model(table).PageCache(
			gdb.CacheOption{Duration: time.Second * 10, Name: "test_page_count"},
			gdb.CacheOption{Duration: time.Second * 10, Name: "test_page_data"},
		).Page(1, 3).All()
		t.AssertNil(err)
		t.Assert(len(all), 3) // cached results

		// Clear cache
		_, err = db.Model(table).Cache(gdb.CacheOption{
			Duration: -1,
			Name:     "test_page_count",
		}).Count()
		t.AssertNil(err)

		// Query without cache - should return updated count
		count, err := db.Model(table).Count()
		t.AssertNil(err)
		t.Assert(count, 11)
	})
}

// Test_Model_Cache_DifferentNames tests different cache names for same query
func Test_Model_Cache_DifferentNames(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		// Cache with name1
		one, err := db.Model(table).Cache(gdb.CacheOption{
			Duration: time.Second * 10,
			Name:     "cache_name1",
		}).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["passport"], "user_1")

		// Cache same query with name2
		one, err = db.Model(table).Cache(gdb.CacheOption{
			Duration: time.Second * 10,
			Name:     "cache_name2",
		}).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["passport"], "user_1")

		// Update record
		_, err = db.Model(table).Data(g.Map{"passport": "diff_name"}).Where("id", 1).Update()
		t.AssertNil(err)

		// Clear only cache_name1
		_, err = db.Model(table).Cache(gdb.CacheOption{
			Duration: -1,
			Name:     "cache_name1",
		}).Where("id", 1).One()
		t.AssertNil(err)

		// Query with cache_name1 - should get fresh data
		one, err = db.Model(table).Cache(gdb.CacheOption{
			Duration: time.Second * 10,
			Name:     "cache_name1",
		}).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["passport"], "diff_name")

		// Query with cache_name2 - should still have cached old value
		one, err = db.Model(table).Cache(gdb.CacheOption{
			Duration: time.Second * 10,
			Name:     "cache_name2",
		}).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["passport"], "user_1") // still cached
	})
}
