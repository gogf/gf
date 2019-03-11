package gbase64

import (
	"testing"
)

func TestEncodeFile(t *testing.T) {
	//测试文件路径
	filepath := []map[string]string{
		{"file": "./xx.png", "retult": "nil"},
		{"file": "./xx2.png", "retult": "error"},
		{"file": "./xx.exe", "retult": "error"},
	}

	rflag := ""

	for _, v := range filepath {
		if _, flag := EncodeFile(v["file"]); flag != nil {
			rflag = "error"
		} else {
			rflag = "nil"
		}

		if rflag != v["retult"] {
			t.Errorf(v["file"] + "=>编码失败")
		}

	}

}
