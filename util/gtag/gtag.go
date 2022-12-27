// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package gtag providing tag content storing for struct.
//
// Note that calling functions of this package is not concurrently safe,
// which means you cannot call them in runtime but in boot procedure.
package gtag

const (
	Default           = "default"      // Default value tag of struct field for receiving parameters from HTTP request.
	DefaultShort      = "d"            // Short name of Default.
	Param             = "param"        // Parameter name for converting certain parameter to specified struct field.
	ParamShort        = "p"            // Short name of Param.
	Valid             = "valid"        // Validation rule tag for struct of field.
	ValidShort        = "v"            // Short name of Valid.
	NoValidation      = "nv"           // No validation for specified struct/field.
	ORM               = "orm"          // ORM tag for ORM feature, which performs different features according scenarios.
	Arg               = "arg"          // Arg tag for struct, usually for command argument option.
	Brief             = "brief"        // Brief tag for struct, usually be considered as summary.
	Root              = "root"         // Root tag for struct, usually for nested commands management.
	Additional        = "additional"   // Additional tag for struct, usually for additional description of command.
	AdditionalShort   = "ad"           // Short name of Additional.
	Path              = `path`         // Route path for HTTP request.
	Method            = `method`       // Route method for HTTP request.
	Domain            = `domain`       // Route domain for HTTP request.
	Mime              = `mime`         // MIME type for HTTP request/response.
	Consumes          = `consumes`     // MIME type for HTTP request.
	Summary           = `summary`      // Summary for struct, usually for OpenAPI in request struct.
	SummaryShort      = `sm`           // Short name of Summary.
	SummaryShort2     = `sum`          // Short name of Summary.
	Description       = `description`  // Description for struct, usually for OpenAPI in request struct.
	DescriptionShort  = `dc`           // Short name of Description.
	DescriptionShort2 = `des`          // Short name of Description.
	Example           = `example`      // Example for struct, usually for OpenAPI in request struct.
	ExampleShort      = `eg`           // Short name of Example.
	Examples          = `examples`     // Examples for struct, usually for OpenAPI in request struct.
	ExamplesShort     = `egs`          // Short name of Examples.
	ExternalDocs      = `externalDocs` // External docs for struct, always for OpenAPI in request struct.
	ExternalDocsShort = `ed`           // Short name of ExternalDocs.
	GConv             = "gconv"        // GConv defines the converting target name for specified struct field.
	GConvShort        = "c"            // GConv defines the converting target name for specified struct field.
	Json              = "json"         // Json tag is supported by stdlib.
	Security          = "security"     // security schema.
)
