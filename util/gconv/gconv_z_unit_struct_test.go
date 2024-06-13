// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gconv_test

import (
	"testing"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/gconv"
)

type structExpect struct {
	PlanetName   string
	Planet_Place string
	planetTime   string
}

type structTagGconvExpect struct {
	PlanetNameGconv  string `gconv:"PlanetName"`
	PlanetPlaceGconv string `gconv:"-"`
}
type structTagParamExpect struct {
	PlanetNameParam  string `param:"PlanetName"`
	PlanetPlaceParam string `param:"-"`
}
type structTagCExpect struct {
	PlanetNameC  string `c:"PlanetName"`
	PlanetPlaceC string `c:"-"`
}
type structTagPExpect struct {
	PlanetNameP  string `p:"PlanetName"`
	PlanetPlaceP string `p:"-"`
}
type structTagJsonExpect struct {
	PlanetNameJson  string `json:"PlanetName"`
	PlanetPlaceJson string `json:"-"`
}

var structValueTests = []map[string]string{
	{
		"planetname":  "Earth",
		"planetplace": "亚马逊雨林",
		"planettime":  "2021-01-01",
	},
	{
		"planetName":  "Earth",
		"planetPlace": "亚马逊雨林",
		"planetTime":  "2021-01-01",
	},
	{
		"planet-name":  "Earth",
		"planet-place": "亚马逊雨林",
		"planet-time":  "2021-01-01",
	},
	{
		"planet_name":  "Earth",
		"planet_place": "亚马逊雨林",
		"planet_time":  "2021-01-01",
	},
	{
		"planet name":  "Earth",
		"planet place": "亚马逊雨林",
		"planet time":  "2021-01-01",
	},
	{
		"PLANETNAME":  "Earth",
		"PLANETPLACE": "亚马逊雨林",
		"PLANETTIME":  "2021-01-01",
	},
	{
		"PLANETnAME":  "Earth",
		"PLANETpLACE": "亚马逊雨林",
		"PLANETtIME":  "2021-01-01",
	},
	{
		"PLANET-NAME":  "Earth",
		"PLANET-PLACE": "亚马逊雨林",
		"PLANET-TIME":  "2021-01-01",
	},
	{
		"PLANET_NAME":  "Earth",
		"PLANET_PLACE": "亚马逊雨林",
		"PLANET_TIME":  "2021-01-01",
	},
	{
		"PLANET NAME":  "Earth",
		"PLANET PLACE": "亚马逊雨林",
		"PLANET TIME":  "2021-01-01",
	},
	{
		"PlanetName":  "Earth",
		"PlanetPlace": "亚马逊雨林",
		"PlanetTime":  "2021-01-01",
	},
	{
		"Planet-Name":  "Earth",
		"Planet-Place": "亚马逊雨林",
		"Planet-Time":  "2021-01-01",
	},
	{
		"Planet_Name":  "Earth",
		"Planet_Place": "亚马逊雨林",
		"Planet_Time":  "2021-01-01",
	},
	{
		"Planet Name":  "Earth",
		"Planet Place": "亚马逊雨林",
		"Planet Time":  "2021-01-01",
	},
}

func TestStruct(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		for _, test := range structValueTests {
			var (
				err    error
				expect = new(structExpect)
			)
			err = gconv.Struct(test, expect)
			t.AssertNil(err)
			t.Assert(expect.PlanetName, "Earth")
			t.Assert(expect.Planet_Place, "亚马逊雨林")
			t.Assert(expect.planetTime, "")

			tagTestValue, ok := test["PlanetName"]
			if !ok {
				continue
			}
			var (
				expectTagGconv = new(structTagGconvExpect)
				expectTagParam = new(structTagParamExpect)
				expectTagC     = new(structTagCExpect)
				expectTagP     = new(structTagPExpect)
				expectTagJson  = new(structTagJsonExpect)
			)
			err = gconv.Struct(test, expectTagGconv)
			t.AssertNil(err)
			t.Assert(expectTagGconv.PlanetNameGconv, tagTestValue)
			t.Assert(expectTagGconv.PlanetPlaceGconv, "")

			err = gconv.Struct(test, expectTagParam)
			t.AssertNil(err)
			t.Assert(expectTagParam.PlanetNameParam, tagTestValue)
			t.Assert(expectTagParam.PlanetPlaceParam, "")

			err = gconv.Struct(test, expectTagC)
			t.AssertNil(err)
			t.Assert(expectTagC.PlanetNameC, tagTestValue)
			t.Assert(expectTagC.PlanetPlaceC, "")

			err = gconv.Struct(test, expectTagP)
			t.AssertNil(err)
			t.Assert(expectTagP.PlanetNameP, tagTestValue)
			t.Assert(expectTagP.PlanetPlaceP, "")

			err = gconv.Struct(test, expectTagJson)
			t.AssertNil(err)
			t.Assert(expectTagJson.PlanetNameJson, tagTestValue)
			t.Assert(expectTagJson.PlanetPlaceJson, "")
		}
	})

	// Test for nil.
	gtest.C(t, func(t *gtest.T) {
		var (
			err    error
			expect = new(structExpect)
		)

		err = gconv.Struct(nil, nil)
		t.AssertNil(err)
		t.Assert(expect.PlanetName, "")
		t.Assert(expect.Planet_Place, "")
		t.Assert(expect.planetTime, "")
	})
}

func TestStructErr(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		type Score struct {
			Name   string
			Result int
		}
		type User struct {
			Score Score
		}

		user := new(User)
		scores := map[string]interface{}{
			"Score": 1,
		}
		err := gconv.Struct(scores, user)
		t.AssertNE(err, nil)
	})

	gtest.C(t, func(t *gtest.T) {
		type CustomString string
		type CustomStruct struct {
			S string
		}
		var (
			a CustomString = "abc"
			b *CustomStruct
		)
		err := gconv.Scan(a, &b)
		t.AssertNE(err, nil)
		t.Assert(b, nil)
	})

	gtest.C(t, func(t *gtest.T) {
		var i *int = nil
		err := gconv.Struct(map[string]string{}, i)
		t.AssertNE(err, nil)
	})
}

// Test for Struct containing time.Time attribute.
func TestStructWithTime(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		type S struct {
			T *gtime.Time
		}
		var (
			err error
			now = time.Now()
			s   = new(S)
		)
		err = gconv.Struct(g.Map{
			"t": &now,
		}, s)
		t.AssertNil(err)
		t.Assert(s.T.UTC().Time.String(), now.UTC().String())
	})
}

func TestStructs(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		for _, test := range structValueTests {
			var (
				err     error
				tests   = []map[string]string{test, test}
				expects []*structExpect
			)
			err = gconv.SliceStruct(tests, &expects)
			t.AssertNil(err)
			t.Assert(len(expects), 2)
			for _, expect := range expects {
				t.Assert(expect.PlanetName, "Earth")
				t.Assert(expect.Planet_Place, "亚马逊雨林")
				t.Assert(expect.planetTime, "")
			}

			tagTestValue, ok := test["PlanetName"]
			if !ok {
				continue
			}
			var (
				expectTagGconvs = []*structTagGconvExpect{}
				expectTagParams = []*structTagParamExpect{}
				expectTagCs     = []*structTagCExpect{}
				expectTagPs     = []*structTagPExpect{}
				expectTagJsons  = []*structTagJsonExpect{}
			)

			err = gconv.SliceStruct(tests, &expectTagGconvs)
			t.AssertNil(err)
			t.Assert(len(expectTagGconvs), 2)
			for _, expect := range expectTagGconvs {
				t.Assert(expect.PlanetNameGconv, tagTestValue)
				t.Assert(expect.PlanetPlaceGconv, "")
			}

			err = gconv.SliceStruct(tests, &expectTagParams)
			t.AssertNil(err)
			t.Assert(len(expectTagParams), 2)
			for _, expect := range expectTagParams {
				t.Assert(expect.PlanetNameParam, tagTestValue)
				t.Assert(expect.PlanetPlaceParam, "")
			}

			err = gconv.SliceStruct(tests, &expectTagCs)
			t.AssertNil(err)
			t.Assert(len(expectTagCs), 2)
			for _, expect := range expectTagCs {
				t.Assert(expect.PlanetNameC, tagTestValue)
				t.Assert(expect.PlanetPlaceC, "")
			}

			err = gconv.SliceStruct(tests, &expectTagPs)
			t.AssertNil(err)
			t.Assert(len(expectTagPs), 2)
			for _, expect := range expectTagPs {
				t.Assert(expect.PlanetNameP, tagTestValue)
				t.Assert(expect.PlanetPlaceP, "")
			}

			err = gconv.SliceStruct(tests, &expectTagJsons)
			t.AssertNil(err)
			t.Assert(len(expectTagJsons), 2)
			for _, expect := range expectTagJsons {
				t.Assert(expect.PlanetNameJson, tagTestValue)
				t.Assert(expect.PlanetPlaceJson, "")
			}
		}
	})
}
