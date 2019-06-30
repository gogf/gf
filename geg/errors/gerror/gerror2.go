package main

import (
	"github.com/gogf/gf/g/os/glog"

	"github.com/gogf/gf/g/errors/gerror"
)

func OpenFile() error {
	return gerror.New("permission denied")
}

func OpenConfig() error {
	return gerror.Wrap(OpenFile(), "configuration file opening failed")
}

func ReadConfig() error {
	return gerror.Wrap(OpenConfig(), "reading configuration failed")
}

func main() {
	//glog.Println(OpenConfig())
	glog.Printf("unexpected error: %+s", ReadConfig())
	//glog.Errorf("unexpected error: %+s", OpenConfig())

}
