// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package httpclient

import (
	"github.com/gogf/gf/v2/net/gclient"
	"github.com/gogf/gf/v2/os/glog"
)

// Config is the configuration struct for SDK client.
type Config struct {
	URL     string          `v:"required"` // Service address. Eg: user.svc.local, http://user.svc.local
	Client  *gclient.Client // Custom underlying client.
	Handler Handler         // Custom response handler.
	Logger  *glog.Logger    // Custom logger.
	RawDump bool            // Whether auto dump request&response in stdout.
}
