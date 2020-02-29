// Copyright 2019 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package ghttp

import (
	"errors"
	"github.com/gogf/gf/os/gfile"
	"github.com/gogf/gf/os/gtime"
	"github.com/gogf/gf/util/grand"
	"io"
	"mime/multipart"
	"os"
	"strconv"
	"strings"
)

// UploadFile wraps the multipart uploading file with more and convenient features.
type UploadFile struct {
	*multipart.FileHeader
}

// UploadFiles is array type for *UploadFile.
type UploadFiles []*UploadFile

// Save saves the single uploading file to specified path.
// The parameter path can be either a directory or a file path. If <path> is a directory,
// it saves the uploading file to the directory using its original name. If <path> is a
// file path, it saves the uploading file to the file path.
//
// The parameter <randomlyRename> specifies whether randomly renames the file name, which
// make sense if the <path> is a directory.
//
// Note that it will overwrite the target file if there's already a same name file exist.
func (f *UploadFile) Save(path string, randomlyRename ...bool) error {
	if f == nil {
		return nil
	}
	file, err := f.Open()
	if err != nil {
		return err
	}
	defer file.Close()

	var newFile *os.File
	if gfile.IsDir(path) {
		filename := gfile.Basename(f.Filename)
		if len(randomlyRename) > 0 && randomlyRename[0] {
			filename = strings.ToLower(strconv.FormatInt(gtime.TimestampNano(), 36) + grand.S(6))
			filename = filename + gfile.Ext(f.Filename)
		}
		newFile, err = gfile.Create(gfile.Join(path, filename))
	} else {
		newFile, err = gfile.Create(path)
	}
	if err != nil {
		return err
	}
	defer newFile.Close()

	if _, err := io.Copy(newFile, file); err != nil {
		return err
	}
	return nil
}

// Save saves all uploading files to specified directory path.
//
// The parameter <dirPath> should be a directory path or it returns error.
//
// The parameter <randomlyRename> specifies whether randomly renames all the file names.
func (fs UploadFiles) Save(dirPath string, randomlyRename ...bool) error {
	if len(fs) == 0 {
		return nil
	}
	if !gfile.IsDir(dirPath) {
		return errors.New(`parameter "dirPath" should be a directory path`)
	}
	var err error
	for _, f := range fs {
		if err = f.Save(dirPath, randomlyRename...); err != nil {
			return err
		}
	}
	return nil
}

// GetUploadFile retrieves and returns the uploading file with specified form name.
// This function is used for retrieving single uploading file object, which is
// uploaded using multipart form content type.
//
// Note that the <name> is the file field name of the multipart form from client.
func (r *Request) GetUploadFile(name string) *UploadFile {
	uploadFiles := r.GetUploadFiles(name)
	if len(uploadFiles) > 0 {
		return uploadFiles[0]
	}
	return nil
}

// GetUploadFiles retrieves and returns multiple uploading files with specified form name.
// This function is used for retrieving multiple uploading file objects, which are
// uploaded using multipart form content type.
//
// Note that the <name> is the file field name of the multipart form from client.
func (r *Request) GetUploadFiles(name string) UploadFiles {
	multipartFiles := r.GetMultipartFiles(name)
	if len(multipartFiles) > 0 {
		uploadFiles := make(UploadFiles, len(multipartFiles))
		for k, v := range multipartFiles {
			uploadFiles[k] = &UploadFile{
				FileHeader: v,
			}
		}
		return uploadFiles
	}
	return nil
}
