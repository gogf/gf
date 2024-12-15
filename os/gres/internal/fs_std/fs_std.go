// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package fs_std

import (
	"io/fs"
	"os"
	"path/filepath"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/os/gres/internal/defines"
)

// FS implements the FS interface using the standard library fs.FS.
type FS struct {
	fs fs.FS
}

var _ defines.FS = (*FS)(nil)

func NewFS(fs fs.FS) *FS {
	return &FS{
		fs: fs,
	}
}

// Get returns the file with given path.
func (fs *FS) Get(path string) defines.File {
	f, err := fs.fs.Open(path)
	if err != nil {
		return nil
	}
	file := &FileImp{
		path: path,
		file: f,
		fs:   fs,
	}
	return file
}

// IsEmpty checks and returns whether the resource is empty.
func (fs *FS) IsEmpty() bool {
	if dir, ok := fs.fs.(interface {
		ReadDir(name string) ([]os.DirEntry, error)
	}); ok {
		entries, err := dir.ReadDir(".")
		if err != nil {
			return true
		}
		return len(entries) == 0
	}
	return true
}

// ScanDir returns the files under the given path,
// the parameter `path` should be a folder type.
func (fs *FS) ScanDir(path string, pattern string, recursive ...bool) []defines.File {
	var (
		files       = make([]defines.File, 0)
		isRecursive = len(recursive) > 0 && recursive[0]
	)
	err := fs.walkDir(path, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			matched, err := filepath.Match(pattern, filepath.Base(path))
			if err != nil {
				return err
			}
			if matched {
				if file := fs.Get(path); file != nil {
					files = append(files, file)
				}
			}
		}
		if !isRecursive && d.IsDir() && path != "." {
			return filepath.SkipDir
		}
		return nil
	})

	if err != nil {
		return nil
	}
	return files
}

// walkDir walks the file tree rooted at path, calling fn for each file or
// directory in the tree, including path.
func (fs *FS) walkDir(path string, fn func(path string, d os.DirEntry, err error) error) error {
	if dir, ok := fs.fs.(interface {
		ReadDir(name string) ([]os.DirEntry, error)
	}); ok {
		entries, err := dir.ReadDir(path)
		if err != nil {
			err = fn(path, nil, err)
			if err != nil {
				return err
			}
			return nil
		}

		for _, entry := range entries {
			var (
				fileName = entry.Name()
				filePath = filepath.Join(path, fileName)
			)
			err = fn(filePath, entry, nil)
			if err != nil {
				if gerror.Is(err, filepath.SkipDir) {
					if entry.IsDir() {
						continue
					}
					return nil
				}
				return err
			}
			if entry.IsDir() {
				err = fs.walkDir(filePath, fn)
				if err != nil {
					if gerror.Is(err, filepath.SkipDir) {
						continue
					}
					return err
				}
			}
		}
		return nil
	}
	return gerror.New("filesystem does not implement ReadDir")
}

func (fs *FS) ListAll() []defines.File {
	return fs.ScanDir(".", "*", true)
}
