// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package oracle_test

import (
	"fmt"
	"testing"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gconv"
)

// modelSupplementAssertOra asserts that `err` carries the Oracle error `code`.
func modelSupplementAssertOra(t *gtest.T, err error, code string) {
	t.AssertNE(err, nil)
	if !gstr.Contains(err.Error(), code) {
		t.Error(fmt.Sprintf("expected %s, got: %v", code, err))
	}
}

// Test_Model_Batch_InsertWithoutId tests batch insert without primary key, retrieving the ids filled by the trigger.
func Test_Model_Batch_InsertWithoutId(t *testing.T) {
	table := createTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Data(g.List{
			{"passport": "t1", "password": "p1", "nickname": "n1"},
			{"passport": "t2", "password": "p2", "nickname": "n2"},
			{"passport": "t3", "password": "p3", "nickname": "n3"},
			{"passport": "t4", "password": "p4", "nickname": "n4"},
			{"passport": "t5", "password": "p5", "nickname": "n5"},
		}).Batch(2).Insert()
		t.AssertNil(err)
		n, _ := result.RowsAffected()
		t.Assert(n, 5)

		all, err := db.Model(table).Order("id").All()
		t.AssertNil(err)
		t.Assert(gconv.Ints(all.Array("ID")), []int{1, 2, 3, 4, 5})
		t.Assert(all.Array("PASSPORT"), g.Slice{"t1", "t2", "t3", "t4", "t5"})
	})
}

// Test_Model_Batch_SaveResult tests saving back a whole query result, where MERGE counts one
// affected row per updated record.
func Test_Model_Batch_SaveResult(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).All()
		t.AssertNil(err)
		t.Assert(len(result), TableSize)
		for _, v := range result {
			v["NICKNAME"].Set(v["NICKNAME"].String() + v["ID"].String())
		}
		r, e := db.Model(table).Data(result).Save()
		t.Assert(e, nil)
		n, e := r.RowsAffected()
		t.Assert(e, nil)
		t.Assert(n, TableSize)

		array, err := db.Model(table).Order("id").Array("nickname")
		t.AssertNil(err)
		t.Assert(array[0], "name_11")
		t.Assert(array[9], "name_1010")
	})
}

// Test_Model_Batch_ReplaceResult tests replacing with a whole query result, where MERGE counts one
// affected row per updated record.
func Test_Model_Batch_ReplaceResult(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).All()
		t.AssertNil(err)
		t.Assert(len(result), TableSize)
		for _, v := range result {
			v["NICKNAME"].Set(v["NICKNAME"].String() + v["ID"].String())
		}
		r, e := db.Model(table).Data(result).Replace()
		t.Assert(e, nil)
		n, e := r.RowsAffected()
		t.Assert(e, nil)
		t.Assert(n, TableSize)

		count, err := db.Model(table).Count()
		t.AssertNil(err)
		t.Assert(count, TableSize)
		value, err := db.Model(table).Where("id", 5).Value("nickname")
		t.AssertNil(err)
		t.Assert(value, "name_55")
	})
}

// Test_Model_Having_Placeholder tests Having with placeholder arguments, which Oracle only accepts
// with a GROUP BY for the selected columns.
func Test_Model_Having_Placeholder(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(table).Where("id > 1").Having("id > ?", 8).All()
		modelSupplementAssertOra(t, err, "ORA-00979")
	})
	gtest.C(t, func(t *gtest.T) {
		all, err := db.Model(table).Fields("id").Where("id > 1").Group("id").Having("id > ?", 8).Order("id").All()
		t.AssertNil(err)
		t.Assert(len(all), 2)
		t.Assert(gconv.Ints(all.Array("ID")), []int{9, 10})
	})
	gtest.C(t, func(t *gtest.T) {
		all, err := db.Model(table).Fields("id").Where("id > ?", 1).Group("id").Having("id > ?", 8).Order("id").All()
		t.AssertNil(err)
		t.Assert(len(all), 2)
		t.Assert(gconv.Ints(all.Array("ID")), []int{9, 10})
	})
}

// Test_Model_Having_FieldValue tests Having with a field name and its value.
func Test_Model_Having_FieldValue(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		all, err := db.Model(table).Fields("id").Where("id > ?", 1).Group("id").Having("id", 8).All()
		t.AssertNil(err)
		t.Assert(len(all), 1)
		t.Assert(all[0]["ID"], 8)
	})
}

// Test_Model_Delete_Limit tests DELETE with LIMIT, which the driver does not rewrite for Oracle.
func Test_Model_Delete_Limit(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	// DELETE...LIMIT
	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(table).Where(1).Limit(2).Delete()
		modelSupplementAssertOra(t, err, "ORA-00920")

		_, err = db.Model(table).Where("1=1").Limit(2).Delete()
		modelSupplementAssertOra(t, err, "ORA-00933")

		count, err := db.Model(table).Count()
		t.AssertNil(err)
		t.Assert(count, TableSize)
	})

	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Where("1=1").Delete()
		t.AssertNil(err)
		n, _ := result.RowsAffected()
		t.Assert(n, TableSize)
	})
}

// Test_Model_Distinct_AliasHaving tests DISTINCT on an aliased table together with HAVING, which
// Oracle only accepts with a GROUP BY for the selected column.
func Test_Model_Distinct_AliasHaving(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(table, "t").Fields("distinct t.id").Where("id > 1").Having("id > 8").All()
		modelSupplementAssertOra(t, err, "ORA-00979")

		all, err := db.Model(table, "t").Fields("distinct t.id").Where("id > 1").Group("t.id").Having("t.id > 8").Order("t.id").All()
		t.AssertNil(err)
		t.Assert(len(all), 2)
		t.Assert(gconv.Ints(all.Array("ID")), []int{9, 10})
	})
}

// Test_Model_InsertAndGetId_WithPrimaryKey tests InsertAndGetId with the primary key in the data.
func Test_Model_InsertAndGetId_WithPrimaryKey(t *testing.T) {
	table := createTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(table).Data(g.Map{
			"id":       1,
			"passport": "user_1",
			"password": "pass_1",
			"nickname": "name_1",
		}).InsertAndGetId()
		t.AssertNil(err)

		one, err := db.Model(table).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["PASSPORT"], "user_1")
		t.Assert(one["NICKNAME"], "name_1")
	})

	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(table).Data(g.Map{
			"passport": "user_2",
			"password": "pass_2",
			"nickname": "name_2",
		}).InsertAndGetId()
		modelSupplementAssertOra(t, err, "ORA-00001")

		count, err := db.Model(table).Count()
		t.AssertNil(err)
		t.Assert(count, 1)
	})
}

// Test_Model_Update_Limit tests UPDATE with ORDER BY and LIMIT, which the driver does not rewrite for Oracle.
func Test_Model_Update_Limit(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	// UPDATE...LIMIT
	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(table).Data("nickname", "T100").Where(1).Order("id desc").Limit(2).Update()
		modelSupplementAssertOra(t, err, "ORA-00920")

		_, err = db.Model(table).Data("nickname", "T100").Where("1=1").Order("id desc").Limit(2).Update()
		modelSupplementAssertOra(t, err, "ORA-00933")

		v1, err := db.Model(table).Fields("nickname").Where("id", 10).Value()
		t.AssertNil(err)
		t.Assert(v1.String(), "name_10")

		count, err := db.Model(table).Where("nickname", "T100").Count()
		t.AssertNil(err)
		t.Assert(count, 0)
	})
}

// Test_DB_BatchInsert_DifferentFields tests batch insert of records with different fields.
func Test_DB_BatchInsert_DifferentFields(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		table := createTable()
		defer dropTable(table)
		r, err := db.Insert(ctx, table, g.List{
			{
				"ID":          2,
				"PASSPORT":    "t2",
				"PASSWORD":    "25d55ad283aa400af464c76d713c07ac",
				"NICKNAME":    "name_2",
				"CREATE_TIME": gtime.Now().String(),
			},
			{
				"ID":          3,
				"PASSPORT":    "user_3",
				"PASSWORD":    "25d55ad283aa400af464c76d713c07ad",
				"NICKNAME":    "name_3",
				"SALARY":      2653.35,
				"CREATE_TIME": gtime.Now().String(),
			},
		}, 1)
		t.AssertNil(err)
		n, err := r.RowsAffected()
		t.AssertNil(err)
		t.Assert(n, 2)

		all, err := db.Model(table).Order("id").All()
		t.AssertNil(err)
		t.Assert(len(all), 2)
		t.Assert(all[0]["SALARY"].IsNil(), true)
		t.Assert(all[1]["SALARY"].Float64(), 2653.35)
	})
}

// Test_Model_Where_MapSliceWhereOr tests a map condition with slice values combined with WhereOr and Where.
func Test_Model_Where_MapSliceWhereOr(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Where(g.Map{
			"id":       g.Slice{1, 2, 3},
			"passport": g.Slice{"user_2", "user_3"},
		}).WhereOr("nickname=?", g.Slice{"name_4"}).Where("id", 3).Order("id").One()
		t.AssertNil(err)
		t.Assert(len(result), 6)
		t.Assert(result["ID"].Int(), 2)

		all, err := db.Model(table).Where(g.Map{
			"id":       g.Slice{1, 2, 3},
			"passport": g.Slice{"user_2", "user_3"},
		}).WhereOr("nickname=?", g.Slice{"name_4"}).Where("id", 3).Order("id").All()
		t.AssertNil(err)
		t.Assert(gconv.Ints(all.Array("ID")), []int{2, 3})
	})
}
