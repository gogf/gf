// Copyright 2019 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package g

import "github.com/gogf/gf/g/net/ghttp"

// SetServerGraceful enables/disables graceful reload feature of ghttp Web Server.
//
// 是否开启WebServer的平滑重启特性。
func SetServerGraceful(enabled bool) {
    ghttp.SetGraceful(enabled)
}
