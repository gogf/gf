// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package mysql

import (
	"testing"

	"github.com/gogf/gf/v2/test/gtest"
)

func TestFormatPartitionClause(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		d := &Driver{}
		t.Assert(d.FormatPartitionClause([]string{"p0"}), "PARTITION (p0)")
		t.Assert(d.FormatPartitionClause([]string{"p0", "p1"}), "PARTITION (p0,p1)")
		t.Assert(d.FormatPartitionClause([]string{"p0", "p1", "p2"}), "PARTITION (p0,p1,p2)")
	})
}
