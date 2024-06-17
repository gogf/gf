// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gendao

import (
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/text/gstr"

	"github.com/gogf/gf/cmd/gf/v2/internal/utility/mlog"
)

func doClear(items *CGenDaoInternalGenItems) {
	var allGeneratedFilePaths = make([]string, 0)
	for _, item := range items.Items {
		allGeneratedFilePaths = append(allGeneratedFilePaths, item.GeneratedFilePaths...)
	}
	for i, v := range allGeneratedFilePaths {
		allGeneratedFilePaths[i] = gfile.RealPath(v)
	}
	for _, item := range items.Items {
		if !item.Clear {
			continue
		}
		doClearItem(item, allGeneratedFilePaths)
	}
}

func doClearItem(item CGenDaoInternalGenItem, allGeneratedFilePaths []string) {
	var generatedFilePaths = make([]string, 0)
	for _, dirPath := range item.StorageDirPaths {
		filePaths, err := gfile.ScanDirFile(dirPath, "*.go", true)
		if err != nil {
			mlog.Fatal(err)
		}
		generatedFilePaths = append(generatedFilePaths, filePaths...)
	}
	for _, filePath := range generatedFilePaths {
		if !gstr.InArray(allGeneratedFilePaths, filePath) {
			if err := gfile.Remove(filePath); err != nil {
				mlog.Print(err)
			}
		}
	}
}
