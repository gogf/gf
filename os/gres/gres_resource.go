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

// New creates and returns a new resource object.
func New() *Resource {
	return &Resource{
		tree: gtree.NewBTree(gDEFAULT_TREE_M, func(v1, v2 interface{}) int {
			return strings.Compare(v1.(string), v2.(string))
		}),
	}
}

// Add unpacks and adds the <content> into current resource object.
// The unnecessary parameter <prefix> indicates the prefix
// for each file storing into current resource object.
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
		files[i].resource = r
		r.tree.Set(namePrefix+files[i].file.Name, files[i])
	}
	return nil
}

// Load loads, unpacks and adds the data from <path> into current resource object.
// The unnecessary parameter <prefix> indicates the prefix
// for each file storing into current resource object.
func (r *Resource) Load(path string, prefix ...string) error {
	realPath, err := gfile.Search(path)
	if err != nil {
		return err
	}
	return r.Add(gfile.GetBytes(realPath), prefix...)
}

// Get returns the file with given path.
func (r *Resource) Get(path string) *File {
	if path != "/" {
		for path[len(path)-1] == '/' {
			path = path[:len(path)-1]
		}
	}
	result := r.tree.Get(path)
	if result != nil {
		return result.(*File)
	}
	return nil
}

// GetWithIndex searches file with <path>, if the file is directory
// it then does index files searching under this directory.
//
// GetWithIndex is usually used for http static file service.
func (r *Resource) GetWithIndex(path string, indexFiles []string) *File {
	if path != "/" {
		for path[len(path)-1] == '/' {
			path = path[:len(path)-1]
		}
	}
	if file := r.Get(path); file != nil {
		if len(indexFiles) > 0 && file.FileInfo().IsDir() {
			var f *File
			for _, name := range indexFiles {
				if f = r.Get(path + "/" + name); f != nil {
					return f
				}
			}
		}
		return file
	}
	return nil
}

// GetContent directly returns the content of <path>.
func (r *Resource) GetContent(path string) []byte {
	file := r.Get(path)
	if file != nil {
		return file.Content()
	}
	return nil
}

// Contains checks whether the <path> exists in current resource object.
func (r *Resource) Contains(path string) bool {
	return r.Get(path) != nil
}

// Scan returns the files under the given path, the parameter <path> should be a folder type.
//
// The pattern parameter <pattern> supports multiple file name patterns,
// using the ',' symbol to separate multiple patterns.
//
// It scans directory recursively if given parameter <recursive> is true.
func (r *Resource) Scan(path string, pattern string, recursive ...bool) []*File {
	if path != "/" {
		for path[len(path)-1] == '/' {
			path = path[:len(path)-1]
		}
	}
	name := ""
	files := make([]*File, 0)
	length := len(path)
	patterns := strings.Split(pattern, ",")
	for i := 0; i < len(patterns); i++ {
		patterns[i] = strings.TrimSpace(patterns[i])
	}
	// Used for type checking for first entry.
	first := true
	r.tree.IteratorFrom(path, true, func(key, value interface{}) bool {
		if first {
			if !value.(*File).FileInfo().IsDir() {
				return false
			}
			first = false
		}
		name = key.(string)
		if len(name) <= length {
			return true
		}
		if path != name[:length] {
			return false
		}
		if len(recursive) == 0 || !recursive[0] {
			if strings.IndexByte(name[length+1:], '/') != -1 {
				return true
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

// Dump prints the files of current resource object.
func (r *Resource) Dump() {
	r.tree.Iterator(func(key, value interface{}) bool {
		fmt.Printf("%7s %s\n", gfile.FormatSize(value.(*File).FileInfo().Size()), key)
		return true
	})
	fmt.Printf("TOTAL FILES: %d\n", r.tree.Size())
}
