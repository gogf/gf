package gbinary

import (
    "bytes"
    "encoding/binary"
    "math"
)

// (通用)二进制打包
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

// (通用)二进制解包，注意第二个参数之后的变量是变量的指针地址
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

func DecodeToInt8(b []byte) int8 {
    return int8(b[0])
}

func DecodeToUint8(b []byte) uint8 {
    return uint8(b[0])
}

func DecodeToInt16(b []byte) int16 {
    return int16(binary.LittleEndian.Uint16(b))
}

func DecodeToUint16(b []byte) uint16 {
    return binary.LittleEndian.Uint16(b)
}

func DecodeToInt32(b []byte) int32 {
    return int32(binary.LittleEndian.Uint32(b))
}

func DecodeToUint32(b []byte) uint32 {
    return binary.LittleEndian.Uint32(b)
}

func DecodeToInt64(b []byte) int64 {
    return int64(binary.LittleEndian.Uint64(b))
}

func DecodeToUint64(b []byte) uint64 {
    return binary.LittleEndian.Uint64(b)
}

func DecodeToFloat32(b []byte) float32 {
    return math.Float32frombits(binary.LittleEndian.Uint32(b))
}

func DecodeToFloat64(b []byte) float64 {
    return math.Float64frombits(binary.LittleEndian.Uint64(b))
}