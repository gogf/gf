package gfile

//func TestSearch(t *testing.T) {
//	gtest.Case(t, func() {
//		var (
//			paths1  string = "/testfiless"
//			paths2  string = "./testfile/dirfiles_no"
//			tpath   string
//			tpath2  string
//			tempstr string
//			err     error
//		)
//
//		CreateDir(paths1)
//		defer DelTestFiles(paths1)
//
//		tpath, err = Search(Testpath() + paths1)
//		gtest.Assert(err, nil)
//
//		tpath = filepath.ToSlash(tpath)
//
//		//==================自定义优先路径
//		tpath2, err = Search(Testpath() + paths1)
//		gtest.Assert(err, nil)
//		tpath2 = filepath.ToSlash(tpath2)
//
//		//tempstr, _ = filepath.Abs("./")
//		tempstr = Testpath()
//		paths1 = tempstr + paths1
//		paths1 = filepath.ToSlash(paths1)
//		//paths1 = strings.Replace(paths1, "./", "/", 1)
//
//		gtest.Assert(tpath, paths1)
//
//		gtest.Assert(tpath2, tpath)
//
//		//测试当前目录
//		tpath2, err = Search(Testpath()+paths1, "./")
//		tpath2 = filepath.ToSlash(tpath2)
//		tempss:=filepath.ToSlash(Testpath()+paths1)
//		gtest.Assert(tpath2, tempss)
//
//
//		//测试当前目录
//		tempstr, _ = filepath.Abs("./")
//		tempstr = Testpath()
//		paths1 = tempstr + paths1
//		paths1 = filepath.ToSlash(paths1)
//
//		gtest.Assert(tpath2, paths1)
//
//		//测试目录不存在时
//		_, err = Search(paths2)
//		gtest.AssertNE(err, nil)
//
//	})
//}
