// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gendao

type (
	// CGenDaoInternalGenItems tracks generation state across multiple configuration entries.
	// Each configuration entry (e.g., different database links in the config array)
	// gets its own CGenDaoInternalGenItem via Scale(). The index field points to the
	// current active item.
	CGenDaoInternalGenItems struct {
		index int                      // Index of the current active generation item.
		Items []CGenDaoInternalGenItem // List of all generation items, one per config entry.
	}

	// CGenDaoInternalGenItem tracks generated files and directories for a single
	// configuration entry. Used by the Clear feature to identify and remove stale files.
	CGenDaoInternalGenItem struct {
		Clear              bool     // Whether to clear stale files for this item.
		StorageDirPaths    []string // Directories where generated files are stored (dao, do, entity, table).
		GeneratedFilePaths []string // All file paths generated in this run.
	}
)

// newCGenDaoInternalGenItems creates a new generation items tracker with an empty item list.
func newCGenDaoInternalGenItems() *CGenDaoInternalGenItems {
	return &CGenDaoInternalGenItems{
		index: -1,
		Items: make([]CGenDaoInternalGenItem, 0),
	}
}

// Scale adds a new generation item and advances the index to it.
// Must be called once per configuration entry before generating files.
func (i *CGenDaoInternalGenItems) Scale() {
	i.Items = append(i.Items, CGenDaoInternalGenItem{
		StorageDirPaths:    make([]string, 0),
		GeneratedFilePaths: make([]string, 0),
		Clear:              false,
	})
	i.index++
}

// SetClear enables or disables the clear (stale file removal) flag for the current item.
func (i *CGenDaoInternalGenItems) SetClear(clear bool) {
	i.Items[i.index].Clear = clear
}

// AppendDirPath records a directory path used for storing generated files in the current item.
func (i *CGenDaoInternalGenItems) AppendDirPath(storageDirPath string) {
	i.Items[i.index].StorageDirPaths = append(
		i.Items[i.index].StorageDirPaths,
		storageDirPath,
	)
}

// AppendGeneratedFilePath records a file path that was generated in the current item.
func (i *CGenDaoInternalGenItems) AppendGeneratedFilePath(generatedFilePath string) {
	i.Items[i.index].GeneratedFilePaths = append(
		i.Items[i.index].GeneratedFilePaths,
		generatedFilePath,
	)
}
