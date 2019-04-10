package gfile

import (
	"github.com/gogf/gf/g/test/gtest"
	"path/filepath"
	"strings"
	"testing"
)

func TestSearch(t *testing.T) {
	gtest.Case(t, func(){
		var(
			paths1 string ="./testfile/dirfiles"
			tpath string
			tempstr string
			err error
		)

		tpath,err=Search(paths1)
		gtest.Assert(err,nil)

		tpath=filepath.ToSlash(tpath)



		tempstr,_=filepath.Abs("./")
		paths1=tempstr+paths1
		paths1=filepath.ToSlash(paths1)
		paths1=strings.Replace(paths1,"./","/",1)




		gtest.Assert(tpath,paths1)

	})
}
