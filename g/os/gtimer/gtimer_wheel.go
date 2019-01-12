// Copyright 2019 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

package gtimer

import (
    "gitee.com/johng/gf/g/container/glist"
    "gitee.com/johng/gf/g/container/gtype"
)

// 单层时间轮
type wheel struct {
    wheels     *Timer          // 所属定时器
    level      int             // 所属分层索引号
    slots      []*glist.List   // 所有的循环任务项, 按照Slot Number进行分组
    number     int64           // Slot Number
    closed     chan struct{}   // 停止事件
    ticks      *gtype.Int64    // 当前时间轮已转动的刻度数量
    totalMs    int64           // 整个时间轮的时间长度(毫秒)=number*interval
    createMs   int64           // 创建时间(毫秒)
    intervalMs int64           // 时间间隔(slot时间长度, 毫秒)
}

// 关闭循环任务
func (w *wheel) Close() {
    w.closed <- struct{}{}
}