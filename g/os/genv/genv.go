// Copyright 2017 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

// 环境变量管理
package genv

import "os"

func All() []string {
    return os.Environ()
}

func Get(k string) string {
    return os.Getenv(k)
}

func Set(k, v string) error {
    return os.Setenv(k, v)
}

func Remove(k string) error {
    return os.Unsetenv(k)
}