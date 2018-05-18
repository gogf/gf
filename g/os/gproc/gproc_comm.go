// Copyright 2018 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.
// @todo 后期改为tcp进程通信形式


package gproc

import (
    "io"
    "os"
    "fmt"
    "time"
    "errors"
    "gitee.com/johng/gf/g/os/glog"
    "gitee.com/johng/gf/g/os/gfile"
    "gitee.com/johng/gf/g/os/gtime"
    "gitee.com/johng/gf/g/os/gflock"
    "gitee.com/johng/gf/g/util/gconv"
    "gitee.com/johng/gf/g/os/gfsnotify"
    "gitee.com/johng/gf/g/container/gqueue"
    "gitee.com/johng/gf/g/encoding/gbinary"
)

const (
    // 由于子进程的temp dir有可能会和父进程不一致(特别是windows下)，影响进程间通信，这里统一使用环境变量设置
    gPROC_TEMP_DIR_ENV_KEY           = "gproc.tempdir"
    // 写入通信数据失败时候的重试次数
    gPROC_COMM_FAILURE_RETRY_COUNT   = 3
    // (毫秒)主动通信内容检查时间间隔
    gPROC_COMM_ACTIVE_CHECK_INTERVAL = 500*time.Millisecond
)

// 全局通信文件清理文件锁(同一时刻只能存在一个进程进行通信文件清理)
var commClearLocker = gflock.New("comm.clear.lock")
// 当前进程的文件锁
var commLocker      = gflock.New(fmt.Sprintf("%d.lock", os.Getpid()))
// 进程通信消息队列
var commQueue       = gqueue.New()

// TCP通信数据结构定义
type Msg struct {
    Pid  int    // PID，哪个进程发送的消息
    Data []byte // 参数，消息附带的参数
}

// 进程管理/通信初始化操作
func init() {
    path := getCommFilePath(os.Getpid())
    checkAndInitCommFile(path)
    commLocker.Lock()
    fileMtime := gfile.MTime(path)
    commLocker.UnLock()
    if gtime.Second() - fileMtime < 10 {
        // 初始化时读取已有数据(文件修改时间在10秒以内)
        checkCommBuffer(path)
    } else {
        // 否则清空旧的数据内容
        commLocker.Lock()
        os.Truncate(path, 0)
        commLocker.UnLock()
    }
    // 文件事件监听，如果通信数据文件有任何变化，读取文件并添加到消息队列
    err := gfsnotify.Add(path, func(event *gfsnotify.Event) {
        checkCommBuffer(path)
    })
    if err != nil {
        glog.Error(err)
    }

    go autoActiveCheckComm()
}

// 检测并初始化通信文件
func checkAndInitCommFile(path string) {
    commLocker.Lock()
    defer commLocker.UnLock()
    if !gfile.Exists(path) {
        err := gfile.Create(path)
        if err != nil {
            glog.Error(err)
            os.Exit(1)
        }
    }
    // 检测写入权限
    if !gfile.IsWritable(path) {
        glog.Errorfln("%s is not writable for gproc", path)
        os.Exit(1)
    }
}

// 主动通信内容检测
func autoActiveCheckComm() {
    path := getCommFilePath(Pid())
    for {
        time.Sleep(gPROC_COMM_ACTIVE_CHECK_INTERVAL)
        checkCommBuffer(path)
    }
}

// 手动检查进程通信消息，如果存在消息曾推送到进程消息队列
func checkCommBuffer(path string) {
    commLocker.Lock()
    // 读取进程间通信数据
    buffer := gfile.GetBinContents(path)
    if len(buffer) > 0 {
        os.Truncate(path, 0)
    }
    commLocker.UnLock()
    if len(buffer) > 0 {
        for _, v := range bufferToMsgs(buffer) {
            commQueue.PushBack(v)
        }
    }
}

// 获取其他进程传递到当前进程的消息包，阻塞执行
func Receive() *Msg {
    if v := commQueue.PopFront(); v != nil {
        return v.(*Msg)
    }
    return nil
}

// 向指定gproc进程发送数据
// 数据格式：总长度(32bit) | PID(32bit) | 校验(32bit) | 参数(变长)
func Send(pid int, data []byte) error {
    var err error = nil
    b := make([]byte, 0)
    b  = append(b, gbinary.EncodeInt32(int32(len(data) + 12))...)
    b  = append(b, gbinary.EncodeInt32(int32(os.Getpid()))...)
    b  = append(b, gbinary.EncodeUint32(checksum(data))...)
    b  = append(b, data...)
    l := gflock.New(fmt.Sprintf("%d.lock", pid))
    l.Lock()
    for i := gPROC_COMM_FAILURE_RETRY_COUNT; i > 0; i-- {
        err = doSend(pid, b)
        if err == nil {
            break
        }
    }
    l.UnLock()
    //glog.Printfln("%d to %d, %v, %d, %v", Pid(), pid, data, gfile.Size(getCommFilePath(pid)), err)
    return err
}

// 执行进程间通信数据写入
func doSend(pid int, buffer []byte) error {
    file, err := gfile.OpenWithFlag(getCommFilePath(pid), os.O_RDWR|os.O_CREATE|os.O_APPEND)
    if err != nil{
        return err
    }
    // 必须要Close才会更新文件的修改时间
    defer file.Close()
    // 获取原有文件内容大小
    stat, err := file.Stat()
    if err != nil {
        return err
    }
    oldSize := stat.Size()
    // 执行数据写入
    writeSize, err := file.Write(buffer)
    if err != nil {
        return err
    }
    if writeSize < len(buffer) {
        return io.ErrShortWrite
    }
    // 写入成功之后获取最新文件内容大小，执行对比
    if stat, err := file.Stat(); err != nil {
        return err
    } else {
        // 由于文件锁机制的保证，同一时刻只会有一个进程(&协程)在执行写入，不会出现数据粘包情况
        // 这里从严谨性考虑增加大小判断，更进一步避免粘包，或者丢包情况
        if stat.Size() - int64(writeSize) != oldSize {
            return errors.New("error writing data")
        }
    }
    return nil
}

// 获取指定进程的通信文件地址
func getCommFilePath(pid int) string {
    return getCommDirPath() + gfile.Separator + gconv.String(pid)
}

// 获取进程间通信目录地址
func getCommDirPath() string {
    tempDir := os.Getenv("gproc.tempdir")
    if tempDir == "" {
        tempDir = gfile.TempDir()
    }
    return tempDir + gfile.Separator + "gproc"
}

// 数据解包，防止黏包
func bufferToMsgs(buffer []byte) []*Msg {
    s    := 0
    msgs := make([]*Msg, 0)
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
        msgs = append(msgs, &Msg {
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