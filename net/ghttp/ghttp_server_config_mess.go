// Copyright 2019 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package ghttp

func (s *Server) SetGzipContentTypes(types []string) {
	s.config.GzipContentTypes = types
}

// 服务注册时对象和方法名称转换为URI时的规则
func (s *Server) SetNameToUriType(t int) {
	s.config.NameToUriType = t
}

// 是否在程序启动时打印路由表信息
func (s *Server) SetDumpRouteMap(enabled bool) {
	s.config.DumpRouteMap = enabled
}

// 设置路由缓存过期时间(秒)
func (s *Server) SetRouterCacheExpire(expire int) {
	s.config.RouterCacheExpire = expire
}

func (s *Server) SetFormParsingMemory(maxMemory int64) {
	s.config.FormParsingMemory = maxMemory
}
