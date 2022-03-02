package cmd

import (
	"context"
	"fmt"

	"github.com/gogf/gf/cmd/gf/v2/internal/utility/mlog"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/os/gproc"
	"github.com/gogf/gf/v2/util/gtag"
)

var (
	Docker = cDocker{}
)

type cDocker struct {
	g.Meta `name:"docker" usage:"{cDockerUsage}" brief:"{cDockerBrief}" eg:"{cDockerEg}" dc:"{cDockerDc}"`
}

const (
	cDockerUsage = `gf docker [MAIN] [OPTION]`
	cDockerBrief = `build docker image for current GoFrame project`
	cDockerEg    = `
gf docker 
gf docker -t hub.docker.com/john/image:tag
gf docker -p -t hub.docker.com/john/image:tag
gf docker main.go
gf docker main.go -t hub.docker.com/john/image:tag
gf docker main.go -t hub.docker.com/john/image:tag
gf docker main.go -p -t hub.docker.com/john/image:tag
`
	cDockerDc = `
The "docker" command builds the GF project to a docker images.
It runs "gf build" firstly to compile the project to binary file.
It then runs "docker build" command automatically to generate the docker image.
You should have docker installed, and there must be a Dockerfile in the root of the project.
`
	cDockerMainBrief  = `main file path for "gf build", it's "main.go" in default. empty string for no binary build`
	cDockerBuildBrief = `binary build options before docker image build, it's "-a amd64 -s linux" in default`
	cDockerFileBrief  = `file path of the Dockerfile. it's "manifest/docker/Dockerfile" in default`
	cDockerShellBrief = `path of the shell file which is executed before docker build`
	cDockerPushBrief  = `auto push the docker image to docker registry if "-t" option passed`
	cDockerTagBrief   = `tag name for this docker, which is usually used for docker push`
	cDockerExtraBrief = `extra build options passed to "docker image"`
)

func init() {
	gtag.Sets(g.MapStrStr{
		`cDockerUsage`:      cDockerUsage,
		`cDockerBrief`:      cDockerBrief,
		`cDockerEg`:         cDockerEg,
		`cDockerDc`:         cDockerDc,
		`cDockerMainBrief`:  cDockerMainBrief,
		`cDockerFileBrief`:  cDockerFileBrief,
		`cDockerShellBrief`: cDockerShellBrief,
		`cDockerBuildBrief`: cDockerBuildBrief,
		`cDockerPushBrief`:  cDockerPushBrief,
		`cDockerTagBrief`:   cDockerTagBrief,
		`cDockerExtraBrief`: cDockerExtraBrief,
	})
}

type cDockerInput struct {
	g.Meta `name:"docker" config:"gfcli.docker"`
	Main   string `name:"MAIN"  arg:"true" brief:"{cDockerMainBrief}"  d:"main.go"`
	File   string `name:"file"  short:"f"  brief:"{cDockerFileBrief}"  d:"manifest/docker/Dockerfile"`
	Shell  string `name:"shell" short:"s"  brief:"{cDockerShellBrief}" d:"manifest/docker/docker.sh"`
	Build  string `name:"build" short:"b"  brief:"{cDockerBuildBrief}" d:"-a amd64 -s linux"`
	Tag    string `name:"tag"   short:"t"  brief:"{cDockerTagBrief}"`
	Push   bool   `name:"push"  short:"p"  brief:"{cDockerPushBrief}" orphan:"true"`
	Extra  string `name:"extra" short:"e"  brief:"{cDockerExtraBrief}"`
}
type cDockerOutput struct{}

func (c cDocker) Index(ctx context.Context, in cDockerInput) (out *cDockerOutput, err error) {
	// Necessary check.
	if gproc.SearchBinary("docker") == "" {
		mlog.Fatalf(`command "docker" not found in your environment, please install docker first to proceed this command`)
	}

	// Binary build.
	in.Build += " --exit"
	if in.Main != "" {
		if err = gproc.ShellRun(fmt.Sprintf(`gf build %s %s`, in.Main, in.Build)); err != nil {
			return
		}
	}

	// Shell executing.
	if gfile.Exists(in.Shell) {
		if err = gproc.ShellRun(gfile.GetContents(in.Shell)); err != nil {
			return
		}
	}
	// Docker build.
	dockerBuildOptions := ""
	if in.Tag != "" {
		dockerBuildOptions = fmt.Sprintf(`-t %s`, in.Tag)
	}
	if in.Extra != "" {
		dockerBuildOptions = fmt.Sprintf(`%s %s`, dockerBuildOptions, in.Extra)
	}
	if err = gproc.ShellRun(fmt.Sprintf(`docker build -f %s . %s`, in.File, dockerBuildOptions)); err != nil {
		return
	}
	// Docker push.
	if in.Tag == "" || !in.Push {
		return
	}
	if err = gproc.ShellRun(fmt.Sprintf(`docker push %s`, in.Tag)); err != nil {
		return
	}
	return
}
