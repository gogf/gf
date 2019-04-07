package gfile

import (
	"github.com/gogf/gf/g/test/gtest"
	"os"
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
		// gtest.Assert(IsFile(topath),true)


	 })


 }



