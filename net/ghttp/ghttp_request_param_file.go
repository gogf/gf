// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package ghttp

import (
	"context"
	"io"
	"mime/multipart"
	"strconv"
	"strings"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/internal/intlog"
	"github.com/gogf/gf/v2/internal/json"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/gogf/gf/v2/util/grand"
)

// UploadFile wraps the multipart uploading file with more and convenient features.
type UploadFile struct {
	*multipart.FileHeader `json:"-"`
	ctx                   context.Context
}

// MarshalJSON implements the interface MarshalJSON for json.Marshal.
func (f UploadFile) MarshalJSON() ([]byte, error) {
	return json.Marshal(gconv.Map(f))
}

// UploadFiles is an array type of *UploadFile.
type UploadFiles []*UploadFile

// Save saves the single uploading file to directory path and returns the saved file name.
//
// The parameter `dirPath` should be a directory path, or it returns error.
//
// Note that it will OVERWRITE the target file if there's already a same name file exist.
func (f *UploadFile) Save(dirPath string, randomlyRename ...bool) (filename string, err error) {
	if f == nil {
		return "", gerror.NewCode(
			gcode.CodeMissingParameter,
			"file is empty, maybe you retrieve it from invalid field name or form enctype",
		)
	}
	if !gfile.Exists(dirPath) {
		if err = gfile.Mkdir(dirPath); err != nil {
			return
		}
	} else if !gfile.IsDir(dirPath) {
		return "", gerror.NewCode(gcode.CodeInvalidParameter, `parameter "dirPath" should be a directory path`)
	}

	file, err := f.Open()
	if err != nil {
		err = gerror.Wrapf(err, `UploadFile.Open failed`)
		return "", err
	}
	defer file.Close()

	name := gfile.Basename(f.Filename)
	if len(randomlyRename) > 0 && randomlyRename[0] {
		name = strings.ToLower(strconv.FormatInt(gtime.TimestampNano(), 36) + grand.S(6))
		name = name + gfile.Ext(f.Filename)
	}
	filePath := gfile.Join(dirPath, name)
	newFile, err := gfile.Create(filePath)
	if err != nil {
		return "", err
	}
	defer newFile.Close()
	intlog.Printf(f.ctx, `save upload file: %s`, filePath)
	if _, err = io.Copy(newFile, file); err != nil {
		err = gerror.Wrapf(err, `io.Copy failed from "%s" to "%s"`, f.Filename, filePath)
		return "", err
	}
	return gfile.Basename(filePath), nil
}

// Save saves all uploading files to specified directory path and returns the saved file names.
//
// The parameter `dirPath` should be a directory path or it returns error.
//
// The parameter `randomlyRename` specifies whether randomly renames all the file names.
func (fs UploadFiles) Save(dirPath string, randomlyRename ...bool) (filenames []string, err error) {
	if len(fs) == 0 {
		return nil, gerror.NewCode(
			gcode.CodeMissingParameter,
			"file array is empty, maybe you retrieve it from invalid field name or form enctype",
		)
	}
	for _, f := range fs {
		if filename, err := f.Save(dirPath, randomlyRename...); err != nil {
			return filenames, err
		} else {
			filenames = append(filenames, filename)
		}
	}
	return
}

// GetUploadFile retrieves and returns the uploading file with specified form name.
// This function is used for retrieving single uploading file object, which is
// uploaded using multipart form content type.
//
// It returns nil if retrieving failed or no form file with given name posted.
//
// Note that the `name` is the file field name of the multipart form from client.
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
// It returns nil if retrieving failed or no form file with given name posted.
//
// Note that the `name` is the file field name of the multipart form from client.
func (r *Request) GetUploadFiles(name string) UploadFiles {
	multipartFiles := r.GetMultipartFiles(name)
	if len(multipartFiles) > 0 {
		uploadFiles := make(UploadFiles, len(multipartFiles))
		for k, v := range multipartFiles {
			uploadFiles[k] = &UploadFile{
				ctx:        r.Context(),
				FileHeader: v,
			}
		}
		return uploadFiles
	}
	return nil
}
