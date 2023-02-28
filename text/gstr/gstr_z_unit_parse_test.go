// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// go test *.go -bench=".*"

package gstr_test

import (
	"net/url"
	"testing"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/text/gstr"
)

func Test_Parse(t *testing.T) {
	// cover test
	gtest.C(t, func(t *gtest.T) {
		// empty
		m, err := gstr.Parse("")
		t.AssertNil(err)
		t.Assert(m, nil)
		// invalid
		m, err = gstr.Parse("a&b")
		t.AssertNil(err)
		t.Assert(m, make(map[string]interface{}))
		// special key
		m, err = gstr.Parse(" =1& b=2&   c =3")
		t.AssertNil(err)
		t.Assert(m, map[string]interface{}{"b": "2", "c_": "3"})
		m, err = gstr.Parse("c[=3")
		t.AssertNil(err)
		t.Assert(m, map[string]interface{}{"c_": "3"})
		m, err = gstr.Parse("v[a][a]a=m")
		t.AssertNil(err)
		t.Assert(m, g.Map{
			"v": g.Map{
				"a": g.Map{
					"a": "m",
				},
			},
		})
		// v[][a]=m&v[][b]=b => map["v"]:[{"a":"m","b":"b"}]
		m, err = gstr.Parse("v[][a]=m&v[][b]=b")
		t.AssertNil(err)
		t.Assert(m, g.Map{
			"v": g.Slice{
				g.Map{
					"a": "m",
					"b": "b",
				},
			},
		})
		// v[][a]=m&v[][a]=b => map["v"]:[{"a":"m"},{"a":"b"}]
		m, err = gstr.Parse("v[][a]=m&v[][a]=b")
		t.AssertNil(err)
		t.Assert(m, g.Map{
			"v": g.Slice{
				g.Map{
					"a": "m",
				},
				g.Map{
					"a": "b",
				},
			},
		})
		// error
		m, err = gstr.Parse("v=111&v[]=m&v[]=a&v[]=b")
		t.Log(err)
		t.AssertNE(err, nil)
		m, err = gstr.Parse("v=111&v[a]=m&v[a]=a")
		t.Log(err)
		t.AssertNE(err, nil)
		_, err = gstr.Parse("%Q=%Q&b")
		t.Log(err)
		t.AssertNE(err, nil)
		_, err = gstr.Parse("a=%Q&b")
		t.Log(err)
		t.AssertNE(err, nil)
		_, err = gstr.Parse("v[a][a]=m&v[][a]=b")
		t.Log(err)
		t.AssertNE(err, nil)
	})

	// url
	gtest.C(t, func(t *gtest.T) {
		s := "goframe.org/index?name=john&score=100"
		u, err := url.Parse(s)
		t.AssertNil(err)
		m, err := gstr.Parse(u.RawQuery)
		t.AssertNil(err)
		t.Assert(m["name"], "john")
		t.Assert(m["score"], "100")

		// name overwrite
		m, err = gstr.Parse("a=1&a=2")
		t.AssertNil(err)
		t.Assert(m, g.Map{
			"a": 2,
		})
		// slice
		m, err = gstr.Parse("a[]=1&a[]=2")
		t.AssertNil(err)
		t.Assert(m, g.Map{
			"a": g.Slice{"1", "2"},
		})
		// map
		m, err = gstr.Parse("a=1&b=2&c=3")
		t.AssertNil(err)
		t.Assert(m, g.Map{
			"a": "1",
			"b": "2",
			"c": "3",
		})
		m, err = gstr.Parse("a=1&a=2&c=3")
		t.AssertNil(err)
		t.Assert(m, g.Map{
			"a": "2",
			"c": "3",
		})
		// map
		m, err = gstr.Parse("m[a]=1&m[b]=2&m[c]=3")
		t.AssertNil(err)
		t.Assert(m, g.Map{
			"m": g.Map{
				"a": "1",
				"b": "2",
				"c": "3",
			},
		})
		m, err = gstr.Parse("m[a]=1&m[a]=2&m[b]=3")
		t.AssertNil(err)
		t.Assert(m, g.Map{
			"m": g.Map{
				"a": "2",
				"b": "3",
			},
		})
		// map - slice
		m, err = gstr.Parse("m[a][]=1&m[a][]=2")
		t.AssertNil(err)
		t.Assert(m, g.Map{
			"m": g.Map{
				"a": g.Slice{"1", "2"},
			},
		})
		m, err = gstr.Parse("m[a][b][]=1&m[a][b][]=2")
		t.AssertNil(err)
		t.Assert(m, g.Map{
			"m": g.Map{
				"a": g.Map{
					"b": g.Slice{"1", "2"},
				},
			},
		})
		// map - complicated
		m, err = gstr.Parse("m[a1][b1][c1][d1]=1&m[a2][b2]=2&m[a3][b3][c3]=3")
		t.AssertNil(err)
		t.Assert(m, g.Map{
			"m": g.Map{
				"a1": g.Map{
					"b1": g.Map{
						"c1": g.Map{
							"d1": "1",
						},
					},
				},
				"a2": g.Map{
					"b2": "2",
				},
				"a3": g.Map{
					"b3": g.Map{
						"c3": "3",
					},
				},
			},
		})
	})

}
