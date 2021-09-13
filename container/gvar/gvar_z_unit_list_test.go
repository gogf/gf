// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gvar_test

import (
	"github.com/gogf/gf/container/gvar"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/test/gtest"
	"testing"
)

func TestVar_ListItemValues_Map(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		listMap := g.List{
			g.Map{"id": 1, "score": 100},
			g.Map{"id": 2, "score": 99},
			g.Map{"id": 3, "score": 99},
		}
		t.Assert(gvar.New(listMap).ListItemValues("id"), g.Slice{1, 2, 3})
		t.Assert(gvar.New(listMap).ListItemValues("score"), g.Slice{100, 99, 99})
	})
	gtest.C(t, func(t *gtest.T) {
		listMap := g.List{
			g.Map{"id": 1, "score": 100},
			g.Map{"id": 2, "score": nil},
			g.Map{"id": 3, "score": 0},
		}
		t.Assert(gvar.New(listMap).ListItemValues("id"), g.Slice{1, 2, 3})
		t.Assert(gvar.New(listMap).ListItemValues("score"), g.Slice{100, nil, 0})
	})
}

func TestVar_ListItemValues_Struct(t *testing.T) {
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
		t.Assert(gvar.New(listStruct).ListItemValues("Id"), g.Slice{1, 2, 3})
		t.Assert(gvar.New(listStruct).ListItemValues("Score"), g.Slice{100, 99, 0})
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
		t.Assert(gvar.New(listStruct).ListItemValues("Id"), g.Slice{1, 2, 3})
		t.Assert(gvar.New(listStruct).ListItemValues("Score"), g.Slice{100, 99, 0})
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
		t.Assert(gvar.New(listStruct).ListItemValues("Id"), g.Slice{1, 2, 3})
		t.Assert(gvar.New(listStruct).ListItemValues("Score"), g.Slice{100, nil, 0})
	})
}

func TestVar_ListItemValuesUnique(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		listMap := g.List{
			g.Map{"id": 1, "score": 100},
			g.Map{"id": 2, "score": 100},
			g.Map{"id": 3, "score": 100},
			g.Map{"id": 4, "score": 100},
			g.Map{"id": 5, "score": 100},
		}
		t.Assert(gvar.New(listMap).ListItemValuesUnique("id"), g.Slice{1, 2, 3, 4, 5})
		t.Assert(gvar.New(listMap).ListItemValuesUnique("score"), g.Slice{100})
	})
	gtest.C(t, func(t *gtest.T) {
		listMap := g.List{
			g.Map{"id": 1, "score": 100},
			g.Map{"id": 2, "score": 100},
			g.Map{"id": 3, "score": 100},
			g.Map{"id": 4, "score": 100},
			g.Map{"id": 5, "score": 99},
		}
		t.Assert(gvar.New(listMap).ListItemValuesUnique("id"), g.Slice{1, 2, 3, 4, 5})
		t.Assert(gvar.New(listMap).ListItemValuesUnique("score"), g.Slice{100, 99})
	})
	gtest.C(t, func(t *gtest.T) {
		listMap := g.List{
			g.Map{"id": 1, "score": 100},
			g.Map{"id": 2, "score": 100},
			g.Map{"id": 3, "score": 0},
			g.Map{"id": 4, "score": 100},
			g.Map{"id": 5, "score": 99},
		}
		t.Assert(gvar.New(listMap).ListItemValuesUnique("id"), g.Slice{1, 2, 3, 4, 5})
		t.Assert(gvar.New(listMap).ListItemValuesUnique("score"), g.Slice{100, 0, 99})
	})
}
