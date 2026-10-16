// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// This file tests how struct converting matches the parameter keys that differ from the tags of
// the fields in letter case, like the column names of a database folding unquoted identifiers to
// upper case.

package gconv_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/gconv"
)

// Test_Struct_CaseInsensitiveTag asserts that, with StructOption.CaseInsensitiveTag, a key matches
// the field whose tag it equals ignoring letter case, rather than any field whose name it equals
// once symbols are removed. Each round converts to a new struct type, one more field making it
// differ, so that no match cached by a previous round decides the result.
func Test_Struct_CaseInsensitiveTag(t *testing.T) {
	var (
		converter = gconv.NewConverter()
		option    = gconv.StructOption{PriorityTag: "orm", CaseInsensitiveTag: true}
	)
	gtest.C(t, func(t *gtest.T) {
		var intType = reflect.TypeOf(0)
		for i := 0; i < 100; i++ {
			structType := reflect.StructOf([]reflect.StructField{
				{Name: "Id", Type: intType, Tag: `orm:"id"`},
				{Name: "Userid", Type: intType, Tag: `orm:"user_id"`},
				{Name: "UserId", Type: intType, Tag: `orm:"UserId"`},
				{Name: fmt.Sprintf("Round%d", i), Type: intType, Tag: `orm:"-"`},
			})
			pointer := reflect.New(structType)
			err := converter.Struct(map[string]any{
				"ID":      1,
				"USER_ID": 2,
				"USERID":  3,
			}, pointer.Interface(), option)
			t.AssertNil(err)
			t.Assert(pointer.Elem().FieldByName("Id").Int(), 1)
			t.Assert(pointer.Elem().FieldByName("Userid").Int(), 2)
			t.Assert(pointer.Elem().FieldByName("UserId").Int(), 3)
		}
	})
	gtest.C(t, func(t *gtest.T) {
		type Instance struct {
			ID   int    `orm:"f_id"`
			Name string `orm:"name"`
		}
		for i := 0; i < 3; i++ {
			var instance *Instance
			err := converter.Struct(map[string]any{
				"F_ID": 1,
				"NAME": "john",
			}, &instance, option)
			t.AssertNil(err)
			t.Assert(instance.ID, 1)
			t.Assert(instance.Name, "john")
		}
	})
	// The key equal to the tag is preferred to the key equal to it in case-insensitive.
	gtest.C(t, func(t *gtest.T) {
		type User struct {
			Id   int    `orm:"id"`
			Name string `orm:"name"`
		}
		var user *User
		err := converter.Struct(map[string]any{
			"ID":   2,
			"id":   1,
			"name": "john",
		}, &user, option)
		t.AssertNil(err)
		t.Assert(user.Id, 1)
		t.Assert(user.Name, "john")
	})
}

// Test_Struct_CaseSensitiveTagByDefault asserts that, without StructOption.CaseInsensitiveTag,
// a key differing from the tag of a field in letter case does not match the field by the tag.
func Test_Struct_CaseSensitiveTagByDefault(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		type Instance struct {
			ID int `orm:"f_id"`
		}
		var instance *Instance
		err := gconv.StructTag(map[string]any{"F_ID": 1}, &instance, "orm")
		t.AssertNil(err)
		t.Assert(instance.ID, 0)
	})
}

// Test_Struct_FuzzyKeyNotReused asserts that the key a field matched in a previous converting is
// not used again for the field when another field of the current converting has used it.
func Test_Struct_FuzzyKeyNotReused(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		type Article struct {
			Title string
			Name  string
		}
		var article *Article
		err := gconv.Struct(map[string]any{"NAME": "a"}, &article)
		t.AssertNil(err)
		t.Assert(article.Title, "")
		t.Assert(article.Name, "a")

		article = nil
		err = gconv.Struct(map[string]any{"NAME": "b", "OTHER": "c"}, &article, map[string]string{"NAME": "Title"})
		t.AssertNil(err)
		t.Assert(article.Title, "b")
		t.Assert(article.Name, "")
	})
}
