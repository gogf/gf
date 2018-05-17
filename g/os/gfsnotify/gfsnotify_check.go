// Copyright 2018 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

package gfsnotify

import (
    "time"
    "gitee.com/johng/gf/g/os/gtime"
    "gitee.com/johng/gf/g/os/gfile"
    "fmt"
)

const (
    gDEFAULT_ACTIVE_CHECK_INTERVAL = 500 // (毫秒)默认主动检测时间间隔，未开启时什么也不做
)

// 文件定时监听器循环
func (w *Watcher) startActiveCheckLoop() {
    go func() {
        for {
            select {
                // 关闭事件
                case <- w.closeChan:
                    return

                default:
                    if w.activeCheckInterval.Val() > 0 {
                        paths := w.watchUpdateTimeMap.Keys()
                        for _, path := range paths {
                            lastUpdateTime := w.watchUpdateTimeMap.Get(path)
                            if int(gtime.Millisecond()) - lastUpdateTime > w.activeCheckInterval.Val() {

                                fileUpdateTime := int(gfile.MTimeMillisecond(path))
                                fmt.Println("check:", path, fileUpdateTime, lastUpdateTime)
                                if fileUpdateTime > lastUpdateTime {
                                    fmt.Println("update:", path)
                                    w.watchUpdateTimeMap.Set(path, fileUpdateTime)
                                    w.events.PushBack(&Event{
                                        Path : path,
                                        Op   : Op(WRITE),
                                    })
                                }
                            }
                        }
                        time.Sleep(time.Duration(w.activeCheckInterval.Val())*time.Millisecond)
                    } else {
                        time.Sleep(gDEFAULT_ACTIVE_CHECK_INTERVAL*time.Millisecond)
                    }
            }
        }
    }()
}