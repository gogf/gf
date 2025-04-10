package packed

import "github.com/gogf/gf/v3/os/gres"

func init() {
	if err := gres.Add("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"); err != nil {
		panic("add binary content to resource manager failed: " + err.Error())
	}
}
