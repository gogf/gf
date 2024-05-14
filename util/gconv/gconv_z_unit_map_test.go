// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gconv_test

import (
	"reflect"
	"testing"

	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/gconv"
)

type SubMapTest struct {
	Name string
}

var mapTests = []struct {
	value  interface{}
	expect map[string]interface{}
}{
	{map[int]int{1: 1}, map[string]interface{}{"1": 1}},
	{map[float64]int{1.1: 1}, map[string]interface{}{"1.1": 1}},
	{map[string]int{"k1": 1}, map[string]interface{}{"k1": 1}},

	{map[string]int{"k1": 1}, map[string]interface{}{"k1": 1}},
	{map[string]float64{"k1": 1.1}, map[string]interface{}{"k1": 1.1}},
	{map[string]string{"k1": "v1"}, map[string]interface{}{"k1": "v1"}},

	{`{"earth": "亚马逊雨林"}`,
		map[string]interface{}{"earth": "亚马逊雨林"}},
	{[]byte(`{"earth": "撒哈拉沙漠"}`),
		map[string]interface{}{"earth": "撒哈拉沙漠"}},

	{&struct {
		Earth string
	}{
		Earth: "大峡谷",
	}, map[string]interface{}{"Earth": "大峡谷"}},

	{struct {
		Earth string
	}{
		Earth: "马里亚纳海沟",
	}, map[string]interface{}{"Earth": "马里亚纳海沟"}},

	{struct {
		Earth string
		mars  string
	}{
		Earth: "大堡礁",
		mars:  "奥林帕斯山",
	}, map[string]interface{}{"Earth": "大堡礁"}},

	{struct {
		Earth string
		SubMapTest
	}{
		Earth: "中国",
		SubMapTest: SubMapTest{
			Name: "长江",
		},
	}, map[string]interface{}{"Earth": "中国", "Name": "长江"}},

	{struct {
		Earth string
		China SubMapTest
	}{
		Earth: "中国",
		China: SubMapTest{
			Name: "黄河",
		},
	}, map[string]interface{}{"Earth": "中国", "China": map[string]interface{}{"Name": "黄河"}}},

	{struct {
		China         string `c:"中国"`
		America       string `c:"-"`
		UnitedKingdom string `c:"UK,omitempty"`
	}{
		China:         "长城",
		America:       "Statue of Liberty",
		UnitedKingdom: "",
	}, map[string]interface{}{"中国": "长城", "UK": ""}},

	{struct {
		China         string `gconv:"中国"`
		America       string `gconv:"-"`
		UnitedKingdom string `c:"UK,omitempty"`
	}{
		China:         "故宫",
		America:       "White House",
		UnitedKingdom: "",
	}, map[string]interface{}{"中国": "故宫", "UK": ""}},

	{struct {
		China         string `json:"中国"`
		America       string `json:"-"`
		UnitedKingdom string `json:"UK,omitempty"`
	}{
		China:         "东方明珠",
		America:       "Empire State Building",
		UnitedKingdom: "",
	}, map[string]interface{}{"中国": "东方明珠", "UK": ""}},

	{[]int{1, 2, 3}, map[string]interface{}{"1": 2, "3": nil}},
	{[]int{1, 2, 3, 4}, map[string]interface{}{"1": 2, "3": 4}},
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
				maps = []interface{}{
					test.value,
					test.value,
				}
				expects = []interface{}{
					test.expect,
					test.expect,
				}
			)
			t.Assert(gconv.Maps(maps), expects)
			t.Assert(gconv.SliceMap(maps), expects)
		}
	})
}

func TestMapStrStr(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		for _, test := range mapTests {
			for k, v := range test.expect {
				test.expect[k] = gconv.String(v)
			}
			t.Assert(gconv.MapStrStr(test.value), test.expect)
		}
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
			Mercury interface{} `json:",omitempty"`
		}{
			Earth:   "死海",
			Venus:   0,
			Mars:    "",
			Mercury: nil,
		}
		r := gconv.Map(testMapOmitEmpty, gconv.MapOption{OmitEmpty: true})
		t.Assert(r, map[string]interface{}{"Earth": "死海"})
	})

	// Test for option: Tags.
	gtest.C(t, func(t *gtest.T) {
		var testMapOmitEmpty = struct {
			Earth string `gconv:"errEarth" chinese:"地球" french:"Terre"`
		}{
			Earth: "尼莫点",
		}
		c := gconv.Map(testMapOmitEmpty, gconv.MapOption{Tags: []string{"chinese", "french"}})
		t.Assert(c, map[string]interface{}{"地球": "尼莫点"})

		f := gconv.Map(testMapOmitEmpty, gconv.MapOption{Tags: []string{"french", "chinese"}})
		t.Assert(f, map[string]interface{}{"Terre": "尼莫点"})
	})
}

// See gconv_test.TestScan for more.
func TestMapToMap(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			err    error
			value  = map[string]string{"k1": "v1"}
			expect = make(map[string]interface{})
		)
		err = gconv.MapToMap(value, &expect)
		t.Assert(err, nil)
		t.Assert(value["k1"], expect["k1"])
	})
}
