// Per: https://github.com/clbanning/mxj/issues/24
// Per: https://github.com/clbanning/mxj/issues/25

package main

import (
	"fmt"
	"io"
	"os"
	"github.com/clbanning/mxj/x2j"
)

func main() {
	for {
		_, _, err := x2j.XmlReaderToJsonWriter(os.Stdin, os.Stdout)
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println(err)
			break
		}
	}
}
