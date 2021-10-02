// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package ghttp

import (
	"github.com/gogf/gf/os/gfile"
	"github.com/gogf/gf/text/gstr"
)

const (
	swaggerUIDefaultURL = `https://petstore.swagger.io/v2/swagger.json`
)

// swaggerUI is a build-in hook handler for replace default swagger json URL to local openapi json file path.
// This handler makes sense only if the openapi specification automatic producing configuration is enabled.
func (s *Server) swaggerUI(r *Request) {
	if s.config.OpenApiPath == "" {
		return
	}
	if r.StaticFile != nil && r.StaticFile.File != nil && gfile.Basename(r.StaticFile.File.Name()) == "index.html" {
		r.Response.Write(gstr.Replace(
			string(r.StaticFile.File.Content()),
			swaggerUIDefaultURL,
			s.config.OpenApiPath,
		))
		r.ExitAll()
	}
}
