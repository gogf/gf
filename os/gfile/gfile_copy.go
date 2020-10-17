// Copyright 2019 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gfile

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

// Copy file/directory from <src> to <dst>.
//
// If <src> is file, it calls CopyFile to implements copy feature,
// or else it calls CopyDir.
func Copy(src string, dst string) error {
	if src == "" {
		return errors.New("source path cannot be empty")
	}
	if dst == "" {
		return errors.New("destination path cannot be empty")
	}
	if IsFile(src) {
		return CopyFile(src, dst)
	}
	return CopyDir(src, dst)
}

// CopyFile copies the contents of the file named <src> to the file named
// by <dst>. The file will be created if it does not exist. If the
// destination file exists, all it's contents will be replaced by the contents
// of the source file. The file mode will be copied from the source and
// the copied data is synced/flushed to stable storage.
// Thanks: https://gist.github.com/r0l1/92462b38df26839a3ca324697c8cba04
func CopyFile(src, dst string) (err error) {
	if src == "" {
		return errors.New("source file cannot be empty")
	}
	if dst == "" {
		return errors.New("destination file cannot be empty")
	}
	// If src and dst are the same path, it does nothing.
	if src == dst {
		return nil
	}
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
	err = os.Chmod(dst, DefaultPermCopy)
	if err != nil {
		return
	}
	return
}

// CopyDir recursively copies a directory tree, attempting to preserve permissions.
//
// Note that, the Source directory must exist and symlinks are ignored and skipped.
func CopyDir(src string, dst string) (err error) {
	if src == "" {
		return errors.New("source directory cannot be empty")
	}
	if dst == "" {
		return errors.New("destination directory cannot be empty")
	}
	// If src and dst are the same path, it does nothing.
	if src == dst {
		return nil
	}
	src = filepath.Clean(src)
	dst = filepath.Clean(dst)
	si, err := os.Stat(src)
	if err != nil {
		return err
	}
	if !si.IsDir() {
		return fmt.Errorf("source is not a directory")
	}
	if !Exists(dst) {
		err = os.MkdirAll(dst, DefaultPermCopy)
		if err != nil {
			return
		}
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
