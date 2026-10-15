// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package clickhouse_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/test/gtest"
)

// chCCreateJSONTable creates a table storing JSON documents in String columns, as the JSON type
// of ClickHouse 24.11 is experimental.
func chCCreateJSONTable() string {
	name := fmt.Sprintf(`json_table_%d`, gtime.TimestampNano())
	dropTable(name)
	if _, err := db.Exec(ctx, fmt.Sprintf(`
		CREATE TABLE %s (
			id       UInt64,
			name     Nullable(String),
			config   Nullable(String),
			metadata Nullable(String)
		) ENGINE = MergeTree()
		ORDER BY id
	`, name)); err != nil {
		gtest.Fatal(err)
	}
	return name
}

// Test_JSON_Insert_Map tests inserting maps into JSON columns.
// Note: ClickHouse reports no last insert id, so the inserted record is checked instead.
func Test_JSON_Insert_Map(t *testing.T) {
	table := chCCreateJSONTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		data := g.Map{
			"id":   1,
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
		_, err := db.Model(table).Data(data).Insert()
		t.AssertNil(err)

		one, err := db.Model(table).WherePri(1).One()
		t.AssertNil(err)
		t.Assert(one["name"], "user1")
		t.AssertNE(one["config"], nil)
		t.AssertNE(one["metadata"], nil)
		t.Assert(one["config"].Map(), g.Map{"theme": "dark", "lang": "zh-CN"})
		t.Assert(one["metadata"].Map(), g.Map{"tags": g.Slice{"admin", "developer"}, "level": 5})
	})
}

// Test_JSON_Insert_String tests inserting JSON strings into JSON columns.
// Note: ClickHouse reports no last insert id, so the inserted record is checked instead.
func Test_JSON_Insert_String(t *testing.T) {
	table := chCCreateJSONTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		data := g.Map{
			"id":       1,
			"name":     "user2",
			"config":   `{"theme":"light","lang":"en-US"}`,
			"metadata": `{"tags":["user"],"level":1}`,
		}
		_, err := db.Model(table).Data(data).Insert()
		t.AssertNil(err)

		one, err := db.Model(table).WherePri(1).One()
		t.AssertNil(err)
		t.Assert(one["name"], "user2")
		t.Assert(one["config"], `{"theme":"light","lang":"en-US"}`)
		t.Assert(one["metadata"], `{"tags":["user"],"level":1}`)
	})
}

// Test_JSON_Insert_Null tests inserting NULL into JSON columns.
// Note: ClickHouse reports no last insert id, so the inserted record is checked instead.
func Test_JSON_Insert_Null(t *testing.T) {
	table := chCCreateJSONTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		data := g.Map{
			"id":       1,
			"name":     "user3",
			"config":   nil,
			"metadata": nil,
		}
		_, err := db.Model(table).Data(data).Insert()
		t.AssertNil(err)

		one, err := db.Model(table).WherePri(1).One()
		t.AssertNil(err)
		t.Assert(one["name"], "user3")
		t.Assert(one["config"], nil)
		t.Assert(one["metadata"], nil)
	})
}

// Test_JSON_Update tests updating a JSON column with a map.
// Note: ClickHouse reports no affected rows for a mutation, so the updated record is checked instead.
func Test_JSON_Update(t *testing.T) {
	table := chCCreateJSONTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(table).Data(g.Map{
			"id":   1,
			"name": "user1",
			"config": g.Map{
				"theme": "dark",
			},
		}).Insert()
		t.AssertNil(err)

		_, err = db.Model(table).Data(g.Map{
			"config": g.Map{
				"theme": "light",
				"lang":  "en-US",
			},
		}).WherePri(1).Update()
		t.AssertNil(err)

		one, err := db.Model(table).WherePri(1).One()
		t.AssertNil(err)
		t.AssertNE(one["config"], nil)
		t.Assert(one["config"].Map(), g.Map{"theme": "light", "lang": "en-US"})
	})
}

// Test_JSON_Extract_Where tests filtering on JSON fields with JSONExtractString.
func Test_JSON_Extract_Where(t *testing.T) {
	table := chCCreateJSONTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		data := g.Slice{
			g.Map{
				"id":   1,
				"name": "user1",
				"config": g.Map{
					"theme": "dark",
					"lang":  "zh-CN",
				},
			},
			g.Map{
				"id":   2,
				"name": "user2",
				"config": g.Map{
					"theme": "light",
					"lang":  "en-US",
				},
			},
			g.Map{
				"id":   3,
				"name": "user3",
				"config": g.Map{
					"theme": "dark",
					"lang":  "en-US",
				},
			},
		}
		_, err := db.Model(table).Data(data).Insert()
		t.AssertNil(err)

		all, err := db.Model(table).Where("JSONExtractString(config, 'theme') = ?", "dark").OrderAsc("id").All()
		t.AssertNil(err)
		t.Assert(len(all), 2)
		t.Assert(all[0]["name"], "user1")
		t.Assert(all[1]["name"], "user3")

		all, err = db.Model(table).Where("JSONExtractString(config, 'lang') = ?", "en-US").OrderAsc("id").All()
		t.AssertNil(err)
		t.Assert(len(all), 2)
		t.Assert(all[0]["name"], "user2")
		t.Assert(all[1]["name"], "user3")
	})
}

// Test_JSON_Extract_Select tests selecting JSON fields with JSONExtract functions.
func Test_JSON_Extract_Select(t *testing.T) {
	table := chCCreateJSONTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(table).Data(g.Map{
			"id":   1,
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

		one, err := db.Model(table).Fields(
			"name, JSONExtractString(config, 'theme') as theme, JSONExtractInt(metadata, 'level') as level",
		).WherePri(1).One()
		t.AssertNil(err)
		t.Assert(one["name"], "user1")
		t.AssertNE(one["theme"], nil)
		t.AssertNE(one["level"], nil)
		t.Assert(one["theme"], "dark")
		t.Assert(one["level"], 5)
	})
}

// Test_JSON_Array_Query tests filtering on the elements of a JSON array.
func Test_JSON_Array_Query(t *testing.T) {
	table := chCCreateJSONTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		data := g.Slice{
			g.Map{
				"id":   1,
				"name": "user1",
				"metadata": g.Map{
					"tags": g.Slice{"admin", "developer"},
				},
			},
			g.Map{
				"id":   2,
				"name": "user2",
				"metadata": g.Map{
					"tags": g.Slice{"user"},
				},
			},
			g.Map{
				"id":   3,
				"name": "user3",
				"metadata": g.Map{
					"tags": g.Slice{"admin", "user"},
				},
			},
		}
		_, err := db.Model(table).Data(data).Insert()
		t.AssertNil(err)

		all, err := db.Model(table).Where("has(JSONExtract(ifNull(metadata, ''), 'tags', 'Array(String)'), ?)", "admin").OrderAsc("id").All()
		t.AssertNil(err)
		t.Assert(len(all), 2)
		t.Assert(all[0]["name"], "user1")
		t.Assert(all[1]["name"], "user3")
	})
}

// Test_JSON_Batch_Insert tests inserting a batch of records with JSON columns.
// Note: ClickHouse reports no affected rows, so the inserted records are counted instead.
func Test_JSON_Batch_Insert(t *testing.T) {
	table := chCCreateJSONTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		data := g.Slice{
			g.Map{
				"id":   1,
				"name": "user1",
				"config": g.Map{
					"theme": "dark",
				},
			},
			g.Map{
				"id":   2,
				"name": "user2",
				"config": g.Map{
					"theme": "light",
				},
			},
			g.Map{
				"id":     3,
				"name":   "user3",
				"config": nil,
			},
		}
		_, err := db.Model(table).Data(data).Insert()
		t.AssertNil(err)

		all, err := db.Model(table).OrderAsc("id").All()
		t.AssertNil(err)
		t.Assert(len(all), 3)
		t.Assert(all[0]["config"].Map(), g.Map{"theme": "dark"})
		t.Assert(all[1]["config"].Map(), g.Map{"theme": "light"})
		t.Assert(all[2]["config"], nil)
	})
}

// Test_JSON_Scan_To_Struct tests scanning a JSON column into a struct attribute.
func Test_JSON_Scan_To_Struct(t *testing.T) {
	table := chCCreateJSONTable()
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
			"id":   1,
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
}

// Test_JSON_Complex_Structure tests extracting nested JSON paths.
func Test_JSON_Complex_Structure(t *testing.T) {
	table := chCCreateJSONTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		data := g.Map{
			"id":   1,
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
		_, err := db.Model(table).Data(data).Insert()
		t.AssertNil(err)

		one, err := db.Model(table).Fields(
			"JSONExtractString(config, 'ui', 'theme') as theme, JSONExtractInt(config, 'ui', 'fontSize', 'base') as base_font",
		).WherePri(1).One()
		t.AssertNil(err)
		t.AssertNE(one["theme"], nil)
		t.AssertNE(one["base_font"], nil)
		t.Assert(one["theme"], "dark")
		t.Assert(one["base_font"], 14)
	})
}

// Test_JSON_Transaction tests writing JSON columns inside a transaction.
// Note: ClickHouse has no transactions, so Transaction fails as unsupported before any write.
func Test_JSON_Transaction(t *testing.T) {
	table := chCCreateJSONTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		err := db.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
			_, err := tx.Model(table).Ctx(ctx).Data(g.Map{
				"id":   1,
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
		t.AssertNE(err, nil)
		t.Assert(err.Error(), "unsupported method:Transaction")

		count, err := db.Model(table).Count()
		t.AssertNil(err)
		t.Assert(count, 0)
	})
}
