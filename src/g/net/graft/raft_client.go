package graft

import (
    "net"
    "encoding/json"
    "log"
    "time"
    "io"
    "g/util/gutil"
)

// 获取数据
func Receive(conn net.Conn) []byte {
    conn.SetReadDeadline(time.Now().Add(gTCP_READ_TIMEOUT * time.Millisecond))
    retry      := 0
    buffersize := 1024
    data       := make([]byte, 0)
    for {
        buffer      := make([]byte, buffersize)
        length, err := conn.Read(buffer)
        if err != nil {
            if retry > gTCP_RETRY_COUNT - 1 {
                break;
            }
            if err != io.EOF {
                //log.Println("receive err:", err, "retry:", retry)
            }
            retry ++
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
    data := Receive(conn)
    //log.Println(string(data))
    if data != nil && len(data) > 0 {
        var msg Msg
        err := json.Unmarshal(data, &msg)
        if err != nil {
            log.Println("receive msg parse err:", err)
            return nil
        }
        return &msg
    }
    return nil
}

// 发送数据
func Send(conn net.Conn, data []byte) error {
    conn.SetReadDeadline(time.Now().Add(gTCP_WRITE_TIMEOUT * time.Millisecond))
    retry := 0
    for {
        _, err := conn.Write(data)
        if err != nil {
            if retry > gTCP_RETRY_COUNT - 1 {
                return err
            }
            //log.Println("data send:", err, "try:", retry)
            retry ++
            time.Sleep(100 * time.Millisecond)
        } else {
            return nil
        }
    }
}

// 发送Msg
func SendMsg(conn net.Conn, head int, body string) error {
    var msg = Msg{
        Head : head,
        Body : body,
    }
    s, err := json.Marshal(msg)
    if err != nil {
        log.Println("send msg parse err:", err)
        return err
    }
    return Send(conn, s)
}