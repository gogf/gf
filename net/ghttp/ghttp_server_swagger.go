// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package ghttp

import (
	"fmt"

	"github.com/gogf/gf/v2/text/gstr"
)

const (
	swaggerUIDocName            = `redoc.standalone.js`
	swaggerUIDocNamePlaceHolder = `{SwaggerUIDocName}`
	swaggerUIDocURLPlaceHolder  = `{SwaggerUIDocUrl}`
	swaggerUIDocJsPlaceHolder   = `{SwaggerUIDocJs}`
	swaggerUITemplate           = `
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
		<script src="{SwaggerUIDocJs}"> </script>
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

	if s.config.SwaggerJsURL == "" {
		s.config.SwaggerJsURL = "https://unpkg.com/redoc@2.0.0-rc.70/bundles/redoc.standalone.js"
	}

	if r.StaticFile != nil && r.StaticFile.File != nil && r.StaticFile.IsDir {
		content := gstr.ReplaceByMap(swaggerUITemplate, map[string]string{
			swaggerUIDocURLPlaceHolder:  s.config.OpenApiPath,
			swaggerUIDocNamePlaceHolder: gstr.TrimRight(fmt.Sprintf(`//%s%s`, r.Host, r.Server.config.SwaggerPath), "/") + "/" + swaggerUIDocName,
			swaggerUIDocJsPlaceHolder:   s.config.SwaggerJsURL,
		})
		r.Response.Write(content)
		r.ExitAll()
	}
}
