// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gres

import (
	"io/fs"
	"sort"
)

// MixedFS implements the FS interface by combining StdFS and ResFS.
// It prioritizes using StdFS and falls back to ResFS when file not found.
type MixedFS struct {
	stdFS *StdFS
	resFS *ResFS
}

var _ FS = (*MixedFS)(nil)

// NewMixedFS creates and returns a new MixedFS.
func NewMixedFS(stdFs fs.FS, resFS *ResFS) *MixedFS {
	return &MixedFS{
		resFS: resFS,
		stdFS: NewStdFS(stdFs),
	}
}

// Get returns the file with given path.
func (fs *MixedFS) Get(path string) File {
	if file := fs.resFS.Get(path); file != nil {
		return file
	}
	return fs.stdFS.Get(path)
}

// IsEmpty checks and returns whether the resource is empty.
func (fs *MixedFS) IsEmpty() bool {
	return fs.resFS.IsEmpty() && fs.stdFS.IsEmpty()
}

// ScanDir returns the files under the given path,
// the parameter `path` should be a folder type.
func (fs *MixedFS) ScanDir(path string, pattern string, recursive ...bool) []File {
	var (
		filesMap = make(map[string]File)
		files    = make([]File, 0)
	)

	// Get files from ResFS
	resFiles := fs.resFS.ScanDir(path, pattern, recursive...)
	for _, file := range resFiles {
		if _, exists := filesMap[file.Path()]; !exists {
			filesMap[file.Path()] = file
		}
	}

	// Get files from StdFS
	stdFiles := fs.stdFS.ScanDir(path, pattern, recursive...)
	for _, file := range stdFiles {
		filesMap[file.Path()] = file
	}

	// Convert map to slice and sort by path
	paths := make([]string, 0, len(filesMap))
	for filePath := range filesMap {
		paths = append(paths, filePath)
	}
	sort.Strings(paths)

	// Build sorted result
	for _, filePath := range paths {
		files = append(files, filesMap[filePath])
	}

	return files
}
