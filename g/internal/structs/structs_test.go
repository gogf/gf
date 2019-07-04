// Copyright 2019 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package structs_test

import (
	"testing"

	"github.com/gogf/gf/g/internal/structs"

	"github.com/gogf/gf/g"

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
		gtest.Assert(structs.TagMapName(user, []string{"params"}, true), g.Map{"name": "Name", "pass": "Pass"})
		gtest.Assert(structs.TagMapName(&user, []string{"params"}, true), g.Map{"name": "Name", "pass": "Pass"})

		gtest.Assert(structs.TagMapName(&user, []string{"params", "my-tag1"}, true), g.Map{"name": "Name", "pass": "Pass"})
		gtest.Assert(structs.TagMapName(&user, []string{"my-tag1", "params"}, true), g.Map{"name": "Name", "pass1": "Pass"})
		gtest.Assert(structs.TagMapName(&user, []string{"my-tag2", "params"}, true), g.Map{"name": "Name", "pass2": "Pass"})
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
		gtest.Assert(structs.TagMapName(user, []string{"params"}, true), g.Map{
			"base":      "Base",
			"password1": "Pass1",
			"password2": "Pass2",
		})
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
		gtest.Assert(structs.TagMapName(user1, []string{"params"}, true), g.Map{"password1": "Pass1", "password2": "Pass2"})
		gtest.Assert(structs.TagMapName(user2, []string{"params"}, true), g.Map{"password1": "Pass1", "password2": "Pass2"})
	})
}
