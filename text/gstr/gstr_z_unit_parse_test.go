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
	// url
	gtest.C(t, func(t *gtest.T) {
		s := "goframe.org/index?name=john&score=100"
		u, err := url.Parse(s)
		t.Assert(err, nil)
		m, err := gstr.Parse(u.RawQuery)
		t.Assert(err, nil)
		t.Assert(m["name"], "john")
		t.Assert(m["score"], "100")

		// name overwrite
		m, err = gstr.Parse("a=1&a=2")
		t.Assert(err, nil)
		t.Assert(m, g.Map{
			"a": 2,
		})
		// slice
		m, err = gstr.Parse("a[]=1&a[]=2")
		t.Assert(err, nil)
		t.Assert(m, g.Map{
			"a": g.Slice{"1", "2"},
		})
		// map
		m, err = gstr.Parse("a=1&b=2&c=3")
		t.Assert(err, nil)
		t.Assert(m, g.Map{
			"a": "1",
			"b": "2",
			"c": "3",
		})
		m, err = gstr.Parse("a=1&a=2&c=3")
		t.Assert(err, nil)
		t.Assert(m, g.Map{
			"a": "2",
			"c": "3",
		})
		// map
		m, err = gstr.Parse("m[a]=1&m[b]=2&m[c]=3")
		t.Assert(err, nil)
		t.Assert(m, g.Map{
			"m": g.Map{
				"a": "1",
				"b": "2",
				"c": "3",
			},
		})
		m, err = gstr.Parse("m[a]=1&m[a]=2&m[b]=3")
		t.Assert(err, nil)
		t.Assert(m, g.Map{
			"m": g.Map{
				"a": "2",
				"b": "3",
			},
		})
		// map - slice
		m, err = gstr.Parse("m[a][]=1&m[a][]=2")
		t.Assert(err, nil)
		t.Assert(m, g.Map{
			"m": g.Map{
				"a": g.Slice{"1", "2"},
			},
		})
		m, err = gstr.Parse("m[a][b][]=1&m[a][b][]=2")
		t.Assert(err, nil)
		t.Assert(m, g.Map{
			"m": g.Map{
				"a": g.Map{
					"b": g.Slice{"1", "2"},
				},
			},
		})
		// map - complicated
		m, err = gstr.Parse("m[a1][b1][c1][d1]=1&m[a2][b2]=2&m[a3][b3][c3]=3")
		t.Assert(err, nil)
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
