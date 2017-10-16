package main

import (
    "g/os/gfile"
    "os"
    "fmt"
    "g/database/gkvdb"
)



func main() {
    db, err := gkvdb.New("/tmp/test2", "t")
    fmt.Println(err)
    fmt.Println(db.Set("name", []byte("1234567890")))
    //fmt.Println(db.Set("name", []byte("1111111111")))
    //fmt.Println(db.Set("name", []byte("222")))
    return
    //for i := 0; i < 10000000; i++ {
    //    gfile.PutContentsAppend("/tmp/test", "1234567890")
    //}

    file, _ := gfile.OpenWithFlag("/tmp/test", os.O_WRONLY|os.O_CREATE)
    //fmt.Println(string(gfile.GetBinContentByTwoOffsets(file, 10000000, 10000010)))
    //n, err := file.WriteAt([]byte("123"), 2286445522*8)
    n, err := file.WriteAt([]byte("123"), 42864454*(4 + 8 + 8))
    fmt.Println(n)
    fmt.Println(err)
    defer file.Close()
    //fmt.Println(gcrc32.EncodeString("123"))
}