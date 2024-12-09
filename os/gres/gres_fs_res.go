// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gres

import (
	"context"
	"os"
	"path/filepath"
	"strings"

	"github.com/gogf/gf/v2/container/gtree"
	"github.com/gogf/gf/v2/internal/intlog"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/os/gtime"
)

// ResFS implements the FS interface using the default resource implementation.
type ResFS struct {
	tree *gtree.BTree // The tree storing all resource files.
}

var _ FS = (*ResFS)(nil)

const (
	defaultTreeM = 100
)

// NewResFS creates and returns a new ResFS.
func NewResFS() *ResFS {
	return &ResFS{
		tree: gtree.NewBTree(defaultTreeM, func(v1, v2 interface{}) int {
			return strings.Compare(v1.(string), v2.(string))
		}),
	}
}

// Get returns the file with given path.
func (fs *ResFS) Get(path string) File {
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
		return result.(File)
	}
	return nil
}

// IsEmpty checks and returns whether the resource is empty.
func (fs *ResFS) IsEmpty() bool {
	return fs.tree.IsEmpty()
}

// ScanDir returns the files under the given path,
// the parameter `path` should be a folder type.
func (fs *ResFS) ScanDir(path string, pattern string, recursive ...bool) []File {
	isRecursive := false
	if len(recursive) > 0 {
		isRecursive = recursive[0]
	}
	return fs.doScanDir(path, pattern, isRecursive)
}

// doScanDir is an internal method which scans directory
// and returns the absolute path list of files that are not sorted.
//
// The pattern parameter `pattern` supports multiple file name patterns,
// using the ',' symbol to separate multiple patterns.
//
// It scans directory recursively if given parameter `recursive` is true.
func (fs *ResFS) doScanDir(path string, pattern string, recursive bool) []File {
	path = strings.ReplaceAll(path, "\\", "/")
	path = strings.ReplaceAll(path, "//", "/")
	if path != "/" {
		for path[len(path)-1] == '/' {
			path = path[:len(path)-1]
		}
	}
	var (
		files    = make([]File, 0)
		patterns = strings.Split(pattern, ",")
	)
	for i := 0; i < len(patterns); i++ {
		patterns[i] = strings.TrimSpace(patterns[i])
	}

	// Get root directory
	rootFile := fs.Get(path)
	if rootFile == nil || !rootFile.FileInfo().IsDir() {
		return files
	}

	// Walk through the tree to find matching files
	fs.tree.IteratorAsc(func(key, value interface{}) bool {
		var (
			file     = value.(File)
			filePath = key.(string)
		)

		// Skip if not under the target path
		if !strings.HasPrefix(filePath, path) {
			return true
		}

		// Skip if not recursive and file is in subdirectory
		if !recursive {
			relPath := strings.TrimPrefix(filePath, path)
			if strings.Contains(relPath, "/") {
				return true
			}
		}

		// Check if file matches any pattern
		name := gfile.Basename(filePath)
		for _, p := range patterns {
			if match, _ := filepath.Match(p, name); match {
				files = append(files, file)
				break
			}
		}
		return true
	})
	return files
}

// Add adds the `content` into current ResFS with given `prefix`.
func (fs *ResFS) Add(content string, prefix ...string) error {
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
		files[i].(*localFile).fs = fs
		fs.tree.Set(namePrefix+files[i].Path(), files[i])
	}
	intlog.Printf(context.TODO(), "Add %d files to resource manager", fs.tree.Size())
	return nil
}

// Load loads, unpacks and adds the data from `path` into ResFS.
func (fs *ResFS) Load(path string, prefix ...string) error {
	realPath, err := gfile.Search(path)
	if err != nil {
		return err
	}
	return fs.Add(gfile.GetContents(realPath), prefix...)
}

// Dump prints the files of ResFS.
func (fs *ResFS) Dump() {
	var info os.FileInfo
	fs.tree.Iterator(func(key, value interface{}) bool {
		info = value.(File).FileInfo()
		intlog.Printf(
			context.TODO(),
			"%v %8s %s",
			gtime.New(info.ModTime()).ISO8601(),
			gfile.FormatSize(info.Size()),
			key,
		)
		return true
	})
	intlog.Printf(context.TODO(), "TOTAL FILES: %d", fs.tree.Size())
}
