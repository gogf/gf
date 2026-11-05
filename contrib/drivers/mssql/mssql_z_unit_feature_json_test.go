// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package mssql_test

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
		name = fmt.Sprintf(`json_table_%d`, gtime.TimestampNano())
	}
	dropTable(name)
	if _, err := db.Exec(ctx, fmt.Sprintf(`
		CREATE TABLE [%s] (
			ID          int IDENTITY(1,1) NOT NULL,
			NAME        varchar(45) NULL,
			CONFIG      nvarchar(max) NULL,
			METADATA    nvarchar(max) NULL,
			PRIMARY KEY (ID)
		)`, name)); err != nil {
		gtest.Fatal(err)
	}
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

		valid, err := db.Model(table).Fields("ISJSON(CONFIG) AS V1, ISJSON(METADATA) AS V2").WherePri(1).One()
		t.AssertNil(err)
		t.Assert(valid["V1"], 1)
		t.Assert(valid["V2"], 1)
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

func Test_JSON_Extract_Where(t *testing.T) {
	table := createJSONTable()
	defer dropTable(table)
	table2 := createJSONTable()
	defer dropTable(table2)

	gtest.C(t, func(t *gtest.T) {
		data := g.Slice{
			g.Map{
				"name": "user1",
				"config": g.Map{
					"theme": "dark",
					"lang":  "zh-CN",
				},
			},
			g.Map{
				"name": "user2",
				"config": g.Map{
					"theme": "light",
					"lang":  "en-US",
				},
			},
			g.Map{
				"name": "user3",
				"config": g.Map{
					"theme": "dark",
					"lang":  "en-US",
				},
			},
		}
		_, err := db.Model(table).Data(data).Insert()
		t.AssertNil(err)

		all, err := db.Model(table).Where("JSON_VALUE(config, '$.theme') = ?", "dark").All()
		t.AssertNil(err)
		t.Assert(len(all), 2)

		all, err = db.Model(table).Where("JSON_VALUE(config, '$.lang') = ?", "en-US").All()
		t.AssertNil(err)
		t.Assert(len(all), 2)
	})

	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(table2).Data(g.Slice{
			g.Map{"name": "user1", "config": `{"theme":"dark","lang":"zh-CN"}`},
			g.Map{"name": "user2", "config": `{"theme":"light","lang":"en-US"}`},
			g.Map{"name": "user3", "config": `{"theme":"dark","lang":"en-US"}`},
		}).Insert()
		t.AssertNil(err)

		names, err := db.Model(table2).Fields("name").Where("JSON_VALUE(config, '$.theme') = ?", "dark").Order("id").Array()
		t.AssertNil(err)
		t.Assert(names, g.Slice{"user1", "user3"})

		names, err = db.Model(table2).Fields("name").Where("JSON_VALUE(config, '$.lang') = ?", "en-US").Order("id").Array()
		t.AssertNil(err)
		t.Assert(names, g.Slice{"user2", "user3"})
	})
}

func Test_JSON_Extract_Select(t *testing.T) {
	table := createJSONTable()
	defer dropTable(table)
	table2 := createJSONTable()
	defer dropTable(table2)

	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(table).Data(g.Map{
			"name": "user1",
			"config": g.Map{
				"theme": "dark",
				"lang":  "zh-CN",
			},
			"metadata": g.Map{
				"level": 5,
			},
		}).Insert()
		t.AssertNil(err)

		one, err := db.Model(table).Fields("name, JSON_VALUE(config, '$.theme') as theme, JSON_VALUE(metadata, '$.level') as level").WherePri(1).One()
		t.AssertNil(err)
		t.Assert(one["name"], "user1")
		t.AssertNE(one["theme"], nil)
		t.AssertNE(one["level"], nil)
		t.Assert(one["theme"], "dark")
		t.Assert(one["level"], 5)
	})

	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(table2).Data(g.Map{
			"name":     "user1",
			"config":   `{"theme":"dark","lang":"zh-CN"}`,
			"metadata": `{"level":5}`,
		}).Insert()
		t.AssertNil(err)

		one, err := db.Model(table2).Fields("name, JSON_VALUE(config, '$.theme') as theme, JSON_VALUE(metadata, '$.level') as level").WherePri(1).One()
		t.AssertNil(err)
		t.Assert(one["name"], "user1")
		t.Assert(one["theme"], "dark")
		t.Assert(one["level"], 5)
	})
}

func Test_JSON_Array_Query(t *testing.T) {
	table := createJSONTable()
	defer dropTable(table)
	table2 := createJSONTable()
	defer dropTable(table2)

	gtest.C(t, func(t *gtest.T) {
		data := g.Slice{
			g.Map{
				"name": "user1",
				"metadata": g.Map{
					"tags": g.Slice{"admin", "developer"},
				},
			},
			g.Map{
				"name": "user2",
				"metadata": g.Map{
					"tags": g.Slice{"user"},
				},
			},
			g.Map{
				"name": "user3",
				"metadata": g.Map{
					"tags": g.Slice{"admin", "user"},
				},
			},
		}
		_, err := db.Model(table).Data(data).Insert()
		t.AssertNil(err)

		all, err := db.Model(table).Where("EXISTS (SELECT 1 FROM OPENJSON(metadata, '$.tags') WHERE [value] = ?)", "admin").All()
		t.AssertNil(err)
		t.Assert(len(all), 2)
	})

	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(table2).Data(g.Slice{
			g.Map{"name": "user1", "metadata": `{"tags":["admin","developer"]}`},
			g.Map{"name": "user2", "metadata": `{"tags":["user"]}`},
			g.Map{"name": "user3", "metadata": `{"tags":["admin","user"]}`},
		}).Insert()
		t.AssertNil(err)

		names, err := db.Model(table2).Fields("name").
			Where("EXISTS (SELECT 1 FROM OPENJSON(metadata, '$.tags') WHERE [value] = ?)", "admin").
			Order("id").Array()
		t.AssertNil(err)
		t.Assert(names, g.Slice{"user1", "user3"})
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
		t.Assert(json.Valid(all[1]["CONFIG"].Bytes()), true)
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
			"name":   "user2",
			"config": `{"theme":"light","lang":"en-US"}`,
		}).Insert()
		t.AssertNil(err)

		var user User
		err = db.Model(table).Where("name", "user2").Scan(&user)
		t.AssertNil(err)
		t.Assert(user.Name, "user2")
		t.AssertNE(user.Config, nil)
		if user.Config != nil {
			t.Assert(user.Config.Theme, "light")
			t.Assert(user.Config.Lang, "en-US")
		}
	})
}

func Test_JSON_Complex_Structure(t *testing.T) {
	table := createJSONTable()
	defer dropTable(table)
	table2 := createJSONTable()
	defer dropTable(table2)

	gtest.C(t, func(t *gtest.T) {
		data := g.Map{
			"name": "user1",
			"config": g.Map{
				"ui": g.Map{
					"theme": "dark",
					"fontSize": g.Map{
						"base": 14,
						"code": 12,
					},
				},
				"editor": g.Map{
					"tabSize":  4,
					"wordWrap": true,
				},
			},
		}
		result, err := db.Model(table).Data(data).Insert()
		t.AssertNil(err)
		n, _ := result.LastInsertId()
		t.Assert(n, 1)

		one, err := db.Model(table).Fields("JSON_VALUE(config, '$.ui.theme') as theme, JSON_VALUE(config, '$.ui.fontSize.base') as base_font").WherePri(1).One()
		t.AssertNil(err)
		t.AssertNE(one["theme"], nil)
		t.AssertNE(one["base_font"], nil)
		t.Assert(one["theme"], "dark")
		t.Assert(one["base_font"], 14)
	})

	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(table2).Data(g.Map{
			"name":   "user1",
			"config": `{"ui":{"theme":"dark","fontSize":{"base":14,"code":12}},"editor":{"tabSize":4,"wordWrap":true}}`,
		}).Insert()
		t.AssertNil(err)

		one, err := db.Model(table2).Fields(
			"JSON_VALUE(config, '$.ui.theme') as theme, JSON_VALUE(config, '$.ui.fontSize.base') as base_font, " +
				"JSON_VALUE(config, '$.editor.wordWrap') as word_wrap, JSON_QUERY(config, '$.ui.fontSize') as font_size",
		).WherePri(1).One()
		t.AssertNil(err)
		t.Assert(one["theme"], "dark")
		t.Assert(one["base_font"], 14)
		t.Assert(one["word_wrap"], "true")
		t.Assert(one["font_size"].Map(), g.Map{"base": 14, "code": 12})
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
