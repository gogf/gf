// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gaussdb_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/test/gtest"
)

func Test_Model_Hook_Select(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		m := db.Model(table).Hook(
			gdb.AfterSelect(func(ctx context.Context, in *gdb.HookSelectInput, result gdb.Result, err error) (gdb.Result, error) {
				if err != nil {
					return result, err
				}
				for i, record := range result {
					record["test"] = gvar.New(100 + record["id"].Int())
					result[i] = record
				}
				return result, nil
			}),
		)
		all, err := m.Where("id > ?", 6).OrderAsc("id").All()
		t.AssertNil(err)
		t.Assert(len(all), 4)
		t.Assert(all[0]["id"].Int(), 7)
		t.Assert(all[0]["test"].Int(), 107)
		t.Assert(all[1]["test"].Int(), 108)
		t.Assert(all[2]["test"].Int(), 109)
		t.Assert(all[3]["test"].Int(), 110)
	})
}

func Test_Model_Hook_Insert(t *testing.T) {
	table := createTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		m := db.Model(table).Hook(
			gdb.BeforeInsert(func(ctx context.Context, in *gdb.HookInsertInput) error {
				for i, item := range in.Data {
					item["passport"] = fmt.Sprintf(`test_port_%d`, item["id"])
					item["nickname"] = fmt.Sprintf(`test_name_%d`, item["id"])
					item["password"] = fmt.Sprintf(`test_pass_%d`, item["id"])
					item["create_time"] = CreateTime
					in.Data[i] = item
				}
				return nil
			}),
		)
		_, err := m.Insert(g.Map{
			"id":       1,
			"nickname": "name_1",
		})
		t.AssertNil(err)
		one, err := m.One()
		t.AssertNil(err)
		t.Assert(one["id"].Int(), 1)
		t.Assert(one["passport"], `test_port_1`)
		t.Assert(one["nickname"], `test_name_1`)
	})
}

func Test_Model_Hook_Update(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		m := db.Model(table).Hook(
			gdb.BeforeUpdate(func(ctx context.Context, in *gdb.HookUpdateInput) error {
				switch value := in.Data.(type) {
				case gdb.List:
					for i, data := range value {
						data["passport"] = `port`
						data["nickname"] = `name`
						value[i] = data
					}
					in.Data = value

				case gdb.Map:
					value["passport"] = `port`
					value["nickname"] = `name`
					in.Data = value
				}
				return nil
			}),
		)
		_, err := m.Data(g.Map{
			"nickname": "name_1",
		}).WherePri(1).Update()
		t.AssertNil(err)

		one, err := m.One()
		t.AssertNil(err)
		t.Assert(one["id"].Int(), 1)
		t.Assert(one["passport"], `port`)
		t.Assert(one["nickname"], `name`)
	})
}

func Test_Model_Hook_Delete(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		m := db.Model(table).Hook(
			gdb.BeforeDelete(func(ctx context.Context, in *gdb.HookDeleteInput) error {
				origCondition := in.Condition
				in.Condition = "1=0"
				_, err := in.Model.Data(g.Map{
					"nickname": `deleted`,
				}).Where(origCondition).Update()
				return err
			}),
		)
		_, err := m.Where("1=1").Delete()
		t.AssertNil(err)

		all, err := m.All()
		t.AssertNil(err)
		for _, item := range all {
			t.Assert(item["nickname"].String(), `deleted`)
		}
	})
}

func Test_Model_Hook_Select_Count(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		m := db.Model(table).Hook(
			gdb.AfterSelect(func(ctx context.Context, in *gdb.HookSelectInput, result gdb.Result, err error) (gdb.Result, error) {
				if err != nil {
					return result, err
				}
				// Adding extra fields should not affect Count operations.
				for i, record := range result {
					record["extra"] = gvar.New("extra_value")
					result[i] = record
				}
				return result, nil
			}),
		)
		count, err := m.Count()
		t.AssertNil(err)
		t.Assert(count, TableSize)
	})
}

func Test_Model_Hook_Chain(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	// Normal chain: two hooks both modify data
	gtest.C(t, func(t *gtest.T) {
		m := db.Model(table).
			Hook(
				gdb.AfterSelect(func(ctx context.Context, in *gdb.HookSelectInput, result gdb.Result, err error) (gdb.Result, error) {
					if err != nil {
						return result, err
					}
					for i, record := range result {
						record["hook1"] = gvar.New("value1")
						result[i] = record
					}
					return result, nil
				}),
			).
			Hook(
				gdb.AfterSelect(func(ctx context.Context, in *gdb.HookSelectInput, result gdb.Result, err error) (gdb.Result, error) {
					if err != nil {
						return result, err
					}
					for i, record := range result {
						record["hook2"] = gvar.New("value2")
						result[i] = record
					}
					return result, nil
				}),
			)
		all, err := m.Where("id", 1).All()
		t.AssertNil(err)
		t.Assert(len(all), 1)
		t.Assert(all[0]["id"].Int(), 1)
		t.Assert(all[0]["hook1"].String(), "value1")
		t.Assert(all[0]["hook2"].String(), "value2")
	})

	// Error chain: hook returns error
	gtest.C(t, func(t *gtest.T) {
		m := db.Model(table).Hook(
			gdb.AfterSelect(func(ctx context.Context, in *gdb.HookSelectInput, result gdb.Result, err error) (gdb.Result, error) {
				return nil, gerror.New("hook error")
			}),
		)
		_, err := m.Where("id", 1).All()
		t.AssertNE(err, nil)
		t.Assert(err.Error(), "hook error")
	})
}
