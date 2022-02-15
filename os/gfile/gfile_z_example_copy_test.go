// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gfile_test

import (
	"fmt"

	"github.com/gogf/gf/v2/os/gfile"
)

func ExampleCopy() {
	// init
	var (
		srcFileName = "gflie_example.txt"
		srcTempDir  = gfile.Temp("gfile_example_copy_src")
		srcTempFile = gfile.Join(srcTempDir, srcFileName)

		// copy file
		dstFileName = "gflie_example_copy.txt"
		dstTempFile = gfile.Join(srcTempDir, dstFileName)

		// copy dir
		dstTempDir = gfile.Temp("gfile_example_copy_dst")
	)

	// write contents
	gfile.PutContents(srcTempFile, "goframe example copy")

	// copy file
	gfile.Copy(srcTempFile, dstTempFile)

	// read contents after copy file
	fmt.Println(gfile.GetContents(dstTempFile))

	// copy dir
	gfile.Copy(srcTempDir, dstTempDir)

	// list copy dir file
	fList, _ := gfile.ScanDir(dstTempDir, "*", false)
	for _, v := range fList {
		fmt.Println(gfile.Basename(v))
	}

	// Output:
	// goframe example copy
	// gflie_example.txt
	// gflie_example_copy.txt
}
