package main

import (
    "fmt"
    "net/http"
    "github.com/tabalt/gracehttp"
    "gitee.com/johng/gf/g/encoding/gbinary"
    "gitee.com/johng/gf/g/os/gproc"
)


func test() {
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "hello world")
    })
    s := gracehttp.NewServer(":8888", nil, gracehttp.DEFAULT_READ_TIMEOUT, gracehttp.DEFAULT_WRITE_TIMEOUT)
    err := s.ListenAndServe()
    if err != nil {
        fmt.Println(err)
    }
}

// 常见的二进制数据校验方式，生成校验结果
func checksum(buffer []byte) uint32 {
    var checksum uint32
    for _, b := range buffer {
        checksum += uint32(b)
    }
    return checksum
}

// 数据解包，防止黏包
func bufferToMsgs(buffer []byte) []*gproc.Msg {
    s    := 0
    msgs := make([]*gproc.Msg, 0)
    for s < len(buffer) {
        fmt.Println(s)
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


func main() {
    b := []byte{26, 0, 0, 0, 60, 109, 0, 0, 84, 5, 0, 0, 104, 101, 108, 108, 111, 32, 112, 114, 111, 99, 101, 115, 115, 33, 26, 0, 0, 0, 60, 109, 0, 0, 84, 5, 0, 0, 104, 101, 108, 108, 111, 32, 112, 114, 111, 99, 101, 115, 115, 33, 26, 0, 0, 0, 60, 109, 0, 0, 84, 5, 0, 0, 104, 101, 108, 108, 111, 32, 112, 114, 111, 99, 101, 115, 115, 33, 26, 0, 0, 0, 60, 109, 0, 0, 84, 5, 0, 0, 104, 101, 108, 108, 111, 32, 112, 114, 111, 99, 101, 115, 115, 33, 26, 0, 0, 0, 60, 109, 0, 0, 84, 5, 0, 0, 104, 101, 108, 108, 111, 32, 112, 114, 111, 99, 101, 115, 115, 33, 26, 0, 0, 0, 60, 109, 0, 0, 84, 5, 0, 0, 104, 101, 108, 108, 111, 32, 112, 114, 111, 99, 101, 115, 115, 33, 26, 0, 0, 0, 60, 109, 0, 0, 84, 5, 0, 0, 104, 101, 108, 108, 111, 32, 112, 114, 111, 99, 101, 115, 115, 33, 26, 0, 0, 0, 60, 109, 0, 0, 84, 5, 0, 0, 104, 101, 108, 108, 111, 32, 112, 114, 111, 99, 101, 115, 115, 33, 26, 0, 0, 0, 60, 109, 0, 0, 84, 5, 0, 0, 104, 101, 108, 108, 111, 32, 112, 114, 111, 99, 101, 115, 115, 33, 26, 0, 0, 0, 60, 109, 0, 0, 84, 5, 0, 0, 104, 101, 108, 108, 111, 32, 112, 114, 111, 99, 101, 115, 115, 33, 26, 0, 0, 0, 60, 109, 0, 0, 84, 5, 0, 0, 104, 101, 108, 108, 111, 32, 112, 114, 111, 99, 101, 115, 115, 33}
    m := bufferToMsgs(b)
    fmt.Println(m)
}
