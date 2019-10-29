package main

import (
	"fmt"

	"github.com/gogf/gf/encoding/gbinary"
	"github.com/gogf/gf/os/glog"
)

func main() {
	// 使用gbinary.Encoded对基本数据类型进行二进制打包
	fmt.Println(gbinary.Encode(18, 300, 1.01))

	// 使用gbinary.Decode对整形二进制解包，注意第二个及其后参数为字长确定的整形变量的指针地址，字长确定的类型，
	// 例如：int8/16/32/64、uint8/16/32/64、float32/64
	// 这里的1.01默认为float64类型(64位系统下)
	buffer := gbinary.Encode(18, 300, 1.01)
	var i1 int8
	var i2 int16
	var f3 float64
	if err := gbinary.Decode(buffer, &i1, &i2, &f3); err != nil {
		glog.Error(err)
	} else {
		fmt.Println(i1, i2, f3)
	}

	// 编码/解析 int，自动识别变量长度
	fmt.Println(gbinary.DecodeToInt(gbinary.EncodeInt(1)))
	fmt.Println(gbinary.DecodeToInt(gbinary.EncodeInt(300)))
	fmt.Println(gbinary.DecodeToInt(gbinary.EncodeInt(70000)))
	fmt.Println(gbinary.DecodeToInt(gbinary.EncodeInt(2000000000)))
	fmt.Println(gbinary.DecodeToInt(gbinary.EncodeInt(500000000000)))

	// 编码/解析 uint，自动识别变量长度
	fmt.Println(gbinary.DecodeToUint(gbinary.EncodeUint(1)))
	fmt.Println(gbinary.DecodeToUint(gbinary.EncodeUint(300)))
	fmt.Println(gbinary.DecodeToUint(gbinary.EncodeUint(70000)))
	fmt.Println(gbinary.DecodeToUint(gbinary.EncodeUint(2000000000)))
	fmt.Println(gbinary.DecodeToUint(gbinary.EncodeUint(500000000000)))

	// 编码/解析 int8/16/32/64
	fmt.Println(gbinary.DecodeToInt8(gbinary.EncodeInt8(int8(100))))
	fmt.Println(gbinary.DecodeToInt16(gbinary.EncodeInt16(int16(100))))
	fmt.Println(gbinary.DecodeToInt32(gbinary.EncodeInt32(int32(100))))
	fmt.Println(gbinary.DecodeToInt64(gbinary.EncodeInt64(int64(100))))

	// 编码/解析 uint8/16/32/64
	fmt.Println(gbinary.DecodeToUint8(gbinary.EncodeUint8(uint8(100))))
	fmt.Println(gbinary.DecodeToUint16(gbinary.EncodeUint16(uint16(100))))
	fmt.Println(gbinary.DecodeToUint32(gbinary.EncodeUint32(uint32(100))))
	fmt.Println(gbinary.DecodeToUint64(gbinary.EncodeUint64(uint64(100))))

	// 编码/解析 string
	fmt.Println(gbinary.DecodeToString(gbinary.EncodeString("I'm string!")))
}
