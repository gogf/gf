// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gendao

func newCGenDaoInternalGenItems() *CGenDaoInternalGenItems {
	return &CGenDaoInternalGenItems{
		Items: make([]CGenDaoInternalGenItem, 0),
	}
}

func (i *CGenDaoInternalGenItems) Scale() {
	i.Items = append(i.Items, CGenDaoInternalGenItem{
		StorageDirPaths:    make([]string, 0),
		GeneratedFilePaths: make([]string, 0),
		Clear:              false,
	})
}

func (i *CGenDaoInternalGenItems) SetClear(clear bool) {
	var (
		index  = 0
		length = len(i.Items)
	)
	if length > 0 {
		index = length - 1
	}
	i.Items[index].Clear = clear
}

func (i CGenDaoInternalGenItems) AppendDirPath(storageDirPath string) {
	var (
		index  = 0
		length = len(i.Items)
	)
	if length > 0 {
		index = length - 1
	}
	i.Items[index].StorageDirPaths = append(
		i.Items[index].StorageDirPaths,
		storageDirPath,
	)
}

func (i CGenDaoInternalGenItems) AppendGeneratedFilePath(generatedFilePath string) {
	var (
		index  = 0
		length = len(i.Items)
	)
	if length > 0 {
		index = length - 1
	}
	i.Items[index].GeneratedFilePaths = append(
		i.Items[index].GeneratedFilePaths,
		generatedFilePath,
	)
}
