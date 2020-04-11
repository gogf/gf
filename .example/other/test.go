package main

import (
	"fmt"
	"github.com/gogf/gf/text/gregex"
	"github.com/gogf/gf/text/gstr"
)

func main() {
	s := `user u LEFT JOIN user_detail ud ON(ud.id=u.id) LEFT JOIN user_stats us ON(us.id=u.id)`
	split := " JOIN "
	as := "my-as"
	if gstr.Contains(s, split) {
		// For join table.
		array := gstr.Split(s, split)
		array[len(array)-1], _ = gregex.ReplaceString(`(.+) ON`, fmt.Sprintf(`$1 AS %s ON`, as), array[len(array)-1])
		s = gstr.Join(array, split)
	} else {
		// For base table.
		s = gstr.TrimRight(s) + " AS " + as
	}
	fmt.Println(s)
}
