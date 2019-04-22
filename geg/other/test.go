package main

import (
	"fmt"
	"github.com/gogf/gf/g/crypto/gmd5"
	"github.com/gogf/gf/g/encoding/gbase64"
	"github.com/gogf/gf/g/text/gstr"
)
func main() {
	//fmt.Println([]byte{152})
	//fmt.Println(len([]byte{152}))
	//fmt.Println(len(string([]byte{152})))
	//fmt.Println(len(string(byte(152))))
	//os.Exit(1)
	//fmt.Println(gstr.Chr(152))
	//fmt.Println(len(gstr.Chr(152)))
	//os.Exit(1)
	data := "abcdefg"
	dict := "no"
	key := gmd5.EncryptString(dict)
	x := 0
	lenb := len(data)
	l := len(key)
	char := ""
	strb := ""
	for i := 0; i < lenb; i++ {
		if x == l {
			x = 0
		}
		char += key[x : x+1]
		x++
	}
	for i := 0; i < lenb; i++ {
		fmt.Println((gstr.Ord(data[i:i+1]) + gstr.Ord(char[i:i+1]) % 256))
		//fmt.Println(gstr.Chr((gstr.Ord(data[i:i+1]) + gstr.Ord(char[i:i+1]) % 256) ))
		//fmt.Println(data[i:i+1], gstr.Ord(data[i:i+1]), gstr.Ord(char[i:i+1]))
		strb += gstr.Chr((gstr.Ord(data[i:i+1]) + gstr.Ord(char[i:i+1]) % 256) )
		fmt.Println(len(strb))
		fmt.Println("=============")
	}
	fmt.Println(strb)
	fmt.Println(len(strb))
	result := gbase64.Encode(strb)
	fmt.Println(result)
}