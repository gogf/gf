package main

import (
    "fmt"
    "bytes"
    "encoding/binary"
    "g/encoding/gcrc32"
    "strconv"
)

// 二进制打包
func encode(vs ...interface{}) []byte {
    buf := new(bytes.Buffer)
    for i := 0; i < len(vs); i++ {
        binary.Write(buf, binary.LittleEndian, vs[i])
    }
    return buf.Bytes()
}

// 二进制解包
func decode(b []byte, vs ...interface{}) {
    buf := bytes.NewBuffer(b)
    for i := 0; i < len(vs); i++ {
        binary.Read(buf, binary.LittleEndian, vs[i])
    }
}

func main() {
    //hash := crc32.NewIEEE()

    check_str := "Discard medicine more than two years old."
    fmt.Println(gcrc32.EncodeString(check_str))
    fmt.Printf("%x\n", gcrc32.EncodeString(check_str))
    fmt.Println(strconv.ParseInt("6b9cdfe7", 16, 32))

    return
    //fmt.Println("\a")
    ////gfile.PutContents("/tmp/test", "123\0456\0789")
    ////fmt.Println(gfile.GetContents("/tmp/test"))
    //return
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

//    s := `74142374,300,{"key_99999":"value_99999"}`
//    reg, _ := regexp.Compile(`^(\d+),.+$`)
//    results := reg.FindStringSubmatch(s)
//    fmt.Println(results)
//    return
//    //path1 := "/home/john/temp/temp"
//    path2 := "/home/john/temp/temp2"
//    //path3 := "/home/john/temp/gluster"
//    gfile.PutBinContents(path2, []byte("123456\n"))
//    //gfile.PutBinContents(path1, gcompress.Zlib(gfile.GetBinContents(path3)))
//    //file, _ := gfile.Open(path1)
//    //fmt.Println(gfile.GetNextCharOffset(file, "\0000", 0))
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

//return
//    start   := gtime.Second()
//
//    for n := 0; n < 10; n++ {
//        content := ""
//        path    := fmt.Sprintf("/home/john/temp/gluster.db/gluster.entry.%d.db", n)
//        for i := n*100000; i < (n+1)*100000; i++ {
//            id      := i*10000+grand.Rand(0, 9999)
//            content += fmt.Sprintf("{\"Id\":%d,\"Act\":300,\"Items\":{\"key_%d\":\"value_%d\"}}\n", id, i, i)
//        }
//        gfile.PutContents(path, content)
//        fmt.Println("done:", n)
//    }

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

    //fmt.Println(gtime.Second() - start)

}