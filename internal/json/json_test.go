// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package json_test

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/internal/json"
	"github.com/gogf/gf/test/gtest"
	"testing"
)

type User struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

var (
	user = &User{
		Id:   1265476890672672808,
		Name: "john",
	}
	userBytes, _ = json.Marshal(user)
)

func TestMarshal(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		b, err := json.Marshal(user)
		t.AssertNil(err)
		t.Assert(b, `{"id":1265476890672672808,"name":"john"}`)
	})
}

func TestMarshalIndent(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		b, err := json.MarshalIndent(user, "#", "@")
		t.AssertNil(err)
		t.Assert(b, `{
#@"id": 1265476890672672808,
#@"name": "john"
#}`)
	})
}

func TestUnmarshal(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var m map[string]interface{}
		b, _ := json.Marshal(g.Map{
			"user": user,
		})
		err := json.Unmarshal(b, &m)
		t.AssertNil(err)
		// precision lost for big int.
		t.Assert(m["user"], g.Map{
			"id":   1265476890672672800,
			"name": user.Name,
		})
	})
}

func TestUnmarshalUseNumber(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var m map[string]interface{}
		b, _ := json.Marshal(g.Map{
			"user": user,
		})
		err := json.UnmarshalUseNumber(b, &m)
		t.AssertNil(err)
		t.Assert(m["user"], g.Map{
			"id":   user.Id,
			"name": user.Name,
		})
	})
}

func TestValid(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		m := g.Map{
			`{}`:                     true,
			`{"id":1,"name":"john"}`: true,
			`1`:                      true,
			`"john"`:                 true,
			`"`:                      false,
			`<xml></xml>`:            false,
		}
		for k, v := range m {
			t.Assert(json.Valid([]byte(k)), v)
		}
	})
}
