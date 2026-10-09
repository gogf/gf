// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package sqlitecgo_test

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/test/gtest"
)

func Test_Model_Value_WithCache(t *testing.T) {
	table := createTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		value, err := db.Model(table).Where("id", 1).Cache(gdb.CacheOption{
			Duration: time.Second * 10,
			Force:    false,
		}).Value()
		t.AssertNil(err)
		t.Assert(value.Int(), 0)
	})

	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Data(g.MapStrAny{
			"id":       1,
			"passport": fmt.Sprintf(`passport_%d`, 1),
			"password": fmt.Sprintf(`password_%d`, 1),
			"nickname": fmt.Sprintf(`nickname_%d`, 1),
		}).Insert()
		t.AssertNil(err)
		n, _ := result.RowsAffected()
		t.Assert(n, 1)
	})
	gtest.C(t, func(t *gtest.T) {
		value, err := db.Model(table).Where("id", 1).Cache(gdb.CacheOption{
			Duration: time.Second * 10,
			Force:    false,
		}).Value("id")
		t.AssertNil(err)
		t.Assert(value.Int(), 1)
	})
}

func Test_Model_Count_WithCache(t *testing.T) {
	table := createTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		count, err := db.Model(table).Where("id", 1).Cache(gdb.CacheOption{
			Duration: time.Second * 10,
			Force:    false,
		}).Count()
		t.AssertNil(err)
		t.Assert(count, int64(0))
	})
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Data(g.MapStrAny{
			"id":       1,
			"passport": fmt.Sprintf(`passport_%d`, 1),
			"password": fmt.Sprintf(`password_%d`, 1),
			"nickname": fmt.Sprintf(`nickname_%d`, 1),
		}).Insert()
		t.AssertNil(err)
		n, _ := result.RowsAffected()
		t.Assert(n, 1)
	})
	gtest.C(t, func(t *gtest.T) {
		count, err := db.Model(table).Where("id", 1).Cache(gdb.CacheOption{
			Duration: time.Second * 10,
			Force:    false,
		}).Count()
		t.AssertNil(err)
		t.Assert(count, int64(1))
	})
}

func Test_Model_Count_All_WithCache(t *testing.T) {
	table := createTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		count, err := db.Model(table).Cache(gdb.CacheOption{
			Duration: time.Second * 10,
			Force:    false,
		}).Count()
		t.AssertNil(err)
		t.Assert(count, int64(0))
	})
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Data(g.MapStrAny{
			"id":       1,
			"passport": fmt.Sprintf(`passport_%d`, 1),
			"password": fmt.Sprintf(`password_%d`, 1),
			"nickname": fmt.Sprintf(`nickname_%d`, 1),
		}).Insert()
		t.AssertNil(err)
		n, _ := result.RowsAffected()
		t.Assert(n, 1)
	})
	gtest.C(t, func(t *gtest.T) {
		count, err := db.Model(table).Cache(gdb.CacheOption{
			Duration: time.Second * 10,
			Force:    false,
		}).Count()
		t.AssertNil(err)
		t.Assert(count, int64(1))
	})
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Data(g.MapStrAny{
			"id":       2,
			"passport": fmt.Sprintf(`passport_%d`, 2),
			"password": fmt.Sprintf(`password_%d`, 2),
			"nickname": fmt.Sprintf(`nickname_%d`, 2),
		}).Insert()
		t.AssertNil(err)
		n, _ := result.RowsAffected()
		t.Assert(n, 1)
	})
	gtest.C(t, func(t *gtest.T) {
		count, err := db.Model(table).Cache(gdb.CacheOption{
			Duration: time.Second * 10,
			Force:    false,
		}).Count()
		t.AssertNil(err)
		t.Assert(count, int64(1))
	})
}

func Test_Model_CountColumn_WithCache(t *testing.T) {
	table := createTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		count, err := db.Model(table).Where("id", 1).Cache(gdb.CacheOption{
			Duration: time.Second * 10,
			Force:    false,
		}).CountColumn("id")
		t.AssertNil(err)
		t.Assert(count, int64(0))
	})
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Data(g.MapStrAny{
			"id":       1,
			"passport": fmt.Sprintf(`passport_%d`, 1),
			"password": fmt.Sprintf(`password_%d`, 1),
			"nickname": fmt.Sprintf(`nickname_%d`, 1),
		}).Insert()
		t.AssertNil(err)
		n, _ := result.RowsAffected()
		t.Assert(n, 1)
	})
	gtest.C(t, func(t *gtest.T) {
		count, err := db.Model(table).Where("id", 1).Cache(gdb.CacheOption{
			Duration: time.Second * 10,
			Force:    false,
		}).CountColumn("id")
		t.AssertNil(err)
		t.Assert(count, int64(1))
	})
}

func Test_Model_Option_List(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		table := createTable()
		defer dropTable(table)
		r, err := db.Model(table).Fields("id, password").Data(g.List{
			g.Map{
				"id":       1,
				"passport": "1",
				"password": "1",
				"nickname": "1",
			},
			g.Map{
				"id":       2,
				"passport": "2",
				"password": "2",
				"nickname": "2",
			},
		}).Save()
		t.AssertNil(err)
		n, _ := r.RowsAffected()
		t.Assert(n, 2)
		list, err := db.Model(table).Order("id asc").All()
		t.AssertNil(err)
		t.Assert(len(list), 2)
		t.Assert(list[0]["id"].String(), "1")
		t.Assert(list[0]["nickname"].String(), "")
		t.Assert(list[0]["passport"].String(), "passport")
		t.Assert(list[0]["password"].String(), "1")

		t.Assert(list[1]["id"].String(), "2")
		t.Assert(list[1]["nickname"].String(), "")
		t.Assert(list[1]["passport"].String(), "passport")
		t.Assert(list[1]["password"].String(), "2")
	})

	gtest.C(t, func(t *gtest.T) {
		table := createTable()
		defer dropTable(table)
		r, err := db.Model(table).OmitEmpty().Fields("id, password").Data(g.List{
			g.Map{
				"id":       1,
				"passport": "1",
				"password": 0,
				"nickname": "1",
			},
			g.Map{
				"id":       2,
				"passport": "2",
				"password": "2",
				"nickname": "2",
			},
		}).Save()
		t.AssertNil(err)
		n, _ := r.RowsAffected()
		t.Assert(n, 2)
		list, err := db.Model(table).Order("id asc").All()
		t.AssertNil(err)
		t.Assert(len(list), 2)
		t.Assert(list[0]["id"].String(), "1")
		t.Assert(list[0]["nickname"].String(), "")
		t.Assert(list[0]["passport"].String(), "passport")
		t.Assert(list[0]["password"].String(), "0")

		t.Assert(list[1]["id"].String(), "2")
		t.Assert(list[1]["nickname"].String(), "")
		t.Assert(list[1]["passport"].String(), "passport")
		t.Assert(list[1]["password"].String(), "2")
	})
}

func Test_Model_OmitEmpty(t *testing.T) {
	table := fmt.Sprintf(`table_%s`, gtime.TimestampNanoStr())
	if _, err := db.Exec(ctx, fmt.Sprintf(`
    CREATE TABLE IF NOT EXISTS %s (
        id   INTEGER PRIMARY KEY AUTOINCREMENT,
        name VARCHAR(45) NOT NULL
    );
    `, table)); err != nil {
		gtest.Error(err)
	}
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(table).OmitEmpty().Data(g.Map{
			"id":   1,
			"name": "",
		}).Save()
		t.AssertNE(err, nil)
	})
	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(table).OmitEmptyData().Data(g.Map{
			"id":   1,
			"name": "",
		}).Save()
		t.AssertNE(err, nil)
	})
	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(table).OmitEmptyWhere().Data(g.Map{
			"id":   1,
			"name": "",
		}).Save()
		t.AssertNil(err)
	})
}

func Test_Model_OmitNil(t *testing.T) {
	table := fmt.Sprintf(`table_%s`, gtime.TimestampNanoStr())
	if _, err := db.Exec(ctx, fmt.Sprintf(`
    CREATE TABLE IF NOT EXISTS %s (
        id   INTEGER PRIMARY KEY AUTOINCREMENT,
        name VARCHAR(45) NOT NULL
    );
    `, table)); err != nil {
		gtest.Error(err)
	}
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(table).OmitNil().Data(g.Map{
			"id":   1,
			"name": nil,
		}).Save()
		t.AssertNE(err, nil)
	})
	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(table).OmitNil().Data(g.Map{
			"id":   1,
			"name": "",
		}).Save()
		t.AssertNil(err)
	})
	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(table).OmitNilWhere().Data(g.Map{
			"id":   1,
			"name": "",
		}).Save()
		t.AssertNil(err)
	})
}

func Test_Model_FieldsEx_WithReservedWords(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var table = "fieldsex_test_table"
		dropTable(table)
		if _, err := db.Exec(ctx, fmt.Sprintf(`
CREATE TABLE %s (
  `+"`id`"+`          INTEGER PRIMARY KEY AUTOINCREMENT,
  `+"`key`"+`         VARCHAR(45) DEFAULT NULL,
  `+"`category_id`"+` INTEGER NOT NULL,
  `+"`user_id`"+`     INTEGER NOT NULL,
  `+"`title`"+`       VARCHAR(255) NOT NULL,
  `+"`content`"+`     TEXT NOT NULL,
  `+"`sort`"+`        INTEGER DEFAULT 0,
  `+"`brief`"+`       VARCHAR(255) DEFAULT NULL,
  `+"`thumb`"+`       VARCHAR(255) DEFAULT NULL,
  `+"`tags`"+`        VARCHAR(900) DEFAULT NULL,
  `+"`referer`"+`     VARCHAR(255) DEFAULT NULL,
  `+"`status`"+`      SMALLINT DEFAULT 0,
  `+"`view_count`"+`  INTEGER DEFAULT 0,
  `+"`zan_count`"+`   INTEGER DEFAULT NULL,
  `+"`cai_count`"+`   INTEGER DEFAULT NULL,
  `+"`created_at`"+`  DATETIME DEFAULT NULL,
  `+"`updated_at`"+`  DATETIME DEFAULT NULL
);
`, table)); err != nil {
			t.AssertNil(err)
		}
		defer dropTable(table)
		_, err := db.Model(table).FieldsEx("content").One()
		t.AssertNil(err)
	})
}

func Test_Model_NullField(t *testing.T) {
	table := fmt.Sprintf(`table_%s`, gtime.TimestampNanoStr())
	if _, err := db.Exec(ctx, fmt.Sprintf(`
    CREATE TABLE IF NOT EXISTS %s (
        id       INTEGER PRIMARY KEY AUTOINCREMENT,
        passport VARCHAR(45)
    );
    `, table)); err != nil {
		gtest.Error(err)
	}
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		type User struct {
			Id       int
			Passport *string
		}
		data := g.Map{
			"id":       1,
			"passport": nil,
		}
		result, err := db.Model(table).Data(data).Insert()
		t.AssertNil(err)
		n, _ := result.RowsAffected()
		t.Assert(n, 1)
		one, err := db.Model(table).WherePri(1).One()
		t.AssertNil(err)

		var user *User
		err = one.Struct(&user)
		t.AssertNil(err)
		t.Assert(user.Id, data["id"])
		t.Assert(user.Passport, data["passport"])
	})
}

func Test_Model_OnDuplicate(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	// string type 1.
	gtest.C(t, func(t *gtest.T) {
		data := g.Map{
			"id":          1,
			"passport":    "pp1",
			"password":    "pw1",
			"nickname":    "n1",
			"create_time": "2016-06-06",
		}
		_, err := db.Model(table).OnDuplicate("passport,password").Data(data).Save()
		t.AssertNil(err)
		one, err := db.Model(table).WherePri(1).One()
		t.AssertNil(err)
		t.Assert(one["passport"], data["passport"])
		t.Assert(one["password"], data["password"])
		t.Assert(one["nickname"], "name_1")
	})

	// string type 2.
	gtest.C(t, func(t *gtest.T) {
		data := g.Map{
			"id":          1,
			"passport":    "pp1",
			"password":    "pw1",
			"nickname":    "n1",
			"create_time": "2016-06-06",
		}
		_, err := db.Model(table).OnDuplicate("passport", "password").Data(data).Save()
		t.AssertNil(err)
		one, err := db.Model(table).WherePri(1).One()
		t.AssertNil(err)
		t.Assert(one["passport"], data["passport"])
		t.Assert(one["password"], data["password"])
		t.Assert(one["nickname"], "name_1")
	})

	// slice.
	gtest.C(t, func(t *gtest.T) {
		data := g.Map{
			"id":          1,
			"passport":    "pp1",
			"password":    "pw1",
			"nickname":    "n1",
			"create_time": "2016-06-06",
		}
		_, err := db.Model(table).OnDuplicate(g.Slice{"passport", "password"}).Data(data).Save()
		t.AssertNil(err)
		one, err := db.Model(table).WherePri(1).One()
		t.AssertNil(err)
		t.Assert(one["passport"], data["passport"])
		t.Assert(one["password"], data["password"])
		t.Assert(one["nickname"], "name_1")
	})

	// map.
	gtest.C(t, func(t *gtest.T) {
		data := g.Map{
			"id":          1,
			"passport":    "pp1",
			"password":    "pw1",
			"nickname":    "n1",
			"create_time": "2016-06-06",
		}
		_, err := db.Model(table).OnDuplicate(g.Map{
			"passport": "nickname",
			"password": "nickname",
		}).Data(data).Save()
		t.AssertNil(err)
		one, err := db.Model(table).WherePri(1).One()
		t.AssertNil(err)
		t.Assert(one["passport"], data["nickname"])
		t.Assert(one["password"], data["nickname"])
		t.Assert(one["nickname"], "name_1")
	})

	// map+raw.
	gtest.C(t, func(t *gtest.T) {
		data := g.MapStrStr{
			"id":          "1",
			"passport":    "pp1",
			"password":    "pw1",
			"nickname":    "n1",
			"create_time": "2016-06-06",
		}
		_, err := db.Model(table).OnDuplicate(g.Map{
			"passport": gdb.Raw("EXCLUDED.`passport` || '1'"),
			"password": gdb.Raw("EXCLUDED.`password` || '2'"),
		}).Data(data).Save()
		t.AssertNil(err)
		one, err := db.Model(table).WherePri(1).One()
		t.AssertNil(err)
		t.Assert(one["passport"], data["passport"]+"1")
		t.Assert(one["password"], data["password"]+"2")
		t.Assert(one["nickname"], "name_1")
	})
}

func Test_Model_OnDuplicateEx(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	// string type 1.
	gtest.C(t, func(t *gtest.T) {
		data := g.Map{
			"id":          1,
			"passport":    "pp1",
			"password":    "pw1",
			"nickname":    "n1",
			"create_time": "2016-06-06",
		}
		_, err := db.Model(table).OnDuplicateEx("nickname,create_time").Data(data).Save()
		t.AssertNil(err)
		one, err := db.Model(table).WherePri(1).One()
		t.AssertNil(err)
		t.Assert(one["passport"], data["passport"])
		t.Assert(one["password"], data["password"])
		t.Assert(one["nickname"], "name_1")
	})

	// string type 2.
	gtest.C(t, func(t *gtest.T) {
		data := g.Map{
			"id":          1,
			"passport":    "pp1",
			"password":    "pw1",
			"nickname":    "n1",
			"create_time": "2016-06-06",
		}
		_, err := db.Model(table).OnDuplicateEx("nickname", "create_time").Data(data).Save()
		t.AssertNil(err)
		one, err := db.Model(table).WherePri(1).One()
		t.AssertNil(err)
		t.Assert(one["passport"], data["passport"])
		t.Assert(one["password"], data["password"])
		t.Assert(one["nickname"], "name_1")
	})

	// slice.
	gtest.C(t, func(t *gtest.T) {
		data := g.Map{
			"id":          1,
			"passport":    "pp1",
			"password":    "pw1",
			"nickname":    "n1",
			"create_time": "2016-06-06",
		}
		_, err := db.Model(table).OnDuplicateEx(g.Slice{"nickname", "create_time"}).Data(data).Save()
		t.AssertNil(err)
		one, err := db.Model(table).WherePri(1).One()
		t.AssertNil(err)
		t.Assert(one["passport"], data["passport"])
		t.Assert(one["password"], data["password"])
		t.Assert(one["nickname"], "name_1")
	})

	// map.
	gtest.C(t, func(t *gtest.T) {
		data := g.Map{
			"id":          1,
			"passport":    "pp1",
			"password":    "pw1",
			"nickname":    "n1",
			"create_time": "2016-06-06",
		}
		_, err := db.Model(table).OnDuplicateEx(g.Map{
			"nickname":    "nickname",
			"create_time": "nickname",
		}).Data(data).Save()
		t.AssertNil(err)
		one, err := db.Model(table).WherePri(1).One()
		t.AssertNil(err)
		t.Assert(one["passport"], data["passport"])
		t.Assert(one["password"], data["password"])
		t.Assert(one["nickname"], "name_1")
	})
}

func Test_Builder_OmitEmptyWhere(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)
	gtest.C(t, func(t *gtest.T) {
		count, err := db.Model(table).Where("id", 1).Count()
		t.AssertNil(err)
		t.Assert(count, int64(1))
	})
	gtest.C(t, func(t *gtest.T) {
		count, err := db.Model(table).Where("id", 0).OmitEmptyWhere().Count()
		t.AssertNil(err)
		t.Assert(count, int64(TableSize))
	})
	gtest.C(t, func(t *gtest.T) {
		builder := db.Model(table).OmitEmptyWhere().Builder()
		count, err := db.Model(table).Where(
			builder.Where("id", 0),
		).Count()
		t.AssertNil(err)
		t.Assert(count, int64(TableSize))
	})
}

func Test_Scan_Nil_Result_Error(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	type S struct {
		Id    int
		Name  string
		Age   int
		Score int
	}
	gtest.C(t, func(t *gtest.T) {
		var s *S
		err := db.Model(table).Where("id", 1).Scan(&s)
		t.AssertNil(err)
		t.Assert(s.Id, 1)
	})
	gtest.C(t, func(t *gtest.T) {
		var s *S
		err := db.Model(table).Where("id", 100).Scan(&s)
		t.AssertNil(err)
		t.Assert(s, nil)
	})
	gtest.C(t, func(t *gtest.T) {
		var s S
		err := db.Model(table).Where("id", 100).Scan(&s)
		t.Assert(err, sql.ErrNoRows)
	})
	gtest.C(t, func(t *gtest.T) {
		var ss []*S
		err := db.Model(table).Scan(&ss)
		t.AssertNil(err)
		t.Assert(len(ss), TableSize)
	})
	// If the result is empty, it returns error.
	gtest.C(t, func(t *gtest.T) {
		var ss = make([]*S, 10)
		err := db.Model(table).WhereGT("id", 100).Scan(&ss)
		t.Assert(err, sql.ErrNoRows)
	})
}

func Test_Model_FixGdbJoin(t *testing.T) {
	for _, v := range fixGdbJoinTables {
		if _, err := db.Exec(ctx, v); err != nil {
			gtest.Error(err)
		}
	}
	defer dropTable(`common_resource`)
	defer dropTable(`managed_resource`)
	defer dropTable(`rules_template`)
	defer dropTable(`resource_mark`)
	gtest.C(t, func(t *gtest.T) {
		t.AssertNil(db.GetCore().ClearCacheAll(ctx))
		sqlSlice, err := gdb.CatchSQL(ctx, func(ctx context.Context) error {
			orm := db.Model(`managed_resource`).Ctx(ctx).
				LeftJoinOnField(`common_resource`, `resource_id`).
				LeftJoinOnFields(`resource_mark`, `resource_mark_id`, `=`, `id`).
				LeftJoinOnFields(`rules_template`, `rule_template_id`, `=`, `template_id`).
				FieldsPrefix(
					`managed_resource`,
					"resource_id", "user", "status", "status_message", "safe_publication", "rule_template_id",
					"created_at", "comments", "expired_at", "resource_mark_id", "instance_id", "resource_name",
					"pay_mode").
				FieldsPrefix(`resource_mark`, "mark_name", "color").
				FieldsPrefix(`rules_template`, "name").
				FieldsPrefix(`common_resource`, `src_instance_id`, "database_kind", "source_type", "ip", "port")
			all, err := orm.OrderAsc("src_instance_id").All()
			t.Assert(err, nil)
			t.Assert(len(all), 4)
			t.Assert(all[0]["pay_mode"], 1)
			t.Assert(all[0]["src_instance_id"], 2)
			t.Assert(all[3]["instance_id"], "dmcins-jxy0x75m")
			t.Assert(all[3]["src_instance_id"], "vdb-6b6m3u1u")
			t.Assert(all[3]["resource_mark_id"], "11")
			return err
		})
		t.AssertNil(err)

		t.Assert(fixGdbJoinExpectSQL, sqlSlice[len(sqlSlice)-1])
	})
}

func Test_Model_Year_Date_Time_DateTime_Timestamp(t *testing.T) {
	table := "date_time_example"
	dropTable(table)
	if _, err := db.Exec(ctx, `
CREATE TABLE date_time_example (
    id        INTEGER PRIMARY KEY AUTOINCREMENT,
    year      year      DEFAULT NULL,
    date      date      DEFAULT NULL,
    time      time      DEFAULT NULL,
    datetime  datetime  DEFAULT NULL,
    timestamp timestamp DEFAULT NULL
);
`); err != nil {
		gtest.Error(err)
	}
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		// insert.
		var now = gtime.Now()
		_, err := db.Model("date_time_example").Insert(g.Map{
			"year":      now,
			"date":      now,
			"time":      now,
			"datetime":  now,
			"timestamp": now,
		})
		t.AssertNil(err)
		// select.
		one, err := db.Model("date_time_example").One()
		t.AssertNil(err)
		t.Assert(one["year"].String(), now.Format("Y"))
		t.Assert(one["date"].String(), now.Format("Y-m-d"))
		t.Assert(one["time"].String(), now.Format("H:i:s"))
		t.AssertLT(one["datetime"].GTime().Sub(now).Seconds(), 5)
		t.AssertLT(one["timestamp"].GTime().Sub(now).Seconds(), 5)
	})
}

func Test_OrderBy_Statement_Generated(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		dropTable(`employee`)
		for _, v := range []string{
			`CREATE TABLE IF NOT EXISTS employee (
				id   INTEGER PRIMARY KEY AUTOINCREMENT,
				name VARCHAR(255) NOT NULL,
				age  INTEGER       NOT NULL
			)`,
			`INSERT INTO employee(name, age) VALUES ('John', 30)`,
			`INSERT INTO employee(name, age) VALUES ('Mary', 28)`,
		} {
			if _, err := db.Exec(ctx, v); err != nil {
				gtest.Error(err)
			}
		}
		defer dropTable(`employee`)
		sqlArray, _ := gdb.CatchSQL(ctx, func(ctx context.Context) error {
			g.DB("default").Ctx(ctx).Model("employee").Order("name asc", "age desc").All()
			return nil
		})
		rawSql := strings.ReplaceAll(sqlArray[len(sqlArray)-1], " ", "")
		expectSql := strings.ReplaceAll("SELECT * FROM `employee` ORDER BY `name` asc, `age` desc", " ", "")
		t.Assert(rawSql, expectSql)
	})
}

func Test_Fields_Raw(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		table := createInitTable()
		defer dropTable(table)
		one, err := db.Model(table).Fields(gdb.Raw("1")).One()
		t.AssertNil(err)
		t.Assert(one["1"], 1)

		one, err = db.Model(table).Fields(gdb.Raw("2")).One()
		t.AssertNil(err)
		t.Assert(one["2"], 2)

		one, err = db.Model(table).Fields(gdb.Raw("2")).Where("id", 2).One()
		t.AssertNil(err)
		t.Assert(one["2"], 2)

		one, err = db.Model(table).Fields(gdb.Raw("2")).Where("id", 10000000000).One()
		t.AssertNil(err)
		t.Assert(len(one), 0)
	})
}

var fixGdbJoinTables = []string{
	`DROP TABLE IF EXISTS ` + "`common_resource`",
	`CREATE TABLE ` + "`common_resource`" + ` (
		` + "`id`" + `              INTEGER PRIMARY KEY AUTOINCREMENT,
		` + "`app_id`" + `          BIGINT NOT NULL,
		` + "`resource_id`" + `     VARCHAR(64) NOT NULL,
		` + "`src_instance_id`" + ` VARCHAR(64) DEFAULT NULL,
		` + "`region`" + `          VARCHAR(36) DEFAULT NULL,
		` + "`zone`" + `            VARCHAR(36) DEFAULT NULL,
		` + "`database_kind`" + `   VARCHAR(20) NOT NULL,
		` + "`source_type`" + `     VARCHAR(64) NOT NULL,
		` + "`ip`" + `              VARCHAR(64) DEFAULT NULL,
		` + "`port`" + `            INTEGER DEFAULT NULL,
		` + "`vpc_id`" + `          VARCHAR(20) DEFAULT NULL,
		` + "`subnet_id`" + `       VARCHAR(20) DEFAULT NULL,
		` + "`proxy_ip`" + `        VARCHAR(64) DEFAULT NULL,
		` + "`proxy_port`" + `      INTEGER DEFAULT NULL,
		` + "`proxy_id`" + `        BIGINT DEFAULT NULL,
		` + "`proxy_snat_ip`" + `   VARCHAR(64) DEFAULT NULL,
		` + "`lease_at`" + `        DATETIME DEFAULT NULL,
		` + "`uin`" + `             VARCHAR(32) NOT NULL
	)`,
	`INSERT INTO ` + "`common_resource`" + ` VALUES
		(1,1,'2','2','2','3','1','1','1',1,'1','1','1',1,1,'1',NULL,''),
		(3,2,'3','3','3','3','3','3','3',3,'3','3','3',3,3,'3',NULL,''),
		(18,1303697168,'dmc-rgnh9qre','vdb-6b6m3u1u','ap-guangzhou','','vdb','cloud','10.0.1.16',80,'vpc-m3dchft7','subnet-9as3a3z2','9.27.72.189',11131,228476,'169.254.128.5, ','2023-11-08 08:13:04',''),
		(20,1303697168,'dmc-4grzi4jg','tdsqlshard-313spncx','ap-guangzhou','','tdsql','cloud','10.255.0.27',3306,'vpc-407k0e8x','subnet-qhkkk3bo','30.86.239.200',24087,0,'',NULL,'')`,

	`DROP TABLE IF EXISTS ` + "`managed_resource`",
	`CREATE TABLE ` + "`managed_resource`" + ` (
		` + "`id`" + `               INTEGER PRIMARY KEY AUTOINCREMENT,
		` + "`instance_id`" + `      VARCHAR(64) NOT NULL,
		` + "`resource_id`" + `      VARCHAR(64) NOT NULL,
		` + "`resource_name`" + `    VARCHAR(64) DEFAULT NULL,
		` + "`status`" + `           VARCHAR(36) NOT NULL DEFAULT 'valid',
		` + "`status_message`" + `   VARCHAR(64) DEFAULT NULL,
		` + "`user`" + `             VARCHAR(64) NOT NULL,
		` + "`password`" + `         VARCHAR(1024) NOT NULL,
		` + "`pay_mode`" + `         TINYINT DEFAULT 0,
		` + "`safe_publication`" + ` INTEGER DEFAULT 0,
		` + "`created_at`" + `       DATETIME NOT NULL DEFAULT (datetime('now')),
		` + "`updated_at`" + `       DATETIME NOT NULL DEFAULT (datetime('now')),
		` + "`expired_at`" + `       DATETIME DEFAULT NULL,
		` + "`deleted`" + `          TINYINT NOT NULL DEFAULT 0,
		` + "`resource_mark_id`" + ` INTEGER DEFAULT NULL,
		` + "`comments`" + `         VARCHAR(64) DEFAULT NULL,
		` + "`rule_template_id`" + ` VARCHAR(64) NOT NULL
	)`,
	`INSERT INTO ` + "`managed_resource`" + ` VALUES
		(1,'2','3','1','1','1','1','1',1,1,'2023-11-06 12:14:21','2023-11-06 12:14:21',NULL,1,1,'1',''),
		(2,'3','2','1','1','1','1','1',1,0,'2023-11-06 12:15:07','2023-11-06 12:15:07',NULL,1,2,'1',''),
		(5,'dmcins-jxy0x75m','dmc-rgnh9qre','erichmao-vdb-test','invalid','The Ip field is required','root','2e39af3d',1,1,'2023-11-08 08:13:20','2023-11-09 05:31:07',NULL,0,11,NULL,'12345'),
		(6,'dmcins-erxms6ya','dmc-4grzi4jg','erichmao-vdb-test','invalid','The Ip field is required','leotaowang','641d846c',1,1,'2023-11-08 22:15:17','2023-11-09 05:31:07',NULL,0,11,NULL,'12345')`,

	`DROP TABLE IF EXISTS ` + "`rules_template`",
	`CREATE TABLE ` + "`rules_template`" + ` (
		` + "`id`" + `               INTEGER PRIMARY KEY AUTOINCREMENT,
		` + "`app_id`" + `           BIGINT DEFAULT NULL,
		` + "`name`" + `             VARCHAR(255) NOT NULL,
		` + "`database_kind`" + `    VARCHAR(64) DEFAULT NULL,
		` + "`is_default`" + `       TINYINT NOT NULL DEFAULT 0,
		` + "`win_rules`" + `        VARCHAR(2048) DEFAULT NULL,
		` + "`inception_rules`" + `  VARCHAR(2048) DEFAULT NULL,
		` + "`auto_exec_rules`" + `  VARCHAR(2048) DEFAULT NULL,
		` + "`order_check_step`" + ` VARCHAR(2048) DEFAULT NULL,
		` + "`template_id`" + `      VARCHAR(64) NOT NULL DEFAULT '',
		` + "`version`" + `          INTEGER NOT NULL DEFAULT 1,
		` + "`deleted`" + `          TINYINT NOT NULL DEFAULT 0,
		` + "`create_at`" + `        DATETIME NOT NULL DEFAULT (datetime('now')),
		` + "`update_at`" + `        DATETIME NOT NULL DEFAULT (datetime('now')),
		` + "`is_system`" + `        TINYINT NOT NULL DEFAULT 0,
		` + "`uin`" + `              VARCHAR(64) DEFAULT NULL,
		` + "`subAccountUin`" + `    VARCHAR(64) DEFAULT NULL
	)`,

	`DROP TABLE IF EXISTS ` + "`resource_mark`",
	`CREATE TABLE ` + "`resource_mark`" + ` (
		` + "`id`" + `         INTEGER PRIMARY KEY AUTOINCREMENT,
		` + "`app_id`" + `     BIGINT NOT NULL,
		` + "`mark_name`" + `  VARCHAR(64) NOT NULL,
		` + "`color`" + `      VARCHAR(11) NOT NULL,
		` + "`creator`" + `    VARCHAR(32) NOT NULL,
		` + "`created_at`" + ` DATETIME NOT NULL DEFAULT (datetime('now')),
		` + "`updated_at`" + ` DATETIME NOT NULL DEFAULT (datetime('now'))
	)`,
	`INSERT INTO ` + "`resource_mark`" + ` VALUES (10,1,'test','red','1','2023-11-06 02:45:46','2023-11-06 02:45:46')`,
}

const fixGdbJoinExpectSQL = "SELECT `managed_resource`.`resource_id`,`managed_resource`.`user`," +
	"`managed_resource`.`status`,`managed_resource`.`status_message`," +
	"`managed_resource`.`safe_publication`,`managed_resource`.`rule_template_id`," +
	"`managed_resource`.`created_at`,`managed_resource`.`comments`," +
	"`managed_resource`.`expired_at`,`managed_resource`.`resource_mark_id`," +
	"`managed_resource`.`instance_id`,`managed_resource`.`resource_name`," +
	"`managed_resource`.`pay_mode`,`resource_mark`.`mark_name`,`resource_mark`.`color`," +
	"`rules_template`.`name`,`common_resource`.`src_instance_id`," +
	"`common_resource`.`database_kind`,`common_resource`.`source_type`," +
	"`common_resource`.`ip`,`common_resource`.`port` " +
	"FROM `managed_resource` " +
	"LEFT JOIN `common_resource` ON (`managed_resource`.`resource_id`=`common_resource`.`resource_id`) " +
	"LEFT JOIN `resource_mark` ON (`managed_resource`.`resource_mark_id` = `resource_mark`.`id`) " +
	"LEFT JOIN `rules_template` ON (`managed_resource`.`rule_template_id` = `rules_template`.`template_id`) " +
	"ORDER BY `src_instance_id` ASC"
