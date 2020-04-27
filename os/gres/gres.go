// Copyright 2019 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package gres provides resource management and packing/unpacking feature between files and bytes.
package gres

const (
	// Separator for directories.
	Separator = "/"
)

var (
	// Default resource object.
	defaultResource = Instance()
)

// Add unpacks and adds the <content> into the default resource object.
// The unnecessary parameter <prefix> indicates the prefix
// for each file storing into current resource object.
func Add(content string, prefix ...string) error {
	return defaultResource.Add(content, prefix...)
}

// Load loads, unpacks and adds the data from <path> into the default resource object.
// The unnecessary parameter <prefix> indicates the prefix
// for each file storing into current resource object.
func Load(path string, prefix ...string) error {
	return defaultResource.Load(path, prefix...)
}

// Get returns the file with given path.
func Get(path string) *File {
	return defaultResource.Get(path)
}

// GetWithIndex searches file with <path>, if the file is directory
// it then does index files searching under this directory.
//
// GetWithIndex is usually used for http static file service.
func GetWithIndex(path string, indexFiles []string) *File {
	return defaultResource.GetWithIndex(path, indexFiles)
}

// GetContent directly returns the content of <path> in default resource object.
func GetContent(path string) []byte {
	return defaultResource.GetContent(path)
}

// Contains checks whether the <path> exists in the default resource object.
func Contains(path string) bool {
	return defaultResource.Contains(path)
}

// IsEmpty checks and returns whether the resource manager is empty.
func IsEmpty() bool {
	return defaultResource.tree.IsEmpty()
}

// ScanDir returns the files under the given path, the parameter <path> should be a folder type.
//
// The pattern parameter <pattern> supports multiple file name patterns,
// using the ',' symbol to separate multiple patterns.
//
// It scans directory recursively if given parameter <recursive> is true.
func ScanDir(path string, pattern string, recursive ...bool) []*File {
	return defaultResource.ScanDir(path, pattern, recursive...)
}

// ScanDirFile returns all sub-files with absolute paths of given <path>,
// It scans directory recursively if given parameter <recursive> is true.
//
// Note that it returns only files, exclusive of directories.
func ScanDirFile(path string, pattern string, recursive ...bool) []*File {
	return defaultResource.ScanDirFile(path, pattern, recursive...)
}

// Dump prints the files of the default resource object.
func Dump() {
	defaultResource.Dump()
}
