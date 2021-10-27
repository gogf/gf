// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gbuild_test

import (
	"github.com/gogf/gf/v2/os/gbuild"
	"github.com/gogf/gf/v2/test/gtest"
	"testing"
)

func Test_Info(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.Assert(gbuild.Info(), map[string]string{
			"gf":   "",
			"go":   "",
			"git":  "",
			"time": "",
		})
	})
}

func Test_Get(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.Assert(gbuild.Get(`none`), nil)
	})
	gtest.C(t, func(t *gtest.T) {
		t.Assert(gbuild.Get(`none`, 1), 1)
	})
}

func Test_Map(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.Assert(gbuild.Map(), map[string]interface{}{})
	})
}
