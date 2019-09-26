// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// 静态文件搜索优先级: Resource > ServerPaths > ServerRoot > SearchPath

package ghttp

import (
	"fmt"
	"strings"

	"github.com/gogf/gf/os/gres"

	"github.com/gogf/gf/container/garray"
	"github.com/gogf/gf/os/gfile"
	"github.com/gogf/gf/os/glog"
	"github.com/gogf/gf/util/gconv"
)

// 静态文件目录映射关系对象
type staticPathItem struct {
	prefix string // 映射的URI前缀
	path   string // 静态文件目录绝对路径
}

// 设置http server参数 - IndexFiles，默认展示文件，如：index.html, index.htm
func (s *Server) SetIndexFiles(index []string) {
	s.config.IndexFiles = index
}

// 允许展示访问目录的文件列表
func (s *Server) SetIndexFolder(enabled bool) {
	s.config.IndexFolder = enabled
}

// 是否开启/关闭静态文件服务，当关闭时仅提供动态接口服务，路由性能会得到一定提升
func (s *Server) SetFileServerEnabled(enabled bool) {
	s.config.FileServerEnabled = enabled
}

// 设置http server参数 - ServerRoot
func (s *Server) SetServerRoot(root string) {
	realPath := root
	if !gres.Contains(realPath) {
		if p, err := gfile.Search(root); err != nil {
			glog.Fatal(fmt.Sprintf(`[ghttp] SetServerRoot failed: %s`, err.Error()))
		} else {
			realPath = p
		}
	}
	glog.Debug("[ghttp] SetServerRoot path:", realPath)
	s.config.SearchPaths = []string{strings.TrimRight(realPath, gfile.Separator)}
	s.config.FileServerEnabled = true
}

// 添加静态文件搜索**目录**，必须给定目录的绝对路径
func (s *Server) AddSearchPath(path string) {
	realPath := path
	if !gres.Contains(realPath) {
		if p, err := gfile.Search(path); err != nil {
			glog.Fatal(fmt.Sprintf(`[ghttp] AddSearchPath failed: %s`, err.Error()))
		} else {
			realPath = p
		}
	}
	s.config.SearchPaths = append(s.config.SearchPaths, realPath)
	s.config.FileServerEnabled = true
}

// 添加URI与静态**目录**的映射
func (s *Server) AddStaticPath(prefix string, path string) {
	realPath := path
	if !gres.Contains(realPath) {
		if p, err := gfile.Search(path); err != nil {
			glog.Fatal(fmt.Sprintf(`[ghttp] AddStaticPath failed: %s`, err.Error()))
		} else {
			realPath = p
		}
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
		})
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
