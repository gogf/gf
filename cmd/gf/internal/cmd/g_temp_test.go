package cmd

import (
	"testing"

	"github.com/gogf/gf/cmd/gf/v2/internal/cmd/genservice"
	"github.com/gogf/gf/v2/test/gtest"
)

func TestTemp(t *testing.T) {
	path := gtest.DataPath("genservice", "logic", "article", "article_extra.go")
	_, err := genservice.CGenService{}.CalculateInterfaceFunctions2(path)
	if err != nil {
		panic(err)
	}
}
