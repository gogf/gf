// Copyright 2017 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

package gudp

import "gitee.com/johng/gf/g/os/glog"

// 执行监听
func (s *Server) Run() {
    if s == nil || s.listener == nil {
        glog.Println("start running failed: socket address bind failed")
        return
    }
    if s.handler == nil {
        glog.Println("start running failed: socket handler not defined")
        return
    }
    for {
        s.handler(s.listener)
    }

}
