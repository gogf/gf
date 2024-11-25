// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package ghttp

// SetSwaggerPath sets the SwaggerPath for server.
func (s *Server) SetSwaggerPath(path string) {
	s.config.SwaggerPath = path
}

// SetSwaggerUITemplate sets the Swagger template for server.
func (s *Server) SetSwaggerUITemplate(swaggerUITemplate string) {
	s.config.SwaggerUITemplate = swaggerUITemplate
}

// SetOpenApiPath sets the OpenApiPath for server.
// For example: SetOpenApiPath("/api.json")
func (s *Server) SetOpenApiPath(path string) {
	s.config.OpenApiPath = path
}

// SetOpenApiAuthUser sets the OpenApiAuthUser for server.
// For example: SetOpenApiAuthUser("gf")
func (s *Server) SetOpenApiAuthUser(user string) {
	s.config.OpenApiAuthUser = user
}

// SetOpenApiAuthPass sets the OpenApiAuthPass for server.
// For example: SetOpenApiAuthPass("123456")
func (s *Server) SetOpenApiAuthPass(pass string) {
	s.config.OpenApiAuthPass = pass
}
