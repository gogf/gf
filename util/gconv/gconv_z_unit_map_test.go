// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gconv_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/gconv"
)

type SubMapTest struct {
	Name string
}

var mapTests = []struct {
	value  any
	expect any
}{
	{map[string]int{"k1": 1}, map[string]any{"k1": 1}},
	{map[string]uint{"k1": 1}, map[string]any{"k1": 1}},
	{map[string]string{"k1": "v1"}, map[string]any{"k1": "v1"}},
	{map[string]float32{"k1": 1.1}, map[string]any{"k1": 1.1}},
	{map[string]float64{"k1": 1.1}, map[string]any{"k1": 1.1}},
	{map[string]bool{"k1": true}, map[string]any{"k1": true}},
	{map[string]any{"k1": "v1"}, map[string]any{"k1": "v1"}},

	{map[any]int{"k1": 1}, map[string]any{"k1": 1}},
	{map[any]uint{"k1": 1}, map[string]any{"k1": 1}},
	{map[any]string{"k1": "v1"}, map[string]any{"k1": "v1"}},
	{map[any]float32{"k1": 1.1}, map[string]any{"k1": 1.1}},
	{map[any]float64{"k1": 1.1}, map[string]any{"k1": 1.1}},
	{map[any]bool{"k1": true}, map[string]any{"k1": true}},
	{map[any]any{"k1": "v1"}, map[string]any{"k1": "v1"}},

	{map[int]int{1: 1}, map[string]any{"1": 1}},
	{map[int]string{1: "v1"}, map[string]any{"1": "v1"}},
	{map[uint]int{1: 1}, map[string]any{"1": 1}},
	{map[uint]string{1: "v1"}, map[string]any{"1": "v1"}},

	{[]int{1, 2, 3}, map[string]any{"1": 2, "3": nil}},
	{[]int{1, 2, 3, 4}, map[string]any{"1": 2, "3": 4}},

	{`{"earth": "亚马逊雨林"}`,
		map[string]any{"earth": "亚马逊雨林"}},
	{[]byte(`{"earth": "撒哈拉沙漠"}`),
		map[string]any{"earth": "撒哈拉沙漠"}},
	{`{Earth}`, nil},

	{"", nil},
	{[]byte(""), nil},
	{`"{earth亚马逊雨林}`, nil},
	{[]byte(`{earth撒哈拉沙漠}`), nil},
	{[]byte(`{Earth}`), nil},

	{nil, nil},

	{&struct {
		Earth string
	}{
		Earth: "大峡谷",
	}, map[string]any{"Earth": "大峡谷"}},

	{struct {
		Earth string
	}{
		Earth: "马里亚纳海沟",
	}, map[string]any{"Earth": "马里亚纳海沟"}},

	{struct {
		Earth string
		mars  string
	}{
		Earth: "大堡礁",
		mars:  "奥林帕斯山",
	}, map[string]any{"Earth": "大堡礁"}},

	{struct {
		Earth string
		SubMapTest
	}{
		Earth: "中国",
		SubMapTest: SubMapTest{
			Name: "长江",
		},
	}, map[string]any{"Earth": "中国", "Name": "长江"}},

	{struct {
		Earth string
		China SubMapTest
	}{
		Earth: "中国",
		China: SubMapTest{
			Name: "黄河",
		},
	}, map[string]any{"Earth": "中国", "China": map[string]any{"Name": "黄河"}}},

	{struct {
		Earth      string
		SubMapTest `json:"sub_map_test"`
	}{
		Earth: "中国",
		SubMapTest: SubMapTest{
			Name: "淮河",
		},
	}, map[string]any{"Earth": "中国", "sub_map_test": map[string]any{"Name": "淮河"}}},

	{struct {
		Earth string
		China SubMapTest `gconv:"中国"`
	}{
		Earth: "中国",
		China: SubMapTest{
			Name: "黄河",
		},
	}, map[string]any{"Earth": "中国", "中国": map[string]any{"Name": "黄河"}}},

	{struct {
		China         string `c:"中国"`
		America       string `c:"-"`
		UnitedKingdom string `c:"UK,omitempty"`
	}{
		China:         "长城",
		America:       "Statue of Liberty",
		UnitedKingdom: "",
	}, map[string]any{"中国": "长城", "UK": ""}},

	{struct {
		China         string `gconv:"中国"`
		America       string `gconv:"-"`
		UnitedKingdom string `c:"UK,omitempty"`
	}{
		China:         "故宫",
		America:       "White House",
		UnitedKingdom: "",
	}, map[string]any{"中国": "故宫", "UK": ""}},

	{struct {
		China         string `json:"中国"`
		America       string `json:"-"`
		UnitedKingdom string `json:"UK,omitempty"`
	}{
		China:         "东方明珠",
		America:       "Empire State Building",
		UnitedKingdom: "",
	}, map[string]any{"中国": "东方明珠", "UK": ""}},

	{struct {
		China   any `json:",omitempty"`
		America string      `json:",omitempty"`
	}{
		China:   "黄山",
		America: "",
	}, map[string]any{"China": "黄山", "America": ""}},
}

func TestMap(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		for _, test := range mapTests {
			t.Assert(gconv.Map(test.value), test.expect)
		}
	})
}

func TestMaps(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		for _, test := range mapTests {
			var (
				maps    any
				expects any
			)

			if v, ok := test.value.(string); ok {
				maps = fmt.Sprintf(`[%s,%s]`, v, v)
			} else if v, ok := test.value.([]byte); ok {
				maps = []byte(fmt.Sprintf(`[%s,%s]`, v, v))
			} else if test.value == nil {
				maps = nil
			} else {
				maps = []any{
					test.value,
					test.value,
				}
			}

			if test.expect == nil {
				expects = test.expect
			} else {
				expects = []any{
					test.expect,
					test.expect,
				}
			}
			t.Assert(gconv.Maps(maps), expects)

			// The following is the same as gconv.Maps.
			t.Assert(gconv.MapsDeep(maps), expects)
			t.Assert(gconv.SliceMap(maps), expects)
			t.Assert(gconv.SliceMapDeep(maps), expects)
		}
	})

	// Test for special types.
	gtest.C(t, func(t *gtest.T) {
		mapStrAny := []map[string]any{
			{"earth": "亚马逊雨林"},
			{"mars": "奥林帕斯山"},
		}
		t.Assert(gconv.Maps(mapStrAny), mapStrAny)

		mapEmpty := []map[string]string{}
		t.AssertNil(gconv.Maps(mapEmpty))

		t.Assert(gconv.Maps(`test`), nil)
		t.Assert(gconv.Maps([]byte(`test`)), nil)
	})
}

func TestMapsDeepExtra(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		type s struct {
			Earth g.Map `c:"earth_map"`
		}

		t.Assert(gconv.MapDeep(&s{
			Earth: g.Map{
				"sea_num": 4,
				"one_sea": g.Map{
					"sea_name": "太平洋",
				},
				"map_sat": g.MapAnyAny{
					1:         "Arctic",
					"Pacific": 2,
					"Indian":  "印度洋",
				},
			},
		}), g.Map{
			"earth_map": g.Map{
				"sea_num": 4,
				"one_sea": g.Map{
					"sea_name": "太平洋",
				},
				"map_sat": g.Map{
					"1":       "Arctic",
					"Pacific": 2,
					"Indian":  "印度洋",
				},
			},
		})
	})

	gtest.C(t, func(t *gtest.T) {
		t.Assert(gconv.MapsDeep(`test`), nil)
		t.Assert(gconv.MapsDeep([]byte(`test`)), nil)
	})
}

func TestMapStrStr(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		for _, test := range mapTests {
			var expect map[string]any
			if v, ok := test.expect.(map[string]any); ok {
				expect = v
			}
			for k, v := range expect {
				expect[k] = gconv.String(v)
			}
			t.Assert(gconv.MapStrStr(test.value), test.expect)
		}
	})
}

func TestMapStrStrDeepExtra(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.Assert(gconv.MapStrStrDeep(map[string]string{"mars": "Syrtis"}), map[string]string{"mars": "Syrtis"})
		t.Assert(gconv.MapStrStrDeep(`{}`), nil)
	})
}

func TestMapWithMapOption(t *testing.T) {
	// Test for option: Deep.
	gtest.C(t, func(t *gtest.T) {
		var testMapDeep = struct {
			Earth      string
			SubMapTest SubMapTest
		}{
			Earth: "中国",
			SubMapTest: SubMapTest{
				Name: "黄山",
			},
		}
		var (
			dt  = gconv.Map(testMapDeep, gconv.MapOption{Deep: true})
			df  = gconv.Map(testMapDeep, gconv.MapOption{Deep: false})
			dtk = reflect.TypeOf(dt["SubMapTest"]).Kind()
			dfk = reflect.TypeOf(df["SubMapTest"]).Kind()
		)
		t.AssertNE(dtk, dfk)
	})

	// Test for option: OmitEmpty.
	gtest.C(t, func(t *gtest.T) {
		var testMapOmitEmpty = struct {
			Earth   string
			Venus   int         `gconv:",omitempty"`
			Mars    string      `c:",omitempty"`
			Mercury any `json:",omitempty"`
		}{
			Earth:   "死海",
			Venus:   0,
			Mars:    "",
			Mercury: nil,
		}
		r := gconv.Map(testMapOmitEmpty, gconv.MapOption{OmitEmpty: true})
		t.Assert(r, map[string]any{"Earth": "死海"})
	})

	// Test for option: Tags.
	gtest.C(t, func(t *gtest.T) {
		var testMapOmitEmpty = struct {
			Earth string `gconv:"errEarth" chinese:"地球" french:"Terre"`
		}{
			Earth: "尼莫点",
		}
		c := gconv.Map(testMapOmitEmpty, gconv.MapOption{Tags: []string{"chinese", "french"}})
		t.Assert(c, map[string]any{"地球": "尼莫点"})

		f := gconv.Map(testMapOmitEmpty, gconv.MapOption{Tags: []string{"french", "chinese"}})
		t.Assert(f, map[string]any{"Terre": "尼莫点"})
	})
}

func TestMapToMapExtra(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			err    error
			value  = map[string]string{"k1": "v1"}
			expect = make(map[string]any)
		)
		err = gconv.MapToMap(value, &expect)
		t.AssertNil(err)
		t.Assert(value["k1"], expect["k1"])
	})

	gtest.C(t, func(t *gtest.T) {
		v := g.Map{
			"k": g.Map{
				"name": "Earth",
			},
		}
		e := make(map[string]SubMapTest)
		err := gconv.MapToMap(v, &e)
		t.AssertNil(err)
		t.Assert(len(e), 1)
		t.Assert(e["k"].Name, "Earth")
	})
}

func TestMaptoMapsExtra(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		v := g.Slice{
			g.Map{"id": 1, "name": "john"},
			g.Map{"id": 2, "name": "smith"},
		}
		var e []*g.Map
		err := gconv.MapToMaps(v, &e)
		t.AssertNil(err)
		t.Assert(len(v), 2)
		t.Assert(v, e)
	})
}
