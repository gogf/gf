package gfile

import (
	"github.com/gogf/gf/g/test/gtest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestIsDir(t *testing.T) {

	gtest.Case(t, func() {
		paths := "/testfile"
		CreateDir(paths)
		defer DelTestFiles(paths)

		gtest.Assert(IsDir(Testpath()+paths), true)
		gtest.Assert(IsDir("./testfile2"), false)
		gtest.Assert(IsDir("./testfile/tt.txt"), false)
		gtest.Assert(IsDir(""), false)

	})

}

func TestCreate(t *testing.T) {
	gtest.Case(t, func() {
		var (
			err       error
			filepaths []string
			fileobj   *os.File
		)

		filepaths = append(filepaths, "/testfile_cc1.txt")
		filepaths = append(filepaths, "/testfile_cc2.txt")

		for _, v := range filepaths {
			fileobj, err = Create(Testpath()+v)
			defer DelTestFiles(v)
			fileobj.Close()
			gtest.Assert(err, nil)

		}

	})

}

func TestOpen(t *testing.T) {
	gtest.Case(t, func() {
		var (
			err     error
			files   []string
			flags   []bool
			fileobj *os.File
		)

		file1 := "/testfile_nc1.txt"
		CreateTestFile(file1, "")
		defer DelTestFiles(file1)

		files = append(files, file1)
		flags = append(flags, true)

		files = append(files, "./testfile/file1/c1.txt")
		flags = append(flags, false)

		for k, v := range files {
			fileobj, err = Open(Testpath() + v)
			fileobj.Close()
			if flags[k] {
				gtest.Assert(err, nil)
			} else {
				gtest.AssertNE(err, nil)
			}

		}

	})
}

func TestOpenFile(t *testing.T) {
	gtest.Case(t, func() {
		var (
			err     error
			files   []string
			flags   []bool
			fileobj *os.File
		)

		files = append(files, "./testfile/file1/nc1.txt")
		flags = append(flags, false)

		f1 := "/testfile_tt.txt"
		CreateTestFile(f1, "")
		defer DelTestFiles(f1)

		files = append(files, f1)
		flags = append(flags, true)

		for k, v := range files {
			fileobj, err = OpenFile(Testpath()+v, os.O_RDWR, 0666)
			fileobj.Close()
			if flags[k] {
				gtest.Assert(err, nil)
			} else {
				gtest.AssertNE(err, nil)
			}

		}

	})
}

func TestOpenWithFlag(t *testing.T) {
	gtest.Case(t, func() {
		var (
			err     error
			files   []string
			flags   []bool
			fileobj *os.File
		)
		file1 := "/testfile_t1.txt"
		CreateTestFile(file1, "")
		defer DelTestFiles(file1)
		files = append(files, file1)
		flags = append(flags, true)

		files = append(files, "./testfile/dirfiles/t1_no.txt")
		flags = append(flags, false)

		for k, v := range files {
			fileobj, err = OpenWithFlag(Testpath()+v, os.O_RDWR)
			fileobj.Close()
			if flags[k] {
				gtest.Assert(err, nil)
			} else {
				gtest.AssertNE(err, nil)
			}

		}

	})
}

func TestOpenWithFlagPerm(t *testing.T) {
	gtest.Case(t, func() {
		var (
			err     error
			files   []string
			flags   []bool
			fileobj *os.File
		)
		file1 := "/testfile_nc1.txt"
		CreateTestFile(file1, "")
		defer DelTestFiles(file1)
		files = append(files, file1)
		flags = append(flags, true)

		files = append(files, "./testfile/tt.txt")
		flags = append(flags, false)

		for k, v := range files {
			fileobj, err = OpenWithFlagPerm(Testpath()+v, os.O_RDWR, 666)
			fileobj.Close()
			if flags[k] {
				gtest.Assert(err, nil)
			} else {
				gtest.AssertNE(err, nil)
			}

		}

	})
}

func TestExists(t *testing.T) {

	gtest.Case(t, func() {
		var (
			flag  bool
			files []string
			flags []bool
		)

		file1 := "/testfile_GetContents.txt"
		CreateTestFile(file1, "")
		defer DelTestFiles(file1)

		files = append(files, file1)
		flags = append(flags, true)

		files = append(files, "./testfile/havefile1/tt_no.txt")
		flags = append(flags, false)

		for k, v := range files {
			flag = Exists(Testpath() + v)
			if flags[k] {
				gtest.Assert(flag, true)
			} else {
				gtest.Assert(flag, false)
			}

		}

	})
}

func TestPwd(t *testing.T) {
	gtest.Case(t, func() {
		paths, err := os.Getwd()
		gtest.Assert(err, nil)
		gtest.Assert(Pwd(), paths)

	})
}

func TestIsFile(t *testing.T) {
	gtest.Case(t, func() {
		var (
			flag  bool
			files []string
			flags []bool
		)

		file1 := "/testfile_tt.txt"
		CreateTestFile(file1, "")
		defer DelTestFiles(file1)
		files = append(files, file1)
		flags = append(flags, true)

		dir1 := "/testfiless"
		CreateDir(dir1)
		defer DelTestFiles(dir1)
		files = append(files, dir1)
		flags = append(flags, false)

		files = append(files, "./testfiledd/tt1.txt")
		flags = append(flags, false)

		for k, v := range files {
			flag = IsFile(Testpath() + v)
			if flags[k] {
				gtest.Assert(flag, true)
			} else {
				gtest.Assert(flag, false)
			}

		}

	})
}

func TestInfo(t *testing.T) {
	gtest.Case(t, func() {
		var (
			err    error
			paths  string = "/testfile_t1.txt"
			files  os.FileInfo
			files2 os.FileInfo
		)

		CreateTestFile(paths, "")
		defer DelTestFiles(paths)
		files, err = Info(Testpath() + paths)
		gtest.Assert(err, nil)

		files2, err = os.Stat(Testpath() + paths)
		gtest.Assert(err, nil)

		gtest.Assert(files, files2)

	})
}

//func TestMove(t *testing.T) {
//	gtest.Case(t, func(){
//		var(
//			paths string ="./testfile/havefile1/ttn1.txt"
//			topath string ="./testfile/havefile1/ttn2.txt"
//		)
//
//		gtest.Assert(Move(paths,topath),nil)
//
//	})
//}
//
// func TestRename(t *testing.T){
//	 gtest.Case(t, func(){
//		 var(
//
//			 paths string ="./testfile/havefile1/ttm1.txt"
//			 topath string ="./testfile/havefile1/ttm2.txt"
//
//		 )
//
//		 gtest.Assert(Rename(paths,topath),nil)
//		 gtest.Assert(IsFile(topath),true)
//
//		 gtest.AssertNE(Rename("",""),nil)
//
//
//	 })
//
//
// }

func TestCopy(t *testing.T) {
	gtest.Case(t, func() {
		var (
			paths  string = "/testfile_copyfile1.txt"
			topath string =  "/testfile_copyfile2.txt"
		)

		CreateTestFile(paths, "")
		defer DelTestFiles(paths)

		gtest.Assert(Copy(Testpath()+paths, Testpath()+topath), nil)
		defer DelTestFiles(topath)

		gtest.Assert(IsFile(Testpath()+topath), true)

		gtest.AssertNE(Copy("", ""), nil)

	})
}

func TestDirNames(t *testing.T) {
	gtest.Case(t, func() {
		var (
			paths    string = "/testdirs"
			err      error
			readlist []string
		)
		havelist := []string{
			"t1.txt",
			"t2.txt",
		}

		//=================创建测试文件
		CreateDir(paths)
		for _, v := range havelist {
			CreateTestFile(paths+"/"+v, "")
		}
		defer DelTestFiles(paths)

		readlist, err = DirNames(Testpath() + paths)

		gtest.Assert(err, nil)
		gtest.Assert(havelist, readlist)

		_, err = DirNames("")
		gtest.AssertNE(err, nil)

	})
}

func TestGlob(t *testing.T) {
	gtest.Case(t, func() {
		var (
			paths      string = "/testfiles/*.txt"
			dirpath    string = "/testfiles"
			err        error
			resultlist []string
		)

		havelist1 := []string{
			"t1.txt",
			"t2.txt",
		}

		havelist2 := []string{
			Testpath() + "/testfiles/t1.txt",
			Testpath() + "/testfiles/t2.txt",
		}

		//===============================构建测试文件
		CreateDir(dirpath)
		for _, v := range havelist1 {
			CreateTestFile(dirpath+"/"+v, "")
		}
		defer DelTestFiles(dirpath)

		resultlist, err = Glob(Testpath()+paths, true)
		gtest.Assert(err, nil)
		gtest.Assert(resultlist, havelist1)

		resultlist, err = Glob(Testpath()+paths, false)

		gtest.Assert(err, nil)
		gtest.Assert(Formatpaths(resultlist), Formatpaths(havelist2))

		_, err = Glob("", true)
		gtest.Assert(err, nil)

	})
}

func TestRemove(t *testing.T) {
	gtest.Case(t, func() {
		var (
			paths string = "/testfile_t1.txt"
		)
		CreateTestFile(paths, "")
		gtest.Assert(Remove(Testpath()+paths), nil)

		gtest.Assert(Remove(""), nil)

		defer DelTestFiles(paths)

	})
}

func TestIsReadable(t *testing.T) {
	gtest.Case(t, func() {
		var (
			paths1 string = "/testfile_GetContents.txt"
			paths2 string = "./testfile_GetContents_no.txt"
		)

		CreateTestFile(paths1, "")
		defer DelTestFiles(paths1)

		gtest.Assert(IsReadable(Testpath()+paths1), true)
		gtest.Assert(IsReadable(paths2), false)

	})
}

func TestIsWritable(t *testing.T) {
	gtest.Case(t, func() {
		var (
			paths1 string = "/testfile_GetContents.txt"
			paths2 string = "./testfile_GetContents_no.txt"
		)

		CreateTestFile(paths1, "")
		defer DelTestFiles(paths1)
		gtest.Assert(IsWritable(Testpath()+paths1), true)
		gtest.Assert(IsWritable(paths2), false)

	})
}

func TestChmod(t *testing.T) {
	gtest.Case(t, func() {
		var (
			paths1 string = "/testfile_GetContents.txt"
			paths2 string = "./testfile_GetContents_no.txt"
		)
		CreateTestFile(paths1, "")
		defer DelTestFiles(paths1)

		gtest.Assert(Chmod(Testpath()+paths1, 0777), nil)
		gtest.AssertNE(Chmod(paths2, 0777), nil)

	})
}

func TestScanDir(t *testing.T) {
	gtest.Case(t, func() {
		var (
			paths1 string = "/testfiledirs"
			files  []string
			err    error
		)

		CreateDir(paths1)
		CreateTestFile(paths1+"/t1.txt", "")
		CreateTestFile(paths1+"/t2.txt", "")
		defer DelTestFiles(paths1)

		files, err = ScanDir(Testpath()+paths1, "t*")

		result := []string{
			Testpath() + paths1 + "/t1.txt",
			Testpath() + paths1 + "/t2.txt",
		}

		gtest.Assert(err, nil)

		gtest.Assert(Formatpaths(files), Formatpaths(result))

		_, err = ScanDir("", "t*")
		gtest.AssertNE(err, nil)

	})
}

//获取绝对目录地址
func TestRealPath(t *testing.T) {
	gtest.Case(t, func() {
		var (
			paths1    string = "/testfile_files"
			readlPath string

			tempstr string
		)

		CreateDir(paths1)
		defer DelTestFiles(paths1)

		readlPath = RealPath("./")
		//readlPath = filepath.ToSlash(readlPath)

		tempstr, _ = filepath.Abs("./")
		//paths1 = tempstr + paths1
		//paths1=Formatpath(tempstr)

		gtest.Assert(readlPath, tempstr)

		gtest.Assert(RealPath("./nodirs"), "")

	})
}

//获取当前执行文件的目录
//注意：当用go test运行测试时，会产生临时的目录文件
func TestSelfPath(t *testing.T) {
	gtest.Case(t, func() {
		var (
			paths1    string
			readlPath string
			tempstr   string
		)
		readlPath = SelfPath()
		readlPath = filepath.ToSlash(readlPath)

		//
		tempstr, _ = filepath.Abs(os.Args[0])
		paths1 = filepath.ToSlash(tempstr)
		paths1 = strings.Replace(paths1, "./", "/", 1)

		gtest.Assert(readlPath, paths1)

	})
}

func TestSelfDir(t *testing.T) {
	gtest.Case(t, func() {
		var (
			paths1    string
			readlPath string
			tempstr   string
		)
		readlPath = SelfDir()

		tempstr, _ = filepath.Abs(os.Args[0])
		paths1 = filepath.Dir(tempstr)

		gtest.Assert(readlPath, paths1)

	})
}

func TestBasename(t *testing.T) {
	gtest.Case(t, func() {
		var (
			paths1    string = "/testfilerr_GetContents.txt"
			readlPath string
		)

		CreateTestFile(paths1, "")
		defer DelTestFiles(paths1)

		readlPath = Basename(Testpath() + paths1)
		gtest.Assert(readlPath, "testfilerr_GetContents.txt")

	})
}

func TestDir(t *testing.T) {
	gtest.Case(t, func() {
		var (
			paths1    string = "/testfiless"
			readlPath string
		)
		CreateDir(paths1)
		defer DelTestFiles(paths1)

		readlPath = Dir(Testpath() + paths1)

		gtest.Assert(readlPath, Testpath())

	})
}

//获取文件名
func TestExt(t *testing.T) {
	gtest.Case(t, func() {
		var (
			paths1   string = "/testfile_GetContents.txt"
			dirpath1        = "/testdirs"
		)
		CreateTestFile(paths1, "")
		defer DelTestFiles(paths1)

		CreateDir(dirpath1)
		defer DelTestFiles(dirpath1)

		gtest.Assert(Ext(Testpath()+paths1), ".txt")
		gtest.Assert(Ext(Testpath()+dirpath1), "")

	})
}

func TestTempDir(t *testing.T) {
	gtest.Case(t, func() {
		var (
			tpath string
		)

		tpath = TempDir()
		gtest.Assert(tpath, os.TempDir())

	})
}

func TestMkdir(t *testing.T) {
	gtest.Case(t, func() {
		var (
			tpath string = "/testfile/createdir"
			err   error
		)

		defer DelTestFiles(tpath)

		err = Mkdir(Testpath() + tpath)
		gtest.Assert(err, nil)

		err = Mkdir("")
		gtest.AssertNE(err, nil)

		err = Mkdir(Testpath() + tpath + "2/t1")
		gtest.Assert(err, nil)

	})
}

func TestStat(t *testing.T) {
	gtest.Case(t, func() {
		var (
			tpath1 string = "/testfile_t1.txt"
			tpath2 string = "./testfile_t1_no.txt"
			err    error
		)

		CreateTestFile(tpath1, "")
		defer DelTestFiles(tpath1)

		_, err = Stat(Testpath() + tpath1)
		gtest.Assert(err, nil)

		_, err = Stat(tpath2)
		gtest.AssertNE(err, nil)

	})
}

func TestMainPkgPath(t *testing.T) {
	gtest.Case(t, func() {
		var (
			reads string
		)

		reads = MainPkgPath()
		gtest.Assert(reads, "")

	})
}
