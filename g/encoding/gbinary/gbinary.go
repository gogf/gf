// Copyright 2017 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

// 二进制及byte操作管理包

package gbinary

import (
    "bytes"
    "encoding/binary"
    "math"
)

// 二进制位(0|1)
type Bit uint8

// (通用,效率较低)二进制打包
func Encode(vs ...interface{}) ([]byte, error) {
    buf := new(bytes.Buffer)
    for i := 0; i < len(vs); i++ {
        err := binary.Write(buf, binary.LittleEndian, vs[i])
        if err != nil {
            return nil, err
        }
    }
    return buf.Bytes(), nil
}

// (通用,效率较低)二进制解包，注意第二个参数之后的变量是变量的指针地址
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
    c := make([]byte, 0)
    c  = append(c, b...)
    for i := 0; i <= l - len(b); i++ {
        c = append(c, 0x00)
    }
    return c
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

// 将ui按位合并到bits数组中，并占length长度位(注意：uis数组中存放的是二进制的0|1数字)
func EncodeBits(bits []Bit, ui uint, l int) []Bit {
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
        b = append(b, byte(DecodeBits(bits[i : i + 8])))
    }
    return b
}

// 从ui字位数组中解析为uint
func DecodeBits(bits []Bit) uint {
    ui := uint(0)
    for _, i := range bits {
        ui = ui << 1 | uint(i)
    }
    return ui
}

// 解析[]byte为字位数组[]uint8
func DecodeBytesToBits(bs []byte) []Bit {
    bits := make([]Bit, 0)
    for _, b := range bs {
        bits = EncodeBits(bits, uint(b), 8)
    }
    return bits
}