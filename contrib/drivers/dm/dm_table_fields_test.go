// Copyright 2019 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// This file tests table metadata helpers for the DM driver.

package dm

import (
	"testing"

	"github.com/gogf/gf/v2/test/gtest"
)

func TestTableNameCandidatesForMetadata(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.Assert(tableNameCandidatesForMetadata("A_TABLES"), []string{"A_TABLES"})
		t.Assert(tableNameCandidatesForMetadata("A_tables"), []string{"A_TABLES", "A_tables"})
		t.Assert(tableNameCandidatesForMetadata(`"A_tables"`), []string{"A_TABLES", "A_tables"})
	})
}
