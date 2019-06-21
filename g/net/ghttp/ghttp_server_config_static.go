// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// 静态文件搜索优先级: ServerPaths > ServerRoot > SearchPath

package ghttp

import (
	"fmt"
	"github.com/gogf/gf/g/container/garray"
	"github.com/gogf/gf/g/os/gfile"
	"github.com/gogf/gf/g/os/glog"
	"github.com/gogf/gf/g/util/gconv"
	"strings"
)

// 静态文件目录映射关系对象
type staticPathItem struct {
	prefix string // 映射的URI前缀
	path   string // 静态文件目录绝对路径
}

// 设置http server参数 - IndexFiles，默认展示文件，如：index.html, index.htm
func (s *Server) SetIndexFiles(index []string) {
	if s.Status() == SERVER_STATUS_RUNNING {
		glog.Error(gCHANGE_CONFIG_WHILE_RUNNING_ERROR)
		return
	}
	s.config.IndexFiles = index
}

// 允许展示访问目录的文件列表
func (s *Server) SetIndexFolder(enabled bool) {
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
func (s *Server) SetServerRoot(root string) {
	if s.Status() == SERVER_STATUS_RUNNING {
		glog.Error(gCHANGE_CONFIG_WHILE_RUNNING_ERROR)
		return
	}
	// RealPath的作用除了校验地址正确性以外，还转换分隔符号为当前系统正确的文件分隔符号
	realPath, err := gfile.Search(root)
	if err != nil {
		glog.Fatal(fmt.Sprintf(`[ghttp] SetServerRoot failed: %s`, err.Error()))
	}
	glog.Debug("[ghttp] SetServerRoot path:", realPath)
	s.config.SearchPaths = []string{strings.TrimRight(realPath, gfile.Separator)}
	s.config.FileServerEnabled = true
}

// 添加静态文件搜索**目录**，必须给定目录的绝对路径
func (s *Server) AddSearchPath(path string) {
	if s.Status() == SERVER_STATUS_RUNNING {
		glog.Error(gCHANGE_CONFIG_WHILE_RUNNING_ERROR)
		return
	}
	realPath, err := gfile.Search(path)
	if err != nil {
		glog.Fatal(fmt.Sprintf(`[ghttp] AddSearchPath failed: %s`, err.Error()))
	}
	s.config.SearchPaths = append(s.config.SearchPaths, realPath)
	s.config.FileServerEnabled = true
}

// 添加URI与静态**目录**的映射
func (s *Server) AddStaticPath(prefix string, path string) {
	if s.Status() == SERVER_STATUS_RUNNING {
		glog.Error(gCHANGE_CONFIG_WHILE_RUNNING_ERROR)
		return
	}
	realPath, err := gfile.Search(path)
	if err != nil {
		glog.Fatal(fmt.Sprintf(`[ghttp] AddStaticPath failed: %s`, err.Error()))
	}
	addItem := staticPathItem{
		prefix: prefix,
		path:   realPath,
	}
	if len(s.config.StaticPaths) > 0 {
		// 先添加item
		s.config.StaticPaths = append(s.config.StaticPaths, addItem)
		// 按照prefix从长到短进行排序
		array := garray.NewSortedArray(func(v1, v2 interface{}) int {
			s1 := gconv.String(v1)
			s2 := gconv.String(v2)
			r := len(s2) - len(s1)
			if r == 0 {
				r = strings.Compare(s1, s2)
			}
			return r
		}, false)
		for _, v := range s.config.StaticPaths {
			array.Add(v.prefix)
		}
		// 按照重新排序的顺序重新添加item
		paths := make([]staticPathItem, 0)
		for _, v := range array.Slice() {
			for _, item := range s.config.StaticPaths {
				if strings.EqualFold(gconv.String(v), item.prefix) {
					paths = append(paths, item)
					break
				}
			}
		}
		s.config.StaticPaths = paths
	} else {
		s.config.StaticPaths = []staticPathItem{addItem}
	}
	s.config.FileServerEnabled = true
}
