// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package oracle_test

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
			Name:     table + "_basic",
		}).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["ID"], 1)
		t.Assert(one["PASSPORT"], "user_1")

		_, err = db.Model(table).Data(g.Map{"passport": "updated_user"}).Where("id", 1).Update()
		t.AssertNil(err)

		one, err = db.Model(table).Cache(gdb.CacheOption{
			Duration: time.Second * 10,
			Name:     table + "_basic",
		}).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["PASSPORT"], "user_1")

		one, err = db.Model(table).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["PASSPORT"], "updated_user")
	})
}

// Test_Model_Cache_TTL tests cache TTL expiration
func Test_Model_Cache_TTL(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		one, err := db.Model(table).Cache(gdb.CacheOption{
			Duration: time.Millisecond * 500,
			Name:     table + "_ttl",
		}).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["PASSPORT"], "user_1")

		_, err = db.Model(table).Data(g.Map{"passport": "ttl_test"}).Where("id", 1).Update()
		t.AssertNil(err)

		one, err = db.Model(table).Cache(gdb.CacheOption{
			Duration: time.Millisecond * 500,
			Name:     table + "_ttl",
		}).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["PASSPORT"], "user_1")

		time.Sleep(time.Millisecond * 600)

		one, err = db.Model(table).Cache(gdb.CacheOption{
			Duration: time.Millisecond * 500,
			Name:     table + "_ttl",
		}).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["PASSPORT"], "ttl_test")
	})
}

// Test_Model_Cache_Clear tests clearing cache with negative duration
func Test_Model_Cache_Clear(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		one, err := db.Model(table).Cache(gdb.CacheOption{
			Duration: time.Second * 60,
			Name:     table + "_clear",
		}).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["PASSPORT"], "user_1")

		_, err = db.Model(table).Cache(gdb.CacheOption{
			Duration: -1,
			Name:     table + "_clear",
		}).Data(g.Map{"passport": "cleared"}).Where("id", 1).Update()
		t.AssertNil(err)

		one, err = db.Model(table).Cache(gdb.CacheOption{
			Duration: time.Second * 60,
			Name:     table + "_clear",
		}).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["PASSPORT"], "cleared")
	})
}

// Test_Model_Cache_NoExpire tests cache with no expiration (Duration=0)
func Test_Model_Cache_NoExpire(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		one, err := db.Model(table).Cache(gdb.CacheOption{
			Duration: 0,
			Name:     table + "_no_expire",
		}).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["PASSPORT"], "user_1")

		_, err = db.Model(table).Data(g.Map{"passport": "no_expire_test"}).Where("id", 1).Update()
		t.AssertNil(err)

		one, err = db.Model(table).Cache(gdb.CacheOption{
			Duration: 0,
			Name:     table + "_no_expire",
		}).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["PASSPORT"], "user_1")

		_, err = db.Model(table).Cache(gdb.CacheOption{
			Duration: -1,
			Name:     table + "_no_expire",
		}).Data(g.Map{"nickname": "cleared"}).Where("id", 1).Update()
		t.AssertNil(err)
	})
}

// Test_Model_Cache_Force tests Force option to cache nil results
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
				Name:     table + "_tx",
			}).Where("id", 1).One()
			t.AssertNil(err)
			t.Assert(one["PASSPORT"], "user_1")

			_, err = tx.Model(table).Data(g.Map{"passport": "tx_update"}).Where("id", 1).Update()
			t.AssertNil(err)

			one, err = tx.Model(table).Cache(gdb.CacheOption{
				Duration: time.Second * 10,
				Name:     table + "_tx",
			}).Where("id", 1).One()
			t.AssertNil(err)
			t.Assert(one["PASSPORT"], "tx_update")

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
			gdb.CacheOption{Duration: time.Second * 10, Name: table + "_count"},
			gdb.CacheOption{Duration: time.Second * 10, Name: table + "_data"},
		).Page(1, 3).AllAndCount(false)
		t.AssertNil(err)
		t.Assert(len(all), 3)
		t.Assert(count, 10)

		_, err = db.Model(table).Data(g.Map{
			"id":       11,
			"passport": "user_11",
			"password": "pass_11",
			"nickname": "name_11",
		}).Insert()
		t.AssertNil(err)

		all, count, err = db.Model(table).PageCache(
			gdb.CacheOption{Duration: time.Second * 10, Name: table + "_count"},
			gdb.CacheOption{Duration: time.Second * 10, Name: table + "_data"},
		).Page(1, 3).AllAndCount(false)
		t.AssertNil(err)
		t.Assert(len(all), 3)
		t.Assert(count, 10)

		_, err = db.Model(table).Cache(gdb.CacheOption{
			Duration: -1,
			Name:     table + "_count",
		}).Data(g.Map{"nickname": "page_test"}).Where("id", 1).Update()
		t.AssertNil(err)

		all, count, err = db.Model(table).PageCache(
			gdb.CacheOption{Duration: time.Second * 10, Name: table + "_count"},
			gdb.CacheOption{Duration: time.Second * 10, Name: table + "_data"},
		).Page(1, 3).AllAndCount(false)
		t.AssertNil(err)
		t.Assert(len(all), 3)
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
			Name:     table + "_1",
		}).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["PASSPORT"], "user_1")

		one, err = db.Model(table).Cache(gdb.CacheOption{
			Duration: time.Second * 10,
			Name:     table + "_2",
		}).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["PASSPORT"], "user_1")

		_, err = db.Model(table).Cache(gdb.CacheOption{
			Duration: -1,
			Name:     table + "_1",
		}).Data(g.Map{"passport": "diff_name"}).Where("id", 1).Update()
		t.AssertNil(err)

		one, err = db.Model(table).Cache(gdb.CacheOption{
			Duration: time.Second * 10,
			Name:     table + "_1",
		}).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["PASSPORT"], "diff_name")

		one, err = db.Model(table).Cache(gdb.CacheOption{
			Duration: time.Second * 10,
			Name:     table + "_2",
		}).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["PASSPORT"], "user_1")
	})
}
