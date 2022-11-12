package main

import (
	"fmt"
	"google.golang.org/protobuf/compiler/protogen"
	"strings"
)

func trimSpaceSlash(str string) string {
	str = strings.TrimSpace(str)
	str = strings.TrimPrefix(str, "//")
	return str
}

func scanSvcComment(svc *protogen.Service, file *protogen.File) string {
	// Comments Leading => Comment Trailing => Special format
	switch {
	case svc.Comments.Leading.String() != "":
		return trimSpaceSlash(svc.Comments.Leading.String())
	case svc.Comments.Trailing.String() != "":
		return trimSpaceSlash(svc.Comments.Trailing.String())
	default:
		return ""
	}
}

func scanMethodComment(method *protogen.Method) string {
	// Comments Leading => Comment Trailing => Special format
	switch {
	case method.Comments.Leading.String() != "":
		return trimSpaceSlash(method.Comments.Leading.String())
	case method.Comments.Trailing.String() != "":
		return trimSpaceSlash(method.Comments.Trailing.String())
	default:
		return ""
	}
}

func scanMessageComment(message *protogen.Message) string {
	// Comments Leading => Comment Trailing => Special format
	switch {
	case message.Comments.Leading.String() != "":
		return trimSpaceSlash(message.Comments.Leading.String())
	case message.Comments.Trailing.String() != "":
		return trimSpaceSlash(message.Comments.Trailing.String())
	default:
		return ""
	}
}

func processFieldComment(message *protogen.Field) string {
	// 集中一起处理tag+comment
	// Comments Leading => Comment Trailing => Special format
	leading := message.Comments.Leading.String()
	trailing := message.Comments.Trailing.String()

	switch {
	case leading == "" && trailing != "":
		return strings.TrimSpace(trailing)
	default:
		gfComment := ""
		fieldComment := ""
		// 判断一下comment当中是否包含规则
		for _, item := range strings.Split(leading, "\n") {
			item = strings.TrimPrefix(item, "//")
			item = strings.TrimPrefix(item, " ")
			if len(item) < 2 {
				continue
			}
			switch {
			case strings.HasPrefix(item, "d:"):
				item = strings.TrimPrefix(item, "d:")
				fallthrough
			case strings.HasPrefix(item, "default:"):
				item = strings.TrimPrefix(item, "default:")
				fallthrough

			case strings.HasPrefix(item, "eg:"):
				item = strings.TrimPrefix(item, "eg:")
				item = strings.TrimSpace(item)
				gfComment += fmt.Sprintf(`d:"%s" `, item)

			case strings.HasPrefix(item, "v:"):
				item = strings.TrimPrefix(item, "v:")
				item = strings.TrimSpace(item)
				gfComment += fmt.Sprintf(`v:"%s" `, item)

			case strings.HasPrefix(item, "p:"):
				item = strings.TrimPrefix(item, "p:")
				item = strings.TrimSpace(item)
				gfComment += fmt.Sprintf(`p:"%s" `, item)

			case strings.HasPrefix(item, "j:"):
				item = strings.TrimPrefix(item, "j:")
				fallthrough

			case strings.HasPrefix(item, "json:"):
				item = strings.TrimPrefix(item, "json:")
				item = strings.TrimSpace(item)
				gfComment += fmt.Sprintf(`json:"%s" `, item)

			default:
				fieldComment = strings.TrimSpace(item)

			}
		}
		if gfComment != "" {
			//if message.Desc.HasJSONName() {
			//	gfComment += fmt.Sprintf(`json:"%s" `, message.Desc.JSONName())
			//}
			gfComment = "`" + gfComment + "` // " + fieldComment
		}
		// 如果comment当中没有规则，那么就直接使用comment
		if gfComment == "" {
			gfComment = leading
		}
		return strings.TrimSpace(gfComment)
	}
}
