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
	gtest.C(t, func(t *gtest.T) {
		j := gjson.New(data)
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
		j := gjson.NewWithTag(data, "xml")
		t.AssertNE(j, nil)
		t.Assert(j.Get("age-xml"), data.Age)
		t.Assert(j.Get("age-json"), nil)
		t.Assert(j.Get("name-xml"), data.Name)
		t.Assert(j.Get("name-json"), nil)
		t.Assert(j.Get("addr-xml"), data.Addr)
		t.Assert(j.Get("addr-json"), nil)
	})
}

func Test_Load_New_CustomStruct(t *testing.T) {
	type Base struct {
		Id int
	}
	type User struct {
		Base
		Name string
	}
	user := new(User)
	user.Id = 1
	user.Name = "john"

	gtest.C(t, func(t *gtest.T) {
		j := gjson.New(user)
		t.AssertNE(j, nil)

		s, err := j.ToJsonString()
		t.Assert(err, nil)
		t.Assert(s == `{"Id":1,"Name":"john"}` || s == `{"Name":"john","Id":1}`, true)
	})
}
