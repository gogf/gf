package main

import (
    "g/os/gfile"
    "fmt"
)

var dirpath1  = "/home/john/Workspace/temp/"
var dirpath2  = "/home/john/Workspace/temp/1"
var filepath1 = "/home/john/Workspace/temp/test.php"
var filepath2 = "/tmp/tmp.test"

func info () {
    fmt.Println(gfile.Info(dirpath1))
}

func scanDir() {
    files := gfile.ScanDir(dirpath1)
    fmt.Println(files)
}

func getContents() {
    fmt.Printf("%s\n", gfile.GetContents(filepath1))
}

func putContents() {
    fmt.Println(gfile.PutContentsAppend(filepath2, "123"))
}

func main() {
    //info()
    //getContents()
    putContents()
    //scanDir()
}