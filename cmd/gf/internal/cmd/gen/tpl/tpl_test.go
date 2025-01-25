package tpl_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/gogf/gf/cmd/gf/v2/internal/cmd/gen/tpl"
)

func TestTpl(t *testing.T) {
	c := tpl.CGenTpl{}
	t.Log(c)
	out, err := c.Tpl(context.Background(), tpl.CGenTplInput{
		Path:         "./output",
		TplPath:      "./testdata",
		Link:         fmt.Sprintf("mysql:root:%s@tcp(127.0.0.1:3306)/focus?loc=Local&parseTime=true", "root123"),
		Tables:       "gf_user",
		ImportPrefix: "github.com/gogf/gf/cmd/gf/v2/internal/cmd/gen/tpl/output",
	})
	if err != nil {
		t.Error(err)
	}
	t.Log(out)
}
