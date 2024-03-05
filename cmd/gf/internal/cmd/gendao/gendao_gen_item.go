// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gendao

type (
	CGenDaoInternalGenItems struct {
		index int
		Items []CGenDaoInternalGenItem
	}
	CGenDaoInternalGenItem struct {
		Clear              bool
		StorageDirPaths    []string
		GeneratedFilePaths []string
	}
)

func newCGenDaoInternalGenItems() *CGenDaoInternalGenItems {
	return &CGenDaoInternalGenItems{
		index: -1,
		Items: make([]CGenDaoInternalGenItem, 0),
	}
}

func (i *CGenDaoInternalGenItems) Scale() {
	i.Items = append(i.Items, CGenDaoInternalGenItem{
		StorageDirPaths:    make([]string, 0),
		GeneratedFilePaths: make([]string, 0),
		Clear:              false,
	})
	i.index++
}

func (i *CGenDaoInternalGenItems) SetClear(clear bool) {
	i.Items[i.index].Clear = clear
}

func (i CGenDaoInternalGenItems) AppendDirPath(storageDirPath string) {
	i.Items[i.index].StorageDirPaths = append(
		i.Items[i.index].StorageDirPaths,
		storageDirPath,
	)
}

func (i CGenDaoInternalGenItems) AppendGeneratedFilePath(generatedFilePath string) {
	i.Items[i.index].GeneratedFilePaths = append(
		i.Items[i.index].GeneratedFilePaths,
		generatedFilePath,
	)
}
