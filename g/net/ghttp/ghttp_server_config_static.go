// Copyright 2017 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

// 静态文件搜索优先级: ServerPaths > ServerRoot > SearchPath

package ghttp

import (
    "gitee.com/johng/gf/g/os/gfile"
    "gitee.com/johng/gf/g/os/glog"
    "strings"
)

// 设置http server参数 - IndexFiles，默认展示文件，如：index.html, index.htm
func (s *Server)SetIndexFiles(index []string) {
    if s.Status() == SERVER_STATUS_RUNNING {
        glog.Error(gCHANGE_CONFIG_WHILE_RUNNING_ERROR)
        return
    }
    s.config.IndexFiles = index
}

// 允许展示访问目录的文件列表
func (s *Server)SetIndexFolder(enabled bool) {
    if s.Status() == SERVER_STATUS_RUNNING {
        glog.Error(gCHANGE_CONFIG_WHILE_RUNNING_ERROR)
        return
    }
    s.config.IndexFolder = enabled
}

// 是否开启/关闭静态文件服务，当关闭时仅提供动态接口服务，路由性能会得到一定提升
func (s *Server) SetFileServerEnabled(enabled bool) {
    if s.Status() == SERVER_STATUS_RUNNING {
        glog.Error(gCHANGE_CONFIG_WHILE_RUNNING_ERROR)
        return
    }
    s.config.FileServerEnabled = enabled
}

// 设置http server参数 - ServerRoot
func (s *Server)SetServerRoot(root string) {
    if s.Status() == SERVER_STATUS_RUNNING {
        glog.Error(gCHANGE_CONFIG_WHILE_RUNNING_ERROR)
        return
    }
    // RealPath的作用除了校验地址正确性以外，还转换分隔符号为当前系统正确的文件分隔符号
    path := gfile.RealPath(root)
    if path == "" {
        glog.Error("invalid root path \"" + root + "\"")
    }
    s.config.ServerRoot = strings.TrimRight(path, gfile.Separator)
}

// 添加静态文件搜索目录，必须给定目录的绝对路径
func (s *Server) AddSearchPath(path string) error {
    if rp, err := s.paths.Add(path); err != nil {
        glog.Error("[ghttp] AddSearchPath failed:", err.Error())
        return err
    } else {
        glog.Debug("[ghttp] AddSearchPath:", rp)
    }
    return nil
}

// 添加URI与静态目录的映射
func (s *Server) AddStaticPath(prefix string, path string) error {
    if rp, err := s.paths.Add(path); err != nil {
        glog.Error("[ghttp] AddSearchPath failed:", err.Error())
        return err
    } else {
        glog.Debug("[ghttp] AddSearchPath:", rp)
    }
    return nil
}

