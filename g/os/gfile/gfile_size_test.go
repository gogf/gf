package gfile

import (
	"github.com/gogf/gf/g/test/gtest"
	"testing"
)

func TestSize(t *testing.T){
	gtest.Case(t, func(){
		var(
			paths1 string ="./testfile/dirfiles/t1.txt"
			sizes int64
		)
		sizes=Size(paths1)
		gtest.Assert(sizes,16)


	})
}


func TestFormatSize(t *testing.T){
	gtest.Case(t, func(){

		gtest.Assert(FormatSize(16),"16.00B")

		gtest.Assert(FormatSize(1600),"1.56K")

		gtest.Assert(FormatSize(16000000),"15.26M")

		gtest.Assert(FormatSize(1600000000),"1.49G")


	})
}



func TestReadableSize(t *testing.T){
	gtest.Case(t, func(){

		gtest.Assert(ReadableSize("./testfile/dirfiles/t1.txt"),"16.00B")

	})
}
