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
			tpath2 string
			tempstr string
			err error
		)

		tpath,err=Search(paths1)
		gtest.Assert(err,nil)

		tpath=filepath.ToSlash(tpath)


		//==================自定义优先路径

		tpath2,err=Search(paths1,"./")
		gtest.Assert(err,nil)
		tpath2=filepath.ToSlash(tpath2)





		tempstr,_=filepath.Abs("./")
		paths1=tempstr+paths1
		paths1=filepath.ToSlash(paths1)
		paths1=strings.Replace(paths1,"./","/",1)

		gtest.Assert(tpath,paths1)

		gtest.Assert(tpath2,paths1)









	})
}
