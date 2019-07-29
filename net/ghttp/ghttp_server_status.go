// Copyright 2018 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.
// 状态码回调函数注册.

package ghttp

import (
	"fmt"
)

// 查询状态码回调函数
func (s *Server) getStatusHandler(status int, r *Request) HandlerFunc {
	domains := []string{r.GetHost(), gDEFAULT_DOMAIN}
	s.statusHandlerMu.RLock()
	defer s.statusHandlerMu.RUnlock()
	for _, domain := range domains {
		if f, ok := s.statusHandlerMap[s.statusHandlerKey(status, domain)]; ok {
			return f
		}
	}
	return nil
}

// 不同状态码下的回调方法处理
// pattern格式：domain#status
func (s *Server) setStatusHandler(pattern string, handler HandlerFunc) {
	s.statusHandlerMu.Lock()
	s.statusHandlerMap[pattern] = handler
	s.statusHandlerMu.Unlock()
}

// 生成状态码回调函数map存储键名
func (s *Server) statusHandlerKey(status int, domain string) string {
	return fmt.Sprintf("%s#%d", domain, status)
}

// 绑定指定的状态码回调函数
func (s *Server) BindStatusHandler(status int, handler HandlerFunc) {
	s.setStatusHandler(s.statusHandlerKey(status, gDEFAULT_DOMAIN), handler)
}

// 通过map批量绑定状态码回调函数
func (s *Server) BindStatusHandlerByMap(handlerMap map[int]HandlerFunc) {
	for k, v := range handlerMap {
		s.BindStatusHandler(k, v)
	}
}
