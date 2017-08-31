package main

import (
    "strings"
    "fmt"
    "regexp"
    "g/os/gconsole"
    "os/exec"
    "sync"
    "g/util/gutil"
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
    src     := gconsole.Value.Get(1)
    name    := gconsole.Option.Get("name")
    version := gconsole.Option.Get("version")
    oses    := strings.Split(gconsole.Option.Get("os"), ",")
    arches  := strings.Split(gconsole.Option.Get("arch"), ",")
    reg     := regexp.MustCompile(`\s+`)
    lines   := strings.Split(strings.TrimSpace(platforms), "\n")
    fmt.Println("building...")
    for _, line := range lines {
        line   = strings.TrimSpace(line)
        line   = reg.ReplaceAllString(line, " ")
        array := strings.Split(line, " ")
        os    := array[0]
        arch  := array[1]
        if len(oses) > 0 && !gutil.StringInArray(oses, os) {
            continue
        }
        if len(arches) > 0 && !gutil.StringInArray(arches, arch) {
            continue
        }
        appname := name + "." + os + "_" + arch
        if os == "windows" {
            appname += ".exe"
        }
        cmd := fmt.Sprintf("CGO_ENABLED=0 GOOS=%s GOARCH=%s go build -o ./bin/%s_%s/%s %s", os, arch, name, version, appname, src)
        wg.Add(1)
        go func(cmd string) {
            defer wg.Done()
            _, err := exec.Command("sh", "-c", cmd).Output()
            if err != nil {
                fmt.Println("build failed:", cmd)
                return
            } else {
                fmt.Println(cmd)
            }
        }(cmd)
    }
    wg.Wait()

}
