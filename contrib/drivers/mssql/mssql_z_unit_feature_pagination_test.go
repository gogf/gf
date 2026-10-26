// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package mssql_test

import (
	"testing"
	"time"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/gconv"
)

// Test_Model_AllAndCount_Basic tests basic AllAndCount functionality
func Test_Model_AllAndCount_Basic(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		result, count, err := db.Model(table).AllAndCount(false)
		t.AssertNil(err)
		t.Assert(len(result), TableSize)
		t.Assert(count, TableSize)
	})

	gtest.C(t, func(t *gtest.T) {
		result, count, err := db.Model(table).AllAndCount(true)
		t.AssertNil(err)
		t.Assert(len(result), TableSize)
		t.Assert(count, TableSize)
	})
}

// Test_Model_AllAndCount_WithWhere tests AllAndCount with WHERE conditions
func Test_Model_AllAndCount_WithWhere(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		result, count, err := db.Model(table).Where("id > ?", 5).AllAndCount(false)
		t.AssertNil(err)
		t.Assert(len(result), 5)
		t.Assert(count, 5)
		t.Assert(result[0]["ID"], 6)
	})

	gtest.C(t, func(t *gtest.T) {
		result, count, err := db.Model(table).Where("id", g.Slice{1, 2, 3}).AllAndCount(false)
		t.AssertNil(err)
		t.Assert(len(result), 3)
		t.Assert(count, 3)
	})
}

// Test_Model_AllAndCount_WithPage tests AllAndCount with pagination
func Test_Model_AllAndCount_WithPage(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		result, count, err := db.Model(table).Page(1, 3).Order("id").AllAndCount(false)
		t.AssertNil(err)
		t.Assert(len(result), 3)
		t.Assert(count, TableSize) // Count should be total, not page size
		t.Assert(gconv.Ints(result.Array("ID")), []int{1, 2, 3})
	})

	gtest.C(t, func(t *gtest.T) {
		result, count, err := db.Model(table).Page(2, 3).Order("id").AllAndCount(false)
		t.AssertNil(err)
		t.Assert(len(result), 3)
		t.Assert(count, TableSize)
		t.Assert(gconv.Ints(result.Array("ID")), []int{4, 5, 6})
	})

	gtest.C(t, func(t *gtest.T) {
		result, count, err := db.Model(table).Page(4, 3).OrderDesc("id").AllAndCount(false)
		t.AssertNil(err)
		t.Assert(count, TableSize)
		t.Assert(gconv.Ints(result.Array("ID")), []int{1})
	})

	gtest.C(t, func(t *gtest.T) {
		result, count, err := db.Model(table).Where("id > ?", 2).Page(2, 3).Order("id").AllAndCount(true)
		t.AssertNil(err)
		t.Assert(count, 8)
		t.Assert(gconv.Ints(result.Array("ID")), []int{6, 7, 8})
	})
}

// Test_Model_AllAndCount_WithFields tests AllAndCount with specific fields
// Related: https://github.com/gogf/gf/issues/4698
func Test_Model_AllAndCount_WithFields(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		result, count, err := db.Model(table).Fields("id, nickname").AllAndCount(false)
		t.AssertNil(err)
		t.Assert(len(result), TableSize)
		t.Assert(count, TableSize)
		t.Assert(len(result[0]), 2) // Only 2 fields
	})

	// Regression test for #4698: AllAndCount(true) with multiple fields should work correctly
	// https://github.com/gogf/gf/issues/4698
	gtest.C(t, func(t *gtest.T) {
		result, count, err := db.Model(table).Fields("id, nickname").AllAndCount(true)
		t.AssertNil(err)
		t.Assert(len(result), TableSize)
		t.Assert(count, TableSize)
		t.Assert(len(result[0]), 2)
	})

	gtest.C(t, func(t *gtest.T) {
		result, count, err := db.Model(table).Fields("id, nickname").Order("id").Page(2, 4).AllAndCount(false)
		t.AssertNil(err)
		t.Assert(count, TableSize)
		t.Assert(len(result), 4)
		t.Assert(len(result[0]), 2)
		t.Assert(gconv.Ints(result.Array("ID")), []int{5, 6, 7, 8})
		t.Assert(result.Array("NICKNAME"), g.Slice{"name_5", "name_6", "name_7", "name_8"})
	})
}

// Test_Model_AllAndCount_Empty tests AllAndCount with no results
func Test_Model_AllAndCount_Empty(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		result, count, err := db.Model(table).Where("id > ?", 1000).AllAndCount(false)
		t.AssertNil(err)
		t.Assert(len(result), 0)
		t.Assert(count, 0)
	})

	gtest.C(t, func(t *gtest.T) {
		result, count, err := db.Model(table).Where("id < ?", 0).AllAndCount(true)
		t.AssertNil(err)
		t.Assert(len(result), 0)
		t.Assert(count, 0)
	})

	gtest.C(t, func(t *gtest.T) {
		result, count, err := db.Model(table).Order("id").Page(5, 3).AllAndCount(false)
		t.AssertNil(err)
		t.Assert(len(result), 0)
		t.Assert(count, TableSize)
	})
}

// Test_Model_AllAndCount_WithCache tests AllAndCount with cache
func Test_Model_AllAndCount_WithCache(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		result1, count1, err := db.Model(table).PageCache(gdb.CacheOption{
			Duration: time.Second * 10,
			Force:    false,
		}, gdb.CacheOption{
			Duration: time.Second * 10,
			Force:    false,
		}).Page(1, 5).AllAndCount(false)
		t.AssertNil(err)
		t.Assert(len(result1), 5)
		t.Assert(count1, TableSize)

		// Second call should use cache
		result2, count2, err := db.Model(table).PageCache(gdb.CacheOption{
			Duration: time.Second * 10,
			Force:    false,
		}, gdb.CacheOption{
			Duration: time.Second * 10,
			Force:    false,
		}).Page(1, 5).AllAndCount(false)
		t.AssertNil(err)
		t.Assert(len(result2), 5)
		t.Assert(count2, count1)
	})
}

// Test_Model_AllAndCount_Distinct tests AllAndCount with DISTINCT
func Test_Model_AllAndCount_Distinct(t *testing.T) {
	table := createTable()
	defer dropTable(table)

	// Insert duplicate nicknames
	for i := 1; i <= 10; i++ {
		nickname := "name_" + gconv.String((i-1)/2) // Creates duplicates
		db.Model(table).Data(g.Map{
			"id":       i,
			"passport": "pass_" + gconv.String(i),
			"password": "pwd",
			"nickname": nickname,
		}).Insert()
	}

	gtest.C(t, func(t *gtest.T) {
		result, count, err := db.Model(table).Fields("DISTINCT nickname").AllAndCount(true)
		t.AssertNil(err)
		t.Assert(count, 5) // 10 records / 2 = 5 distinct nicknames
		t.Assert(len(result), 5)
	})

	gtest.C(t, func(t *gtest.T) {
		result, count, err := db.Model(table).Fields("DISTINCT nickname").Order("nickname").Page(2, 2).AllAndCount(true)
		t.AssertNil(err)
		t.Assert(count, 5)
		t.Assert(result.Array("nickname"), g.Slice{"name_2", "name_3"})
	})

	gtest.C(t, func(t *gtest.T) {
		result, count, err := db.Model(table).
			Fields("nickname", "COUNT(1) AS cnt").
			Group("nickname").
			Order("nickname").
			Page(2, 2).
			AllAndCount(false)
		t.AssertNil(err)
		t.Assert(count, 5)
		t.Assert(result.Array("NICKNAME"), g.Slice{"name_2", "name_3"})
		t.Assert(gconv.Ints(result.Array("cnt")), []int{2, 2})
	})
}

// Test_Model_ScanAndCount_Basic tests basic ScanAndCount functionality
func Test_Model_ScanAndCount_Basic(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	type User struct {
		Id       int
		Passport string
		Password string
		Nickname string
	}

	gtest.C(t, func(t *gtest.T) {
		var users []User
		var count int
		err := db.Model(table).ScanAndCount(&users, &count, false)
		t.AssertNil(err)
		t.Assert(len(users), TableSize)
		t.Assert(count, TableSize)
	})

	gtest.C(t, func(t *gtest.T) {
		var users []User
		var count int
		err := db.Model(table).ScanAndCount(&users, &count, true)
		t.AssertNil(err)
		t.Assert(len(users), TableSize)
		t.Assert(count, TableSize)
	})
}

// Test_Model_ScanAndCount_WithWhere tests ScanAndCount with WHERE conditions
func Test_Model_ScanAndCount_WithWhere(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	type User struct {
		Id       int
		Passport string
		Nickname string
	}

	gtest.C(t, func(t *gtest.T) {
		var users []User
		var count int
		err := db.Model(table).Where("id <= ?", 5).ScanAndCount(&users, &count, false)
		t.AssertNil(err)
		t.Assert(len(users), 5)
		t.Assert(count, 5)
		t.Assert(users[0].Id, 1)
	})
}

// Test_Model_ScanAndCount_WithPage tests ScanAndCount with pagination
func Test_Model_ScanAndCount_WithPage(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	type User struct {
		Id       int
		Nickname string
	}

	gtest.C(t, func(t *gtest.T) {
		var users []User
		var count int
		err := db.Model(table).Page(2, 3).Order("id").ScanAndCount(&users, &count, false)
		t.AssertNil(err)
		t.Assert(len(users), 3)
		t.Assert(count, TableSize) // Total count, not page count
		t.Assert(users[0].Id, 4)
		t.Assert(users[1].Id, 5)
		t.Assert(users[2].Id, 6)
		t.Assert(users[2].Nickname, "name_6")
	})

	gtest.C(t, func(t *gtest.T) {
		var users []User
		var count int
		err := db.Model(table).Where("id > ?", 3).Page(3, 3).Order("id").ScanAndCount(&users, &count, false)
		t.AssertNil(err)
		t.Assert(count, 7)
		t.Assert(len(users), 1)
		t.Assert(users[0].Id, 10)
	})

	gtest.C(t, func(t *gtest.T) {
		var users []User
		var count int
		err := db.Model(table).Limit(1, 2).OrderDesc("id").ScanAndCount(&users, &count, false)
		t.AssertNil(err)
		t.Assert(count, TableSize)
		t.Assert(len(users), 2)
		t.Assert(users[0].Id, 9)
		t.Assert(users[1].Id, 8)
	})
}

// Test_Model_ScanAndCount_Single tests ScanAndCount for single record
func Test_Model_ScanAndCount_Single(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	type User struct {
		Id       int
		Passport string
	}

	gtest.C(t, func(t *gtest.T) {
		var user User
		var count int
		err := db.Model(table).Where("id", 1).ScanAndCount(&user, &count, false)
		t.AssertNil(err)
		t.Assert(user.Id, 1)
		t.Assert(count, 1)
	})

	gtest.C(t, func(t *gtest.T) {
		var user User
		var count int
		err := db.Model(table).Where("id > ?", 7).OrderDesc("id").ScanAndCount(&user, &count, false)
		t.AssertNil(err)
		t.Assert(user.Id, 10)
		t.Assert(user.Passport, "user_10")
		t.Assert(count, 3)
	})
}

// Test_Model_ScanAndCount_Empty tests ScanAndCount with no results
func Test_Model_ScanAndCount_Empty(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	type User struct {
		Id int
	}

	gtest.C(t, func(t *gtest.T) {
		var users []User
		var count int
		err := db.Model(table).Where("id > ?", 1000).ScanAndCount(&users, &count, false)
		t.AssertNil(err)
		t.Assert(len(users), 0)
		t.Assert(count, 0)
	})
}

// Test_Model_ScanAndCount_WithFields tests ScanAndCount with specific fields
func Test_Model_ScanAndCount_WithFields(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	type User struct {
		Id       int
		Nickname string
	}

	gtest.C(t, func(t *gtest.T) {
		var users []User
		var count int
		err := db.Model(table).Fields("id, nickname").ScanAndCount(&users, &count, false)
		t.AssertNil(err)
		t.Assert(len(users), TableSize)
		t.Assert(count, TableSize)
		t.Assert(users[0].Id > 0, true)
		t.AssertNE(users[0].Nickname, "")
	})
}

// Test_Model_ScanAndCount_WithCache tests ScanAndCount with cache
func Test_Model_ScanAndCount_WithCache(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	type User struct {
		Id int
	}

	gtest.C(t, func(t *gtest.T) {
		var users1 []User
		var count1 int
		err := db.Model(table).PageCache(gdb.CacheOption{
			Duration: time.Second * 10,
			Force:    false,
		}, gdb.CacheOption{
			Duration: time.Second * 10,
			Force:    false,
		}).Page(1, 5).ScanAndCount(&users1, &count1, false)
		t.AssertNil(err)
		t.Assert(len(users1), 5)
		t.Assert(count1, TableSize)

		// Second call should use cache
		var users2 []User
		var count2 int
		err = db.Model(table).PageCache(gdb.CacheOption{
			Duration: time.Second * 10,
			Force:    false,
		}, gdb.CacheOption{
			Duration: time.Second * 10,
			Force:    false,
		}).Page(1, 5).ScanAndCount(&users2, &count2, false)
		t.AssertNil(err)
		t.Assert(len(users2), 5)
		t.Assert(count2, count1)
	})
}

// Test_Model_Chunk_Basic tests basic Chunk functionality
func Test_Model_Chunk_Basic(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		var (
			total  int
			chunks int
			ids    []int
			sizes  []int
		)
		db.Model(table).Order("id").Chunk(3, func(result gdb.Result, err error) bool {
			t.AssertNil(err)
			chunks++
			total += len(result)
			sizes = append(sizes, len(result))
			ids = append(ids, gconv.Ints(result.Array("ID"))...)
			return true
		})
		t.Assert(chunks, 4) // 10 records / 3 = 4 chunks (3+3+3+1)
		t.Assert(total, TableSize)
		t.Assert(sizes, []int{3, 3, 3, 1})
		t.Assert(ids, []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10})
	})

	gtest.C(t, func(t *gtest.T) {
		var sizes []int
		db.Model(table).Order("id").Chunk(5, func(result gdb.Result, err error) bool {
			t.AssertNil(err)
			sizes = append(sizes, len(result))
			return true
		})
		t.Assert(sizes, []int{5, 5})
	})
}

// Test_Model_Chunk_StopEarly tests Chunk with early stop
func Test_Model_Chunk_StopEarly(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		var (
			chunks int
			ids    []int
		)
		db.Model(table).Order("id").Chunk(3, func(result gdb.Result, err error) bool {
			t.AssertNil(err)
			chunks++
			ids = append(ids, gconv.Ints(result.Array("ID"))...)
			return chunks < 2 // Stop after 2nd chunk
		})
		t.Assert(chunks, 2)
		t.Assert(ids, []int{1, 2, 3, 4, 5, 6})
	})
}

// Test_Model_Chunk_WithWhere tests Chunk with WHERE conditions
func Test_Model_Chunk_WithWhere(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		var (
			total  int
			chunks int
			ids    []int
		)
		db.Model(table).Where("id <= ?", 5).Order("id").Chunk(2, func(result gdb.Result, err error) bool {
			t.AssertNil(err)
			chunks++
			total += len(result)
			ids = append(ids, gconv.Ints(result.Array("ID"))...)
			return true
		})
		t.Assert(chunks, 3) // 5 records / 2 = 3 chunks (2+2+1)
		t.Assert(total, 5)
		t.Assert(ids, []int{1, 2, 3, 4, 5})
	})

	gtest.C(t, func(t *gtest.T) {
		var ids []int
		db.Model(table).Where("id > ?", 3).OrderDesc("id").Chunk(4, func(result gdb.Result, err error) bool {
			t.AssertNil(err)
			ids = append(ids, gconv.Ints(result.Array("ID"))...)
			return true
		})
		t.Assert(ids, []int{10, 9, 8, 7, 6, 5, 4})
	})
}

// Test_Model_Chunk_ErrorHandling tests Chunk error handling
func Test_Model_Chunk_ErrorHandling(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var errorReceived bool
		db.Model("non_existent_table").Chunk(10, func(result gdb.Result, err error) bool {
			if err != nil {
				errorReceived = true
				return false
			}
			return true
		})
		t.Assert(errorReceived, true)
	})
}

// Test_Model_Chunk_Empty tests Chunk with no results
func Test_Model_Chunk_Empty(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		var chunks int
		db.Model(table).Where("id > ?", 1000).Chunk(10, func(result gdb.Result, err error) bool {
			chunks++
			return true
		})
		t.Assert(chunks, 0) // No chunks for empty result
	})
}

// Test_Model_Page_Boundary tests Page with boundary values
// Related: https://github.com/gogf/gf/issues/4699
func Test_Model_Page_Boundary(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	// Page 0 should be treated as page 1
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Page(0, 3).Order("id").All()
		t.AssertNil(err)
		t.Assert(len(result), 3)
		t.Assert(gconv.Ints(result.Array("ID")), []int{1, 2, 3})
	})

	// Negative page should be treated as page 1
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Page(-1, 3).Order("id").All()
		t.AssertNil(err)
		t.Assert(len(result), 3)
		t.Assert(gconv.Ints(result.Array("ID")), []int{1, 2, 3})
	})

	// Size 0: framework treats limit=0 as "no limit", returns all records
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Page(1, 0).All()
		t.AssertNil(err)
		t.Assert(len(result), TableSize)
	})

	// Negative size: normalized to 0, same as Page(1, 0)
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Page(1, -1).All()
		t.AssertNil(err)
		t.Assert(len(result), TableSize)
	})

	// Very large page number (beyond available data)
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Page(100, 3).All()
		t.AssertNil(err)
		t.Assert(len(result), 0)
	})

	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Page(2, 3).Order("id").All()
		t.AssertNil(err)
		t.Assert(gconv.Ints(result.Array("ID")), []int{4, 5, 6})
		t.Assert(len(result[0]), 7)
	})

	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Page(4, 3).Order("id").All()
		t.AssertNil(err)
		t.Assert(gconv.Ints(result.Array("ID")), []int{10})
	})

	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Page(5, 2).Order("id").All()
		t.AssertNil(err)
		t.Assert(gconv.Ints(result.Array("ID")), []int{9, 10})
	})

	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Page(6, 2).Order("id").All()
		t.AssertNil(err)
		t.Assert(len(result), 0)
	})

	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Page(1, 1).OrderDesc("id").All()
		t.AssertNil(err)
		t.Assert(gconv.Ints(result.Array("ID")), []int{10})
	})
}

// Test_Model_Limit_Boundary tests Limit with boundary values
// Related: https://github.com/gogf/gf/issues/4699
func Test_Model_Limit_Boundary(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	// Limit 0: framework treats limit=0 as "no limit", returns all records
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Limit(0).All()
		t.AssertNil(err)
		t.Assert(len(result), TableSize)
	})

	// Negative limit: normalized to 0, same as Limit(0)
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Limit(-1).All()
		t.AssertNil(err)
		t.Assert(len(result), TableSize)
	})

	// Limit larger than available data
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Limit(1000).All()
		t.AssertNil(err)
		t.Assert(len(result), TableSize)
	})

	// Limit(offset, size): offset=5 skips 5 rows, size=100 takes up to 100
	// With 10 rows total, skipping 5 returns remaining 5 rows
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Limit(5, 100).All()
		t.AssertNil(err)
		t.Assert(len(result), TableSize-5)
	})

	// Offset beyond data: returns empty result
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Limit(100, 5).All()
		t.AssertNil(err)
		t.Assert(len(result), 0)
	})

	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Order("id").Limit(3).All()
		t.AssertNil(err)
		t.Assert(gconv.Ints(result.Array("ID")), []int{1, 2, 3})
		t.Assert(len(result[0]), 7)
	})

	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).OrderDesc("id").Limit(1).All()
		t.AssertNil(err)
		t.Assert(gconv.Ints(result.Array("ID")), []int{10})
	})

	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Order("id").Limit(0, 3).All()
		t.AssertNil(err)
		t.Assert(gconv.Ints(result.Array("ID")), []int{1, 2, 3})
	})

	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Order("id").Limit(5, 100).All()
		t.AssertNil(err)
		t.Assert(gconv.Ints(result.Array("ID")), []int{6, 7, 8, 9, 10})
	})

	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Order("id").Limit(2, 3).All()
		t.AssertNil(err)
		t.Assert(gconv.Ints(result.Array("ID")), []int{3, 4, 5})
	})

	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).OrderDesc("id").Limit(1, 2).All()
		t.AssertNil(err)
		t.Assert(gconv.Ints(result.Array("ID")), []int{9, 8})
	})

	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Order("id").Limit(3).Offset(2).All()
		t.AssertNil(err)
		t.Assert(gconv.Ints(result.Array("ID")), []int{3, 4, 5})
	})

	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Order("id").Offset(7).All()
		t.AssertNil(err)
		t.Assert(gconv.Ints(result.Array("ID")), []int{8, 9, 10})
	})
}

// Test_Model_Page_Limit_Combination tests Page and Limit used together
func Test_Model_Page_Limit_Combination(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		// Page should override Limit
		result, err := db.Model(table).Limit(5).Page(1, 3).Order("id").All()
		t.AssertNil(err)
		t.Assert(len(result), 3)
		t.Assert(gconv.Ints(result.Array("ID")), []int{1, 2, 3})
	})

	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Page(2, 3).Limit(4, 2).Order("id").All()
		t.AssertNil(err)
		t.Assert(gconv.Ints(result.Array("ID")), []int{5, 6})
	})

	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Page(2, 3).OrderDesc("id").All()
		t.AssertNil(err)
		t.Assert(gconv.Ints(result.Array("ID")), []int{7, 6, 5})
	})

	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Order("nickname").Order("id").Page(1, 3).All()
		t.AssertNil(err)
		t.Assert(gconv.Ints(result.Array("ID")), []int{1, 10, 2})
	})

	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Where("id > ?", 2).Order("id").Page(2, 3).All()
		t.AssertNil(err)
		t.Assert(gconv.Ints(result.Array("ID")), []int{6, 7, 8})
	})

	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).
			Fields("id % 3 m, COUNT(1) c").
			Group("id % 3").
			Order("m").
			Page(2, 1).
			All()
		t.AssertNil(err)
		t.Assert(len(result), 1)
		t.Assert(result[0]["m"], 1)
		t.Assert(result[0]["c"], 4)
	})

	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Distinct().Fields("password").Order("password").Page(2, 2).All()
		t.AssertNil(err)
		t.Assert(result.Array("PASSWORD"), g.Slice{"pass_2", "pass_3"})
	})

	gtest.C(t, func(t *gtest.T) {
		table2 := createInitTable()
		defer dropTable(table2)

		result, err := db.Model(table+" a").
			InnerJoin(table2+" b", "a.id = b.id").
			Fields("a.id, b.nickname").
			Order("id").
			Page(2, 3).
			All()
		t.AssertNil(err)
		t.Assert(gconv.Ints(result.Array("id")), []int{4, 5, 6})
		t.Assert(result.Array("nickname"), g.Slice{"name_4", "name_5", "name_6"})
	})

	gtest.C(t, func(t *gtest.T) {
		table2 := createInitTable()
		defer dropTable(table2)

		result, err := db.Model(table+" a").
			InnerJoin(table2+" b", "a.id = b.id").
			Fields("a.id, b.nickname").
			Order("a.id").
			Page(2, 3).
			All()
		t.AssertNil(err)
		t.Assert(gconv.Ints(result.Array("id")), []int{4, 5, 6})
		t.Assert(result.Array("nickname"), g.Slice{"name_4", "name_5", "name_6"})
	})

	gtest.C(t, func(t *gtest.T) {
		table2 := createInitTable()
		defer dropTable(table2)

		result, err := db.Model(table+" a").
			InnerJoin(table2+" b", "a.id = b.id").
			Order("a.id").
			Page(2, 3).
			All()
		t.AssertNil(err)
		t.Assert(gconv.Ints(result.Array("ID")), []int{4, 5, 6})
	})

	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model("(?) AS sub", db.Model(table).Where("id <=", 7)).
			Order("id").
			Page(2, 3).
			All()
		t.AssertNil(err)
		t.Assert(gconv.Ints(result.Array("ID")), []int{4, 5, 6})
	})

	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model("(?) AS sub", db.Model(table).Where("id <=", 7)).
			Order("sub.id").
			Page(2, 3).
			All()
		t.AssertNil(err)
		t.Assert(gconv.Ints(result.Array("ID")), []int{4, 5, 6})
	})

	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).
			Where("id IN(?)", db.Model(table).Fields("id").Where("id >", 3)).
			Order("id").
			Page(2, 2).
			All()
		t.AssertNil(err)
		t.Assert(gconv.Ints(result.Array("ID")), []int{6, 7})
	})
}
