// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package oracle_test

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/test/gtest"
)

func createJSONTable(table ...string) string {
	var name string
	if len(table) > 0 {
		name = table[0]
	} else {
		name = fmt.Sprintf(`t_json_%d`, gtime.TimestampMicro()%1e9)
	}
	dropTable(name)
	if _, err := db.Exec(ctx, fmt.Sprintf(`
		CREATE TABLE %s (
			ID          NUMBER(10) NOT NULL,
			NAME        VARCHAR2(45) NULL,
			CONFIG      CLOB NULL,
			METADATA    VARCHAR2(4000) NULL,
			PRIMARY KEY (ID)
		)`, name)); err != nil {
		gtest.Fatal(err)
	}
	createAutoIncrement(name, "ID", 1)
	return name
}

func Test_JSON_Insert_Map(t *testing.T) {
	table := createJSONTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		data := g.Map{
			"name": "user1",
			"config": g.Map{
				"theme": "dark",
				"lang":  "zh-CN",
			},
			"metadata": g.Map{
				"tags":  g.Slice{"admin", "developer"},
				"level": 5,
			},
		}
		result, err := db.Model(table).Data(data).Insert()
		t.AssertNil(err)
		n, _ := result.LastInsertId()
		t.Assert(n, 1)

		one, err := db.Model(table).WherePri(1).One()
		t.AssertNil(err)
		t.Assert(one["NAME"], "user1")
		t.AssertNE(one["CONFIG"], nil)
		t.AssertNE(one["METADATA"], nil)
		t.Assert(json.Valid(one["CONFIG"].Bytes()), true)
		t.Assert(json.Valid(one["METADATA"].Bytes()), true)
		t.Assert(one["CONFIG"].Map(), g.Map{"theme": "dark", "lang": "zh-CN"})
		metadata := one["METADATA"].Map()
		t.Assert(metadata["tags"], g.Slice{"admin", "developer"})
		t.Assert(metadata["level"], 5)
	})
}

func Test_JSON_Insert_String(t *testing.T) {
	table := createJSONTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		data := g.Map{
			"name":     "user2",
			"config":   `{"theme":"light","lang":"en-US"}`,
			"metadata": `{"tags":["user"],"level":1}`,
		}
		result, err := db.Model(table).Data(data).Insert()
		t.AssertNil(err)
		n, _ := result.LastInsertId()
		t.Assert(n, 1)

		one, err := db.Model(table).WherePri(1).One()
		t.AssertNil(err)
		t.Assert(one["NAME"], "user2")
		t.AssertNE(one["CONFIG"], nil)
		t.AssertNE(one["METADATA"], nil)
		t.Assert(one["CONFIG"].String(), `{"theme":"light","lang":"en-US"}`)
		t.Assert(one["METADATA"].String(), `{"tags":["user"],"level":1}`)
	})
}

func Test_JSON_Insert_Null(t *testing.T) {
	table := createJSONTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		data := g.Map{
			"name":     "user3",
			"config":   nil,
			"metadata": nil,
		}
		result, err := db.Model(table).Data(data).Insert()
		t.AssertNil(err)
		n, _ := result.LastInsertId()
		t.Assert(n, 1)

		one, err := db.Model(table).WherePri(1).One()
		t.AssertNil(err)
		t.Assert(one["NAME"], "user3")
		t.Assert(one["CONFIG"], nil)
		t.Assert(one["METADATA"], nil)
		t.Assert(one["CONFIG"].IsNil(), true)
		t.Assert(one["METADATA"].IsNil(), true)
	})
}

func Test_JSON_Update(t *testing.T) {
	table := createJSONTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(table).Data(g.Map{
			"name": "user1",
			"config": g.Map{
				"theme": "dark",
			},
		}).Insert()
		t.AssertNil(err)

		result, err := db.Model(table).Data(g.Map{
			"config": g.Map{
				"theme": "light",
				"lang":  "en-US",
			},
		}).WherePri(1).Update()
		t.AssertNil(err)
		n, _ := result.RowsAffected()
		t.Assert(n, 1)

		one, err := db.Model(table).WherePri(1).One()
		t.AssertNil(err)
		t.AssertNE(one["CONFIG"], nil)
		t.Assert(json.Valid(one["CONFIG"].Bytes()), true)
		t.Assert(one["CONFIG"].Map(), g.Map{"theme": "light", "lang": "en-US"})
	})
}

func Test_JSON_Batch_Insert(t *testing.T) {
	table := createJSONTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		data := g.Slice{
			g.Map{
				"name": "user1",
				"config": g.Map{
					"theme": "dark",
				},
			},
			g.Map{
				"name": "user2",
				"config": g.Map{
					"theme": "light",
				},
			},
			g.Map{
				"name":   "user3",
				"config": nil,
			},
		}
		result, err := db.Model(table).Data(data).Insert()
		t.AssertNil(err)
		n, _ := result.RowsAffected()
		t.Assert(n, 3)

		all, err := db.Model(table).Order("id").All()
		t.AssertNil(err)
		t.Assert(len(all), 3)
		t.Assert(all[2]["CONFIG"].IsNil(), true)
		t.Assert(json.Valid(all[0]["CONFIG"].Bytes()), true)
		t.Assert(all[0]["CONFIG"].Map(), g.Map{"theme": "dark"})
		t.Assert(all[1]["CONFIG"].Map(), g.Map{"theme": "light"})
	})
}

func Test_JSON_Scan_To_Struct(t *testing.T) {
	table := createJSONTable()
	defer dropTable(table)

	type Config struct {
		Theme string `json:"theme"`
		Lang  string `json:"lang"`
	}
	type User struct {
		Id     int
		Name   string
		Config *Config
	}

	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(table).Data(g.Map{
			"name": "user1",
			"config": g.Map{
				"theme": "dark",
				"lang":  "zh-CN",
			},
		}).Insert()
		t.AssertNil(err)

		var user User
		err = db.Model(table).WherePri(1).Scan(&user)
		t.AssertNil(err)
		t.Assert(user.Name, "user1")
		t.AssertNE(user.Config, nil)
		if user.Config != nil {
			t.Assert(user.Config.Theme, "dark")
			t.Assert(user.Config.Lang, "zh-CN")
		}
	})

	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(table).Data(g.Map{
			"id":     2,
			"name":   "user2",
			"config": `{"theme":"light","lang":"en-US"}`,
		}).Insert()
		t.AssertNil(err)

		var user User
		err = db.Model(table).WherePri(2).Scan(&user)
		t.AssertNil(err)
		t.Assert(user.Id, 2)
		t.Assert(user.Name, "user2")
		t.AssertNE(user.Config, nil)
		if user.Config != nil {
			t.Assert(user.Config.Theme, "light")
			t.Assert(user.Config.Lang, "en-US")
		}
	})
}

func Test_JSON_Transaction(t *testing.T) {
	table := createJSONTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		err := db.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
			_, err := tx.Model(table).Ctx(ctx).Data(g.Map{
				"name": "user1",
				"config": g.Map{
					"theme": "dark",
				},
			}).Insert()
			if err != nil {
				return err
			}

			_, err = tx.Model(table).Ctx(ctx).Data(g.Map{
				"config": g.Map{
					"theme": "light",
				},
			}).WherePri(1).Update()
			return err
		})
		t.AssertNil(err)

		one, err := db.Model(table).WherePri(1).One()
		t.AssertNil(err)
		t.Assert(one["NAME"], "user1")
		t.AssertNE(one["CONFIG"], nil)
		t.Assert(json.Valid(one["CONFIG"].Bytes()), true)
		t.Assert(one["CONFIG"].Map(), g.Map{"theme": "light"})
	})
}
