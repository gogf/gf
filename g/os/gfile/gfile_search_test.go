package gfile

import (
	"github.com/gogf/gf/g/test/gtest"
	"os"
	"path/filepath"
	"testing"
)

func TestSearch(t *testing.T) {
	gtest.Case(t, func() {
		var (
			paths1  string = "/testfiless"
			paths2  string = "./testfile/dirfiles_no"
			tpath   string
			tpath2  string
			tempstr string
			err     error
		)

		CreateDir(paths1)
		defer DelTestFiles(paths1)

		tpath, err = Search(os.TempDir() + paths1)
		gtest.Assert(err, nil)

		tpath = filepath.ToSlash(tpath)

		//==================自定义优先路径
		tpath2, err = Search(os.TempDir() + paths1)
		gtest.Assert(err, nil)
		tpath2 = filepath.ToSlash(tpath2)

		//tempstr, _ = filepath.Abs("./")
		tempstr = os.TempDir()
		paths1 = tempstr + paths1
		paths1 = filepath.ToSlash(paths1)
		//paths1 = strings.Replace(paths1, "./", "/", 1)

		gtest.Assert(tpath, paths1)

		gtest.Assert(tpath2, tpath)

		//测试当前目录
		tpath2, err = Search(os.TempDir()+paths1, "./")
		gtest.Assert(err, nil)
		tpath2 = filepath.ToSlash(tpath2)

		//测试当前目录
		tempstr, _ = filepath.Abs("./")
		tempstr = os.TempDir()
		paths1 = tempstr + paths1
		paths1 = filepath.ToSlash(paths1)

		gtest.Assert(tpath2, paths1)

		//测试目录不存在时
		_, err = Search(paths2)
		gtest.AssertNE(err, nil)

	})
}
