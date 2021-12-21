// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gcompress

import (
	"bytes"
	"compress/gzip"
	"io"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/os/gfile"
)

// Gzip compresses `data` using gzip algorithm.
// The optional parameter `level` specifies the compression level from
// 1 to 9 which means from none to the best compression.
//
// Note that it returns error if given `level` is invalid.
func Gzip(data []byte, level ...int) ([]byte, error) {
	var (
		writer *gzip.Writer
		buf    bytes.Buffer
		err    error
	)
	if len(level) > 0 {
		writer, err = gzip.NewWriterLevel(&buf, level[0])
		if err != nil {
			err = gerror.Wrap(err, `gzip.NewWriterLevel failed`)
			return nil, err
		}
	} else {
		writer = gzip.NewWriter(&buf)
	}
	if _, err = writer.Write(data); err != nil {
		err = gerror.Wrap(err, `writer.Write failed`)
		return nil, err
	}
	if err = writer.Close(); err != nil {
		err = gerror.Wrap(err, `writer.Close failed`)
		return nil, err
	}
	return buf.Bytes(), nil
}

// GzipFile compresses the file `src` to `dst` using gzip algorithm.
func GzipFile(src, dst string, level ...int) error {
	var (
		writer *gzip.Writer
		err    error
	)
	srcFile, err := gfile.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()
	dstFile, err := gfile.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	if len(level) > 0 {
		writer, err = gzip.NewWriterLevel(dstFile, level[0])
		if err != nil {
			err = gerror.Wrap(err, `gzip.NewWriterLevel failed`)
			return err
		}
	} else {
		writer = gzip.NewWriter(dstFile)
	}
	defer writer.Close()

	_, err = io.Copy(writer, srcFile)
	if err != nil {
		err = gerror.Wrap(err, `io.Copy failed`)
		return err
	}
	return nil
}

// UnGzip decompresses `data` with gzip algorithm.
func UnGzip(data []byte) ([]byte, error) {
	var buf bytes.Buffer
	reader, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		err = gerror.Wrap(err, `gzip.NewReader failed`)
		return nil, err
	}
	if _, err = io.Copy(&buf, reader); err != nil {
		err = gerror.Wrap(err, `io.Copy failed`)
		return nil, err
	}
	if err = reader.Close(); err != nil {
		err = gerror.Wrap(err, `reader.Close failed`)
		return buf.Bytes(), err
	}
	return buf.Bytes(), nil
}

// UnGzipFile decompresses file `src` to `dst` using gzip algorithm.
func UnGzipFile(src, dst string) error {
	srcFile, err := gfile.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()
	dstFile, err := gfile.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	reader, err := gzip.NewReader(srcFile)
	if err != nil {
		err = gerror.Wrap(err, `gzip.NewReader failed`)
		return err
	}
	defer reader.Close()

	if _, err = io.Copy(dstFile, reader); err != nil {
		err = gerror.Wrap(err, `io.Copy failed`)
		return err
	}
	return nil
}
