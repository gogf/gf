package main

import (
	"fmt"

	"github.com/gogf/gf/v2/errors/gerror"
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
	fmt.Println(gerror.Cause(ReadConfig()))
}
