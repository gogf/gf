// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package pgsql_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/test/gtest"
)

func Test_DB_Query(t *testing.T) {
	table := createTable("name")
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		_, err := db.Query(ctx, fmt.Sprintf("select * from %s ", table))
		t.AssertNil(err)
	})
}

func Test_DB_Exec(t *testing.T) {
	table := createTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		_, err := db.Exec(ctx, fmt.Sprintf("select * from %s ", table))
		t.AssertNil(err)
	})
}

func Test_DB_Insert(t *testing.T) {
	table := createTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		_, err := db.Insert(ctx, table, g.Map{
			"id":          1,
			"passport":    "t1",
			"password":    "25d55ad283aa400af464c76d713c07ad",
			"nickname":    "T1",
			"create_time": gtime.Now().String(),
		})
		t.AssertNil(err)
		answer, err := db.GetAll(ctx, fmt.Sprintf("SELECT * FROM %s WHERE id=?", table), 1)
		t.AssertNil(err)
		t.Assert(len(answer), 1)
		t.Assert(answer[0]["passport"], "t1")
		t.Assert(answer[0]["password"], "25d55ad283aa400af464c76d713c07ad")
		t.Assert(answer[0]["nickname"], "T1")

		// normal map
		result, err := db.Insert(ctx, table, g.Map{
			"id":          "2",
			"passport":    "t2",
			"password":    "25d55ad283aa400af464c76d713c07ad",
			"nickname":    "name_2",
			"create_time": gtime.Now().String(),
		})
		t.AssertNil(err)
		n, _ := result.RowsAffected()
		t.Assert(n, 1)
		answer, err = db.GetAll(ctx, fmt.Sprintf("SELECT * FROM %s WHERE id=?", table), 2)
		t.AssertNil(err)
		t.Assert(len(answer), 1)
		t.Assert(answer[0]["passport"], "t2")
		t.Assert(answer[0]["password"], "25d55ad283aa400af464c76d713c07ad")
		t.Assert(answer[0]["nickname"], "name_2")
	})
}

func Test_DB_Save(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		createTable("t_user")
		defer dropTable("t_user")

		i := 10
		data := g.Map{
			"id":          i,
			"passport":    fmt.Sprintf(`t%d`, i),
			"password":    fmt.Sprintf(`p%d`, i),
			"nickname":    fmt.Sprintf(`T%d`, i),
			"create_time": gtime.Now().String(),
		}
		_, err := db.Save(ctx, "t_user", data, 10)
		gtest.AssertNE(err, nil)
	})
}

func Test_DB_Replace(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		createTable("t_user")
		defer dropTable("t_user")

		i := 10
		data := g.Map{
			"id":          i,
			"passport":    fmt.Sprintf(`t%d`, i),
			"password":    fmt.Sprintf(`p%d`, i),
			"nickname":    fmt.Sprintf(`T%d`, i),
			"create_time": gtime.Now().String(),
		}
		_, err := db.Replace(ctx, "t_user", data, 10)
		gtest.AssertNE(err, nil)
	})
}

func Test_DB_GetAll(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		result, err := db.GetAll(ctx, fmt.Sprintf("SELECT * FROM %s WHERE id=?", table), 1)
		t.AssertNil(err)
		t.Assert(len(result), 1)
		t.Assert(result[0]["id"].Int(), 1)
	})
	gtest.C(t, func(t *gtest.T) {
		result, err := db.GetAll(ctx, fmt.Sprintf("SELECT * FROM %s WHERE id=?", table), g.Slice{1})
		t.AssertNil(err)
		t.Assert(len(result), 1)
		t.Assert(result[0]["id"].Int(), 1)
	})
	gtest.C(t, func(t *gtest.T) {
		result, err := db.GetAll(ctx, fmt.Sprintf("SELECT * FROM %s WHERE id in(?)", table), g.Slice{1, 2, 3})
		t.AssertNil(err)
		t.Assert(len(result), 3)
		t.Assert(result[0]["id"].Int(), 1)
		t.Assert(result[1]["id"].Int(), 2)
		t.Assert(result[2]["id"].Int(), 3)
	})
	gtest.C(t, func(t *gtest.T) {
		result, err := db.GetAll(ctx, fmt.Sprintf("SELECT * FROM %s WHERE id in(?,?,?)", table), g.Slice{1, 2, 3})
		t.AssertNil(err)
		t.Assert(len(result), 3)
		t.Assert(result[0]["id"].Int(), 1)
		t.Assert(result[1]["id"].Int(), 2)
		t.Assert(result[2]["id"].Int(), 3)
	})
	gtest.C(t, func(t *gtest.T) {
		result, err := db.GetAll(ctx, fmt.Sprintf("SELECT * FROM %s WHERE id in(?,?,?)", table), g.Slice{1, 2, 3}...)
		t.AssertNil(err)
		t.Assert(len(result), 3)
		t.Assert(result[0]["id"].Int(), 1)
		t.Assert(result[1]["id"].Int(), 2)
		t.Assert(result[2]["id"].Int(), 3)
	})
	gtest.C(t, func(t *gtest.T) {
		result, err := db.GetAll(ctx, fmt.Sprintf("SELECT * FROM %s WHERE id>=? AND id <=?", table), g.Slice{1, 3})
		t.AssertNil(err)
		t.Assert(len(result), 3)
		t.Assert(result[0]["id"].Int(), 1)
		t.Assert(result[1]["id"].Int(), 2)
		t.Assert(result[2]["id"].Int(), 3)
	})
}

func Test_DB_GetOne(t *testing.T) {
	table := createTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		type User struct {
			Id         int
			Passport   string
			Password   string
			Nickname   string
			CreateTime string
		}
		data := User{
			Id:         1,
			Passport:   "user_1",
			Password:   "pass_1",
			Nickname:   "name_1",
			CreateTime: "2020-10-10 12:00:01",
		}
		_, err := db.Insert(ctx, table, data)
		t.AssertNil(err)

		one, err := db.GetOne(ctx, fmt.Sprintf("SELECT * FROM %s WHERE id=?", table), 1)
		t.AssertNil(err)
		t.Assert(one["passport"], data.Passport)
		t.Assert(one["create_time"], data.CreateTime)
		t.Assert(one["nickname"], data.Nickname)
	})
}

func Test_DB_GetValue(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)
	gtest.C(t, func(t *gtest.T) {
		value, err := db.GetValue(ctx, fmt.Sprintf("SELECT id FROM %s WHERE passport=?", table), "user_3")
		t.AssertNil(err)
		t.Assert(value.Int(), 3)
	})
}

func Test_DB_GetCount(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)
	gtest.C(t, func(t *gtest.T) {
		count, err := db.GetCount(ctx, fmt.Sprintf("SELECT * FROM %s", table))
		t.AssertNil(err)
		t.Assert(count, TableSize)
	})
}

func Test_DB_GetArray(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)
	gtest.C(t, func(t *gtest.T) {
		array, err := db.GetArray(ctx, fmt.Sprintf("SELECT password FROM %s", table))
		t.AssertNil(err)
		arrays := make([]string, 0)
		for i := 1; i <= TableSize; i++ {
			arrays = append(arrays, fmt.Sprintf(`pass_%d`, i))
		}
		t.Assert(array, arrays)
	})
}

func Test_DB_GetScan(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)
	gtest.C(t, func(t *gtest.T) {
		type User struct {
			Id         int
			Passport   string
			Password   string
			NickName   string
			CreateTime gtime.Time
		}
		user := new(User)
		err := db.GetScan(ctx, user, fmt.Sprintf("SELECT * FROM %s WHERE id=?", table), 3)
		t.AssertNil(err)
		t.Assert(user.NickName, "name_3")
	})
}

func Test_DB_Update(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		result, err := db.Update(ctx, table, "password='987654321'", "id=3")
		t.AssertNil(err)
		n, _ := result.RowsAffected()
		t.Assert(n, 1)

		one, err := db.Model(table).Where("id", 3).One()
		t.AssertNil(err)
		t.Assert(one["id"].Int(), 3)
		t.Assert(one["passport"].String(), "user_3")
		t.Assert(one["password"].String(), "987654321")
		t.Assert(one["nickname"].String(), "name_3")
	})
}

func Test_DB_Delete(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		result, err := db.Delete(ctx, table, "id>3")
		t.AssertNil(err)
		n, _ := result.RowsAffected()
		t.Assert(n, 7)
	})
}

func Test_DB_Tables(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		tables := []string{"t_user1", "pop", "haha"}
		for _, v := range tables {
			createTable(v)
		}
		result, err := db.Tables(ctx)
		gtest.Assert(err, nil)
		for i := 0; i < len(tables); i++ {
			find := false
			for j := 0; j < len(result); j++ {
				if tables[i] == result[j] {
					find = true
					break
				}
			}
			gtest.AssertEQ(find, true)
		}
	})
}

func Test_DB_TableFields(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		table := createTable()
		defer dropTable(table)

		var expect = map[string][]interface{}{
			//[]string: Index Type Null Key Default Comment
			//id is bigserial so the default is a pgsql function
			"id":          {0, "int8", false, "pri", fmt.Sprintf("nextval('%s_id_seq'::regclass)", table), ""},
			"passport":    {1, "varchar", false, "", nil, ""},
			"password":    {2, "varchar", false, "", nil, ""},
			"nickname":    {3, "varchar", false, "", nil, ""},
			"create_time": {4, "timestamp", false, "", nil, ""},
		}

		res, err := db.TableFields(ctx, table)
		gtest.Assert(err, nil)

		for k, v := range expect {
			_, ok := res[k]
			gtest.AssertEQ(ok, true)

			gtest.AssertEQ(res[k].Index, v[0])
			gtest.AssertEQ(res[k].Name, k)
			gtest.AssertEQ(res[k].Type, v[1])
			gtest.AssertEQ(res[k].Null, v[2])
			gtest.AssertEQ(res[k].Key, v[3])
			gtest.AssertEQ(res[k].Default, v[4])
			gtest.AssertEQ(res[k].Comment, v[5])
		}
	})
}

func Test_NoFields_Error(t *testing.T) {
	createSql := `CREATE TABLE IF NOT EXISTS %s (
id bigint PRIMARY KEY,
int_col INT);`

	type Data struct {
		Id     int64
		IntCol int64
	}
	// pgsql converts table names to lowercase
	tableName := "Error_table"
	errStr := fmt.Sprintf(`The table "%s" may not exist, or the table contains no fields`, tableName)
	_, err := db.Exec(ctx, fmt.Sprintf(createSql, tableName))
	gtest.AssertNil(err)
	defer dropTable(tableName)

	gtest.C(t, func(t *gtest.T) {
		var data = Data{
			Id:     2,
			IntCol: 2,
		}
		_, err = db.Model(tableName).Data(data).Insert()
		t.Assert(err, errStr)

		// Insert a piece of test data using lowercase
		_, err = db.Model(strings.ToLower(tableName)).Data(data).Insert()
		t.AssertNil(err)

		_, err = db.Model(tableName).Where("id", 1).Data(g.Map{
			"int_col": 9999,
		}).Update()
		t.Assert(err, errStr)

	})
	// The inserted field does not exist in the table
	gtest.C(t, func(t *gtest.T) {
		data := map[string]any{
			"id1":        22,
			"int_col_22": 11111,
		}
		_, err = db.Model(tableName).Data(data).Insert()
		t.Assert(err, errStr)

		lowerTableName := strings.ToLower(tableName)
		_, err = db.Model(lowerTableName).Data(data).Insert()
		t.Assert(err, fmt.Errorf(`input data match no fields in table "%s"`, lowerTableName))

		_, err = db.Model(lowerTableName).Where("id", 1).Data(g.Map{
			"int_col-2": 9999,
		}).Update()
		t.Assert(err, fmt.Errorf(`input data match no fields in table "%s"`, lowerTableName))
	})

}
