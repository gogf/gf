// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/wangyougui/gf.

package httpclient

import (
	"github.com/wangyougui/gf/v2/net/gclient"
	"github.com/wangyougui/gf/v2/os/glog"
)

// Config is the configuration struct for SDK client.
type Config struct {
	URL     string          `v:"required"` // Service address. Eg: user.svc.local, http://user.svc.local
	Client  *gclient.Client // Custom underlying client.
	Logger  *glog.Logger    // Custom logger.
	RawDump bool            // Whether auto dump request&response in stdout.
}
