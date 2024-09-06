package packed

import "github.com/gogf/gf/v2/os/gres"

func init() {
	if err := gres.Add("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"); err != nil {
		panic("add binary content to resource manager failed: " + err.Error())
	}
}
