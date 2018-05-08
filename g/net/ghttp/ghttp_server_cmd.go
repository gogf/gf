// Copyright 2017 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

package ghttp

import (
    "fmt"
    "gitee.com/johng/gf/g/net/gtcp"
)

// 开启命令监听端口
func (s *Server) startCmdService() {
    s.BindHandler("/heartbeat", func(r *Request) {

    })
    s.BindHandler("/restart", func(r *Request) {

    })
    server := s.newGracefulServer(fmt.Sprintf("127.0.0.1:%d", s.cmdPort))
    if err := server.ListenAndServe(); err != nil {

    }
}