// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gconv_test

import (
	"github.com/gogf/gf/container/gvar"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/test/gtest"
	"github.com/gogf/gf/util/gconv"
	"testing"
)

func Test_Scan_StructStructs(t *testing.T) {
	type User struct {
		Uid   int
		Name  string
		Pass1 string `gconv:"password1"`
		Pass2 string `gconv:"password2"`
	}
	gtest.C(t, func(t *gtest.T) {
		var (
			user   = new(User)
			params = g.Map{
				"uid":   1,
				"name":  "john",
				"PASS1": "123",
				"PASS2": "456",
			}
		)
		err := gconv.Scan(params, user)
		t.Assert(err, nil)
		t.Assert(user, &User{
			Uid:   1,
			Name:  "john",
			Pass1: "123",
			Pass2: "456",
		})
	})
	gtest.C(t, func(t *gtest.T) {
		var (
			users  []User
			params = g.Slice{
				g.Map{
					"uid":   1,
					"name":  "john1",
					"PASS1": "111",
					"PASS2": "222",
				},
				g.Map{
					"uid":   2,
					"name":  "john2",
					"PASS1": "333",
					"PASS2": "444",
				},
			}
		)
		err := gconv.Scan(params, &users)
		t.AssertNil(err)
		t.Assert(users, g.Slice{
			&User{
				Uid:   1,
				Name:  "john1",
				Pass1: "111",
				Pass2: "222",
			},
			&User{
				Uid:   2,
				Name:  "john2",
				Pass1: "333",
				Pass2: "444",
			},
		})
	})
}

func Test_Scan_StructStr(t *testing.T) {
	type User struct {
		Uid   int
		Name  string
		Pass1 string `gconv:"password1"`
		Pass2 string `gconv:"password2"`
	}
	gtest.C(t, func(t *gtest.T) {
		var (
			user   = new(User)
			params = `{"uid":1,"name":"john", "pass1":"123","pass2":"456"}`
		)
		err := gconv.Scan(params, user)
		t.Assert(err, nil)
		t.Assert(user, &User{
			Uid:   1,
			Name:  "john",
			Pass1: "123",
			Pass2: "456",
		})
	})
	gtest.C(t, func(t *gtest.T) {
		var (
			users  []User
			params = `[
{"uid":1,"name":"john1", "pass1":"111","pass2":"222"},
{"uid":2,"name":"john2", "pass1":"333","pass2":"444"}
]`
		)
		err := gconv.Scan(params, &users)
		t.Assert(err, nil)
		t.Assert(users, g.Slice{
			&User{
				Uid:   1,
				Name:  "john1",
				Pass1: "111",
				Pass2: "222",
			},
			&User{
				Uid:   2,
				Name:  "john2",
				Pass1: "333",
				Pass2: "444",
			},
		})
	})
}

func Test_Scan_Map(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var m map[string]string
		data := g.Map{
			"k1": "v1",
			"k2": "v2",
		}
		err := gconv.Scan(data, &m)
		t.AssertNil(err)
		t.Assert(data, m)
	})
	gtest.C(t, func(t *gtest.T) {
		var m map[int]int
		data := g.Map{
			"1": "11",
			"2": "22",
		}
		err := gconv.Scan(data, &m)
		t.AssertNil(err)
		t.Assert(data, m)
	})
	// json string parameter.
	gtest.C(t, func(t *gtest.T) {
		var m map[string]string
		data := `{"k1":"v1","k2":"v2"}`
		err := gconv.Scan(data, &m)
		t.AssertNil(err)
		t.Assert(m, g.Map{
			"k1": "v1",
			"k2": "v2",
		})
	})
}

func Test_Scan_Maps(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var maps []map[string]string
		data := g.Slice{
			g.Map{
				"k1": "v1",
				"k2": "v2",
			},
			g.Map{
				"k3": "v3",
				"k4": "v4",
			},
		}
		err := gconv.Scan(data, &maps)
		t.AssertNil(err)
		t.Assert(data, maps)
	})
	// json string parameter.
	gtest.C(t, func(t *gtest.T) {
		var maps []map[string]string
		data := `[{"k1":"v1","k2":"v2"},{"k3":"v3","k4":"v4"}]`
		err := gconv.Scan(data, &maps)
		t.AssertNil(err)
		t.Assert(maps, g.Slice{
			g.Map{
				"k1": "v1",
				"k2": "v2",
			},
			g.Map{
				"k3": "v3",
				"k4": "v4",
			},
		})
	})
}

func Test_Scan_JsonAttributes(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		type Sku struct {
			GiftId      int64  `json:"gift_id"`
			Name        string `json:"name"`
			ScorePrice  int    `json:"score_price"`
			MarketPrice int    `json:"market_price"`
			CostPrice   int    `json:"cost_price"`
			Stock       int    `json:"stock"`
		}
		v := gvar.New(`
[
{"name": "red", "stock": 10, "gift_id": 1, "cost_price": 80, "score_price": 188, "market_price": 188}, 
{"name": "blue", "stock": 100, "gift_id": 2, "cost_price": 81, "score_price": 200, "market_price": 288}
]`)
		type Product struct {
			Skus []Sku
		}
		var p *Product
		err := gconv.Scan(g.Map{
			"Skus": v,
		}, &p)
		t.AssertNil(err)
		t.Assert(len(p.Skus), 2)

		t.Assert(p.Skus[0].Name, "red")
		t.Assert(p.Skus[0].Stock, 10)
		t.Assert(p.Skus[0].GiftId, 1)
		t.Assert(p.Skus[0].CostPrice, 80)
		t.Assert(p.Skus[0].ScorePrice, 188)
		t.Assert(p.Skus[0].MarketPrice, 188)

		t.Assert(p.Skus[1].Name, "blue")
		t.Assert(p.Skus[1].Stock, 100)
		t.Assert(p.Skus[1].GiftId, 2)
		t.Assert(p.Skus[1].CostPrice, 81)
		t.Assert(p.Skus[1].ScorePrice, 200)
		t.Assert(p.Skus[1].MarketPrice, 288)
	})
}

func Test_Scan_JsonAttributes_StringArray(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		type S struct {
			Array []string
		}
		var s *S
		err := gconv.Scan(g.Map{
			"Array": `["a", "b"]`,
		}, &s)
		t.AssertNil(err)
		t.Assert(len(s.Array), 2)
		t.Assert(s.Array[0], "a")
		t.Assert(s.Array[1], "b")
	})

	gtest.C(t, func(t *gtest.T) {
		type S struct {
			Array []string
		}
		var s *S
		err := gconv.Scan(g.Map{
			"Array": `[]`,
		}, &s)
		t.AssertNil(err)
		t.Assert(len(s.Array), 0)
	})

	gtest.C(t, func(t *gtest.T) {
		type S struct {
			Array []int64
		}
		var s *S
		err := gconv.Scan(g.Map{
			"Array": `[]`,
		}, &s)
		t.AssertNil(err)
		t.Assert(len(s.Array), 0)
	})
}
