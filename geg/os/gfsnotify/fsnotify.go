package main

import (
    "github.com/fsnotify/fsnotify"
    "log"
)

func main() {
    //创建一个监控对象
    watch, err := fsnotify.NewWatcher()
    if err != nil {
        log.Fatal(err)
    }
    defer watch.Close()
    //添加要监控的对象，文件或文件夹
    err = watch.Add("/home/john/Documents/temp")
    if err != nil {
        log.Fatal(err)
    }
    //我们另启一个goroutine来处理监控对象的事件
    go func() {
        for {
            select {
                case ev := <-watch.Events:
                    //判断事件发生的类型，如下5种
                    // Create 创建
                    // Write 写入
                    // Remove 删除
                    // Rename 重命名
                    // Chmod 修改权限
                    if ev.Op&fsnotify.Create == fsnotify.Create {
                        log.Println("创建文件 : ", ev.Name)
                    }
                    if ev.Op&fsnotify.Write == fsnotify.Write {
                        log.Println("写入文件 : ", ev.Name)
                    }
                    if ev.Op&fsnotify.Remove == fsnotify.Remove {
                        log.Println("删除文件 : ", ev.Name)
                    }
                    if ev.Op&fsnotify.Rename == fsnotify.Rename {
                        log.Println("重命名文件 : ", ev.Name)
                    }
                    if ev.Op&fsnotify.Chmod == fsnotify.Chmod {
                        log.Println("修改权限 : ", ev.Name)
                    }

                case err := <-watch.Errors:
                    log.Println("error : ", err)
                    return

            }
        }
    }()

    //循环
    select {}
}