// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package mysql_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/gogf/gf/v3/container/gvar"
	"github.com/gogf/gf/v3/database/gdb"
	"github.com/gogf/gf/v3/frame/g"
	"github.com/gogf/gf/v3/test/gtest"
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
		all, err := m.Where(`id > 6`).OrderAsc(`id`).All(ctx)
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
					in.Data[i] = item
				}
				return nil
			}),
		)
		_, err := m.Data(g.Map{
			"id":       1,
			"nickname": "name_1",
		}).Insert(ctx)
		t.AssertNil(err)
		one, err := m.One(ctx)
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
		}).WherePri(1).Update(ctx)
		t.AssertNil(err)

		one, err := m.One(ctx)
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
				// Make delete a no-op, then execute the intended update.
				in.Condition = "1=0"
				_, err := in.Model.Data(g.Map{
					"nickname": `deleted`,
				}).Where(origCondition).Update(ctx)
				return err
			}),
		)
		_, err := m.Where(1).Delete(ctx)
		t.AssertNil(err)

		all, err := m.All(ctx)
		t.AssertNil(err)
		for _, item := range all {
			t.Assert(item["nickname"].String(), `deleted`)
		}
	})
}
// Test_Model_Hook_Multiple tests multiple hooks execution order
func Test_Model_Hook_Multiple(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		var afterCalls []string
		m := db.Model(table).
			Hook(
				gdb.AfterSelect(func(ctx context.Context, in *gdb.HookSelectInput, result gdb.Result, err error) (gdb.Result, error) {
					afterCalls = append(afterCalls, "hook1")
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
					afterCalls = append(afterCalls, "hook2")
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

		_, err := m.Where("id", 1).One(ctx)
		t.AssertNil(err)

		one, err := m.One(ctx)
		t.AssertNil(err)
		t.Assert(one["hook1"].String(), "value1")
		t.Assert(one["hook2"].String(), "value2")
		t.Assert(afterCalls, g.Slice{"hook1", "hook2"})
	})
}

// Test_Model_Hook_Error_Abort tests hook returning error aborts operation
func Test_Model_Hook_Error_Abort(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		m := db.Model(table).Hook(
			gdb.BeforeInsert(func(ctx context.Context, in *gdb.HookInsertInput) error {
				// Return error to abort insert.
				return fmt.Errorf("hook aborted insert")
			}),
		)

		_, err := m.Data(g.Map{
			"passport": "test_abort",
			"password": "pass",
			"nickname": "name",
		}).Insert(ctx)
		t.AssertNE(err, nil)
		t.Assert(err.Error(), "hook aborted insert")

		// Verify record was not inserted
		count, err := db.Model(table).Where("passport", "test_abort").Count(ctx)
		t.AssertNil(err)
		t.Assert(count, 0)
	})
}

// Test_Model_Hook_Modify_Data tests hook modifying data before insert
func Test_Model_Hook_Modify_Data(t *testing.T) {
	table := createTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		m := db.Model(table).Hook(
			gdb.BeforeInsert(func(ctx context.Context, in *gdb.HookInsertInput) error {
				// Modify all data items
				for i := range in.Data {
					in.Data[i]["password"] = "encrypted_" + fmt.Sprint(in.Data[i]["password"])
					in.Data[i]["nickname"] = "verified_" + fmt.Sprint(in.Data[i]["nickname"])
				}
				return nil
			}),
		)

		_, err := m.Data(g.Map{
			"passport": "test_user",
			"password": "plain123",
			"nickname": "john",
		}).Insert(ctx)
		t.AssertNil(err)

		// Verify data was modified by hook
		one, err := db.Model(table).Where("passport", "test_user").One(ctx)
		t.AssertNil(err)
		t.Assert(one["password"].String(), "encrypted_plain123")
		t.Assert(one["nickname"].String(), "verified_john")
	})
}
