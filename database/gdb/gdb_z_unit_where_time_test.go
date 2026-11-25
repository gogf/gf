// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Tests for driver conversion of time arguments in WHERE conditions.

package gdb

import (
	"bytes"
	"context"
	"errors"
	"testing"
	"time"

	"github.com/gogf/gf/v2/container/gmap"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/test/gtest"
)

// whereTimeDB records metadata lookups and exposes the column type used for conversion.
type whereTimeDB struct {
	*dialectFakeDB
	err error // Conversion error returned by the driver.
}

// newWhereTimeDB installs the test driver's conversion method on the core.
func newWhereTimeDB() *whereTimeDB {
	db := &whereTimeDB{dialectFakeDB: newDialectFakeDB("id", "created_at", "updated_at")}
	db.Core.db = db
	return db
}

// ConvertValueForField exposes the supplied field type and maps a zero time to nil as Core does.
func (d *whereTimeDB) ConvertValueForField(_ context.Context, fieldType string, value any) (any, error) {
	if stamp, ok := value.(time.Time); ok && stamp.IsZero() {
		return nil, d.err
	}
	return fieldType, d.err
}

// Test_WhereBuilder_TimeField checks driver binding and metadata lookups per condition group.
func Test_WhereBuilder_TimeField(t *testing.T) {
	when := gtime.New("2018-10-24 11:22:33")
	for _, tc := range []struct {
		name    string
		where   any
		args    []any
		want    []any
		lookups int
	}{
		{"integer field", "id", []any{1}, []any{1}, 0},
		{"map with integer field", Map{"id": 1}, nil, []any{1}, 0},
		{"struct field mapping", struct{ Id int }{1}, nil, []any{1}, 1},
		{"SQL expression", "DATE(created_at) = ?", []any{when}, []any{when.Time}, 0},
		{"unknown field", "other_date", []any{when}, []any{when.Time}, 1},
		{"time field", "created_at", []any{when.Time}, []any{"datetime"}, 1},
		{"ordered map with time", gmap.NewListMapFrom(map[any]any{"created_at": when}), nil, []any{"datetime"}, 1},
		{"two time fields", Map{"created_at": *when, "updated_at": &when.Time}, nil, []any{"datetime", "datetime"}, 1},
	} {
		t.Run(tc.name, func(t *testing.T) {
			gtest.C(t, func(t *gtest.T) {
				db := newWhereTimeDB()
				_, args := db.Model("events").Builder().Where(tc.where, tc.args...).Build()
				t.AssertEQ(args, tc.want)
				t.AssertEQ(len(db.lookups), tc.lookups)
			})
		})
	}
}

// Test_WhereBuilder_NullTime checks that only a nil time pointer produces IS NULL.
func Test_WhereBuilder_NullTime(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		builder := newWhereTimeDB().Model("events").Builder()
		where, args := builder.Where("created_at", (*time.Time)(nil)).Build()
		t.AssertEQ(where, `"created_at" IS NULL`)
		t.AssertEQ(len(args), 0)

		where, args = builder.Where("created_at", time.Time{}).Build()
		t.AssertEQ(where, `"created_at"=?`)
		t.AssertEQ(args, []any{time.Time{}})
	})
}

// Test_WhereBuilder_TimeConversionFallback checks that a failed conversion keeps the original time argument.
func Test_WhereBuilder_TimeConversionFallback(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			logs    bytes.Buffer
			db      = newWhereTimeDB()
			when    = time.Now()
			builder = db.Model("events").Builder().Where("created_at", when)
		)
		db.err = errors.New("driver rejected time")
		db.SetLogger(glog.NewWithWriter(&logs))
		where, args := builder.Build()
		t.AssertEQ(where, `"created_at"=?`)
		t.AssertEQ(args, []any{when})
		t.AssertIN("WARN", logs.String())
		t.AssertIN("created_at", logs.String())
		t.AssertIN(db.err.Error(), logs.String())

		db.SetLogger(nil)
		_, args = builder.Build()
		t.AssertEQ(args, []any{when})
	})
}
