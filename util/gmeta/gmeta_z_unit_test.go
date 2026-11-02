// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gmeta_test

import (
	"testing"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/internal/json"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/gogf/gf/v2/util/gmeta"
	"github.com/gogf/gf/v2/util/gtag"
)

func TestMeta_Basic(t *testing.T) {
	type A struct {
		gmeta.Meta `tag:"123" orm:"456"`
		Id         int
		Name       string
	}

	gtest.C(t, func(t *gtest.T) {
		a := &A{
			Id:   100,
			Name: "john",
		}
		t.Assert(len(gmeta.Data(a)), 2)
		t.AssertEQ(gmeta.Get(a, "tag").String(), "123")
		t.AssertEQ(gmeta.Get(a, "orm").String(), "456")
		t.AssertEQ(gmeta.Get(a, "none"), nil)

		b, err := json.Marshal(a)
		t.AssertNil(err)
		t.Assert(b, `{"Id":100,"Name":"john"}`)
	})
}

func TestMeta_Convert_Map(t *testing.T) {
	type A struct {
		gmeta.Meta `tag:"123" orm:"456"`
		Id         int
		Name       string
	}

	gtest.C(t, func(t *gtest.T) {
		a := &A{
			Id:   100,
			Name: "john",
		}
		m := gconv.Map(a)
		t.Assert(len(m), 2)
		t.Assert(m[`Meta`], nil)
	})
}

func TestMeta_Json(t *testing.T) {
	type A struct {
		gmeta.Meta `tag:"123" orm:"456"`
		Id         int
	}

	gtest.C(t, func(t *gtest.T) {
		a := &A{
			Id: 100,
		}
		b, err := json.Marshal(a)
		t.AssertNil(err)
		t.Assert(string(b), `{"Id":100}`)
	})
}

func TestMeta_Cache(t *testing.T) {
	type A struct {
		gmeta.Meta `tag:"123" orm:"456" path:"/a" method:"get" summary:"s" dc:"d"`
		Id         int
	}
	type B struct {
		Id int
	}

	// Modifying the returned map must not affect later calls.
	gtest.C(t, func(t *gtest.T) {
		data := gmeta.Data(&A{})
		data["tag"] = "changed"
		delete(data, "orm")
		t.Assert(gmeta.Data(A{}), g.MapStrStr{
			"tag": "123", "orm": "456", "path": "/a", "method": "get", "summary": "s", "dc": "d",
		})
		t.AssertEQ(gmeta.Get(&A{}, "tag").String(), "123")
		t.AssertEQ(gmeta.Get(&A{}, "orm").String(), "456")
	})
	// Struct without Meta and non-struct objects.
	gtest.C(t, func(t *gtest.T) {
		t.AssertEQ(len(gmeta.Data(B{})), 0)
		t.AssertNE(gmeta.Data(B{}), nil)
		t.AssertEQ(gmeta.Get(B{}, "tag"), nil)
		t.AssertNil(gmeta.Data(1))
		t.AssertEQ(gmeta.Get(1, "tag"), nil)
	})
	// Get reads the cached metadata instead of parsing the tag on every call,
	// so its allocations do not grow with the number of meta tags.
	gtest.C(t, func(t *gtest.T) {
		type Small struct {
			gmeta.Meta `tag:"1"`
		}
		type Large struct {
			gmeta.Meta `tag:"1" k1:"1" k2:"2" k3:"3" k4:"4" k5:"5" k6:"6" k7:"7" k8:"8" k9:"9" k10:"10" k11:"11" k12:"12" k13:"13" k14:"14" k15:"15" k16:"16"`
		}
		var (
			small       = &Small{}
			large       = &Large{}
			smallAllocs = testing.AllocsPerRun(100, func() { gmeta.Get(small, "tag") })
			largeAllocs = testing.AllocsPerRun(100, func() { gmeta.Get(large, "tag") })
		)
		t.AssertEQ(largeAllocs, smallAllocs)
	})
}

func TestMeta_Cache_TagPlaceholder(t *testing.T) {
	type A struct {
		gmeta.Meta `path:"/{TestMetaCachePrefix}/user" method:"get"`
	}
	gtest.C(t, func(t *gtest.T) {
		t.Assert(gmeta.Get(A{}, "path").String(), "/{TestMetaCachePrefix}/user")
		t.Assert(gmeta.Data(A{})["path"], "/{TestMetaCachePrefix}/user")

		gtag.Set("TestMetaCachePrefix", "v1")
		t.Assert(gmeta.Get(A{}, "path").String(), "/v1/user")
		t.Assert(gmeta.Data(A{}), g.MapStrStr{"path": "/v1/user", "method": "get"})

		gtag.SetOver("TestMetaCachePrefix", "v2")
		t.Assert(gmeta.Get(&A{}, "path").String(), "/v2/user")
		t.Assert(gmeta.Data(&A{}), g.MapStrStr{"path": "/v2/user", "method": "get"})
	})
}
