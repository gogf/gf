package main

import (
    "fmt"
    "strings"
)

func main() {
    s1 := `C:\Documents and Settings\Claymore\桌面\gf.test`
    s2 := `C:\Documents and Settings\Claymore\桌面\gf.tes`




    fmt.Println(len(s2) >= len(s1) && strings.EqualFold(s2[0 : len(s1)], s1) )
}