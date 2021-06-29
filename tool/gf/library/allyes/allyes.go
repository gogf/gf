package allyes

import (
	"github.com/gogf/gf/os/gcmd"
	"github.com/gogf/gf/os/genv"
)

const (
	EnvName = "GF_CLI_ALL_YES"
)

// Init initializes the package manually.
func Init() {
	if gcmd.ContainsOpt("y") {
		genv.Set(EnvName, "1")
	}
}

// Check checks whether option allow all yes for command.
func Check() bool {
	return genv.Get(EnvName) == "1"
}
