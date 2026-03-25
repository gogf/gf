// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package genenums

import (
	"go/constant"
	"path/filepath"
	"testing"

	"golang.org/x/tools/go/packages"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/guid"
)

func Test_NewEnumsParser(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// Test creating parser without prefixes
		p := NewEnumsParser(nil)
		t.AssertNE(p, nil)
		t.Assert(len(p.enums), 0)
		t.Assert(len(p.prefixes), 0)
		t.AssertNE(p.parsedPkg, nil)
		t.AssertNE(p.standardPackages, nil)
	})
}

func Test_NewEnumsParser_WithPrefixes(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// Test creating parser with prefixes
		prefixes := []string{"github.com/gogf", "github.com/test"}
		p := NewEnumsParser(prefixes)
		t.AssertNE(p, nil)
		t.Assert(len(p.prefixes), 2)
		t.Assert(p.prefixes[0], "github.com/gogf")
		t.Assert(p.prefixes[1], "github.com/test")
	})
}

func Test_EnumsParser_Export_Empty(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// Test exporting empty enums
		p := NewEnumsParser(nil)
		result := p.Export()
		t.Assert(result, "{}")
	})
}

func Test_EnumsParser_Export_WithEnums(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// Test exporting with manually added enums
		p := NewEnumsParser(nil)

		// Add some test enums
		p.enums = []EnumItem{
			{
				Name:  "StatusActive",
				Value: "1",
				Type:  "pkg.Status",
				Kind:  constant.Int,
			},
			{
				Name:  "StatusInactive",
				Value: "0",
				Type:  "pkg.Status",
				Kind:  constant.Int,
			},
			{
				Name:  "TypeA",
				Value: "type_a",
				Type:  "pkg.Type",
				Kind:  constant.String,
			},
		}

		result := p.Export()
		t.AssertNE(result, "")

		// Parse the result to verify - use raw map to avoid gjson path issues with "."
		var resultMap map[string][]interface{}
		err := gjson.DecodeTo(result, &resultMap)
		t.AssertNil(err)

		// Verify Status type has 2 values
		statusValues := resultMap["pkg.Status"]
		t.Assert(len(statusValues), 2)

		// Verify Type type has 1 value
		typeValues := resultMap["pkg.Type"]
		t.Assert(len(typeValues), 1)
		t.Assert(typeValues[0], "type_a")
	})
}

func Test_EnumsParser_Export_IntValues(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		p := NewEnumsParser(nil)
		p.enums = []EnumItem{
			{Name: "One", Value: "1", Type: "pkg.Int", Kind: constant.Int},
			{Name: "Two", Value: "2", Type: "pkg.Int", Kind: constant.Int},
			{Name: "Negative", Value: "-5", Type: "pkg.Int", Kind: constant.Int},
		}

		result := p.Export()
		var resultMap map[string][]interface{}
		err := gjson.DecodeTo(result, &resultMap)
		t.AssertNil(err)

		values := resultMap["pkg.Int"]
		t.Assert(len(values), 3)
		// Int values should be exported as integers (stored as float64 in JSON)
		t.Assert(values[0], float64(1))
		t.Assert(values[1], float64(2))
		t.Assert(values[2], float64(-5))
	})
}

func Test_EnumsParser_Export_FloatValues(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		p := NewEnumsParser(nil)
		p.enums = []EnumItem{
			{Name: "Pi", Value: "3.14159", Type: "pkg.Float", Kind: constant.Float},
			{Name: "E", Value: "2.71828", Type: "pkg.Float", Kind: constant.Float},
		}

		result := p.Export()
		var resultMap map[string][]interface{}
		err := gjson.DecodeTo(result, &resultMap)
		t.AssertNil(err)

		values := resultMap["pkg.Float"]
		t.Assert(len(values), 2)
	})
}

func Test_EnumsParser_Export_BoolValues(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		p := NewEnumsParser(nil)
		p.enums = []EnumItem{
			{Name: "True", Value: "true", Type: "pkg.Bool", Kind: constant.Bool},
			{Name: "False", Value: "false", Type: "pkg.Bool", Kind: constant.Bool},
		}

		result := p.Export()
		var resultMap map[string][]interface{}
		err := gjson.DecodeTo(result, &resultMap)
		t.AssertNil(err)

		values := resultMap["pkg.Bool"]
		t.Assert(len(values), 2)
		t.Assert(values[0], true)
		t.Assert(values[1], false)
	})
}

func Test_EnumsParser_Export_StringValues(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		p := NewEnumsParser(nil)
		p.enums = []EnumItem{
			{Name: "Hello", Value: "hello", Type: "pkg.Str", Kind: constant.String},
			{Name: "World", Value: "world", Type: "pkg.Str", Kind: constant.String},
		}

		result := p.Export()
		var resultMap map[string][]interface{}
		err := gjson.DecodeTo(result, &resultMap)
		t.AssertNil(err)

		values := resultMap["pkg.Str"]
		t.Assert(len(values), 2)
		t.Assert(values[0], "hello")
		t.Assert(values[1], "world")
	})
}

func Test_EnumsParser_Export_MixedTypes(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		p := NewEnumsParser(nil)
		p.enums = []EnumItem{
			{Name: "IntVal", Value: "42", Type: "pkg.IntType", Kind: constant.Int},
			{Name: "StrVal", Value: "test", Type: "pkg.StrType", Kind: constant.String},
			{Name: "BoolVal", Value: "true", Type: "pkg.BoolType", Kind: constant.Bool},
		}

		result := p.Export()
		var resultMap map[string][]interface{}
		err := gjson.DecodeTo(result, &resultMap)
		t.AssertNil(err)

		// Each type should have its own array
		t.Assert(len(resultMap["pkg.IntType"]), 1)
		t.Assert(len(resultMap["pkg.StrType"]), 1)
		t.Assert(len(resultMap["pkg.BoolType"]), 1)
	})
}

func Test_EnumItem_Structure(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// Test EnumItem structure
		item := EnumItem{
			Name:  "TestEnum",
			Value: "test_value",
			Type:  "github.com/test/pkg.EnumType",
			Kind:  constant.String,
		}

		t.Assert(item.Name, "TestEnum")
		t.Assert(item.Value, "test_value")
		t.Assert(item.Type, "github.com/test/pkg.EnumType")
		t.Assert(item.Kind, constant.String)
	})
}

func Test_EnumsParser_ParsePackages_Integration(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// Create a temporary directory with a Go package containing enums
		// Note: The module path must contain "/" for enums to be parsed
		// (the parser skips std types without "/" in the type name)
		tempDir := gfile.Temp(guid.S())
		err := gfile.Mkdir(tempDir)
		t.AssertNil(err)
		defer gfile.Remove(tempDir)

		// Create go.mod with a path containing "/"
		goModContent := `module github.com/test/enumtest

go 1.21
`
		err = gfile.PutContents(filepath.Join(tempDir, "go.mod"), goModContent)
		t.AssertNil(err)

		// Create a Go file with enum definitions
		enumsContent := `package enumtest

type Status int

const (
	StatusActive   Status = 1
	StatusInactive Status = 0
)

type Color string

const (
	ColorRed   Color = "red"
	ColorGreen Color = "green"
	ColorBlue  Color = "blue"
)
`
		err = gfile.PutContents(filepath.Join(tempDir, "enums.go"), enumsContent)
		t.AssertNil(err)

		// Load the package
		cfg := &packages.Config{
			Dir:   tempDir,
			Mode:  pkgLoadMode,
			Tests: false,
		}
		pkgs, err := packages.Load(cfg)
		t.AssertNil(err)
		t.Assert(len(pkgs) > 0, true)

		// Parse the packages
		p := NewEnumsParser(nil)
		p.ParsePackages(pkgs)

		// Export and verify - result should contain parsed enums
		result := p.Export()
		// Verify the export contains some data
		t.Assert(len(result) > 2, true) // More than just "{}"

		// Parse result as raw map to handle keys with "/"
		var resultMap map[string][]interface{}
		err = gjson.DecodeTo(result, &resultMap)
		t.AssertNil(err)

		// Verify Status enum was parsed (type will be "github.com/test/enumtest.Status")
		statusKey := "github.com/test/enumtest.Status"
		statusValues, hasStatus := resultMap[statusKey]
		t.Assert(hasStatus, true)
		t.Assert(len(statusValues), 2)

		// Verify Color enum was parsed
		colorKey := "github.com/test/enumtest.Color"
		colorValues, hasColor := resultMap[colorKey]
		t.Assert(hasColor, true)
		t.Assert(len(colorValues), 3)
	})
}

func Test_EnumsParser_ParsePackages_WithPrefixes(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// Create a temporary directory with a Go package
		tempDir := gfile.Temp(guid.S())
		err := gfile.Mkdir(tempDir)
		t.AssertNil(err)
		defer gfile.Remove(tempDir)

		// Create go.mod with a specific module name
		goModContent := `module github.com/allowed/pkg

go 1.21
`
		err = gfile.PutContents(filepath.Join(tempDir, "go.mod"), goModContent)
		t.AssertNil(err)

		// Create a Go file with enum definitions
		enumsContent := `package pkg

type Status int

const (
	StatusOK Status = 1
)
`
		err = gfile.PutContents(filepath.Join(tempDir, "enums.go"), enumsContent)
		t.AssertNil(err)

		// Load the package
		cfg := &packages.Config{
			Dir:   tempDir,
			Mode:  pkgLoadMode,
			Tests: false,
		}
		pkgs, err := packages.Load(cfg)
		t.AssertNil(err)

		// Parse with prefix filter that matches
		p := NewEnumsParser([]string{"github.com/allowed"})
		p.ParsePackages(pkgs)

		result := p.Export()
		// Should have enums because prefix matches
		t.AssertNE(result, "{}")

		// Parse with prefix filter that doesn't match
		p2 := NewEnumsParser([]string{"github.com/other"})
		p2.ParsePackages(pkgs)

		result2 := p2.Export()
		// Should be empty because prefix doesn't match
		t.Assert(result2, "{}")
	})
}

func Test_getStandardPackages(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		stdPkgs := getStandardPackages()
		t.AssertNE(stdPkgs, nil)
		t.Assert(len(stdPkgs) > 0, true)

		// Verify some common standard packages are included
		_, hasFmt := stdPkgs["fmt"]
		t.Assert(hasFmt, true)

		_, hasOs := stdPkgs["os"]
		t.Assert(hasOs, true)

		_, hasContext := stdPkgs["context"]
		t.Assert(hasContext, true)
	})
}
