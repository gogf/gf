package main

import (
    "fmt"
    "gitee.com/johng/gf/g/os/gfile"
    "gitee.com/johng/gf/g/os/gtime"
    "gitee.com/johng/gf/g/encoding/gbinary"
    "gitee.com/johng/gf/g/os/gproc"
)

// 数据解包，防止黏包
func bufferToMsgs(buffer []byte) []*gproc.Msg {
    s    := 0
    msgs := make([]*gproc.Msg, 0)
    for s < len(buffer) {
        length := gbinary.DecodeToInt(buffer[s : s + 4])
        if length < 0 || length > len(buffer) {
            s++
            continue
        }
        checksum1 := gbinary.DecodeToUint32(buffer[s + 8 : s + 12])
        checksum2 := checksum(buffer[s + 12 : s + length])
        if checksum1 != checksum2 {
            s++
            continue
        }
        msgs = append(msgs, &gproc.Msg {
            Pid  : gbinary.DecodeToInt(buffer[s + 4 : s + 8]),
            Data : buffer[s + 12 : s + length],
        })
        s += length
    }
    return msgs
}

// 常见的二进制数据校验方式，生成校验结果
func checksum(buffer []byte) uint32 {
    var checksum uint32
    for _, b := range buffer {
        checksum += uint32(b)
    }
    return checksum
}

func main(){
    b := gfile.GetBinContents("/home/john/Documents/11248")
    for _, msg := range bufferToMsgs(b) {
        fmt.Println(msg.Pid)
        fmt.Println(msg.Data)
    }

    return
    t1 := gfile.MTime("/home/john/Workspace/Go/GOPATH/src/gitee.com/johng/gf/geg/other/test.go")
    t2 := gtime.Second()
    fmt.Println(t1)
    fmt.Println(t2)
}