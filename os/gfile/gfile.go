// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package gfile provides easy-to-use operations for file system.
package gfile

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/gogf/gf/v2/container/gtype"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gconv"
)

const (
	// Separator for file system.
	// It here defines the separator as variable
	// to allow it modified by developer if necessary.
	Separator = string(filepath.Separator)

	// DefaultPermOpen is the default perm for file opening.
	DefaultPermOpen = os.FileMode(0666)

	// DefaultPermCopy is the default perm for file/folder copy.
	DefaultPermCopy = os.FileMode(0755)
)

var (
	// The absolute file path for main package.
	// It can be only checked and set once.
	mainPkgPath = gtype.NewString()

	// selfPath is the current running binary path.
	// As it is most commonly used, it is so defined as an internal package variable.
	selfPath = ""
)

func init() {
	// Initialize internal package variable: selfPath.
	selfPath, _ = exec.LookPath(os.Args[0])
	if selfPath != "" {
		selfPath, _ = filepath.Abs(selfPath)
	}
	if selfPath == "" {
		selfPath, _ = filepath.Abs(os.Args[0])
	}
}

// Mkdir creates directories recursively with given `path`.
// The parameter `path` is suggested to be an absolute path instead of relative one.
func Mkdir(path string) (err error) {
	if err = os.MkdirAll(path, os.ModePerm); err != nil {
		err = gerror.Wrapf(err, `os.MkdirAll failed for path "%s" with perm "%d"`, path, os.ModePerm)
		return err
	}
	return nil
}

// Create creates a file with given `path` recursively.
// The parameter `path` is suggested to be absolute path.
func Create(path string) (*os.File, error) {
	dir := Dir(path)
	if !Exists(dir) {
		if err := Mkdir(dir); err != nil {
			return nil, err
		}
	}
	file, err := os.Create(path)
	if err != nil {
		err = gerror.Wrapf(err, `os.Create failed for name "%s"`, path)
	}
	return file, err
}

// Open opens file/directory READONLY.
func Open(path string) (*os.File, error) {
	file, err := os.Open(path)
	if err != nil {
		err = gerror.Wrapf(err, `os.Open failed for name "%s"`, path)
	}
	return file, err
}

// OpenFile opens file/directory with custom `flag` and `perm`.
// The parameter `flag` is like: O_RDONLY, O_RDWR, O_RDWR|O_CREATE|O_TRUNC, etc.
func OpenFile(path string, flag int, perm os.FileMode) (*os.File, error) {
	file, err := os.OpenFile(path, flag, perm)
	if err != nil {
		err = gerror.Wrapf(err, `os.OpenFile failed with name "%s", flag "%d", perm "%d"`, path, flag, perm)
	}
	return file, err
}

// OpenWithFlag opens file/directory with default perm and custom `flag`.
// The default `perm` is 0666.
// The parameter `flag` is like: O_RDONLY, O_RDWR, O_RDWR|O_CREATE|O_TRUNC, etc.
func OpenWithFlag(path string, flag int) (*os.File, error) {
	file, err := OpenFile(path, flag, DefaultPermOpen)
	if err != nil {
		return nil, err
	}
	return file, nil
}

// OpenWithFlagPerm opens file/directory with custom `flag` and `perm`.
// The parameter `flag` is like: O_RDONLY, O_RDWR, O_RDWR|O_CREATE|O_TRUNC, etc.
// The parameter `perm` is like: 0600, 0666, 0777, etc.
func OpenWithFlagPerm(path string, flag int, perm os.FileMode) (*os.File, error) {
	file, err := OpenFile(path, flag, perm)
	if err != nil {
		return nil, err
	}
	return file, nil
}

// Join joins string array paths with file separator of current system.
func Join(paths ...string) string {
	var s string
	for _, path := range paths {
		if s != "" {
			s += Separator
		}
		s += gstr.TrimRight(path, Separator)
	}
	return s
}

// Exists checks whether given `path` exist.
func Exists(path string) bool {
	if stat, err := os.Stat(path); stat != nil && !os.IsNotExist(err) {
		return true
	}
	return false
}

// IsDir checks whether given `path` a directory.
// Note that it returns false if the `path` does not exist.
func IsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

// Pwd returns absolute path of current working directory.
// Note that it returns an empty string if retrieving current
// working directory failed.
func Pwd() string {
	path, err := os.Getwd()
	if err != nil {
		return ""
	}
	return path
}

// Chdir changes the current working directory to the named directory.
// If there is an error, it will be of type *PathError.
func Chdir(dir string) (err error) {
	err = os.Chdir(dir)
	if err != nil {
		err = gerror.Wrapf(err, `os.Chdir failed with dir "%s"`, dir)
	}
	return
}

// IsFile checks whether given `path` a file, which means it's not a directory.
// Note that it returns false if the `path` does not exist.
func IsFile(path string) bool {
	s, err := Stat(path)
	if err != nil {
		return false
	}
	return !s.IsDir()
}

// Stat returns a FileInfo describing the named file.
// If there is an error, it will be of type *PathError.
func Stat(path string) (os.FileInfo, error) {
	info, err := os.Stat(path)
	if err != nil {
		err = gerror.Wrapf(err, `os.Stat failed for file "%s"`, path)
	}
	return info, err
}

// Move renames (moves) `src` to `dst` path.
// If `dst` already exists and is not a directory, it'll be replaced.
func Move(src string, dst string) (err error) {
	err = os.Rename(src, dst)
	if err != nil {
		err = gerror.Wrapf(err, `os.Rename failed from "%s" to "%s"`, src, dst)
	}
	return
}

// Rename is alias of Move.
// See Move.
func Rename(src string, dst string) error {
	return Move(src, dst)
}

// DirNames returns sub-file names of given directory `path`.
// Note that the returned names are NOT absolute paths.
func DirNames(path string) ([]string, error) {
	f, err := Open(path)
	if err != nil {
		return nil, err
	}
	list, err := f.Readdirnames(-1)
	_ = f.Close()
	if err != nil {
		err = gerror.Wrapf(err, `Read dir files failed from path "%s"`, path)
		return nil, err
	}
	return list, nil
}

// Glob returns the names of all files matching pattern or nil
// if there is no matching file. The syntax of patterns is the same
// as in Match. The pattern may describe hierarchical names such as
// /usr/*/bin/ed (assuming the Separator is '/').
//
// Glob ignores file system errors such as I/O errors reading directories.
// The only possible returned error is ErrBadPattern, when pattern
// is malformed.
func Glob(pattern string, onlyNames ...bool) ([]string, error) {
	list, err := filepath.Glob(pattern)
	if err != nil {
		err = gerror.Wrapf(err, `filepath.Glob failed for pattern "%s"`, pattern)
		return nil, err
	}
	if len(onlyNames) > 0 && onlyNames[0] && len(list) > 0 {
		array := make([]string, len(list))
		for k, v := range list {
			array[k] = Basename(v)
		}
		return array, nil
	}
	return list, nil
}

// Remove deletes all file/directory with `path` parameter.
// If parameter `path` is directory, it deletes it recursively.
//
// It does nothing if given `path` does not exist or is empty.
func Remove(path string) (err error) {
	// It does nothing if `path` is empty.
	if path == "" {
		return nil
	}
	if err = os.RemoveAll(path); err != nil {
		err = gerror.Wrapf(err, `os.RemoveAll failed for path "%s"`, path)
	}
	return
}

// IsReadable checks whether given `path` is readable.
func IsReadable(path string) bool {
	result := true
	file, err := os.OpenFile(path, os.O_RDONLY, DefaultPermOpen)
	if err != nil {
		result = false
	}
	file.Close()
	return result
}

// IsWritable checks whether given `path` is writable.
//
// TODO improve performance; use golang.org/x/sys to cross-plat-form
func IsWritable(path string) bool {
	result := true
	if IsDir(path) {
		// If it's a directory, create a temporary file to test whether it's writable.
		tmpFile := strings.TrimRight(path, Separator) + Separator + gconv.String(time.Now().UnixNano())
		if f, err := Create(tmpFile); err != nil || !Exists(tmpFile) {
			result = false
		} else {
			_ = f.Close()
			_ = Remove(tmpFile)
		}
	} else {
		// If it's a file, check if it can open it.
		file, err := os.OpenFile(path, os.O_WRONLY, DefaultPermOpen)
		if err != nil {
			result = false
		}
		_ = file.Close()
	}
	return result
}

// Chmod is alias of os.Chmod.
// See os.Chmod.
func Chmod(path string, mode os.FileMode) (err error) {
	err = os.Chmod(path, mode)
	if err != nil {
		err = gerror.Wrapf(err, `os.Chmod failed with path "%s" and mode "%s"`, path, mode)
	}
	return
}

// Abs returns an absolute representation of path.
// If the path is not absolute it will be joined with the current
// working directory to turn it into an absolute path. The absolute
// path name for a given file is not guaranteed to be unique.
// Abs calls Clean on the result.
func Abs(path string) string {
	p, _ := filepath.Abs(path)
	return p
}

// RealPath converts the given `path` to its absolute path
// and checks if the file path exists.
// If the file does not exist, return an empty string.
func RealPath(path string) string {
	p, err := filepath.Abs(path)
	if err != nil {
		return ""
	}
	if !Exists(p) {
		return ""
	}
	return p
}

// SelfPath returns absolute file path of current running process(binary).
func SelfPath() string {
	return selfPath
}

// SelfName returns file name of current running process(binary).
func SelfName() string {
	return Basename(SelfPath())
}

// SelfDir returns absolute directory path of current running process(binary).
func SelfDir() string {
	return filepath.Dir(SelfPath())
}

// Basename returns the last element of path, which contains file extension.
// Trailing path separators are removed before extracting the last element.
// If the path is empty, Base returns ".".
// If the path consists entirely of separators, Basename returns a single separator.
//
// Example:
// Basename("/var/www/file.js") -> file.js
// Basename("file.js")          -> file.js
func Basename(path string) string {
	return filepath.Base(path)
}

// Name returns the last element of path without file extension.
//
// Example:
// Name("/var/www/file.js") -> file
// Name("file.js")          -> file
func Name(path string) string {
	base := filepath.Base(path)
	if i := strings.LastIndexByte(base, '.'); i != -1 {
		return base[:i]
	}
	return base
}

// Dir returns all but the last element of path, typically the path's directory.
// After dropping the final element, Dir calls Clean on the path and trailing
// slashes are removed.
// If the `path` is empty, Dir returns ".".
// If the `path` is ".", Dir treats the path as current working directory.
// If the `path` consists entirely of separators, Dir returns a single separator.
// The returned path does not end in a separator unless it is the root directory.
//
// Example:
// Dir("/var/www/file.js") -> "/var/www"
// Dir("file.js")          -> "."
func Dir(path string) string {
	if path == "." {
		return filepath.Dir(RealPath(path))
	}
	return filepath.Dir(path)
}

// IsEmpty checks whether the given `path` is empty.
// If `path` is a folder, it checks if there's any file under it.
// If `path` is a file, it checks if the file size is zero.
//
// Note that it returns true if `path` does not exist.
func IsEmpty(path string) bool {
	stat, err := Stat(path)
	if err != nil {
		return true
	}
	if stat.IsDir() {
		file, err := os.Open(path)
		if err != nil {
			return true
		}
		defer file.Close()
		names, err := file.Readdirnames(-1)
		if err != nil {
			return true
		}
		return len(names) == 0
	} else {
		return stat.Size() == 0
	}
}

// Ext returns the file name extension used by path.
// The extension is the suffix beginning at the final dot
// in the final element of path; it is empty if there is
// no dot.
// Note: the result contains symbol '.'.
//
// Example:
// Ext("main.go")  => .go
// Ext("api.json") => .json
func Ext(path string) string {
	ext := filepath.Ext(path)
	if p := strings.IndexByte(ext, '?'); p != -1 {
		ext = ext[0:p]
	}
	return ext
}

// ExtName is like function Ext, which returns the file name extension used by path,
// but the result does not contain symbol '.'.
//
// Example:
// ExtName("main.go")  => go
// ExtName("api.json") => json
func ExtName(path string) string {
	return strings.TrimLeft(Ext(path), ".")
}

// Temp retrieves and returns the temporary directory of current system.
//
// The optional parameter `names` specifies the sub-folders/sub-files,
// which will be joined with current system separator and returned with the path.
func Temp(names ...string) string {
	path := os.TempDir()
	for _, name := range names {
		path = Join(path, name)
	}
	return path
}
