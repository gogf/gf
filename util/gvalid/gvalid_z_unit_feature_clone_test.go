// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gvalid_test

import (
	"context"
	"testing"

	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/gvalid"
)

// It tests that the chaining operations of Validator work on independent
// validators, and the base validator is not affected by them.
func Test_Validator_Clone_Independence(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			ctx  = context.Background()
			base = gvalid.New()
			err1 = base.Data(1).Rules("min:5").Run(ctx)
			err2 = base.Data(10).Rules("min:5").Run(ctx)
			err3 = base.Data(1).Rules("min:5").Messages("too small").Run(ctx)
			errB = base.Run(ctx)
		)
		// Validations chained from the same base validator are independent.
		t.AssertNE(err1, nil)
		t.AssertNil(err2)
		// Custom messages of one validation are not leaked into another.
		t.Assert(err3.Error(), "too small")
		// The base validator itself keeps no data and rules.
		t.AssertNE(errB, nil)
		t.Assert(errB.Error(), "no data passed for validation")
	})
}
