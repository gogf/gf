// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package goai

import (
	"github.com/gogf/gf/v2/util/gtag"
)

type EnumXExtensionInput struct {
	TypeID string          // TypeID is the full type id, eg: github.com/gogf/gf/v2/net/goai_test.Status.
	Items  []gtag.EnumItem // Items are all enum values and comments of this type.
}

// EnumXExtensionFunc builds x-extension values for the whole enum type.
// Returned map key should be extension key like "x-apifox-enum".
type EnumXExtensionFunc func(in EnumXExtensionInput) map[string]any

// Config provides extra configuration feature for OpenApiV3 implements.
type Config struct {
	ReadContentTypes        []string // ReadContentTypes specifies the default MIME types for consuming if MIME types are not configured.
	WriteContentTypes       []string // WriteContentTypes specifies the default MIME types for producing if MIME types are not configured.
	CommonRequest           any      // Common request structure for all paths.
	CommonRequestDataField  string   // Common request field name to be replaced with certain business request structure. Eg: `Data`, `Request.`.
	CommonResponse          any      // Common response structure for all paths.
	CommonResponseDataField string   // Common response field name to be replaced with certain business response structure. Eg: `Data`, `Response.`.
	IgnorePkgPath           bool     // Ignores package name for schema name.
	// TypeMapping customizes OpenAPI type mapping for given golang types.
	// Map key supports both short type name and full type id:
	// 1. `carbon.Carbon`
	// 2. `github.com/golang-module/carbon/v2.Carbon`
	TypeMapping map[string]string
	// EnumXExtensionFunc is called for each enum type.
	// Returned extension values are written to schema x-extensions.
	EnumXExtensionFunc EnumXExtensionFunc
}

// fillWithDefaultValue fills configuration object of `oai` with default values if these are not configured.
func (oai *OpenApiV3) fillWithDefaultValue() {
	if oai.OpenAPI == "" {
		oai.OpenAPI = `3.0.0`
	}
	if len(oai.Config.ReadContentTypes) == 0 {
		oai.Config.ReadContentTypes = defaultReadContentTypes
	}
	if len(oai.Config.WriteContentTypes) == 0 {
		oai.Config.WriteContentTypes = defaultWriteContentTypes
	}
}

// AddEnumXExtensionFunc adds enum x-extension callback.
func (oai *OpenApiV3) AddEnumXExtensionFunc(fn EnumXExtensionFunc) {
	oai.Config.EnumXExtensionFunc = fn
}
