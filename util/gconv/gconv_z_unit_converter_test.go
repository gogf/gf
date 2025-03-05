// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gconv_test

import (
	"database/sql"
	"fmt"
	"reflect"
	"testing"

	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/gconv"
)

type converterStructInTest struct {
	Name string
}

type converterStructOutTest struct {
	Place string
}

func TestRegisterConverter(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		err := gconv.RegisterConverter(
			func(in converterStructInTest) (*converterStructOutTest, error) {
				return &converterStructOutTest{
					Place: in.Name,
				}, nil
			},
		)
		t.AssertNil(err)
	})

	// Test failure cases.
	gtest.C(t, func(t *gtest.T) {
		var err error
		err = gconv.RegisterConverter(123)
		t.AssertNE(err, nil)

		err = gconv.RegisterConverter(func() {})
		t.AssertNE(err, nil)

		err = gconv.RegisterConverter(
			func(in *converterStructInTest) (*converterStructOutTest, error) {
				return nil, nil
			},
		)
		t.AssertNE(err, nil)

		err = gconv.RegisterConverter(
			func(in converterStructInTest) (converterStructOutTest, error) {
				return converterStructOutTest{}, nil
			},
		)
		t.AssertNE(err, nil)

		err = gconv.RegisterConverter(
			func(in converterStructInTest) (*converterStructOutTest, error) {
				return nil, nil
			},
		)
		t.AssertNE(err, nil)
	})

	gtest.C(t, func(t *gtest.T) {
		var (
			converterStructIn  = converterStructInTest{"小行星带"}
			converterStructOut converterStructOutTest
		)
		err := gconv.Scan(converterStructIn, &converterStructOut)
		t.AssertNil(err)
		t.Assert(converterStructOut.Place, converterStructIn.Name)
	})
}

func TestConvertWithRefer(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.AssertEQ(gconv.ConvertWithRefer("1", 100), 1)
		t.AssertEQ(gconv.ConvertWithRefer("1.01", 1.111), 1.01)
		t.AssertEQ(gconv.ConvertWithRefer("1.01", "1.111"), "1.01")
		t.AssertEQ(gconv.ConvertWithRefer("1.01", false), true)
		t.AssertNE(gconv.ConvertWithRefer("1.01", false), false)
	})
}

func testAnyToMyInt(from any, to reflect.Value) error {
	switch x := from.(type) {
	case int:
		to.SetInt(123456)
	default:
		return fmt.Errorf("unsupported type %T(%v)", x, x)
	}
	return nil
}

func testAnyToSqlNullType(_ any, to reflect.Value) error {
	if to.Kind() != reflect.Ptr {
		to = to.Addr()
	}
	return to.Interface().(sql.Scanner).Scan(123456)
}

func TestNewConverter(t *testing.T) {
	type Dst[T any] struct {
		A T
	}
	gtest.C(t, func(t *gtest.T) {
		conv := gconv.NewConverter()
		conv.RegisterAnyConverterFunc(testAnyToMyInt, reflect.TypeOf((*myInt)(nil)))
		var dst Dst[myInt]
		err := conv.Struct(map[string]any{
			"a": 1200,
		}, &dst, gconv.StructOption{})
		t.AssertNil(err)
		t.Assert(dst, Dst[myInt]{
			A: 123456,
		})
	})
	gtest.C(t, func(t *gtest.T) {
		conv := gconv.NewConverter()
		conv.RegisterAnyConverterFunc(testAnyToMyInt, reflect.TypeOf((myInt)(0)))
		var dst Dst[*myInt]
		err := conv.Struct(map[string]any{
			"a": 1200,
		}, &dst, gconv.StructOption{})
		t.AssertNil(err)
		t.Assert(*dst.A, 123456)
	})

	gtest.C(t, func(t *gtest.T) {
		conv := gconv.NewConverter()
		conv.RegisterAnyConverterFunc(testAnyToSqlNullType, reflect.TypeOf((*sql.Scanner)(nil)))
		type sqlNullDst struct {
			A sql.Null[int]
			B sql.Null[float32]
			C sql.NullInt64
			D sql.NullString
		}
		var dst sqlNullDst
		err := conv.Struct(map[string]any{
			"a": 12,
			"b": 34,
			"c": 56,
			"d": "sqlNull",
		}, &dst, gconv.StructOption{})
		t.AssertNil(err)
		t.Assert(dst, sqlNullDst{
			A: sql.Null[int]{V: 123456, Valid: true},
			B: sql.Null[float32]{V: 123456, Valid: true},
			C: sql.NullInt64{Int64: 123456, Valid: true},
			D: sql.NullString{String: "123456", Valid: true},
		})
	})
}
