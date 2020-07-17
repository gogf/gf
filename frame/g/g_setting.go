// Copyright 2019 gf Author(https://github.com/jin502437344/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/jin502437344/gf.

package g

import "github.com/jin502437344/gf/net/ghttp"

// SetServerGraceful enables/disables graceful reload feature of http Web Server.
// This feature is disabled in default.
// Deprecated, use configuration of ghttp.Server for controlling this feature.
func SetServerGraceful(enabled bool) {
	ghttp.SetGraceful(enabled)
}
