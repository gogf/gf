// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gconv_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/gconv"
)

func TestConverter_ConvertWithRefer(t *testing.T) {
	type tA struct {
		Val int
	}

	type tB struct {
		Val1 int32
		Val2 string
	}

	gtest.C(t, func(t *gtest.T) {
		err := gconv.RegisterConverter(func(a tA) (b *tB, err error) {
			b = &tB{
				Val1: int32(a.Val),
				Val2: "abcd",
			}
			return
		})
		t.AssertNil(err)
	})

	gtest.C(t, func(t *gtest.T) {
		a := &tA{
			Val: 1,
		}
		var b tB
		result := gconv.ConvertWithRefer(a, &b)
		t.Assert(result.(*tB), &tB{
			Val1: 1,
			Val2: "abcd",
		})
	})

	gtest.C(t, func(t *gtest.T) {
		a := &tA{
			Val: 1,
		}
		var b tB
		result := gconv.ConvertWithRefer(a, b)
		t.Assert(result.(tB), tB{
			Val1: 1,
			Val2: "abcd",
		})
	})
}

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

	type tEE struct {
		Val1 time.Time  `json:"val1"`
		Val2 *time.Time `json:"val2"`
		Val3 *time.Time `json:"val3"`
	}

	type tFF struct {
		Val1 json.RawMessage            `json:"val1"`
		Val2 []json.RawMessage          `json:"val2"`
		Val3 map[string]json.RawMessage `json:"val3"`
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

	// fix: https://github.com/gogf/gf/issues/2665
	gtest.C(t, func(t *gtest.T) {
		aa := &tEE{}

		var tmp = map[string]any{
			"val1": "2023-04-15 19:10:00 +0800 CST",
			"val2": "2023-04-15 19:10:00 +0800 CST",
			"val3": "2006-01-02T15:04:05Z07:00",
		}
		err := gconv.Struct(tmp, aa)
		t.AssertNil(err)
		t.AssertNE(aa, nil)
		t.Assert(aa.Val1.Local(), gtime.New("2023-04-15 19:10:00 +0800 CST").Local().Time)
		t.Assert(aa.Val2.Local(), gtime.New("2023-04-15 19:10:00 +0800 CST").Local().Time)
		t.Assert(aa.Val3.Local(), gtime.New("2006-01-02T15:04:05Z07:00").Local().Time)
	})

	// fix: https://github.com/gogf/gf/issues/3006
	gtest.C(t, func(t *gtest.T) {
		ff := &tFF{}
		var tmp = map[string]any{
			"val1": map[string]any{"hello": "world"},
			"val2": []any{map[string]string{"hello": "world"}},
			"val3": map[string]map[string]string{"val3": {"hello": "world"}},
		}

		err := gconv.Struct(tmp, ff)
		t.AssertNil(err)
		t.AssertNE(ff, nil)
		t.Assert(ff.Val1, []byte(`{"hello":"world"}`))
		t.AssertEQ(len(ff.Val2), 1)
		t.Assert(ff.Val2[0], []byte(`{"hello":"world"}`))
		t.AssertEQ(len(ff.Val3), 1)
		t.Assert(ff.Val3["val3"], []byte(`{"hello":"world"}`))
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

// fix: https://github.com/gogf/gf/issues/3099
func TestConverter_CustomTimeType_ToStruct(t *testing.T) {
	type timestamppb struct {
		S string
	}
	type CustomGTime struct {
		T *gtime.Time
	}
	type CustomPbTime struct {
		T *timestamppb
	}
	gtest.C(t, func(t *gtest.T) {
		var (
			a = CustomGTime{
				T: gtime.NewFromStrFormat("2023-10-26", "Y-m-d"),
			}
			b *CustomPbTime
		)
		err := gconv.Scan(a, &b)
		t.AssertNil(err)
		t.AssertNE(b, nil)
		t.Assert(b.T.S, "")
	})

	gtest.C(t, func(t *gtest.T) {
		err := gconv.RegisterConverter(func(in gtime.Time) (*timestamppb, error) {
			return &timestamppb{
				S: in.Local().Format("Y-m-d"),
			}, nil
		})
		t.AssertNil(err)
		err = gconv.RegisterConverter(func(in timestamppb) (*gtime.Time, error) {
			return gtime.NewFromStr(in.S), nil
		})
		t.AssertNil(err)
	})
	gtest.C(t, func(t *gtest.T) {
		var (
			a = CustomGTime{
				T: gtime.NewFromStrFormat("2023-10-26", "Y-m-d"),
			}
			b *CustomPbTime
			c *CustomGTime
		)
		err := gconv.Scan(a, &b)
		t.AssertNil(err)
		t.AssertNE(b, nil)
		t.AssertNE(b.T, nil)

		err = gconv.Scan(b, &c)
		t.AssertNil(err)
		t.AssertNE(c, nil)
		t.AssertNE(c.T, nil)
		t.AssertEQ(a.T.Timestamp(), c.T.Timestamp())
	})
}
