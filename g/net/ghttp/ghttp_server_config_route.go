// Copyright 2018 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.


package ghttp

import "gitee.com/johng/gf/g/os/glog"

func (s *Server) SetDenyIps(ips []string) {
    if s.Status() == SERVER_STATUS_RUNNING {
        glog.Error(gCHANGE_CONFIG_WHILE_RUNNING_ERROR)
        return
    }
    s.config.DenyIps = ips
}

func (s *Server) SetAllowIps(ips []string) {
    if s.Status() == SERVER_STATUS_RUNNING {
        glog.Error(gCHANGE_CONFIG_WHILE_RUNNING_ERROR)
        return
    }
    s.config.AllowIps = ips
}

func (s *Server) SetDenyRoutes(routes []string) {
    if s.Status() == SERVER_STATUS_RUNNING {
        glog.Error(gCHANGE_CONFIG_WHILE_RUNNING_ERROR)
        return
    }
    s.config.DenyRoutes = routes
}

// 设置URI重写规则
func (s *Server) SetRewrite(uri string, rewrite string) {
    if s.Status() == SERVER_STATUS_RUNNING {
        glog.Error(gCHANGE_CONFIG_WHILE_RUNNING_ERROR)
        return
    }
    s.config.Rewrites[uri] = rewrite
}

// 设置URI重写规则（批量）
func (s *Server) SetRewriteMap(rewrites map[string]string) {
    if s.Status() == SERVER_STATUS_RUNNING {
        glog.Error(gCHANGE_CONFIG_WHILE_RUNNING_ERROR)
        return
    }
    for k, v := range rewrites {
        s.config.Rewrites[k] = v
    }
}