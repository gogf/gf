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

func TestConverter_Struct(t *testing.T) {
	type tA struct {
		Val int
	}

	type tB struct {
		Val1 int32
		Val2 string
	}

	type tAA struct {
		ValTop int
		ValTA  tA
	}

	type tBB struct {
		ValTop int32
		ValTB  tB
	}

	type tCC struct {
		ValTop string
		ValTa  *tB
	}

	type tDD struct {
		ValTop string
		ValTa  tB
	}

	gtest.C(t, func(t *gtest.T) {
		a := &tA{
			Val: 1,
		}
		var b *tB
		err := gconv.Scan(a, &b)
		t.AssertNil(err)
		t.AssertNE(b, nil)
		t.Assert(b.Val1, 0)
		t.Assert(b.Val2, "")
	})

	gtest.C(t, func(t *gtest.T) {
		err := gconv.RegisterConverter(func(a tA) (b *tB, err error) {
			b = &tB{
				Val1: int32(a.Val),
				Val2: "abc",
			}
			return
		})
		t.AssertNil(err)
	})

	gtest.C(t, func(t *gtest.T) {
		a := &tA{
			Val: 1,
		}
		var b *tB
		err := gconv.Scan(a, &b)
		t.AssertNil(err)
		t.AssertNE(b, nil)
		t.Assert(b.Val1, 1)
		t.Assert(b.Val2, "abc")
	})

	gtest.C(t, func(t *gtest.T) {
		a := &tA{
			Val: 1,
		}
		var b *tB
		err := gconv.Scan(a, &b)
		t.AssertNil(err)
		t.AssertNE(b, nil)
		t.Assert(b.Val1, 1)
		t.Assert(b.Val2, "abc")
	})

	gtest.C(t, func(t *gtest.T) {
		a := &tA{
			Val: 1,
		}
		var b *tB
		err := gconv.Scan(a, &b)
		t.AssertNil(err)
		t.AssertNE(b, nil)
		t.Assert(b.Val1, 1)
		t.Assert(b.Val2, "abc")
	})

	gtest.C(t, func(t *gtest.T) {
		a := &tA{
			Val: 1,
		}
		var b *tB
		err := gconv.Scan(a, &b)
		t.AssertNil(err)
		t.AssertNE(b, nil)
		t.Assert(b.Val1, 1)
		t.Assert(b.Val2, "abc")
	})

	gtest.C(t, func(t *gtest.T) {
		aa := &tAA{
			ValTop: 123,
			ValTA:  tA{Val: 234},
		}
		var bb *tBB

		err := gconv.Scan(aa, &bb)
		t.AssertNil(err)
		t.AssertNE(bb, nil)
		t.Assert(bb.ValTop, 123)
		t.AssertNE(bb.ValTB.Val1, 234)

		err = gconv.RegisterConverter(func(a tAA) (b *tBB, err error) {
			b = &tBB{
				ValTop: int32(a.ValTop) + 2,
			}
			err = gconv.Scan(a.ValTA, &b.ValTB)
			return
		})
		t.AssertNil(err)

		err = gconv.Scan(aa, &bb)
		t.AssertNil(err)
		t.AssertNE(bb, nil)
		t.Assert(bb.ValTop, 125)
		t.Assert(bb.ValTB.Val1, 234)
		t.Assert(bb.ValTB.Val2, "abc")

	})

	gtest.C(t, func(t *gtest.T) {
		aa := &tAA{
			ValTop: 123,
			ValTA:  tA{Val: 234},
		}
		var cc *tCC
		err := gconv.Scan(aa, &cc)
		t.AssertNil(err)
		t.AssertNE(cc, nil)
		t.Assert(cc.ValTop, "123")
		t.AssertNE(cc.ValTa, nil)
		t.Assert(cc.ValTa.Val1, 234)
		t.Assert(cc.ValTa.Val2, "abc")
	})

	gtest.C(t, func(t *gtest.T) {
		aa := &tAA{
			ValTop: 123,
			ValTA:  tA{Val: 234},
		}

		var dd *tDD
		err := gconv.Scan(aa, &dd)
		t.AssertNil(err)
		t.AssertNE(dd, nil)
		t.Assert(dd.ValTop, "123")
		t.Assert(dd.ValTa.Val1, 234)
		t.Assert(dd.ValTa.Val2, "abc")
	})
}

func TestConverter_CustomBasicType_ToStruct(t *testing.T) {
	type CustomString string
	type CustomStruct struct {
		S string
	}
	gtest.C(t, func(t *gtest.T) {
		var (
			a CustomString = "abc"
			b *CustomStruct
		)
		err := gconv.Scan(a, &b)
		t.AssertNE(err, nil)
		t.Assert(b, nil)
	})

	gtest.C(t, func(t *gtest.T) {
		err := gconv.RegisterConverter(func(a CustomString) (b *CustomStruct, err error) {
			b = &CustomStruct{
				S: string(a),
			}
			return
		})
		t.AssertNil(err)
	})
	gtest.C(t, func(t *gtest.T) {
		var (
			a CustomString = "abc"
			b *CustomStruct
		)
		err := gconv.Scan(a, &b)
		t.AssertNil(err)
		t.AssertNE(b, nil)
		t.Assert(b.S, a)
	})
	gtest.C(t, func(t *gtest.T) {
		var (
			a CustomString = "abc"
			b *CustomStruct
		)
		err := gconv.Scan(&a, &b)
		t.AssertNil(err)
		t.AssertNE(b, nil)
		t.Assert(b.S, a)
	})
}
