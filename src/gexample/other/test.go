package main

import (
    "fmt"
    "g/util/gtime"
    "g/core/types/gmap"
    "g/os/gfile"
    "g/util/grand"
    "sync/atomic"
)

type ST struct {
    I int64
}

//func check(id int64) bool {
//    path      := "/home/john/Workspace/Go/gluster/bin/gluster_0.8/gluster.db/gluster.entry.1.db"
//    file, err := gfile.OpenWithFlag(path, os.O_RDONLY)
//    if err == nil {
//        defer file.Close()
//        buffer := bufio.NewReader(file)
//        for {
//            line, _, err := buffer.ReadLine()
//            if err == nil {
//                var entry gluster.LogEntry
//                if json.Unmarshal(line, &entry) == nil {
//                    if entry.Id == id {
//                        return true
//                    } else if entry.Id > id {
//                        return false
//                    }
//                }
//            } else {
//                break;
//            }
//        }
//    }
//    return false
//}
//
//func check2(id int64) bool {
//    path      := "/home/john/Workspace/Go/gluster/bin/gluster_0.8/gluster.db/gluster.entry.1.db"
//    content   := gfile.GetBinContents(path)
//    slices    := bytes.SplitN(content, []byte("\n"), -1)
//    for _, line := range slices {
//        var entry gluster.LogEntry
//        if json.Unmarshal(line, &entry) == nil {
//            if entry.Id == id {
//                return true
//            } else if entry.Id > id {
//                return false
//            }
//        }
//    }
//    return false
//}
//
//func getLogEntryListFromFileById(start int64, checkid int64, max int) []gluster.LogEntry {
//    id    := start
//    match := false
//    array := make([]gluster.LogEntry, 0)
//    for {
//        path      := "/home/john/Workspace/Go/gluster/bin/gluster_0.8/gluster.db/gluster.entry.1.db"
//        file, err := gfile.OpenWithFlag(path, os.O_RDONLY)
//        if err == nil {
//            defer file.Close()
//            buffer := bufio.NewReader(file)
//            for {
//                if len(array) == max {
//                    return array
//                }
//                line, _, err := buffer.ReadLine()
//                if err == nil {
//                    var entry gluster.LogEntry
//                    if err := json.Unmarshal(line, &entry); err == nil {
//                        if entry.Id == checkid {
//                            match = true
//                        } else if entry.Id > checkid {
//                            if match {
//                                array = append(array, entry)
//                            } else {
//                                break;
//                            }
//                        }
//                    } else {
//                        return array
//                    }
//                } else {
//                    return array
//                }
//            }
//        } else {
//            break;
//        }
//        // 下一批次
//        id += 100000
//    }
//    return array
//}

type T1 struct {
    m *gmap.StringInterfaceMap
}



func main() {
    var i int64 = 1
    atomic.SwapInt64(&i, 2)
    fmt.Println(atomic.LoadInt64(&i))
    return
    //j := gjson.DecodeToJson(gfile.GetContents("/home/john/Workspace/Go/gluster/src/gluster/gluster_server.json"))
    //fmt.Println(j.GetBool("Scan2"))
    //return
    //a := []int{1,2,3}

    //b := []int{4,5,6}
    //a = append(a, b...)
    //fmt.Println(a)
    //return
    //start1 := gtime.Millisecond()
    //fmt.Println(check(999973173))
    //fmt.Println(check(999893892))
    //fmt.Println(gtime.Millisecond() - start1)

//    start2 := gtime.Millisecond()
//    fmt.Println(check2(999973173))
//    fmt.Println(check2(999893892))
//    fmt.Println(gtime.Millisecond() - start2)
//    return
//    path := "/home/john/temp/index.html"
//    //gfile.PutBinContents(path, gcompress.Zlib(gfile.GetBinContents(path)))
//    file, _ := gfile.Open(path)
//    fmt.Println(gfile.GetNextCharOffset(file, "\n", 0))
//return
//    var wg sync.WaitGroup

    // {"DataMap":{"name2":"john2"},"LastLogId":45766}
    //m := make(map[string]string)
    //for i := 0; i<1000000; i++ {
    //    key   := fmt.Sprintf("key_%d", i)
    //    value := fmt.Sprintf("value_%d", i)
    //    m[key] = value
    //}
    //path    := "/home/john/temp/gluster.data.db"
    //content := map[string]interface{}{
    //    "DataMap"   : m,
    //    "LastLogId" : 9999991721,
    //}
    //gfile.PutBinContents(path, gcompress.Zlib([]byte(gjson.Encode(content))))

return
    start   := gtime.Second()

    for n := 0; n < 10; n++ {
        content := ""
        path    := fmt.Sprintf("/home/john/temp/gluster.db/gluster.entry.%d.db", n)
        for i := n*100000; i < (n+1)*100000; i++ {
            id      := i*10000+grand.Rand(0, 9999)
            content += fmt.Sprintf("{\"Id\":%d,\"Act\":300,\"Items\":{\"key_%d\":\"value_%d\"}}\n", id, i, i)
        }
        gfile.PutContents(path, content)
        fmt.Println("done:", n)
    }

    //for i := 0; i< 500; i++ {
    //    wg.Add(1)
    //    go func() {
    //        for i := 0; i< 100; i ++ {
    //            r := ghttp.Post("http://127.0.0.1:4168/kv", fmt.Sprintf("{\"key_%d_1\":\"value_%d\"}", i, i))
    //            r.Close()
    //        }
    //        wg.Done()
    //    }()
    //}
    //wg.Wait()

    fmt.Println(gtime.Second() - start)

}