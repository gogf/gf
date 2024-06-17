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
