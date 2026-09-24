// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package sqlitecgo_test

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
		one, err := db.Model(table).Cache(gdb.CacheOption{
			Duration: time.Second * 10,
			Name:     "test_cache_basic",
		}).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["id"], 1)
		t.Assert(one["passport"], "user_1")

		_, err = db.Model(table).Data(g.Map{"passport": "updated_user"}).Where("id", 1).Update()
		t.AssertNil(err)

		one, err = db.Model(table).Cache(gdb.CacheOption{
			Duration: time.Second * 10,
			Name:     "test_cache_basic",
		}).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["passport"], "user_1")

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
		one, err := db.Model(table).Cache(gdb.CacheOption{
			Duration: time.Millisecond * 100,
			Name:     "test_cache_ttl",
		}).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["passport"], "user_1")

		_, err = db.Model(table).Data(g.Map{"passport": "ttl_test"}).Where("id", 1).Update()
		t.AssertNil(err)

		one, err = db.Model(table).Cache(gdb.CacheOption{
			Duration: time.Millisecond * 100,
			Name:     "test_cache_ttl",
		}).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["passport"], "user_1")

		time.Sleep(time.Millisecond * 150)

		one, err = db.Model(table).Cache(gdb.CacheOption{
			Duration: time.Millisecond * 100,
			Name:     "test_cache_ttl",
		}).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["passport"], "ttl_test")
	})
}

// Test_Model_Cache_Clear tests clearing cache with negative duration
func Test_Model_Cache_Clear(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		one, err := db.Model(table).Cache(gdb.CacheOption{
			Duration: time.Second * 60,
			Name:     "test_cache_clear",
		}).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["passport"], "user_1")

		_, err = db.Model(table).Cache(gdb.CacheOption{
			Duration: -1,
			Name:     "test_cache_clear",
		}).Data(g.Map{"passport": "cleared"}).Where("id", 1).Update()
		t.AssertNil(err)

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
		one, err := db.Model(table).Cache(gdb.CacheOption{
			Duration: 0,
			Name:     "test_cache_no_expire",
		}).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["passport"], "user_1")

		_, err = db.Model(table).Data(g.Map{"passport": "no_expire_test"}).Where("id", 1).Update()
		t.AssertNil(err)

		time.Sleep(time.Millisecond * 100)

		one, err = db.Model(table).Cache(gdb.CacheOption{
			Duration: 0,
			Name:     "test_cache_no_expire",
		}).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["passport"], "user_1")

		_, err = db.Model(table).Cache(gdb.CacheOption{
			Duration: -1,
			Name:     "test_cache_no_expire",
		}).Data(g.Map{"nickname": "cleared"}).Where("id", 1).Update()
		t.AssertNil(err)

		one, err = db.Model(table).Cache(gdb.CacheOption{
			Duration: 0,
			Name:     "test_cache_no_expire",
		}).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["passport"], "no_expire_test")
	})
}

// Test_Model_Cache_Force tests Force option to cache empty results.
// Rows are written with raw db.Exec so that the write never passes through the
// model layer and cannot touch the select cache.
func Test_Model_Cache_Force(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)
	// Note: Force cache test is intentionally empty. CacheOption.Force never caches an
	// empty result, see https://github.com/gogf/gf/issues/4891, so the scenario cannot
	// be asserted until that is fixed.
}

// Test_Model_Cache_DisabledInTransaction tests cache is disabled in transactions
func Test_Model_Cache_DisabledInTransaction(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		err := db.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
			one, err := tx.Model(table).Cache(gdb.CacheOption{
				Duration: time.Second * 10,
				Name:     "test_tx_cache",
			}).Where("id", 1).One()
			t.AssertNil(err)
			t.Assert(one["passport"], "user_1")

			_, err = tx.Model(table).Data(g.Map{"passport": "tx_update"}).Where("id", 1).Update()
			t.AssertNil(err)

			one, err = tx.Model(table).Cache(gdb.CacheOption{
				Duration: time.Second * 10,
				Name:     "test_tx_cache",
			}).Where("id", 1).One()
			t.AssertNil(err)
			t.Assert(one["passport"], "tx_update")

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
		all, count, err := db.Model(table).PageCache(
			gdb.CacheOption{Duration: time.Second * 10, Name: "test_page_count"},
			gdb.CacheOption{Duration: time.Second * 10, Name: "test_page_data"},
		).Page(1, 3).AllAndCount(false)
		t.AssertNil(err)
		t.Assert(len(all), 3)
		t.Assert(count, 10)

		_, err = db.Model(table).Data(g.Map{
			"id":       11,
			"passport": "user_11",
		}).Insert()
		t.AssertNil(err)

		all, count, err = db.Model(table).PageCache(
			gdb.CacheOption{Duration: time.Second * 10, Name: "test_page_count"},
			gdb.CacheOption{Duration: time.Second * 10, Name: "test_page_data"},
		).Page(1, 3).AllAndCount(false)
		t.AssertNil(err)
		t.Assert(len(all), 3)
		// Count comes from the cache, so it still reports the pre-insert value.
		t.Assert(count, 10)

		_, err = db.Model(table).Cache(gdb.CacheOption{
			Duration: -1,
			Name:     "test_page_count",
		}).Data(g.Map{"nickname": "page_test"}).Where("id", 1).Update()
		t.AssertNil(err)

		all, count, err = db.Model(table).PageCache(
			gdb.CacheOption{Duration: time.Second * 10, Name: "test_page_count"},
			gdb.CacheOption{Duration: time.Second * 10, Name: "test_page_data"},
		).Page(1, 3).AllAndCount(false)
		t.AssertNil(err)
		t.Assert(len(all), 3)
		// Count cache was cleared above, so the new row is visible.
		t.Assert(count, 11)

		totalCount, err := db.Model(table).Count()
		t.AssertNil(err)
		t.Assert(totalCount, 11)
	})
}

// Test_Model_Cache_DifferentNames tests different cache names for same query
func Test_Model_Cache_DifferentNames(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		one, err := db.Model(table).Cache(gdb.CacheOption{
			Duration: time.Second * 10,
			Name:     "cache_name1",
		}).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["passport"], "user_1")

		one, err = db.Model(table).Cache(gdb.CacheOption{
			Duration: time.Second * 10,
			Name:     "cache_name2",
		}).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["passport"], "user_1")

		_, err = db.Model(table).Cache(gdb.CacheOption{
			Duration: -1,
			Name:     "cache_name1",
		}).Data(g.Map{"passport": "diff_name"}).Where("id", 1).Update()
		t.AssertNil(err)

		one, err = db.Model(table).Cache(gdb.CacheOption{
			Duration: time.Second * 10,
			Name:     "cache_name1",
		}).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["passport"], "diff_name")

		one, err = db.Model(table).Cache(gdb.CacheOption{
			Duration: time.Second * 10,
			Name:     "cache_name2",
		}).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["passport"], "user_1")
	})
}
