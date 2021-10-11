// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gcode_test

import (
	"github.com/gogf/gf/v2/errors/gcode"
	"testing"

	"github.com/gogf/gf/v2/test/gtest"
)

func Test_Nil(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		c := gcode.New(1, "custom error", "detailed description")
		t.Assert(c.Code(), 1)
		t.Assert(c.Message(), "custom error")
		t.Assert(c.Detail(), "detailed description")
	})
}
