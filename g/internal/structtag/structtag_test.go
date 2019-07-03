// Copyright 2019 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package structtag_test

import (
	"testing"

	"github.com/gogf/gf/g"

	"github.com/gogf/gf/g/internal/structtag"

	"github.com/gogf/gf/g/test/gtest"
)

func Test_Basic(t *testing.T) {
	gtest.Case(t, func() {
		type User struct {
			Id   int
			Name string `params:"name"`
			Pass string `my-tag1:"pass1" my-tag2:"pass2" params:"pass"`
		}
		var user User
		gtest.Assert(structtag.Map(user, []string{"params"}), g.Map{"name": "Name", "pass": "Pass"})
		gtest.Assert(structtag.Map(&user, []string{"params"}), g.Map{"name": "Name", "pass": "Pass"})

		gtest.Assert(structtag.Map(&user, []string{"params", "my-tag1"}), g.Map{"name": "Name", "pass": "Pass"})
		gtest.Assert(structtag.Map(&user, []string{"my-tag1", "params"}), g.Map{"name": "Name", "pass1": "Pass"})
		gtest.Assert(structtag.Map(&user, []string{"my-tag2", "params"}), g.Map{"name": "Name", "pass2": "Pass"})
	})

	gtest.Case(t, func() {
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
		gtest.Assert(structtag.Map(user, []string{"params"}), g.Map{"base": "Base"})
	})

	gtest.Case(t, func() {
		type Base struct {
			Pass1 string `params:"password1"`
			Pass2 string `params:"password2"`
		}
		type UserWithBase1 struct {
			Id   int
			Name string
			Base
		}
		type UserWithBase2 struct {
			Id   int
			Name string
			Pass Base
		}
		user1 := new(UserWithBase1)
		user2 := new(UserWithBase2)
		gtest.Assert(structtag.Map(user1, []string{"params"}), g.Map{"password1": "Pass1", "password2": "Pass2"})
		gtest.Assert(structtag.Map(user2, []string{"params"}), g.Map{"password1": "Pass1", "password2": "Pass2"})
	})
}
