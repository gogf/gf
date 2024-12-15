// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package fs_mixed

import (
	"io/fs"
	"sort"

	"github.com/gogf/gf/v2/os/gres/internal/defines"
	"github.com/gogf/gf/v2/os/gres/internal/fs_res"
	"github.com/gogf/gf/v2/os/gres/internal/fs_std"
)

// FS implements the FS interface by combining fs_res.FS and StdFS.
// It prioritizes using fs_res.FS and falls back to StdFS when file not found.
type FS struct {
	resFS *fs_res.FS
	stdFS *fs_std.FS
}

var _ defines.FS = (*FS)(nil)

// NewFS creates and returns a new FS.
func NewFS(resFS *fs_res.FS, stdFs fs.FS) *FS {
	return &FS{
		resFS: resFS,
		stdFS: fs_std.NewFS(stdFs),
	}
}

// Get returns the file with given path.
func (fs *FS) Get(path string) defines.File {
	if file := fs.resFS.Get(path); file != nil {
		return file
	}
	return fs.stdFS.Get(path)
}

// IsEmpty checks and returns whether the resource is empty.
func (fs *FS) IsEmpty() bool {
	return fs.resFS.IsEmpty() && fs.stdFS.IsEmpty()
}

// ScanDir returns the files under the given path,
// the parameter `path` should be a folder type.
func (fs *FS) ScanDir(path string, pattern string, recursive ...bool) []defines.File {
	var (
		filesMap = make(map[string]defines.File)
		files    = make([]defines.File, 0)
	)
	// Get files from StdFS
	stdFiles := fs.stdFS.ScanDir(path, pattern, recursive...)
	for _, file := range stdFiles {
		filesMap[file.Name()] = file
	}

	// Get files from fs_res.FS
	resFiles := fs.resFS.ScanDir(path, pattern, recursive...)
	for _, file := range resFiles {
		if _, exists := filesMap[file.Name()]; !exists {
			filesMap[file.Name()] = file
		}
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

func (fs *FS) ListAll() []defines.File {
	var (
		resAll   = fs.resFS.ListAll()
		stdAll   = fs.resFS.ListAll()
		filesMap = make(map[string]defines.File)
		files    = make([]defines.File, 0)
	)
	for _, file := range stdAll {
		filesMap[file.Name()] = file
	}
	for _, file := range resAll {
		filesMap[file.Name()] = file
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
