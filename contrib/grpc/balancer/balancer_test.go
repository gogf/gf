// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package balancer_test

import (
	"testing"

	"github.com/gogf/gf/v2/net/gsel"
)

func Test_Register(t *testing.T) {
	Register("test", gsel.NewSelectorRandom())
}
