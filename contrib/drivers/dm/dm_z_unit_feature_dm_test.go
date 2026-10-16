// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package dm_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/text/gregex"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gconv"
)

func Test_DB_Insert_ToSQL(t *testing.T) {
	table := createUserTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		sql, err := gdb.ToSQL(ctx, func(ctx context.Context) error {
			_, err := db.Ctx(ctx).Model(table).Data(g.Map{
				"id":       1,
				"passport": "t1",
				"password": "p1",
				"nickname": "n1",
			}).Insert()
			return err
		})
		t.AssertNil(err)
		t.Assert(sql, fmt.Sprintf("INSERT INTO %s(ID,PASSPORT,PASSWORD,NICKNAME) VALUES(1,'t1','p1','n1') ", table))

		count, err := db.Model(table).Count()
		t.AssertNil(err)
		t.Assert(count, 0)
	})

	gtest.C(t, func(t *gtest.T) {
		var id int64
		sql, err := gdb.ToSQL(ctx, func(ctx context.Context) (err error) {
			id, err = db.Ctx(ctx).Model(table).Data(g.Map{
				"passport": "t2",
				"password": "p2",
				"nickname": "n2",
			}).InsertAndGetId()
			return err
		})
		t.AssertNil(err)
		t.Assert(id, 0)
		t.Assert(sql, fmt.Sprintf("INSERT INTO %s(PASSPORT,PASSWORD,NICKNAME) VALUES('t2','p2','n2') ", table))

		count, err := db.Model(table).Count()
		t.AssertNil(err)
		t.Assert(count, 0)
	})
}

func Test_DB_Insert_CatchSQL(t *testing.T) {
	table := createUserTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		var id int64
		sqlArray, err := gdb.CatchSQL(ctx, func(ctx context.Context) (err error) {
			id, err = db.Ctx(ctx).Model(table).Data(g.Map{
				"passport": "t1",
				"password": "p1",
				"nickname": "n1",
			}).InsertAndGetId()
			return err
		})
		t.AssertNil(err)
		t.Assert(id, 1)
		t.AssertGT(len(sqlArray), 0)
		t.Assert(
			sqlArray[len(sqlArray)-1],
			fmt.Sprintf("INSERT INTO %s(PASSPORT,PASSWORD,NICKNAME) VALUES('t1','p1','n1') ", table),
		)

		one, err := db.Model(table).WherePri(1).One()
		t.AssertNil(err)
		t.Assert(one["PASSPORT"], "t1")
	})
}

func Test_DB_Save_ToSQL(t *testing.T) {
	table := createInitUserTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		sql, err := gdb.ToSQL(ctx, func(ctx context.Context) error {
			_, err := db.Ctx(ctx).Model(table).Data(g.Map{
				"id":       1,
				"passport": "t1",
				"password": "p1",
				"nickname": "n1",
			}).Save()
			return err
		})
		t.AssertNil(err)
		t.Assert(sql, fmt.Sprintf(
			"MERGE INTO %s T1 USING (SELECT 1 AS ID,'n1' AS NICKNAME,'t1' AS PASSPORT,'p1' AS PASSWORD FROM DUAL) T2 ON (T1.ID = T2.ID) "+
				"WHEN NOT MATCHED THEN INSERT(ID,NICKNAME,PASSPORT,PASSWORD) VALUES (T2.ID,T2.NICKNAME,T2.PASSPORT,T2.PASSWORD) "+
				"WHEN MATCHED THEN UPDATE SET T1.NICKNAME = T2.NICKNAME,T1.PASSPORT = T2.PASSPORT,T1.PASSWORD = T2.PASSWORD",
			table,
		))

		one, err := db.Model(table).WherePri(1).One()
		t.AssertNil(err)
		t.Assert(one["NICKNAME"], "name_1")
	})
}

// Test_Model_Batch_Save_DifferentFields tests batch save and insert ignore with records carrying different fields
func Test_Model_Batch_Save_DifferentFields(t *testing.T) {
	table := createInitUserTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Data(g.List{
			g.Map{"id": 1, "nickname": "saved_1"},
			g.Map{"id": 11, "passport": "user_11", "password": "pass_11", "nickname": "name_11"},
		}).Save()
		t.AssertNil(err)
		// Note: DM MERGE counts one affected row per merged record, also for an updated one.
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

func Test_OnDuplicateKeyUpdate_Raw(t *testing.T) {
	table := createDuplicateTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(table).Data(g.Map{
			"email":       "user1@example.com",
			"username":    "user1",
			"score":       100,
			"login_count": 1,
		}).OnConflict("email").Save()
		t.AssertNil(err)

		_, err = db.Model(table).Data(g.Map{
			"email":       "user1@example.com",
			"username":    "user1_v2",
			"score":       50,
			"login_count": 1,
		}).OnConflict("email").OnDuplicate(gdb.Raw("T1.login_count = T1.login_count + T2.login_count")).Save()
		t.AssertNil(err)

		one, err := db.Model(table).Where("email", "user1@example.com").One()
		t.AssertNil(err)
		t.Assert(one["USERNAME"], "user1")
		t.Assert(one["SCORE"], 100)
		t.Assert(one["LOGIN_COUNT"], 2)
	})

	gtest.C(t, func(t *gtest.T) {
		for _, score := range []int{50, 150} {
			_, err := db.Model(table).Data(g.Map{
				"email":    "user1@example.com",
				"username": "user1_v3",
				"score":    score,
			}).OnConflict("email").OnDuplicate(g.Map{
				"score": gdb.Raw("GREATEST(T1.score, T2.score)"),
			}).Save()
			t.AssertNil(err)
		}

		one, err := db.Model(table).Where("email", "user1@example.com").One()
		t.AssertNil(err)
		t.Assert(one["USERNAME"], "user1")
		t.Assert(one["SCORE"], 150)

		count, err := db.Model(table).Count()
		t.AssertNil(err)
		t.Assert(count, 1)
	})
}

// Test_Model_LockUpdateSkipLocked tests the LockUpdateSkipLocked convenience method
// Note: MySQL comments this test out for servers before 8.0; DM supports SKIP LOCKED.
func Test_Model_LockUpdateSkipLocked(t *testing.T) {
	table := createInitUserTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		// Test LockUpdateSkipLocked basic usage
		one, err := db.Model(table).LockUpdateSkipLocked().Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["ID"], 1)

		// Test LockUpdateSkipLocked with All()
		all, err := db.Model(table).LockUpdateSkipLocked().Where("id>?", 7).Order("id").All()
		t.AssertNil(err)
		t.Assert(len(all), 3)
	})
}

// Test_Model_Lock_WithLimit tests Lock methods on queries paginated with LIMIT
func Test_Model_Lock_WithLimit(t *testing.T) {
	table := createInitUserTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		tx, err := db.Begin(ctx)
		t.AssertNil(err)
		defer tx.Rollback()

		one, err := tx.Model(table).LockUpdate().Where("id>?", 5).One()
		t.AssertNil(err)
		locked := one["ID"].Int()
		t.AssertGE(locked, 6)
		for id := 6; id <= TableSize; id++ {
			if id == locked {
				lockAssertHeld(t, table, id)
			} else {
				t.AssertNil(lockTry(table, gdb.LockForUpdateNowait, id))
			}
		}
	})

	gtest.C(t, func(t *gtest.T) {
		tx, err := db.Begin(ctx)
		t.AssertNil(err)
		defer tx.Rollback()

		one, err := tx.Model(table).LockUpdate().Where("id>?", 5).Order("id").One()
		t.AssertNil(err)
		t.Assert(one["ID"], 6)

		all, err := tx.Model(table).LockUpdate().Where("id>?", 5).Order("id").Page(2, 2).All()
		t.AssertNil(err)
		t.Assert(len(all), 2)
		t.Assert(all[0]["ID"], 8)
		t.Assert(all[1]["ID"], 9)
		for id := 6; id <= TableSize; id++ {
			if id == 6 || id == 8 || id == 9 {
				lockAssertHeld(t, table, id)
			} else {
				t.AssertNil(lockTry(table, gdb.LockForUpdateNowait, id))
			}
		}
	})
}

// Test_Model_Join_AsJoinedTable tests As() aliasing the last joined table
func Test_Model_Join_AsJoinedTable(t *testing.T) {
	var (
		table1 = createInitUserTable()
		table2 = createInitUserTable()
	)
	defer dropTable(table1)
	defer dropTable(table2)

	gtest.C(t, func(t *gtest.T) {
		r, err := db.Model(table1+" AS t1").
			FieldsPrefix("t1", "id").
			FieldsPrefix("t2", "nickname").
			LeftJoin(table2, "t1.id = t2.id").As("t2").
			Where("t1.id IN(?)", g.Slice{1, 2}).
			Order("t1.id asc").
			All()
		t.AssertNil(err)

		t.Assert(len(r), 2)
		t.Assert(r[0]["ID"], "1")
		t.Assert(r[0]["NICKNAME"], "name_1")
		t.Assert(r[1]["ID"], "2")
		t.Assert(r[1]["NICKNAME"], "name_2")
	})
}

// Test_Model_Page_SingleRow tests Page with a page size of one row
func Test_Model_Page_SingleRow(t *testing.T) {
	table := createInitUserTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		for page := 1; page <= TableSize; page++ {
			result, err := db.Model(table).Order("id").Page(page, 1).All()
			t.AssertNil(err)
			t.Assert(gconv.Ints(result.Array("ID")), []int{page})
		}

		result, err := db.Model(table).Order("id").Page(TableSize+1, 1).All()
		t.AssertNil(err)
		t.Assert(len(result), 0)
	})

	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).OrderDesc("id").Page(1, 1).All()
		t.AssertNil(err)
		t.Assert(gconv.Ints(result.Array("ID")), []int{10})

		one, err := db.Model(table).Order("id").Page(3, 1).One()
		t.AssertNil(err)
		t.Assert(one["ID"], 3)
		t.Assert(len(one), 6)
	})
}

// Test_Model_Limit_SingleRow tests Limit taking one row, with and without offset
func Test_Model_Limit_SingleRow(t *testing.T) {
	table := createInitUserTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Order("id").Limit(0, 1).All()
		t.AssertNil(err)
		t.Assert(gconv.Ints(result.Array("ID")), []int{1})

		result, err = db.Model(table).Order("id").Limit(1).All()
		t.AssertNil(err)
		t.Assert(gconv.Ints(result.Array("ID")), []int{1})

		result, err = db.Model(table).OrderDesc("id").Limit(1).All()
		t.AssertNil(err)
		t.Assert(gconv.Ints(result.Array("ID")), []int{10})
	})

	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Order("id").Limit(1, 1).All()
		t.AssertNil(err)
		t.Assert(gconv.Ints(result.Array("ID")), []int{2})

		result, err = db.Model(table).Order("id").Limit(9, 1).All()
		t.AssertNil(err)
		t.Assert(gconv.Ints(result.Array("ID")), []int{10})

		result, err = db.Model(table).Order("id").Limit(10, 1).All()
		t.AssertNil(err)
		t.Assert(len(result), 0)
	})
}

// Test_Model_Chunk_SizeOne tests Chunk with a chunk size of one row
func Test_Model_Chunk_SizeOne(t *testing.T) {
	table := createInitUserTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		var (
			sizes []int
			ids   []int
		)
		db.Model(table).Order("id").Chunk(1, func(result gdb.Result, err error) bool {
			t.AssertNil(err)
			sizes = append(sizes, len(result))
			ids = append(ids, gconv.Ints(result.Array("ID"))...)
			return true
		})
		t.Assert(sizes, []int{1, 1, 1, 1, 1, 1, 1, 1, 1, 1})
		t.Assert(ids, []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10})
	})

	gtest.C(t, func(t *gtest.T) {
		var ids []int
		db.Model(table).Where("id > ?", 6).OrderDesc("id").Chunk(1, func(result gdb.Result, err error) bool {
			t.AssertNil(err)
			t.Assert(len(result), 1)
			ids = append(ids, gconv.Ints(result.Array("ID"))...)
			return true
		})
		t.Assert(ids, []int{10, 9, 8, 7})
	})

	gtest.C(t, func(t *gtest.T) {
		var ids []int
		db.Model(table).Order("id").Chunk(1, func(result gdb.Result, err error) bool {
			t.AssertNil(err)
			ids = append(ids, gconv.Ints(result.Array("ID"))...)
			return len(ids) < 3
		})
		t.Assert(ids, []int{1, 2, 3})
	})
}

// Test_Model_Page_ZeroPage tests that Page(0, n) returns the first page for various sizes
func Test_Model_Page_ZeroPage(t *testing.T) {
	table := createInitUserTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Order("id").Page(0, 1).All()
		t.AssertNil(err)
		t.Assert(gconv.Ints(result.Array("ID")), []int{1})

		result, err = db.Model(table).OrderDesc("id").Page(0, 4).All()
		t.AssertNil(err)
		t.Assert(gconv.Ints(result.Array("ID")), []int{10, 9, 8, 7})

		result, err = db.Model(table).Order("id").Page(0, 20).All()
		t.AssertNil(err)
		t.Assert(gconv.Ints(result.Array("ID")), []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10})
	})

	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Where("id > ?", 5).Order("id").Page(0, 2).All()
		t.AssertNil(err)
		t.Assert(gconv.Ints(result.Array("ID")), []int{6, 7})

		first, err := db.Model(table).Order("id").Page(0, 3).All()
		t.AssertNil(err)
		second, err := db.Model(table).Order("id").Page(1, 3).All()
		t.AssertNil(err)
		t.Assert(gconv.Ints(first.Array("ID")), gconv.Ints(second.Array("ID")))
	})
}

// Test_Model_Limit_UnalignedOffset tests Limit with an offset that is not a multiple of the size
func Test_Model_Limit_UnalignedOffset(t *testing.T) {
	table := createInitUserTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Order("id").Limit(4, 3).All()
		t.AssertNil(err)
		t.Assert(gconv.Ints(result.Array("ID")), []int{5, 6, 7})

		result, err = db.Model(table).Order("id").Limit(7, 3).All()
		t.AssertNil(err)
		t.Assert(gconv.Ints(result.Array("ID")), []int{8, 9, 10})

		result, err = db.Model(table).Order("id").Limit(9, 3).All()
		t.AssertNil(err)
		t.Assert(gconv.Ints(result.Array("ID")), []int{10})

		result, err = db.Model(table).Order("id").Limit(1, 4).All()
		t.AssertNil(err)
		t.Assert(gconv.Ints(result.Array("ID")), []int{2, 3, 4, 5})

		result, err = db.Model(table).Order("id").Limit(8, 5).All()
		t.AssertNil(err)
		t.Assert(gconv.Ints(result.Array("ID")), []int{9, 10})
	})

	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).OrderDesc("id").Limit(5, 2).All()
		t.AssertNil(err)
		t.Assert(gconv.Ints(result.Array("ID")), []int{5, 4})

		result, err = db.Model(table).Order("id").Limit(3).Offset(4).All()
		t.AssertNil(err)
		t.Assert(gconv.Ints(result.Array("ID")), []int{5, 6, 7})

		result, err = db.Model(table).Where("id > ?", 2).Order("id").Limit(1, 3).All()
		t.AssertNil(err)
		t.Assert(gconv.Ints(result.Array("ID")), []int{4, 5, 6})
	})
}

// Test_Raw_Save tests raw values in the data of Save.
func Test_Raw_Save(t *testing.T) {
	table := createInitUserTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(table).Data(g.Map{
			"id":          1,
			"passport":    gdb.Raw("'port_' || 'raw'"),
			"password":    "pass_1",
			"nickname":    gdb.Raw("UPPER('name_raw')"),
			"create_time": gdb.Raw("now()"),
		}).Save()
		t.AssertNil(err)

		one, err := db.Model(table).WherePri(1).One()
		t.AssertNil(err)
		t.Assert(one["PASSPORT"], "port_raw")
		// Note: DM pads CHAR(32) values with blanks.
		t.Assert(gstr.TrimRight(one["PASSWORD"].String()), "pass_1")
		t.Assert(one["NICKNAME"], "NAME_RAW")
		t.Assert(gregex.IsMatchString(`^\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}$`, one["CREATE_TIME"].String()), true)

		count, err := db.Model(table).Count()
		t.AssertNil(err)
		t.Assert(count, TableSize)
	})
}

// Test_DataType_Text_Large tests character LOB values longer than a VARCHAR bind.
func Test_DataType_Text_Large(t *testing.T) {
	table := "test_text_large_" + gtime.TimestampMicroStr()
	_, err := db.Exec(ctx, "CREATE TABLE "+table+" (id INT PRIMARY KEY AUTO_INCREMENT, data CLOB, data_text TEXT)")
	if err != nil {
		t.Fatal(err)
	}
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		var (
			text      = strings.Repeat("0123456789", 5000)
			nationals = strings.Repeat("中文", 3000)
		)
		_, err := db.Model(table).Data(g.Map{
			"data":      text,
			"data_text": nationals,
		}).Insert()
		t.AssertNil(err)

		one, err := db.Model(table).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["DATA"].String(), text)
		t.Assert(one["DATA_TEXT"].String(), nationals)

		_, err = db.Model(table).Data(g.Map{
			"data":      nationals,
			"data_text": text,
		}).Where("id", 1).Update()
		t.AssertNil(err)

		one, err = db.Model(table).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["DATA"].String(), nationals)
		t.Assert(one["DATA_TEXT"].String(), text)
	})
}

func Test_Model_Fields_SimilarColumns(t *testing.T) {
	table := "user_" + gtime.Now().TimestampNanoStr()
	if _, err := db.Exec(ctx, fmt.Sprintf(`
	    CREATE TABLE %s (
	        id         INT NOT NULL AUTO_INCREMENT,
	        user_id    INT NOT NULL,
	        UserId     INT NOT NULL,
	        PRIMARY KEY (id)
	    )
	    `, table,
	)); err != nil {
		gtest.AssertNil(err)
	}
	defer dropTable(table)

	_, err := db.Model(table).Data(g.Map{"ID": 1, "USER_ID": 2, "USERID": 3}).Insert()
	gtest.AssertNil(err)

	gtest.C(t, func(t *gtest.T) {
		for i := 0; i < 20; i++ {
			one, err := db.Model(table).Fields("user_id").WherePri(1).One()
			t.AssertNil(err)
			t.Assert(one, g.Map{"USER_ID": 2})

			one, err = db.Model(table).Fields("UserId").WherePri(1).One()
			t.AssertNil(err)
			t.Assert(one, g.Map{"USERID": 3})
		}
	})
}

func Test_Model_Where_Struct_SimilarColumns(t *testing.T) {
	table := "user_" + gtime.Now().TimestampNanoStr()
	if _, err := db.Exec(ctx, fmt.Sprintf(`
	    CREATE TABLE %s (
	        id         INT NOT NULL AUTO_INCREMENT,
	        user_id    INT NOT NULL,
	        UserId     INT NOT NULL,
	        PRIMARY KEY (id)
	    )
	    `, table,
	)); err != nil {
		gtest.AssertNil(err)
	}
	defer dropTable(table)

	_, err := db.Model(table).Data(g.List{
		{"ID": 1, "USER_ID": 2, "USERID": 3},
		{"ID": 2, "USER_ID": 2, "USERID": 9},
		{"ID": 3, "USER_ID": 9, "USERID": 3},
	}).Insert()
	gtest.AssertNil(err)

	gtest.C(t, func(t *gtest.T) {
		type User struct {
			Userid int `orm:"user_id"`
			UserId int `orm:"UserId"`
		}
		for i := 0; i < 20; i++ {
			all, err := db.Model(table).Where(User{Userid: 2, UserId: 3}).All()
			t.AssertNil(err)
			t.Assert(len(all), 1)
			t.Assert(all[0]["ID"], 1)
		}
	})
}
