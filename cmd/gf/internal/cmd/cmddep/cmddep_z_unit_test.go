// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package cmddep

import (
	"testing"

	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/test/gtest"
)

// Test data model creation and classification
func Test_PackageInfo_Creation(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		pkg := &PackageInfo{
			ImportPath:   "github.com/gogf/gf/v2/os/gfile",
			ModulePath:   "github.com/gogf/gf/v2",
			Kind:         KindInternal,
			Tier:         2,
			Imports:      []string{"fmt", "os"},
			IsStdLib:     false,
			IsModuleRoot: false,
		}
		t.Assert(pkg != nil, true)
		t.Assert(pkg.Kind, KindInternal)
		t.Assert(pkg.Tier, 2)
		t.Assert(len(pkg.Imports), 2)
	})
}

// Test FilterOptions normalization
func Test_FilterOptions_Normalize(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		opts := &FilterOptions{
			IncludeInternal: false,
			IncludeExternal: false,
			IncludeStdLib:   false,
		}
		opts.Normalize("github.com/gogf/gf/v2")

		// After normalization, internal should be included by default
		t.Assert(opts.IncludeInternal, true)
		t.Assert(opts.IncludeExternal, false)
	})
}

// Test ShouldInclude decision logic
func Test_FilterOptions_ShouldInclude(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// Test internal package inclusion
		opts := &FilterOptions{
			IncludeInternal: true,
			IncludeExternal: false,
			IncludeStdLib:   true,
		}

		internalPkg := &PackageInfo{
			Kind:     KindInternal,
			IsStdLib: false,
		}
		externalPkg := &PackageInfo{
			Kind:     KindExternal,
			IsStdLib: false,
		}
		stdlibPkg := &PackageInfo{
			Kind:     KindStdLib,
			IsStdLib: true,
		}

		t.Assert(opts.ShouldInclude(internalPkg), true)
		t.Assert(opts.ShouldInclude(externalPkg), false)
		t.Assert(opts.ShouldInclude(stdlibPkg), true)
	})
}

// Test TraversalContext Visit tracking
func Test_TraversalContext_Visit(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := &TraversalContext{
			visited: make(map[string]bool),
		}

		// First visit should return false
		t.Assert(ctx.Visit("pkg1"), false)
		// Second visit should return true
		t.Assert(ctx.Visit("pkg1"), true)
		// New package should return false
		t.Assert(ctx.Visit("pkg2"), false)
	})
}

// Test guessModuleRoot function
func Test_GuessModuleRoot(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		tests := []struct {
			input  string
			expect string
		}{
			{"github.com/gogf/gf", "github.com/gogf/gf"},
			{"github.com/gogf/gf/v2", "github.com/gogf/gf/v2"},
			{"github.com/gogf/gf/v2/os/gfile", "github.com/gogf/gf/v2"},
			{"github.com/gogf/gf/v2/os", "github.com/gogf/gf/v2"},
		}

		for _, test := range tests {
			result := guessModuleRoot(test.input)
			t.AssertEQ(result, test.expect)
		}
	})
}

// Test input to filter options conversion
func Test_ConvertInputToFilterOptions(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		a := newAnalyzer()

		input := Input{
			Internal: true,
			External: false,
			NoStd:    true,
			Module:   false,
			Direct:   false,
			Depth:    3,
		}

		opts := a.convertInputToFilterOptions(input)

		t.Assert(opts.IncludeInternal, true)
		t.Assert(opts.IncludeExternal, false)
		t.Assert(opts.IncludeStdLib, false) // NoStd=true means !IncludeStdLib
		t.Assert(opts.Depth, 3)
	})
}

func Test_Dep_Tree(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := gctx.New()
		_, err := Dep.Index(ctx, Input{
			Package:  "./",
			Format:   "tree",
			Depth:    1,
			Internal: true,
			NoStd:    true,
		})
		t.AssertNil(err)
	})
}

func Test_Dep_List(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := gctx.New()
		_, err := Dep.Index(ctx, Input{
			Package:  "./",
			Format:   "list",
			Depth:    1,
			Internal: true,
			NoStd:    true,
		})
		t.AssertNil(err)
	})
}

func Test_Dep_Mermaid(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := gctx.New()
		_, err := Dep.Index(ctx, Input{
			Package:  "./",
			Format:   "mermaid",
			Depth:    1,
			Internal: true,
			NoStd:    true,
		})
		t.AssertNil(err)
	})
}

func Test_Dep_Dot(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := gctx.New()
		_, err := Dep.Index(ctx, Input{
			Package:  "./",
			Format:   "dot",
			Depth:    1,
			Internal: true,
			NoStd:    true,
		})
		t.AssertNil(err)
	})
}

func Test_Dep_JSON(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := gctx.New()
		_, err := Dep.Index(ctx, Input{
			Package:  "./",
			Format:   "json",
			Depth:    1,
			Internal: true,
			NoStd:    true,
		})
		t.AssertNil(err)
	})
}

func Test_Dep_Reverse(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := gctx.New()
		_, err := Dep.Index(ctx, Input{
			Package:  "./",
			Format:   "tree",
			Depth:    1,
			Internal: true,
			NoStd:    true,
			Reverse:  true,
		})
		t.AssertNil(err)
	})
}

func Test_Dep_Group(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := gctx.New()
		_, err := Dep.Index(ctx, Input{
			Package:  "./",
			Format:   "mermaid",
			Depth:    1,
			Internal: true,
			NoStd:    true,
			Group:    true,
		})
		t.AssertNil(err)
	})
}

func Test_ModuleLevel_Direct(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := gctx.New()

		// Test module level with direct only
		in := Input{
			Module: true,
			Direct: true,
			Format: "list",
		}

		t.Logf("Input.Module: %v, Input.Direct: %v", in.Module, in.Direct)

		_, err := Dep.Index(ctx, in)
		t.AssertNil(err)
	})
}
