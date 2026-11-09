// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gtest_test

import (
	"testing"
	"time"

	"github.com/gogf/gf/v2/test/gtest"
)

// https://github.com/gogf/gf/issues/4953
func Test_Issue4953(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			big   = int64(1) << 33
			small = int64(1)<<32 + 5
		)
		t.AssertGT(big, small)
		t.AssertGE(big, small)
		t.AssertLT(small, big)
		t.AssertLE(small, big)
		t.AssertGT(uint64(big), uint64(small))
		t.AssertGE(uint64(big), uint64(small))
		t.AssertLT(uint64(small), uint64(big))
		t.AssertLE(uint64(small), uint64(big))
		t.AssertGT(3*time.Second, 0)
		t.AssertLE(2*time.Second, 6*time.Second)
	})
}
