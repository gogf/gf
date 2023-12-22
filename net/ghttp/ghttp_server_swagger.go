// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package ghttp

import (
	"github.com/gogf/gf/v2/text/gstr"
)

const (
	swaggerUIDocURLPlaceHolder = `{SwaggerUIDocUrl}`
	swaggerUITemplate          = `
<!DOCTYPE html>
<html>
	<head>
	<title>API Reference</title>
	<meta charset="utf-8"/>
	<meta name="viewport" content="width=device-width, initial-scale=1">
	<style>
		body {
			margin:  0;
			padding: 0;
		}
	</style>
	</head>
	<body>
		<redoc spec-url="{SwaggerUIDocUrl}" show-object-schema-examples="true"></redoc>
		<script src="https://cdn.redoc.ly/redoc/latest/bundles/redoc.standalone.js"> </script>
	</body>
</html>
`
)

// swaggerUI is a build-in hook handler for replace default swagger json URL to local openapi json file path.
// This handler makes sense only if the openapi specification automatic producing configuration is enabled.
func (s *Server) swaggerUI(r *Request) {
	if s.config.OpenApiPath == "" {
		return
	}
	var templateContent = swaggerUITemplate
	if s.config.SwaggerUITemplate != "" {
		templateContent = s.config.SwaggerUITemplate
	}

	if r.StaticFile != nil && r.StaticFile.File != nil && r.StaticFile.IsDir {
		content := gstr.ReplaceByMap(templateContent, map[string]string{
			swaggerUIDocURLPlaceHolder: s.config.OpenApiPath,
		})
		r.Response.Write(content)
		r.ExitAll()
	}
}
