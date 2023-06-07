// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gtag_test

import (
	"fmt"
	"testing"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/gtag"
	"github.com/gogf/gf/v2/util/guid"
)

func Test_Set_Get(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		k := guid.S()
		v := guid.S()
		gtag.Set(k, v)
		t.Assert(gtag.Get(k), v)
	})
}

func Test_SetOver_Get(t *testing.T) {
	// panic by Set
	gtest.C(t, func(t *gtest.T) {
		var (
			k  = guid.S()
			v1 = guid.S()
			v2 = guid.S()
		)
		gtag.Set(k, v1)
		t.Assert(gtag.Get(k), v1)
		defer func() {
			t.AssertNE(recover(), nil)
		}()
		gtag.Set(k, v2)
	})
	gtest.C(t, func(t *gtest.T) {
		var (
			k  = guid.S()
			v1 = guid.S()
			v2 = guid.S()
		)
		gtag.SetOver(k, v1)
		t.Assert(gtag.Get(k), v1)
		gtag.SetOver(k, v2)
		t.Assert(gtag.Get(k), v2)
	})
}

func Test_Sets_Get(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			k1 = guid.S()
			k2 = guid.S()
			v1 = guid.S()
			v2 = guid.S()
		)
		gtag.Sets(g.MapStrStr{
			k1: v1,
			k2: v2,
		})
		t.Assert(gtag.Get(k1), v1)
		t.Assert(gtag.Get(k2), v2)
	})
}

func Test_SetsOver_Get(t *testing.T) {
	// panic by Sets
	gtest.C(t, func(t *gtest.T) {
		var (
			k1 = guid.S()
			k2 = guid.S()
			v1 = guid.S()
			v2 = guid.S()
			v3 = guid.S()
		)
		gtag.Sets(g.MapStrStr{
			k1: v1,
			k2: v2,
		})
		t.Assert(gtag.Get(k1), v1)
		t.Assert(gtag.Get(k2), v2)
		defer func() {
			t.AssertNE(recover(), nil)
		}()
		gtag.Sets(g.MapStrStr{
			k1: v3,
			k2: v3,
		})
	})
	gtest.C(t, func(t *gtest.T) {
		var (
			k1 = guid.S()
			k2 = guid.S()
			v1 = guid.S()
			v2 = guid.S()
			v3 = guid.S()
		)
		gtag.SetsOver(g.MapStrStr{
			k1: v1,
			k2: v2,
		})
		t.Assert(gtag.Get(k1), v1)
		t.Assert(gtag.Get(k2), v2)
		gtag.SetsOver(g.MapStrStr{
			k1: v3,
			k2: v3,
		})
		t.Assert(gtag.Get(k1), v3)
		t.Assert(gtag.Get(k2), v3)
	})
}

func Test_Parse(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			k1      = guid.S()
			k2      = guid.S()
			v1      = guid.S()
			v2      = guid.S()
			content = fmt.Sprintf(`this is {%s} and {%s}`, k1, k2)
			expect  = fmt.Sprintf(`this is %s and %s`, v1, v2)
		)
		gtag.Sets(g.MapStrStr{
			k1: v1,
			k2: v2,
		})
		t.Assert(gtag.Parse(content), expect)
	})
}

func Test_SetGlobalEnums(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		oldEnumsJson, err := gtag.GetGlobalEnums()
		t.AssertNil(err)

		err = gtag.SetGlobalEnums(`{"k8s.io/apimachinery/pkg/api/resource.Format": [
        "BinarySI",
        "DecimalExponent",
        "DecimalSI"
    ]}`)
		t.AssertNil(err)
		t.Assert(gtag.GetEnumsByType("k8s.io/apimachinery/pkg/api/resource.Format"), `[
        "BinarySI",
        "DecimalExponent",
        "DecimalSI"
    ]`)
		t.AssertNil(gtag.SetGlobalEnums(oldEnumsJson))
	})
}
