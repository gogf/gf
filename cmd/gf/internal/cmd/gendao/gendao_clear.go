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

// doClear performs cleanup of stale generated files across all generation items.
// It collects all generated file paths from all items, then for each item with
// Clear enabled, removes any .go files in its directories that are NOT in the
// generated file list. This ensures files for dropped/removed tables are cleaned up.
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

// doClearItem removes stale .go files for a single generation item.
// It scans all storage directories for .go files and deletes any file
// that is not in the allGeneratedFilePaths list (i.e., no longer corresponds
// to an existing database table).
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
			if err := gfile.RemoveFile(filePath); err != nil {
				mlog.Print(err)
			}
		}
	}
}
