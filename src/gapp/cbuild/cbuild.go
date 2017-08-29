package main

import (
    "strings"
    "fmt"
    "regexp"
    "g/os/gconsole"
    "os/exec"
    "sync"
)

const platforms = `
    android   arm
    darwin    386
    darwin    amd64
    darwin    arm
    darwin    arm64
    dragonfly amd64
    freebsd   386
    freebsd   amd64
    freebsd   arm
    linux     386
    linux     amd64
    linux     arm
    linux     arm64
    linux     ppc64
    linux     ppc64le
    linux     mips
    linux     mipsle
    linux     mips64
    linux     mips64le
    netbsd    386
    netbsd    amd64
    netbsd    arm
    openbsd   386
    openbsd   amd64
    openbsd   arm
    plan9     386
    plan9     amd64
    solaris   amd64
    windows   386
    windows   amd64
`
func main() {
    var wg sync.WaitGroup
    app   := gconsole.Value.Get(1)
    param := gconsole.Value.Get(2)
    reg   := regexp.MustCompile(`\s+`)
    lines := strings.Split(strings.TrimSpace(platforms), "\n")
    fmt.Println("building...")
    for _, line := range lines {
        line   = strings.TrimSpace(line)
        line   = reg.ReplaceAllString(line, " ")
        array := strings.Split(line, " ")
        os    := array[0]
        arch  := array[1]
        name  := os + "_" + arch
        if os == "windows" {
            name += ".exe"
        }
        cmd   := fmt.Sprintf("CGO_ENABLED=0 GOOS=%s GOARCH=%s go build -o ./bin/%s/%s %s", os, arch, app, name, param)
        wg.Add(1)
        go func(cmd string) {
            defer wg.Done()
            _, err := exec.Command("sh", "-c", cmd).Output()
            if err != nil {
                fmt.Println("build failed:", cmd)
                fmt.Println("build failed:", err)
                return
            }
        }(cmd)
    }
    wg.Wait()

}
