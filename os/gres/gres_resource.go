// Copyright 2019 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gres

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/gogf/gf/container/gtree"
	"github.com/gogf/gf/os/gfile"
)

type Resource struct {
	tree *gtree.BTree
}

const (
	gDEFAULT_TREE_M = 100
)

func New() *Resource {
	return &Resource{
		tree: gtree.NewBTree(gDEFAULT_TREE_M, func(v1, v2 interface{}) int {
			return strings.Compare(v1.(string), v2.(string))
		}),
	}
}

func (r *Resource) Add(content []byte, prefix ...string) error {
	files, err := UnpackContent(content)
	if err != nil {
		return err
	}
	namePrefix := ""
	if len(prefix) > 0 {
		namePrefix = prefix[0]
	}
	for i := 0; i < len(files); i++ {
		r.tree.Set(namePrefix+files[i].zipFile.Name, files[i])
	}
	return nil
}

func (r *Resource) Load(path string, prefix ...string) error {
	realPath, err := gfile.Search(path)
	if err != nil {
		return err
	}
	return r.Add(gfile.GetBytes(realPath), prefix...)
}

func (r *Resource) Get(path string) *File {
	result := r.tree.Get(path)
	if result != nil {
		return result.(*File)
	}
	return nil
}

func (r *Resource) Scan(path string, pattern string, recursive ...bool) []*File {
	if path != "/" {
		path = strings.TrimRight(path, "/\\")
	}
	name := ""
	files := make([]*File, 0)
	length := len(path)
	patterns := strings.Split(pattern, ",")
	for i := 0; i < len(patterns); i++ {
		patterns[i] = strings.TrimSpace(patterns[i])
	}
	r.tree.IteratorFrom(path, func(key, value interface{}) bool {
		name = key.(string)
		if path != name[:length] {
			return false
		}
		if len(recursive) == 0 || !recursive[0] {
			if strings.IndexByte(name[length:], '/') != -1 {
				return false
			}
		}
		for _, p := range patterns {
			if match, err := filepath.Match(p, gfile.Basename(name)); err == nil && match {
				files = append(files, value.(*File))
				return true
			}
		}
		return true
	})
	return files
}

func (r *Resource) Dump() {
	r.tree.Iterator(func(key, value interface{}) bool {
		fmt.Printf("%7s %s\n", gfile.FormatSize(value.(*File).FileInfo().Size()), key)
		return true
	})
	fmt.Printf("TOTAL %d FILES", r.tree.Size())
}
