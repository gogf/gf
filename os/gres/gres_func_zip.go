// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gres

import (
	"archive/zip"
	"io"
	"os"
	"strings"
	"time"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/internal/fileinfo"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/text/gregex"
)

// ZipPathWriter compresses `paths` to `writer` using zip compressing algorithm.
// The unnecessary parameter `prefix` indicates the path prefix for zip file.
//
// Note that the parameter `paths` can be either a directory or a file, which
// supports multiple paths join with ','.
func zipPathWriter(paths string, writer io.Writer, option ...Option) error {
	zipWriter := zip.NewWriter(writer)
	defer zipWriter.Close()
	for _, path := range strings.Split(paths, ",") {
		path = strings.TrimSpace(path)
		if err := doZipPathWriter(path, zipWriter, option...); err != nil {
			return err
		}
	}
	return nil
}

// doZipPathWriter compresses the file of given `path` and writes the content to `zipWriter`.
// The parameter `exclude` specifies the exclusive file path that is not compressed to `zipWriter`,
// commonly the destination zip file path.
// The unnecessary parameter `prefix` indicates the path prefix for zip file.
func doZipPathWriter(srcPath string, zipWriter *zip.Writer, option ...Option) error {
	var (
		err          error
		files        []string
		usedOption   Option
		absolutePath string
	)
	if len(option) > 0 {
		usedOption = option[0]
	}
	absolutePath, err = gfile.Search(srcPath)
	if err != nil {
		return err
	}
	if gfile.IsDir(absolutePath) {
		files, err = gfile.ScanDir(absolutePath, "*", true)
		if err != nil {
			return err
		}
	} else {
		files = []string{absolutePath}
	}
	headerPrefix := strings.TrimRight(usedOption.Prefix, `\/`)
	if headerPrefix != "" && gfile.IsDir(absolutePath) {
		headerPrefix += "/"
	}

	if headerPrefix == "" {
		if usedOption.KeepPath {
			// It keeps the path from file system to zip info in resource manager.
			// Usually for relative path, it makes little sense for absolute path.
			headerPrefix = srcPath
		} else {
			headerPrefix = gfile.Basename(absolutePath)
		}
	}
	headerPrefix = strings.ReplaceAll(headerPrefix, `//`, `/`)
	for _, file := range files {
		// It here calculates the file name prefix, especially packing the directory.
		// Eg:
		// path: dir1
		// file: dir1/dir2/file
		// file[len(absolutePath):] => /dir2/file
		// gfile.Dir(subFilePath)   => /dir2
		var subFilePath string
		// Normal handling: remove the `absolutePath`(source directory path) for file.
		subFilePath = file[len(absolutePath):]
		if subFilePath != "" {
			subFilePath = gfile.Dir(subFilePath)
		}
		if err = zipFile(file, headerPrefix+subFilePath, zipWriter); err != nil {
			return err
		}
	}
	// Add all directories to zip archive.
	if headerPrefix != "" {
		var (
			name    string
			tmpPath = headerPrefix
		)
		for {
			name = strings.ReplaceAll(gfile.Basename(tmpPath), `\`, `/`)
			err = zipFileVirtual(fileinfo.New(name, 0, os.ModeDir|os.ModePerm, time.Now()), tmpPath, zipWriter)
			if err != nil {
				return err
			}
			if tmpPath == `/` || !strings.Contains(tmpPath, `/`) {
				break
			}
			tmpPath = gfile.Dir(tmpPath)
		}
	}
	return nil
}

// zipFile compresses the file of given `path` and writes the content to `zw`.
// The parameter `prefix` indicates the path prefix for zip file.
func zipFile(path string, prefix string, zw *zip.Writer) error {
	prefix = strings.ReplaceAll(prefix, `//`, `/`)

	file, err := os.Open(path)
	if err != nil {
		err = gerror.Wrapf(err, `os.Open failed for path "%s"`, path)
		return nil
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		err = gerror.Wrapf(err, `read file stat failed for path "%s"`, path)
		return err
	}

	header, err := createFileHeader(info, prefix)
	if err != nil {
		return err
	}
	if !info.IsDir() {
		// Default compression level.
		header.Method = zip.Deflate
	}
	// Zip header containing the info of a zip file.
	writer, err := zw.CreateHeader(header)
	if err != nil {
		err = gerror.Wrapf(err, `create zip header failed for %#v`, header)
		return err
	}
	if !info.IsDir() {
		if _, err = io.Copy(writer, file); err != nil {
			err = gerror.Wrapf(err, `io.Copy failed for file "%s"`, path)
			return err
		}
	}
	return nil
}

func zipFileVirtual(info os.FileInfo, path string, zw *zip.Writer) error {
	header, err := createFileHeader(info, "")
	if err != nil {
		return err
	}
	header.Name = path
	if _, err = zw.CreateHeader(header); err != nil {
		err = gerror.Wrapf(err, `create zip header failed for %#v`, header)
		return err
	}
	return nil
}

func createFileHeader(info os.FileInfo, prefix string) (*zip.FileHeader, error) {
	header, err := zip.FileInfoHeader(info)
	if err != nil {
		err = gerror.Wrapf(err, `create file header failed for name "%s"`, info.Name())
		return nil, err
	}
	if len(prefix) > 0 {
		header.Name = prefix + `/` + header.Name
		header.Name = strings.ReplaceAll(header.Name, `\`, `/`)
		header.Name, _ = gregex.ReplaceString(`/{2,}`, `/`, header.Name)
	}
	return header, nil
}
