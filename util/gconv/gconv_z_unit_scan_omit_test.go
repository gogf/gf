// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gconv_test

import (
	"testing"

	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/gconv"
)

type User struct {
	Name  string
	Age   int
	Email string
}

type User2 struct {
	Name  *string
	Age   int
	Email string
}

type Person struct {
	Name  string
	Age   int
	Email string
}

func TestScan_OmitEmpty(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		user := User{Name: "", Age: 20, Email: ""}
		person := Person{Name: "zhangsan", Age: 0, Email: "old@example.com"}

		err := gconv.ScanWithOptions(user, &person, gconv.ScanOption{
			OmitEmpty: true,
		})
		t.AssertNil(err)
		t.Assert(person.Name, "zhangsan")
		t.Assert(person.Age, 20)
		t.Assert(person.Email, "old@example.com")
	})
}

func TestScan_OmitNil(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		data := map[string]any{
			"Name":  nil,
			"Age":   30,
			"Email": nil,
		}
		person := Person{Name: "lisi", Age: 0, Email: "old@example.com"}

		err := gconv.ScanWithOptions(data, &person, gconv.ScanOption{
			OmitNil: true,
		})
		t.AssertNil(err)
		t.Assert(person.Name, "lisi")
		t.Assert(person.Age, 30)
		t.Assert(person.Email, "old@example.com")
	})
}

func TestScan_OmitEmptyAndOmitNil(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		data := map[string]any{
			"Name":  "",
			"Age":   25,
			"Email": nil,
		}
		person := Person{Name: "wangwu", Age: 0, Email: "old2@example.com"}

		err := gconv.ScanWithOptions(data, &person, gconv.ScanOption{
			OmitEmpty: true,
			OmitNil:   true,
		})
		t.AssertNil(err)
		t.Assert(person.Name, "wangwu")
		t.Assert(person.Age, 25)
		t.Assert(person.Email, "old2@example.com")
	})
}

func TestScan_NoOmitOptions(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		user := User{Name: "", Age: 20, Email: ""}
		person := Person{Name: "zhangsan", Age: 30, Email: "old@example.com"}

		err := gconv.ScanWithOptions(user, &person, gconv.ScanOption{
			OmitEmpty: false,
			OmitNil:   false,
		})
		t.AssertNil(err)
		t.Assert(person.Name, "")
		t.Assert(person.Age, 20)
		t.Assert(person.Email, "")
	})
}

func TestScan_OriginalBehavior(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		user := User{Name: "newname", Age: 25, Email: "new@example.com"}
		person := Person{Name: "", Age: 0, Email: ""}

		err := gconv.Scan(user, &person)
		t.AssertNil(err)
		t.Assert(person.Name, "newname")
		t.Assert(person.Age, 25)
		t.Assert(person.Email, "new@example.com")
	})
}

func TestScan_StructOmitEmptyAndOmitNilOptions(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		user2 := User2{Name: nil, Age: 25, Email: ""}
		person := Person{Name: "wangwu", Age: 0, Email: "old2@example.com"}

		err := gconv.ScanWithOptions(user2, &person, gconv.ScanOption{
			OmitEmpty: true,
			OmitNil:   true,
		})
		t.AssertNil(err)
		t.Assert(person.Name, "wangwu")
		t.Assert(person.Age, 25)
		t.Assert(person.Email, "old2@example.com")
	})
}
