// Copyright 2017 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

package gtcp

import (
    "gitee.com/johng/gf/g/os/glog"
)

// 执行监听
func (s *gTcpServer) Run() {
    if s == nil || s.listener == nil {
        glog.Println("start running failed: socket address bind failed")
        return
    }
    for  {
        conn, err := s.listener.Accept()
        if err != nil {
            glog.Error(err)
        } else if conn != nil {
            go s.handler(conn)
        }
    }
}
