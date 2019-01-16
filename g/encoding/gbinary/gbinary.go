// Copyright 2017 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

// Package gbinary provides useful API for handling binary/bytes data.
package gbinary

import (
    "fmt"
    "math"
    "bytes"
    "encoding/binary"
)

// 二进制位(0|1)
type Bit int8

// 针对基本类型进行二进制打包，支持的基本数据类型包括:int/8/16/32/64、uint/8/16/32/64、float32/64、bool、string、[]byte
// 其他未知类型使用 fmt.Sprintf("%v", value) 转换为字符串之后处理
func Encode(vs ...interface{}) []byte {
    buf := new(bytes.Buffer)
    for i := 0; i < len(vs); i++ {
        switch value := vs[i].(type) {
            case int:     buf.Write(EncodeInt(value))
            case int8:    buf.Write(EncodeInt8(value))
            case int16:   buf.Write(EncodeInt16(value))
            case int32:   buf.Write(EncodeInt32(value))
            case int64:   buf.Write(EncodeInt64(value))
            case uint:    buf.Write(EncodeUint(value))
            case uint8:   buf.Write(EncodeUint8(value))
            case uint16:  buf.Write(EncodeUint16(value))
            case uint32:  buf.Write(EncodeUint32(value))
            case uint64:  buf.Write(EncodeUint64(value))
            case bool:    buf.Write(EncodeBool(value))
            case string:  buf.Write(EncodeString(value))
            case []byte:  buf.Write(value)
            case float32: buf.Write(EncodeFloat32(value))
            case float64: buf.Write(EncodeFloat64(value))
            default:
                if err := binary.Write(buf, binary.LittleEndian, value); err != nil {
                    buf.Write(EncodeString(fmt.Sprintf("%v", value)))
                }
        }
    }
    return buf.Bytes()
}

// 将变量转换为二进制[]byte，并指定固定的[]byte长度返回，长度单位为字节(byte)；
// 如果转换的二进制长度超过指定长度，那么进行截断处理
func EncodeByLength(length int, vs ...interface{}) []byte {
    b := Encode(vs...)
    if len(b) < length {
        b = append(b, make([]byte, length - len(b))...)
    } else if len(b) > length {
        b = b[0 : length]
    }
    return b
}

// 整形二进制解包，注意第二个及其后参数为字长确定的整形变量的指针地址，以便确定解析的[]byte长度，
// 例如：int8/16/32/64、uint8/16/32/64、float32/64等等
func Decode(b []byte, vs ...interface{}) error {
    buf := bytes.NewBuffer(b)
    for i := 0; i < len(vs); i++ {
        err := binary.Read(buf, binary.LittleEndian, vs[i])
        if err != nil {
            return err
        }
    }
    return nil
}

func EncodeString(s string) []byte {
    return []byte(s)
}

func DecodeToString(b []byte) string {
    return string(b)
}

func EncodeBool(b bool) []byte {
    if b == true {
        return []byte{1}
    } else {
        return []byte{0}
    }
}

// 自动识别int类型长度，转换为[]byte
func EncodeInt(i int) []byte {
    if i <= math.MaxInt8 {
        return EncodeInt8(int8(i))
    } else if i <= math.MaxInt16 {
        return EncodeInt16(int16(i))
    } else if i <= math.MaxInt32 {
        return EncodeInt32(int32(i))
    } else {
        return EncodeInt64(int64(i))
    }
}

// 自动识别uint类型长度，转换为[]byte
func EncodeUint(i uint) []byte {
    if i <= math.MaxUint8 {
        return EncodeUint8(uint8(i))
    } else if i <= math.MaxUint16 {
        return EncodeUint16(uint16(i))
    } else if i <= math.MaxUint32 {
        return EncodeUint32(uint32(i))
    } else {
        return EncodeUint64(uint64(i))
    }
}

func EncodeInt8(i int8) []byte {
    return []byte{byte(i)}
}

func EncodeUint8(i uint8) []byte {
    return []byte{byte(i)}
}

func EncodeInt16(i int16) []byte {
    bytes := make([]byte, 2)
    binary.LittleEndian.PutUint16(bytes, uint16(i))
    return bytes
}

func EncodeUint16(i uint16) []byte {
    bytes := make([]byte, 2)
    binary.LittleEndian.PutUint16(bytes, i)
    return bytes
}

func EncodeInt32(i int32) []byte {
    bytes := make([]byte, 4)
    binary.LittleEndian.PutUint32(bytes, uint32(i))
    return bytes
}

func EncodeUint32(i uint32) []byte {
    bytes := make([]byte, 4)
    binary.LittleEndian.PutUint32(bytes, i)
    return bytes
}

func EncodeInt64(i int64) []byte {
    bytes := make([]byte, 8)
    binary.LittleEndian.PutUint64(bytes, uint64(i))
    return bytes
}

func EncodeUint64(i uint64) []byte {
    bytes := make([]byte, 8)
    binary.LittleEndian.PutUint64(bytes, i)
    return bytes
}

func EncodeFloat32(f float32) []byte {
    bits  := math.Float32bits(f)
    bytes := make([]byte, 4)
    binary.LittleEndian.PutUint32(bytes, bits)
    return bytes
}

func EncodeFloat64(f float64) []byte {
    bits  := math.Float64bits(f)
    bytes := make([]byte, 8)
    binary.LittleEndian.PutUint64(bytes, bits)
    return bytes
}

// 当b位数不够时，进行高位补0
func fillUpSize(b []byte, l int) []byte {
    if len(b) >= l {
        return b
    }
    c := make([]byte, 0)
    c  = append(c, b...)
    for i := 0; i < l - len(b); i++ {
        c = append(c, 0x00)
    }
    return c
}

// 将二进制解析为int类型，根据[]byte的长度进行自动转换.
// 注意内部使用的是uint*，使用int会造成位丢失。
func DecodeToInt(b []byte) int {
    if len(b) < 2 {
        return int(DecodeToUint8(b))
    } else if len(b) < 3 {
        return int(DecodeToUint16(b))
    } else if len(b) < 5 {
        return int(DecodeToUint32(b))
    } else {
        return int(DecodeToUint64(b))
    }
}

// 将二进制解析为uint类型，根据[]byte的长度进行自动转换
func DecodeToUint(b []byte) uint {
    if len(b) < 2 {
        return uint(DecodeToUint8(b))
    } else if len(b) < 3 {
        return uint(DecodeToUint16(b))
    } else if len(b) < 5 {
        return uint(DecodeToUint32(b))
    } else {
        return uint(DecodeToUint64(b))
    }
}

// 将二进制解析为bool类型，识别标准是判断二进制中数值是否都为0，或者为空
func DecodeToBool(b []byte) bool {
    if len(b) == 0 {
        return false
    }
    if bytes.Compare(b, make([]byte, len(b))) == 0 {
        return false
    }
    return true
}

func DecodeToInt8(b []byte) int8 {
    return int8(b[0])
}

func DecodeToUint8(b []byte) uint8 {
    return uint8(b[0])
}

func DecodeToInt16(b []byte) int16 {
    return int16(binary.LittleEndian.Uint16(fillUpSize(b, 2)))
}

func DecodeToUint16(b []byte) uint16 {
    return binary.LittleEndian.Uint16(fillUpSize(b, 2))
}

func DecodeToInt32(b []byte) int32 {
    return int32(binary.LittleEndian.Uint32(fillUpSize(b, 4)))
}

func DecodeToUint32(b []byte) uint32 {
    return binary.LittleEndian.Uint32(fillUpSize(b, 4))
}

func DecodeToInt64(b []byte) int64 {
    return int64(binary.LittleEndian.Uint64(fillUpSize(b, 8)))
}

func DecodeToUint64(b []byte) uint64 {
    return binary.LittleEndian.Uint64(fillUpSize(b, 8))
}

func DecodeToFloat32(b []byte) float32 {
    return math.Float32frombits(binary.LittleEndian.Uint32(fillUpSize(b, 4)))
}

func DecodeToFloat64(b []byte) float64 {
    return math.Float64frombits(binary.LittleEndian.Uint64(fillUpSize(b, 8)))
}

// 默认编码
func EncodeBits(bits []Bit, i int, l int) []Bit {
    return EncodeBitsWithUint(bits, uint(i), l)
}

// 将ui按位合并到bits数组中，并占length长度位(注意：uis数组中存放的是二进制的0|1数字)
func EncodeBitsWithUint(bits []Bit, ui uint, l int) []Bit {
    a := make([]Bit, l)
    for i := l - 1; i >= 0; i-- {
        a[i] = Bit(ui & 1)
        ui >>= 1
    }
    if bits != nil {
        return append(bits, a...)
    } else {
        return a
    }
}
// 将bits转换为[]byte，从左至右进行编码，不足1 byte按0往末尾补充
func EncodeBitsToBytes(bits []Bit) []byte {
    if len(bits)%8 != 0 {
        for i := 0; i < len(bits)%8; i++ {
            bits = append(bits, 0)
        }
    }
    b := make([]byte, 0)
    for i := 0; i < len(bits); i += 8 {
        b = append(b, byte(DecodeBitsToUint(bits[i : i + 8])))
    }
    return b
}

// 解析为int
func DecodeBits(bits []Bit) int {
    v := int(0)
    for _, i := range bits {
        v = v << 1 | int(i)
    }
    return v
}

// 解析为uint
func DecodeBitsToUint(bits []Bit) uint {
    v := uint(0)
    for _, i := range bits {
        v = v << 1 | uint(i)
    }
    return v
}

// 解析[]byte为字位数组[]uint8
func DecodeBytesToBits(bs []byte) []Bit {
    bits := make([]Bit, 0)
    for _, b := range bs {
        bits = EncodeBitsWithUint(bits, uint(b), 8)
    }
    return bits
}