// Copyright 2019 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package dm_test

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/test/gtest"
)

func Test_DB_Ping(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		err1 := dblink.PingMaster()
		err2 := dblink.PingSlave()
		t.Assert(err1, nil)
		t.Assert(err2, nil)
	})
}

func TestTables(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		tables := []string{"A_tables", "A_tables2"}

		for _, v := range tables {
			createInitTable(v)
			// createTable(v)
		}
		result, err := db.Tables(ctx)
		gtest.Assert(err, nil)

		for i := 0; i < len(tables); i++ {
			find := false
			for j := 0; j < len(result); j++ {
				if strings.ToUpper(tables[i]) == result[j] {
					find = true
					break
				}
			}
			gtest.AssertEQ(find, true)
		}

		result, err = dblink.Tables(ctx)
		gtest.Assert(err, nil)
		for i := 0; i < len(tables); i++ {
			find := false
			for j := 0; j < len(result); j++ {
				if strings.ToUpper(tables[i]) == result[j] {
					find = true
					break
				}
			}
			gtest.AssertEQ(find, true)
		}
	})
}

func TestTableFields(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		tables := "A_tables"
		var expect = map[string][]interface{}{
			"ID":           {"BIGINT", false},
			"ACCOUNT_NAME": {"VARCHAR", false},
			"PWD_RESET":    {"TINYINT", false},
			"DELETED":      {"INT", false},
			"CREATED_TIME": {"TIMESTAMP", false},
		}

		res, err := db.TableFields(ctx, tables)
		gtest.Assert(err, nil)

		for k, v := range expect {
			_, ok := res[k]
			gtest.AssertEQ(ok, true)

			gtest.AssertEQ(res[k].Name, k)
			gtest.Assert(res[k].Type, v[0])
			gtest.Assert(res[k].Null, v[1])
		}

	})

	gtest.C(t, func(t *gtest.T) {
		_, err := db.TableFields(ctx, "t_user t_user2")
		gtest.AssertNE(err, nil)
	})
}

func Test_DB_Query(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		tableName := "A_tables"
		// createTable(tableName)
		_, err := db.Query(ctx, fmt.Sprintf("SELECT * from %s", tableName))
		t.AssertNil(err)

		resTwo := make([]User, 0)
		err = db.Model(tableName).Scan(&resTwo)
		t.AssertNil(err)

		resThree := make([]User, 0)
		model := db.Model(tableName)
		model.Where("id", g.Slice{1, 2, 3, 4})
		// model.Where("account_name like ?", "%"+"list"+"%")
		model.Where("deleted", 0).Order("pwd_reset desc")
		_, err = model.Count()
		t.AssertNil(err)
		err = model.Page(2, 2).Scan(&resThree)
		t.AssertNil(err)
	})
}

func TestModelSave(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// createTable("A_tables")
		data := []User{
			{
				ID:          100,
				AccountName: "user_100",
				CreatedTime: time.Now(),
			},
		}
		_, err := db.Model("A_tables").Data(data).Save()
		gtest.Assert(err, nil)

		data2 := []User{
			{
				ID:          101,
				AccountName: "user_101",
			},
		}
		_, err = db.Model("A_tables").Data(&data2).Save()
		gtest.Assert(err, nil)

		data3 := []User{
			{
				ID:          10,
				AccountName: "user_10",
				PwdReset:    10,
			},
		}
		_, err = db.Model("A_tables").Save(data3)
		gtest.Assert(err, nil)

		data4 := []User{
			{
				ID:          9,
				AccountName: "user_9",
				CreatedTime: time.Now(),
			},
		}
		_, err = db.Model("A_tables").Save(&data4)
		gtest.Assert(err, nil)

		// TODO:: Should be Supported 'Replace' Operation
		// _, err = db.Schema(TestDbName).Replace(ctx, "DoInsert", data, 10)
		// gtest.Assert(err, nil)
	})
}

func TestModelInsert(t *testing.T) {
	// g.Model.insert not lost default not null coloumn
	gtest.C(t, func(t *gtest.T) {
		// createTable("A_tables")
		i := 200
		data := User{
			ID:          int64(i),
			AccountName: fmt.Sprintf(`A%dtwo`, i),
			PwdReset:    0,
			// CreatedTime: time.Now(),
			UpdatedTime: time.Now(),
		}
		// _, err := db.Schema(TestDbName).Model("A_tables").Data(data).Insert()
		_, err := db.Model("A_tables").Insert(&data)
		gtest.Assert(err, nil)
	})

	gtest.C(t, func(t *gtest.T) {
		// createTable("A_tables")
		i := 201
		data := User{
			ID:          int64(i),
			AccountName: fmt.Sprintf(`A%dtwoONE`, i),
			PwdReset:    1,
			CreatedTime: time.Now(),
			// UpdatedTime: time.Now(),
		}
		// _, err := db.Schema(TestDbName).Model("A_tables").Data(data).Insert()
		_, err := db.Model("A_tables").Data(&data).Insert()
		gtest.Assert(err, nil)
	})
}

func TestDBInsert(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// createTable("A_tables")
		i := 300
		data := g.Map{
			"ID":           i,
			"ACCOUNT_NAME": fmt.Sprintf(`A%dthress`, i),
			"PWD_RESET":    3,
		}
		_, err := db.Insert(ctx, "A_tables", &data)
		gtest.Assert(err, nil)
	})
}

func Test_DB_Exec(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		_, err := db.Exec(ctx, "SELECT ? from dual", 1)
		t.AssertNil(err)

		_, err = db.Exec(ctx, "ERROR")
		t.AssertNE(err, nil)
	})
}

func Test_DB_Insert(t *testing.T) {
	// table := createTable()
	// defer dropTable(table)
	gtest.C(t, func(t *gtest.T) {
		// normal map
		_, err := db.Insert(ctx, "A_tables", g.Map{
			"ID":           1000,
			"ACCOUNT_NAME": "map1",
			"CREATED_TIME": gtime.Now().String(),
		})
		t.AssertNil(err)

		result, err := db.Insert(ctx, "A_tables", g.Map{
			"ID":           "2000",
			"ACCOUNT_NAME": "map2",
			"CREATED_TIME": gtime.Now(),
		})
		t.AssertNil(err)
		n, _ := result.RowsAffected()
		t.Assert(n, 1)

		result, err = db.Insert(ctx, "A_tables", g.Map{
			"ID":           3000,
			"ACCOUNT_NAME": "map3",
			// "CREATED_TIME": gtime.Now().String(),
		})
		t.AssertNil(err)
		n, _ = result.RowsAffected()
		t.Assert(n, 1)

		// struct
		result, err = db.Insert(ctx, "A_tables", User{
			ID:          4000,
			AccountName: "struct_4",
			// CreatedTime: timeStr,
			// UpdatedTime: timeStr,
		})
		t.AssertNil(err)
		n, _ = result.RowsAffected()
		t.Assert(n, 1)

		ones, err := db.Model("A_tables").Where("ID", 4000).All()
		t.AssertNil(err)
		t.Assert(ones[0]["ID"].Int(), 4000)
		t.Assert(ones[0]["ACCOUNT_NAME"].String(), "struct_4")
		// TODO Question2
		// this is DM bug.
		// t.Assert(one["CREATED_TIME"].GTime().String(), timeStr)

		// *struct
		timeStr := time.Now()
		result, err = db.Insert(ctx, "A_tables", &User{
			ID:          5000,
			AccountName: "struct_5",
			CreatedTime: timeStr,
			// UpdatedTime: timeStr,
		})
		t.AssertNil(err)
		n, _ = result.RowsAffected()
		t.Assert(n, 1)

		one, err := db.Model("A_tables").Where("ID", 5000).One()
		t.AssertNil(err)
		t.Assert(one["ID"].Int(), 5000)
		t.Assert(one["ACCOUNT_NAME"].String(), "struct_5")

		// batch with Insert
		r, err := db.Insert(ctx, "A_tables", g.Slice{
			g.Map{
				"ID":           6000,
				"ACCOUNT_NAME": "t6000",
			},
			g.Map{
				"ID":           6001,
				"ACCOUNT_NAME": "t6001",
			},
		})
		t.AssertNil(err)
		n, _ = r.RowsAffected()
		t.Assert(n, 2)

		one, err = db.Model("A_tables").Where("ID", 6000).One()
		t.AssertNil(err)
		t.Assert(one["ID"].Int(), 6000)
		t.Assert(one["ACCOUNT_NAME"].String(), "t6000")
	})
}

func Test_DB_BatchInsert(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		table := "A_tables"
		r, err := db.Insert(ctx, table, g.List{
			{
				"ID":           400,
				"ACCOUNT_NAME": "list_400",
				// "CREATE_TIME":  gtime.Now().String(),
			},
			{
				"ID":           401,
				"ACCOUNT_NAME": "list_401",
				"CREATE_TIME":  gtime.Now().String(),
			},
		}, 1)
		t.AssertNil(err)
		n, _ := r.RowsAffected()
		t.Assert(n, 2)
	})

	gtest.C(t, func(t *gtest.T) {
		table := "A_tables"
		// table := createTable()
		// defer dropTable(table)
		// []interface{}
		r, err := db.Insert(ctx, table, g.Slice{
			g.Map{
				"ID":           500,
				"ACCOUNT_NAME": "500_batch_500",
				"CREATE_TIME":  gtime.Now().String(),
			},
			g.Map{
				"ID":           501,
				"ACCOUNT_NAME": "501_batch_501",
				// "CREATE_TIME":  gtime.Now().String(),
			},
		}, 1)
		t.AssertNil(err)
		n, _ := r.RowsAffected()
		t.Assert(n, 2)
	})

	// batch insert map
	gtest.C(t, func(t *gtest.T) {
		table := "A_tables"
		// table := createTable()
		// defer dropTable(table)
		result, err := db.Insert(ctx, table, g.Map{
			"ID":           600,
			"ACCOUNT_NAME": "600_batch_600",
			"CREATE_TIME":  gtime.Now().String(),
		})
		t.AssertNil(err)
		n, _ := result.RowsAffected()
		t.Assert(n, 1)
	})
}

func Test_DB_BatchInsert_Struct(t *testing.T) {
	// batch insert struct
	gtest.C(t, func(t *gtest.T) {
		table := "A_tables"
		// table := createTable()
		// defer dropTable(table)
		user := &User{
			ID:          700,
			AccountName: "BatchInsert_Struct_700",
			// CreatedTime: time.Now(),
		}
		result, err := db.Model(table).Insert(user)
		t.AssertNil(err)
		n, _ := result.RowsAffected()
		t.Assert(n, 1)
	})
}

func Test_DB_Update(t *testing.T) {
	table := "A_tables"
	// table := createInitTable()

	gtest.C(t, func(t *gtest.T) {
		result, err := db.Update(ctx, table, "pwd_reset=7", "id=700")
		t.AssertNil(err)
		n, _ := result.RowsAffected()
		t.Assert(n, 1)

		one, err := db.Model(table).Where("ID", 700).One()
		t.AssertNil(err)
		t.Assert(one["ID"].Int(), 700)
		t.Assert(one["ACCOUNT_NAME"].String(), "BatchInsert_Struct_700")
		t.Assert(one["PWD_RESET"].String(), "7")
	})
}

func Test_DB_GetAll(t *testing.T) {
	table := "A_tables"
	// table := createInitTable()
	// defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		result, err := db.GetAll(ctx, fmt.Sprintf("SELECT * FROM %s WHERE id=?", table), 1)
		t.AssertNil(err)
		t.Assert(len(result), 1)
		t.Assert(result[0]["ID"].Int(), 1)
	})
	gtest.C(t, func(t *gtest.T) {
		result, err := db.GetAll(ctx, fmt.Sprintf("SELECT * FROM %s WHERE id=?", table), g.Slice{1})
		t.AssertNil(err)
		t.Assert(len(result), 1)
		t.Assert(result[0]["ID"].Int(), 1)
	})
	gtest.C(t, func(t *gtest.T) {
		result, err := db.GetAll(ctx, fmt.Sprintf("SELECT * FROM %s WHERE id in(?)", table), g.Slice{1, 2, 3})
		t.AssertNil(err)
		t.Assert(len(result), 3)
		t.Assert(result[0]["ID"].Int(), 1)
		t.Assert(result[1]["ID"].Int(), 2)
		t.Assert(result[2]["ID"].Int(), 3)
	})
	gtest.C(t, func(t *gtest.T) {
		result, err := db.GetAll(ctx, fmt.Sprintf("SELECT * FROM %s WHERE id in(?,?,?)", table), g.Slice{1, 2, 3})
		t.AssertNil(err)
		t.Assert(len(result), 3)
		t.Assert(result[0]["ID"].Int(), 1)
		t.Assert(result[1]["ID"].Int(), 2)
		t.Assert(result[2]["ID"].Int(), 3)
	})
	gtest.C(t, func(t *gtest.T) {
		result, err := db.GetAll(ctx, fmt.Sprintf("SELECT * FROM %s WHERE id in(?,?,?)", table), g.Slice{1, 2, 3}...)
		t.AssertNil(err)
		t.Assert(len(result), 3)
		t.Assert(result[0]["ID"].Int(), 1)
		t.Assert(result[1]["ID"].Int(), 2)
		t.Assert(result[2]["ID"].Int(), 3)
	})
	gtest.C(t, func(t *gtest.T) {
		result, err := db.GetAll(ctx, fmt.Sprintf("SELECT * FROM %s WHERE id>=? AND id <=?", table), g.Slice{1, 3})
		t.AssertNil(err)
		t.Assert(len(result), 3)
		t.Assert(result[0]["ID"].Int(), 1)
		t.Assert(result[1]["ID"].Int(), 2)
		t.Assert(result[2]["ID"].Int(), 3)
	})
}

func Test_DB_GetOne(t *testing.T) {
	// table := createInitTable()
	table := "A_tables"
	gtest.C(t, func(t *gtest.T) {
		record, err := db.GetOne(ctx, fmt.Sprintf("SELECT * FROM %s WHERE account_name=?", table), "struct_4")
		t.AssertNil(err)
		t.Assert(record["ACCOUNT_NAME"].String(), "struct_4")
	})
}

func Test_DB_GetValue(t *testing.T) {
	table := "A_tables"
	// table := createInitTable()
	// defer dropTable(table)
	gtest.C(t, func(t *gtest.T) {
		value, err := db.GetValue(ctx, fmt.Sprintf("SELECT id FROM %s WHERE account_name=?", table), "map2")
		t.AssertNil(err)
		t.Assert(value.Int(), 2000)
	})
}

func Test_DB_GetCount(t *testing.T) {
	table := "A_tables"
	// table := createInitTable()
	// defer dropTable(table)
	gtest.C(t, func(t *gtest.T) {
		count, err := db.GetCount(ctx, fmt.Sprintf("SELECT * FROM %s", table))
		t.AssertNil(err)
		t.Assert(count, 28)
	})
}

func Test_DB_GetStruct(t *testing.T) {
	table := "A_tables"
	// table := createInitTable()
	// defer dropTable(table)
	gtest.C(t, func(t *gtest.T) {
		user := new(User)
		err := db.GetScan(ctx, user, fmt.Sprintf("SELECT * FROM %s WHERE id=?", table), 3)
		t.AssertNil(err)
		t.Assert(user.AccountName, "name_3")
	})
	gtest.C(t, func(t *gtest.T) {
		user := new(User)
		err := db.GetScan(ctx, user, fmt.Sprintf("SELECT * FROM %s WHERE id=?", table), 200)
		t.AssertNil(err)
		t.Assert(user.AccountName, "A200two")
	})
}

func Test_DB_GetStructs(t *testing.T) {
	table := "A_tables"
	// table := createInitTable()
	// defer dropTable(table)
	gtest.C(t, func(t *gtest.T) {
		var users []User
		err := db.GetScan(ctx, &users, fmt.Sprintf("SELECT * FROM %s WHERE id>?", table), 4000)
		t.AssertNil(err)
		t.Assert(users[0].ID, 5000)
		t.Assert(users[1].ID, 6000)
		t.Assert(users[2].ID, 6001)
		t.Assert(users[0].AccountName, "struct_5")
		t.Assert(users[1].AccountName, "t6000")
		t.Assert(users[2].AccountName, "t6001")
	})
}

func Test_DB_GetScan(t *testing.T) {
	table := "A_tables"
	// table := createInitTable()
	// defer dropTable(table)
	gtest.C(t, func(t *gtest.T) {
		user := new(User)
		err := db.GetScan(ctx, user, fmt.Sprintf("SELECT * FROM %s WHERE id=?", table), 3)
		t.AssertNil(err)
		t.Assert(user.AccountName, "name_3")
	})
	gtest.C(t, func(t *gtest.T) {
		var user *User
		err := db.GetScan(ctx, &user, fmt.Sprintf("SELECT * FROM %s WHERE id=?", table), 3)
		t.AssertNil(err)
		t.Assert(user.AccountName, "name_3")
	})
	gtest.C(t, func(t *gtest.T) {
		var users []User
		err := db.GetScan(ctx, &users, fmt.Sprintf("SELECT * FROM %s WHERE id<?", table), 4)
		t.AssertNil(err)
		t.Assert(users[0].ID, 1)
		t.Assert(users[1].ID, 2)
		t.Assert(users[2].ID, 3)
		t.Assert(users[0].AccountName, "name_1")
		t.Assert(users[1].AccountName, "name_2")
		t.Assert(users[2].AccountName, "name_3")
	})
}

func Test_DB_Delete(t *testing.T) {
	// table := createInitTable()
	// defer dropTable(table)
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Delete(ctx, "A_tables", "id=32")
		t.AssertNil(err)
		n, _ := result.RowsAffected()
		t.Assert(n, 0)
	})

	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model("A_tables").Where("id", 33).Delete()
		t.AssertNil(err)
		n, _ := result.RowsAffected()
		t.Assert(n, 0)
	})
}

func Test_Empty_Slice_Argument(t *testing.T) {
	table := "A_tables"
	// table := createInitTable()
	// defer dropTable(table)
	gtest.C(t, func(t *gtest.T) {
		result, err := db.GetAll(ctx, fmt.Sprintf(`select * from %s where id in(?)`, table), g.Slice{})
		t.AssertNil(err)
		t.Assert(len(result), 0)
	})
}

// func Test_GROUP_CONCAT(t *testing.T) {
// 	gtest.C(t, func(t *gtest.T) {
// 		type GroupIdAndUserIDsInfo struct {
// 			GroupID int64
// 			UserIDs string
// 		}
// 		result := make([]GroupIdAndUserIDsInfo, 0)

// 		model := db.Model("t_inf_group", "groupinfo").Fields("groupinfo.id as group_id", "GROUP_CONCAT(userinfo.id) as user_ids")
// 		model.InnerJoin("t_lin_user_group", "lin", "groupinfo.id = lin.group_id")
// 		model.InnerJoin("t_inf_user", "userinfo", "lin.user_id = userinfo.id")
// 		model.Where("groupinfo.enabled", 1).Where("groupinfo.deleted", 0)
// 		model.Where("userinfo.enabled", 1).Where("userinfo.deleted", 0)
// 		model.Group("groupinfo.id")

// 		err := model.Scan(&result)
// 		gtest.Assert(err, nil)
// 		g.Dump(result)
// 	})
// }

// func TestGroup(t *testing.T) {
// 	gtest.C(t, func(t *gtest.T) {
// 		type GroupListResult struct {
// 			ID           int64    `json:"group_id"`
// 			GroupName    string   `json:"group_name"`
// 			CategoryName string   `json:"category_name"`
// 			Description  string   `json:"description"`
// 			RoleName     string   `json:"role_name"`
// 			UserIDs      []string `json:"user_ids"`
// 			Enabled      int64    `json:"enabled"`
// 			CreatedTime  string   `json:"created_time"`
// 			UpdateTime   string   `json:"updated_time"`
// 		}
// 		result := make([]GroupListResult, 0)

// 		model := db.Model("t_inf_group", "groupinfo")
// 		model.LeftJoin("t_inf_group_category", "category", "groupinfo.category_id=category.id and (category.enabled = 1) and (category.deleted = 0)").
// 			Where("groupinfo.deleted", 0).
// 			Where("groupinfo.enabled", 1)

// 		total, err := model.Count()
// 		gtest.Assert(err, nil)
// 		model.Fields("distinct groupinfo.id, groupinfo.group_name, groupinfo.enabled, ifnull(category.category_name,'') as category_name", "groupinfo.created_time", "groupinfo.updated_time", "groupinfo.description")
// 		err = model.Order("groupinfo.updated_time desc").Page(1, 100).Scan(&result)
// 		gtest.Assert(err, nil)
// 		g.Dump(result)
// 		g.Dump(total)
// 	})

// 	gtest.C(t, func(t *gtest.T) {
// 		type GroupListByUserIdResult struct {
// 			ID           int64  `json:"group_id"`
// 			GroupName    string `json:"group_name"`
// 			CategoryName string `json:"category_name"`
// 			RoleName     string `json:"role_name"`
// 		}
// 		result := make([]*GroupListByUserIdResult, 0)

// 		model := db.Model("t_inf_group", "groupinfo").Fields("distinct groupinfo.id, groupinfo.group_name, groupinfo.enabled, category.category_name,groupinfo.updated_time")

// 		model.LeftJoin("t_inf_group_category", "category", "groupinfo.category_id=category.id and (category.enabled = 1) and (category.deleted = 0)")

// 		// if userId != 0 {
// 		// 	model.InnerJoin(grouptype.TLINUSERGROUP, "ug", "groupinfo.id=ug.group_id")
// 		// 	model.InnerJoin(grouptype.TINFUSER, "u", "u.id=ug.user_id")
// 		// 	model.Where("u.id = ?", userId).Where("u.deleted", consts.DataDeletedFalse)
// 		// }

// 		model.Where("groupinfo.enabled", 1).Where("groupinfo.deleted", 0)

// 		err := model.Order("groupinfo.updated_time desc").Scan(&result)
// 		//
// 		gtest.Assert(err, nil)
// 		g.Dump(result)
// 	})
// 	gtest.C(t, func(t *gtest.T) {

// 		model := db.Model("t_inf_role", "role").Fields("role.role_name", "role.id")
// 		model.RightJoin("t_lin_group_role", "link", "link.role_id=role.id")
// 		// model.Where("link.group_id", gid)
// 		model.Where("role.deleted", 0)

// 		record, err := model.One()
// 		gtest.Assert(err, nil)
// 		g.Dump(record)

// 	})
// 	gtest.C(t, func(t *gtest.T) {
// 		type GroupInfos struct {
// 			RoleName string `orm:"role_name"`
// 			RoleID   int64  `orm:"id"`
// 			GroupID  int64  `orm:"group_id"`
// 		}
// 		result := make([]GroupInfos, 0)

// 		model := db.Model("t_inf_role", "role").Fields("role.id", "role.role_name", "link.group_id")
// 		model.RightJoin("t_lin_group_role", "link", "link.role_id=role.id")
// 		model.Where("role.enabled", 1).Where("role.deleted", 0)
// 		err := model.Scan(&result)
// 		gtest.Assert(err, nil)
// 		g.Dump(result)
// 	})
// 	gtest.C(t, func(t *gtest.T) {
// 		type GroupIdAndRoleNameInfo struct {
// 			GroupID  int64
// 			RoleID   int64
// 			RoleName string
// 		}
// 		result := make([]GroupIdAndRoleNameInfo, 0)
// 		model := db.Model("t_inf_group", "groupinfo").Fields("groupinfo.id as group_id", "lin.role_id", "role.role_name")
// 		model.InnerJoin("t_lin_group_role", "lin", "groupinfo.id = lin.group_id")
// 		model.InnerJoin("t_inf_role", "role", "lin.role_id = role.id")
// 		model.Where("groupinfo.enabled", 1).Where("groupinfo.deleted", 0)
// 		model.Where("role.enabled", 1).Where("role.deleted", 0)

// 		err2 := model.Scan(&result)
// 		gtest.Assert(err2, nil)
// 		g.Dump(result)
// 	})
// }
