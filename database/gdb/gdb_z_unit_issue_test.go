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
