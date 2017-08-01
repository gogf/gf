package graft

import (
    "net"
    "encoding/json"
    "log"
    "time"
    "io"
    "g/util/gutil"
    "g/encoding/gjson"
)

// 获取数据
func Receive(conn net.Conn) []byte {
    try        := 0
    buffersize := 1024
    data       := make([]byte, 0)
    for {
        buffer      := make([]byte, buffersize)
        length, err := conn.Read(buffer)
        if err != nil {
            if err != io.EOF {
                log.Println("node receive:", err, "try:", try)
            }
            if try > 2 {
                break;
            }
            try ++
            time.Sleep(100 * time.Millisecond)
        } else {
            if length == buffersize {
                data = gutil.MergeSlice(data, buffer)
            } else {
                data = gutil.MergeSlice(data, buffer[0:length])
                break;
            }
        }
    }
    return data
}

// 获取Msg
func RecieveMsg(conn net.Conn) *Msg {
    response := Receive(conn)
    //log.Println(response)
    if response != nil && len(response) > 0 {
        var msg Msg
        err := json.Unmarshal(response, &msg)
        if err != nil {
            log.Println(err)
            return nil
        }
        return &msg
    }
    return nil
}

// 发送数据
func Send(conn net.Conn, data []byte) error {
    try := 0
    for {
        _, err := conn.Write(data)
        if err != nil {
            log.Println("data send:", err, "try:", try)
            if try > 2 {
                return err
            }
            try ++
            time.Sleep(100 * time.Millisecond)
        } else {
            return nil
        }
    }
}

// 发送Msg
func SendMsg(conn net.Conn, head int, body interface{}) error {
    var msg = Msg{
        Head : head,
        Body : *gjson.Encode(body),
    }
    s, err := json.Marshal(msg)
    if err != nil {
        return err
    }
    return Send(conn, s)
}