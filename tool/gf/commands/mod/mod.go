package mod

import (
	"fmt"
	"github.com/gogf/gf/container/gmap"
	"github.com/gogf/gf/os/gcmd"
	"github.com/gogf/gf/os/genv"
	"github.com/gogf/gf/os/gfile"
	"github.com/gogf/gf/text/gstr"
	"github.com/gogf/gf/tool/gf/library/mlog"
)

func Help() {
	mlog.Print(gstr.TrimLeft(`
USAGE 
    gf mod ARGUMENT

ARGUMENT
    path  copy all packages with its latest version in Go modules, which does not exist 
          in GOPATH, to GOPATH. This enables your project using GOPATH building, but you 
          should have GOPATH environment variable configured.

EXAMPLES
    gf mod path
`))
}

func Run() {
	argument := gcmd.GetArg(2)
	switch argument {
	case "path":
		doPath()

	default:
		mlog.Print("argument cannot be empty")
		Help()
	}
}

// doPath copies all packages in Go modules, which does not exist in GOPATH, to GOPATH.
// This enables your project using GOPATH building, but you should have GOPATH
// environment variable configured.
func doPath() {
	goPathEnv := genv.Get("GOPATH")
	if goPathEnv == "" {
		mlog.Fatal("GOPATH is not found in your environment")
	}
	mlog.Print("scanning...")
	var (
		copied    = false
		haveCount = 0
	)
	for _, goPath := range gstr.SplitAndTrim(goPathEnv, ";") {
		goModPath := gfile.Join(goPath, "pkg", "mod")
		if !gfile.Exists(goModPath) {
			continue
		}
		pathMap := gmap.NewStrStrMap()
		_, err := gfile.ScanDirFunc(goModPath, "*.*", true, func(path string) string {
			// Ignore the cache folder.
			if gstr.Contains(path, gfile.Join(goModPath, "cache")) {
				return ""
			}
			name := gfile.Name(path)
			if name == "" {
				return ""
			}
			if !gstr.Contains(name, "@") {
				return ""
			}
			if n := gstr.Count(path, "@"); n > 1 {
				return ""
			}
			if !gfile.IsDir(path) {
				return ""
			}
			array := gstr.Split(path, "@")
			if v := pathMap.Get(array[0]); v == "" {
				pathMap.Set(array[0], array[1])
			} else {
				if gstr.CompareVersionGo(v, array[1]) < 0 {
					pathMap.Set(array[0], array[1])
				}
			}
			return path
		})
		if err != nil {
			mlog.Fatal(err)
		}
		haveCount += pathMap.Size()
		pathMap.Iterator(func(k string, v string) bool {
			src := fmt.Sprintf(`%s@%s`, k, v)
			dst := gfile.Join(goPath, "src", gstr.Trim(gstr.Replace(k, goModPath, ""), "\\/"))
			if !gfile.Exists(dst) {
				mlog.Printf(`copying %s to %s`, src, dst)
				if err := gfile.Copy(src, dst); err != nil {
					mlog.Fatal(err)
				}
				copied = true
			}
			return true
		})
	}
	if !copied {
		if haveCount > 0 {
			mlog.Print(`all packages of go modules already exist in GOPATH`)
		} else {
			mlog.Printf(`no packages found in go module path: %s`, goPathEnv)
		}
		return
	}
	mlog.Print("done!")
}
