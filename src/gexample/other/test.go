package main

import (
    "fmt"
    "os/exec"
    "strings"
)

type ST struct {
    I int64
}

func main() {
    cmd := "echo -c 1"
    parts := strings.Fields(cmd)

    r, e := exec.Command(parts[0], parts[1:]...).Output()
    fmt.Println(string(r))
    fmt.Println(e)

}