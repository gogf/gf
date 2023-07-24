// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gutil_test

import (
	"testing"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/gutil"
)

func Test_ListItemValues_Map(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		listMap := g.List{
			g.Map{"id": 1, "score": 100},
			g.Map{"id": 2, "score": 99},
			g.Map{"id": 3, "score": 99},
		}
		t.Assert(gutil.ListItemValues(listMap, "id"), g.Slice{1, 2, 3})
		t.Assert(gutil.ListItemValues(&listMap, "id"), g.Slice{1, 2, 3})
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
	gtest.C(t, func(t *gtest.T) {
		listMap := g.List{}
		t.Assert(len(gutil.ListItemValues(listMap, "id")), 0)
	})
}

func Test_ListItemValues_Map_SubKey(t *testing.T) {
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

func Test_ListItemValues_Map_Array_SubKey(t *testing.T) {
	type Scores struct {
		Math    int
		English int
	}
	gtest.C(t, func(t *gtest.T) {
		listMap := g.List{
			g.Map{"id": 1, "scores": []Scores{{1, 2}, {3, 4}}},
			g.Map{"id": 2, "scores": []Scores{{5, 6}, {7, 8}}},
			g.Map{"id": 3, "scores": []Scores{{9, 10}, {11, 12}}},
		}
		t.Assert(gutil.ListItemValues(listMap, "scores", "Math"), g.Slice{1, 3, 5, 7, 9, 11})
		t.Assert(gutil.ListItemValues(listMap, "scores", "English"), g.Slice{2, 4, 6, 8, 10, 12})
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

func Test_ListItemValues_Struct_SubKey(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		type Student struct {
			Id    int
			Score float64
		}
		type Class struct {
			Total    int
			Students []Student
		}
		listStruct := g.Slice{
			Class{2, []Student{{1, 1}, {2, 2}}},
			Class{3, []Student{{3, 3}, {4, 4}, {5, 5}}},
			Class{1, []Student{{6, 6}}},
		}
		t.Assert(gutil.ListItemValues(listStruct, "Total"), g.Slice{2, 3, 1})
		t.Assert(gutil.ListItemValues(listStruct, "Students"), `[[{"Id":1,"Score":1},{"Id":2,"Score":2}],[{"Id":3,"Score":3},{"Id":4,"Score":4},{"Id":5,"Score":5}],[{"Id":6,"Score":6}]]`)
		t.Assert(gutil.ListItemValues(listStruct, "Students", "Id"), g.Slice{1, 2, 3, 4, 5, 6})
	})
	gtest.C(t, func(t *gtest.T) {
		type Student struct {
			Id    int
			Score float64
		}
		type Class struct {
			Total    int
			Students []*Student
		}
		listStruct := g.Slice{
			&Class{2, []*Student{{1, 1}, {2, 2}}},
			&Class{3, []*Student{{3, 3}, {4, 4}, {5, 5}}},
			&Class{1, []*Student{{6, 6}}},
		}
		t.Assert(gutil.ListItemValues(listStruct, "Total"), g.Slice{2, 3, 1})
		t.Assert(gutil.ListItemValues(listStruct, "Students"), `[[{"Id":1,"Score":1},{"Id":2,"Score":2}],[{"Id":3,"Score":3},{"Id":4,"Score":4},{"Id":5,"Score":5}],[{"Id":6,"Score":6}]]`)
		t.Assert(gutil.ListItemValues(listStruct, "Students", "Id"), g.Slice{1, 2, 3, 4, 5, 6})
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

func Test_ListItemValuesUnique_Struct_SubKey(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		type Student struct {
			Id    int
			Score float64
		}
		type Class struct {
			Total    int
			Students []Student
		}
		listStruct := g.Slice{
			Class{2, []Student{{1, 1}, {1, 2}}},
			Class{3, []Student{{2, 3}, {2, 4}, {5, 5}}},
			Class{1, []Student{{6, 6}}},
		}
		t.Assert(gutil.ListItemValuesUnique(listStruct, "Total"), g.Slice{2, 3, 1})
		t.Assert(gutil.ListItemValuesUnique(listStruct, "Students", "Id"), g.Slice{1, 2, 5, 6})
	})
	gtest.C(t, func(t *gtest.T) {
		type Student struct {
			Id    int
			Score float64
		}
		type Class struct {
			Total    int
			Students []*Student
		}
		listStruct := g.Slice{
			&Class{2, []*Student{{1, 1}, {1, 2}}},
			&Class{3, []*Student{{2, 3}, {2, 4}, {5, 5}}},
			&Class{1, []*Student{{6, 6}}},
		}
		t.Assert(gutil.ListItemValuesUnique(listStruct, "Total"), g.Slice{2, 3, 1})
		t.Assert(gutil.ListItemValuesUnique(listStruct, "Students", "Id"), g.Slice{1, 2, 5, 6})
	})
}

func Test_ListItemValuesUnique_Map_Array_SubKey(t *testing.T) {
	type Scores struct {
		Math    int
		English int
	}
	gtest.C(t, func(t *gtest.T) {
		listMap := g.List{
			g.Map{"id": 1, "scores": []Scores{{1, 2}, {1, 2}}},
			g.Map{"id": 2, "scores": []Scores{{5, 8}, {5, 8}}},
			g.Map{"id": 3, "scores": []Scores{{9, 10}, {11, 12}}},
		}
		t.Assert(gutil.ListItemValuesUnique(listMap, "scores", "Math"), g.Slice{1, 5, 9, 11})
		t.Assert(gutil.ListItemValuesUnique(listMap, "scores", "English"), g.Slice{2, 8, 10, 12})
		t.Assert(gutil.ListItemValuesUnique(listMap, "scores", "PE"), g.Slice{})
	})
}
