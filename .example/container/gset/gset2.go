package main

import (
	"fmt"

	"github.com/gogf/gf/v2/container/gset"
	"github.com/gogf/gf/v2/frame/g"
)

func main() {
	s1 := gset.NewFrom(g.Slice{1, 2, 3})
	s2 := gset.NewFrom(g.Slice{4, 5, 6})
	s3 := gset.NewFrom(g.Slice{1, 2, 3, 4, 5, 6, 7})

	// 交集
	fmt.Println(s3.Intersect(s1).Slice())
	// 差集
	fmt.Println(s3.Diff(s1).Slice())
	// 并集
	fmt.Println(s1.Union(s2).Slice())
	// 补集
	fmt.Println(s1.Complement(s3).Slice())
}
