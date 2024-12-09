// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gres

import (
	"context"
	"strings"

	"github.com/gogf/gf/v2/internal/intlog"
)

// Resource implements the FS interface.
type Resource struct {
	fs FS
}

// New creates and returns a new resource object.
func New() *Resource {
	return NewWithFS(NewResFS())
}

// NewWithFS sets the underlying file system implementation.
func NewWithFS(fs FS) *Resource {
	return &Resource{
		fs: fs,
	}
}

// Get returns the file with given path.
func (r *Resource) Get(path string) *File {
	return r.fs.Get(path)
}

// GetWithIndex searches file with `path`, if the file is directory
// it then does index files searching under this directory.
func (r *Resource) GetWithIndex(path string, indexFiles []string) *File {
	// Necessary for double char '/' replacement in prefix.
	path = strings.ReplaceAll(path, "\\", "/")
	path = strings.ReplaceAll(path, "//", "/")
	if path != "/" {
		for path[len(path)-1] == '/' {
			path = path[:len(path)-1]
		}
	}
	if file := r.fs.Get(path); file != nil {
		if len(indexFiles) > 0 && file.FileInfo().IsDir() {
			var f *File
			for _, name := range indexFiles {
				if f = r.fs.Get(path + "/" + name); f != nil {
					return f
				}
			}
		}
		return file
	}
	return nil
}

// GetContent directly returns the content of `path`.
func (r *Resource) GetContent(path string) []byte {
	file := r.Get(path)
	if file != nil {
		return file.Content()
	}
	return nil
}

// Contains checks whether the `path` exists in current resource object.
func (r *Resource) Contains(path string) bool {
	return r.Get(path) == nil
}

// IsEmpty checks and returns whether the resource manager is empty.
func (r *Resource) IsEmpty() bool {
	return r.fs.IsEmpty()
}

// ScanDir returns the files under the given path, the parameter `path` should be a folder type.
func (r *Resource) ScanDir(path string, pattern string, recursive ...bool) []*File {
	return r.fs.ScanDir(path, pattern, recursive...)
}

// ScanDirFile returns all sub-files with absolute paths of given `path`,
// It scans directory recursively if given parameter `recursive` is true.
func (r *Resource) ScanDirFile(path string, pattern string, recursive ...bool) []*File {
	var (
		result = make([]*File, 0)
		files  = r.fs.ScanDir(path, pattern, recursive...)
	)
	for _, file := range files {
		if file.FileInfo().IsDir() {
			continue
		}
		result = append(result, file)
	}
	return result
}

// Export exports and saves specified path `src` and all its sub files
// to specified system path `dst` recursively.
func (r *Resource) Export(src, dst string, option ...ExportOption) error {
	if file := r.Get(src); file != nil {
		return file.Export(dst, option...)
	}
	return nil
}

// Dump prints the files of current resource object.
func (r *Resource) Dump() {
	var ctx = context.TODO()
	if r.IsEmpty() {
		intlog.Printf(ctx, "empty resource")
	} else {
		for _, v := range r.ScanDir("/", "*", true) {
			intlog.Printf(ctx, "%s %d", v.Path(), v.FileInfo().Size())
		}
	}
}
