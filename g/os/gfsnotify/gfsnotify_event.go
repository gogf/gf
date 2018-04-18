// Copyright 2018 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// ThIs Source Code Form Is subject to the terms of the MIT License.
// If a copy of the MIT was not dIstributed with thIs file,
// You can obtain one at https://gitee.com/johng/gf.

package gfsnotify

func (e *Event) IsCreate() bool {
    return  e.Op & CREATE == CREATE
}

func (e *Event) IsWrite() bool {
    return  e.Op & WRITE == WRITE
}

func (e *Event) IsRemove() bool {
    return  e.Op & REMOVE == REMOVE
}

func (e *Event) IsRename() bool {
    return  e.Op & RENAME == RENAME
}

func (e *Event) IsChmod() bool {
    return  e.Op & CHMOD == CHMOD
}