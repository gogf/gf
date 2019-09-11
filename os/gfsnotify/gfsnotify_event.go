// Copyright 2018 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// ThIs Source Code Form Is subject to the terms of the MIT License.
// If a copy of the MIT was not dIstributed with thIs file,
// You can obtain one at https://github.com/gogf/gf.

package gfsnotify

func (e *Event) String() string {
	return e.event.String()
}

// 文件/目录创建
func (e *Event) IsCreate() bool {
	return e.Op == 1 || e.Op&CREATE == CREATE
}

// 文件/目录修改
func (e *Event) IsWrite() bool {
	return e.Op&WRITE == WRITE
}

// 文件/目录删除
func (e *Event) IsRemove() bool {
	return e.Op&REMOVE == REMOVE
}

// 文件/目录重命名
func (e *Event) IsRename() bool {
	return e.Op&RENAME == RENAME
}

// 文件/目录修改权限
func (e *Event) IsChmod() bool {
	return e.Op&CHMOD == CHMOD
}
