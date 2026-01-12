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
		analyzer.packages = map[string]*goPackage{
			"github.com/other/package": {
				ImportPath: "github.com/other/package",
				Standard:   false,
			},
			"github.com/gogf/gf/cmd/gf/v2/internal": {
				ImportPath: "github.com/gogf/gf/cmd/gf/v2/internal",
				Standard:   false,
			},
			"fmt": {
				ImportPath: "fmt",
				Standard:   true,
			},
		}

		// Test using new FilterOptions system
		in := Input{
			Internal: false,
			External: true,
			NoStd:    true,
		}

		opts := analyzer.convertInputToFilterOptions(in)
		opts.Normalize(analyzer.modulePrefix)
		store := analyzer.buildPackageStore()

		// Test external package (should be included)
		externalPkg, ok := store.packages["github.com/other/package"]
		t.Assert(ok, true)
		t.Assert(opts.ShouldInclude(externalPkg), true)

		// Test internal package (should not be included)
		internalPkg, ok := store.packages["github.com/gogf/gf/cmd/gf/v2/internal"]
		t.Assert(ok, true)
		t.Assert(opts.ShouldInclude(internalPkg), false)

		// Test standard library (should not be included due to NoStd)
		stdPkg, ok := store.packages["fmt"]
		t.Assert(ok, true)
		t.Assert(opts.ShouldInclude(stdPkg), false)
	})
}

func TestExternalGrouping(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		analyzer := newAnalyzer()

		// Test external group extraction using shortName
		t.Assert(analyzer.shortName("github.com/user/repo", false), "github.com/user/repo")
		t.Assert(analyzer.shortName("golang.org/x/tools", false), "golang.org/x/tools")
		t.Assert(analyzer.shortName("fmt", false), "fmt")
		t.Assert(analyzer.shortName("simple", false), "simple")
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
