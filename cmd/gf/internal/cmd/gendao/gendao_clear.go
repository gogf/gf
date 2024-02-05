// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gendao

import (
	"context"

	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/text/gstr"

	"github.com/gogf/gf/cmd/gf/v2/internal/utility/mlog"
)

func doClear(ctx context.Context, in CGenDaoInput) {
	filePaths, err := gfile.ScanDirFile(in.Path, "*.go", true)
	if err != nil {
		mlog.Fatal(err)
	}
	var allGeneratedFilePaths = make([]string, 0)
	allGeneratedFilePaths = append(allGeneratedFilePaths, in.generatedFilePaths.DaoFilePaths...)
	allGeneratedFilePaths = append(allGeneratedFilePaths, in.generatedFilePaths.DaoInternalFilePaths...)
	allGeneratedFilePaths = append(allGeneratedFilePaths, in.generatedFilePaths.EntityFilePaths...)
	allGeneratedFilePaths = append(allGeneratedFilePaths, in.generatedFilePaths.DoFilePaths...)
	for _, filePath := range filePaths {
		if !gstr.InArray(allGeneratedFilePaths, filePath) {
			if err = gfile.Remove(filePath); err != nil {
				mlog.Print(err)
			}
		}
	}
}
