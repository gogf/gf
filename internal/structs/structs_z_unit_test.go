// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package structs_test

import (
	"testing"

	"github.com/gogf/gf/internal/structs"

	"github.com/gogf/gf/frame/g"

	"github.com/gogf/gf/test/gtest"
)

func Test_Basic(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		type User struct {
			Id   int
			Name string `params:"name"`
			Pass string `my-tag1:"pass1" my-tag2:"pass2" params:"pass"`
		}
		var user User
		m, _ := structs.TagMapName(user, []string{"params"})
		t.Assert(m, g.Map{"name": "Name", "pass": "Pass"})
		m, _ = structs.TagMapName(&user, []string{"params"})
		t.Assert(m, g.Map{"name": "Name", "pass": "Pass"})

		m, _ = structs.TagMapName(&user, []string{"params", "my-tag1"})
		t.Assert(m, g.Map{"name": "Name", "pass": "Pass"})
		m, _ = structs.TagMapName(&user, []string{"my-tag1", "params"})
		t.Assert(m, g.Map{"name": "Name", "pass1": "Pass"})
		m, _ = structs.TagMapName(&user, []string{"my-tag2", "params"})
		t.Assert(m, g.Map{"name": "Name", "pass2": "Pass"})
	})

	gtest.C(t, func(t *gtest.T) {
		type Base struct {
			Pass1 string `params:"password1"`
			Pass2 string `params:"password2"`
		}
		type UserWithBase struct {
			Id   int
			Name string
			Base `params:"base"`
		}
		user := new(UserWithBase)
		m, _ := structs.TagMapName(user, []string{"params"})
		t.Assert(m, g.Map{
			"base":      "Base",
			"password1": "Pass1",
			"password2": "Pass2",
		})
	})

	gtest.C(t, func(t *gtest.T) {
		type Base struct {
			Pass1 string `params:"password1"`
			Pass2 string `params:"password2"`
		}
		type UserWithEmbeddedAttribute struct {
			Id   int
			Name string
			Base
		}
		type UserWithoutEmbeddedAttribute struct {
			Id   int
			Name string
			Pass Base
		}
		user1 := new(UserWithEmbeddedAttribute)
		user2 := new(UserWithoutEmbeddedAttribute)
		m, _ := structs.TagMapName(user1, []string{"params"})
		t.Assert(m, g.Map{"password1": "Pass1", "password2": "Pass2"})
		m, _ = structs.TagMapName(user2, []string{"params"})
		t.Assert(m, g.Map{})
	})
}

func Test_StructOfNilPointer(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		type User struct {
			Id   int
			Name string `params:"name"`
			Pass string `my-tag1:"pass1" my-tag2:"pass2" params:"pass"`
		}
		var user *User
		m, _ := structs.TagMapName(user, []string{"params"})
		t.Assert(m, g.Map{"name": "Name", "pass": "Pass"})
		m, _ = structs.TagMapName(&user, []string{"params"})
		t.Assert(m, g.Map{"name": "Name", "pass": "Pass"})

		m, _ = structs.TagMapName(&user, []string{"params", "my-tag1"})
		t.Assert(m, g.Map{"name": "Name", "pass": "Pass"})
		m, _ = structs.TagMapName(&user, []string{"my-tag1", "params"})
		t.Assert(m, g.Map{"name": "Name", "pass1": "Pass"})
		m, _ = structs.TagMapName(&user, []string{"my-tag2", "params"})
		t.Assert(m, g.Map{"name": "Name", "pass2": "Pass"})
	})
}

func Test_FieldMap(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		type User struct {
			Id   int
			Name string `params:"name"`
			Pass string `my-tag1:"pass1" my-tag2:"pass2" params:"pass"`
		}
		var user *User
		m, _ := structs.FieldMap(structs.FieldMapInput{
			Pointer:          user,
			PriorityTagArray: []string{"params"},
			RecursiveOption:  structs.RecursiveOptionEmbedded,
		})
		t.Assert(len(m), 3)
		_, ok := m["Id"]
		t.Assert(ok, true)
		_, ok = m["Name"]
		t.Assert(ok, false)
		_, ok = m["name"]
		t.Assert(ok, true)
		_, ok = m["Pass"]
		t.Assert(ok, false)
		_, ok = m["pass"]
		t.Assert(ok, true)
	})
	gtest.C(t, func(t *gtest.T) {
		type User struct {
			Id   int
			Name string `params:"name"`
			Pass string `my-tag1:"pass1" my-tag2:"pass2" params:"pass"`
		}
		var user *User
		m, _ := structs.FieldMap(structs.FieldMapInput{
			Pointer:          user,
			PriorityTagArray: nil,
			RecursiveOption:  structs.RecursiveOptionEmbedded,
		})
		t.Assert(len(m), 3)
		_, ok := m["Id"]
		t.Assert(ok, true)
		_, ok = m["Name"]
		t.Assert(ok, true)
		_, ok = m["name"]
		t.Assert(ok, false)
		_, ok = m["Pass"]
		t.Assert(ok, true)
		_, ok = m["pass"]
		t.Assert(ok, false)
	})
}

func Test_StructType(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		type B struct {
			Name string
		}
		type A struct {
			B
		}
		r, err := structs.StructType(new(A))
		t.AssertNil(err)
		t.Assert(r.Signature(), `github.com/gogf/gf/internal/structs_test/structs_test.A`)
	})
	gtest.C(t, func(t *gtest.T) {
		type B struct {
			Name string
		}
		type A struct {
			B
		}
		r, err := structs.StructType(new(A).B)
		t.AssertNil(err)
		t.Assert(r.Signature(), `github.com/gogf/gf/internal/structs_test/structs_test.B`)
	})
	gtest.C(t, func(t *gtest.T) {
		type B struct {
			Name string
		}
		type A struct {
			*B
		}
		r, err := structs.StructType(new(A).B)
		t.AssertNil(err)
		t.Assert(r.String(), `structs_test.B`)
	})
	// Error.
	gtest.C(t, func(t *gtest.T) {
		type B struct {
			Name string
		}
		type A struct {
			*B
			Id int
		}
		_, err := structs.StructType(new(A).Id)
		t.AssertNE(err, nil)
	})
}

func Test_StructTypeBySlice(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		type B struct {
			Name string
		}
		type A struct {
			Array []*B
		}
		r, err := structs.StructType(new(A).Array)
		t.AssertNil(err)
		t.Assert(r.Signature(), `github.com/gogf/gf/internal/structs_test/structs_test.B`)
	})
	gtest.C(t, func(t *gtest.T) {
		type B struct {
			Name string
		}
		type A struct {
			Array []B
		}
		r, err := structs.StructType(new(A).Array)
		t.AssertNil(err)
		t.Assert(r.Signature(), `github.com/gogf/gf/internal/structs_test/structs_test.B`)
	})
	gtest.C(t, func(t *gtest.T) {
		type B struct {
			Name string
		}
		type A struct {
			Array *[]B
		}
		r, err := structs.StructType(new(A).Array)
		t.AssertNil(err)
		t.Assert(r.Signature(), `github.com/gogf/gf/internal/structs_test/structs_test.B`)
	})
}

func TestType_FieldKeys(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		type B struct {
			Id   int
			Name string
		}
		type A struct {
			Array []*B
		}
		r, err := structs.StructType(new(A).Array)
		t.AssertNil(err)
		t.Assert(r.FieldKeys(), g.Slice{"Id", "Name"})
	})
}
