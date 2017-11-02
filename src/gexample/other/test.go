package main

import (
    "fmt"
    "g/os/gfile"
    "os"
)

func main() {
    slice := []int{1,2,3,4,5,6,7,8,9}
    index := 1
    //fmt.Println(append(slice[:index], slice[index+1:]...))

    //rear:=append([]int{}, slice[index:]...)
    //slice=append(slice[0:index], 88)
    //slice=append(slice, rear...)
    //
    //fmt.Println(slice)

    fmt.Println(append(append(slice[0 : index], 88), append([]int{}, slice[index : ]...)...))
    return
    //a := gbinary.EncodeBits(nil, 100, 10)
    //fmt.Println(a)
    //b := gbinary.EncodeBitsToBytes(a)
    //fmt.Println(b)
    //fmt.Println(gbinary.EncodeInt32(1))
    //return
    //return

    //fmt.Println(gbinary.DecodeToInt64([]byte{1}))
    //return
    //fmt.Println(gbinary.EncodeInt32(1)[0:3])
    //b := []int{1,2,3}
    //c := []int{4}
    //copy(b[1:], c)
    //fmt.Println(b)
    //return
    //space, err := gfilespace.New("/tmp/test")
    //if err != nil {
    //    fmt.Println(err)
    //}
    //for i := 0; i < 10; i++ {
    //    space.AddBlock(int64(i), uint32((i + 1)*10))
    //}
    //fmt.Println(space.GetBlock(50))
    //return


    //db.Set([]byte("1"), []byte(grand.RandStr(10)))
    //grand.RandStr(10)
    //db.Set([]byte("r88U89b6Vv"), []byte("john211111111111111111111111"))
    //db.Get([]byte("name2"))
    //fmt.Println(e)

    //fmt.Println(string(v))
    //r := int32(binary.LittleEndian.Uint32(b))
    //fmt.Println(int32(r))
    //binary.BigEndian.Uint16(b)
    //gbinary.DecodeToInt32([]byte{1,2,3,4})
    //fmt.Println(gtime.Microsecond() - t1)
    //fmt.Println([]byte{byte(i)})

    //b := make([]byte, 0)
    //a := ghash.BKDRHash([]byte("john"))
    //for i := 0; i < 1000; i++ {
    //    r, e := gbinary.Encode([]byte("key_" + strconv.Itoa(i)), a, a)
    //    if e != nil {
    //        fmt.Println(e)
    //        return
    //    }
    //    b = append(b, r...)
    //}
    //fmt.Printf("length:     %d\n", len(b)/1024)
    //fmt.Printf("compressed: %d\n", len(gcompress.Zlib(b))/1024)
    //t1 := gtime.Microsecond()
    ////gcompress.Zlib(b)
    //gbinary.Encode([]byte("key_" + strconv.Itoa(100)), a, a)
    //fmt.Println(gtime.Microsecond() - t1)
    //return
    //t1 := gtime.Second()
    //m := make(map[uint64]bool)
    //c := 0
    //for i := 0; i < 100000000; i++ {
    //    key := ghash.SDBMHash64([]byte("this is test key" + strconv.Itoa(i)))
    //    if _, ok := m[key]; ok {
    //        c++
    //    } else {
    //        m[key] = true
    //    }
    //}
    //fmt.Println(gtime.Second() - t1)
    //fmt.Println("conflicts:", c)
    //fmt.Println(ghash.BKDRHash([]byte("johnWRWEREWREWRWEREWRRRRRRRRRRRRRRRRRRRRRRRRRRRRR")))
    //t1 := gtime.Microsecond()
    //ghash.BKDRHash64([]byte("john"))
    ////fmt.Println(ghash.ELFHash([]byte("john")))
    ////fmt.Println(ghash.JSHash([]byte("john")))
    //fmt.Println(gtime.Microsecond() - t1)
    //fmt.Println(ghash.BKDRHash64([]byte("john29384723894723894789sdkjfhsjkdh")))
    //return
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

    //return
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


    //return
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
    //
    file, _ := gfile.OpenWithFlag("/tmp/test", os.O_RDWR|os.O_CREATE)
    fmt.Println(gfile.GetBinContentByTwoOffsets(file, 100, 110))
    //n, err := file.WriteAt([]byte("123"), 2286 445 522*8)
    //n, err := file.WriteAt([]byte("123"), 1000000*(16))
    //fmt.Println(n)
    //fmt.Println(err)
    defer file.Close()
    //fmt.Println(gcrc32.EncodeString("123"))
}