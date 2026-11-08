// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gdb

import (
	"testing"

	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/test/gtest"
)

// Test_Issue4699 tests negative values for Limit/Page/Offset should be treated as zero.
// See https://github.com/gogf/gf/issues/4699
func Test_Issue4699(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// Create a base model for testing
		m := &Model{}

		// Test Limit with single negative parameter
		m1 := m.Limit(-1)
		t.AssertEQ(m1.limit, 0)

		// Test Limit with two parameters (start, limit) where both are negative
		m2 := m.Limit(-10, -5)
		t.AssertEQ(m2.start, 0)
		t.AssertEQ(m2.limit, 0)

		// Test Limit with mixed parameters (negative start, positive limit)
		m3 := m.Limit(-10, 5)
		t.AssertEQ(m3.start, 0)
		t.AssertEQ(m3.limit, 5)

		// Test Page with negative limit
		m4 := m.Page(1, -10)
		t.AssertEQ(m4.start, 0)
		t.AssertEQ(m4.limit, 0)

		// Test Page with negative limit on page 2
		m5 := m.Page(2, -10)
		t.AssertEQ(m5.start, 0) // (2-1) * 0 = 0
		t.AssertEQ(m5.limit, 0)

		// Test Offset with negative value
		m6 := m.Offset(-5)
		t.AssertEQ(m6.offset, 0)

		// Test Offset with positive value (sanity check)
		m7 := m.Offset(10)
		t.AssertEQ(m7.offset, 10)
	})
}

// Test_ScanValidateSingleFieldSpecified tests that validateSingleFieldSpecified
// and isSingleFieldSpecified correctly accept/reject various field configurations,
// especially the gdb.Raw case which bypasses len(m.fields) checks.
func Test_ScanValidateSingleFieldSpecified(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// No fields at all → reject.
		m0 := &Model{}
		t.AssertNE(m0.validateSingleFieldSpecified(), nil)
		t.Assert(m0.isSingleFieldSpecified(), false)

		// Exactly one normal field → accept.
		m1 := &Model{fields: []any{"name"}}
		t.AssertNil(m1.validateSingleFieldSpecified())
		t.Assert(m1.isSingleFieldSpecified(), true)

		// Two normal fields → reject.
		m2 := &Model{fields: []any{"name", "age"}}
		t.AssertNE(m2.validateSingleFieldSpecified(), nil)
		t.Assert(m2.isSingleFieldSpecified(), false)

		// gdb.Raw("name,age") stored as one field entry → reject,
		// because it may expand to multiple columns at the SQL level.
		m3 := &Model{fields: []any{Raw("name,age")}}
		t.AssertNE(m3.validateSingleFieldSpecified(), nil)
		t.Assert(m3.isSingleFieldSpecified(), false)

		// FieldsEx only → accepted. FieldsEx can narrow the result to a single
		// column; with a nil db the column-count check is skipped, so validation
		// passes. isSingleFieldSpecified still reports false because FieldsEx is
		// not an explicit single-field declaration.
		m4 := &Model{fieldsEx: []any{"id"}}
		t.AssertNil(m4.validateSingleFieldSpecified())
		t.Assert(m4.isSingleFieldSpecified(), false)

		// gdb.Raw("name") as single field → accept, because it resolves to a
		// single column and is not an expanding field.
		m5 := &Model{fields: []any{Raw("name")}}
		t.AssertNil(m5.validateSingleFieldSpecified())
		t.Assert(m5.isSingleFieldSpecified(), true)

		// Wildcard "*" → reject (expands to all columns).
		m6 := &Model{fields: []any{"*"}}
		t.AssertNE(m6.validateSingleFieldSpecified(), nil)
		t.Assert(m6.isSingleFieldSpecified(), false)

		// Table-qualified wildcard "a.*" → reject (expands to all columns of alias a).
		m7 := &Model{fields: []any{"a.*"}}
		t.AssertNE(m7.validateSingleFieldSpecified(), nil)
		t.Assert(m7.isSingleFieldSpecified(), false)

		// COUNT(*) → accept (aggregate function, always produces one column).
		m8 := &Model{fields: []any{"COUNT(*)"}}
		t.AssertNil(m8.validateSingleFieldSpecified())
		t.Assert(m8.isSingleFieldSpecified(), true)
	})
}

// Test_Issue4747 tests that Record.Struct and Result.Structs convert nested
// Result values wrapped by gvar.Var at every level.
// See https://github.com/gogf/gf/issues/4747
func Test_Issue4747(t *testing.T) {
	type BuyTrack struct {
		BuyTrackId     *uint    `json:"buy_track_id"`
		BuyTrackName   *string  `json:"buy_track_name"`
		BuyTrackWeight *float64 `json:"buy_track_weight"`
	}
	type Track struct {
		TrackId      *uint      `json:"track_id"`
		TrackWeight  *float64   `json:"track_weight"`
		BuyTrackList []BuyTrack `json:"buy_track_list"`
	}
	type Item struct {
		ItemId    *uint   `json:"item_id"`
		TrackList []Track `json:"track_list"`
	}
	var (
		buyTrackList = Result{{
			"buy_track_id":     gvar.New(1),
			"buy_track_name":   gvar.New("1111"),
			"buy_track_weight": gvar.New("2.5"),
		}}
		trackList = Result{{
			"track_id":       gvar.New(2),
			"track_weight":   gvar.New("3.5"),
			"buy_track_list": gvar.New(buyTrackList),
		}}
		item = Record{
			"item_id":    gvar.New(3),
			"track_list": gvar.New(trackList),
		}
	)
	gtest.C(t, func(t *gtest.T) {
		var info Item
		t.AssertNil(item.Struct(&info))
		t.Assert(*info.ItemId, 3)
		t.Assert(len(info.TrackList), 1)
		t.Assert(*info.TrackList[0].TrackId, 2)
		t.Assert(*info.TrackList[0].TrackWeight, 3.5)
		t.Assert(len(info.TrackList[0].BuyTrackList), 1)
		t.Assert(*info.TrackList[0].BuyTrackList[0].BuyTrackId, 1)
		t.Assert(*info.TrackList[0].BuyTrackList[0].BuyTrackName, "1111")
		t.Assert(*info.TrackList[0].BuyTrackList[0].BuyTrackWeight, 2.5)
	})
	gtest.C(t, func(t *gtest.T) {
		var list []Item
		t.AssertNil(Result{item}.Structs(&list))
		t.Assert(len(list), 1)
		t.Assert(*list[0].ItemId, 3)
		t.Assert(len(list[0].TrackList), 1)
		t.Assert(*list[0].TrackList[0].TrackWeight, 3.5)
		t.Assert(len(list[0].TrackList[0].BuyTrackList), 1)
		t.Assert(*list[0].TrackList[0].BuyTrackList[0].BuyTrackWeight, 2.5)
	})
	// Record.Map keeps the nested Result as it is.
	gtest.C(t, func(t *gtest.T) {
		_, ok := item.Map()["track_list"].(Result)
		t.Assert(ok, true)
	})
}

// Test_Issue4787 tests that Result.Structs fills map slices with the record
// values, the same as Result.List, instead of leaving empty maps.
// See https://github.com/gogf/gf/issues/4787
func Test_Issue4787(t *testing.T) {
	var result = Result{
		{"id": gvar.New(1), "name": gvar.New("john")},
		{"id": gvar.New(2), "name": gvar.New("smith")},
	}
	gtest.C(t, func(t *gtest.T) {
		var list []map[string]any
		t.AssertNil(result.Structs(&list))
		t.Assert(len(list), 2)
		t.AssertEQ(list[0]["id"], 1)
		t.AssertEQ(list[0]["name"], "john")
		t.AssertEQ(list[1]["id"], 2)
		t.AssertEQ(list[1]["name"], "smith")
	})
	gtest.C(t, func(t *gtest.T) {
		var list []*map[string]any
		t.AssertNil(result.Structs(&list))
		t.Assert(len(list), 2)
		t.AssertEQ((*list[0])["name"], "john")
		t.AssertEQ((*list[1])["name"], "smith")
	})
	gtest.C(t, func(t *gtest.T) {
		var list []map[string]string
		t.AssertNil(result.Structs(&list))
		t.Assert(len(list), 2)
		t.AssertEQ(list[0]["id"], "1")
		t.AssertEQ(list[1]["name"], "smith")
	})
}
