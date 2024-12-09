// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gres

import (
	"io/fs"
)

// MixedFS implements the FS interface by combining StdFS and ResFS.
// It prioritizes using StdFS and falls back to ResFS when file not found.
type MixedFS struct {
	stdFS *StdFS
	resFS *ResFS
}

// NewMixedFS creates and returns a new MixedFS.
func NewMixedFS(stdFs fs.FS, resFS *ResFS) *MixedFS {
	return &MixedFS{
		resFS: resFS,
		stdFS: NewStdFS(stdFs),
	}
}

// Get returns the file with given path.
func (fs *MixedFS) Get(path string) *File {
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
func (fs *MixedFS) ScanDir(path string, pattern string, recursive ...bool) []*File {
	var (
		filesMap = make(map[string]*File)
		files    = make([]*File, 0)
	)

	// Get files from ResFS
	defaultFiles := fs.resFS.ScanDir(path, pattern, recursive...)
	for _, file := range defaultFiles {
		if _, exists := filesMap[file.Path()]; !exists {
			filesMap[file.Path()] = file
		}
	}

	// Get files from StdFS
	stdFiles := fs.stdFS.ScanDir(path, pattern, recursive...)
	if len(stdFiles) > 0 {

	}
	for _, file := range stdFiles {
		filesMap[file.Path()] = file
	}

	// Convert map to slice
	for _, file := range filesMap {
		files = append(files, file)
	}

	return files
}

// ScanDirFile returns all sub-files with absolute paths of given `path`,
// It scans directory recursively if given parameter `recursive` is true.
func (fs *MixedFS) ScanDirFile(path string, pattern string, recursive ...bool) []*File {
	return fs.ScanDir(path, pattern, recursive...)
}

// Export exports and saves specified path `src` and all its sub files
// to specified system path `dst` recursively.
func (fs *MixedFS) Export(src, dst string, option ...ExportOption) error {
	if file := fs.Get(src); file != nil {
		return file.Export(dst, option...)
	}
	return nil
}
