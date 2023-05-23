// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gfile

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
)

// Copy file/directory from `src` to `dst`.
//
// If `src` is file, it calls CopyFile to implements copy feature,
// or else it calls CopyDir.
func Copy(src string, dst string) error {
	if src == "" {
		return gerror.NewCode(gcode.CodeInvalidParameter, "source path cannot be empty")
	}
	if dst == "" {
		return gerror.NewCode(gcode.CodeInvalidParameter, "destination path cannot be empty")
	}
	if IsFile(src) {
		return CopyFile(src, dst)
	}
	return CopyDir(src, dst)
}

// CopyFile copies the contents of the file named `src` to the file named
// by `dst`. The file will be created if it does not exist. If the
// destination file exists, all it's contents will be replaced by the contents
// of the source file. The file mode will be copied from the source and
// the copied data is synced/flushed to stable storage.
// Thanks: https://gist.github.com/r0l1/92462b38df26839a3ca324697c8cba04
func CopyFile(src, dst string) (err error) {
	if src == "" {
		return gerror.NewCode(gcode.CodeInvalidParameter, "source file cannot be empty")
	}
	if dst == "" {
		return gerror.NewCode(gcode.CodeInvalidParameter, "destination file cannot be empty")
	}
	// If src and dst are the same path, it does nothing.
	if src == dst {
		return nil
	}
	var inFile *os.File
	inFile, err = Open(src)
	if err != nil {
		return
	}
	defer func() {
		if e := inFile.Close(); e != nil {
			err = gerror.Wrapf(e, `file close failed for "%s"`, src)
		}
	}()
	var outFile *os.File
	outFile, err = Create(dst)
	if err != nil {
		return
	}
	defer func() {
		if e := outFile.Close(); e != nil {
			err = gerror.Wrapf(e, `file close failed for "%s"`, dst)
		}
	}()
	if _, err = io.Copy(outFile, inFile); err != nil {
		err = gerror.Wrapf(err, `io.Copy failed from "%s" to "%s"`, src, dst)
		return
	}
	if err = outFile.Sync(); err != nil {
		err = gerror.Wrapf(err, `file sync failed for file "%s"`, dst)
		return
	}
	if err = Chmod(dst, DefaultPermCopy); err != nil {
		return
	}
	return
}

// CopyDir recursively copies a directory tree, attempting to preserve permissions.
//
// Note that, the Source directory must exist and symlinks are ignored and skipped.
func CopyDir(src string, dst string) (err error) {
	if src == "" {
		return gerror.NewCode(gcode.CodeInvalidParameter, "source directory cannot be empty")
	}
	if dst == "" {
		return gerror.NewCode(gcode.CodeInvalidParameter, "destination directory cannot be empty")
	}
	// If src and dst are the same path, it does nothing.
	if src == dst {
		return nil
	}
	src = filepath.Clean(src)
	dst = filepath.Clean(dst)
	si, err := Stat(src)
	if err != nil {
		return err
	}
	if !si.IsDir() {
		return gerror.NewCode(gcode.CodeInvalidParameter, "source is not a directory")
	}
	if !Exists(dst) {
		if err = os.MkdirAll(dst, DefaultPermCopy); err != nil {
			err = gerror.Wrapf(err, `create directory failed for path "%s", perm "%s"`, dst, DefaultPermCopy)
			return
		}
	}
	entries, err := ioutil.ReadDir(src)
	if err != nil {
		err = gerror.Wrapf(err, `read directory failed for path "%s"`, src)
		return
	}
	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())
		if entry.IsDir() {
			if err = CopyDir(srcPath, dstPath); err != nil {
				return
			}
		} else {
			// Skip symlinks.
			if entry.Mode()&os.ModeSymlink != 0 {
				continue
			}
			if err = CopyFile(srcPath, dstPath); err != nil {
				return
			}
		}
	}
	return
}
