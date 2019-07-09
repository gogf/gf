package main

import (
	"fmt"

	"github.com/gogf/gf/g/text/gregex"
)

func main() {
	query := "SELECT * FROM user where status=1 LIMIT 10, 100"
	query, _ = gregex.ReplaceString(` LIMIT (\d+),\s*(\d+)`, ` LIMIT $1 OFFSET $2`, query)
	fmt.Println(query)
}
