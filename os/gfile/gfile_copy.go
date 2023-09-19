// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gfile

import (
	"io"
	"os"
	"path/filepath"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
)

// CopyOption is the option for Copy* functions.
type CopyOption struct {
	// Auto call file sync after source file content copied to target file.
	Sync bool

	// Preserve the mode of the original file to the target file.
	// If true, the Mode attribute will make no sense.
	PreserveMode bool

	// Destination created file mode.
	// The default file mode is DefaultPermCopy if PreserveMode is false.
	Mode os.FileMode
}

// Copy file/directory from `src` to `dst`.
//
// If `src` is file, it calls CopyFile to implements copy feature,
// or else it calls CopyDir.
//
// If `src` is file, but `dst` already exists and is a folder,
// it then creates a same name file of `src` in folder `dst`.
//
// Eg:
// Copy("/tmp/file1", "/tmp/file2") => /tmp/file1 copied to /tmp/file2
// Copy("/tmp/dir1",  "/tmp/dir2")  => /tmp/dir1  copied to /tmp/dir2
// Copy("/tmp/file1", "/tmp/dir2")  => /tmp/file1 copied to /tmp/dir2/file1
// Copy("/tmp/dir1",  "/tmp/file2") => error
func Copy(src string, dst string, option ...CopyOption) error {
	if src == "" {
		return gerror.NewCode(gcode.CodeInvalidParameter, "source path cannot be empty")
	}
	if dst == "" {
		return gerror.NewCode(gcode.CodeInvalidParameter, "destination path cannot be empty")
	}
	srcStat, srcStatErr := os.Stat(src)
	if srcStatErr != nil {
		if os.IsNotExist(srcStatErr) {
			return gerror.WrapCodef(
				gcode.CodeInvalidParameter,
				srcStatErr,
				`the src path "%s" does not exist`,
				src,
			)
		}
		return gerror.WrapCodef(
			gcode.CodeInternalError, srcStatErr, `call os.Stat on "%s" failed`, src,
		)
	}
	dstStat, dstStatErr := os.Stat(dst)
	if dstStatErr != nil && !os.IsNotExist(dstStatErr) {
		return gerror.WrapCodef(
			gcode.CodeInternalError, dstStatErr, `call os.Stat on "%s" failed`, dst)
	}

	if IsFile(src) {
		var isDstExist = false
		if dstStat != nil && !os.IsNotExist(dstStatErr) {
			isDstExist = true
		}
		if isDstExist && dstStat.IsDir() {
			var (
				srcName = Basename(src)
				dstPath = Join(dst, srcName)
			)
			return CopyFile(src, dstPath, option...)
		}
		return CopyFile(src, dst, option...)
	}
	if !srcStat.IsDir() && dstStat != nil && dstStat.IsDir() {
		return gerror.NewCodef(
			gcode.CodeInvalidParameter,
			`Copy failed: the src path "%s" is file, but the dst path "%s" is folder`,
			src, dst,
		)
	}
	return CopyDir(src, dst, option...)
}

// CopyFile copies the contents of the file named `src` to the file named
// by `dst`. The file will be created if it does not exist. If the
// destination file exists, all it's contents will be replaced by the contents
// of the source file. The file mode will be copied from the source and
// the copied data is synced/flushed to stable storage.
// Thanks: https://gist.github.com/r0l1/92462b38df26839a3ca324697c8cba04
func CopyFile(src, dst string, option ...CopyOption) (err error) {
	var usedOption = getCopyOption(option...)
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
	// file state check.
	srcStat, srcStatErr := os.Stat(src)
	if srcStatErr != nil {
		if os.IsNotExist(srcStatErr) {
			return gerror.WrapCodef(
				gcode.CodeInvalidParameter,
				srcStatErr,
				`the src path "%s" does not exist`,
				src,
			)
		}
		return gerror.WrapCodef(
			gcode.CodeInternalError, srcStatErr, `call os.Stat on "%s" failed`, src,
		)
	}
	dstStat, dstStatErr := os.Stat(dst)
	if dstStatErr != nil && !os.IsNotExist(dstStatErr) {
		return gerror.WrapCodef(
			gcode.CodeInternalError, dstStatErr, `call os.Stat on "%s" failed`, dst,
		)
	}
	if !srcStat.IsDir() && dstStat != nil && dstStat.IsDir() {
		return gerror.NewCodef(
			gcode.CodeInvalidParameter,
			`CopyFile failed: the src path "%s" is file, but the dst path "%s" is folder`,
			src, dst,
		)
	}
	// copy file logic.
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
	if usedOption.Sync {
		if err = outFile.Sync(); err != nil {
			err = gerror.Wrapf(err, `file sync failed for file "%s"`, dst)
			return
		}
	}
	if usedOption.PreserveMode {
		usedOption.Mode = srcStat.Mode().Perm()
	}
	if err = Chmod(dst, usedOption.Mode); err != nil {
		return
	}
	return
}

// CopyDir recursively copies a directory tree, attempting to preserve permissions.
//
// Note that, the Source directory must exist and symlinks are ignored and skipped.
func CopyDir(src string, dst string, option ...CopyOption) (err error) {
	var usedOption = getCopyOption(option...)
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
	if usedOption.PreserveMode {
		usedOption.Mode = si.Mode().Perm()
	}
	if !Exists(dst) {
		if err = os.MkdirAll(dst, usedOption.Mode); err != nil {
			err = gerror.Wrapf(
				err,
				`create directory failed for path "%s", perm "%s"`,
				dst,
				usedOption.Mode,
			)
			return
		}
	}
	entries, err := os.ReadDir(src)
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
			if entry.Type()&os.ModeSymlink != 0 {
				continue
			}
			if err = CopyFile(srcPath, dstPath, option...); err != nil {
				return
			}
		}
	}
	return
}

func getCopyOption(option ...CopyOption) CopyOption {
	var usedOption CopyOption
	if len(option) > 0 {
		usedOption = option[0]
	}
	if usedOption.Mode == 0 {
		usedOption.Mode = DefaultPermCopy
	}
	return usedOption
}
