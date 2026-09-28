// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package pgsql

import (
	"context"
	"strings"
	"testing"

	"github.com/lib/pq/oid"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/test/gtest"
)

// Test_LocalTypeCoverage asserts that every type name the underlying driver can report has
// an explicit local type. A name missing from localTypeMap reaches the keyword matching of
// the core, which infers the type from substrings of the name and is wrong for any name
// that merely embeds a keyword, such as `point` embedding "int".
//
// When this test fails, the driver has gained type names since the map was written. Add
// them to localTypeMap rather than relying on the core to guess them.
func Test_LocalTypeCoverage(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var missing []string
		for _, typeName := range oid.TypeName {
			name := strings.ToLower(typeName)
			if _, ok := localTypeMap[name]; !ok {
				missing = append(missing, name)
			}
		}
		t.Assert(missing, nil)
	})
}

// Test_LocalTypeMapHasNoUnknownName asserts the reverse direction, so that the map does not
// keep entries for names the driver can no longer report. An entry that is deliberately not
// reportable belongs in extraTypeNames, which documents why it is kept.
func Test_LocalTypeMapHasNoUnknownName(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var reportable = make(map[string]bool, len(oid.TypeName))
		for _, typeName := range oid.TypeName {
			reportable[strings.ToLower(typeName)] = true
		}
		var unknown []string
		for name := range localTypeMap {
			if _, ok := extraTypeNames[name]; ok {
				continue
			}
			if !reportable[name] {
				unknown = append(unknown, name)
			}
		}
		t.Assert(unknown, nil)
	})
}

// Test_ExtraTypeNamesAreMapped asserts every documented extra actually has a mapping, so
// that the two stay in step.
func Test_ExtraTypeNamesAreMapped(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		for name := range extraTypeNames {
			_, ok := localTypeMap[name]
			t.Assert(ok, true)
		}
	})
}

// Test_LocalTypeBitKeepsPrecision asserts that bit keeps being decided by its precision,
// which the lookup by name in localTypeMap cannot see: bit(1) is a boolean.
func Test_LocalTypeBitKeepsPrecision(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			ctx    = context.Background()
			driver = &Driver{Core: &gdb.Core{}}
		)
		localType, err := driver.CheckLocalTypeForField(ctx, "bit(1)", nil)
		t.AssertNil(err)
		t.Assert(localType, gdb.LocalTypeBool)

		localType, err = driver.CheckLocalTypeForField(ctx, "bit(8)", nil)
		t.AssertNil(err)
		t.Assert(localType, gdb.LocalTypeInt64Bytes)
	})
}
