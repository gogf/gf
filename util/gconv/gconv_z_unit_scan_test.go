// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gconv_test

import (
	"fmt"
	"testing"

	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/gconv"
)

type scanStructTest struct {
	Name  string
	Place string
}

type scanExpectTest struct {
	mapStrStr map[string]string
	mapStrAny map[string]interface{}
	mapAnyAny map[interface{}]interface{}

	structSub    scanStructTest
	structSubPtr *scanStructTest
}

var scanValueMapsTest = []map[string]interface{}{
	{"Name": false, "Place": true},
	{"Name": int(0), "Place": int(1)},
	{"Name": int8(0), "Place": int8(1)},
	{"Name": int16(0), "Place": int16(1)},
	{"Name": int32(0), "Place": int32(1)},
	{"Name": int64(0), "Place": int64(1)},
	{"Name": uint(0), "Place": uint(1)},
	{"Name": uint8(0), "Place": uint8(1)},
	{"Name": uint16(0), "Place": uint16(1)},
	{"Name": uint32(0), "Place": uint32(1)},
	{"Name": uint64(0), "Place": uint64(1)},
	{"Name": float32(0), "Place": float32(1)},
	{"Name": float64(0), "Place": float64(1)},
	{"Name": "Mercury", "Place": "卡罗利斯盆地"},
	{"Name": []byte("Saturn"), "Place": []byte("土星环")},
	{"Name": complex64(0), "Place": complex64(1 + 2i)},
	{"Name": complex128(0), "Place": complex128(1 + 2i)},
	{"Name": interface{}(0), "Place": interface{}("1")},
	{"Name": gvar.New("Jupiter"), "Place": gvar.New("大红斑")},
	{"Name": gtime.New("2024-01-01 01:01:01"), "Place": gtime.New("2021-01-01 01:01:01")},
	{"Name": map[string]string{"Name": "Sun"}, "Place": map[string]string{"Place": "太阳黑子"}},
	{"Name": []string{"Earth", "Moon"}, "Place": []string{"好望角", "万户环形山"}},
}

var scanValueStructsTest = []scanStructTest{
	{"Venus", "阿佛洛狄特高原"},
}

var scanValueJsonTest = []string{
	`{"Name": "Mars", "Place": "奥林帕斯山"}`,
}

var scanExpects = scanExpectTest{
	mapStrStr: make(map[string]string),
	mapStrAny: make(map[string]interface{}),
	mapAnyAny: make(map[interface{}]interface{}),

	structSub:    scanStructTest{},
	structSubPtr: &scanStructTest{},
}

func TestScan(t *testing.T) {
	// Test for map converting.
	gtest.C(t, func(t *gtest.T) {
		scanValuesTest := scanValueMapsTest
		for _, test := range scanValuesTest {
			var (
				err         error
				scanExpects = scanExpects
			)

			err = gconv.Scan(test, &scanExpects.mapStrStr)
			t.AssertNil(err)
			t.Assert(test["Name"], scanExpects.mapStrStr["Name"])
			t.Assert(test["Place"], scanExpects.mapStrStr["Place"])

			err = gconv.Scan(test, &scanExpects.mapStrAny)
			t.AssertNil(err)
			t.Assert(test["Name"], scanExpects.mapStrAny["Name"])
			t.Assert(test["Place"], scanExpects.mapStrAny["Place"])

			err = gconv.Scan(test, &scanExpects.mapAnyAny)
			t.AssertNil(err)
			t.Assert(test["Name"], scanExpects.mapAnyAny["Name"])
			t.Assert(test["Place"], scanExpects.mapAnyAny["Place"])

			err = gconv.Scan(test, &scanExpects.structSub)
			t.AssertNil(err)
			t.Assert(test["Name"], scanExpects.structSub.Name)
			t.Assert(test["Place"], scanExpects.structSub.Place)

			err = gconv.Scan(test, &scanExpects.structSubPtr)
			t.AssertNil(err)
			t.Assert(test["Name"], scanExpects.structSubPtr.Name)
			t.Assert(test["Place"], scanExpects.structSubPtr.Place)

		}
	})

	// Test for slice map converting.
	gtest.C(t, func(t *gtest.T) {
		scanValuesTest := scanValueMapsTest
		for _, test := range scanValuesTest {
			var (
				err         error
				scanExpects = scanExpects
				maps        = []map[string]interface{}{test, test}
			)

			var mss = []map[string]string{scanExpects.mapStrStr, scanExpects.mapStrStr}
			err = gconv.Scan(maps, &mss)
			t.AssertNil(err)
			t.Assert(len(mss), len(maps))
			for k, _ := range maps {
				t.Assert(maps[k]["Name"], mss[k]["Name"])
				t.Assert(maps[k]["Place"], mss[k]["Place"])
			}

			var msa = []map[string]interface{}{scanExpects.mapStrAny, scanExpects.mapStrAny}
			err = gconv.Scan(maps, &msa)
			t.AssertNil(err)
			t.Assert(len(msa), len(maps))
			for k, _ := range maps {
				t.Assert(maps[k]["Name"], msa[k]["Name"])
				t.Assert(maps[k]["Place"], msa[k]["Place"])
			}

			var maa = []map[interface{}]interface{}{scanExpects.mapAnyAny, scanExpects.mapAnyAny}
			err = gconv.Scan(maps, &maa)
			t.AssertNil(err)
			t.Assert(len(maa), len(maps))
			for k, _ := range maps {
				t.Assert(maps[k]["Name"], maa[k]["Name"])
				t.Assert(maps[k]["Place"], maa[k]["Place"])
			}

			var ss = []scanStructTest{scanExpects.structSub, scanExpects.structSub}
			err = gconv.Scan(maps, &ss)
			t.AssertNil(err)
			t.Assert(len(ss), len(maps))
			for k, _ := range maps {
				t.Assert(maps[k]["Name"], ss[k].Name)
				t.Assert(maps[k]["Place"], ss[k].Place)
			}

			var ssp = []*scanStructTest{scanExpects.structSubPtr, scanExpects.structSubPtr}
			err = gconv.Scan(maps, &ssp)
			t.AssertNil(err)
			t.Assert(len(ssp), len(maps))
			for k, _ := range maps {
				t.Assert(maps[k]["Name"], ssp[k].Name)
				t.Assert(maps[k]["Place"], ssp[k].Place)
			}
		}
	})

	// Test for struct converting.
	gtest.C(t, func(t *gtest.T) {
		scanValuesTest := scanValueStructsTest
		for _, test := range scanValuesTest {
			var (
				err         error
				scanExpects = scanExpects
			)

			err = gconv.Scan(test, &scanExpects.mapStrStr)
			t.AssertNil(err)
			t.Assert(test.Name, scanExpects.mapStrStr["Name"])
			t.Assert(test.Place, scanExpects.mapStrStr["Place"])

			err = gconv.Scan(test, &scanExpects.mapStrAny)
			t.AssertNil(err)
			t.Assert(test.Name, scanExpects.mapStrAny["Name"])
			t.Assert(test.Place, scanExpects.mapStrAny["Place"])

			err = gconv.Scan(test, &scanExpects.mapAnyAny)
			t.AssertNil(err)
			t.Assert(test.Name, scanExpects.mapAnyAny["Name"])
			t.Assert(test.Place, scanExpects.mapAnyAny["Place"])

			err = gconv.Scan(test, &scanExpects.structSub)
			t.AssertNil(err)
			t.Assert(test.Name, scanExpects.structSub.Name)
			t.Assert(test.Place, scanExpects.structSub.Place)

			err = gconv.Scan(test, &scanExpects.structSubPtr)
			t.AssertNil(err)
			t.Assert(test.Name, scanExpects.structSubPtr.Name)
			t.Assert(test.Place, scanExpects.structSubPtr.Place)
		}
	})

	// Test for slice struct converting.
	gtest.C(t, func(t *gtest.T) {
		scanValuesTest := scanValueStructsTest
		for _, test := range scanValuesTest {
			var (
				err         error
				scanExpects = scanExpects
				structs     = []scanStructTest{test, test}
			)

			var mss = []map[string]string{scanExpects.mapStrStr, scanExpects.mapStrStr}
			err = gconv.Scan(structs, &mss)
			t.AssertNil(err)
			t.Assert(len(mss), len(structs))
			for k, _ := range structs {
				t.Assert(structs[k].Name, mss[k]["Name"])
				t.Assert(structs[k].Place, mss[k]["Place"])
			}

			var msa = []map[string]interface{}{scanExpects.mapStrAny, scanExpects.mapStrAny}
			err = gconv.Scan(structs, &msa)
			t.AssertNil(err)
			t.Assert(len(msa), len(structs))
			for k, _ := range structs {
				t.Assert(structs[k].Name, msa[k]["Name"])
				t.Assert(structs[k].Place, msa[k]["Place"])
			}

			var maa = []map[interface{}]interface{}{scanExpects.mapAnyAny, scanExpects.mapAnyAny}
			err = gconv.Scan(structs, &maa)
			t.AssertNil(err)
			t.Assert(len(maa), len(structs))
			for k, _ := range structs {
				t.Assert(structs[k].Name, maa[k]["Name"])
				t.Assert(structs[k].Place, maa[k]["Place"])
			}

			var ss = []scanStructTest{scanExpects.structSub, scanExpects.structSub}
			err = gconv.Scan(structs, &ss)
			t.AssertNil(err)
			t.Assert(len(ss), len(structs))
			for k, _ := range structs {
				t.Assert(structs[k].Name, ss[k].Name)
				t.Assert(structs[k].Place, ss[k].Place)
			}

			var ssp = []*scanStructTest{scanExpects.structSubPtr, scanExpects.structSubPtr}
			err = gconv.Scan(structs, &ssp)
			t.AssertNil(err)
			t.Assert(len(ssp), len(structs))
			for k, _ := range structs {
				t.Assert(structs[k].Name, ssp[k].Name)
				t.Assert(structs[k].Place, ssp[k].Place)
			}
		}
	})

	// Test for json converting.
	gtest.C(t, func(t *gtest.T) {
		scanValuesTest := scanValueJsonTest
		for _, test := range scanValuesTest {
			var (
				err         error
				scanExpects = scanExpects
			)

			err = gconv.Scan(test, &scanExpects.mapStrStr)
			t.AssertNil(err)
			t.Assert("Mars", scanExpects.mapStrStr["Name"])
			t.Assert("奥林帕斯山", scanExpects.mapStrStr["Place"])

			err = gconv.Scan(test, &scanExpects.mapStrAny)
			t.AssertNil(err)
			t.Assert("Mars", scanExpects.mapStrAny["Name"])
			t.Assert("奥林帕斯山", scanExpects.mapStrAny["Place"])

			err = gconv.Scan(test, &scanExpects.mapAnyAny)
			t.Assert(err, gerror.New(
				"json.UnmarshalUseNumber failed: json: cannot unmarshal object into Go value of type map[interface {}]interface {}",
			))

			err = gconv.Scan(test, &scanExpects.structSub)
			t.AssertNil(err)
			t.Assert("Mars", scanExpects.structSub.Name)
			t.Assert("奥林帕斯山", scanExpects.structSub.Place)

			err = gconv.Scan(test, &scanExpects.structSubPtr)
			t.AssertNil(err)
			t.Assert("Mars", scanExpects.structSubPtr.Name)
			t.Assert("奥林帕斯山", scanExpects.structSubPtr.Place)
		}
	})

	// Test for slice json converting.
	gtest.C(t, func(t *gtest.T) {
		scanValuesTest := scanValueJsonTest
		for _, test := range scanValuesTest {
			var (
				err         error
				scanExpects = scanExpects
				jsons       = fmt.Sprintf("[%s, %s]", test, test)
			)

			var mss = []map[string]string{scanExpects.mapStrStr, scanExpects.mapStrStr}
			err = gconv.Scan(jsons, &mss)
			t.AssertNil(err)
			t.Assert(len(mss), 2)
			for k, _ := range mss {
				t.Assert("Mars", mss[k]["Name"])
				t.Assert("奥林帕斯山", mss[k]["Place"])
			}

			var msa = []map[string]interface{}{scanExpects.mapStrAny, scanExpects.mapStrAny}
			err = gconv.Scan(jsons, &msa)
			t.AssertNil(err)
			t.Assert(len(msa), 2)
			for k, _ := range msa {
				t.Assert("Mars", msa[k]["Name"])
				t.Assert("奥林帕斯山", msa[k]["Place"])
			}

			var maa = []map[interface{}]interface{}{scanExpects.mapAnyAny, scanExpects.mapAnyAny}
			err = gconv.Scan(jsons, &maa)
			t.Assert(err, gerror.New(
				"json.UnmarshalUseNumber failed: json: cannot unmarshal object into Go value of type map[interface {}]interface {}",
			))

			var ss = []scanStructTest{scanExpects.structSub, scanExpects.structSub}
			err = gconv.Scan(jsons, &ss)
			t.AssertNil(err)
			t.Assert(len(ss), 2)
			for k, _ := range ss {
				t.Assert("Mars", ss[k].Name)
				t.Assert("奥林帕斯山", ss[k].Place)
			}

			var ssp = []*scanStructTest{scanExpects.structSubPtr, scanExpects.structSubPtr}
			err = gconv.Scan(jsons, &ssp)
			t.AssertNil(err)
			t.Assert(len(ssp), 2)
			for k, _ := range ssp {
				t.Assert("Mars", ssp[k].Name)
				t.Assert("奥林帕斯山", ssp[k].Place)
			}
		}
	})

	// Test for paramKeyToAttrMap
	gtest.C(t, func(t *gtest.T) {
		scanValuesTest := scanValueMapsTest
		for _, test := range scanValuesTest {
			var (
				err          error
				scanExpects  = scanExpects
				mapParameter = map[string]string{"Name": "Place", "Place": "Name"}
			)

			// TODO: The following test cases should be working, but they are not.
			//err = gconv.Scan(test, &scanExpects.mapStrStr, mapParameter)
			//t.AssertNil(err)
			//t.Assert(test["Name"], scanExpects.mapStrStr["Place"])
			//t.Assert(test["Place"], scanExpects.mapStrStr["Name"])
			//
			//err = gconv.Scan(test, &scanExpects.mapStrAny, mapParameter)
			//t.AssertNil(err)
			//t.Assert(test["Name"], scanExpects.mapStrAny["Place"])
			//t.Assert(test["Place"], scanExpects.mapStrAny["Name"])
			//
			//err = gconv.Scan(test, &scanExpects.mapAnyAny, mapParameter)
			//t.AssertNil(err)
			//t.Assert(test["Name"], scanExpects.mapAnyAny["Place"])
			//t.Assert(test["Place"], scanExpects.mapAnyAny["Name"])

			err = gconv.Scan(test, &scanExpects.structSub, mapParameter)
			t.AssertNil(err)
			t.Assert(test["Name"], scanExpects.structSub.Place)
			t.Assert(test["Place"], scanExpects.structSub.Name)

			err = gconv.Scan(test, &scanExpects.structSubPtr, mapParameter)
			t.AssertNil(err)
			t.Assert(test["Name"], scanExpects.structSubPtr.Place)
			t.Assert(test["Place"], scanExpects.structSubPtr.Name)
		}
	})

	// Test for special types.
	gtest.C(t, func(t *gtest.T) {
		var (
			err error
			src = "Sun"
			dst = "日冕"
		)

		err = gconv.Scan(nil, &dst)
		t.AssertNil(err)
		t.Assert(dst, "日冕")

		err = gconv.Scan(src, nil)
		t.Assert(err, gerror.New("destination pointer should not be nil"))

		// Test for non-pointer.
		err = gconv.Scan(src, dst)
		t.Assert(err, gerror.New(
			"destination pointer should be type of pointer, but got type: string",
		))
	})
}

func TestScanEmptyStringToCustomType(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		type Status string
		type Req struct {
			Name     string
			Statuses []Status
			Types    []string
		}
		var (
			req  *Req
			data = g.Map{
				"Name":     "john",
				"Statuses": "",
				"Types":    "",
			}
		)
		err := gconv.Scan(data, &req)
		t.AssertNil(err)
		t.Assert(len(req.Statuses), 0)
		t.Assert(len(req.Types), 0)
	})
}
