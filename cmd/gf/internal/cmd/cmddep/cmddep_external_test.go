// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package cmddep

import (
	"testing"

	"github.com/gogf/gf/v2/test/gtest"
)

func TestExternalDependencyAnalysis(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		analyzer := newAnalyzer()
		analyzer.modulePrefix = "github.com/gogf/gf/cmd/gf/v2"

		// Test shouldInclude with external dependencies
		in := Input{
			Internal: false,
			External: true,
			NoStd:    true,
		}

		// Test external package (should be included)
		t.Assert(analyzer.shouldInclude("github.com/other/package", in), true)

		// Test internal package (should not be included)
		t.Assert(analyzer.shouldInclude("github.com/gogf/gf/cmd/gf/v2/internal", in), false)

		// Test standard library (should not be included due to NoStd)
		t.Assert(analyzer.shouldInclude("fmt", in), false)
	})
}

func TestExternalGrouping(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		analyzer := newAnalyzer()

		// Test external group extraction
		t.Assert(analyzer.getExternalGroup("github.com/user/repo"), "github.com/user")
		t.Assert(analyzer.getExternalGroup("golang.org/x/tools"), "golang.org")
		t.Assert(analyzer.getExternalGroup("fmt"), "stdlib")
		t.Assert(analyzer.getExternalGroup("simple"), "simple")
	})
}

func TestDependencyStats(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		analyzer := newAnalyzer()
		analyzer.modulePrefix = "github.com/gogf/gf/cmd/gf/v2"

		// Add test packages
		analyzer.packages = map[string]*goPackage{
			"github.com/gogf/gf/cmd/gf/v2/internal": {
				ImportPath: "github.com/gogf/gf/cmd/gf/v2/internal",
				Standard:   false,
			},
			"github.com/external/package": {
				ImportPath: "github.com/external/package",
				Standard:   false,
			},
			"fmt": {
				ImportPath: "fmt",
				Standard:   true,
			},
		}

		in := Input{
			Internal: true,
			External: true,
			NoStd:    false,
		}

		stats := analyzer.getDependencyStats(in)
		t.Assert(stats["total"], 3)
		t.Assert(stats["internal"], 1)
		t.Assert(stats["external"], 1)
		t.Assert(stats["stdlib"], 1)
	})
}