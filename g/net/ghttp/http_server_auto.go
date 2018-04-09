// Copyright 2018 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

package ghttp

// 开启异步队列处理循环，该异步线程与Server同生命周期
func (s *Server) startCloseQueueLoop() {
    go func() {
        for {
            if v := s.closeQueue.PopFront(); v != nil {
                r := v.(*Request)
                s.callHookHandler(r, "BeforeRequestClose")

                // 关闭当前会话的Cookie
                r.Cookie.Close()
                // 更新Session会话超时时间
                r.Session.UpdateExpire()

                s.callHookHandler(r, "AfterRequestClose")
            }
        }
    }()
}