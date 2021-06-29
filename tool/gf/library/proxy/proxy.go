package proxy

import (
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/os/genv"
	"github.com/gogf/gf/tool/gf/library/mlog"
	"time"
)

var (
	httpClient = ghttp.NewClient()
)

func init() {
	httpClient.SetTimeout(time.Second)
}

// AutoSet automatically checks and sets the golang proxy.
func AutoSet() {
	SetGoModuleEnabled(true)
	genv.Set("GOPROXY", "https://goproxy.cn")
}

// SetGoModuleEnabled enables/disables the go module feature.
func SetGoModuleEnabled(enabled bool) {
	if enabled {
		mlog.Debug("set GO111MODULE=on")
		genv.Set("GO111MODULE", "on")
	} else {
		mlog.Debug("set GO111MODULE=off")
		genv.Set("GO111MODULE", "off")
	}
}
