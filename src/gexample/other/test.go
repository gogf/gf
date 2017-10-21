package main

import (
    "g/os/gfile"
    "os"
    "fmt"
    "g/encoding/ghash"
    "g/util/gtime"
    "strconv"
)



func main() {
    t1 := gtime.Second()
    m := make(map[uint64]bool)
    c := 0
    for i := 0; i < 10000000; i++ {
        key := ghash.SDBMHash64([]byte("this is test key" + strconv.Itoa(i)))
        if _, ok := m[key]; ok {
            c++
        } else {
            m[key] = true
        }
    }
    fmt.Println(gtime.Second() - t1)
    fmt.Println("conflicts:", c)
    //fmt.Println(ghash.BKDRHash([]byte("johnWRWEREWREWRWEREWRRRRRRRRRRRRRRRRRRRRRRRRRRRRR")))
    //t1 := gtime.Microsecond()
    //fmt.Println(ghash.BKDRHash([]byte("john")))
    //fmt.Println(ghash.ELFHash([]byte("john")))
    //fmt.Println(ghash.JSHash([]byte("john")))
    //fmt.Println(gtime.Microsecond() - t1)
    return
    //btree := gbtree.New(3)
    //t1 := gtime.Microsecond()
    //for i := 1; i <= 11; i++ {
    //    btree.Set([]byte{byte(i)}, []byte{byte(i)})
    //}
    //fmt.Println(gtime.Microsecond() - t1)
    //btree.Print()
    //fmt.Println()
    //fmt.Println()
    //btree.Remove([]byte{11})
    //btree.Print()

    //t2 := gtime.Microsecond()
    //btree.Get([]byte("key2"))
    //fmt.Println(btree.Get([]byte{200}))
    //fmt.Println(gtime.Microsecond() - t2)

    return
    ////m := gmap.NewStringInterfaceMap()
    //t1 := gtime.Microsecond()
    //gcrc32.EncodeString("123")
    //fmt.Println(gtime.Microsecond() - t1)
    //return
    //db, err := leveldb.OpenFile("/tmp/lv.db", nil)
    //defer db.Close()
    //t1 := gtime.Microsecond()
    //err = db.Put([]byte("key"), []byte("value"), nil)
    //fmt.Println(gtime.Microsecond() - t1)
    //
    //t2 := gtime.Microsecond()
    //
    //fmt.Println(db.Get([]byte("key"), nil))
    //fmt.Println(gtime.Microsecond() - t2)
    //
    //return
    //db, err := bolt.Open("/tmp/my.db", 0600, nil)
    //if err != nil {
    //    log.Fatal(err)
    //}
    //defer db.Close()
    //
    //tx, err := db.Begin(true)
    //if err != nil {
    //    log.Fatal(err)
    //}
    //defer tx.Rollback()

    // Use the transaction...
    //_, err = tx.CreateBucket([]byte("MyBucket"))
    //if err != nil {
    //    log.Fatal(err)
    //}

    // Commit the transaction and check for error.
    //if err := tx.Commit(); err != nil {
    //    log.Fatal(err)
    //}
    //t1 := gtime.Microsecond()
    //db.Update(func(tx *bolt.Tx) error {
    //    b := tx.Bucket([]byte("MyBucket"))
    //    err := b.Put([]byte("answer"), []byte("11"))
    //    return err
    //})
    //fmt.Println(gtime.Microsecond() - t1)
    //
    //t2 := gtime.Microsecond()
    //db.View(func(tx *bolt.Tx) error {
    //    b := tx.Bucket([]byte("MyBucket"))
    //    v := b.Get([]byte("answer"))
    //    fmt.Printf("The answer is: %s\n", v)
    //    return nil
    //})
    //fmt.Println(gtime.Microsecond() - t2)


    return
    ////db, err := gkvdb.New("/tmp/test2", "t")
    //fmt.Println(err)
    ////fmt.Println(db.Set("1", []byte("1")))
    //t1 := gtime.Microsecond()
    ////fmt.Println(db.Get("1"))
    //fmt.Println(db.Set("1", []byte("1")))
    //fmt.Println(gtime.Microsecond() - t1)
    ////fmt.Println(db.Set("name", []byte("222")))
    //return
    //for i := 0; i < 10000000; i++ {
    //    gfile.PutContentsAppend("/tmp/test", "1234567890")
    //}

    file, _ := gfile.OpenWithFlag("/tmp/test", os.O_WRONLY|os.O_CREATE)
    fmt.Println(string(gfile.GetBinContentByTwoOffsets(file, 10000000, 10000010)))
    //n, err := file.WriteAt([]byte("123"), 2286 445 522*8)
    n, err := file.WriteAt([]byte("123"), 42864454*(4 + 8 + 8))
    fmt.Println(n)
    fmt.Println(err)
    defer file.Close()
    //fmt.Println(gcrc32.EncodeString("123"))
}