package api

import (
	"github.com/gogf/gf/v2/net/gclient"
	"github.com/gogf/gf/v2/os/glog"
)

// Config is the configuration struct for SDK client.
type Config struct {
	URL     string          `v:"required"` // Service address. Eg: user.svc.local, http://user.svc.local
	Client  *gclient.Client // Custom underlying client.
	Logger  *glog.Logger    // Custom logger.
	RawDump bool            // Whether auto dump request&response in stdout.
}
