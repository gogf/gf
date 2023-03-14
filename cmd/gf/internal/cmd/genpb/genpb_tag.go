package genpb

import (
	"context"
	"fmt"

	"github.com/gogf/gf/cmd/gf/v2/internal/utility/utils"
	"github.com/gogf/gf/v2/container/gmap"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/text/gregex"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gconv"
)

type generateStructTagInput struct {
	OutputApiPath string
}

func (c CGenPb) generateStructTag(ctx context.Context, in generateStructTagInput) (err error) {
	files, err := gfile.ScanDirFile(in.OutputApiPath, "*.pb.go", true)
	if err != nil {
		return err
	}
	var content string
	for _, file := range files {
		content = gfile.GetContents(file)
		content, err = c.doTagReplacement(ctx, content)
		if err != nil {
			return err
		}
		if err = gfile.PutContents(file, content); err != nil {
			return err
		}
		utils.GoFmt(file)
	}
	return
}

func (c CGenPb) doTagReplacement(ctx context.Context, content string) (string, error) {
	content, err := gregex.ReplaceStringFuncMatch(`type (\w+) struct {([\s\S]+?)}`, content, func(match []string) string {
		var (
			topCommentMatch  []string
			tailCommentMatch []string
			lines            = gstr.Split(match[2], "\n")
			lineTagMap       = gmap.NewListMap()
		)
		for index, line := range lines {
			line = gstr.Trim(line)
			if line == "" {
				continue
			}
			// Top comment.
			topCommentMatch, _ = gregex.MatchString(`^/[/|\*](.+)`, line)
			if len(topCommentMatch) > 1 {
				c.tagCommentIntoListMap(gstr.Trim(topCommentMatch[1]), lineTagMap)
				continue
			}
			// Tail comment.
			tailCommentMatch, _ = gregex.MatchString(".+?`.+?`.+?//(.+)", line)
			if len(tailCommentMatch) > 1 {
				c.tagCommentIntoListMap(gstr.Trim(tailCommentMatch[1]), lineTagMap)
			}
			// Tag injection.
			if !lineTagMap.IsEmpty() {
				tagContent := c.listMapToStructTag(lineTagMap)
				lineTagMap.Clear()
				line, _ = gregex.ReplaceString("`(.+)`", fmt.Sprintf("`$1 %s`", tagContent), line)
			}
			lines[index] = line
		}
		match[2] = gstr.Join(lines, "\n")
		return fmt.Sprintf("type %s struct {%s}", match[1], match[2])
	})
	return content, err
}

func (c CGenPb) tagCommentIntoListMap(comment string, lineTagMap *gmap.ListMap) {
	tagCommentMatch, _ := gregex.MatchString(`^(\w+):(.+)`, comment)
	if len(tagCommentMatch) > 1 {
		var (
			tagName    = gstr.Trim(tagCommentMatch[1])
			tagContent = gstr.Trim(tagCommentMatch[2])
		)
		lineTagMap.Set(tagName, lineTagMap.GetVar(tagName).String()+tagContent)
	} else {
		var (
			tagName    = "dc"
			tagContent = comment
		)
		lineTagMap.Set(tagName, lineTagMap.GetVar(tagName).String()+tagContent)
	}
}

func (c CGenPb) listMapToStructTag(lineTagMap *gmap.ListMap) string {
	var tag string
	lineTagMap.Iterator(func(key, value interface{}) bool {
		if tag != "" {
			tag += " "
		}
		tag += fmt.Sprintf(
			`%s:"%s"`,
			key, gstr.Replace(gconv.String(value), `"`, `\"`),
		)
		return true
	})
	return tag
}
