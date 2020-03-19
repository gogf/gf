// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gparser_test

import (
	"github.com/gogf/gf/encoding/gparser"
	"testing"

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
	gtest.C(t, func(t *gtest.T) {
		j := gparser.New(data)
		t.AssertNE(j, nil)
		t.Assert(j.Get("age-xml"), nil)
		t.Assert(j.Get("age-json"), data.Age)
		t.Assert(j.Get("name-xml"), nil)
		t.Assert(j.Get("name-json"), data.Name)
		t.Assert(j.Get("addr-xml"), nil)
		t.Assert(j.Get("addr-json"), data.Addr)
	})
	// XML
	gtest.C(t, func(t *gtest.T) {
		j := gparser.NewWithTag(data, "xml")
		t.AssertNE(j, nil)
		t.Assert(j.Get("age-xml"), data.Age)
		t.Assert(j.Get("age-json"), nil)
		t.Assert(j.Get("name-xml"), data.Name)
		t.Assert(j.Get("name-json"), nil)
		t.Assert(j.Get("addr-xml"), data.Addr)
		t.Assert(j.Get("addr-json"), nil)
	})
}
