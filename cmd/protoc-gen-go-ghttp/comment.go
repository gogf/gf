package main

import (
	"fmt"
	"google.golang.org/protobuf/compiler/protogen"
	"strings"
)

func scanSvcComment(svc *protogen.Service, file *protogen.File) string {
	// Comments Leading => Comment Trailing => Special format
	comment := ""
	switch {
	case svc.Comments.Leading.String() != "":
		comment = svc.Comments.Leading.String()
	case svc.Comments.Trailing.String() != "":
		comment = svc.Comments.Trailing.String()
	default:
		comment = string(svc.Desc.Name())
	}
	return strings.TrimSpace(comment)
}

func scanMethodComment(method *protogen.Method, file *protogen.File) string {
	// Comments Leading => Comment Trailing => Special format
	comment := ""
	switch {
	case method.Comments.Leading.String() != "":
		comment = method.Comments.Leading.String()
	case method.Comments.Trailing.String() != "":
		comment = method.Comments.Trailing.String()
	}
	return string(method.Desc.Name()) + " " + strings.TrimSpace(comment)
}

func scanMessageComment(message *protogen.Message) string {
	// Comments Leading => Comment Trailing => Special format
	comment := ""
	switch {
	case message.Comments.Leading.String() != "":
		comment = message.Comments.Leading.String()
	case message.Comments.Trailing.String() != "":
		comment = message.Comments.Trailing.String()
	}
	return string(message.Desc.Name()) + " " + strings.TrimSpace(comment)
}

func processFieldComment(message *protogen.Field) string {
	// 集中一起处理tag+comment
	// Comments Leading => Comment Trailing => Special format
	comment := message.Comments.Leading.String()
	if comment == "" && message.Comments.Trailing.String() != "" {
		return message.Comments.Trailing.String()
	}

	gfComment := ""
	fieldComment := ""
	// 判断一下comment当中是否包含规则
	for _, item := range strings.Split(comment, "\n") {
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
		gfComment = comment
	}
	return gfComment
}
