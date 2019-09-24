// Copyright 2019 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package ghttp

import "github.com/gogf/gf/os/glog"

func (s *Server) SetGzipContentTypes(types []string) {
	if s.Status() == SERVER_STATUS_RUNNING {
		glog.Error(gCHANGE_CONFIG_WHILE_RUNNING_ERROR)
		return
	}
	s.config.GzipContentTypes = types
}

// 服务注册时对象和方法名称转换为URI时的规则
func (s *Server) SetNameToUriType(t int) {
	if s.Status() == SERVER_STATUS_RUNNING {
		glog.Error(gCHANGE_CONFIG_WHILE_RUNNING_ERROR)
		return
	}
	s.config.NameToUriType = t
}

// 是否在程序启动时打印路由表信息
func (s *Server) SetDumpRouteMap(enabled bool) {
	if s.Status() == SERVER_STATUS_RUNNING {
		glog.Error(gCHANGE_CONFIG_WHILE_RUNNING_ERROR)
		return
	}
	s.config.DumpRouteMap = enabled
}

// 设置路由缓存过期时间(秒)
func (s *Server) SetRouterCacheExpire(expire int) {
	if s.Status() == SERVER_STATUS_RUNNING {
		glog.Error(gCHANGE_CONFIG_WHILE_RUNNING_ERROR)
		return
	}
	s.config.RouterCacheExpire = expire
}

func (s *Server) SetFormParsingMemory(maxMemory int64) {
	if s.Status() == SERVER_STATUS_RUNNING {
		glog.Error(gCHANGE_CONFIG_WHILE_RUNNING_ERROR)
		return
	}
	s.config.FormParsingMemory = maxMemory
}
