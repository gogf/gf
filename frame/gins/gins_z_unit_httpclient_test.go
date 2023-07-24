// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gins_test

import (
	"fmt"
	"testing"

	"github.com/gogf/gf/v2/frame/gins"
	"github.com/gogf/gf/v2/test/gtest"
)

func Test_Client(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			c  = gins.HttpClient()
			c1 = gins.HttpClient("c1")
			c2 = gins.HttpClient("c2")
		)
		c.SetAgent("test1")
		c.SetAgent("test2")
		t.AssertNE(fmt.Sprintf(`%p`, c), fmt.Sprintf(`%p`, c1))
		t.AssertNE(fmt.Sprintf(`%p`, c), fmt.Sprintf(`%p`, c2))
		t.AssertNE(fmt.Sprintf(`%p`, c1), fmt.Sprintf(`%p`, c2))
	})
}
