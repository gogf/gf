package main

import (
	"fmt"

	"github.com/jin502437344/gf/errors/gerror"
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
	//err := ReadConfig()
	//glog.Printf("%s\n%+s", err, err)
	//glog.Printf("%+v", err)
	fmt.Printf("%+v", ReadConfig())
}
