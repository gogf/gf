// Copyright 2018 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

package gproc

// 获取其他进程传递到当前进程的消息包，阻塞执行
func Receive() *Msg {
    if v := commReceiveQueue.PopFront(); v != nil {
        return v.(*Msg)
    }
    return nil
}