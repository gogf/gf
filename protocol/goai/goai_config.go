// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package goai

// Config provides extra configuration feature for OpenApiV3 implements.
type Config struct {
	ReadContentTypes        []string    // ReadContentTypes specifies the default MIME types for consuming if MIME types are not configured.
	WriteContentTypes       []string    // WriteContentTypes specifies the default MIME types for producing if MIME types are not configured.
	CommonRequest           interface{} // Common request structure for all paths.
	CommonRequestDataField  string      // Common request field name to be replaced with certain business request structure. Eg: `Data`, `Request.`.
	CommonResponse          interface{} // Common response structure for all paths.
	CommonResponseDataField string      // Common response field name to be replaced with certain business response structure. Eg: `Data`, `Response.`.
	IgnorePkgPath           bool        // Ignores package name for schema name.
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
