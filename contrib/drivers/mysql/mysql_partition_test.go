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
		var d = &Driver{}
		t.Assert(d.FormatPartitionClause([]string{"p0"}), "PARTITION (`p0`)")
		t.Assert(d.FormatPartitionClause([]string{"p0", "p1"}), "PARTITION (`p0`,`p1`)")
		t.Assert(d.FormatPartitionClause([]string{"p0", "p1", "p2"}), "PARTITION (`p0`,`p1`,`p2`)")
	})
	// No names leaves the table expression without a clause.
	gtest.C(t, func(t *gtest.T) {
		var d = &Driver{}
		t.Assert(d.FormatPartitionClause(nil), "")
		t.Assert(d.FormatPartitionClause([]string{}), "")
	})
	// A name that is not a plain identifier stays inside the clause.
	gtest.C(t, func(t *gtest.T) {
		var d = &Driver{}
		t.Assert(
			d.FormatPartitionClause([]string{"p0) UNION SELECT * FROM user --"}),
			"PARTITION (`p0) UNION SELECT * FROM user --`)",
		)
		t.Assert(
			d.FormatPartitionClause([]string{"p0`,(SELECT 1)) --"}),
			"PARTITION (`p0``,(SELECT 1)) --`)",
		)
		t.Assert(d.FormatPartitionClause([]string{"p 0"}), "PARTITION (`p 0`)")
	})
	// A name is trimmed, so a stray space does not name a different partition.
	gtest.C(t, func(t *gtest.T) {
		var d = &Driver{}
		t.Assert(d.FormatPartitionClause([]string{" p0 "}), "PARTITION (`p0`)")
	})
	// An empty name is kept as an empty identifier for the server to reject, instead of
	// being dropped: dropping it would silently read every partition.
	gtest.C(t, func(t *gtest.T) {
		var d = &Driver{}
		t.Assert(d.FormatPartitionClause([]string{""}), "PARTITION (``)")
		t.Assert(d.FormatPartitionClause([]string{"p0", ""}), "PARTITION (`p0`,``)")
	})
}
