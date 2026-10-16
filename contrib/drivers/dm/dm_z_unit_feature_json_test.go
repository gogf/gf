// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package dm_test

import (
	"context"
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
	// Note: DM has no JSON column type, JSON documents are stored in CLOB columns.
	if _, err := db.Exec(ctx, fmt.Sprintf(`
		CREATE TABLE %s (
			id          INT NOT NULL AUTO_INCREMENT,
			name        VARCHAR(45) NULL,
			config      CLOB NULL,
			metadata    CLOB NULL,
			PRIMARY KEY (id)
		)
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
		t.Assert(one["NAME"], "user1")
		t.AssertNE(one["CONFIG"], nil)
		t.AssertNE(one["METADATA"], nil)
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
	})
}

func Test_JSON_Update(t *testing.T) {
	table := createJSONTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		// Insert initial data
		_, err := db.Model(table).Data(g.Map{
			"name": "user1",
			"config": g.Map{
				"theme": "dark",
			},
		}).Insert()
		t.AssertNil(err)

		// Update JSON column
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
	})
}

func Test_JSON_Extract_Where(t *testing.T) {
	table := createJSONTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		// Insert test data
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

		// Query by JSON field using JSON_EXTRACT
		// Note: DM has no JSON_EXTRACT; JSON_VALUE extracts the scalar at the path.
		all, err := db.Model(table).Where("JSON_VALUE(config, '$.theme') = ?", "dark").All()
		t.AssertNil(err)
		t.Assert(len(all), 2)

		all, err = db.Model(table).Where("JSON_VALUE(config, '$.lang') = ?", "en-US").All()
		t.AssertNil(err)
		t.Assert(len(all), 2)
	})
}

func Test_JSON_Extract_Select(t *testing.T) {
	table := createJSONTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		// Insert test data
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

		// Select with JSON_EXTRACT
		// Note: DM has no JSON_EXTRACT; JSON_VALUE extracts the scalar at the path.
		one, err := db.Model(table).Fields("name, JSON_VALUE(config, '$.theme') as theme, JSON_VALUE(metadata, '$.level') as level").WherePri(1).One()
		t.AssertNil(err)
		t.Assert(one["NAME"], "user1")
		t.AssertNE(one["THEME"], nil)
		t.AssertNE(one["LEVEL"], nil)
	})
}

func Test_JSON_Array_Query(t *testing.T) {
	table := createJSONTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		// Insert data with JSON array
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

		// Query by JSON array contains
		// Note: DM has no JSON_CONTAINS; JSON_TABLE lists the array elements as plain values.
		all, err := db.Model(table).Where(
			"? IN (SELECT tag FROM JSON_TABLE(metadata, '$.tags[*]' COLUMNS (tag VARCHAR(45) PATH '$')))", "admin",
		).All()
		t.AssertNil(err)
		t.Assert(len(all), 2)
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
		// Insert data
		_, err := db.Model(table).Data(g.Map{
			"name": "user1",
			"config": g.Map{
				"theme": "dark",
				"lang":  "zh-CN",
			},
		}).Insert()
		t.AssertNil(err)

		// Scan to struct
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

func Test_JSON_Complex_Structure(t *testing.T) {
	table := createJSONTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		// Insert complex nested JSON
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

		// Query nested JSON path
		// Note: DM has no JSON_EXTRACT; JSON_VALUE extracts the scalar at the path.
		one, err := db.Model(table).Fields("JSON_VALUE(config, '$.ui.theme') as theme, JSON_VALUE(config, '$.ui.fontSize.base') as base_font").WherePri(1).One()
		t.AssertNil(err)
		t.AssertNE(one["THEME"], nil)
		t.AssertNE(one["BASE_FONT"], nil)
	})
}

func Test_JSON_Transaction(t *testing.T) {
	table := createJSONTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		err := db.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
			// Insert in transaction
			_, err := tx.Model(table).Ctx(ctx).Data(g.Map{
				"name": "user1",
				"config": g.Map{
					"theme": "dark",
				},
			}).Insert()
			if err != nil {
				return err
			}

			// Update in transaction
			_, err = tx.Model(table).Ctx(ctx).Data(g.Map{
				"config": g.Map{
					"theme": "light",
				},
			}).WherePri(1).Update()
			return err
		})
		t.AssertNil(err)

		// Verify data
		one, err := db.Model(table).WherePri(1).One()
		t.AssertNil(err)
		t.Assert(one["NAME"], "user1")
		t.AssertNE(one["CONFIG"], nil)
	})
}
