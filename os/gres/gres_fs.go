// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gres

// FS is the interface that defines a virtual file system.
type FS interface {
	// Get returns the file with given path.
	Get(path string) File

	// IsEmpty checks and returns whether the resource is empty.
	IsEmpty() bool

	// ScanDir returns the files under the given path,
	// the parameter `path` should be a folder type.
	ScanDir(path string, pattern string, recursive ...bool) []File
}

// ExportOption contains options for Export.
type ExportOption struct {
	RemovePrefix string // Remove the prefix from source file before export.
}
