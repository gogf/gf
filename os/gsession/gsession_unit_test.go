// Copyright 2019 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gsession

import (
	"testing"

	"github.com/gogf/gf/test/gtest"
)

func Test_NewSessionId(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		id1 := NewSessionId()
		id2 := NewSessionId()
		t.AssertNE(id1, id2)
		t.Assert(len(id1), 18)
	})
}
