// Copyright 2019 gf Author(https://github.com/jin502437344/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/jin502437344/gf.

package gutil_test

import (
	"github.com/jin502437344/gf/frame/g"
	"testing"

	"github.com/jin502437344/gf/test/gtest"
	"github.com/jin502437344/gf/util/gutil"
)

func Test_ListItemValues_Map(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		listMap := g.List{
			g.Map{"id": 1, "score": 100},
			g.Map{"id": 2, "score": 99},
			g.Map{"id": 3, "score": 99},
		}
		t.Assert(gutil.ListItemValues(listMap, "id"), g.Slice{1, 2, 3})
		t.Assert(gutil.ListItemValues(listMap, "score"), g.Slice{100, 99, 99})
	})
	gtest.C(t, func(t *gtest.T) {
		listMap := g.List{
			g.Map{"id": 1, "score": 100},
			g.Map{"id": 2, "score": nil},
			g.Map{"id": 3, "score": 0},
		}
		t.Assert(gutil.ListItemValues(listMap, "id"), g.Slice{1, 2, 3})
		t.Assert(gutil.ListItemValues(listMap, "score"), g.Slice{100, nil, 0})
	})
}

func Test_ListItemValues_SubKey(t *testing.T) {
	type Scores struct {
		Math    int
		English int
	}
	gtest.C(t, func(t *gtest.T) {
		listMap := g.List{
			g.Map{"id": 1, "scores": Scores{100, 60}},
			g.Map{"id": 2, "scores": Scores{0, 100}},
			g.Map{"id": 3, "scores": Scores{59, 99}},
		}
		t.Assert(gutil.ListItemValues(listMap, "scores", "Math"), g.Slice{100, 0, 59})
		t.Assert(gutil.ListItemValues(listMap, "scores", "English"), g.Slice{60, 100, 99})
		t.Assert(gutil.ListItemValues(listMap, "scores", "PE"), g.Slice{})
	})
}

func Test_ListItemValues_Struct(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		type T struct {
			Id    int
			Score float64
		}
		listStruct := g.Slice{
			T{1, 100},
			T{2, 99},
			T{3, 0},
		}
		t.Assert(gutil.ListItemValues(listStruct, "Id"), g.Slice{1, 2, 3})
		t.Assert(gutil.ListItemValues(listStruct, "Score"), g.Slice{100, 99, 0})
	})
	// Pointer items.
	gtest.C(t, func(t *gtest.T) {
		type T struct {
			Id    int
			Score float64
		}
		listStruct := g.Slice{
			&T{1, 100},
			&T{2, 99},
			&T{3, 0},
		}
		t.Assert(gutil.ListItemValues(listStruct, "Id"), g.Slice{1, 2, 3})
		t.Assert(gutil.ListItemValues(listStruct, "Score"), g.Slice{100, 99, 0})
	})
	// Nil element value.
	gtest.C(t, func(t *gtest.T) {
		type T struct {
			Id    int
			Score interface{}
		}
		listStruct := g.Slice{
			T{1, 100},
			T{2, nil},
			T{3, 0},
		}
		t.Assert(gutil.ListItemValues(listStruct, "Id"), g.Slice{1, 2, 3})
		t.Assert(gutil.ListItemValues(listStruct, "Score"), g.Slice{100, nil, 0})
	})
}

func Test_ListItemValuesUnique(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		listMap := g.List{
			g.Map{"id": 1, "score": 100},
			g.Map{"id": 2, "score": 100},
			g.Map{"id": 3, "score": 100},
			g.Map{"id": 4, "score": 100},
			g.Map{"id": 5, "score": 100},
		}
		t.Assert(gutil.ListItemValuesUnique(listMap, "id"), g.Slice{1, 2, 3, 4, 5})
		t.Assert(gutil.ListItemValuesUnique(listMap, "score"), g.Slice{100})
	})
	gtest.C(t, func(t *gtest.T) {
		listMap := g.List{
			g.Map{"id": 1, "score": 100},
			g.Map{"id": 2, "score": 100},
			g.Map{"id": 3, "score": 100},
			g.Map{"id": 4, "score": 100},
			g.Map{"id": 5, "score": 99},
		}
		t.Assert(gutil.ListItemValuesUnique(listMap, "id"), g.Slice{1, 2, 3, 4, 5})
		t.Assert(gutil.ListItemValuesUnique(listMap, "score"), g.Slice{100, 99})
	})
	gtest.C(t, func(t *gtest.T) {
		listMap := g.List{
			g.Map{"id": 1, "score": 100},
			g.Map{"id": 2, "score": 100},
			g.Map{"id": 3, "score": 0},
			g.Map{"id": 4, "score": 100},
			g.Map{"id": 5, "score": 99},
		}
		t.Assert(gutil.ListItemValuesUnique(listMap, "id"), g.Slice{1, 2, 3, 4, 5})
		t.Assert(gutil.ListItemValuesUnique(listMap, "score"), g.Slice{100, 0, 99})
	})
}
