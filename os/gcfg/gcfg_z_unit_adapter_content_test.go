// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// go test *.go -bench=".*" -benchmem

package gcfg_test

import (
	"testing"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcfg"
	"github.com/gogf/gf/v2/test/gtest"
)

func TestAdapterContent_Available_Get_Data(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		adapter, err := gcfg.NewAdapterContent()
		t.AssertNil(err)
		t.Assert(adapter.Available(ctx), false)
	})
	gtest.C(t, func(t *gtest.T) {
		content := `{"a": 1, "b": 2, "c": {"d": 3}}`
		adapter, err := gcfg.NewAdapterContent(content)
		t.AssertNil(err)

		c := gcfg.NewWithAdapter(adapter)
		t.Assert(c.Available(ctx), true)
		t.Assert(c.MustGet(ctx, "a"), 1)
		t.Assert(c.MustGet(ctx, "b"), 2)
		t.Assert(c.MustGet(ctx, "c.d"), 3)
		t.Assert(c.MustGet(ctx, "d"), nil)
		t.Assert(c.MustData(ctx), g.Map{
			"a": 1,
			"b": 2,
			"c": g.Map{
				"d": 3,
			},
		})
	})
}

func TestAdapterContent_SetContent(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		adapter, err := gcfg.NewAdapterContent()
		t.AssertNil(err)
		t.Assert(adapter.Available(ctx), false)

		content := `{"a": 1, "b": 2, "c": {"d": 3}}`
		err = adapter.SetContent(content)
		t.AssertNil(err)
		c := gcfg.NewWithAdapter(adapter)
		t.Assert(c.Available(ctx), true)
		t.Assert(c.MustGet(ctx, "a"), 1)
		t.Assert(c.MustGet(ctx, "b"), 2)
		t.Assert(c.MustGet(ctx, "c.d"), 3)
		t.Assert(c.MustGet(ctx, "d"), nil)
		t.Assert(c.MustData(ctx), g.Map{
			"a": 1,
			"b": 2,
			"c": g.Map{
				"d": 3,
			},
		})
	})

}
