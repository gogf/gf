// Copyright 2017 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

package ghttp

import (
    "gitee.com/johng/gf/g/os/glog"
)


// 设置http server参数 - CookieMaxAge
func (s *Server)SetCookieMaxAge(age int) {
    if s.Status() == SERVER_STATUS_RUNNING {
        glog.Error(gCHANGE_CONFIG_WHILE_RUNNING_ERROR)
        return
    }
    s.config.CookieMaxAge = age
}

// 设置http server参数 - CookiePath
func (s *Server)SetCookiePath(path string) {
    if s.Status() == SERVER_STATUS_RUNNING {
        glog.Error(gCHANGE_CONFIG_WHILE_RUNNING_ERROR)
        return
    }
    s.config.CookiePath = path
}

// 设置http server参数 - CookieDomain
func (s *Server)SetCookieDomain(domain string) {
    if s.Status() == SERVER_STATUS_RUNNING {
        glog.Error(gCHANGE_CONFIG_WHILE_RUNNING_ERROR)
        return
    }
    s.config.CookieDomain = domain
}

// 获取http server参数 - CookieMaxAge
func (s *Server)GetCookieMaxAge() int {
    return s.config.CookieMaxAge
}

// 获取http server参数 - CookiePath
func (s *Server)GetCookiePath() string {
    return s.config.CookiePath
}

// 获取http server参数 - CookieDomain
func (s *Server)GetCookieDomain() string {
    return s.config.CookieDomain
}
