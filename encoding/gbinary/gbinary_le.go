// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gbinary

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
)

// LeEncode encodes one or multiple <values> into bytes using LittleEndian.
// It uses type asserting checking the type of each value of <values> and internally
// calls corresponding converting function do the bytes converting.
//
// It supports common variable type asserting, and finally it uses fmt.Sprintf converting
// value to string and then to bytes.
func LeEncode(values ...interface{}) []byte {
	buf := new(bytes.Buffer)
	for i := 0; i < len(values); i++ {
		if values[i] == nil {
			return buf.Bytes()
		}

		switch value := values[i].(type) {
		case int:
			buf.Write(LeEncodeInt(value))
		case int8:
			buf.Write(LeEncodeInt8(value))
		case int16:
			buf.Write(LeEncodeInt16(value))
		case int32:
			buf.Write(LeEncodeInt32(value))
		case int64:
			buf.Write(LeEncodeInt64(value))
		case uint:
			buf.Write(LeEncodeUint(value))
		case uint8:
			buf.Write(LeEncodeUint8(value))
		case uint16:
			buf.Write(LeEncodeUint16(value))
		case uint32:
			buf.Write(LeEncodeUint32(value))
		case uint64:
			buf.Write(LeEncodeUint64(value))
		case bool:
			buf.Write(LeEncodeBool(value))
		case string:
			buf.Write(LeEncodeString(value))
		case []byte:
			buf.Write(value)
		case float32:
			buf.Write(LeEncodeFloat32(value))
		case float64:
			buf.Write(LeEncodeFloat64(value))
		default:
			if err := binary.Write(buf, binary.LittleEndian, value); err != nil {
				buf.Write(LeEncodeString(fmt.Sprintf("%v", value)))
			}
		}
	}
	return buf.Bytes()
}

// 将变量转换为二进制[]byte，并指定固定的[]byte长度返回，长度单位为字节(byte)；
// 如果转换的二进制长度超过指定长度，那么进行截断处理
func LeEncodeByLength(length int, values ...interface{}) []byte {
	b := LeEncode(values...)
	if len(b) < length {
		b = append(b, make([]byte, length-len(b))...)
	} else if len(b) > length {
		b = b[0:length]
	}
	return b
}

// 整形二进制解包，注意第二个及其后参数为字长确定的整形变量的指针地址，以便确定解析的[]byte长度，
// 例如：int8/16/32/64、uint8/16/32/64、float32/64等等
func LeDecode(b []byte, values ...interface{}) error {
	buf := bytes.NewBuffer(b)
	for i := 0; i < len(values); i++ {
		err := binary.Read(buf, binary.LittleEndian, values[i])
		if err != nil {
			return err
		}
	}
	return nil
}

func LeEncodeString(s string) []byte {
	return []byte(s)
}

func LeDecodeToString(b []byte) string {
	return string(b)
}

func LeEncodeBool(b bool) []byte {
	if b == true {
		return []byte{1}
	} else {
		return []byte{0}
	}
}

// 自动识别int类型长度，转换为[]byte
func LeEncodeInt(i int) []byte {
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
func LeEncodeUint(i uint) []byte {
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

func LeEncodeInt8(i int8) []byte {
	return []byte{byte(i)}
}

func LeEncodeUint8(i uint8) []byte {
	return []byte{i}
}

func LeEncodeInt16(i int16) []byte {
	b := make([]byte, 2)
	binary.LittleEndian.PutUint16(b, uint16(i))
	return b
}

func LeEncodeUint16(i uint16) []byte {
	b := make([]byte, 2)
	binary.LittleEndian.PutUint16(b, i)
	return b
}

func LeEncodeInt32(i int32) []byte {
	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, uint32(i))
	return b
}

func LeEncodeUint32(i uint32) []byte {
	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, i)
	return b
}

func LeEncodeInt64(i int64) []byte {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, uint64(i))
	return b
}

func LeEncodeUint64(i uint64) []byte {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, i)
	return b
}

func LeEncodeFloat32(f float32) []byte {
	bits := math.Float32bits(f)
	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, bits)
	return b
}

func LeEncodeFloat64(f float64) []byte {
	bits := math.Float64bits(f)
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, bits)
	return b
}

// 将二进制解析为int类型，根据[]byte的长度进行自动转换.
// 注意内部使用的是uint*，使用int会造成位丢失。
func LeDecodeToInt(b []byte) int {
	if len(b) < 2 {
		return int(LeDecodeToUint8(b))
	} else if len(b) < 3 {
		return int(LeDecodeToUint16(b))
	} else if len(b) < 5 {
		return int(LeDecodeToUint32(b))
	} else {
		return int(LeDecodeToUint64(b))
	}
}

// 将二进制解析为uint类型，根据[]byte的长度进行自动转换
func LeDecodeToUint(b []byte) uint {
	if len(b) < 2 {
		return uint(LeDecodeToUint8(b))
	} else if len(b) < 3 {
		return uint(LeDecodeToUint16(b))
	} else if len(b) < 5 {
		return uint(LeDecodeToUint32(b))
	} else {
		return uint(LeDecodeToUint64(b))
	}
}

// 将二进制解析为bool类型，识别标准是判断二进制中数值是否都为0，或者为空。
func LeDecodeToBool(b []byte) bool {
	if len(b) == 0 {
		return false
	}
	if bytes.Compare(b, make([]byte, len(b))) == 0 {
		return false
	}
	return true
}

func LeDecodeToInt8(b []byte) int8 {
	return int8(b[0])
}

func LeDecodeToUint8(b []byte) uint8 {
	return uint8(b[0])
}

func LeDecodeToInt16(b []byte) int16 {
	return int16(binary.LittleEndian.Uint16(LeFillUpSize(b, 2)))
}

func LeDecodeToUint16(b []byte) uint16 {
	return binary.LittleEndian.Uint16(LeFillUpSize(b, 2))
}

func LeDecodeToInt32(b []byte) int32 {
	return int32(binary.LittleEndian.Uint32(LeFillUpSize(b, 4)))
}

func LeDecodeToUint32(b []byte) uint32 {
	return binary.LittleEndian.Uint32(LeFillUpSize(b, 4))
}

func LeDecodeToInt64(b []byte) int64 {
	return int64(binary.LittleEndian.Uint64(LeFillUpSize(b, 8)))
}

func LeDecodeToUint64(b []byte) uint64 {
	return binary.LittleEndian.Uint64(LeFillUpSize(b, 8))
}

func LeDecodeToFloat32(b []byte) float32 {
	return math.Float32frombits(binary.LittleEndian.Uint32(LeFillUpSize(b, 4)))
}

func LeDecodeToFloat64(b []byte) float64 {
	return math.Float64frombits(binary.LittleEndian.Uint64(LeFillUpSize(b, 8)))
}

// 当b位数不够时，进行高位补0。
// 注意这里为了不影响原有输入参数，是采用的值复制设计。
func LeFillUpSize(b []byte, l int) []byte {
	if len(b) >= l {
		return b[:l]
	}
	c := make([]byte, l)
	copy(c, b)
	return c
}
