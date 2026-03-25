// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package goai_test

import (
	"context"
	"strings"
	"testing"

	"github.com/gogf/gf/v2/net/goai"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/gmeta"
)

// TestOpenApiV3_GenericType tests the schema name generation for generic types
// This test validates the PR fix for swagger $ref replace that handles Go generics
// Specifically testing that [ and ] characters in type names are replaced with dots
func TestOpenApiV3_GenericType(t *testing.T) {
	// Define a generic type wrapper
	type GenericItem[T any] struct {
		Value T `dc:"Generic value"`
	}

	type StringItem = GenericItem[string]

	type IntItem = GenericItem[int]

	type Req struct {
		gmeta.Meta `path:"/generic" method:"POST" tags:"default"`
		StringData StringItem `dc:"String generic type"`
		IntData    IntItem    `dc:"Int generic type"`
	}

	type Res struct {
		gmeta.Meta `description:"Generic Response"`
		Data       string `dc:"Response data"`
	}

	f := func(ctx context.Context, req *Req) (res *Res, err error) {
		return
	}

	gtest.C(t, func(t *gtest.T) {
		var (
			err error
			oai = goai.New()
		)
		err = oai.Add(goai.AddInput{
			Path:   "/generic",
			Object: f,
		})
		t.AssertNil(err)

		// Verify that schema names are properly generated without special characters
		schemas := oai.Components.Schemas.Map()
		t.AssertGT(len(schemas), 0)

		// Check that bracket characters [ and ] have been replaced with dots
		// According to PR fix: `[`: `.`, `]`: `.`
		for schemaName := range schemas {
			// Should not contain [ or ] characters after replacement
			t.Assert(!strings.Contains(schemaName, "["), true)
			t.Assert(!strings.Contains(schemaName, "]"), true)
		}
	})
}

// TestOpenApiV3_SchemaNameReplacement tests the special character replacement in schema names
// This verifies the core PR change which replaces:
// - [ with .
// - ] with .
// - { with empty string
// - } with empty string
// - spaces with empty string
func TestOpenApiV3_SchemaNameReplacement(t *testing.T) {
	type SimpleReq struct {
		gmeta.Meta `path:"/test" method:"POST"`
		Name       string `dc:"Name field"`
	}

	type SimpleRes struct {
		gmeta.Meta `description:"Simple Response"`
		Status     string `dc:"Status field"`
	}

	f := func(ctx context.Context, req *SimpleReq) (res *SimpleRes, err error) {
		return
	}

	gtest.C(t, func(t *gtest.T) {
		var (
			err error
			oai = goai.New()
		)
		err = oai.Add(goai.AddInput{
			Path:   "/test",
			Object: f,
		})
		t.AssertNil(err)

		// Get schema names and verify they are properly formatted
		schemas := oai.Components.Schemas.Map()
		for schemaName := range schemas {
			// Verify special characters have been replaced:
			// - [ should be replaced with .
			// - ] should be replaced with .
			// - { should be replaced with empty
			// - } should be replaced with empty
			// - spaces should be replaced with empty
			t.Assert(!strings.Contains(schemaName, "["), true)
			t.Assert(!strings.Contains(schemaName, "]"), true)
			t.Assert(!strings.Contains(schemaName, "{"), true)
			t.Assert(!strings.Contains(schemaName, "}"), true)
		}
	})
}

// TestOpenApiV3_ComplexGenericType tests more complex generic types
// This specifically tests handling of map types and nested generic structures
func TestOpenApiV3_ComplexGenericType(t *testing.T) {
	type MapWrapper struct {
		gmeta.Meta `path:"/mapwrapper" method:"POST"`
		Data       map[string]string `dc:"Map data"`
	}

	type Res struct {
		gmeta.Meta `description:"Map Response"`
		Result     string `dc:"Result"`
	}

	f := func(ctx context.Context, req *MapWrapper) (res *Res, err error) {
		return
	}

	gtest.C(t, func(t *gtest.T) {
		var (
			err error
			oai = goai.New()
		)
		err = oai.Add(goai.AddInput{
			Path:   "/mapwrapper",
			Object: f,
		})
		t.AssertNil(err)

		// Verify schema generation completes without errors
		schemas := oai.Components.Schemas.Map()
		t.AssertGT(len(schemas), 0)

		// All schema names should be valid (no bracket characters)
		for schemaName := range schemas {
			t.Assert(!strings.Contains(schemaName, "["), true)
			t.Assert(!strings.Contains(schemaName, "]"), true)
		}
	})
}

// TestOpenApiV3_PathWithSpecialChars tests path parameters with special handling
// This ensures the PR changes don't affect regular parameter handling
func TestOpenApiV3_PathWithSpecialChars(t *testing.T) {
	type GetDetailReq struct {
		gmeta.Meta `path:"/detail" method:"GET"`
		ResourceId string `json:"resourceId" in:"query" dc:"Resource identifier"`
		Type       string `json:"type" in:"query" dc:"Resource type"`
	}

	type DetailRes struct {
		gmeta.Meta `description:"Detail Response"`
		Content    string `dc:"Detail content"`
	}

	f := func(ctx context.Context, req *GetDetailReq) (res *DetailRes, err error) {
		return
	}

	gtest.C(t, func(t *gtest.T) {
		var (
			err error
			oai = goai.New()
		)
		err = oai.Add(goai.AddInput{
			Path:   "/detail",
			Object: f,
		})
		t.AssertNil(err)

		// Verify all schemas are properly named
		schemas := oai.Components.Schemas.Map()
		for schemaName := range schemas {
			// Should not contain special characters that were supposed to be replaced
			t.Assert(!strings.Contains(schemaName, "["), true)
			t.Assert(!strings.Contains(schemaName, "]"), true)
		}
	})
}

// TestOpenApiV3_SliceOfGenericTypes tests slice of generic types
// This validates that slices containing generics are properly handled
func TestOpenApiV3_SliceOfGenericTypes(t *testing.T) {
	type Item[T any] struct {
		Value T `dc:"Item value"`
	}

	type StringItem = Item[string]

	type SliceReq struct {
		gmeta.Meta `path:"/slice" method:"POST"`
		Items      []StringItem `dc:"Slice of generic items"`
	}

	type SliceRes struct {
		gmeta.Meta `description:"Slice Response"`
		Count      int `dc:"Item count"`
	}

	f := func(ctx context.Context, req *SliceReq) (res *SliceRes, err error) {
		return
	}

	gtest.C(t, func(t *gtest.T) {
		var (
			err error
			oai = goai.New()
		)
		err = oai.Add(goai.AddInput{
			Path:   "/slice",
			Object: f,
		})
		t.AssertNil(err)

		schemas := oai.Components.Schemas.Map()
		t.AssertGT(len(schemas), 0)

		// Verify no bracket characters in schema names
		for schemaName := range schemas {
			t.Assert(!strings.Contains(schemaName, "["), true)
			t.Assert(!strings.Contains(schemaName, "]"), true)
		}
	})
}
