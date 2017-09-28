package gbinary

import (
    "bytes"
    "encoding/binary"
)

// 二进制打包
func Encode(vs ...interface{}) []byte {
    buf := new(bytes.Buffer)
    for i := 0; i < len(vs); i++ {
        binary.Write(buf, binary.LittleEndian, vs[i])
    }
    return buf.Bytes()
}

// 二进制解包
func Decode(b []byte, vs ...interface{}) {
    buf := bytes.NewBuffer(b)
    for i := 0; i < len(vs); i++ {
        binary.Read(buf, binary.LittleEndian, vs[i])
    }
}