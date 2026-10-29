// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package sqlite_test

import (
	"context"
	"testing"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/test/gtest"
)

// Test_Model_OmitEmpty_Comprehensive tests OmitEmpty filtering for both data and where parameters
func Test_Model_OmitEmpty_Comprehensive(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		// Test OmitEmpty with empty string in Data
		result, err := db.Model(table).OmitEmpty().Data(g.Map{
			"nickname": "",
			"passport": "new_user",
		}).Where("id", 1).Update()
		t.AssertNil(err)
		n, _ := result.RowsAffected()
		t.Assert(n, 1)

		// Verify nickname was not updated (omitted)
		one, err := db.Model(table).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["nickname"], "name_1")
		t.Assert(one["passport"], "new_user")

		// Test OmitEmpty with empty slice in Where
		all, err := db.Model(table).OmitEmpty().Where(g.Map{
			"id":       []int{},
			"passport": "new_user",
		}).All()
		t.AssertNil(err)
		t.Assert(len(all), 1)

		// Without OmitEmpty, empty slice causes WHERE 0=1
		all, err = db.Model(table).Where(g.Map{
			"id": []int{},
		}).All()
		t.AssertNil(err)
		t.Assert(len(all), 0)
	})
}

// Test_Model_OmitEmptyWhere_Extended tests OmitEmpty filtering only for where parameters
func Test_Model_OmitEmptyWhere_Extended(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		// OmitEmptyWhere only affects Where, not Data
		result, err := db.Model(table).OmitEmptyWhere().Data(g.Map{
			"nickname": "",
		}).Where(g.Map{
			"id":       1,
			"passport": "",
		}).Update()
		t.AssertNil(err)
		n, _ := result.RowsAffected()
		t.Assert(n, 1)

		// Verify nickname was updated to empty (Data is not affected by OmitEmptyWhere)
		one, err := db.Model(table).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["nickname"], "")

		// Test with empty slice in Where
		all, err := db.Model(table).OmitEmptyWhere().Where(g.Map{
			"id": []int{},
		}).Order("id").Limit(3).All()
		t.AssertNil(err)
		t.Assert(len(all), 3)

		// Test with zero value in Where (zero is considered empty)
		all, err = db.Model(table).OmitEmptyWhere().Where(g.Map{
			"id": 0,
		}).Order("id").Limit(3).All()
		t.AssertNil(err)
		t.Assert(len(all), 3)
	})
}

// Test_Model_OmitEmptyData tests OmitEmpty filtering only for data parameters
func Test_Model_OmitEmptyData(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		// OmitEmptyData only affects Data, not Where
		result, err := db.Model(table).OmitEmptyData().Data(g.Map{
			"nickname": "",
			"passport": "test_user",
		}).Where(g.Map{
			"id": 1,
		}).Update()
		t.AssertNil(err)
		n, _ := result.RowsAffected()
		t.Assert(n, 1)

		// Verify nickname was not updated (omitted), passport was updated
		one, err := db.Model(table).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["nickname"], "name_1")
		t.Assert(one["passport"], "test_user")

		// Test Insert with OmitEmptyData
		result, err = db.Model(table).OmitEmptyData().Data(g.Map{
			"id":       100,
			"passport": "user_100",
			"nickname": "",
			"password": "pass_100",
		}).Insert()
		t.AssertNil(err)
		n, _ = result.RowsAffected()
		t.Assert(n, 1)

		// Verify nickname is NULL (was omitted from INSERT)
		one, err = db.Model(table).Where("id", 100).One()
		t.AssertNil(err)
		t.Assert(one["passport"], "user_100")
		t.Assert(one["nickname"].IsNil(), true)
	})
}

// Test_Model_OmitNil_Comprehensive tests OmitNil filtering for both data and where parameters
func Test_Model_OmitNil_Comprehensive(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		// Test OmitNil with nil value in Data
		result, err := db.Model(table).OmitNil().Data(g.Map{
			"nickname": nil,
			"passport": "nil_test",
		}).Where("id", 1).Update()
		t.AssertNil(err)
		n, _ := result.RowsAffected()
		t.Assert(n, 1)

		// Verify nickname was not updated (omitted)
		one, err := db.Model(table).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["nickname"], "name_1")
		t.Assert(one["passport"], "nil_test")

		// Test OmitNil with nil in Where
		all, err := db.Model(table).OmitNil().Where(g.Map{
			"passport": nil,
		}).Order("id").Limit(5).All()
		t.AssertNil(err)
		t.Assert(len(all), 5)

		// Without OmitNil, WHERE passport=NULL (which won't match anything)
		all, err = db.Model(table).Where(g.Map{
			"passport": nil,
		}).All()
		t.AssertNil(err)
		t.Assert(len(all), 0)
	})
}

// Test_Model_OmitNilWhere tests OmitNil filtering only for where parameters
func Test_Model_OmitNilWhere(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		// OmitNilWhere only affects Where, not Data
		result, err := db.Model(table).OmitNilWhere().Data(g.Map{
			"nickname": nil,
		}).Where(g.Map{
			"id":       1,
			"passport": nil,
		}).Update()
		t.AssertNil(err)
		n, _ := result.RowsAffected()
		t.Assert(n, 1)

		// Verify nickname was set to NULL (Data is not affected by OmitNilWhere)
		one, err := db.Model(table).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["nickname"].IsNil(), true)

		// Test with nil in Where
		all, err := db.Model(table).OmitNilWhere().Where(g.Map{
			"passport": nil,
		}).Order("id").Limit(3).All()
		t.AssertNil(err)
		t.Assert(len(all), 3)
	})
}

// Test_Model_OmitNilData tests OmitNil filtering only for data parameters
func Test_Model_OmitNilData(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		// OmitNilData only affects Data, not Where
		result, err := db.Model(table).OmitNilData().Data(g.Map{
			"nickname": nil,
			"passport": "omitnil_test",
		}).Where(g.Map{
			"id": 1,
		}).Update()
		t.AssertNil(err)
		n, _ := result.RowsAffected()
		t.Assert(n, 1)

		// Verify nickname was not updated (omitted), passport was updated
		one, err := db.Model(table).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["nickname"], "name_1")
		t.Assert(one["passport"], "omitnil_test")

		// Test Insert with OmitNilData
		result, err = db.Model(table).OmitNilData().Data(g.Map{
			"id":       101,
			"passport": "user_101",
			"nickname": nil,
			"password": "pass_101",
		}).Insert()
		t.AssertNil(err)
		n, _ = result.RowsAffected()
		t.Assert(n, 1)

		// Verify insert
		one, err = db.Model(table).Where("id", 101).One()
		t.AssertNil(err)
		t.Assert(one["passport"], "user_101")
	})
}

// Test_Model_OmitEmpty_WithStruct tests OmitEmpty with struct data
func Test_Model_OmitEmpty_WithStruct(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	type User struct {
		Id       int
		Passport string
		Nickname string
		Password string
	}

	gtest.C(t, func(t *gtest.T) {
		// Test OmitEmptyData with struct
		user := User{
			Passport: "struct_user",
			Nickname: "",
			Password: "struct_pass",
		}
		result, err := db.Model(table).OmitEmptyData().Data(user).Where("id", 1).Update()
		t.AssertNil(err)
		n, _ := result.RowsAffected()
		t.Assert(n, 1)

		// Verify nickname was not updated
		one, err := db.Model(table).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["nickname"], "name_1")
		t.Assert(one["passport"], "struct_user")
	})
}

// Test_Model_OmitNil_WithPointerStruct tests OmitNil with pointer struct data
func Test_Model_OmitNil_WithPointerStruct(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	type User struct {
		Id       int
		Passport *string
		Nickname *string
		Password string
	}

	gtest.C(t, func(t *gtest.T) {
		sqlArray2, err := gdb.CatchSQL(ctx, func(ctx context.Context) error {
			_, err := db.Ctx(ctx).Model(table).OmitNilData().Data(g.Map{
				"passport": "map_user",
				"nickname": nil,
				"password": "map_pass",
			}).Where("id", 2).Update()
			return err
		})
		t.AssertNil(err)
		t.Logf("Map SQL: %v", sqlArray2)

		one2, err := db.Model(table).Where("id", 2).One()
		t.AssertNil(err)
		t.Logf("Map result - nickname: %v, passport: %v", one2["nickname"], one2["passport"])
		t.Assert(one2["nickname"], "name_2")
		t.Assert(one2["passport"], "map_user")
	})
}

// Test_Model_OmitEmpty_ZeroValues tests OmitEmpty with various zero values
func Test_Model_OmitEmpty_ZeroValues(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		// Test OmitEmptyData with various zero values
		result, err := db.Model(table).OmitEmptyData().Data(g.Map{
			"id":       0,
			"passport": "zero_test",
			"nickname": "",
		}).Insert()
		t.AssertNil(err)
		n, _ := result.RowsAffected()
		t.Assert(n, 1)

		// Verify the insert (id should be auto-generated since 0 was omitted)
		one, err := db.Model(table).Where("passport", "zero_test").One()
		t.AssertNil(err)
		t.Assert(one["passport"], "zero_test")
		t.AssertNE(one["id"], 0)
	})
}

// Test_Model_OmitEmpty_ComplexWhere tests OmitEmpty with complex where conditions
func Test_Model_OmitEmpty_ComplexWhere(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		// Test OmitEmptyWhere with multiple conditions
		all, err := db.Model(table).OmitEmptyWhere().Where(g.Map{
			"id >":     0,
			"passport": "",
			"nickname": "?",
		}).Order("id").Limit(3).All()
		t.AssertNil(err)

		// Test with all empty conditions
		all, err = db.Model(table).OmitEmptyWhere().Where(g.Map{
			"passport": "",
			"nickname": "",
		}).Order("id").Limit(5).All()
		t.AssertNil(err)
		t.Assert(len(all), 5)
	})
}

// Test_Model_Omit_ChainedMethods tests Omit methods with other chained methods
func Test_Model_Omit_ChainedMethods(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		// Test OmitEmpty with Fields and Order
		result, err := db.Model(table).
			OmitEmptyData().
			Fields("passport", "nickname").
			Data(g.Map{
				"passport": "chain_test",
				"nickname": "",
			}).
			Where("id", 1).
			Update()
		t.AssertNil(err)
		n, _ := result.RowsAffected()
		t.Assert(n, 1)

		one, err := db.Model(table).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["passport"], "chain_test")
		t.Assert(one["nickname"], "name_1")

		// Test OmitNilWhere with multiple Where clauses
		all, err := db.Model(table).
			OmitNilWhere().
			Where("id>?", 5).
			Where(g.Map{
				"passport": nil,
			}).
			Order("id").
			All()
		t.AssertNil(err)
		t.Assert(len(all), 5)
	})
}
