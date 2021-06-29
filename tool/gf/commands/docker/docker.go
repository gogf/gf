package docker

import (
	"fmt"
	"github.com/gogf/gf/container/garray"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gcmd"
	"github.com/gogf/gf/os/gfile"
	"github.com/gogf/gf/os/gproc"
	"github.com/gogf/gf/text/gstr"
	"github.com/gogf/gf/tool/gf/library/mlog"
	"os"
	"strings"
)

func Help() {
	mlog.Print(gstr.TrimLeft(`
USAGE    
    gf docker [FILE] [OPTION]

ARGUMENT
    FILE      file path for "gf build", it's "main.go" in default.
    OPTION    the same options as "docker build" except some options as follows defined

OPTION
    -p, --push  auto push the docker image to docker registry if "-t" option passed

EXAMPLES
    gf docker 
    gf docker -t hub.docker.com/john/image:tag
    gf docker -p -t hub.docker.com/john/image:tag
    gf docker main.go
    gf docker main.go -t hub.docker.com/john/image:tag
    gf docker main.go -t hub.docker.com/john/image:tag
    gf docker main.go -p -t hub.docker.com/john/image:tag

DESCRIPTION
    The "docker" command builds the GF project to a docker images.
    It runs "gf build" firstly to compile the project to binary file.
    It then runs "docker build" command automatically to generate the docker image.
    You should have docker installed, and there must be a Dockerfile in the root of the project.

`))
}

func Run() {
	var err error
	autoPush := false
	array := garray.NewStrArrayFromCopy(os.Args)
	index := array.Search("--push")
	if index < 0 {
		index = array.Search("-p")
	}
	if index != -1 {
		array.Remove(index)
		autoPush = true
	}
	file := "main.go"
	extraOptions := ""
	if array.Len() > 2 {
		v, _ := array.Get(2)
		if gfile.ExtName(v) == "go" {
			file, _ = array.Get(2)
			if array.Len() > 3 {
				extraOptions = strings.Join(array.SubSlice(3), " ")
			}
		} else {
			extraOptions = strings.Join(array.SubSlice(2), " ")
		}
	}
	// Binary build.
	err = gproc.ShellRun(fmt.Sprintf(`gf build %s -a amd64 -s linux`, file))
	if err != nil {
		return
	}
	// Docker build.
	err = gproc.ShellRun(fmt.Sprintf(`docker build . %s`, extraOptions))
	if err != nil {
		return
	}
	// Docker push.
	if !autoPush {
		return
	}
	parser, err := gcmd.Parse(g.MapStrBool{
		"t,tag": true,
	})
	if err != nil {
		mlog.Fatal(err)
	}
	tag := parser.GetOpt("t")
	if tag == "" {
		return
	}
	err = gproc.ShellRun(fmt.Sprintf(`docker push %s`, tag))
	if err != nil {
		return
	}
}
