// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// This file tests that records convert to structs by orm tags differing from the field names of
// the records in letter case, as the column names of a database folding identifiers to upper case.

package gdb

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/test/gtest"
)

// Test_Record_Struct_UpperCaseFieldNames asserts that Record.Struct and Result.Structs bind each
// field of a record to the attribute whose orm tag it equals in case-insensitive, whatever other
// attribute name it equals once symbols are removed. Each round converts to a new struct type, one
// more attribute making it differ, so that no match cached by a previous round decides the result.
func Test_Record_Struct_UpperCaseFieldNames(t *testing.T) {
	var (
		intType    = reflect.TypeOf(0)
		structType = func(round int) reflect.Type {
			return reflect.StructOf([]reflect.StructField{
				{Name: "Id", Type: intType, Tag: `orm:"id"`},
				{Name: "Userid", Type: intType, Tag: `orm:"user_id"`},
				{Name: "UserId", Type: intType, Tag: `orm:"UserId"`},
				{Name: fmt.Sprintf("Round%d", round), Type: intType, Tag: `orm:"-"`},
			})
		}
		record = Record{
			"ID":      gvar.New(1),
			"USER_ID": gvar.New(2),
			"USERID":  gvar.New(3),
		}
		assertUser = func(t *gtest.T, user reflect.Value) {
			t.Assert(user.FieldByName("Id").Int(), 1)
			t.Assert(user.FieldByName("Userid").Int(), 2)
			t.Assert(user.FieldByName("UserId").Int(), 3)
		}
	)
	gtest.C(t, func(t *gtest.T) {
		for i := 0; i < 100; i++ {
			user := reflect.New(structType(i))
			t.AssertNil(record.Struct(user.Interface()))
			assertUser(t, user.Elem())
		}
	})
	gtest.C(t, func(t *gtest.T) {
		for i := 100; i < 200; i++ {
			users := reflect.New(reflect.SliceOf(structType(i)))
			t.AssertNil(Result{record, record}.Structs(users.Interface()))
			t.Assert(users.Elem().Len(), 2)
			assertUser(t, users.Elem().Index(0))
			assertUser(t, users.Elem().Index(1))
		}
	})
}
