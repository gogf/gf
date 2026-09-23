// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package oracle

import (
	"strings"
	"testing"

	go_ora "github.com/sijms/go-ora/v2"

	"github.com/gogf/gf/v2/test/gtest"
)

// reportableTypeNames returns the type names the driver can report, by asking every TNSType
// code for its name and keeping the ones it knows.
func reportableTypeNames() []string {
	var (
		names []string
		seen  = map[string]bool{}
	)
	for code := 0; code < 256; code++ {
		name := go_ora.TNSType(code).String()
		if strings.HasPrefix(name, "TNSType(") {
			continue
		}
		name = strings.ToLower(name)
		if !seen[name] {
			seen[name] = true
			names = append(names, name)
		}
	}
	return names
}

// Test_LocalTypeCoverage asserts that every type name the underlying driver can report has
// an explicit local type. A name missing from localTypeMap reaches the keyword matching of
// the core, which infers the type from substrings of the name and is wrong for any name
// that merely embeds a keyword, such as `IntervalDS_DTY` embedding "int".
//
// When this test fails, the driver has gained type names since the map was written. Add
// them to localTypeMap rather than relying on the core to guess them.
func Test_LocalTypeCoverage(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var missing []string
		for _, name := range reportableTypeNames() {
			if _, ok := localTypeMap[name]; !ok {
				missing = append(missing, name)
			}
		}
		t.Assert(missing, nil)
	})
}

// Test_LocalTypeMapHasNoUnknownName asserts the reverse direction, so that the map does not
// keep entries for names the driver can no longer report.
func Test_LocalTypeMapHasNoUnknownName(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var reportable = map[string]bool{}
		for _, name := range reportableTypeNames() {
			reportable[name] = true
		}
		var unknown []string
		for name := range localTypeMap {
			if !reportable[name] {
				unknown = append(unknown, name)
			}
		}
		t.Assert(unknown, nil)
	})
}
