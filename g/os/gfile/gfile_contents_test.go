package gfile

import(
	"io/ioutil"
	"testing"
	"github.com/gogf/gf/g/test/gtest"

)


func TestGetContents(t *testing.T) {
	gtest.Case(t,func(){
		var(
			filepaths string= "./testfile/havefile1/GetContents.txt"
		)

		gtest.Assert(GetContents(filepaths),"abcdefghijkmln")

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
		//文件不存在时@todo:等断言功能优化后，再来修改这里
		//if readcontent!=nil{
		//	t.Error("文件应不存在")
		//}
		gtest.Assert(GetBinContents(filepaths2),nil)



	})
}

//暂时不知道用途 @todo:继续添加
func TestTruncate(t *testing.T) {
	gtest.Case(t , func() {


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
	})

}












