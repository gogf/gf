// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package sqlite_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
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
		CREATE TABLE %s (
			id          INTEGER PRIMARY KEY AUTOINCREMENT,
			name        VARCHAR(45) NULL,
			config      JSON NULL,
			metadata    JSON NULL
		);
	`, name)); err != nil {
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
		t.Assert(one["name"], "user1")
		t.AssertNE(one["config"], nil)
		t.AssertNE(one["metadata"], nil)

		valid, err := db.Model(table).
			Fields("json_valid(config) as c, json_valid(metadata) as m").WherePri(1).One()
		t.AssertNil(err)
		t.Assert(valid["c"].Int(), 1)
		t.Assert(valid["m"].Int(), 1)

		theme, err := db.Model(table).Fields("json_extract(config, '$.theme')").WherePri(1).Value()
		t.AssertNil(err)
		t.Assert(theme.String(), "dark")

		level, err := db.Model(table).Fields("json_extract(metadata, '$.level')").WherePri(1).Value()
		t.AssertNil(err)
		t.Assert(level.Int(), 5)
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
		t.Assert(one["name"], "user2")
		t.Assert(one["config"].String(), `{"theme":"light","lang":"en-US"}`)
		t.Assert(one["metadata"].String(), `{"tags":["user"],"level":1}`)
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
		t.Assert(one["name"], "user3")
		t.Assert(one["config"], nil)
		t.Assert(one["metadata"], nil)

		count, err := db.Model(table).Where("config IS NULL").Count()
		t.AssertNil(err)
		t.Assert(count, 1)
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
		t.AssertNE(one["config"], nil)

		theme, err := db.Model(table).Fields("json_extract(config, '$.theme')").WherePri(1).Value()
		t.AssertNil(err)
		t.Assert(theme.String(), "light")
	})

	gtest.C(t, func(t *gtest.T) {
		_, err := db.Exec(ctx, fmt.Sprintf(
			"UPDATE %s SET config = json_set(config, '$.lang', 'ja-JP') WHERE id = 1", table,
		))
		t.AssertNil(err)

		lang, err := db.Model(table).Fields("json_extract(config, '$.lang')").WherePri(1).Value()
		t.AssertNil(err)
		t.Assert(lang.String(), "ja-JP")
	})
}

func Test_JSON_Extract_Where(t *testing.T) {
	table := createJSONTable()
	defer dropTable(table)

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

		all, err := db.Model(table).Where("json_extract(config, '$.theme') = ?", "dark").All()
		t.AssertNil(err)
		t.Assert(len(all), 2)

		all, err = db.Model(table).Where("json_extract(config, '$.lang') = ?", "en-US").All()
		t.AssertNil(err)
		t.Assert(len(all), 2)

		all, err = db.Model(table).
			Where("json_extract(config, '$.theme') = ?", "dark").
			Where("json_extract(config, '$.lang') = ?", "en-US").All()
		t.AssertNil(err)
		t.Assert(len(all), 1)
		t.Assert(all[0]["name"].String(), "user3")
	})
}

func Test_JSON_Extract_Select(t *testing.T) {
	table := createJSONTable()
	defer dropTable(table)

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

		one, err := db.Model(table).
			Fields("name, json_extract(config, '$.theme') as theme, json_extract(metadata, '$.level') as level").
			WherePri(1).One()
		t.AssertNil(err)
		t.Assert(one["name"], "user1")
		t.Assert(one["theme"].String(), "dark")
		t.Assert(one["level"].Int(), 5)

		one, err = db.Model(table).
			Fields("json_quote(json_extract(config, '$.theme')) as theme").
			WherePri(1).One()
		t.AssertNil(err)
		t.Assert(one["theme"].String(), `"dark"`)

		one, err = db.Model(table).
			Fields("json_extract(config, '$.missing') as missing").
			WherePri(1).One()
		t.AssertNil(err)
		t.Assert(one["missing"].IsNil(), true)
	})
}

func Test_JSON_Array_Query(t *testing.T) {
	table := createJSONTable()
	defer dropTable(table)

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

		all, err := db.Model(table).
			Where("EXISTS (SELECT 1 FROM json_each(metadata, '$.tags') WHERE value = ?)", "admin").All()
		t.AssertNil(err)
		t.Assert(len(all), 2)

		all, err = db.Model(table).
			Where("EXISTS (SELECT 1 FROM json_each(metadata, '$.tags') WHERE value = ?)", "developer").All()
		t.AssertNil(err)
		t.Assert(len(all), 1)
		t.Assert(all[0]["name"].String(), "user1")

		one, err := db.Model(table).
			Fields("json_array_length(metadata, '$.tags') as n").WherePri(1).One()
		t.AssertNil(err)
		t.Assert(one["n"].Int(), 2)
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

		all, err := db.Model(table).All()
		t.AssertNil(err)
		t.Assert(len(all), 3)

		count, err := db.Model(table).Where("config IS NOT NULL").Count()
		t.AssertNil(err)
		t.Assert(count, 2)
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
			"config": nil,
		}).Insert()
		t.AssertNil(err)

		var user User
		err = db.Model(table).WherePri(2).Scan(&user)
		t.AssertNil(err)
		t.Assert(user.Name, "user2")
		t.Assert(user.Config, nil)
	})
}

func Test_JSON_Complex_Structure(t *testing.T) {
	table := createJSONTable()
	defer dropTable(table)

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

		one, err := db.Model(table).
			Fields("json_extract(config, '$.ui.theme') as theme, json_extract(config, '$.ui.fontSize.base') as base_font").
			WherePri(1).One()
		t.AssertNil(err)
		t.Assert(one["theme"].String(), "dark")
		t.Assert(one["base_font"].Int(), 14)

		one, err = db.Model(table).
			Fields("json_type(config, '$.editor') as tp, json_extract(config, '$.editor.tabSize') as tab").
			WherePri(1).One()
		t.AssertNil(err)
		t.Assert(one["tp"].String(), "object")
		t.Assert(one["tab"].Int(), 4)
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
		t.Assert(one["name"], "user1")
		t.AssertNE(one["config"], nil)

		theme, err := db.Model(table).Fields("json_extract(config, '$.theme')").WherePri(1).Value()
		t.AssertNil(err)
		t.Assert(theme.String(), "light")
	})

	gtest.C(t, func(t *gtest.T) {
		err := db.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
			_, err := tx.Model(table).Ctx(ctx).Data(g.Map{
				"name": "user2",
				"config": g.Map{
					"theme": "rolled-back",
				},
			}).Insert()
			if err != nil {
				return err
			}
			return gerror.New("rollback")
		})
		t.AssertNE(err, nil)

		count, err := db.Model(table).Count()
		t.AssertNil(err)
		t.Assert(count, 1)
	})
}
