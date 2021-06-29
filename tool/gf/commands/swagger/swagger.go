package swagger

import (
	"errors"
	"fmt"
	"github.com/gogf/gf/container/gtype"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gcmd"
	"github.com/gogf/gf/os/gfile"
	"github.com/gogf/gf/os/gfsnotify"
	"github.com/gogf/gf/os/gproc"
	"github.com/gogf/gf/os/gtime"
	"github.com/gogf/gf/os/gtimer"
	"github.com/gogf/gf/text/gstr"
	"github.com/gogf/gf/tool/gf/library/mlog"
	"github.com/gogf/swagger"
)

const (
	defaultOutput    = "./swagger"
	swaggoRepoPath   = "github.com/swaggo/swag/cmd/swag"
	PackedGoFileName = "swagger.go"
)

func Help() {
	mlog.Print(gstr.TrimLeft(`
USAGE    
    gf swagger [OPTION]

OPTION
    -s, --server  start a swagger server at specified address after swagger files
                  produced
    -o, --output  the output directory for storage parsed swagger files,
                  the default output directory is "./swagger"
    -/--pack      auto parses and packs swagger into packed/swagger.go. 

EXAMPLES
    gf swagger
    gf swagger --pack
    gf swagger -s 8080
    gf swagger -s 127.0.0.1:8080
    gf swagger -o ./document/swagger


DESCRIPTION
    The "swagger" command parses the current project and produces swagger API description 
    files, which can be used in swagger API server. If used with "-s/--server" option, it
    watches the changes of go files of current project and reproduces the swagger files,
    which is quite convenient for local API development.
    If it fails in command "swag", please firstly check your system PATH whether containing 
    go binary path, or you can install the "swag" tool manually referring to: 
    https://github.com/swaggo/swag
`))
}

func Run() {
	mlog.SetHeaderPrint(true)
	parser, err := gcmd.Parse(g.MapStrBool{
		"s,server": true,
		"o,output": true,
		"pack":     false,
	})
	if err != nil {
		mlog.Fatal(err)
	}
	server := parser.GetOpt("server")
	output := parser.GetOpt("output", defaultOutput)
	// Generate swagger files.
	if err := generateSwaggerFiles(output, parser.ContainsOpt("pack")); err != nil {
		mlog.Print(err)
	}
	// Watch the go file changes and regenerate the swagger files.
	dirty := gtype.NewBool()
	_, err = gfsnotify.Add(gfile.RealPath("."), func(event *gfsnotify.Event) {
		if gfile.ExtName(event.Path) != "go" || gstr.Contains(event.Path, "swagger") {
			return
		}
		// Variable <dirty> is used for running the changes only one in one second.
		if !dirty.Cas(false, true) {
			return
		}
		// With some delay in case of multiple code changes in very short interval.
		gtimer.SetTimeout(1500*gtime.MS, func() {
			mlog.Printf(`go file changes: %s`, event.String())
			mlog.Print(`reproducing swagger files...`)
			if err := generateSwaggerFiles(output, parser.ContainsOpt("pack")); err != nil {
				mlog.Print(err)
			} else {
				mlog.Print(`done!`)
			}
			dirty.Set(false)
		})
	})
	if err != nil {
		mlog.Fatal(err)
	}
	// Swagger server starts.
	if server != "" {
		if gstr.IsNumeric(server) {
			server = ":" + server
		}
		s := g.Server()
		s.Plugin(&swagger.Swagger{})
		s.SetAddr(server)
		s.Run()
	}
}

// generateSwaggerFiles generates necessary swagger files.
func generateSwaggerFiles(output string, pack bool) error {
	mlog.Print(`producing swagger files...`)
	// Temporary storing swagger files directory.
	tempOutputPath := gfile.Join(gfile.TempDir(), "swagger")
	if gfile.Exists(tempOutputPath) {
		gfile.Remove(tempOutputPath)
	}
	gfile.Mkdir(tempOutputPath)
	// Check and install swag tool.
	swag := gproc.SearchBinary("swag")
	if swag == "" {
		err := gproc.ShellRun(fmt.Sprintf(`go get -u -v %s`, swaggoRepoPath))
		if err != nil {
			return err
		}
	}
	// Generate swagger files using swag.
	command := fmt.Sprintf(`swag init -o %s`, tempOutputPath)
	result, err := gproc.ShellExec(command)
	if err != nil {
		return errors.New(result + err.Error())
	}
	if !gfile.Exists(gfile.Join(tempOutputPath, "swagger.json")) {
		return errors.New("make swagger files failed")
	}
	if !gfile.Exists(output) {
		gfile.Mkdir(output)
	}
	if err = gfile.CopyFile(
		gfile.Join(tempOutputPath, "swagger.json"),
		gfile.Join(output, "swagger.json"),
	); err != nil {
		return err
	}
	mlog.Print(`done!`)
	// Auto pack into go file.
	if pack && gfile.Exists("swagger") {
		packCmd := fmt.Sprintf(`gf pack %s packed/%s -n packed`, "swagger", PackedGoFileName)
		mlog.Print(packCmd)
		if err := gproc.ShellRun(packCmd); err != nil {
			return err
		}
	}
	return nil
}
