package packed

import "github.com/gogf/gf/v2/os/gres"

func init() {
	if err := gres.Add("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"); err != nil {
		panic("add binary content to resource manager failed: " + err.Error())
	}
}
