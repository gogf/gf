// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gjson_test

import (
	"encoding/json"
	"github.com/gogf/gf/util/gconv"
	"testing"

	"github.com/gogf/gf/encoding/gjson"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/test/gtest"
)

func TestJson_UnmarshalJSON(t *testing.T) {
	data := []byte(`{"n":123456789, "m":{"k":"v"}, "a":[1,2,3]}`)
	gtest.Case(t, func() {
		j := gjson.New(nil)
		err := json.Unmarshal(data, j)
		gtest.Assert(err, nil)
		gtest.Assert(j.Get("n"), "123456789")
		gtest.Assert(j.Get("m"), g.Map{"k": "v"})
		gtest.Assert(j.Get("m.k"), "v")
		gtest.Assert(j.Get("a"), g.Slice{1, 2, 3})
		gtest.Assert(j.Get("a.1"), 2)
	})
}

func TestJson_UnmarshalValue(t *testing.T) {
	type T struct {
		Name string
		Json *gjson.Json
	}
	// JSON
	gtest.Case(t, func() {
		var t *T
		err := gconv.Struct(g.Map{
			"name": "john",
			"json": []byte(`{"n":123456789, "m":{"k":"v"}, "a":[1,2,3]}`),
		}, &t)
		gtest.Assert(err, nil)
		gtest.Assert(t.Name, "john")
		gtest.Assert(t.Json.Get("n"), "123456789")
		gtest.Assert(t.Json.Get("m"), g.Map{"k": "v"})
		gtest.Assert(t.Json.Get("m.k"), "v")
		gtest.Assert(t.Json.Get("a"), g.Slice{1, 2, 3})
		gtest.Assert(t.Json.Get("a.1"), 2)
	})
	// Map
	gtest.Case(t, func() {
		var t *T
		err := gconv.Struct(g.Map{
			"name": "john",
			"json": g.Map{
				"n": 123456789,
				"m": g.Map{"k": "v"},
				"a": g.Slice{1, 2, 3},
			},
		}, &t)
		gtest.Assert(err, nil)
		gtest.Assert(t.Name, "john")
		gtest.Assert(t.Json.Get("n"), "123456789")
		gtest.Assert(t.Json.Get("m"), g.Map{"k": "v"})
		gtest.Assert(t.Json.Get("m.k"), "v")
		gtest.Assert(t.Json.Get("a"), g.Slice{1, 2, 3})
		gtest.Assert(t.Json.Get("a.1"), 2)
	})
}
