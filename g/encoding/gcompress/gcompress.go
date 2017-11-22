package gcompress

import (
    "bytes"
    "compress/zlib"
    "io"
)

// 进行zlib压缩
func Zlib(data []byte) []byte {
    if data == nil || len(data) < 13 {
        return data
    }
    var in bytes.Buffer
    w := zlib.NewWriter(&in)
    w.Write(data)
    w.Close()
    return in.Bytes()
}

// 进行zlib解压缩
func UnZlib(data []byte) []byte {
    if data == nil || len(data) < 13 {
        return data
    }
    b := bytes.NewReader(data)
    var out bytes.Buffer
    r, err := zlib.NewReader(b)
    if err != nil {
        return nil
    }
    io.Copy(&out, r)
    return out.Bytes()
}