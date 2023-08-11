package cmd

import (
	"fmt"
	"testing"

	"github.com/gogf/gf/v2/test/gtest"
)

func Test_Fix_doFixV25Content(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			content = gtest.DataContent(`fix25_content.go.txt`)
			f       = cFix{}
		)
		newContent, err := f.doFixV25Content(content)
		t.AssertNil(err)
		fmt.Println(newContent)
	})
}
