// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package gfile provides easy-to-use operations for file system.
package gfile

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"time"

	"github.com/gogf/gf/g/container/gtype"
	"github.com/gogf/gf/g/text/gregex"
	"github.com/gogf/gf/g/text/gstr"
	"github.com/gogf/gf/g/util/gconv"
)

const (
	// Separator for file system.
	Separator = string(filepath.Separator)
	// Default perm for file opening.
	gDEFAULT_PERM = 0666
)

var (
	// The absolute file path for main package.
	// It can be only checked and set once.
	mainPkgPath = gtype.NewString()
)

// Mkdir creates directories recursively with given <path>.
// The parameter <path> is suggested to be absolute path.
func Mkdir(path string) error {
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

// Create creates file with given <path> recursively.
// The parameter <path> is suggested to be absolute path.
func Create(path string) (*os.File, error) {
	dir := Dir(path)
	if !Exists(dir) {
		Mkdir(dir)
	}
	return os.Create(path)
}

// Open opens file/directory readonly.
func Open(path string) (*os.File, error) {
	return os.Open(path)
}

// OpenFile opens file/directory with given <flag> and <perm>.
func OpenFile(path string, flag int, perm os.FileMode) (*os.File, error) {
	return os.OpenFile(path, flag, perm)
}

// OpenWithFlag opens file/directory with default perm and given <flag>.
func OpenWithFlag(path string, flag int) (*os.File, error) {
	f, err := os.OpenFile(path, flag, gDEFAULT_PERM)
	if err != nil {
		return nil, err
	}
	return f, nil
}

// OpenWithFlagPerm opens file/directory with given <flag> and <perm>.
func OpenWithFlagPerm(path string, flag int, perm int) (*os.File, error) {
	f, err := os.OpenFile(path, flag, os.FileMode(perm))
	if err != nil {
		return nil, err
	}
	return f, nil
}

// Exists checks whether given <path> exist.
func Exists(path string) bool {
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		return true
	}
	return false
}

// IsDir checks whether given <path> a directory.
func IsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

// Pwd returns absolute path of current working directory.
func Pwd() string {
	path, _ := os.Getwd()
	return path
}

// IsFile checks whether given <path> a file, which means it's not a directory.
func IsFile(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !s.IsDir()
}

// Alias of Stat.
// See Stat.
func Info(path string) (os.FileInfo, error) {
	return Stat(path)
}

// Stat returns a FileInfo describing the named file.
// If there is an error, it will be of type *PathError.
func Stat(path string) (os.FileInfo, error) {
	return os.Stat(path)
}

// Move renames (moves) <src> to <dst> path.
func Move(src string, dst string) error {
	return os.Rename(src, dst)
}

// Alias of Move.
// See Move.
func Rename(src string, dst string) error {
	return Move(src, dst)
}

// Copy file/directory from <src> to <dst>.
//
// If <src> is file, it calls CopyFile to implements copy feature,
// or else it calls CopyDir.
func Copy(src string, dst string) error {
	if IsFile(src) {
		return CopyFile(src, dst)
	}
	return CopyDir(src, dst)
}

// CopyFile copies the contents of the file named src to the file named
// by dst. The file will be created if it does not already exist. If the
// destination file exists, all it's contents will be replaced by the contents
// of the source file. The file mode will be copied from the source and
// the copied data is synced/flushed to stable storage.
// Thanks: https://gist.github.com/r0l1/92462b38df26839a3ca324697c8cba04
func CopyFile(src, dst string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return
	}
	defer func() {
		if e := in.Close(); e != nil {
			err = e
		}
	}()
	out, err := os.Create(dst)
	if err != nil {
		return
	}
	defer func() {
		if e := out.Close(); e != nil {
			err = e
		}
	}()
	_, err = io.Copy(out, in)
	if err != nil {
		return
	}
	err = out.Sync()
	if err != nil {
		return
	}
	si, err := os.Stat(src)
	if err != nil {
		return
	}
	err = os.Chmod(dst, si.Mode())
	if err != nil {
		return
	}
	return
}

// CopyDir recursively copies a directory tree, attempting to preserve permissions.
// Source directory must exist, destination directory must *not* exist.
// Symlinks are ignored and skipped.
func CopyDir(src string, dst string) (err error) {
	src = filepath.Clean(src)
	dst = filepath.Clean(dst)
	si, err := os.Stat(src)
	if err != nil {
		return err
	}
	if !si.IsDir() {
		return fmt.Errorf("source is not a directory")
	}
	_, err = os.Stat(dst)
	if err != nil && !os.IsNotExist(err) {
		return
	}
	if err == nil {
		return fmt.Errorf("destination already exists")
	}
	err = os.MkdirAll(dst, si.Mode())
	if err != nil {
		return
	}
	entries, err := ioutil.ReadDir(src)
	if err != nil {
		return
	}
	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())
		if entry.IsDir() {
			err = CopyDir(srcPath, dstPath)
			if err != nil {
				return
			}
		} else {
			// Skip symlinks.
			if entry.Mode()&os.ModeSymlink != 0 {
				continue
			}
			err = CopyFile(srcPath, dstPath)
			if err != nil {
				return
			}
		}
	}
	return
}

// DirNames returns sub-file names of given directory <path>.
func DirNames(path string) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	list, err := f.Readdirnames(-1)
	f.Close()
	if err != nil {
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
	if list, err := filepath.Glob(pattern); err == nil {
		if len(onlyNames) > 0 && onlyNames[0] && len(list) > 0 {
			array := make([]string, len(list))
			for k, v := range list {
				array[k] = Basename(v)
			}
			return array, nil
		}
		return list, nil
	} else {
		return nil, err
	}
}

// Remove deletes all file/directory with <path> parameter.
// If parameter <path> is directory, it deletes it recursively.
func Remove(path string) error {
	return os.RemoveAll(path)
}

// IsReadable checks whether given <path> is readable.
func IsReadable(path string) bool {
	result := true
	file, err := os.OpenFile(path, os.O_RDONLY, gDEFAULT_PERM)
	if err != nil {
		result = false
	}
	file.Close()
	return result
}

// IsWritable checks whether given <path> is writable.
//
// @TODO improve performance; use golang.org/x/sys to cross-plat-form
func IsWritable(path string) bool {
	result := true
	if IsDir(path) {
		// If it's a directory, create a temporary file to test whether it's writable.
		tmpFile := strings.TrimRight(path, Separator) + Separator + gconv.String(time.Now().UnixNano())
		if f, err := Create(tmpFile); err != nil || !Exists(tmpFile) {
			result = false
		} else {
			f.Close()
			Remove(tmpFile)
		}
	} else {
		// 如果是文件，那么判断文件是否可打开
		file, err := os.OpenFile(path, os.O_WRONLY, gDEFAULT_PERM)
		if err != nil {
			result = false
		}
		file.Close()
	}
	return result
}

// See os.Chmod.
func Chmod(path string, mode os.FileMode) error {
	return os.Chmod(path, mode)
}

// ScanDir returns all sub-files with absolute paths of given <path>,
// It scans directory recursively if given parameter <recursive> is true.
func ScanDir(path string, pattern string, recursive ...bool) ([]string, error) {
	list, err := doScanDir(path, pattern, recursive...)
	if err != nil {
		return nil, err
	}
	if len(list) > 0 {
		sort.Strings(list)
	}
	return list, nil
}

// doScanDir is an internal method which scans directory
// and returns the absolute path list of files that are not sorted.
//
// The pattern parameter <pattern> supports multiple file name patterns,
// using the ',' symbol to separate multiple patterns.
//
// It scans directory recursively if given parameter <recursive> is true.
func doScanDir(path string, pattern string, recursive ...bool) ([]string, error) {
	list := ([]string)(nil)
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	names, err := file.Readdirnames(-1)
	if err != nil {
		return nil, err
	}
	for _, name := range names {
		path := fmt.Sprintf("%s%s%s", path, Separator, name)
		if IsDir(path) && len(recursive) > 0 && recursive[0] {
			array, _ := doScanDir(path, pattern, true)
			if len(array) > 0 {
				list = append(list, array...)
			}
		}
		// If it meets pattern, then add it to the result list.
		for _, p := range strings.Split(pattern, ",") {
			if match, err := filepath.Match(strings.TrimSpace(p), name); err == nil && match {
				list = append(list, path)
			}
		}
	}
	return list, nil
}

// RealPath converts the given <path> to its absolute path
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
	p, _ := filepath.Abs(os.Args[0])
	return p
}

// SelfName returns file name of current running process(binary).
func SelfName() string {
	return Basename(SelfPath())
}

// SelfDir returns absolute directory path of current running process(binary).
func SelfDir() string {
	return filepath.Dir(SelfPath())
}

// Basename returns the last element of path.
// Trailing path separators are removed before extracting the last element.
// If the path is empty, Base returns ".".
// If the path consists entirely of separators, Base returns a single separator.
func Basename(path string) string {
	return filepath.Base(path)
}

// Dir returns all but the last element of path, typically the path's directory.
// After dropping the final element, Dir calls Clean on the path and trailing
// slashes are removed.
// If the path is empty, Dir returns ".".
// If the path consists entirely of separators, Dir returns a single separator.
// The returned path does not end in a separator unless it is the root directory.
func Dir(path string) string {
	return filepath.Dir(path)
}

// Ext returns the file name extension used by path.
// The extension is the suffix beginning at the final dot
// in the final element of path; it is empty if there is
// no dot.
//
// Note: the result contains symbol '.'.
func Ext(path string) string {
	return filepath.Ext(path)
}

// Home returns absolute path of current user's home directory.
func Home() (string, error) {
	u, err := user.Current()
	if nil == err {
		return u.HomeDir, nil
	}
	if "windows" == runtime.GOOS {
		return homeWindows()
	}
	return homeUnix()
}

func homeUnix() (string, error) {
	if home := os.Getenv("HOME"); home != "" {
		return home, nil
	}
	var stdout bytes.Buffer
	cmd := exec.Command("sh", "-c", "eval echo ~$USER")
	cmd.Stdout = &stdout
	if err := cmd.Run(); err != nil {
		return "", err
	}

	result := strings.TrimSpace(stdout.String())
	if result == "" {
		return "", errors.New("blank output when reading home directory")
	}

	return result, nil
}

func homeWindows() (string, error) {
	drive := os.Getenv("HOMEDRIVE")
	path := os.Getenv("HOMEPATH")
	home := drive + path
	if drive == "" || path == "" {
		home = os.Getenv("USERPROFILE")
	}
	if home == "" {
		return "", errors.New("HOMEDRIVE, HOMEPATH, and USERPROFILE are blank")
	}

	return home, nil
}

// MainPkgPath returns absolute file path of package main,
// which contains the entrance function main.
//
// It's only available in develop environment.
//
// Note1: Only valid for source development environments,
// IE only valid for systems that generate this executable.
// Note2: When the method is called for the first time, if it is in an asynchronous goroutine,
// the method may not get the main package path.
func MainPkgPath() string {
	path := mainPkgPath.Val()
	if path != "" {
		if path == "-" {
			return ""
		}
		return path
	}
	for i := 1; i < 10000; i++ {
		if _, file, _, ok := runtime.Caller(i); ok {
			// <file> is separated by '/'
			if gstr.Contains(file, "/gf/g/") {
				continue
			}
			if Ext(file) != ".go" {
				continue
			}
			// separator of <file> '/' will be converted to Separator.
			for path = Dir(file); len(path) > 1 && Exists(path) && path[len(path)-1] != os.PathSeparator; {
				files, _ := ScanDir(path, "*.go")
				for _, v := range files {
					if gregex.IsMatchString(`package\s+main`, GetContents(v)) {
						mainPkgPath.Set(path)
						return path
					}
				}
				path = Dir(path)
			}

		} else {
			break
		}
	}
	// If it fails finding the path, then mark it as "-",
	// which means it will never do this search again.
	mainPkgPath.Set("-")
	return ""
}

// See os.TempDir().
func TempDir() string {
	return os.TempDir()
}
