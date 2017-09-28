package gbinary

import (
    "bytes"
    "encoding/binary"
)

// 二进制打包
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

// 二进制解包，注意第二个参数之后的变量是变量的指针地址
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

func DecodeToInt32(b []byte) (int32, error) {
    var i int32
    err := Decode(b, &i)
    if err != nil {
        return 0, err
    }
    return i, nil
}

func DecodeToInt64(b []byte) (int64, error) {
    var i int64
    err := Decode(b, &i)
    if err != nil {
        return 0, err
    }
    return i, nil
}

func DecodeToBytes(b []byte, size int) ([]byte, error) {
    r := make([]byte, size)
    err := Decode(b, &r)
    if err != nil {
        return nil, err
    }
    return r, nil
}

func DecodeToString(b []byte, size int) (string, error) {
    r, err := DecodeToBytes(b, size)
    if err != nil {
        return "", err
    }
    return string(r), nil
}