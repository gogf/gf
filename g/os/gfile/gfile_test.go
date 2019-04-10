package gfile

import (
	"github.com/gogf/gf/g/test/gtest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)



func TestIsDir(t *testing.T){


	gtest.Case(t, func() {
		gtest.Assert(IsDir("./testfile"), true)
		gtest.Assert(IsDir("./testfile2"), false)
		gtest.Assert(IsDir("./testfile/tt.txt"), false)
	})

}

func TestCreate(t *testing.T){
	gtest.Case(t, func() {
		var (
			err error
			filepaths []string
		)

		filepaths=append(filepaths,"./testfile/file/c1.txt")
		filepaths=append(filepaths,"./testfile/file1/c2.txt")


		 for _,v:=range filepaths{
			 _,err=Create(v)
			 gtest.Assert(err,nil)

		 }


	})



}

func TestOpen(t *testing.T)  {
	gtest.Case(t, func(){
		var(
			err error
			files []string
			flags []bool
		)

		files=append(files,"./testfile/file1/nc1.txt")
		flags=append(flags,false)

		files=append(files,"./testfile/tt.txt")
		flags=append(flags,true)


		for k,v:=range files{
			_,err=Open(v)


			if flags[k]{
				gtest.Assert(err,nil)
			}else{
				gtest.AssertNE(err,nil)
			}

		}


	})
}


func TestOpenFile(t *testing.T)  {
	gtest.Case(t, func(){
		var(
			err error
			files []string
			flags []bool
		)

		files=append(files,"./testfile/file1/nc1.txt")
		flags=append(flags,false)

		files=append(files,"./testfile/tt.txt")
		flags=append(flags,true)


		for k,v:=range files{
			_,err=OpenFile(v,os.O_RDWR,0666)
			if flags[k]{
				gtest.Assert(err,nil)
			}else{
				gtest.AssertNE(err,nil)
			}

		}


	})
}



func TestOpenWithFlag(t *testing.T) {
	gtest.Case(t, func(){
		var(
			err error
			files []string
			flags []bool
		)

		files=append(files,"./testfile/file1/nc1.txt")
		flags=append(flags,false)

		files=append(files,"./testfile/tt.txt")
		flags=append(flags,true)


		for k,v:=range files{
			_,err=OpenWithFlag(v,os.O_RDWR)
			if flags[k]{
				gtest.Assert(err,nil)
			}else{
				gtest.AssertNE(err,nil)
			}

		}


	})
}


func TestOpenWithFlagPerm(t *testing.T) {
	gtest.Case(t, func(){
		var(
			err error
			files []string
			flags []bool
		)

		files=append(files,"./testfile/file1/nc1.txt")
		flags=append(flags,false)

		files=append(files,"./testfile/tt.txt")
		flags=append(flags,true)


		for k,v:=range files{
			_,err=OpenWithFlagPerm(v,os.O_RDWR,666)
			if flags[k]{
				gtest.Assert(err,nil)
			}else{
				gtest.AssertNE(err,nil)
			}

		}


	})
}




func TestExists(t *testing.T) {

	gtest.Case(t, func(){
		var(
			flag bool
			files []string
			flags []bool
		)

		files=append(files,"./testfile/file1/nc1.txt")
		flags=append(flags,false)

		files=append(files,"./testfile/tt.txt")
		flags=append(flags,true)


		for k,v:=range files{
			flag=Exists(v)
			if flags[k]{
				gtest.Assert(flag,true)
			}else{
				gtest.Assert(flag,false)
			}

		}


	})
}


func TestPwd(t *testing.T) {
	gtest.Case(t, func(){
		paths,err:=os.Getwd()
		gtest.Assert(err,nil)
		gtest.Assert(Pwd(),paths)

	})
}

func TestIsFile(t *testing.T) {
	gtest.Case(t, func(){
		var(
			flag bool
			files []string
			flags []bool
		)

		files=append(files,"./testfile/file1/nc1.txt")
		flags=append(flags,false)

		files=append(files,"./testfile/tt.txt")
		flags=append(flags,true)

		files=append(files,"./testfile")
		flags=append(flags,false)


		for k,v:=range files{
			flag=IsFile(v)
			if flags[k]{
				gtest.Assert(flag,true)
			}else{
				gtest.Assert(flag,false)
			}

		}


	})
}


func TestInfo(t *testing.T) {
	gtest.Case(t, func(){
		var(
			err error
			paths string ="./testfile/tt.txt"
			files os.FileInfo
			files2 os.FileInfo
		)

		files,err=Info(paths)
		gtest.Assert(err,nil)


		files2,err=os.Stat(paths)
		gtest.Assert(err,nil)

		gtest.Assert(files,files2)

	})
}


func TestMove(t *testing.T) {
	gtest.Case(t, func(){
		var(
			paths string ="./testfile/havefile1/ttn1.txt"
			topath string ="./testfile/havefile1/ttn2.txt"
		)

		gtest.Assert(Move(paths,topath),nil)

	})
}

 func TestRename(t *testing.T){
	 gtest.Case(t, func(){
		 var(

			 paths string ="./testfile/havefile1/ttm1.txt"
			 topath string ="./testfile/havefile1/ttm2.txt"

		 )

		 gtest.Assert(Rename(paths,topath),nil)
		gtest.Assert(IsFile(topath),true)


	 })


 }

func TestCopy(t *testing.T) {
	gtest.Case(t, func(){
		var(
			paths string ="./testfile/havefile1/copyfile1.txt"
			topath string ="./testfile/havefile1/copyfile2.txt"
		)

		gtest.Assert(Copy(paths,topath),nil)
		gtest.Assert(IsFile(topath),true)


	})
}

func  TestDirNames(t *testing.T)  {
	gtest.Case(t, func(){
		var(
			paths string ="./testfile/dirfiles"
			err error
			readlist []string

		)
		havelist:=[]string{
			"t1.txt",
			"t2.txt",
		}
		readlist,err=DirNames(paths)

		gtest.Assert(err,nil)
		gtest.Assert(havelist,readlist)



	})
}


func TestGlob(t *testing.T) {
	gtest.Case(t, func(){
		var(
			paths string ="./testfile/dirfiles/*.txt"
			err error
			resultlist []string

		)

		havelist1:=[]string{
			"t1.txt",
			"t2.txt",
		}

		havelist2:=[]string{
			"testfile/dirfiles/t1.txt",
			"testfile/dirfiles/t2.txt",
		}

		resultlist,err=Glob(paths,true)
		gtest.Assert(err,nil)
		gtest.Assert(resultlist,havelist1)


		resultlist,err=Glob(paths,false)

		//转换成统一的目录分隔符
		for k,v:=range resultlist{
			resultlist[k]=filepath.ToSlash(v)
		}
		gtest.Assert(err,nil)
		gtest.Assert(resultlist,havelist2)

	})
}

func TestRemove(t *testing.T) {
	gtest.Case(t, func(){
		var(
			paths string ="./testfile/delfile/t1.txt"

		)

		gtest.Assert(Remove(paths),nil)


	})
}

func TestIsReadable(t *testing.T){
	gtest.Case(t, func(){
		var(
			paths1 string ="./testfile/havefile1/GetContents.txt"
			paths2 string ="./testfile/havefile1/GetContents_no.txt"
		)
		gtest.Assert(IsReadable(paths1),true)
		gtest.Assert(IsReadable(paths2),false)

	})
}

func TestIsWritable(t *testing.T){
	gtest.Case(t, func(){
		var(
			paths1 string ="./testfile/havefile1/GetContents.txt"
			paths2 string ="./testfile/havefile1/GetContents_no.txt"
		)
		gtest.Assert(IsWritable(paths1),true)
		gtest.Assert(IsWritable(paths2),false)

	})
}


func TestChmod(t *testing.T) {
	gtest.Case(t, func(){
		var(
			paths1 string ="./testfile/havefile1/GetContents.txt"
			paths2 string ="./testfile/havefile1/GetContents_no.txt"
		)


		gtest.Assert(Chmod(paths1,0777),nil)
		gtest.AssertNE(Chmod(paths2,0777),nil)

	})
}


func TestScanDir(t *testing.T){
	gtest.Case(t, func(){
		var(
			paths1 string ="./testfile/dirfiles"
			files []string
			err error
		)
		files,err=ScanDir(paths1,"t*")

		result:=[]string{
			"./testfile/dirfiles/t1.txt",
			"./testfile/dirfiles/t2.txt",
		}

		gtest.Assert(err,nil)

		for k,v:=range files{
			files[k]=filepath.ToSlash(v)
		}

		gtest.Assert(files,result)


	})
}

//获取绝对目录地址
func TestRealPath(t *testing.T){
	gtest.Case(t, func(){
		var(
			paths1 string ="./testfile/dirfiles"
			readlPath string

			tempstr string
		)
		readlPath=RealPath(paths1)
		readlPath=filepath.ToSlash(readlPath)

		tempstr,_=filepath.Abs("./")
		paths1=tempstr+paths1
		paths1=filepath.ToSlash(paths1)
		paths1=strings.Replace(paths1,"./","/",1)


		gtest.Assert(readlPath,paths1)


	})
}


//获取当前执行文件的目录
//注意：当用go test运行测试时，会产生临时的目录文件
func TestSelfPath(t *testing.T){
	gtest.Case(t, func(){
		var(
			paths1 string
			readlPath string
			tempstr string
		)
		readlPath=SelfPath()
		readlPath=filepath.ToSlash(readlPath)

		//
		tempstr,_=filepath.Abs(os.Args[0])
		paths1=filepath.ToSlash(tempstr)
		paths1=strings.Replace(paths1,"./","/",1)


		gtest.Assert(readlPath,paths1)

	})
}


func TestSelfDir(t *testing.T){
	gtest.Case(t, func(){
		var(
			paths1 string
			readlPath string
			tempstr string
		)
		readlPath=SelfDir()


		tempstr,_=filepath.Abs(os.Args[0])
		paths1=filepath.Dir(tempstr)

		gtest.Assert(readlPath,paths1)

	})
}


func TestBasename(t *testing.T){
	gtest.Case(t, func(){
		var(
			paths1 string ="./testfile/havefile1/GetContents.txt"
			readlPath string

		)
		readlPath=Basename(paths1)
		gtest.Assert(readlPath,"GetContents.txt")

	})
}

func TestDir(t *testing.T){
	gtest.Case(t, func(){
		var(
			paths1 string ="./testfile/havefile1"
			readlPath string
		)
		readlPath=Dir(paths1)

		gtest.Assert(readlPath,"testfile")

	})
}


//获取文件名
func TestExt(t *testing.T) {
	gtest.Case(t, func(){
		var(
			paths1 string ="./testfile/havefile1/GetContents.txt"
		)

		gtest.Assert(Ext(paths1),".txt")

	})
}


func TestTempDir(t *testing.T){
	gtest.Case(t, func(){
		var(
			tpath string
		)

		tpath=TempDir()
		gtest.Assert(tpath,os.TempDir())

	})
}







