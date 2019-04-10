package gfile

import (
	"github.com/gogf/gf/g/test/gtest"
	"io/ioutil"
	"strings"
	"testing"
)


func TestGetContents(t *testing.T) {
	gtest.Case(t,func(){
		var(
			filepaths string= "./testfile/havefile1/GetContents.txt"
		)

		gtest.Assert(GetContents(filepaths),"abcdefghijkmln")
		gtest.Assert(GetContents(""),"")

	})
}


func TestGetBinContents(t *testing.T) {
	gtest.Case(t , func() {
		var(
			filepaths1 string="./testfile/havefile1/GetContents.txt" //存在文件
			filepaths2 string="./testfile/havefile1/GetContents_no.txt"  //不存大文件
			readcontent []byte
		)
		readcontent=GetBinContents(filepaths1)
		gtest.Assert(readcontent,[]byte("abcdefghijkmln"))


		readcontent=GetBinContents(filepaths2)
		gtest.Assert(readcontent,nil)

		//if readcontent!=nil{
		//	t.Error("文件应不存在")
		//}
		gtest.Assert(GetBinContents(filepaths2),nil)



	})
}

//截断文件为指定的大小
func TestTruncate(t *testing.T) {
	gtest.Case(t , func() {
		var(
			filepaths1 string="./testfile/havefile1/GetContents.txt" //存在文件
			err error
		)
		err=Truncate(filepaths1,200)
		gtest.Assert(err,nil)

		err=Truncate("",200)
		gtest.AssertNE(err,nil)

	})
}

func TestPutContents(t *testing.T) {
	gtest.Case(t , func() {
		var(
			filepaths string="./testfile/havefile1/PutContents.txt"
			err error
			readcontent []byte
		)

		err=PutContents(filepaths,"test!")
		gtest.Assert(err,nil)

		//==================判断是否真正写入
		readcontent, err=ioutil.ReadFile(filepaths)
		gtest.Assert(err,nil)
		gtest.Assert(string(readcontent),"test!")


		err=PutContents("","test!")
		gtest.AssertNE(err,nil)



	})
}





func TestPutContentsAppend(t *testing.T) {
	gtest.Case(t , func() {
		var(
			filepaths string="./testfile/havefile1/PutContents.txt"
			err error
			readcontent []byte
		)

		err=PutContentsAppend(filepaths,"hello")
		gtest.Assert(err,nil)

		//==================判断是否真正写入
		readcontent, err=ioutil.ReadFile(filepaths)
		gtest.Assert(err,nil)
		gtest.Assert(string(readcontent),"test!hello")


		err=PutContentsAppend("","hello")
		gtest.AssertNE(err,nil)



	})


}


func TestPutBinContents(t *testing.T){
	gtest.Case(t , func() {
		var(
			filepaths string="./testfile/havefile1/PutContents.txt"
			err error
			readcontent []byte
		)

		err=PutBinContents(filepaths,[]byte("test!!"))
		gtest.Assert(err,nil)

		//==================判断是否真正写入
		readcontent, err=ioutil.ReadFile(filepaths)
		gtest.Assert(err,nil)
		gtest.Assert(string(readcontent),"test!!")


		err=PutBinContents("",[]byte("test!!"))
		gtest.AssertNE(err,nil)



	})
}


func TestPutBinContentsAppend(t *testing.T) {
	gtest.Case(t , func() {
		var(
			filepaths string="./testfile/havefile1/PutContents.txt"  //原文件内容: yy
			err error
			readcontent []byte
		)

		err=PutBinContentsAppend(filepaths,[]byte("word"))
		gtest.Assert(err,nil)

		//==================判断是否真正写入
		readcontent, err=ioutil.ReadFile(filepaths)
		gtest.Assert(err,nil)
		gtest.Assert(string(readcontent),"test!!word")


		err=PutBinContentsAppend("",[]byte("word"))
		gtest.AssertNE(err,nil)


	})
}

func TestGetBinContentsByTwoOffsetsByPath(t *testing.T) {
	gtest.Case(t, func() {
		var (
			filepaths   string = "./testfile/havefile1/GetContents.txt" //原文件内容: abcdefghijk
			readcontent []byte
		)

		readcontent = GetBinContentsByTwoOffsetsByPath(filepaths, 2, 5)

		gtest.Assert(string(readcontent), "cde")

		readcontent = GetBinContentsByTwoOffsetsByPath("", 2, 5)
		gtest.Assert(len(readcontent),0)

	})

}


func TestGetNextCharOffsetByPath(t *testing.T) {
	gtest.Case(t, func() {
		var (
			filepaths   string = "./testfile/havefile1/GetContents.txt" //原文件内容: abcdefghijk
			localindex int64

		)

		localindex = GetNextCharOffsetByPath(filepaths,'d', 1)
		gtest.Assert(localindex, 3)

		localindex = GetNextCharOffsetByPath("",'d', 1)
		gtest.Assert(localindex, -1)

	})
}


func TestGetNextCharOffset(t *testing.T) {
	gtest.Case(t, func() {
		var (
			localindex int64

		)
		reader:=strings.NewReader("helloword")

		localindex = GetNextCharOffset(reader,'w', 1)
		gtest.Assert(localindex,5)

		localindex = GetNextCharOffset(reader,'j', 1)
		gtest.Assert(localindex,-1)


	})
}

func TestGetBinContentsByTwoOffsets(t *testing.T) {
	gtest.Case(t, func() {
		var (
			reads []byte

		)
		reader:=strings.NewReader("helloword")

		reads = GetBinContentsByTwoOffsets(reader,1, 3)
		gtest.Assert(string(reads),"el")

		reads = GetBinContentsByTwoOffsets(reader,10, 30)
		gtest.Assert(string(reads),"")

	})
}

func TestGetBinContentsTilChar(t *testing.T) {
	gtest.Case(t, func() {
		var (
			reads []byte
			indexs int64

		)
		reader:=strings.NewReader("helloword")

		reads,_ = GetBinContentsTilChar(reader,'w', 2)
		gtest.Assert(string(reads),"llow")

		_,indexs = GetBinContentsTilChar(reader,'w', 20)
		gtest.Assert(indexs,-1)

	})
}

func TestGetBinContentsTilCharByPath(t *testing.T) {
	gtest.Case(t, func() {
		var (
			reads []byte
			indexs int64
			filepaths   string = "./testfile/havefile1/GetContents.txt"

		)


		reads,_ = GetBinContentsTilCharByPath(filepaths,'c',2)
		gtest.Assert(string(reads),"c")

		reads,_ = GetBinContentsTilCharByPath(filepaths,'y',1)
		gtest.Assert(string(reads),"")


		_,indexs = GetBinContentsTilCharByPath(filepaths,'x',1)
		gtest.Assert(indexs,-1)



	})
}

func TestHome(t *testing.T) {
	gtest.Case(t, func() {
		var (
			reads string
			err error

		)

 		reads,err=Home()

 		gtest.Assert(err,nil)
 		gtest.AssertNE(reads,"")


	})
}











