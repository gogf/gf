package main

import (
    "gitee.com/johng/gf/g/encoding/gbinary"
    "gitee.com/johng/gf/g/os/glog"
    "fmt"
)

func main() {
    // 使用gbinary.Encoded对整形二进制打包，注意参数必须为字长确定的类型：int8/16/32/64、uint8/16/32/64
    if buffer, err := gbinary.Encode(int32(18), int64(24)); err != nil {
        glog.Error(err)
    } else {
        fmt.Println(buffer)
    }


    // 使用gbinary.Decode对整形二进制解包，注意第二个及其后参数为字长确定的整形变量的指针地址，字长确定的类型：int8/16/32/64、uint8/16/32/64
    if buffer, err := gbinary.Encode(int32(18), int64(24)); err != nil {
        glog.Error(err)
    } else {
        var i1 int32
        var i2 int64
        if err := gbinary.Decode(buffer, &i1, &i2); err != nil {
            glog.Error(err)
        } else {
            fmt.Println(i1, i2)
        }
    }

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
