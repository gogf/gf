package cmd

import (
	"testing"

	"github.com/gogf/gf/cmd/gf/v2/internal/cmd/genservice"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/gutil"
)

func TestTemp(t *testing.T) {
	path := gtest.DataPath("genservice", "logic", "article", "article_extra.go")
	pkgs, items, err := genservice.CGenService{}.CalculateItemsInSrc(path)
	if err != nil {
		panic(err)
	}
	gutil.Dump(pkgs)
	gutil.Dump(items)
}
