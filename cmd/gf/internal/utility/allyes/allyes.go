package allyes

import (
	"github.com/gogf/gf/v2/os/gcmd"
	"github.com/gogf/gf/v2/os/genv"
)

const (
	EnvName = "GF_CLI_ALL_YES"
)

// Init initializes the package manually.
func Init() {
	if gcmd.GetOpt("y") != nil {
		genv.MustSet(EnvName, "1")
	}
}

// Check checks whether option allow all yes for command.
func Check() bool {
	return genv.Get(EnvName).String() == "1"
}
