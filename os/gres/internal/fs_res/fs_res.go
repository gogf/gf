// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package fs_res

import (
	"context"
	"path/filepath"
	"strings"

	"github.com/gogf/gf/v2/container/gtree"
	"github.com/gogf/gf/v2/internal/intlog"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/os/gres/internal/defines"
)

// FS implements the FS interface using the default resource implementation.
type FS struct {
	tree *gtree.BTree // The tree storing all resource files.
}

var _ defines.FS = (*FS)(nil)

const (
	defaultTreeM = 100
)

// NewFS creates and returns a new FS using resource manager.
func NewFS() *FS {
	return &FS{
		tree: gtree.NewBTree(defaultTreeM, func(v1, v2 interface{}) int {
			return strings.Compare(v1.(string), v2.(string))
		}),
	}
}

// Get returns the file with given path.
func (fs *FS) Get(path string) defines.File {
	if path == "" {
		return nil
	}
	path = strings.ReplaceAll(path, "\\", "/")
	path = strings.ReplaceAll(path, "//", "/")
	if path != "/" {
		for path[len(path)-1] == '/' {
			path = path[:len(path)-1]
		}
	}
	result := fs.tree.Get(path)
	if result != nil {
		return result.(defines.File)
	}
	return nil
}

// IsEmpty checks and returns whether the resource is empty.
func (fs *FS) IsEmpty() bool {
	return fs.tree.IsEmpty()
}

// ScanDir returns the files under the given path,
// the parameter `path` should be a folder type.
func (fs *FS) ScanDir(path string, pattern string, recursive ...bool) []defines.File {
	isRecursive := false
	if len(recursive) > 0 {
		isRecursive = recursive[0]
	}
	return fs.doScanDir(path, pattern, isRecursive, false)
}

// doScanDir is an internal method which scans directory
// and returns the absolute path list of files that are not sorted.
//
// The pattern parameter `pattern` supports multiple file name patterns,
// using the ',' symbol to separate multiple patterns.
//
// It scans directory recursively if given parameter `recursive` is true.
func (fs *FS) doScanDir(path string, pattern string, recursive bool, onlyFile bool) []defines.File {
	path = strings.ReplaceAll(path, "\\", "/")
	path = strings.ReplaceAll(path, "//", "/")
	if path != "/" {
		for path[len(path)-1] == '/' {
			path = path[:len(path)-1]
		}
	}
	var (
		name     = ""
		files    = make([]defines.File, 0)
		length   = len(path)
		patterns = strings.Split(pattern, ",")
	)
	for i := 0; i < len(patterns); i++ {
		patterns[i] = strings.TrimSpace(patterns[i])
	}

	// Used for type checking for first entry.
	first := true
	fs.tree.IteratorFrom(path, true, func(key, value interface{}) bool {
		if first {
			if !value.(defines.File).FileInfo().IsDir() {
				return false
			}
			first = false
		}
		if onlyFile && value.(defines.File).FileInfo().IsDir() {
			return true
		}
		name = key.(string)
		if len(name) <= length {
			return true
		}
		if path != name[:length] {
			return false
		}
		// To avoid of, eg: /i18n and /i18n-dir
		if !first && name[length] != '/' {
			return true
		}
		if !recursive {
			if strings.IndexByte(name[length+1:], '/') != -1 {
				return true
			}
		}
		for _, p := range patterns {
			if match, err := filepath.Match(p, gfile.Basename(name)); err == nil && match {
				files = append(files, value.(defines.File))
				return true
			}
		}
		return true
	})
	return files
}

func (fs *FS) ListAll() []defines.File {
	files := make([]defines.File, 0)
	fs.tree.Iterator(func(key, value interface{}) bool {
		files = append(files, value.(defines.File))
		return true
	})
	return files
}

// Add adds the `content` into current FS with given `prefix`.
func (fs *FS) Add(content string, prefix ...string) error {
	files, err := UnpackContent(content)
	if err != nil {
		intlog.Printf(context.TODO(), "Add resource files failed: %v", err)
		return err
	}
	namePrefix := ""
	if len(prefix) > 0 {
		namePrefix = prefix[0]
	}
	for i := 0; i < len(files); i++ {
		files[i].(*FileImp).fs = fs
		fs.tree.Set(namePrefix+files[i].Name(), files[i])
	}
	intlog.Printf(context.TODO(), "Add %d files to resource manager", fs.tree.Size())
	return nil
}

// Load loads, unpacks and adds the data from `path` into FS.
func (fs *FS) Load(path string, prefix ...string) error {
	realPath, err := gfile.Search(path)
	if err != nil {
		return err
	}
	return fs.Add(gfile.GetContents(realPath), prefix...)
}
