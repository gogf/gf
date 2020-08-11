// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gjson

import (
	"github.com/gogf/gf/test/gtest"
	"testing"
)

func Test_checkDataType(t *testing.T) {
	data := []byte(`
bb           = """
                   dig := dig;                         END;"""
`)
	gtest.C(t, func(t *gtest.T) {
		t.Assert(checkDataType(data), "toml")
	})
}
