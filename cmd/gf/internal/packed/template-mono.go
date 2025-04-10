package packed

import "github.com/gogf/gf/v3/os/gres"

func init() {
	if err := gres.Add("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"); err != nil {
		panic("add binary content to resource manager failed: " + err.Error())
	}
}
