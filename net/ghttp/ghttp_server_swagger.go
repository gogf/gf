// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package ghttp

import (
	"net/http"

	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/text/gstr"
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
	var (
		indexFileName = `index.html`
	)
	if r.StaticFile != nil && r.StaticFile.File != nil && gfile.Basename(r.StaticFile.File.Name()) == indexFileName {
		if gfile.Basename(r.URL.Path) != indexFileName && r.originUrlPath[len(r.originUrlPath)-1] != '/' {
			r.Response.Header().Set("Location", r.originUrlPath+"/")
			r.Response.WriteHeader(http.StatusMovedPermanently)
			r.ExitAll()
		}
		r.Response.Write(gstr.Replace(
			string(r.StaticFile.File.Content()),
			swaggerUIDefaultURL,
			s.config.OpenApiPath,
		))
		r.ExitAll()
	}
}
