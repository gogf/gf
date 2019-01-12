// Copyright 2019 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

package gvalid

import (
    "gitee.com/johng/gf/g/util/gtest"
    "testing"
)

func Test_Phone(t *testing.T) {
    gtest.Case(t, func() {
        err1 := Check("1361990897", "phone", nil)
        err2 := Check("13619908979", "phone", nil)
        err3 := Check("16719908979", "phone", nil)
        err4 := Check("19719908989", "phone", nil)
        gtest.AssertNE(err1, nil)
        gtest.Assert(err2, nil)
        gtest.Assert(err3, nil)
        gtest.Assert(err4, nil)
    })
}
