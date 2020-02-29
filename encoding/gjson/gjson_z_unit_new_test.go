// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gjson_test

import (
	"testing"

	"github.com/gogf/gf/encoding/gjson"
	"github.com/gogf/gf/test/gtest"
)

func Test_Load_NewWithTag(t *testing.T) {
	type User struct {
		Age  int    `xml:"age-xml"  json:"age-json"`
		Name string `xml:"name-xml" json:"name-json"`
		Addr string `xml:"addr-xml" json:"addr-json"`
	}
	data := User{
		Age:  18,
		Name: "john",
		Addr: "chengdu",
	}
	// JSON
	gtest.Case(t, func() {
		j := gjson.New(data)
		gtest.AssertNE(j, nil)
		gtest.Assert(j.Get("age-xml"), nil)
		gtest.Assert(j.Get("age-json"), data.Age)
		gtest.Assert(j.Get("name-xml"), nil)
		gtest.Assert(j.Get("name-json"), data.Name)
		gtest.Assert(j.Get("addr-xml"), nil)
		gtest.Assert(j.Get("addr-json"), data.Addr)
	})
	// XML
	gtest.Case(t, func() {
		j := gjson.NewWithTag(data, "xml")
		gtest.AssertNE(j, nil)
		gtest.Assert(j.Get("age-xml"), data.Age)
		gtest.Assert(j.Get("age-json"), nil)
		gtest.Assert(j.Get("name-xml"), data.Name)
		gtest.Assert(j.Get("name-json"), nil)
		gtest.Assert(j.Get("addr-xml"), data.Addr)
		gtest.Assert(j.Get("addr-json"), nil)
	})
}
