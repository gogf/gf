package main

import (
    "g/database/gkvdb"
    "fmt"
    "g/util/gtime"
    "strconv"
    "time"
)

func main() {
    t1 := gtime.Microsecond()
    db, err := gkvdb.New("/tmp/test.db", "my")
    if err != nil {
        fmt.Println(err)
    }
    fmt.Println(gtime.Microsecond() - t1)


    //binary.LittleEndian.Uint64(bytes)
    //b, _ := gbinary.Encode(i)
    t2 := gtime.Microsecond()
    //db.Set([]byte{byte(1)}, []byte("1"))
    //db.Set([]byte{byte(1)}, []byte("123456"))
    //db.Set([]byte{byte(1)}, []byte("1"))
    //db.Set([]byte{byte(1)}, []byte("1234567890"))
    //fmt.Println(db.Get([]byte{byte(1)}))
    //fmt.Println(db.Get([]byte{byte(1)}))
    //fmt.Println(db.Set([]byte("name"), []byte("john")))
    //fmt.Println(db.Set([]byte("name2"), []byte("john2")))
    //fmt.Println(db.Get([]byte("name")))
    //fmt.Println(db.Get([]byte("name2")))
    size := 10000000
    gtime.SetInterval(5*time.Second, func() bool {
        db.PrintState()
        return true
    })
    ////////db.Set([]byte{byte(2)}, []byte{byte(2)})
    ////////db.Set([]byte{byte(1)}, []byte{byte(1)})
    ////////db.Set([]byte{byte(0)}, []byte{byte(0)})
    //////
    for i := 0; i < size; i++ {
        //r := []byte(grand.RandStr(10))
        //if err := db.Set(r, r); err != nil {
        if err := db.Set([]byte("key1_" + strconv.Itoa(i)), []byte("value1_" + strconv.Itoa(i))); err != nil {
            //if err := db.Set(gbinary.EncodeInt32(int32(i)), gbinary.EncodeInt32(int32(i))); err != nil {
            fmt.Println(err)
        }
    }
    //for i := 0; i < size; i++ {
    //    r := db.Get([]byte("key1_" + strconv.Itoa(i)))
    //    //r := db.Get(gbinary.EncodeInt32(int32(i)))
    //    if r == nil {
    //        fmt.Println("none for ", i)
    //    }
    //}
    //db.Remove(true)
    db.PrintState()
    //db.Get([]byte("key1_" + strconv.Itoa(99999)))
    //fmt.Println(string(db.Get([]byte("key1_" + strconv.Itoa(99999)))))
    //fmt.Println(gbinary.DecodeToInt32(db.Get(gbinary.EncodeInt32(4253318))))
    fmt.Println(gtime.Microsecond() - t2)
}