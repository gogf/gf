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

type impUnmarshalValue struct{}

func (*impUnmarshalValue) UnmarshalValue(interface{}) error {
	return nil
}

func TestIUnmarshalValue(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var v any = &impUnmarshalValue{}
		_, ok := (v).(gconv.IUnmarshalValue)
		t.AssertEQ(ok, true)
	})
}
