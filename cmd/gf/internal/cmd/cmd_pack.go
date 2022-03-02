package cmd

import (
	"context"
	"strings"

	"github.com/gogf/gf/cmd/gf/v2/internal/utility/allyes"
	"github.com/gogf/gf/cmd/gf/v2/internal/utility/mlog"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcmd"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/os/gres"
	"github.com/gogf/gf/v2/util/gtag"
)

var (
	Pack = cPack{}
)

type cPack struct {
	g.Meta `name:"pack" usage:"{cPackUsage}" brief:"{cPackBrief}" eg:"{cPackEg}"`
}

const (
	cPackUsage = `gf pack SRC DST`
	cPackBrief = `packing any file/directory to a resource file, or a go file`
	cPackEg    = `
gf pack public data.bin
gf pack public,template data.bin
gf pack public,template packed/data.go
gf pack public,template,config packed/data.go
gf pack public,template,config packed/data.go -n=packed -p=/var/www/my-app
gf pack /var/www/public packed/data.go -n=packed
`
	cPackSrcBrief = `source path for packing, which can be multiple source paths.`
	cPackDstBrief = `
destination file path for packed file. if extension of the filename is ".go" and "-n" option is given, 
it enables packing SRC to go file, or else it packs SRC into a binary file.
`
	cPackNameBrief   = `package name for output go file, it's set as its directory name if no name passed`
	cPackPrefixBrief = `prefix for each file packed into the resource file`
)

func init() {
	gtag.Sets(g.MapStrStr{
		`cPackUsage`:       cPackUsage,
		`cPackBrief`:       cPackBrief,
		`cPackEg`:          cPackEg,
		`cPackSrcBrief`:    cPackSrcBrief,
		`cPackDstBrief`:    cPackDstBrief,
		`cPackNameBrief`:   cPackNameBrief,
		`cPackPrefixBrief`: cPackPrefixBrief,
	})
}

type cPackInput struct {
	g.Meta `name:"pack"`
	Src    string `name:"SRC" arg:"true" v:"required" brief:"{cPackSrcBrief}"`
	Dst    string `name:"DST" arg:"true" v:"required" brief:"{cPackDstBrief}"`
	Name   string `name:"name"   short:"n" brief:"{cPackNameBrief}"`
	Prefix string `name:"prefix" short:"p" brief:"{cPackPrefixBrief}"`
}
type cPackOutput struct{}

func (c cPack) Index(ctx context.Context, in cPackInput) (out *cPackOutput, err error) {
	if gfile.Exists(in.Dst) && gfile.IsDir(in.Dst) {
		mlog.Fatalf("DST path '%s' cannot be a directory", in.Dst)
	}
	if !gfile.IsEmpty(in.Dst) && !allyes.Check() {
		s := gcmd.Scanf("path '%s' is not empty, files might be overwrote, continue? [y/n]: ", in.Dst)
		if strings.EqualFold(s, "n") {
			return
		}
	}
	if in.Name == "" && gfile.ExtName(in.Dst) == "go" {
		in.Name = gfile.Basename(gfile.Dir(in.Dst))
	}
	if in.Name != "" {
		if err = gres.PackToGoFile(in.Src, in.Dst, in.Name, in.Prefix); err != nil {
			mlog.Fatalf("pack failed: %v", err)
		}
	} else {
		if err = gres.PackToFile(in.Src, in.Dst, in.Prefix); err != nil {
			mlog.Fatalf("pack failed: %v", err)
		}
	}
	mlog.Print("done!")
	return
}
